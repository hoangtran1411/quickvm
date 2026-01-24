package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"quickvm/internal/hyperv"

	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage VM snapshots/checkpoints",
	Long: `Manage Hyper-V VM snapshots (checkpoints).

Snapshots allow you to save the current state of a VM and restore it later.
This is useful for testing, rollback, and backup purposes.

Available subcommands:
  list    - List all snapshots for a VM
  create  - Create a new snapshot
  restore - Restore a VM to a snapshot
  delete  - Delete a snapshot`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var snapshotListCmd = &cobra.Command{
	Use:   "list <vm-index>",
	Short: "List all snapshots for a VM",
	Long: `List all snapshots (checkpoints) for a specific VM.

Example:
  quickvm snapshot list 1    # List snapshots for VM at index 1`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := hyperv.NewManager()

		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("❌ Invalid VM index: %s\n", args[0])
			return
		}

		// Get VM name for display
		vmName, err := manager.GetVMNameByIndex(cmd.Context(), index)
		if err != nil {
			fmt.Printf("❌ Failed to get VM: %v\n", err)
			return
		}

		fmt.Printf("📸 Snapshots for VM: %s (Index: %d)\n\n", vmName, index)

		snapshots, err := manager.GetSnapshots(cmd.Context(), index)
		if err != nil {
			fmt.Printf("❌ Failed to get snapshots: %v\n", err)
			return
		}

		if len(snapshots) == 0 {
			fmt.Println("📭 No snapshots found for this VM.")
			fmt.Println("\n💡 Tip: Create a snapshot with: quickvm snapshot create", index, "\"Snapshot Name\"")
			return
		}

		// Print header
		fmt.Printf("%-4s %-30s %-20s %-15s\n", "#", "Name", "Created", "Type")
		fmt.Println(strings.Repeat("-", 75))

		// Print snapshots
		for i, snapshot := range snapshots {
			fmt.Printf("%-4d %-30s %-20s %-15s\n",
				i+1,
				truncateString(snapshot.Name, 28),
				snapshot.CreationTime,
				snapshot.SnapshotType,
			)
		}

		fmt.Printf("\n📊 Total: %d snapshot(s)\n", len(snapshots))
	},
}

var snapshotCreateCmd = &cobra.Command{
	Use:   "create <vm-index> <snapshot-name>",
	Short: "Create a new snapshot for a VM",
	Long: `Create a new snapshot (checkpoint) for a specific VM.

The VM can be running or stopped when creating a snapshot.

Examples:
  quickvm snapshot create 1 "Before Update"     # Create snapshot for VM 1
  quickvm snapshot create 2 "Clean State"       # Create snapshot for VM 2`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		manager := hyperv.NewManager()

		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("❌ Invalid VM index: %s\n", args[0])
			return
		}

		snapshotName := args[1]

		// Get VM name for display
		vmName, err := manager.GetVMNameByIndex(cmd.Context(), index)
		if err != nil {
			fmt.Printf("❌ Failed to get VM: %v\n", err)
			return
		}

		fmt.Printf("📸 Creating snapshot '%s' for VM: %s...\n", snapshotName, vmName)

		if err := manager.CreateSnapshot(cmd.Context(), index, snapshotName); err != nil {
			fmt.Printf("❌ Failed to create snapshot: %v\n", err)
			return
		}

		fmt.Printf("✅ Snapshot '%s' created successfully!\n", snapshotName)
		fmt.Println("\n💡 Tip: View snapshots with: quickvm snapshot list", index)
	},
}

var snapshotRestoreCmd = &cobra.Command{
	Use:   "restore <vm-index> <snapshot-name>",
	Short: "Restore a VM to a snapshot",
	Long: `Restore a VM to a previously saved snapshot (checkpoint).

⚠️  Warning: This will revert the VM to the state when the snapshot was taken.
               Any changes made after the snapshot will be lost!

The VM should be stopped before restoring a snapshot.

Examples:
  quickvm snapshot restore 1 "Before Update"   # Restore VM 1 to snapshot
  quickvm snapshot restore 2 "Clean State"     # Restore VM 2 to snapshot`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		manager := hyperv.NewManager()

		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("❌ Invalid VM index: %s\n", args[0])
			return
		}

		snapshotName := args[1]

		// Get VM name for display
		vmName, err := manager.GetVMNameByIndex(cmd.Context(), index)
		if err != nil {
			fmt.Printf("❌ Failed to get VM: %v\n", err)
			return
		}

		fmt.Printf("⏮️  Restoring VM '%s' to snapshot '%s'...\n", vmName, snapshotName)
		fmt.Println("⚠️  Warning: All changes after this snapshot will be lost!")

		if err := manager.RestoreSnapshot(cmd.Context(), index, snapshotName); err != nil {
			fmt.Printf("❌ Failed to restore snapshot: %v\n", err)
			return
		}

		fmt.Printf("✅ VM '%s' restored to snapshot '%s' successfully!\n", vmName, snapshotName)
	},
}

var snapshotDeleteCmd = &cobra.Command{
	Use:   "delete <vm-index> <snapshot-name>",
	Short: "Delete a snapshot from a VM",
	Long: `Delete a snapshot (checkpoint) from a VM.

⚠️  Warning: This action cannot be undone!

Examples:
  quickvm snapshot delete 1 "Old Snapshot"    # Delete snapshot from VM 1
  quickvm snapshot delete 2 "Test Snapshot"   # Delete snapshot from VM 2`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		manager := hyperv.NewManager()

		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("❌ Invalid VM index: %s\n", args[0])
			return
		}

		snapshotName := args[1]

		// Get VM name for display
		vmName, err := manager.GetVMNameByIndex(cmd.Context(), index)
		if err != nil {
			fmt.Printf("❌ Failed to get VM: %v\n", err)
			return
		}

		fmt.Printf("🗑️  Deleting snapshot '%s' from VM '%s'...\n", snapshotName, vmName)

		if err := manager.DeleteSnapshot(cmd.Context(), index, snapshotName); err != nil {
			fmt.Printf("❌ Failed to delete snapshot: %v\n", err)
			return
		}

		fmt.Printf("✅ Snapshot '%s' deleted successfully!\n", snapshotName)
	},
}

// truncateString truncates a string to max length and adds "..." if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func init() {
	// Add subcommands to snapshot command
	snapshotCmd.AddCommand(snapshotListCmd)
	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCmd.AddCommand(snapshotRestoreCmd)
	snapshotCmd.AddCommand(snapshotDeleteCmd)

	// Add snapshot command to root
	rootCmd.AddCommand(snapshotCmd)
}
