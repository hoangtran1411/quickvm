package cmd

import (
	"fmt"
	"os"

	"quickvm/internal/output"
	"quickvm/updater"

	"github.com/spf13/cobra"
)

var (
	autoInstall bool
	checkOnly   bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for updates and install the latest version",
	Long: `Check for new versions of QuickVM from GitHub releases.
If a new version is available, download and install it automatically.`,
	Run: func(_ *cobra.Command, _ []string) {
		if !output.IsJSON() {
			fmt.Println("🔍 Checking for updates...")
		}

		u := updater.NewUpdater(Version)

		release, hasUpdate, err := u.CheckForUpdates()
		if err != nil {
			output.PrintError("UPDATE_CHECK_FAILED", "Failed to check for updates", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to check for updates: %v\n", err)
				fmt.Println("💡 Tip: Check your internet connection and try again")
			}
			os.Exit(1)
		}

		if !hasUpdate {
			if output.IsJSON() {
				output.PrintData(UpdateResult{
					HasUpdate:      false,
					CurrentVersion: Version,
					LatestVersion:  Version,
					Installed:      false,
					Success:        true,
					Message:        "QuickVM is up to date",
				})
				return
			}
			fmt.Println("✅ You are already using the latest version!")
			fmt.Printf("   Current version: %s\n", Version)
			return
		}

		if checkOnly {
			if output.IsJSON() {
				output.PrintData(UpdateResult{
					HasUpdate:      true,
					CurrentVersion: Version,
					LatestVersion:  release.TagName,
					ReleaseNotes:   release.Body,
					Installed:      false,
					Success:        true,
					Message:        "New version available",
				})
				return
			}
			fmt.Printf("🎉 New version available: %s\n", release.TagName)
			fmt.Printf("   Current version: %s\n", Version)
			fmt.Println()
			fmt.Println("📋 Release Notes:")
			fmt.Println(release.Body)
			fmt.Println()
			fmt.Println("💡 Run 'quickvm update' without --check-only to install")
			return
		}

		if output.IsJSON() && !autoInstall {
			output.PrintData(UpdateResult{
				HasUpdate:      true,
				CurrentVersion: Version,
				LatestVersion:  release.TagName,
				ReleaseNotes:   release.Body,
				Installed:      false,
				Success:        true,
				Message:        "Update available. Run with -y/--yes to install non-interactively.",
			})
			return
		}

		if !output.IsJSON() {
			fmt.Printf("🎉 New version available: %s\n", release.TagName)
			fmt.Printf("   Current version: %s\n", Version)
			fmt.Println()
		}

		if !autoInstall {
			fmt.Print("❓ Do you want to install this update? [y/N]: ")
			var response string
			if _, err := fmt.Scanln(&response); err != nil {
				// Default to 'no' if can't read input
				response = "n"
			}

			if response != "y" && response != "Y" {
				fmt.Println("⏭️  Update cancelled")
				return
			}
		}

		if !output.IsJSON() {
			fmt.Println()
		}
		if err := u.DownloadAndInstall(release); err != nil {
			output.PrintError("UPDATE_INSTALL_FAILED", "Update failed", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Update failed: %v\n", err)
			}
			os.Exit(1)
		}

		if output.IsJSON() {
			output.PrintData(UpdateResult{
				HasUpdate:      true,
				CurrentVersion: Version,
				LatestVersion:  release.TagName,
				ReleaseNotes:   release.Body,
				Installed:      true,
				Success:        true,
				Message:        "Update installed successfully",
			})
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVarP(&autoInstall, "yes", "y", false, "Automatically install without prompting")
	updateCmd.Flags().BoolVar(&checkOnly, "check-only", false, "Only check for updates, don't install")
}
