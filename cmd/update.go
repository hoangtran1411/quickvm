package cmd

import (
	"fmt"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Checking for updates...")
		
		u := updater.NewUpdater(Version)
		
		release, hasUpdate, err := u.CheckForUpdates()
		if err != nil {
			fmt.Printf("❌ Failed to check for updates: %v\n", err)
			fmt.Println("💡 Tip: Check your internet connection and try again")
			os.Exit(1)
		}

		if !hasUpdate {
			fmt.Println("✅ You are already using the latest version!")
			fmt.Printf("   Current version: %s\n", Version)
			return
		}

		fmt.Printf("🎉 New version available: %s\n", release.TagName)
		fmt.Printf("   Current version: %s\n", Version)
		fmt.Println()

		if checkOnly {
			fmt.Println("📋 Release Notes:")
			fmt.Println(release.Body)
			fmt.Println()
			fmt.Println("💡 Run 'quickvm update' without --check-only to install")
			return
		}

		if !autoInstall {
			fmt.Print("❓ Do you want to install this update? [y/N]: ")
			var response string
			fmt.Scanln(&response)
			
			if response != "y" && response != "Y" {
				fmt.Println("⏭️  Update cancelled")
				return
			}
		}

		fmt.Println()
		if err := u.DownloadAndInstall(release); err != nil {
			fmt.Printf("❌ Update failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	
	updateCmd.Flags().BoolVarP(&autoInstall, "yes", "y", false, "Automatically install without prompting")
	updateCmd.Flags().BoolVar(&checkOnly, "check-only", false, "Only check for updates, don't install")
}
