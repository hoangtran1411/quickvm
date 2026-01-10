# QuickVM Developer Notes

## Quick Start for Developers

### Prerequisites
- Go 1.21+
- Windows 10/11 with Hyper-V
- PowerShell 5.1+
- Git

### First Time Setup

```powershell
# Clone the repository
git clone <your-repo-url>
cd quickvm

# Install dependencies
go mod download
go mod tidy

# Build
go build -o quickvm.exe

# Test
go test ./...

# Run
.\quickvm.exe list
```

## Architecture Overview

### Component Breakdown

```
┌─────────────────────────────────────────────────────────┐
│                    User Interface                       │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  CLI (Cobra) ────────────┐                              │
│                          │                              │
│                          ▼                              │
│                   ┌─────────────┐                       │
│                   │   Commands  │                       │
│                   │  (cmd/*)    │                       │
│                   └──────┬──────┘                       │
│                          │                              │
│                          ▼                              │
│         ┌────────────────────────────────┐              │
│         │                                │              │
│    ┌────▼────┐                    ┌─────▼─────┐         │
│    │   TUI   │                    │  Hyperv   │         │
│    │ (ui/*) │                    │ (hyperv/*) │         │
│    │         │                    │            │        │
│    │ Bubble  │                    │ PowerShell │        │
│    │  Tea    │                    │ Commands   │        │
│    └─────────┘                    └─────┬──────┘        │
│                                         │               │
│                                         ▼               │
│                                  ┌─────────────┐        │
│                                  │  Hyper-V    │        │
│                                  │   Engine    │        │
│                                  └─────────────┘        │
└─────────────────────────────────────────────────────────┘
```

### Package Responsibilities

#### `cmd/` - CLI Commands Layer
- **Purpose**: Define CLI commands and their behavior
- **Technology**: Cobra framework
- **Files**:
  - `root.go` - Root command, launches TUI if no args
  - `start.go` - Start VM command
  - `stop.go` - Stop VM command
  - `restart.go` - Restart VM command
  - `list.go` - List VMs command
  - `version.go` - Version information

#### `hyperv/` - Hyper-V Integration Layer
- **Purpose**: Interact with Hyper-V through PowerShell
- **Key Functions**:
  - `GetVMs()` - Fetch all VMs and their status
  - `StartVM(index)` - Start a VM by index
  - `StopVM(index)` - Stop a VM by index
  - `RestartVM(index)` - Restart a VM by index
- **Technology**: PowerShell execution via `os/exec`

#### `ui/` - TUI Layer
- **Purpose**: Interactive terminal UI
- **Technology**: Bubble Tea + Bubbles + Lipgloss
- **Components**:
  - Table model for VM list
  - Keyboard navigation
  - Status messages
  - Color-coded states

## Key Design Decisions

### 1. PowerShell Integration
**Decision**: Use PowerShell commands instead of Windows API

**Rationale**:
- ✅ Simpler to implement
- ✅ Easier to maintain
- ✅ More reliable (uses official Hyper-V cmdlets)
- ✅ Better error messages
- ❌ Slightly slower (acceptable for this use case)

**Implementation**:
```go
cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
output, err := cmd.CombinedOutput()
```

### 2. Index-Based VM Reference
**Decision**: Use numeric indices instead of names in CLI

**Rationale**:
- ✅ Faster to type (e.g., `quickvm start 1` vs `quickvm start "My Long VM Name"`)
- ✅ No need to escape special characters
- ✅ Consistent with TUI navigation
- ❌ Indices can change if VMs are added/removed

**Mitigation**: Always show VM name in output for confirmation

### 3. Dual Interface (CLI + TUI)
**Decision**: Provide both command-line and interactive modes

**Rationale**:
- CLI: Fast, scriptable, automation-friendly
- TUI: Visual, exploratory, beginner-friendly

**Usage Patterns**:
- **Automation**: Use CLI commands in scripts
- **Exploration**: Use TUI to browse and manage VMs
- **Quick Actions**: Use CLI for single operations

## Common Development Tasks

### Adding a New Command

1. Create `cmd/mycommand.go`:
```go
package cmd

import (
    "github.com/spf13/cobra"
    "quickvm/hyperv"
)

var myCmd = &cobra.Command{
    Use:   "mycommand [args]",
    Short: "Short description",
    Long:  `Detailed description`,
    Run: func(cmd *cobra.Command, args []string) {
        // Your implementation
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

2. Rebuild and test:
```powershell
go build -o quickvm.exe
.\quickvm.exe mycommand
```

### Adding Hyper-V Functionality

1. Add method to `hyperv/hyperv.go`:
```go
func (m *Manager) MyNewFunction(vmIndex int) error {
    psScript := `Your-PowerShell-Command`
    cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("operation failed: %v\nOutput: %s", err, string(output))
    }
    return nil
}
```

2. Add test in `hyperv/hyperv_test.go`

### Modifying TUI

1. Update model in `ui/table.go`:
```go
// Add to Update function for new keybindings
case "x":
    // Your new action
    return m, m.myNewAction()
```

2. Test manually:
```powershell
go build -o quickvm.exe
.\quickvm.exe  # launches TUI
```

## Testing Strategy

### Unit Tests
```powershell
# Run all tests
go test ./...

# Run specific package
go test ./hyperv

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Manual Testing Checklist

- [ ] `quickvm list` - Shows all VMs
- [ ] `quickvm start 1` - Starts first VM
- [ ] `quickvm stop 1` - Stops first VM
- [ ] `quickvm restart 1` - Restarts first VM
- [ ] `quickvm clone 1 "NewVM"` - Clones a VM
- [ ] `quickvm rdp 1` - Opens RDP connection
- [ ] `quickvm gpu status` - Shows GPU support
- [ ] `quickvm version` - Shows version info
- [ ] `quickvm` (no args) - Launches TUI
- [ ] TUI: Arrow navigation works
- [ ] TUI: Enter to start VM
- [ ] TUI: 's' to stop VM
- [ ] TUI: 't' to restart VM
- [ ] TUI: 'r' to refresh
- [ ] TUI: 'q' to quit

## Performance Considerations

### Current Bottlenecks
1. **PowerShell startup time** (~200-500ms per command)
   - Acceptable for interactive use
   - Could be improved with persistent PowerShell session

2. **JSON parsing** - Minimal overhead

### Optimization Opportunities
1. Cache VM list for short periods (5-10 seconds)
2. Use PowerShell Runspaces for faster execution
3. Batch multiple operations

## Error Handling Patterns

### PowerShell Errors
```go
cmd := exec.Command("powershell", ...)
output, err := cmd.CombinedOutput()
if err != nil {
    return fmt.Errorf("operation failed: %v\nOutput: %s", err, string(output))
}
```

### Validation Errors
```go
if index < 1 || index > len(vms) {
    return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", index, len(vms))
}
```

### User-Friendly Messages
- ✅ Use emojis for visual feedback
- ✅ Include context in error messages
- ✅ Suggest fixes when possible

## Build & Release

### Development Build
```powershell
go build -o quickvm.exe
```

### Optimized Build
```powershell
go build -ldflags="-s -w" -o quickvm.exe
```

### Multi-Architecture Build
```powershell
# AMD64
$env:GOOS="windows"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -o quickvm-amd64.exe

# ARM64
$env:GOOS="windows"; $env:GOARCH="arm64"
go build -ldflags="-s -w" -o quickvm-arm64.exe
```

### Version Management
Update `cmd/version.go`:
```go
var (
    Version   = "1.1.0"
    BuildDate = "2026-01-05"
    GitCommit = "abc1234"
)
```

## Troubleshooting Development Issues

### Import Errors
```powershell
go mod tidy
go mod download
```

### Build Cache Issues
```powershell
go clean -cache
go build -v -o quickvm.exe
```

### PowerShell Script Debugging
Create test scripts in separate `.ps1` files:
```powershell
# test-script.ps1
Get-VM | Select-Object Name, State | ConvertTo-Json

# Run directly
powershell -File test-script.ps1
```

## Future Improvements

### High Priority
- [ ] Bulk Operations Enhancement (Multi-index, --all)
- [ ] Workspace/Profile System (.quickvm/workspaces/*.yaml)
- [ ] VM Config (RAM/CPU management)
- [ ] Better error messages

### Medium Priority
- [ ] Remote Hyper-V server support
- [ ] VM grouping/tagging
- [x] Export/import VM configs ✅
- [x] VM Clone (Full Clone) ✅
- [x] GPU Partitioning (GPU-P) ✅
- [ ] Performance metrics

### Low Priority
- [ ] GUI wrapper
- [ ] Web interface
- [ ] Mobile app integration
- [ ] VM templates

## Resources

### Documentation
- [Cobra Framework](https://github.com/spf13/cobra)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Hyper-V PowerShell](https://docs.microsoft.com/en-us/powershell/module/hyper-v/)
- [Go Documentation](https://golang.org/doc/)

### Similar Projects
- Vagrant (multi-VM management)
- Multipass (Ubuntu VMs)
- Docker (containers)

### Community
- GitHub Discussions
- Stack Overflow
- Reddit r/golang

## Contact & Support

- **Issues**: GitHub Issues
- **Discussions**: GitHub Discussions
- **Email**: [your-email]

---

**Happy Coding! 🚀**
