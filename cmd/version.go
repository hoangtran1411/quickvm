package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version information
	Version   = "1.0.0"
	BuildDate = "2026-01-05"
	GitCommit = "dev"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version, build date, and git commit of QuickVM.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("QuickVM - Fast Hyper-V Virtual Machine Manager\n")
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("\nMade with ❤️  using Go\n")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
