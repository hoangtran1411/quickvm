package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"quickvm/internal/hyperv"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export <vm-index> <path>",
	Short: "Export a VM to a directory",
	Long: `Export a Hyper-V virtual machine to a specified directory.

The export includes all VM configuration files, checkpoints, and virtual hard disks.
The VM can be running or stopped during export.

Examples:
  quickvm export 1 "D:\Backups\VMs"            # Export VM 1 to D:\Backups\VMs
  quickvm export 2 "C:\Export\MyVM"            # Export VM 2 to C:\Export\MyVM
  quickvm export 1 .                            # Export VM 1 to current directory

The exported VM will be placed in a subdirectory named after the VM.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		manager := hyperv.NewManager()

		// Parse VM index
		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("❌ Invalid VM index: %s\n", args[0])
			return
		}

		// Get export path
		exportPath := args[1]

		// Resolve relative path
		if !filepath.IsAbs(exportPath) {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("❌ Failed to get current directory: %v\n", err)
				return
			}
			exportPath = filepath.Join(cwd, exportPath)
		}

		// Get VM name for display
		vmName, err := manager.GetVMNameByIndex(cmd.Context(), index)
		if err != nil {
			fmt.Printf("❌ Failed to get VM: %v\n", err)
			return
		}

		// Check if export path exists, create if not
		if _, err := os.Stat(exportPath); os.IsNotExist(err) {
			fmt.Printf("📁 Creating export directory: %s\n", exportPath)
			if err := os.MkdirAll(exportPath, 0755); err != nil {
				fmt.Printf("❌ Failed to create export directory: %v\n", err)
				return
			}
		}

		fmt.Printf("📦 Exporting VM '%s' to '%s'...\n", vmName, exportPath)
		fmt.Println("⏳ This may take a while depending on VM size...")

		if err := manager.ExportVM(cmd.Context(), index, exportPath); err != nil {
			fmt.Printf("❌ Failed to export VM: %v\n", err)
			return
		}

		// Show success message with export location
		exportedPath := filepath.Join(exportPath, vmName)
		fmt.Printf("\n✅ VM '%s' exported successfully!\n", vmName)
		fmt.Printf("📁 Export location: %s\n", exportedPath)
		fmt.Println("\n💡 Tips:")
		fmt.Printf("   - Import this VM with: quickvm import \"%s\"\n", exportedPath)
		fmt.Println("   - The export contains: VM config, checkpoints, and virtual hard disks")
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
