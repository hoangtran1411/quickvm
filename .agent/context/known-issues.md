# Known Issues & Edge Cases

This document helps AI agents understand common pitfalls and how to avoid them.

## PowerShell Execution

### Issue: Empty JSON Array
PowerShell returns different formats for single vs multiple items:
```powershell
# Single item: returns object
Get-VM | Where State -eq 'Running' | ConvertTo-Json
# {"Name": "VM1", ...}

# Multiple items: returns array
# [{"Name": "VM1"}, {"Name": "VM2"}]
```

**Solution**: Always wrap in array:
```powershell
@(Get-VM) | ConvertTo-Json
```

### Issue: Unicode in VM Names
VM names with special characters can cause parsing issues.

**Solution**: Use `-Compress` and proper encoding:
```go
script := `[Console]::OutputEncoding = [Text.Encoding]::UTF8; Get-VM | ConvertTo-Json -Compress`
```

### Issue: Command Timeout
Long operations (clone, export) may timeout.

**Solution**: Use generous timeouts:
```go
ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
```

## Hyper-V Quirks

### VM States
Valid states and transitions:
```
Off → Running (Start-VM)
Running → Off (Stop-VM)
Running → Saved (Save-VM)
Running → Paused (Suspend-VM)
Saved → Running (Start-VM)
Paused → Running (Resume-VM)
```

### Generation 1 vs 2
- Gen 1: Legacy BIOS, IDE drives
- Gen 2: UEFI, SCSI drives, Secure Boot
- GPU-P only works with Gen 2

### Memory
- Static vs Dynamic memory affects available RAM
- MemoryAssigned is in bytes, divide by 1MB for display

## Testing Gotchas

### Mock JSON Must Match Real Output
PowerShell JSON has specific formatting:
```json
{
  "Name": "TestVM",
  "State": 2,          // Enum as integer
  "CPUUsage": 0,
  "MemoryAssigned": 2147483648  // Bytes, not MB
}
```

Map state integers:
```go
var stateMap = map[int]string{
    2: "Running",
    3: "Off",
    6: "Saved",
    9: "Paused",
}
```

### golangci-lint v2 Migration
Old v1 config will fail. Required changes:
```yaml
# v2 schema (REQUIRED)
version: "2"

# Exclusions moved
linters:
  exclusions:
    rules:
      - path: _test\.go
        linters: [errcheck]
```

## TUI Issues

### Terminal Size
Always handle window resize:
```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    m.table.SetHeight(msg.Height - 5)
```

### Flickering
Use alt screen buffer:
```go
p := tea.NewProgram(model, tea.WithAltScreen())
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| String concat in PowerShell | Use parameterized commands |
| Missing context.Context | Add as first parameter |
| Ignoring errors | Always wrap with `fmt.Errorf` |
| Testing on Linux | Use mocks, not real Hyper-V |
| golangci-lint v1 config | Upgrade to v2 schema |
