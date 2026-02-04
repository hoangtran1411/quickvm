# Code Snippets Reference

Quick copy-paste patterns extracted from the actual quickvm codebase.

## CLI Command Template

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
    Long:  `Detailed description`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        ctx := context.Background()
        manager := hyperv.NewManager()
        
        index, err := parseVMIndex(args[0])
        if err != nil {
            fmt.Printf("❌ %v\n", err)
            return
        }
        
        // Your logic here
        fmt.Printf("✅ Done!\n")
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

## Hyper-V Operation Template

```go
func (m *Manager) MyOperation(ctx context.Context, vmName string) error {
    // Escape single quotes in VM name
    safeName := strings.ReplaceAll(vmName, "'", "''")
    
    script := fmt.Sprintf(`
        $ErrorActionPreference = 'Stop'
        $vm = Get-VM -Name '%s'
        if (-not $vm) {
            throw "VM not found"
        }
        # Your PowerShell here
    `, safeName)
    
    output, err := m.Exec.RunCommand(ctx, script)
    if err != nil {
        return fmt.Errorf("failed to execute operation: %w", err)
    }
    
    // Parse output if needed
    _ = output
    
    return nil
}
```

## Table-Driven Test Template

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    "expected",
            wantErr: false,
        },
        {
            name:    "empty input",
            input:   "",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MyFunction(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

## Mock Executor Template

```go
type MockExecutor struct {
    Response []byte
    Err      error
    Called   bool
    Script   string
}

func (m *MockExecutor) RunCommand(ctx context.Context, script string) ([]byte, error) {
    m.Called = true
    m.Script = script
    return m.Response, m.Err
}

// Usage
func TestWithMock(t *testing.T) {
    mock := &MockExecutor{
        Response: []byte(`{"Name": "TestVM", "State": "Running"}`),
    }
    manager := &Manager{Exec: mock}
    
    result, err := manager.GetVM(context.Background(), "TestVM")
    
    assert.NoError(t, err)
    assert.True(t, mock.Called)
    assert.Contains(t, mock.Script, "TestVM")
}
```

## Bubble Tea Command Template

```go
// Message type
type operationResultMsg struct {
    success bool
    message string
    err     error
}

// Command that returns message
func (m model) performOperation() tea.Cmd {
    return func() tea.Msg {
        ctx := context.Background()
        vm := m.getSelectedVM()
        
        err := m.manager.SomeOperation(ctx, vm.Name)
        if err != nil {
            return operationResultMsg{
                success: false,
                err:     err,
            }
        }
        
        return operationResultMsg{
            success: true,
            message: fmt.Sprintf("Operation completed for %s", vm.Name),
        }
    }
}

// Handle in Update()
case operationResultMsg:
    if msg.err != nil {
        m.message = fmt.Sprintf("❌ Error: %v", msg.err)
    } else {
        m.message = fmt.Sprintf("✅ %s", msg.message)
    }
    return m, m.refreshVMs()
```

## Lipgloss Style Template

```go
var (
    // Colors
    primaryColor = lipgloss.Color("86")   // Cyan
    successColor = lipgloss.Color("82")   // Green
    errorColor   = lipgloss.Color("196")  // Red
    mutedColor   = lipgloss.Color("240")  // Gray
    
    // Styles
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(primaryColor)
    
    successStyle = lipgloss.NewStyle().
        Foreground(successColor)
    
    errorStyle = lipgloss.NewStyle().
        Foreground(errorColor)
    
    boxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(primaryColor).
        Padding(1, 2)
)
```

## Error Wrapping Pattern

```go
// At package boundary
if err != nil {
    return fmt.Errorf("hyperv: failed to start VM '%s': %w", vmName, err)
}

// Check specific errors
if errors.Is(err, ErrVMNotFound) {
    // Handle not found
}

// Custom error types
var ErrVMNotFound = errors.New("VM not found")
var ErrVMNotRunning = errors.New("VM is not running")
```

## Context with Timeout

```go
func (m *Manager) LongOperation(ctx context.Context) error {
    // Add timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
    defer cancel()
    
    // Check context before long operations
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    // Pass context to executor
    output, err := m.Exec.RunCommand(ctx, script)
    // ...
}
```

---

# Real Examples from Codebase

The following are actual patterns extracted from quickvm source code.

## Real: ShellExecutor Interface (internal/hyperv/hyperv.go)

```go
// ShellExecutor defines an interface for executing shell commands with Context
type ShellExecutor interface {
    // RunScript executes a complex PowerShell script (beware of injection, use for static scripts)
    RunScript(ctx context.Context, script string) ([]byte, error)
    // RunCmdlet executes a specific cmdlet with arguments safely
    RunCmdlet(ctx context.Context, cmdlet string, args ...string) ([]byte, error)
}

// PowerShellRunner implements ShellExecutor for actual PowerShell execution
type PowerShellRunner struct{}

// RunScript executes a PowerShell command/script with Context
func (p *PowerShellRunner) RunScript(ctx context.Context, script string) ([]byte, error) {
    cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", script)
    out, err := cmd.CombinedOutput()
    if err != nil {
        return out, fmt.Errorf("execution failed: %w", err)
    }
    return out, nil
}

// RunCmdlet executes a PowerShell cmdlet safely using separate arguments
func (p *PowerShellRunner) RunCmdlet(ctx context.Context, cmdlet string, args ...string) ([]byte, error) {
    psArgs := make([]string, 0, 4+len(args))
    psArgs = append(psArgs, "-NoProfile", "-NonInteractive", "-Command", cmdlet)
    psArgs = append(psArgs, args...)

    cmd := exec.CommandContext(ctx, "powershell", psArgs...)
    out, err := cmd.CombinedOutput()
    if err != nil {
        return out, fmt.Errorf("cmdlet execution failed: %w", err)
    }
    return out, nil
}
```

## Real: Safe VM Operation (internal/hyperv/hyperv.go)

```go
// StartVMByName starts a virtual machine by name
func (m *Manager) StartVMByName(ctx context.Context, name string) error {
    // why: Using RunCmdlet with separate args prevents shell injection attacks
    // where 'name' could contain malicious PowerShell commands.
    output, err := m.Exec.RunCmdlet(ctx, "Start-VM", "-Name", name)
    if err != nil {
        return fmt.Errorf("failed to start VM '%s': %v\nOutput: %s", name, err, string(output))
    }
    return nil
}
```

## Real: CLI Command with Bulk Operations (cmd/start.go)

```go
var (
    startRange string
    startAll   bool
)

var startCmd = &cobra.Command{
    Use:   "start [vm-index]",
    Short: "Start a Hyper-V virtual machine",
    Long: `Start a Hyper-V virtual machine by its index.

Examples:
  quickvm start 1 3 5       # Start VMs at index 1, 3, and 5
  quickvm start --range 1-5 # Start VMs from index 1 to 5
  quickvm start --all       # Start all VMs`,
    Args: cobra.ArbitraryArgs,
    Run: func(cmd *cobra.Command, args []string) {
        runStart(cmd.Context(), hyperv.NewManager(), args, startRange, startAll)
    },
}

func runStart(ctx context.Context, manager hyperv.VMManager, args []string, rangeStr string, all bool) {
    vms, err := manager.GetVMs(ctx)
    if err != nil {
        fmt.Printf("❌ Failed to get VMs: %v\n", err)
        return
    }

    indices, err := getIndices(args, rangeStr, all, len(vms))
    if err != nil {
        fmt.Printf("❌ Error: %v\n", err)
        return
    }

    if len(indices) > 1 {
        fmt.Printf("🚀 Starting %d VMs...\n\n", len(indices))
    }

    successCount := 0
    failCount := 0

    for _, index := range indices {
        vm := vms[index-1]
        fmt.Printf("🚀 Starting VM: %s (Index: %d)...\n", vm.Name, index)

        if err := manager.StartVMByName(ctx, vm.Name); err != nil {
            fmt.Printf("❌ Failed to start VM '%s': %v\n", vm.Name, err)
            failCount++
        } else {
            fmt.Printf("✅ VM '%s' started successfully!\n", vm.Name)
            successCount++
        }
    }

    if len(indices) > 1 {
        fmt.Printf("\n📊 Summary: %d started, %d failed\n", successCount, failCount)
    }
}

func init() {
    startCmd.Flags().StringVarP(&startRange, "range", "r", "", "Range of VM indices")
    startCmd.Flags().BoolVarP(&startAll, "all", "a", false, "Start all virtual machines")
    rootCmd.AddCommand(startCmd)
}
```

## Real: MockManager for Testing (cmd/mock_test.go)

```go
// MockManager is a mock implementation of VMManager for testing
type MockManager struct {
    GetVMsFn          func(ctx context.Context) ([]hyperv.VM, error)
    StartVMFn         func(ctx context.Context, index int) error
    StartVMByNameFn   func(ctx context.Context, name string) error
    StopVMFn          func(ctx context.Context, index int) error
    StopVMByNameFn    func(ctx context.Context, name string) error
    RestartVMFn       func(ctx context.Context, index int) error
    RestartVMByNameFn func(ctx context.Context, name string) error
    GetVMStatusFn     func(ctx context.Context, name string) (string, error)
}

func (m *MockManager) GetVMs(ctx context.Context) ([]hyperv.VM, error) {
    if m.GetVMsFn != nil {
        return m.GetVMsFn(ctx)
    }
    return []hyperv.VM{}, nil
}

func (m *MockManager) StartVMByName(ctx context.Context, name string) error {
    if m.StartVMByNameFn != nil {
        return m.StartVMByNameFn(ctx, name)
    }
    return nil
}

// ... other methods follow same pattern
```

## Real: Handling PowerShell JSON Quirks (internal/hyperv/hyperv.go)

```go
// GetVMs retrieves all Hyper-V virtual machines
func (m *Manager) GetVMs(ctx context.Context) ([]VM, error) {
    psScript := `
        Get-VM | Select-Object @{Name='Name';Expression={$_.Name.ToString()}}, 
        @{Name='State';Expression={$_.State.ToString()}}, 
        @{Name='CPUUsage';Expression={[int]$_.CPUUsage}}, 
        @{Name='MemoryMB';Expression={[int]($_.MemoryAssigned/1MB)}},
        @{Name='Uptime';Expression={$_.Uptime.ToString()}},
        @{Name='Status';Expression={$_.Status.ToString()}},
        @{Name='Version';Expression={$_.Version.ToString()}},
        @{Name='IPAddresses';Expression={($_.NetworkAdapters.IPAddresses | Where-Object { $_ -match '^\d+\.\d+\.\d+\.\d+$' })}} | ConvertTo-Json
    `

    output, err := m.Exec.RunScript(ctx, psScript)
    if err != nil {
        return nil, fmt.Errorf("failed to execute PowerShell command: %v\nOutput: %s", err, string(output))
    }

    var vms []VM
    outputStr := strings.TrimSpace(string(output))

    // Handle single VM case (PowerShell returns object, not array)
    switch {
    case strings.HasPrefix(outputStr, "{"):
        // Single VM - parse as object
        var vmRaw struct {
            Name   string `json:"name"`
            State  string `json:"state"`
            // ... other fields
        }
        if err := json.Unmarshal(output, &vmRaw); err != nil {
            return nil, fmt.Errorf("failed to parse VM data: %v", err)
        }
        vms = append(vms, VM{Name: vmRaw.Name, State: vmRaw.State})
        
    case strings.HasPrefix(outputStr, "["):
        // Multiple VMs - parse as array
        if err := json.Unmarshal(output, &vms); err != nil {
            return nil, fmt.Errorf("failed to parse VMs data: %v", err)
        }
        
    default:
        // Empty output = no VMs
        if outputStr == "" {
            return []VM{}, nil
        }
        return nil, fmt.Errorf("invalid output format")
    }

    // Assign 1-based indices
    for i := range vms {
        vms[i].Index = i + 1
    }

    return vms, nil
}
```

## Real: VM Struct Definition

```go
// VM represents a Hyper-V virtual machine
type VM struct {
    Index       int      `json:"-"`              // 1-based index, not in JSON
    Name        string   `json:"name"`
    State       string   `json:"state"`          // Running, Off, Saved, Paused
    CPUUsage    int      `json:"cpuUsage"`
    MemoryMB    int64    `json:"memoryMB"`
    Uptime      string   `json:"uptime"`
    Status      string   `json:"status"`
    Version     string   `json:"version"`
    IPAddresses []string `json:"ipAddresses"`
}
```

## Real: VMManager Interface

```go
// VMManager defines the interface for Hyper-V operations to allow mocking in tests
type VMManager interface {
    GetVMs(ctx context.Context) ([]VM, error)
    StartVM(ctx context.Context, index int) error
    StartVMByName(ctx context.Context, name string) error
    StopVM(ctx context.Context, index int) error
    StopVMByName(ctx context.Context, name string) error
    RestartVM(ctx context.Context, index int) error
    RestartVMByName(ctx context.Context, name string) error
    GetVMStatus(ctx context.Context, name string) (string, error)
}
```
