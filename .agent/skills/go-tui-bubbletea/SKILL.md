---
name: Go TUI with Bubble Tea
description: Build beautiful, interactive terminal user interfaces in Go using Charmbracelet's Bubble Tea framework with Lipgloss styling.
---

# Go TUI with Bubble Tea

This skill provides patterns for creating premium terminal user interfaces using [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea), [bubbles](https://github.com/charmbracelet/bubbles), and [lipgloss](https://github.com/charmbracelet/lipgloss).

## When to Use

- Interactive CLI applications requiring navigation
- Real-time data displays with refresh capability
- Form inputs and selection interfaces
- Applications needing visual appeal in terminal

## Dependencies

```bash
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss
```

## Architecture Pattern

```
ui/
├── table.go       # Table-based interfaces
├── form.go        # Form inputs
├── styles.go      # Lipgloss style definitions
└── messages.go    # Custom tea.Msg types
```

## Core Templates

### 1. Style Definitions (`ui/styles.go`)

```go
package ui

import "github.com/charmbracelet/lipgloss"

var (
    // Base container styles
    baseStyle = lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("240"))

    // Header style
    headerStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205")).
        Background(lipgloss.Color("235")).
        Padding(0, 1)

    // Selection highlight
    selectedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("229")).
        Background(lipgloss.Color("57")).
        Bold(true)

    // Title style
    titleStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("205")).
        Background(lipgloss.Color("235")).
        Bold(true).
        Padding(0, 1).
        MarginBottom(1)

    // Help text style
    helpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        MarginTop(1)

    // Status styles
    successStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("46")).
        Bold(true)

    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")).
        Bold(true)

    warningStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("226")).
        Bold(true)
)
```

### 2. Table Model (`ui/table.go`)

```go
package ui

import (
    "fmt"
    "strings"

    "github.com/charmbracelet/bubbles/table"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Custom message types
type itemsMsg []Item
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// Item represents a row in the table
type Item struct {
    ID     int
    Name   string
    Status string
    // Add more fields as needed
}

// Model is the Bubble Tea model for our TUI
type Model struct {
    table   table.Model
    items   []Item
    message string
    err     error
}

// NewModel creates a new TUI model
func NewModel() Model {
    columns := []table.Column{
        {Title: "ID", Width: 5},
        {Title: "Name", Width: 30},
        {Title: "Status", Width: 15},
    }

    t := table.New(
        table.WithColumns(columns),
        table.WithFocused(true),
        table.WithHeight(15),
    )

    // Apply styles
    s := table.DefaultStyles()
    s.Header = headerStyle
    s.Selected = selectedStyle
    t.SetStyles(s)

    return Model{
        table: t,
    }
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
    return m.loadItems
}

// loadItems is a command that loads items
func (m Model) loadItems() tea.Msg {
    // Replace with your data loading logic
    items := []Item{
        {ID: 1, Name: "Item 1", Status: "Active"},
        {ID: 2, Name: "Item 2", Status: "Inactive"},
    }
    return itemsMsg(items)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c", "esc":
            return m, tea.Quit

        case "r":
            m.message = "Refreshing..."
            return m, m.loadItems

        case "enter":
            if len(m.items) > 0 && m.table.Cursor() < len(m.items) {
                selected := m.items[m.table.Cursor()]
                m.message = fmt.Sprintf("Selected: %s", selected.Name)
            }
        }

    case itemsMsg:
        m.items = msg
        m.updateTable()
        m.message = "Items loaded!"
        return m, nil

    case errMsg:
        m.err = msg.err
        m.message = fmt.Sprintf("Error: %v", msg.err)
        return m, nil
    }

    m.table, cmd = m.table.Update(msg)
    return m, cmd
}

// updateTable updates the table rows from items
func (m *Model) updateTable() {
    rows := []table.Row{}
    for _, item := range m.items {
        status := item.Status
        switch status {
        case "Active":
            status = successStyle.Render(status)
        case "Inactive":
            status = errorStyle.Render(status)
        default:
            status = warningStyle.Render(status)
        }

        rows = append(rows, table.Row{
            fmt.Sprintf("%d", item.ID),
            item.Name,
            status,
        })
    }
    m.table.SetRows(rows)
}

// View renders the TUI
func (m Model) View() string {
    var b strings.Builder

    // Title
    title := titleStyle.Render("🚀 My Application")
    b.WriteString(title)
    b.WriteString("\n\n")

    // Table
    b.WriteString(baseStyle.Render(m.table.View()))
    b.WriteString("\n")

    // Message
    if m.message != "" {
        b.WriteString("\n")
        if m.err != nil {
            b.WriteString(errorStyle.Render(m.message))
        } else {
            b.WriteString(successStyle.Render(m.message))
        }
        b.WriteString("\n")
    }

    // Help
    help := helpStyle.Render(
        "↑/↓: Navigate • Enter: Select • r: Refresh • q: Quit",
    )
    b.WriteString("\n")
    b.WriteString(help)

    return b.String()
}
```

### 3. Running the TUI

```go
package main

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"
    "yourapp/ui"
)

func main() {
    p := tea.NewProgram(ui.NewModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
    }
}
```

## Color Reference

### Common Lipgloss Colors

| Code | Color | Use Case |
|------|-------|----------|
| `46` | Green | Success, Running |
| `196` | Red | Error, Stopped |
| `226` | Yellow | Warning, Pending |
| `205` | Pink | Accent, Headers |
| `229` | Light Yellow | Selected |
| `57` | Purple | Selection Background |
| `240` | Gray | Borders |
| `241` | Light Gray | Help text |
| `235` | Dark Gray | Background |

## Best Practices

1. **Separate styles**: Keep style definitions in a separate file for maintainability

2. **Use the Elm Architecture**: 
   - `Model` - Application state
   - `Update` - Handle messages, return new model + commands
   - `View` - Render the model to a string

3. **Async operations as Commands**: Return `tea.Cmd` for async operations
   ```go
   func (m Model) doAsyncTask() tea.Cmd {
       return func() tea.Msg {
           result, err := asyncOperation()
           if err != nil {
               return errMsg{err}
           }
           return resultMsg(result)
       }
   }
   ```

4. **Handle window resize**:
   ```go
   case tea.WindowSizeMsg:
       m.width = msg.Width
       m.height = msg.Height
   ```

5. **Use consistent keybindings**:
   - `q` / `Esc` / `Ctrl+C` - Quit
   - `r` - Refresh
   - `Enter` - Select/Confirm
   - `↑/↓` or `j/k` - Navigate

6. **Provide visual feedback** for all actions

## Integration with CLI

To launch TUI when no subcommand is provided:

```go
// In cmd/root.go
var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "My app",
    Run: func(cmd *cobra.Command, args []string) {
        // Launch TUI when no subcommand
        p := tea.NewProgram(ui.NewModel())
        if _, err := p.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
    },
}
```
