---
name: Go PowerShell Integration
description: Safely execute PowerShell commands from Go with proper error handling, output parsing, and testability through interface-based design.
---

# Go PowerShell Integration

This skill provides patterns for safely executing PowerShell commands from Go applications, with emphasis on testability, security, and proper error handling.

## When to Use

- Windows system administration tools
- Hyper-V, Active Directory, or other Windows management
- Any automation requiring PowerShell cmdlets
- Cross-language integration with Windows

## Key Principles

1. **Interface-based design** - Enable mocking for tests
2. **Never use string concatenation** for user input - Prevent command injection
3. **Use context for timeouts** - Handle long-running operations
4. **Comprehensive error handling** - Include PowerShell output in errors

## Architecture Pattern

```
mypackage/
├── executor.go      # ShellExecutor interface + PowerShellRunner
├── manager.go       # Business logic using executor
├── manager_test.go  # Tests with mock executor
└── mock_test.go     # Mock implementations
```

## Core Templates

### 1. Executor Interface (`executor.go`)

```go
package mypackage

import (
    "context"
    "os/exec"
)

// ShellExecutor defines an interface for executing shell commands.
// This abstraction allows for easy mocking in tests.
type ShellExecutor interface {
    RunCommand(script string) ([]byte, error)
    RunCommandContext(ctx context.Context, script string) ([]byte, error)
}

// PowerShellRunner implements ShellExecutor for actual PowerShell execution
type PowerShellRunner struct{}

// RunCommand executes a PowerShell command
func (p *PowerShellRunner) RunCommand(script string) ([]byte, error) {
    cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
    return cmd.CombinedOutput()
}

// RunCommandContext executes a PowerShell command with context for timeout/cancellation
func (p *PowerShellRunner) RunCommandContext(ctx context.Context, script string) ([]byte, error) {
    cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", script)
    return cmd.CombinedOutput()
}
```

### 2. Manager with Dependency Injection (`manager.go`)

```go
package mypackage

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "time"
)

// Manager handles operations using PowerShell
type Manager struct {
    Exec ShellExecutor
}

// NewManager creates a new Manager with default PowerShell runner
func NewManager() *Manager {
    return &Manager{
        Exec: &PowerShellRunner{},
    }
}

// NewManagerWithExecutor creates a Manager with a custom executor (for testing)
func NewManagerWithExecutor(exec ShellExecutor) *Manager {
    return &Manager{
        Exec: exec,
    }
}

// Example: Get system info with JSON parsing
func (m *Manager) GetSystemInfo() (*SystemInfo, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    psScript := `
        @{
            ComputerName = $env:COMPUTERNAME
            OSVersion = [System.Environment]::OSVersion.VersionString
            ProcessorCount = $env:NUMBER_OF_PROCESSORS
        } | ConvertTo-Json
    `

    output, err := m.Exec.RunCommandContext(ctx, psScript)
    if err != nil {
        return nil, fmt.Errorf("failed to get system info: %w\nOutput: %s", err, string(output))
    }

    var info SystemInfo
    if err := json.Unmarshal(output, &info); err != nil {
        return nil, fmt.Errorf("failed to parse system info: %w", err)
    }

    return &info, nil
}

// Example: Execute action with parameter (safe from injection)
func (m *Manager) PerformAction(name string) error {
    // SECURITY: Use proper escaping, never concatenate user input directly
    escapedName := strings.ReplaceAll(name, "'", "''")
    psScript := fmt.Sprintf(`Do-Something -Name '%s'`, escapedName)

    output, err := m.Exec.RunCommand(psScript)
    if err != nil {
        return fmt.Errorf("action failed for '%s': %w\nOutput: %s", name, err, string(output))
    }

    return nil
}

type SystemInfo struct {
    ComputerName   string `json:"ComputerName"`
    OSVersion      string `json:"OSVersion"`
    ProcessorCount string `json:"ProcessorCount"`
}
```

### 3. Mock for Testing (`mock_test.go`)

```go
package mypackage

import (
    "context"
    "fmt"
)

// MockShellExecutor is a mock implementation of ShellExecutor for testing
type MockShellExecutor struct {
    Output       []byte
    Error        error
    CalledScript string
    CallCount    int
}

func (m *MockShellExecutor) RunCommand(script string) ([]byte, error) {
    m.CalledScript = script
    m.CallCount++
    return m.Output, m.Error
}

func (m *MockShellExecutor) RunCommandContext(ctx context.Context, script string) ([]byte, error) {
    return m.RunCommand(script)
}

// NewMockExecutor creates a mock executor with predefined output
func NewMockExecutor(output []byte, err error) *MockShellExecutor {
    return &MockShellExecutor{
        Output: output,
        Error:  err,
    }
}

// MockSequenceExecutor allows defining different outputs for sequential calls
type MockSequenceExecutor struct {
    Outputs   [][]byte
    Errors    []error
    CallIndex int
}

func (m *MockSequenceExecutor) RunCommand(script string) ([]byte, error) {
    if m.CallIndex >= len(m.Outputs) {
        return nil, fmt.Errorf("unexpected call: %s", script)
    }
    output := m.Outputs[m.CallIndex]
    var err error
    if m.CallIndex < len(m.Errors) {
        err = m.Errors[m.CallIndex]
    }
    m.CallIndex++
    return output, err
}

func (m *MockSequenceExecutor) RunCommandContext(ctx context.Context, script string) ([]byte, error) {
    return m.RunCommand(script)
}
```

### 4. Tests (`manager_test.go`)

```go
package mypackage

import (
    "errors"
    "testing"
)

func TestGetSystemInfo_Success(t *testing.T) {
    mockOutput := []byte(`{
        "ComputerName": "TEST-PC",
        "OSVersion": "Microsoft Windows 10.0.22000",
        "ProcessorCount": "8"
    }`)

    mock := NewMockExecutor(mockOutput, nil)
    manager := NewManagerWithExecutor(mock)

    info, err := manager.GetSystemInfo()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if info.ComputerName != "TEST-PC" {
        t.Errorf("expected ComputerName 'TEST-PC', got '%s'", info.ComputerName)
    }
}

func TestGetSystemInfo_Error(t *testing.T) {
    mock := NewMockExecutor([]byte("Access denied"), errors.New("exit status 1"))
    manager := NewManagerWithExecutor(mock)

    _, err := manager.GetSystemInfo()
    if err == nil {
        t.Fatal("expected error, got nil")
    }
}

func TestPerformAction_Escaping(t *testing.T) {
    mock := NewMockExecutor([]byte{}, nil)
    manager := NewManagerWithExecutor(mock)

    // Test with potentially dangerous input
    err := manager.PerformAction("test'; Remove-Item -Recurse /")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Verify the script was properly escaped
    expected := `Do-Something -Name 'test''; Remove-Item -Recurse /'`
    if mock.CalledScript != expected {
        t.Errorf("expected script '%s', got '%s'", expected, mock.CalledScript)
    }
}
```

## Security Best Practices

### ❌ NEVER Do This

```go
// DANGEROUS: Command injection vulnerability!
name := args[0]
psScript := fmt.Sprintf(`Get-VM -Name "%s"`, name)  // User could input: "; Remove-Item -Recurse /"
```

### ✅ Always Do This

```go
// SAFE: Escape single quotes and use single-quoted strings
escapedName := strings.ReplaceAll(name, "'", "''")
psScript := fmt.Sprintf(`Get-VM -Name '%s'`, escapedName)

// OR use validated/sanitized input
if !isValidVMName(name) {
    return fmt.Errorf("invalid VM name: %s", name)
}
```

### Additional Security Measures

1. **Validate input** before passing to PowerShell
2. **Use `-NoProfile`** to prevent profile scripts from running
3. **Use `-NonInteractive`** to prevent interactive prompts
4. **Set timeouts** with `context.Context` to prevent hanging
5. **Log commands** for audit trails (without sensitive data)

## JSON Output Pattern

PowerShell can output structured data as JSON for easy parsing:

```go
psScript := `
    Get-Process | Select-Object Name, Id, CPU | ConvertTo-Json
`
output, _ := m.Exec.RunCommand(psScript)

var processes []Process
json.Unmarshal(output, &processes)
```

**Handling single vs array output:**

```go
outputStr := strings.TrimSpace(string(output))

// PowerShell returns object for single item, array for multiple
if strings.HasPrefix(outputStr, "{") {
    var single Item
    json.Unmarshal(output, &single)
    items = append(items, single)
} else if strings.HasPrefix(outputStr, "[") {
    json.Unmarshal(output, &items)
}
```

## Error Wrapping Pattern

Always use `%w` to wrap errors for proper error chain:

```go
output, err := m.Exec.RunCommand(script)
if err != nil {
    return fmt.Errorf("operation failed: %w\nOutput: %s", err, string(output))
}
```

This allows callers to use `errors.Is()` and `errors.As()` for error handling.
