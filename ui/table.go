package ui

import (
	"context"
	"fmt"
	"strings"

	"quickvm/internal/hyperv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Background(lipgloss.Color("235")).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	statusRunningStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("46")).
				Bold(true)

	statusStoppedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Bold(true)

	statusOtherStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("226")).
				Bold(true)
)

type Model struct {
	table   table.Model
	vms     []hyperv.VM
	manager *hyperv.Manager
	message string
	err     error
}

type vmListMsg []hyperv.VM
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func NewModel() Model {
	columns := []table.Column{
		{Title: "Index", Width: 7},
		{Title: "Name", Width: 25},
		{Title: "State", Width: 10},
		{Title: "IP Address", Width: 15},
		{Title: "CPU%", Width: 6},
		{Title: "Memory(MB)", Width: 12},
		{Title: "Uptime", Width: 18},
		{Title: "Status", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Selected = selectedStyle
	t.SetStyles(s)

	return Model{
		table:   t,
		manager: hyperv.NewManager(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadVMs
}

func (m Model) loadVMs() tea.Msg {
	vms, err := m.manager.GetVMs(context.TODO())
	if err != nil {
		return errMsg{err}
	}
	return vmListMsg(vms)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "r":
			// Refresh VM list
			m.message = "Refreshing VM list..."
			return m, m.loadVMs

		case "enter":
			// Start selected VM
			if len(m.vms) > 0 && m.table.Cursor() < len(m.vms) {
				selectedVM := m.vms[m.table.Cursor()]
				m.message = fmt.Sprintf("Starting VM: %s...", selectedVM.Name)
				return m, m.startVM(selectedVM.Index)
			}

		case "s":
			// Stop selected VM
			if len(m.vms) > 0 && m.table.Cursor() < len(m.vms) {
				selectedVM := m.vms[m.table.Cursor()]
				m.message = fmt.Sprintf("Stopping VM: %s...", selectedVM.Name)
				return m, m.stopVM(selectedVM.Index)
			}

		case "t":
			// Restart selected VM
			if len(m.vms) > 0 && m.table.Cursor() < len(m.vms) {
				selectedVM := m.vms[m.table.Cursor()]
				m.message = fmt.Sprintf("Restarting VM: %s...", selectedVM.Name)
				return m, m.restartVM(selectedVM.Index)
			}
		}

	case vmListMsg:
		m.vms = msg
		m.updateTable()
		if m.message == "Refreshing VM list..." {
			m.message = "VM list refreshed!"
		}
		return m, nil

	case errMsg:
		m.err = msg.err
		m.message = fmt.Sprintf("Error: %v", msg.err)
		return m, nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *Model) updateTable() {
	rows := []table.Row{}
	for _, vm := range m.vms {
		var state string
		switch strings.ToLower(vm.State) {
		case "running":
			state = statusRunningStyle.Render(vm.State)
		case "off":
			state = statusStoppedStyle.Render(vm.State)
		default:
			state = statusOtherStyle.Render(vm.State)
		}

		ip := ""
		if len(vm.IPAddresses) > 0 {
			ip = vm.IPAddresses[0] // Show first IP
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", vm.Index),
			vm.Name,
			state,
			ip,
			fmt.Sprintf("%d%%", vm.CPUUsage),
			fmt.Sprintf("%d", vm.MemoryMB),
			vm.Uptime,
			vm.Status,
		})
	}
	m.table.SetRows(rows)
}

func (m Model) startVM(index int) tea.Cmd {
	return func() tea.Msg {
		if err := m.manager.StartVM(context.TODO(), index); err != nil {
			return errMsg{err}
		}
		// Reload VMs after action
		return m.loadVMs()
	}
}

func (m Model) stopVM(index int) tea.Cmd {
	return func() tea.Msg {
		if err := m.manager.StopVM(context.TODO(), index); err != nil {
			return errMsg{err}
		}
		// Reload VMs after action
		return m.loadVMs()
	}
}

func (m Model) restartVM(index int) tea.Cmd {
	return func() tea.Msg {
		if err := m.manager.RestartVM(context.TODO(), index); err != nil {
			return errMsg{err}
		}
		// Reload VMs after action
		return m.loadVMs()
	}
}

func (m Model) View() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("🖥️  QuickVM - Hyper-V Manager")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Table
	b.WriteString(baseStyle.Render(m.table.View()))
	b.WriteString("\n")

	// Message
	if m.message != "" {
		b.WriteString("\n")
		if m.err != nil {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.message))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(m.message))
		}
		b.WriteString("\n")
	}

	// Help
	help := helpStyle.Render(
		"↑/↓: Navigate • Enter: Start VM • s: Stop VM • t: Restart VM • r: Refresh • q: Quit",
	)
	b.WriteString("\n")
	b.WriteString(help)

	return b.String()
}
