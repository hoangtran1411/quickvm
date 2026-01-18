---
description: Add a new feature to the TUI interface
---

# Add TUI Feature

Follow these steps to add features to the terminal UI:

## 1. Identify the Change

Common TUI changes:
- New keyboard shortcut
- New column in table
- New status message type
- Visual style changes

## 2. Update Model (if needed)

In `ui/table.go`, modify the Model struct:

```go
type Model struct {
    table      table.Model
    vms        []hyperv.VM
    manager    *hyperv.Manager
    message    string
    err        error
    // Add new fields here
    newField   string
}
```

## 3. Add Keyboard Shortcut

In the `Update` function:

```go
case tea.KeyMsg:
    switch msg.String() {
    // ... existing cases ...
    
    case "x":  // New shortcut
        if len(m.vms) > 0 && m.table.Cursor() < len(m.vms) {
            selectedVM := m.vms[m.table.Cursor()]
            m.message = fmt.Sprintf("Performing action on: %s...", selectedVM.Name)
            return m, m.newAction(selectedVM.Index)
        }
    }
```

## 4. Add Command Function

```go
func (m Model) newAction(index int) tea.Cmd {
    return func() tea.Msg {
        if err := m.manager.DoSomething(index); err != nil {
            return errMsg{err}
        }
        // Reload VMs after action
        return m.loadVMs()
    }
}
```

## 5. Add New Column (if needed)

In `NewModel`:
```go
columns := []table.Column{
    // ... existing columns ...
    {Title: "New Column", Width: 15},
}
```

In `updateTable`:
```go
rows = append(rows, table.Row{
    // ... existing fields ...
    vm.NewField,
})
```

## 6. Update Help Text

In the `View` function:
```go
help := helpStyle.Render(
    "↑/↓: Navigate • Enter: Start • s: Stop • t: Restart • x: NewAction • r: Refresh • q: Quit",
)
```

## 7. Add New Style (if needed)

```go
var newStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("205")).
    Bold(true)
```

## 8. Test Manually

// turbo
```powershell
go build -o quickvm.exe
./quickvm.exe
```

Test:
- Press the new keyboard shortcut
- Verify the action works
- Check the status message
- Ensure no visual glitches

## Color Reference

| Code | Color | Use Case |
|------|-------|----------|
| `46` | Green | Success |
| `196` | Red | Error/Stopped |
| `226` | Yellow | Warning |
| `205` | Pink | Accent |
| `57` | Purple | Selected |
| `240` | Gray | Borders |

## Checklist

- [ ] Model updated (if needed)
- [ ] Keyboard shortcut added
- [ ] Command function implemented
- [ ] Table columns updated (if needed)
- [ ] Help text updated
- [ ] Styles added (if needed)
- [ ] Manual testing done
