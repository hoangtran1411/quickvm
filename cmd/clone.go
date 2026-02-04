package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"quickvm/internal/hyperv"
	"quickvm/internal/output"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <vm-index> <new-name>",
	Short: "Clone a VM with a new name (full clone)",
	Long: `Clone a Hyper-V virtual machine with a new name.

This performs a full clone operation:
1. Export the source VM to a temporary directory
2. Import with Copy and GenerateNewId flags
3. Rename to the specified new name
4. Cleanup temporary files

The cloned VM will be completely independent from the source VM.
This may take several minutes depending on the VM disk size.

Examples:
  quickvm clone 1 "WebServer-Copy"            # Clone VM 1 with new name
  quickvm clone 2 "TestVM"                    # Clone VM 2 with new name`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		manager := hyperv.NewManager()

		// Parse VM index
		index, err := strconv.Atoi(args[0])
		if err != nil {
			output.PrintError("INVALID_INDEX", "Invalid VM index", args[0])
			if !output.IsJSON() {
				fmt.Printf("❌ Invalid VM index: %s\n", args[0])
			}
			return
		}

		// Get new name and trim whitespace
		newName := strings.TrimSpace(args[1])
		if newName == "" {
			output.PrintError("INVALID_NAME", "New VM name cannot be empty", "")
			if !output.IsJSON() {
				fmt.Println("❌ New VM name cannot be empty")
			}
			return
		}

		// Get source VM name for display
		sourceName, err := manager.GetVMNameByIndex(cmd.Context(), index)
		if err != nil {
			output.PrintError("VM_GET_FAILED", "Failed to get source VM", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to get source VM: %v\n", err)
			}
			return
		}

		// Check if new name already exists
		exists, err := manager.VMExists(cmd.Context(), newName)
		if err != nil {
			output.PrintError("VM_CHECK_FAILED", "Failed to check VM name", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to check VM name: %v\n", err)
			}
			return
		}
		if exists {
			output.PrintError("VM_EXISTS", "A VM with this name already exists", newName)
			if !output.IsJSON() {
				fmt.Printf("❌ A VM with name '%s' already exists\n", newName)
			}
			return
		}

		if !output.IsJSON() {
			fmt.Printf("🔄 Cloning VM '%s' to '%s'...\n", sourceName, newName)
			fmt.Println("⏳ This may take several minutes depending on VM disk size...")
			fmt.Println()
			fmt.Println("Steps:")
			fmt.Println("  1. Exporting source VM...")
		}

		if err := manager.CloneVM(cmd.Context(), index, newName); err != nil {
			output.PrintError("CLONE_FAILED", "Failed to clone VM", err.Error())
			if !output.IsJSON() {
				fmt.Printf("\n❌ Failed to clone VM: %v\n", err)
			}
			return
		}

		// JSON output for AI agents
		if output.IsJSON() {
			output.PrintData(CloneResult{
				SourceName:  sourceName,
				SourceIndex: index,
				NewName:     newName,
				Success:     true,
				Message:     "VM cloned successfully",
			})
			return
		}

		fmt.Println("  2. Importing with new identity...")
		fmt.Println("  3. Renaming to target name...")
		fmt.Println("  4. Cleaning up temporary files...")
		fmt.Println()
		fmt.Printf("✅ VM '%s' cloned successfully to '%s'!\n", sourceName, newName)
		fmt.Println()
		fmt.Println("💡 Tips:")
		fmt.Printf("   - Start the cloned VM with: quickvm start <index>\n")
		fmt.Printf("   - The cloned VM has a new unique ID\n")
		fmt.Printf("   - The cloned VM is completely independent from the source\n")
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
