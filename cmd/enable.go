package cmd

import (
	"fmt"
	"os"

	"quickvm/internal/hyperv"
	"quickvm/internal/output"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	forceRestart bool
	noRestart    bool
)

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable Hyper-V on this system",
	Long: `Enable Hyper-V feature on this Windows system.

This command requires Administrator privileges and may require a system restart.

Examples:
  quickvm enable              # Enable Hyper-V (prompt for restart)
  quickvm enable -y           # Enable Hyper-V and restart immediately
  quickvm enable --no-restart # Enable Hyper-V without restarting`,
	Run: func(cmd *cobra.Command, _ []string) {
		manager := hyperv.NewManager()

		// First check current status
		info, err := manager.GetSystemInfo(cmd.Context(), false)
		if err != nil {
			output.PrintError("STATUS_CHECK_FAILED", "Error checking Hyper-V status", err.Error())
			if !output.IsJSON() {
				color.Red("❌ Error checking Hyper-V status: %v", err)
			}
			os.Exit(1)
		}

		if info.HyperV.Enabled {
			if output.IsJSON() {
				output.PrintData(EnableResult{
					AlreadyEnabled:   true,
					NeedsRestart:     false,
					RestartScheduled: false,
					Success:          true,
					Message:          fmt.Sprintf("Hyper-V is already enabled (Status: %s)", info.HyperV.Status),
				})
				return
			}
			color.Green("✅ Hyper-V is already enabled on this system!")
			fmt.Println()
			color.Cyan("ℹ️  Status: %s", info.HyperV.Status)
			return
		}

		// Check if running as administrator
		if !hyperv.IsRunningAsAdmin(cmd.Context()) {
			output.PrintError("ADMIN_REQUIRED", "Administrator privileges required", "This command requires Administrator privileges")
			if !output.IsJSON() {
				color.Yellow("⚠️  Hyper-V is currently disabled on this system.")
				fmt.Println()
				color.Red("❌ This command requires Administrator privileges.")
				fmt.Println()
				color.Yellow("💡 Please run this command in an elevated PowerShell or Command Prompt:")
				color.White("   1. Right-click on PowerShell/Terminal")
				color.White("   2. Select 'Run as administrator'")
				color.White("   3. Run 'quickvm enable' again")
			}
			os.Exit(1)
		}

		if !output.IsJSON() {
			color.Yellow("⚠️  Hyper-V is currently disabled on this system.")
			fmt.Println()
			color.Cyan("🔧 Enabling Hyper-V...")
			fmt.Println()
		}

		// Enable Hyper-V
		needsRestart, err := manager.EnableHyperV(cmd.Context())
		if err != nil {
			output.PrintError("ENABLE_FAILED", "Failed to enable Hyper-V", err.Error())
			if !output.IsJSON() {
				color.Red("❌ Failed to enable Hyper-V: %v", err)
			}
			os.Exit(1)
		}

		if output.IsJSON() {
			restartScheduled := false
			if needsRestart && forceRestart {
				_ = manager.ScheduleRestart(cmd.Context(), 10)
				restartScheduled = true
			}
			output.PrintData(EnableResult{
				AlreadyEnabled:   false,
				NeedsRestart:     needsRestart,
				RestartScheduled: restartScheduled,
				Success:          true,
				Message:          "Hyper-V has been enabled successfully",
			})
			return
		}

		color.Green("✅ Hyper-V has been enabled successfully!")
		fmt.Println()

		if needsRestart {
			if noRestart {
				color.Yellow("⚠️  A system restart is required to complete the installation.")
				color.Cyan("ℹ️  Please restart your computer manually when ready.")
			} else if forceRestart {
				color.Yellow("🔄 Restarting your computer in 10 seconds...")
				color.White("   Press Ctrl+C to cancel the restart.")
				fmt.Println()

				if err := manager.ScheduleRestart(cmd.Context(), 10); err != nil {
					color.Red("❌ Failed to schedule restart: %v", err)
					color.Yellow("💡 Please restart your computer manually.")
				}
			} else {
				color.Yellow("⚠️  A system restart is required to complete the installation.")
				fmt.Println()
				fmt.Print("❓ Do you want to restart now? [y/N]: ")

				var response string
				if _, err := fmt.Scanln(&response); err != nil {
					response = "n"
				}

				if response == "y" || response == "Y" {
					color.Yellow("🔄 Restarting your computer in 10 seconds...")
					color.White("   Press Ctrl+C to cancel the restart.")
					fmt.Println()

					if err := manager.ScheduleRestart(cmd.Context(), 10); err != nil {
						color.Red("❌ Failed to schedule restart: %v", err)
						color.Yellow("💡 Please restart your computer manually.")
					}
				} else {
					color.Cyan("ℹ️  Please restart your computer manually when ready.")
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)
	enableCmd.Flags().BoolVarP(&forceRestart, "yes", "y", false, "Restart immediately without prompting")
	enableCmd.Flags().BoolVar(&noRestart, "no-restart", false, "Don't restart after enabling (manual restart required)")
}
