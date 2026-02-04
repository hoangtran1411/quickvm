# Code Snippets Reference

Quick copy-paste patterns for common tasks.

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
