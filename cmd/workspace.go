package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"quickvm/internal/hyperv"
	"quickvm/internal/output"

	"github.com/spf13/cobra"
)

// WorkspaceListResult represents the result of listing workspaces
type WorkspaceListResult struct {
	Workspaces []string `json:"workspaces"`
	Total      int      `json:"total"`
}

// WorkspaceResult represents the result of a single workspace operation
type WorkspaceResult struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// WorkspaceShowResult represents the result of displaying workspace details
type WorkspaceShowResult struct {
	Workspace *hyperv.Workspace `json:"workspace"`
}

// WorkspaceBatchResult represents the result of starting or stopping a workspace
type WorkspaceBatchResult struct {
	Workspace    string              `json:"workspace"`
	Operation    string              `json:"operation"`
	Results      []VMOperationResult `json:"results"`
	SuccessCount int                 `json:"successCount"`
	FailCount    int                 `json:"failCount"`
	TotalCount   int                 `json:"totalCount"`
}

var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"ws"},
	Short:   "Manage VM workspaces (groups)",
	Long:    `Manage groups of virtual machines using workspace profiles.`,
}

var wsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workspaces",
	Run: func(_ *cobra.Command, _ []string) {
		names, err := hyperv.ListWorkspaces()
		if err != nil {
			output.PrintError("WORKSPACE_LIST_FAILED", "Failed to list workspaces", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to list workspaces: %v\n", err)
			}
			os.Exit(1)
		}

		if output.IsJSON() {
			output.PrintData(WorkspaceListResult{
				Workspaces: names,
				Total:      len(names),
			})
			return
		}

		if len(names) == 0 {
			fmt.Println("📭 No workspaces found.")
			fmt.Println("💡 Create one with: quickvm ws create <name> --vms \"VM1,VM2\"")
			return
		}

		fmt.Printf("📋 Available Workspaces (%d):\n", len(names))
		for _, name := range names {
			fmt.Printf("  - %s\n", name)
		}
	},
}

var wsVms string
var wsCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		name := args[0]
		var vms []string
		for _, v := range strings.Split(wsVms, ",") {
			trimmed := strings.TrimSpace(v)
			if trimmed != "" {
				vms = append(vms, trimmed)
			}
		}

		if len(vms) == 0 {
			output.PrintError("INVALID_ARGS", "No valid VMs specified", "Provide at least one non-empty VM name via --vms")
			if !output.IsJSON() {
				fmt.Println("❌ Error: No valid VMs specified. Provide at least one VM name via --vms (e.g., -v 'VM1,VM2')")
			}
			os.Exit(1)
		}

		ws := &hyperv.Workspace{
			Name:        name,
			Description: "Created via CLI",
			VMs:         vms,
		}

		if err := hyperv.SaveWorkspace(ws); err != nil {
			output.PrintError("WORKSPACE_CREATE_FAILED", "Failed to save workspace", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to save workspace: %v\n", err)
			}
			os.Exit(1)
		}

		if output.IsJSON() {
			output.PrintData(WorkspaceResult{
				Name:    name,
				Success: true,
				Message: fmt.Sprintf("Workspace '%s' created successfully", name),
			})
			return
		}

		fmt.Printf("✅ Workspace '%s' created successfully!\n", name)
	},
}

var wsShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show workspace details",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		ws, err := hyperv.LoadWorkspace(args[0])
		if err != nil {
			output.PrintError("WORKSPACE_GET_FAILED", "Failed to load workspace", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to load workspace: %v\n", err)
			}
			os.Exit(1)
		}

		if output.IsJSON() {
			output.PrintData(WorkspaceShowResult{
				Workspace: ws,
			})
			return
		}

		fmt.Printf("📂 Workspace: %s\n", ws.Name)
		fmt.Printf("📝 Description: %s\n", ws.Description)
		fmt.Println("🖥️  Virtual Machines:")
		for _, vm := range ws.VMs {
			fmt.Printf("  - %s\n", vm)
		}
	},
}

var wsDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		if err := hyperv.DeleteWorkspace(args[0]); err != nil {
			output.PrintError("WORKSPACE_DELETE_FAILED", "Failed to delete workspace", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to delete workspace: %v\n", err)
			}
			os.Exit(1)
		}

		if output.IsJSON() {
			output.PrintData(WorkspaceResult{
				Name:    args[0],
				Success: true,
				Message: fmt.Sprintf("Workspace '%s' deleted successfully", args[0]),
			})
			return
		}

		fmt.Printf("✅ Workspace '%s' deleted.\n", args[0])
	},
}

var wsStartCmd = &cobra.Command{
	Use:   "start <name>",
	Short: "Start all VMs in a workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ws, err := hyperv.LoadWorkspace(args[0])
		if err != nil {
			output.PrintError("WORKSPACE_GET_FAILED", "Failed to load workspace", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to load workspace: %v\n", err)
			}
			os.Exit(1)
		}

		manager := hyperv.NewManager()
		if !output.IsJSON() {
			fmt.Printf("🚀 Starting workspace '%s' (%d VMs)...\n", ws.Name, len(ws.VMs))
		}

		results := make([]VMOperationResult, len(ws.VMs))
		successCount := 0
		failCount := 0

		var (
			mu sync.Mutex
			g  errgroup.Group
		)
		sem := make(chan struct{}, defaultBatchConcurrency)

		for i, vmName := range ws.VMs {
			i, vmName := i, vmName
			sem <- struct{}{}
			g.Go(func() error {
				defer func() { <-sem }()

				opCtx, cancel := context.WithTimeout(cmd.Context(), 60*time.Second)
				defer cancel()

				if !output.IsJSON() {
					mu.Lock()
					fmt.Printf("🚀 Starting VM: %s...\n", vmName)
					mu.Unlock()
				}
				opErr := manager.StartVMByName(opCtx, vmName)
				res := VMOperationResult{
					Index: i + 1,
					Name:  vmName,
				}
				mu.Lock()
				defer mu.Unlock()
				if opErr != nil {
					res.Success = false
					res.Error = opErr.Error()
					failCount++
					if !output.IsJSON() {
						fmt.Printf("❌ Failed to start VM '%s': %v\n", vmName, opErr)
					}
				} else {
					res.Success = true
					res.Message = "VM started successfully"
					successCount++
					if !output.IsJSON() {
						fmt.Printf("✅ VM '%s' started.\n", vmName)
					}
				}
				results[i] = res
				return nil
			})
		}
		_ = g.Wait()

		if results == nil {
			results = []VMOperationResult{}
		}

		if output.IsJSON() {
			topSuccess := failCount == 0
			var errInfo *output.ErrorInfo
			if failCount > 0 {
				errInfo = &output.ErrorInfo{
					Code:    "WORKSPACE_START_FAILED",
					Message: fmt.Sprintf("%d of %d VMs failed to start", failCount, len(ws.VMs)),
				}
			}
			output.PrintResponse(output.Response{
				Success: topSuccess,
				Data: WorkspaceBatchResult{
					Workspace:    ws.Name,
					Operation:    "start",
					Results:      results,
					SuccessCount: successCount,
					FailCount:    failCount,
					TotalCount:   len(ws.VMs),
				},
				Error: errInfo,
			})
			if failCount > 0 {
				os.Exit(1)
			}
			return
		}

		if failCount > 0 {
			os.Exit(1)
		}
	},
}

var wsStopCmd = &cobra.Command{
	Use:   "stop <name>",
	Short: "Stop all VMs in a workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ws, err := hyperv.LoadWorkspace(args[0])
		if err != nil {
			output.PrintError("WORKSPACE_GET_FAILED", "Failed to load workspace", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to load workspace: %v\n", err)
			}
			os.Exit(1)
		}

		manager := hyperv.NewManager()
		if !output.IsJSON() {
			fmt.Printf("🛑 Stopping workspace '%s' (%d VMs)...\n", ws.Name, len(ws.VMs))
		}

		results := make([]VMOperationResult, len(ws.VMs))
		successCount := 0
		failCount := 0

		var (
			mu sync.Mutex
			g  errgroup.Group
		)
		sem := make(chan struct{}, defaultBatchConcurrency)

		for i, vmName := range ws.VMs {
			i, vmName := i, vmName
			sem <- struct{}{}
			g.Go(func() error {
				defer func() { <-sem }()

				opCtx, cancel := context.WithTimeout(cmd.Context(), 60*time.Second)
				defer cancel()

				if !output.IsJSON() {
					mu.Lock()
					fmt.Printf("🛑 Stopping VM: %s...\n", vmName)
					mu.Unlock()
				}
				opErr := manager.StopVMByName(opCtx, vmName)
				res := VMOperationResult{
					Index: i + 1,
					Name:  vmName,
				}
				mu.Lock()
				defer mu.Unlock()
				if opErr != nil {
					res.Success = false
					res.Error = opErr.Error()
					failCount++
					if !output.IsJSON() {
						fmt.Printf("❌ Failed to stop VM '%s': %v\n", vmName, opErr)
					}
				} else {
					res.Success = true
					res.Message = "VM stopped successfully"
					successCount++
					if !output.IsJSON() {
						fmt.Printf("✅ VM '%s' stopped.\n", vmName)
					}
				}
				results[i] = res
				return nil
			})
		}
		_ = g.Wait()

		if results == nil {
			results = []VMOperationResult{}
		}

		if output.IsJSON() {
			topSuccess := failCount == 0
			var errInfo *output.ErrorInfo
			if failCount > 0 {
				errInfo = &output.ErrorInfo{
					Code:    "WORKSPACE_STOP_FAILED",
					Message: fmt.Sprintf("%d of %d VMs failed to stop", failCount, len(ws.VMs)),
				}
			}
			output.PrintResponse(output.Response{
				Success: topSuccess,
				Data: WorkspaceBatchResult{
					Workspace:    ws.Name,
					Operation:    "stop",
					Results:      results,
					SuccessCount: successCount,
					FailCount:    failCount,
					TotalCount:   len(ws.VMs),
				},
				Error: errInfo,
			})
			if failCount > 0 {
				os.Exit(1)
			}
			return
		}

		if failCount > 0 {
			os.Exit(1)
		}
	},
}

func init() {
	wsCreateCmd.Flags().StringVarP(&wsVms, "vms", "v", "", "Comma-separated list of VM names")
	_ = wsCreateCmd.MarkFlagRequired("vms")

	workspaceCmd.AddCommand(wsListCmd)
	workspaceCmd.AddCommand(wsCreateCmd)
	workspaceCmd.AddCommand(wsShowCmd)
	workspaceCmd.AddCommand(wsDeleteCmd)
	workspaceCmd.AddCommand(wsStartCmd)
	workspaceCmd.AddCommand(wsStopCmd)
	rootCmd.AddCommand(workspaceCmd)
}
