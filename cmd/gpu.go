package cmd

import (
	"fmt"
	"os"
	"strconv"

	"quickvm/internal/hyperv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var gpuCmd = &cobra.Command{
	Use:   "gpu",
	Short: "Manage GPU passthrough for virtual machines",
	Long: `Manage GPU-P (GPU Partitioning) for Hyper-V virtual machines.

GPU-P allows sharing the host GPU with virtual machines for graphics acceleration.

Examples:
  quickvm gpu status           # Check GPU partitioning support
  quickvm gpu add 1           # Add GPU partition to VM #1
  quickvm gpu remove 1        # Remove GPU partition from VM #1
  quickvm gpu drivers         # Show GPU driver paths for copying`,
}

var gpuStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check GPU partitioning support",
	Long:  `Check if the system has GPUs that support partitioning for VM passthrough.`,
	Run: func(cmd *cobra.Command, _ []string) {
		manager := hyperv.NewManager()

		color.Cyan("🔍 Checking GPU partitioning support...")
		fmt.Println()

		gpus, err := manager.CheckGPUPartitionable(cmd.Context())
		if err != nil {
			color.Red("❌ Error checking GPU support: %v", err)
			os.Exit(1)
		}

		if len(gpus) == 0 {
			color.Yellow("⚠️  No GPUs with partitioning support found.")
			fmt.Println()
			color.White("   Possible reasons:")
			color.White("   • GPU does not support GPU-P")
			color.White("   • GPU drivers are outdated")
			color.White("   • Hyper-V is not enabled")
			return
		}

		color.Green("✅ Found %d GPU(s) with partitioning support:", len(gpus))
		fmt.Println()

		for i, gpu := range gpus {
			color.Cyan("  GPU #%d:", i+1)
			color.White("    Name: %s", gpu.Name)
			color.White("    Partition Count: %d", gpu.PartitionCount)
			if gpu.MaxPartitionVRAM > 0 {
				color.White("    Max VRAM: %d", gpu.MaxPartitionVRAM)
			}
			fmt.Println()
		}
	},
}

var gpuAddCmd = &cobra.Command{
	Use:   "add [vm-index]",
	Short: "Add GPU partition to a virtual machine",
	Long: `Add GPU partition to a Hyper-V virtual machine.

The VM must be stopped before adding GPU passthrough.
After adding, you'll need to copy GPU drivers to the guest VM.

Example: quickvm gpu add 1`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil {
			color.Red("❌ Invalid VM index: %s", args[0])
			return
		}

		manager := hyperv.NewManager()

		// Check admin privileges
		if !hyperv.IsRunningAsAdmin(cmd.Context()) {
			color.Red("❌ This command requires Administrator privileges.")
			fmt.Println()
			color.Yellow("💡 Please run this command in an elevated PowerShell or Command Prompt.")
			os.Exit(1)
		}

		// Check GPU support first
		color.Cyan("🔍 Checking GPU partitioning support...")
		gpus, err := manager.CheckGPUPartitionable(cmd.Context())
		if err != nil {
			color.Red("❌ Error checking GPU support: %v", err)
			os.Exit(1)
		}

		if len(gpus) == 0 {
			color.Red("❌ No GPUs with partitioning support found.")
			color.Yellow("💡 Your GPU may not support GPU-P or drivers need updating.")
			os.Exit(1)
		}

		// Get VMs to validate index
		vms, err := manager.GetVMs(cmd.Context())
		if err != nil {
			color.Red("❌ Failed to get VMs: %v", err)
			os.Exit(1)
		}

		if index < 1 || index > len(vms) {
			color.Red("❌ Invalid VM index: %d (valid range: 1-%d)", index, len(vms))
			return
		}

		vm := vms[index-1]
		color.Cyan("🔧 Adding GPU partition to VM: %s", vm.Name)
		fmt.Println()

		// Check if VM is running
		if vm.State == "Running" {
			color.Red("❌ VM '%s' is currently running.", vm.Name)
			color.Yellow("💡 Please stop the VM first: quickvm stop %d", index)
			os.Exit(1)
		}

		// Add GPU partition with default config
		config := hyperv.DefaultGPUPartitionConfig()
		if err := manager.AddGPUPartition(cmd.Context(), vm.Name, config); err != nil {
			color.Red("❌ Failed to add GPU partition: %v", err)
			os.Exit(1)
		}

		color.Green("✅ GPU partition added successfully to '%s'!", vm.Name)
		fmt.Println()

		// Show driver copy instructions
		color.Yellow("⚠️  Important: You need to copy GPU drivers to the guest VM.")
		fmt.Println()
		color.Cyan("📋 Driver Copy Instructions:")
		fmt.Println()

		driverPaths, _ := manager.GetGPUDriverPaths(cmd.Context())
		if len(driverPaths) > 0 {
			color.White("   1. Copy driver folder from Host to Guest:")
			for _, path := range driverPaths {
				color.White("      FROM: %s", path)
			}
			color.White("      TO:   C:\\Windows\\System32\\HostDriverStore\\FileRepository\\")
			fmt.Println()
		}

		color.White("   2. Copy system files from Host to Guest:")
		color.White("      FROM: C:\\Windows\\System32\\nv*.*")
		color.White("      TO:   C:\\Windows\\System32\\")
		fmt.Println()

		color.Cyan("ℹ️  For detailed instructions, see: docs/GPU_PASSTHROUGH.md")
	},
}

var gpuRemoveCmd = &cobra.Command{
	Use:   "remove [vm-index]",
	Short: "Remove GPU partition from a virtual machine",
	Long: `Remove GPU partition from a Hyper-V virtual machine.

The VM must be stopped before removing GPU passthrough.

Example: quickvm gpu remove 1`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil {
			color.Red("❌ Invalid VM index: %s", args[0])
			return
		}

		manager := hyperv.NewManager()

		// Check admin privileges
		if !hyperv.IsRunningAsAdmin(cmd.Context()) {
			color.Red("❌ This command requires Administrator privileges.")
			fmt.Println()
			color.Yellow("💡 Please run this command in an elevated PowerShell or Command Prompt.")
			os.Exit(1)
		}

		// Get VMs to validate index
		vms, err := manager.GetVMs(cmd.Context())
		if err != nil {
			color.Red("❌ Failed to get VMs: %v", err)
			os.Exit(1)
		}

		if index < 1 || index > len(vms) {
			color.Red("❌ Invalid VM index: %d (valid range: 1-%d)", index, len(vms))
			return
		}

		vm := vms[index-1]
		color.Cyan("🔧 Removing GPU partition from VM: %s", vm.Name)
		fmt.Println()

		// Check if VM is running
		if vm.State == "Running" {
			color.Red("❌ VM '%s' is currently running.", vm.Name)
			color.Yellow("💡 Please stop the VM first: quickvm stop %d", index)
			os.Exit(1)
		}

		// Remove GPU partition
		if err := manager.RemoveGPUPartition(cmd.Context(), vm.Name); err != nil {
			color.Red("❌ Failed to remove GPU partition: %v", err)
			os.Exit(1)
		}

		color.Green("✅ GPU partition removed successfully from '%s'!", vm.Name)
	},
}

var gpuDriversCmd = &cobra.Command{
	Use:   "drivers",
	Short: "Show GPU driver paths for copying to guest",
	Long:  `Display the GPU driver file paths that need to be copied to the guest VM.`,
	Run: func(cmd *cobra.Command, _ []string) {
		manager := hyperv.NewManager()

		color.Cyan("🔍 Searching for GPU driver files...")
		fmt.Println()

		paths, err := manager.GetGPUDriverPaths(cmd.Context())
		if err != nil {
			color.Red("❌ Error getting driver paths: %v", err)
			os.Exit(1)
		}

		if len(paths) == 0 {
			color.Yellow("⚠️  No GPU driver folders found.")
			color.White("   Looking for NVIDIA (nv_dispi.inf_*) or AMD (u0*) drivers.")
			return
		}

		color.Green("✅ Found GPU driver folder(s):")
		fmt.Println()

		for _, path := range paths {
			color.White("   📁 %s", path)
		}

		fmt.Println()
		color.Cyan("📋 Copy Instructions:")
		fmt.Println()
		color.White("   1. Copy the driver folder(s) above to the guest VM:")
		color.White("      C:\\Windows\\System32\\HostDriverStore\\FileRepository\\")
		fmt.Println()
		color.White("   2. Copy all nv*.* files from Host to Guest:")
		color.White("      FROM: C:\\Windows\\System32\\nv*.*")
		color.White("      TO:   C:\\Windows\\System32\\")
	},
}

func init() {
	rootCmd.AddCommand(gpuCmd)
	gpuCmd.AddCommand(gpuStatusCmd)
	gpuCmd.AddCommand(gpuAddCmd)
	gpuCmd.AddCommand(gpuRemoveCmd)
	gpuCmd.AddCommand(gpuDriversCmd)
}
