package cmd

import (
	"fmt"
	"strings"

	"quickvm/internal/hyperv"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all Hyper-V virtual machines",
	Long:    `Display a list of all Hyper-V virtual machines with their status.`,
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		manager := hyperv.NewManager()

		fmt.Println("📋 Fetching Hyper-V virtual machines...")
		vms, err := manager.GetVMs(cmd.Context())
		if err != nil {
			fmt.Printf("❌ Failed to get VMs: %v\n", err)
			return
		}

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
