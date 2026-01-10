package cmd

import (
	"fmt"

	"quickvm/hyperv"

	"github.com/spf13/cobra"
)

var (
	stopRange string
	stopAll   bool
)

var stopCmd = &cobra.Command{
	Use:   "stop [vm-indices]",
	Short: "Stop Hyper-V virtual machines",
	Long: `Stop one or more Hyper-V virtual machines by their indices.

Examples:
  quickvm stop 1 3 5       # Stop VMs at index 1, 3, and 5
  quickvm stop --range 1-5 # Stop VMs from index 1 to 5
  quickvm stop --all       # Stop all VMs`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runStop(hyperv.NewManager(), args, stopRange, stopAll)
	},
}

func runStop(manager hyperv.VMManager, args []string, rangeStr string, all bool) {
	// Get VMs to validate index and get name
	vms, err := manager.GetVMs()
	if err != nil {
		fmt.Printf("❌ Failed to get VMs: %v\n", err)
		return
	}

	// Use shared getIndices logic
	indices, err := getIndices(args, rangeStr, all, len(vms))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(indices) > 1 {
		fmt.Printf("🛑 Stopping %d VMs...\n\n", len(indices))
	}

	successCount := 0
	failCount := 0

	for _, index := range indices {
		vm := vms[index-1]
		fmt.Printf("🛑 Stopping VM: %s (Index: %d)...\n", vm.Name, index)

		if err := manager.StopVM(index); err != nil {
			fmt.Printf("❌ Failed to stop VM '%s': %v\n", vm.Name, err)
			failCount++
		} else {
			fmt.Printf("✅ VM '%s' stopped successfully!\n", vm.Name)
			successCount++
		}
	}

	if len(indices) > 1 {
		fmt.Printf("\n📊 Summary: %d stopped, %d failed\n", successCount, failCount)
	}
}

func init() {
	stopCmd.Flags().StringVarP(&stopRange, "range", "r", "", "Range of VM indices to stop (e.g., '1-5')")
	stopCmd.Flags().BoolVarP(&stopAll, "all", "a", false, "Stop all virtual machines")
	rootCmd.AddCommand(stopCmd)
}
