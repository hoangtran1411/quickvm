package cmd

import (
	"context"

	"quickvm/internal/hyperv"

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
		runStart(cmd.Context(), hyperv.NewManager(), args, startRange, startAll)
	},
}

func runStart(ctx context.Context, manager hyperv.VMManager, args []string, rangeStr string, all bool) {
	runVMBatchOperation(ctx, manager, args, rangeStr, all, VMOperationConfig{
		Operation:   "start",
		ActionVerb:  "Starting",
		ActionEmoji: "🚀",
		SuccessVerb: "started",
		OperationFunc: func(ctx context.Context, mgr hyperv.VMManager, vm hyperv.VM) error {
			return mgr.StartVMByName(ctx, vm.Name)
		},
	})
}

func init() {
	startCmd.Flags().StringVarP(&startRange, "range", "r", "", "Range of VM indices to start (e.g., '1-5' or '1,3,5')")
	startCmd.Flags().BoolVarP(&startAll, "all", "a", false, "Start all virtual machines")
	rootCmd.AddCommand(startCmd)
}
