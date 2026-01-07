package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"quickvm/hyperv"

	"github.com/spf13/cobra"
)

var (
	importCopy          bool
	importGenerateNewID bool
	importVHDPath       string
)

var importCmd = &cobra.Command{
	Use:   "import <path>",
	Short: "Import a VM from an export directory",
	Long: `Import a Hyper-V virtual machine from an exported directory.

The import reads the VM configuration from a previous export and
registers or copies it into Hyper-V.

Examples:
  quickvm import "D:\Backups\VMs\MyVM"              # Import VM from directory
  quickvm import "D:\Backups\VMs\MyVM" --copy       # Copy VM files to default location
  quickvm import "D:\Backups\VMs\MyVM" --new-id     # Generate new VM ID
  quickvm import "D:\Exports\VM" --vhd-path "E:\VHDs"  # Specify VHD destination

Flags:
  --copy       Copy the VM files instead of registering in place
  --new-id     Generate a new unique ID for the imported VM
  --vhd-path   Specify a custom path for virtual hard disk files`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := hyperv.NewManager()

		// Get import path
		importPath := args[0]

		// Resolve relative path
		if !filepath.IsAbs(importPath) {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("❌ Failed to get current directory: %v\n", err)
				return
			}
			importPath = filepath.Join(cwd, importPath)
		}

		// Verify path exists
		if _, err := os.Stat(importPath); os.IsNotExist(err) {
			fmt.Printf("❌ Import path does not exist: %s\n", importPath)
			return
		}

		fmt.Printf("📦 Importing VM from '%s'...\n", importPath)

		// Show options being used
		if importCopy {
			fmt.Println("   📋 Mode: Copy (will copy VM files to default location)")
		} else {
			fmt.Println("   📋 Mode: Register in place")
		}

		if importGenerateNewID {
			fmt.Println("   🔄 Generating new VM ID")
		}

		if importVHDPath != "" {
			fmt.Printf("   💾 VHD destination: %s\n", importVHDPath)
		}

		fmt.Println("⏳ This may take a while depending on VM size...")

		// Build import options
		opts := hyperv.ImportVMOptions{
			Path:          importPath,
			Copy:          importCopy,
			GenerateNewID: importGenerateNewID,
			VHDPath:       importVHDPath,
		}

		vmName, err := manager.ImportVM(opts)
		if err != nil {
			fmt.Printf("❌ Failed to import VM: %v\n", err)
			return
		}

		fmt.Printf("\n✅ VM '%s' imported successfully!\n", vmName)
		fmt.Println("\n💡 Tips:")
		fmt.Println("   - List all VMs with: quickvm list")
		fmt.Printf("   - Start the VM with: quickvm start <index>\n")
	},
}

func init() {
	// Add flags
	importCmd.Flags().BoolVarP(&importCopy, "copy", "c", false, "Copy VM files instead of registering in place")
	importCmd.Flags().BoolVarP(&importGenerateNewID, "new-id", "n", false, "Generate a new unique ID for the VM")
	importCmd.Flags().StringVarP(&importVHDPath, "vhd-path", "v", "", "Custom destination path for VHD files")

	rootCmd.AddCommand(importCmd)
}
