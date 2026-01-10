package cmd

import (
	"fmt"

	"quickvm/hyperv"

	"github.com/spf13/cobra"
)

var (
	startRange string
	startAll   bool
)

var startCmd = &cobra.Command{
	Use:   "start [vm-index]",
	Short: "Start a Hyper-V virtual machine",
	Long: `Start a Hyper-V virtual machine by its index.

Examples:
  quickvm start 1 3 5       # Start VMs at index 1, 3, and 5
  quickvm start --range 1-5 # Start VMs from index 1 to 5
  quickvm start --all       # Start all VMs`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runStart(hyperv.NewManager(), args, startRange, startAll)
	},
}

func runStart(manager hyperv.VMManager, args []string, rangeStr string, all bool) {
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
		fmt.Printf("🚀 Starting %d VMs...\n\n", len(indices))
	}

	successCount := 0
	failCount := 0

	for _, index := range indices {
		vm := vms[index-1]
		fmt.Printf("🚀 Starting VM: %s (Index: %d)...\n", vm.Name, index)

		if err := manager.StartVMByName(vm.Name); err != nil {
			fmt.Printf("❌ Failed to start VM '%s': %v\n", vm.Name, err)
			failCount++
		} else {
			fmt.Printf("✅ VM '%s' started successfully!\n", vm.Name)
			successCount++
		}
	}

	if len(indices) > 1 {
		fmt.Printf("\n📊 Summary: %d started, %d failed\n", successCount, failCount)
	}
}

func init() {
	startCmd.Flags().StringVarP(&startRange, "range", "r", "", "Range of VM indices to start (e.g., '1-5' or '1,3,5')")
	startCmd.Flags().BoolVarP(&startAll, "all", "a", false, "Start all virtual machines")
	rootCmd.AddCommand(startCmd)
}
