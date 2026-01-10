package cmd

import (
	"fmt"

	"quickvm/hyperv"

	"github.com/spf13/cobra"
)

var (
	restartRange string
	restartAll   bool
)

var restartCmd = &cobra.Command{
	Use:   "restart [vm-indices]",
	Short: "Restart Hyper-V virtual machines",
	Long: `Restart one or more Hyper-V virtual machines by their indices.

Examples:
  quickvm restart 1 3 5       # Restart VMs at index 1, 3, and 5
  quickvm restart --range 1-5 # Restart VMs from index 1 to 5
  quickvm restart --all       # Restart all VMs`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runRestart(hyperv.NewManager(), args, restartRange, restartAll)
	},
}

func runRestart(manager hyperv.VMManager, args []string, rangeStr string, all bool) {
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
		fmt.Printf("🔄 Restarting %d VMs...\n\n", len(indices))
	}

	successCount := 0
	failCount := 0

	for _, index := range indices {
		vm := vms[index-1]
		fmt.Printf("🔄 Restarting VM: %s (Index: %d)...\n", vm.Name, index)

		if err := manager.RestartVM(index); err != nil {
			fmt.Printf("❌ Failed to restart VM '%s': %v\n", vm.Name, err)
			failCount++
		} else {
			fmt.Printf("✅ VM '%s' restarted successfully!\n", vm.Name)
			successCount++
		}
	}

	if len(indices) > 1 {
		fmt.Printf("\n📊 Summary: %d restarted, %d failed\n", successCount, failCount)
	}
}

func init() {
	restartCmd.Flags().StringVarP(&restartRange, "range", "r", "", "Range of VM indices to restart (e.g., '1-5')")
	restartCmd.Flags().BoolVarP(&restartAll, "all", "a", false, "Restart all virtual machines")
	rootCmd.AddCommand(restartCmd)
}
