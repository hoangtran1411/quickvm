---
name: Go CLI with Cobra
description: Build production-ready CLI applications in Go using Cobra framework with professional command structure, help generation, and flag handling.
---

# Go CLI with Cobra

This skill provides patterns and templates for building professional CLI applications in Go using [spf13/cobra](https://github.com/spf13/cobra).

## When to Use

- Creating new CLI applications in Go
- Adding commands/subcommands to existing CLI apps  
- Implementing flag handling and validation
- Setting up command hierarchies (e.g., `app resource action`)

## Architecture Pattern

```
cmd/
├── root.go        # Root command, global flags, initialization
├── version.go     # Version information
├── <resource>.go  # Resource-specific commands (e.g., vm.go, user.go)
└── utils.go       # Shared utilities for commands
```

## Core Templates

### 1. Root Command (`cmd/root.go`)

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Global flags
var (
	verbose bool
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "appname",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and explains
your application in detail.`,
	// Uncomment to run TUI when no subcommand is provided
	// Run: func(cmd *cobra.Command, args []string) {
	//     runTUI()
	// },
}

// Execute adds all child commands to the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.appname.yaml)")
}
```

### 2. Resource Command Template (`cmd/<resource>.go`)

```go
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var resourceCmd = &cobra.Command{
	Use:     "resource <index>",
	Short:   "Perform action on a resource",
	Long:    `Detailed description of what this command does.`,
	Args:    cobra.ExactArgs(1), // Require exactly one argument
	Aliases: []string{"rs", "r"}, // Short aliases
	Example: `  appname resource 1
  appname resource 2 --force`,
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("❌ Error: Invalid index. Please provide a number.")
			return
		}

		// Your business logic here
		if err := performAction(index); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Printf("✅ Successfully performed action on resource %d\n", index)
	},
}

// Local flags
var forceFlag bool

func init() {
	resourceCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force the action without confirmation")
	rootCmd.AddCommand(resourceCmd)
}

func performAction(index int) error {
	// Business logic implementation
	return nil
}
```

### 3. Version Command (`cmd/version.go`)

```go
package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// These variables are set during build via -ldflags
var (
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("App Name v%s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
```

### 4. Subcommand Pattern (e.g., `app snapshot create`)

```go
package cmd

import (
	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage snapshots",
	Long:  `Create, list, restore, and delete snapshots.`,
}

var snapshotListCmd = &cobra.Command{
	Use:   "list <vm-index>",
	Short: "List all snapshots for a VM",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation
	},
}

var snapshotCreateCmd = &cobra.Command{
	Use:   "create <vm-index> <name>",
	Short: "Create a new snapshot",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation
	},
}

func init() {
	snapshotCmd.AddCommand(snapshotListCmd)
	snapshotCmd.AddCommand(snapshotCreateCmd)
	rootCmd.AddCommand(snapshotCmd)
}
```

## Best Practices

1. **Use emojis for visual feedback**:
   - ✅ Success
   - ❌ Error  
   - ⏳ In progress
   - ⚠️ Warning

2. **Consistent error handling**:
   ```go
   if err != nil {
       fmt.Printf("❌ Error: %v\n", err)
       os.Exit(1)
   }
   ```

3. **Validate arguments early**:
   ```go
   Args: cobra.MinimumNArgs(1),
   // or
   Args: cobra.ExactArgs(2),
   ```

4. **Use `RunE` for commands that can fail**:
   ```go
   RunE: func(cmd *cobra.Command, args []string) error {
       return someAction()
   },
   ```

5. **Provide helpful examples** in the `Example` field

6. **Use aliases** for frequently used commands

## Build with Version Info

```bash
go build -ldflags="-X 'cmd.Version=1.0.0' -X 'cmd.BuildDate=$(date -I)' -X 'cmd.GitCommit=$(git rev-parse --short HEAD)'" -o app
```

## Entry Point (`main.go`)

```go
package main

import "yourapp/cmd"

func main() {
	cmd.Execute()
}
```
