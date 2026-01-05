package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"quickvm/ui"
)

var rootCmd = &cobra.Command{
	Use:   "quickvm",
	Short: "QuickVM - Fast Hyper-V Virtual Machine Manager",
	Long: `QuickVM is a TUI-based command-line tool for managing Hyper-V virtual machines.
It provides a fast and intuitive interface for starting, stopping, and managing VMs.`,
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
