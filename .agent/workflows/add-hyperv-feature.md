---
description: Add a new Hyper-V feature/operation to the project
---

# Add New Hyper-V Feature

Follow these steps to add a new Hyper-V feature:

## 1. Research PowerShell Command

First, identify the PowerShell cmdlet needed:
```powershell
Get-Help <Cmdlet> -Full
# Example: Get-Help Get-VMSnapshot -Full
```

## 2. Create Feature File

// turbo
Create `hyperv/<feature>.go`:

```go
package hyperv

import (
    "encoding/json"
    "fmt"
    "strings"
)

// <Feature>Info represents the data structure
type <Feature>Info struct {
    Name   string `json:"name"`
    Status string `json:"status"`
    // Add more fields as needed
}

// Get<Feature>s retrieves all <feature>s for a VM
func (m *Manager) Get<Feature>s(vmIndex int) ([]<Feature>Info, error) {
    vms, err := m.GetVMs()
    if err != nil {
        return nil, err
    }

    if vmIndex < 1 || vmIndex > len(vms) {
        return nil, fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", vmIndex, len(vms))
    }

    vm := vms[vmIndex-1]
    
    psScript := fmt.Sprintf(`
        Get-<Feature> -VMName '%s' | Select-Object Name, Status | ConvertTo-Json
    `, strings.ReplaceAll(vm.Name, "'", "''"))

    output, err := m.Exec.RunCommand(psScript)
    if err != nil {
        return nil, fmt.Errorf("failed to get <feature>s: %w\nOutput: %s", err, string(output))
    }

    // Parse JSON output
    var items []<Feature>Info
    if err := json.Unmarshal(output, &items); err != nil {
        // Handle single item case
        var single <Feature>Info
        if err := json.Unmarshal(output, &single); err != nil {
            return nil, fmt.Errorf("failed to parse output: %w", err)
        }
        items = append(items, single)
    }

    return items, nil
}

// Create<Feature> creates a new <feature>
func (m *Manager) Create<Feature>(vmIndex int, name string) error {
    vms, err := m.GetVMs()
    if err != nil {
        return err
    }

    if vmIndex < 1 || vmIndex > len(vms) {
        return fmt.Errorf("invalid VM index: %d", vmIndex)
    }

    vm := vms[vmIndex-1]
    escapedName := strings.ReplaceAll(name, "'", "''")
    
    psScript := fmt.Sprintf(`New-<Feature> -VMName '%s' -Name '%s'`,
        strings.ReplaceAll(vm.Name, "'", "''"),
        escapedName)

    output, err := m.Exec.RunCommand(psScript)
    if err != nil {
        return fmt.Errorf("failed to create <feature>: %w\nOutput: %s", err, string(output))
    }

    return nil
}
```

## 3. Add Tests

// turbo
Create `hyperv/<feature>_test.go`:

```go
package hyperv

import (
    "errors"
    "testing"
)

func TestGet<Feature>s_Success(t *testing.T) {
    // Mock VM list
    vmListOutput := []byte(`[{"name":"TestVM","state":"Running"}]`)
    featureOutput := []byte(`[{"name":"Feature1","status":"Active"}]`)

    mock := &MockSequenceExecutor{
        Outputs: [][]byte{vmListOutput, featureOutput},
    }

    manager := NewManagerWithExecutor(mock)
    items, err := manager.Get<Feature>s(1)

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(items) != 1 {
        t.Errorf("expected 1 item, got %d", len(items))
    }
}

func TestGet<Feature>s_InvalidIndex(t *testing.T) {
    vmListOutput := []byte(`[{"name":"TestVM","state":"Running"}]`)

    mock := &MockShellExecutor{Output: vmListOutput}
    manager := NewManagerWithExecutor(mock)

    _, err := manager.Get<Feature>s(999)

    if err == nil {
        t.Fatal("expected error for invalid index")
    }
}

func TestCreate<Feature>_Success(t *testing.T) {
    vmListOutput := []byte(`[{"name":"TestVM","state":"Running"}]`)
    createOutput := []byte(``)

    mock := &MockSequenceExecutor{
        Outputs: [][]byte{vmListOutput, createOutput},
    }

    manager := NewManagerWithExecutor(mock)
    err := manager.Create<Feature>(1, "NewFeature")

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}
```

## 4. Create CLI Commands

Use workflow `/add-command` to create the CLI interface for this feature.

## 5. Run Tests

// turbo
```powershell
go test -v ./hyperv/...
```

## 6. Update Documentation

Add feature documentation to `docs/` if needed.

## Checklist

- [ ] PowerShell command researched and tested manually
- [ ] Feature file created in `hyperv/`
- [ ] Proper escaping for VM names (prevent injection)
- [ ] Unit tests with mocks
- [ ] CLI commands created
- [ ] All tests pass
