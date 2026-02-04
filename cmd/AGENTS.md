# AGENTS.md - CLI Commands (cmd/)

This directory contains all CLI commands built with [Cobra](https://github.com/spf13/cobra).

## Architecture

```
cmd/
├── root.go         # Root command, launches TUI when no args
├── start.go        # quickvm start <index>
├── stop.go         # quickvm stop <index>
├── restart.go      # quickvm restart <index>
├── list.go         # quickvm list
├── info.go         # quickvm info <index>
├── clone.go        # quickvm clone <index> <name>
├── export.go       # quickvm export <index> <path>
├── import.go       # quickvm import <path>
├── snapshot.go     # quickvm snapshot <subcommand>
├── gpu.go          # quickvm gpu <subcommand>
├── rdp.go          # quickvm rdp <index>
├── workspace.go    # quickvm workspace <subcommand>
├── enable.go       # quickvm enable (Hyper-V feature)
├── update.go       # quickvm update
├── version.go      # quickvm version
└── utils.go        # Shared utilities
```

## Command Pattern

All commands follow this structure:

```go
package cmd

import (
    "context"
    "fmt"
    "quickvm/internal/hyperv"
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand <index>",
    Short: "Brief description",
    Long:  `Detailed description with examples`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        ctx := context.Background()
        manager := hyperv.NewManager()
        
        index, err := parseVMIndex(args[0])
        if err != nil {
            fmt.Printf("❌ %v\n", err)
            return
        }
        
        if err := manager.MyOperation(ctx, index); err != nil {
            fmt.Printf("❌ Failed: %v\n", err)
            return
        }
        
        fmt.Printf("✅ Success!\n")
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
    // Add flags if needed
    myCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force operation")
}
```

## Key Conventions

### Index-Based VM Reference
Commands use 1-based numeric indices instead of VM names:
```powershell
quickvm start 1    # Start first VM
quickvm stop 2     # Stop second VM
```

### Error Output Format
Use emoji prefixes for visual feedback:
- `✅` Success
- `❌` Error
- `⚠️` Warning
- `🔄` In progress

### Shared Utilities (utils.go)

```go
// parseVMIndex - Convert string arg to valid VM index
func parseVMIndex(arg string) (int, error)

// getVMByIndex - Fetch VM list and return VM at index
func getVMByIndex(ctx context.Context, m *hyperv.Manager, index int) (*hyperv.VM, error)
```

## Adding a New Command

1. Create `cmd/<command>.go`
2. Define `var <name>Cmd = &cobra.Command{...}`
3. Add `rootCmd.AddCommand(<name>Cmd)` in `init()`
4. Use `parseVMIndex()` for index arguments
5. Always use context for Hyper-V operations
6. Add tests in `cmd/<command>_test.go`

Or use workflow: `/add-command`

## Testing Commands

```go
func TestMyCommand(t *testing.T) {
    // Commands are tested via integration with hyperv package
    // Mock the ShellExecutor interface for unit tests
}
```

## Subcommand Groups

Some commands have subcommands:

```go
// snapshot.go
var snapshotCmd = &cobra.Command{Use: "snapshot", Short: "Manage snapshots"}
var snapshotCreateCmd = &cobra.Command{Use: "create <index> <name>", ...}
var snapshotListCmd = &cobra.Command{Use: "list <index>", ...}
var snapshotRestoreCmd = &cobra.Command{Use: "restore <index> <name>", ...}

func init() {
    rootCmd.AddCommand(snapshotCmd)
    snapshotCmd.AddCommand(snapshotCreateCmd, snapshotListCmd, snapshotRestoreCmd)
}
```

## Current Commands Reference

| Command | Usage | Description |
|---------|-------|-------------|
| `list` | `quickvm list` | List all VMs with status |
| `start` | `quickvm start <index>` | Start a VM |
| `stop` | `quickvm stop <index>` | Stop a VM |
| `restart` | `quickvm restart <index>` | Restart a VM |
| `info` | `quickvm info <index>` | Show VM details |
| `clone` | `quickvm clone <index> <name>` | Clone a VM |
| `export` | `quickvm export <index> <path>` | Export VM |
| `import` | `quickvm import <path>` | Import VM |
| `rdp` | `quickvm rdp <index>` | Connect via RDP |
| `snapshot` | `quickvm snapshot <sub>` | Manage snapshots |
| `gpu` | `quickvm gpu <sub>` | GPU passthrough |
| `workspace` | `quickvm workspace <sub>` | Workspace profiles |
| `enable` | `quickvm enable` | Enable Hyper-V |
| `update` | `quickvm update` | Auto-update |
| `version` | `quickvm version` | Show version |
