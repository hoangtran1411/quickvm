package hyperv

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// VM represents a Hyper-V virtual machine
type VM struct {
	Index    int    `json:"-"`
	Name     string `json:"name"`
	State    string `json:"state"`
	CPUUsage int    `json:"cpuUsage"`
	MemoryMB int64  `json:"memoryMB"`
	Uptime   string `json:"uptime"`
	Status   string `json:"status"`
	Version  string `json:"version"`
}

// VMManager defines the interface for Hyper-V operations to allow mocking in tests
type VMManager interface {
	GetVMs() ([]VM, error)
	StartVM(index int) error
	StartVMByName(name string) error
	StopVM(index int) error
	StopVMByName(name string) error
	RestartVM(index int) error
	RestartVMByName(name string) error
	GetVMStatus(name string) (string, error)
}

// ShellExecutor defines an interface for executing shell commands
type ShellExecutor interface {
	RunCommand(script string) ([]byte, error)
}

// PowerShellRunner implements ShellExecutor for actual PowerShell execution
type PowerShellRunner struct{}

// RunCommand executes a PowerShell command
func (p *PowerShellRunner) RunCommand(script string) ([]byte, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	return cmd.CombinedOutput()
}

// Manager handles Hyper-V operations
type Manager struct {
	Exec ShellExecutor
}

// NewManager creates a new Hyper-V manager with default PowerShell runner
func NewManager() *Manager {
	return &Manager{
		Exec: &PowerShellRunner{},
	}
}

// GetVMs retrieves all Hyper-V virtual machines
func (m *Manager) GetVMs() ([]VM, error) {
	// PowerShell script to get VM information
	psScript := `
		Get-VM | Select-Object @{Name='Name';Expression={$_.Name.ToString()}}, 
		@{Name='State';Expression={$_.State.ToString()}}, 
		@{Name='CPUUsage';Expression={[int]$_.CPUUsage}}, 
		@{Name='MemoryMB';Expression={[int]($_.MemoryAssigned/1MB)}},
		@{Name='Uptime';Expression={$_.Uptime.ToString()}},
		@{Name='Status';Expression={$_.Status.ToString()}},
		@{Name='Version';Expression={$_.Version.ToString()}} | ConvertTo-Json
	`

	output, err := m.Exec.RunCommand(psScript)
	if err != nil {
		return nil, fmt.Errorf("failed to execute PowerShell command: %v\nOutput: %s", err, string(output))
	}

	// Parse JSON output
	var vms []VM
	outputStr := strings.TrimSpace(string(output))

	// Handle single VM case (PowerShell returns object, not array)
	if strings.HasPrefix(outputStr, "{") {
		var vm VM
		if err := json.Unmarshal(output, &vm); err != nil {
			return nil, fmt.Errorf("failed to parse VM data: %v", err)
		}
		vms = append(vms, vm)
	} else if strings.HasPrefix(outputStr, "[") {
		if err := json.Unmarshal(output, &vms); err != nil {
			return nil, fmt.Errorf("failed to parse VMs data: %v", err)
		}
	} else {
		// If output is empty or doesn't start with JSON structure, check if it's an error or just no VMs.
		// If no VMs are present, PowerShell might return nothing or "[]".
		// Get-VM returns nothing if no VMs exist.
		if outputStr == "" {
			return []VM{}, nil
		}
		return nil, fmt.Errorf("no VMs found or invalid output format")
	}

	// Assign indices
	for i := range vms {
		vms[i].Index = i + 1
	}

	return vms, nil
}

// StartVM starts a virtual machine by index
func (m *Manager) StartVM(index int) error {
	vms, err := m.GetVMs()
	if err != nil {
		return err
	}

	if index < 1 || index > len(vms) {
		return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", index, len(vms))
	}

	vm := vms[index-1]
	return m.StartVMByName(vm.Name)
}

// StartVMByName starts a virtual machine by name
func (m *Manager) StartVMByName(name string) error {
	psScript := fmt.Sprintf(`Start-VM -Name "%s"`, name)
	output, err := m.Exec.RunCommand(psScript)
	if err != nil {
		return fmt.Errorf("failed to start VM '%s': %v\nOutput: %s", name, err, string(output))
	}
	return nil
}

// StopVM stops a virtual machine by index
func (m *Manager) StopVM(index int) error {
	vms, err := m.GetVMs()
	if err != nil {
		return err
	}

	if index < 1 || index > len(vms) {
		return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", index, len(vms))
	}

	vm := vms[index-1]
	return m.StopVMByName(vm.Name)
}

// StopVMByName stops a virtual machine by name
func (m *Manager) StopVMByName(name string) error {
	psScript := fmt.Sprintf(`Stop-VM -Name "%s" -Force`, name)
	output, err := m.Exec.RunCommand(psScript)
	if err != nil {
		return fmt.Errorf("failed to stop VM '%s': %v\nOutput: %s", name, err, string(output))
	}
	return nil
}

// RestartVM restarts a virtual machine by index
func (m *Manager) RestartVM(index int) error {
	vms, err := m.GetVMs()
	if err != nil {
		return err
	}

	if index < 1 || index > len(vms) {
		return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", index, len(vms))
	}

	vm := vms[index-1]
	return m.RestartVMByName(vm.Name)
}

// RestartVMByName restarts a virtual machine by name
func (m *Manager) RestartVMByName(name string) error {
	psScript := fmt.Sprintf(`Restart-VM -Name "%s" -Force`, name)
	output, err := m.Exec.RunCommand(psScript)
	if err != nil {
		return fmt.Errorf("failed to restart VM '%s': %v\nOutput: %s", name, err, string(output))
	}
	return nil
}

// GetVMStatus gets the status of a specific VM by name
func (m *Manager) GetVMStatus(name string) (string, error) {
	psScript := fmt.Sprintf(`(Get-VM -Name "%s").State`, name)
	output, err := m.Exec.RunCommand(psScript)
	if err != nil {
		return "", fmt.Errorf("failed to get VM status: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}
