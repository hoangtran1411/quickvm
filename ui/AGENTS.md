# AGENTS.md - Terminal UI (ui/)

This directory contains the interactive Terminal User Interface built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and styled with [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Architecture

```
ui/
└── table.go        # Main TUI model with VM table
```

## Bubble Tea Model Pattern

```go
type model struct {
    table    table.Model
    vms      []hyperv.VM
    manager  *hyperv.Manager
    message  string
    loading  bool
    quitting bool
}

// Init returns initial command
func (m model) Init() tea.Cmd {
    return m.loadVMs()
}

// Update handles messages and returns new model + commands
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            m.quitting = true
            return m, tea.Quit
        case "enter":
            return m, m.startSelectedVM()
        case "s":
            return m, m.stopSelectedVM()
        case "r":
            return m, m.refreshVMs()
        }
    case vmsLoadedMsg:
        m.vms = msg.vms
        m.loading = false
        return m, nil
    }
    
    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)
    return m, cmd
}

// View renders the UI
func (m model) View() string {
    if m.quitting {
        return ""
    }
    return m.renderTable() + m.renderStatus() + m.renderHelp()
}
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate VM list |
| `Enter` | Start selected VM |
| `s` | Stop selected VM |
| `t` | Restart selected VM |
| `r` | Refresh VM list |
| `i` | Show VM info |
| `c` | Clone VM |
| `d` | Connect via RDP |
| `q` | Quit |

## Lipgloss Styling

```go
var (
    // Color palette
    primaryColor   = lipgloss.Color("86")   // Cyan
    successColor   = lipgloss.Color("82")   // Green
    errorColor     = lipgloss.Color("196")  // Red
    warningColor   = lipgloss.Color("214")  // Orange
    mutedColor     = lipgloss.Color("240")  // Gray
    
    // Styles
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(primaryColor).
        MarginBottom(1)
    
    statusStyle = lipgloss.NewStyle().
        Foreground(mutedColor).
        MarginTop(1)
    
    // VM state colors
    stateStyles = map[string]lipgloss.Style{
        "Running": lipgloss.NewStyle().Foreground(successColor),
        "Off":     lipgloss.NewStyle().Foreground(errorColor),
        "Saved":   lipgloss.NewStyle().Foreground(warningColor),
        "Paused":  lipgloss.NewStyle().Foreground(warningColor),
    }
)
```

## Message Types

Define custom messages for async operations:

```go
// Message types
type vmsLoadedMsg struct {
    vms []hyperv.VM
    err error
}

type vmOperationMsg struct {
    success bool
    message string
    err     error
}

// Commands that return messages
func (m model) loadVMs() tea.Cmd {
    return func() tea.Msg {
        ctx := context.Background()
        vms, err := m.manager.GetVMs(ctx)
        return vmsLoadedMsg{vms: vms, err: err}
    }
}
```

## Table Configuration

Using [bubbles/table](https://github.com/charmbracelet/bubbles/tree/master/table):

```go
columns := []table.Column{
    {Title: "#", Width: 4},
    {Title: "Name", Width: 30},
    {Title: "State", Width: 10},
    {Title: "CPU", Width: 6},
    {Title: "Memory", Width: 10},
    {Title: "Uptime", Width: 15},
}

t := table.New(
    table.WithColumns(columns),
    table.WithRows(rows),
    table.WithFocused(true),
    table.WithHeight(10),
)

// Apply styles
s := table.DefaultStyles()
s.Header = s.Header.
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("240")).
    Bold(true)
s.Selected = s.Selected.
    Foreground(lipgloss.Color("229")).
    Background(lipgloss.Color("57")).
    Bold(true)
t.SetStyles(s)
```

## Adding New Features

1. Add keybinding in `Update()` switch
2. Create command function that returns `tea.Cmd`
3. Define message type for async results
4. Handle message in `Update()`
5. Update help text in `View()`

Example - Add "Export VM" feature:
```go
// In Update()
case "e":
    return m, m.exportSelectedVM()

// Command
func (m model) exportSelectedVM() tea.Cmd {
    return func() tea.Msg {
        vm := m.getSelectedVM()
        err := m.manager.ExportVM(context.Background(), vm.Name, exportPath)
        return vmOperationMsg{
            success: err == nil,
            message: fmt.Sprintf("Exported %s", vm.Name),
            err:     err,
        }
    }
}
```

## Testing TUI

TUI components are typically tested via integration:

```go
func TestModel_Update(t *testing.T) {
    m := initialModel()
    
    // Simulate key press
    newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
    
    model := newModel.(model)
    assert.True(t, model.quitting)
}
```

## Best Practices

1. **Keep Update() pure** - Return commands, don't execute side effects directly
2. **Use commands for I/O** - All Hyper-V operations should be `tea.Cmd`
3. **Handle loading states** - Show spinners/messages during async ops
4. **Graceful errors** - Display errors in UI, don't crash
5. **Consistent styling** - Use predefined Lipgloss styles
6. **Responsive layout** - Handle terminal resize with `tea.WindowSizeMsg`

## Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - UI components (table, spinner)
- `github.com/charmbracelet/lipgloss` - Styling
