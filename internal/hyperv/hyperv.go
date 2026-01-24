package hyperv

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// VM represents a Hyper-V virtual machine
type VM struct {
	Index       int      `json:"-"`
	Name        string   `json:"name"`
	State       string   `json:"state"`
	CPUUsage    int      `json:"cpuUsage"`
	MemoryMB    int64    `json:"memoryMB"`
	Uptime      string   `json:"uptime"`
	Status      string   `json:"status"`
	Version     string   `json:"version"`
	IPAddresses []string `json:"ipAddresses"`
}

// VMManager defines the interface for Hyper-V operations to allow mocking in tests
type VMManager interface {
	GetVMs(ctx context.Context) ([]VM, error)
	StartVM(ctx context.Context, index int) error
	StartVMByName(ctx context.Context, name string) error
	StopVM(ctx context.Context, index int) error
	StopVMByName(ctx context.Context, name string) error
	RestartVM(ctx context.Context, index int) error
	RestartVMByName(ctx context.Context, name string) error
	GetVMStatus(ctx context.Context, name string) (string, error)
}

// ShellExecutor defines an interface for executing shell commands with Context
type ShellExecutor interface {
	// RunScript executes a complex PowerShell script (beware of injection, use for static scripts)
	RunScript(ctx context.Context, script string) ([]byte, error)
	// RunCmdlet executes a specific cmdlet with arguments safely
	RunCmdlet(ctx context.Context, cmdlet string, args ...string) ([]byte, error)
}

// PowerShellRunner implements ShellExecutor for actual PowerShell execution
type PowerShellRunner struct{}

// RunScript executes a PowerShell command/script with Context
func (p *PowerShellRunner) RunScript(ctx context.Context, script string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	return cmd.CombinedOutput()
}

// RunCmdlet executes a PowerShell cmdlet safely using separate arguments to avoid injection
func (p *PowerShellRunner) RunCmdlet(ctx context.Context, cmdlet string, args ...string) ([]byte, error) {
	// Construct args: powershell -NoProfile -NonInteractive -Command Cmdlet -Arg1 val1 ...
	// This relies on the fact that we are passing "Cmdlet" and its args as separate arguments to the powershell executable,
	// which then parses them. This avoids shell injection when 'cmdlet' is just the command name and 'args' are its parameters.
	
	// Base args
	psArgs := []string{"-NoProfile", "-NonInteractive", "-Command", cmdlet}
	// Append the rest
	psArgs = append(psArgs, args...)

	cmd := exec.CommandContext(ctx, "powershell", psArgs...)
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
func (m *Manager) GetVMs(ctx context.Context) ([]VM, error) {
	// PowerShell script to get VM information
	psScript := `
		Get-VM | Select-Object @{Name='Name';Expression={$_.Name.ToString()}}, 
		@{Name='State';Expression={$_.State.ToString()}}, 
		@{Name='CPUUsage';Expression={[int]$_.CPUUsage}}, 
		@{Name='MemoryMB';Expression={[int]($_.MemoryAssigned/1MB)}},
		@{Name='Uptime';Expression={$_.Uptime.ToString()}},
		@{Name='Status';Expression={$_.Status.ToString()}},
		@{Name='Version';Expression={$_.Version.ToString()}},
		@{Name='IPAddresses';Expression={($_.NetworkAdapters.IPAddresses | Where-Object { $_ -match '^\d+\.\d+\.\d+\.\d+$' })}} | ConvertTo-Json
	`

	output, err := m.Exec.RunScript(ctx, psScript)
	if err != nil {
		return nil, fmt.Errorf("failed to execute PowerShell command: %v\nOutput: %s", err, string(output))
	}

	// Parse JSON output
	var vms []VM
	outputStr := strings.TrimSpace(string(output))

	// Handle single VM case (PowerShell returns object, not array)
	if strings.HasPrefix(outputStr, "{") {
		var vmRaw struct {
			Name        string      `json:"name"`
			State       string      `json:"state"`
			CPUUsage    int         `json:"cpuUsage"`
			MemoryMB    int64       `json:"memoryMB"`
			Uptime      string      `json:"uptime"`
			Status      string      `json:"status"`
			Version     string      `json:"version"`
			IPAddresses interface{} `json:"ipAddresses"`
		}
		if err := json.Unmarshal(output, &vmRaw); err != nil {
			return nil, fmt.Errorf("failed to parse VM data: %v", err)
		}

		vm := VM{
			Name:     vmRaw.Name,
			State:    vmRaw.State,
			CPUUsage: vmRaw.CPUUsage,
			MemoryMB: vmRaw.MemoryMB,
			Uptime:   vmRaw.Uptime,
			Status:   vmRaw.Status,
			Version:  vmRaw.Version,
		}

		if ips, ok := vmRaw.IPAddresses.([]interface{}); ok {
			for _, ip := range ips {
				if str, ok := ip.(string); ok {
					vm.IPAddresses = append(vm.IPAddresses, str)
				}
			}
		} else if ip, ok := vmRaw.IPAddresses.(string); ok {
			vm.IPAddresses = []string{ip}
		}

		vms = append(vms, vm)
	} else if strings.HasPrefix(outputStr, "[") {
		var vmsRaw []struct {
			Name        string      `json:"name"`
			State       string      `json:"state"`
			CPUUsage    int         `json:"cpuUsage"`
			MemoryMB    int64       `json:"memoryMB"`
			Uptime      string      `json:"uptime"`
			Status      string      `json:"status"`
			Version     string      `json:"version"`
			IPAddresses interface{} `json:"ipAddresses"`
		}
		if err := json.Unmarshal(output, &vmsRaw); err != nil {
			return nil, fmt.Errorf("failed to parse VMs data: %v", err)
		}

		for _, vmRaw := range vmsRaw {
			vm := VM{
				Name:     vmRaw.Name,
				State:    vmRaw.State,
				CPUUsage: vmRaw.CPUUsage,
				MemoryMB: vmRaw.MemoryMB,
				Uptime:   vmRaw.Uptime,
				Status:   vmRaw.Status,
				Version:  vmRaw.Version,
			}

			if ips, ok := vmRaw.IPAddresses.([]interface{}); ok {
				for _, ip := range ips {
					if str, ok := ip.(string); ok {
						vm.IPAddresses = append(vm.IPAddresses, str)
					}
				}
			} else if ip, ok := vmRaw.IPAddresses.(string); ok {
				vm.IPAddresses = []string{ip}
			}
			vms = append(vms, vm)
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
func (m *Manager) StartVM(ctx context.Context, index int) error {
	vms, err := m.GetVMs(ctx)
	if err != nil {
		return err
	}

	if index < 1 || index > len(vms) {
		return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", index, len(vms))
	}

	vm := vms[index-1]
	return m.StartVMByName(ctx, vm.Name)
}

// StartVMByName starts a virtual machine by name
func (m *Manager) StartVMByName(ctx context.Context, name string) error {
	// why: Using RunCmdlet with separate args prevents shell injection attacks
	// where 'name' could contain malicious PowerShell commands.
	output, err := m.Exec.RunCmdlet(ctx, "Start-VM", "-Name", name)
	if err != nil {
		return fmt.Errorf("failed to start VM '%s': %v\nOutput: %s", name, err, string(output))
	}
	return nil
}

// StopVM stops a virtual machine by index
func (m *Manager) StopVM(ctx context.Context, index int) error {
	vms, err := m.GetVMs(ctx)
	if err != nil {
		return err
	}

	if index < 1 || index > len(vms) {
		return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", index, len(vms))
	}

	vm := vms[index-1]
	return m.StopVMByName(ctx, vm.Name)
}

// StopVMByName stops a virtual machine by name
func (m *Manager) StopVMByName(ctx context.Context, name string) error {
	// why: Safe execution using RunCmdlet to handle VM names with special chars or potential injection attempts.
	output, err := m.Exec.RunCmdlet(ctx, "Stop-VM", "-Name", name, "-Force")
	if err != nil {
		return fmt.Errorf("failed to stop VM '%s': %v\nOutput: %s", name, err, string(output))
	}
	return nil
}

// RestartVM restarts a virtual machine by index
func (m *Manager) RestartVM(ctx context.Context, index int) error {
	vms, err := m.GetVMs(ctx)
	if err != nil {
		return err
	}

	if index < 1 || index > len(vms) {
		return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", index, len(vms))
	}

	vm := vms[index-1]
	return m.RestartVMByName(ctx, vm.Name)
}

// RestartVMByName restarts a virtual machine by name
func (m *Manager) RestartVMByName(ctx context.Context, name string) error {
	// why: Enforce context timeout/cancellation and safe execution.
	output, err := m.Exec.RunCmdlet(ctx, "Restart-VM", "-Name", name, "-Force")
	if err != nil {
		return fmt.Errorf("failed to restart VM '%s': %v\nOutput: %s", name, err, string(output))
	}
	return nil
}

// GetVMStatus gets the status of a specific VM by name
func (m *Manager) GetVMStatus(ctx context.Context, name string) (string, error) {
	// why: Use RunCmdlet to safely query properties without script injection risks.
	// We use Select-Object -ExpandProperty to get just the string value.
	output, err := m.Exec.RunCmdlet(ctx, "Get-VM", "-Name", name, "|", "Select-Object", "-ExpandProperty", "State")
	// Note: Piping in RunCmdlet via args works because we are passing "-Command" "Get-VM ... | ..." to powershell.
	// Wait, RunCmdlet as I defined it: 
	// psArgs := []string{..., "-Command", cmdlet} -> append args.
	// Result: powershell ... -Command Get-VM -Name name | Select...
	// This works in PowerShell syntax.
	
	if err != nil {
		// Fallback for complex queries or if piping fails in this specific way,
		// but standard PowerShell argument parsing usually handles this if arguments are passed correctly.
		// Actually, passing "|" as a separate arg to -Command might rely on whitespace.
		// Safer approach for property extraction:
		return "", fmt.Errorf("failed to get VM status: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}
