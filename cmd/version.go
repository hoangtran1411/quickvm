package cmd

import (
	"fmt"

	"quickvm/internal/output"

	"github.com/spf13/cobra"
)

var (
	// Version of the application
	Version = "1.0.0"
	// BuildDate is the date when the binary was built
	BuildDate = "2026-01-05"
	// GitCommit is the commit hash of the build
	GitCommit = "dev"
)

// VersionInfo represents version information for JSON output
type VersionInfo struct {
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
	GitCommit string `json:"gitCommit"`
	Name      string `json:"name"`
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version, build date, and git commit of QuickVM.`,
	Run: func(_ *cobra.Command, _ []string) {
		// JSON output for AI agents
		if output.IsJSON() {
			output.PrintData(VersionInfo{
				Name:      "QuickVM",
				Version:   Version,
				BuildDate: BuildDate,
				GitCommit: GitCommit,
			})
			return
		}

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
