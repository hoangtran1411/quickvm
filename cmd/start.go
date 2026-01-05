package cmd

import (
	"fmt"
	"strconv"

	"quickvm/hyperv"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [vm-index]",
	Short: "Start a Hyper-V virtual machine",
	Long: `Start a Hyper-V virtual machine by its index.
Example: quickvm start 1`,
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
		fmt.Printf("🚀 Starting VM: %s (Index: %d)...\n", vm.Name, index)

		if err := manager.StartVM(index); err != nil {
			fmt.Printf("❌ Failed to start VM: %v\n", err)
			return
		}

		fmt.Printf("✅ VM '%s' started successfully!\n", vm.Name)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
