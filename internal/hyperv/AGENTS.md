# AGENTS.md - Hyper-V Integration (internal/hyperv/)

This package contains all Hyper-V business logic, abstracting PowerShell commands behind Go interfaces.

## Architecture

```
internal/hyperv/
├── hyperv.go       # Manager struct, VM struct, core operations
├── executor.go     # ShellExecutor interface + PowerShellRunner
├── snapshot.go     # Snapshot operations
├── clone.go        # VM cloning
├── export.go       # VM export
├── import.go       # VM import
├── gpu.go          # GPU passthrough (GPU-P)
├── network.go      # Network operations
└── *_test.go       # Unit tests with mocks
```

## Core Types

```go
// VM represents a Hyper-V virtual machine
type VM struct {
    Name         string
    State        string    // Running, Off, Saved, Paused
    CPUUsage     int
    MemoryMB     int64
    Uptime       string
    Status       string
    Generation   int
    Path         string
    NetworkAdapters []NetworkAdapter
}

// Manager handles all Hyper-V operations
type Manager struct {
    Exec ShellExecutor
}

// ShellExecutor abstracts PowerShell execution for testing
type ShellExecutor interface {
    RunCommand(ctx context.Context, script string) ([]byte, error)
}
```

## PowerShell Execution Pattern

```go
func (m *Manager) GetVMs(ctx context.Context) ([]VM, error) {
    script := `Get-VM | Select-Object Name, State, CPUUsage, 
        @{N='MemoryMB';E={$_.MemoryAssigned/1MB}} | ConvertTo-Json`
    
    output, err := m.Exec.RunCommand(ctx, script)
    if err != nil {
        return nil, fmt.Errorf("failed to get VMs: %w", err)
    }
    
    var vms []VM
    if err := json.Unmarshal(output, &vms); err != nil {
        return nil, fmt.Errorf("failed to parse VM list: %w", err)
    }
    
    return vms, nil
}
```

## Security Rules (CRITICAL)

### ❌ NEVER DO THIS
```go
// Command injection vulnerability!
script := fmt.Sprintf("Get-VM -Name %s", userInput)
```

### ✅ ALWAYS DO THIS
```go
// Use parameterized execution
script := `Get-VM -Name $args[0]`
output, err := m.Exec.RunCommandWithArgs(ctx, script, vmName)
```

Or use PowerShell's built-in escaping:
```go
// Quote and escape the name
safeName := strings.ReplaceAll(vmName, "'", "''")
script := fmt.Sprintf("Get-VM -Name '%s'", safeName)
```

## Testing Pattern

All tests use mock executors to avoid requiring Windows/Hyper-V:

```go
type MockExecutor struct {
    Response []byte
    Err      error
}

func (m *MockExecutor) RunCommand(ctx context.Context, script string) ([]byte, error) {
    return m.Response, m.Err
}

func TestGetVMs(t *testing.T) {
    mock := &MockExecutor{
        Response: []byte(`[{"Name":"TestVM","State":"Running"}]`),
    }
    manager := &Manager{Exec: mock}
    
    vms, err := manager.GetVMs(context.Background())
    
    assert.NoError(t, err)
    assert.Len(t, vms, 1)
    assert.Equal(t, "TestVM", vms[0].Name)
}
```

## Adding New Operations

1. Add method signature to `Manager` struct
2. Write PowerShell script (test in PowerShell ISE first)
3. Implement method with proper error wrapping
4. Add unit test with mock executor
5. Add integration test with `//go:build windows` tag

Example:
```go
// hyperv.go
func (m *Manager) SetVMMemory(ctx context.Context, vmName string, memoryMB int64) error {
    script := fmt.Sprintf(`Set-VMMemory -VMName '%s' -DynamicMemoryEnabled $false -StartupBytes %dMB`,
        strings.ReplaceAll(vmName, "'", "''"),
        memoryMB)
    
    _, err := m.Exec.RunCommand(ctx, script)
    if err != nil {
        return fmt.Errorf("failed to set VM memory: %w", err)
    }
    return nil
}
```

## Common PowerShell Commands

| Operation | PowerShell Command |
|-----------|-------------------|
| List VMs | `Get-VM \| ConvertTo-Json` |
| Start VM | `Start-VM -Name '<name>'` |
| Stop VM | `Stop-VM -Name '<name>' -Force` |
| Get snapshots | `Get-VMSnapshot -VMName '<name>' \| ConvertTo-Json` |
| Create snapshot | `Checkpoint-VM -Name '<name>' -SnapshotName '<snap>'` |
| Clone VM | `Export-VM` + `Import-VM -Copy` |
| GPU info | `Get-VMPartitionableGpu \| ConvertTo-Json` |

## Error Handling

Always wrap errors with context:
```go
if err != nil {
    return nil, fmt.Errorf("operation '%s' failed for VM '%s': %w", 
        operation, vmName, err)
}
```

Use typed errors for specific failure modes:
```go
var (
    ErrVMNotFound    = errors.New("VM not found")
    ErrVMNotRunning  = errors.New("VM is not running")
    ErrInvalidIndex  = errors.New("invalid VM index")
)
```

## Context Usage

All methods MUST accept `context.Context`:
```go
func (m *Manager) StartVM(ctx context.Context, index int) error {
    // Use context for timeouts
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    // Pass to executor
    _, err := m.Exec.RunCommand(ctx, script)
    ...
}
```

## Build Tags

```go
//go:build windows

// Integration tests that require real Hyper-V
func TestIntegration_GetVMs(t *testing.T) {
    // This only runs on Windows with Hyper-V
}
```

## Dependencies

- No external dependencies for core logic
- `encoding/json` for PowerShell JSON parsing
- `context` for timeout/cancellation
- `os/exec` in `PowerShellRunner` only
