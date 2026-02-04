package cmd

import (
	"fmt"
	"strings"

	"quickvm/internal/hyperv"
	"quickvm/internal/output"

	"github.com/spf13/cobra"
)

// VMListResponse represents the JSON response for VM list
type VMListResponse struct {
	VMs   []hyperv.VM `json:"vms"`
	Total int         `json:"total"`
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all Hyper-V virtual machines",
	Long:    `Display a list of all Hyper-V virtual machines with their status.`,
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, _ []string) {
		manager := hyperv.NewManager()

		if !output.IsJSON() {
			fmt.Println("📋 Fetching Hyper-V virtual machines...")
		}

		vms, err := manager.GetVMs(cmd.Context())
		if err != nil {
			output.PrintError("VM_LIST_FAILED", "Failed to get VMs", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to get VMs: %v\n", err)
			}
			return
		}

		// JSON output for AI agents
		if output.IsJSON() {
			output.PrintData(VMListResponse{
				VMs:   vms,
				Total: len(vms),
			})
			return
		}

		// Human-readable output
		if len(vms) == 0 {
			fmt.Println("ℹ️  No virtual machines found.")
			return
		}

		// Print header
		fmt.Println("\n" + strings.Repeat("=", 110))
		fmt.Printf("%-7s %-30s %-12s %-8s %-12s %-20s %-15s\n",
			"Index", "Name", "State", "CPU%", "Memory(MB)", "Uptime", "Status")
		fmt.Println(strings.Repeat("=", 110))

		// Print VMs
		for _, vm := range vms {
			stateIcon := "⚪"
			switch strings.ToLower(vm.State) {
			case "running":
				stateIcon = "🟢"
			case "off":
				stateIcon = "🔴"
			case "paused":
				stateIcon = "🟡"
			}

			fmt.Printf("%-7d %-30s %s %-10s %-8d %-12d %-20s %-15s\n",
				vm.Index,
				vm.Name,
				stateIcon,
				vm.State,
				vm.CPUUsage,
				vm.MemoryMB,
				vm.Uptime,
				vm.Status,
			)
		}

		fmt.Println(strings.Repeat("=", 110))
		fmt.Printf("\nTotal VMs: %d\n", len(vms))
		fmt.Println("\n💡 Tip: Use 'quickvm start <index>' to start a VM")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
