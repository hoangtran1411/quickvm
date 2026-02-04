package cmd

import (
	"context"

	"quickvm/internal/hyperv"

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
		runStop(cmd.Context(), hyperv.NewManager(), args, stopRange, stopAll)
	},
}

func runStop(ctx context.Context, manager hyperv.VMManager, args []string, rangeStr string, all bool) {
	runVMBatchOperation(ctx, manager, args, rangeStr, all, VMOperationConfig{
		Operation:   "stop",
		ActionVerb:  "Stopping",
		ActionEmoji: "🛑",
		SuccessVerb: "stopped",
		OperationFunc: func(ctx context.Context, mgr hyperv.VMManager, vm hyperv.VM) error {
			return mgr.StopVMByName(ctx, vm.Name)
		},
	})
}

func init() {
	stopCmd.Flags().StringVarP(&stopRange, "range", "r", "", "Range of VM indices to stop (e.g., '1-5')")
	stopCmd.Flags().BoolVarP(&stopAll, "all", "a", false, "Stop all virtual machines")
	rootCmd.AddCommand(stopCmd)
}
