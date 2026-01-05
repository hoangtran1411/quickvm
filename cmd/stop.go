package cmd

import (
	"fmt"
	"strconv"

	"quickvm/hyperv"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [vm-index]",
	Short: "Stop a Hyper-V virtual machine",
	Long: `Stop a Hyper-V virtual machine by its index.
Example: quickvm stop 1`,
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
		fmt.Printf("🛑 Stopping VM: %s (Index: %d)...\n", vm.Name, index)

		if err := manager.StopVM(index); err != nil {
			fmt.Printf("❌ Failed to stop VM: %v\n", err)
			return
		}

		fmt.Printf("✅ VM '%s' stopped successfully!\n", vm.Name)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
