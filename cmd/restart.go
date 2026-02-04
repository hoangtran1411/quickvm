package cmd

import (
	"context"

	"quickvm/internal/hyperv"

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
		runRestart(cmd.Context(), hyperv.NewManager(), args, restartRange, restartAll)
	},
}

func runRestart(ctx context.Context, manager hyperv.VMManager, args []string, rangeStr string, all bool) {
	runVMBatchOperation(ctx, manager, args, rangeStr, all, VMOperationConfig{
		Operation:   "restart",
		ActionVerb:  "Restarting",
		ActionEmoji: "🔄",
		SuccessVerb: "restarted",
		OperationFunc: func(ctx context.Context, mgr hyperv.VMManager, vm hyperv.VM) error {
			return mgr.RestartVMByName(ctx, vm.Name)
		},
	})
}

func init() {
	restartCmd.Flags().StringVarP(&restartRange, "range", "r", "", "Range of VM indices to restart (e.g., '1-5')")
	restartCmd.Flags().BoolVarP(&restartAll, "all", "a", false, "Restart all virtual machines")
	rootCmd.AddCommand(restartCmd)
}
