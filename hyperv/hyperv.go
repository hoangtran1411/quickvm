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

// Manager handles Hyper-V operations
type Manager struct{}

// NewManager creates a new Hyper-V manager
func NewManager() *Manager {
	return &Manager{}
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

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
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
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
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
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
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
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restart VM '%s': %v\nOutput: %s", name, err, string(output))
	}
	return nil
}

// GetVMStatus gets the status of a specific VM by name
func (m *Manager) GetVMStatus(name string) (string, error) {
	psScript := fmt.Sprintf(`(Get-VM -Name "%s").State`, name)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get VM status: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}
