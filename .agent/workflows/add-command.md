---
description: Add a new CLI command to the project
---

# Add New CLI Command

Follow these steps to add a new command to the QuickVM CLI:

## 1. Determine Command Type

First, identify the type of command:
- **Simple command**: `quickvm <command>` (e.g., `quickvm list`)
- **Subcommand**: `quickvm <parent> <child>` (e.g., `quickvm snapshot create`)

## 2. Create Command File

// turbo
Create the command file at `cmd/<command>.go`:

```go
package cmd

import (
    "fmt"
    "strconv"

    "quickvm/hyperv"
    "github.com/spf13/cobra"
)

var <command>Cmd = &cobra.Command{
    Use:     "<command> <args>",
    Short:   "Brief description",
    Long:    `Detailed description of what this command does.`,
    Args:    cobra.ExactArgs(1),
    Aliases: []string{"<alias>"},
    Example: `  quickvm <command> 1
  quickvm <command> 2 --flag`,
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(<command>Cmd)
}
```

## 3. Add Business Logic (if needed)

If the command requires Hyper-V operations, add method to `hyperv/<feature>.go`:

```go
func (m *Manager) <Operation>(param int) error {
    psScript := `Your-PowerShell-Command`
    output, err := m.Exec.RunCommand(psScript)
    if err != nil {
        return fmt.Errorf("operation failed: %w\nOutput: %s", err, string(output))
    }
    return nil
}
```

## 4. Add Unit Tests

// turbo
Create test file `cmd/<command>_test.go`:

```go
package cmd

import (
    "testing"
)

func Test<Command>_ValidArgs(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {"valid index", []string{"1"}, false},
        {"invalid index", []string{"abc"}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 5. Update Documentation

Update `README.md` to include the new command in:
- Usage section
- Architecture section (if new file added)

## 6. Verify

// turbo
Run tests and build:
```powershell
go test ./cmd/...
go build -o quickvm.exe
./quickvm.exe <command> --help
```

## Checklist

- [ ] Command file created in `cmd/`
- [ ] Business logic added to `hyperv/` (if needed)
- [ ] Unit tests written
- [ ] README.md updated
- [ ] Command works with `--help`
- [ ] All tests pass
