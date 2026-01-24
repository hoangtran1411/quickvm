package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"quickvm/ui"
	"quickvm/updater"
)

var autoUpdate bool

var rootCmd = &cobra.Command{
	Use:   "quickvm",
	Short: "QuickVM - Fast Hyper-V Virtual Machine Manager",
	Long: `QuickVM is a TUI-based command-line tool for managing Hyper-V virtual machines.
It provides a fast and intuitive interface for starting, stopping, and managing VMs.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Check for updates if --update flag is set
		if autoUpdate && cmd.Name() != "update" {
			checkAndUpdate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Launch TUI
		p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&autoUpdate, "update", false, "Check for updates before running")
}

// checkAndUpdate checks for updates and prompts to install if available
func checkAndUpdate() {
	u := updater.NewUpdater(Version)

	release, hasUpdate, err := u.CheckForUpdates()
	if err != nil {
		// Silently fail on update check errors
		return
	}

	if !hasUpdate {
		fmt.Println("✅ The current version is the latest!")
		fmt.Println()
		return
	}

	fmt.Printf("🎉 New version available: %s (current: %s)\n", release.TagName, Version)
	fmt.Print("❓ Do you want to update now? [Y/n]: ")

	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		// Default to 'yes' if can't read input
		response = ""
	}

	if response == "n" || response == "N" {
		fmt.Println("⏭️  Continuing with current version...")
		fmt.Println()
		return
	}

	fmt.Println()
	if err := u.DownloadAndInstall(release); err != nil {
		fmt.Printf("❌ Update failed: %v\n", err)
		fmt.Println("⏭️  Continuing with current version...")
		fmt.Println()
		return
	}

	fmt.Println("✅ Update complete! Please restart QuickVM.")
	os.Exit(0)
}
