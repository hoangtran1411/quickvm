package cmd

import (
	"fmt"
	"strings"

	"quickvm/hyperv"

	"github.com/spf13/cobra"
)

var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"ws"},
	Short:   "Manage VM workspaces (groups)",
	Long:    `Manage groups of virtual machines using workspace profiles.`,
}

var wsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		names, err := hyperv.ListWorkspaces()
		if err != nil {
			fmt.Printf("❌ Failed to list workspaces: %v\n", err)
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
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		vms := strings.Split(wsVms, ",")
		for i, v := range vms {
			vms[i] = strings.TrimSpace(v)
		}

		ws := &hyperv.Workspace{
			Name:        name,
			Description: "Created via CLI",
			VMs:         vms,
		}

		if err := hyperv.SaveWorkspace(ws); err != nil {
			fmt.Printf("❌ Failed to save workspace: %v\n", err)
			return
		}

		fmt.Printf("✅ Workspace '%s' created successfully!\n", name)
	},
}

var wsShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show workspace details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ws, err := hyperv.LoadWorkspace(args[0])
		if err != nil {
			fmt.Printf("❌ Failed to load workspace: %v\n", err)
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
	Run: func(cmd *cobra.Command, args []string) {
		if err := hyperv.DeleteWorkspace(args[0]); err != nil {
			fmt.Printf("❌ Failed to delete workspace: %v\n", err)
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
			fmt.Printf("❌ Failed to load workspace: %v\n", err)
			return
		}

		manager := hyperv.NewManager()
		fmt.Printf("🚀 Starting workspace '%s' (%d VMs)...\n", ws.Name, len(ws.VMs))

		for _, vmName := range ws.VMs {
			fmt.Printf("🚀 Starting VM: %s...\n", vmName)
			if err := manager.StartVMByName(vmName); err != nil {
				fmt.Printf("❌ Failed to start VM '%s': %v\n", vmName, err)
			} else {
				fmt.Printf("✅ VM '%s' started.\n", vmName)
			}
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
			fmt.Printf("❌ Failed to load workspace: %v\n", err)
			return
		}

		manager := hyperv.NewManager()
		fmt.Printf("🛑 Stopping workspace '%s' (%d VMs)...\n", ws.Name, len(ws.VMs))

		for _, vmName := range ws.VMs {
			fmt.Printf("🛑 Stopping VM: %s...\n", vmName)
			// Need StopVMByName
			if err := manager.StopVMByName(vmName); err != nil {
				fmt.Printf("❌ Failed to stop VM '%s': %v\n", vmName, err)
			} else {
				fmt.Printf("✅ VM '%s' stopped.\n", vmName)
			}
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
