package cmd

import (
	"fmt"
	"strconv"

	"quickvm/hyperv"

	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart [vm-index]",
	Short: "Restart a Hyper-V virtual machine",
	Long: `Restart a Hyper-V virtual machine by its index.
Example: quickvm restart 1`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("❌ Invalid VM index: %s\n", args[0])
			return
		}

		manager := hyperv.NewManager()
		
		// Get VMs to validate index and get name
		vms, err := manager.GetVMs()
		if err != nil {
			fmt.Printf("❌ Failed to get VMs: %v\n", err)
			return
		}

		if index < 1 || index > len(vms) {
			fmt.Printf("❌ Invalid VM index: %d (valid range: 1-%d)\n", index, len(vms))
			return
		}

		vm := vms[index-1]
		fmt.Printf("🔄 Restarting VM: %s (Index: %d)...\n", vm.Name, index)

		if err := manager.RestartVM(index); err != nil {
			fmt.Printf("❌ Failed to restart VM: %v\n", err)
			return
		}

		fmt.Printf("✅ VM '%s' restarted successfully!\n", vm.Name)
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
