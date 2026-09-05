// Package hyperv provides interfaces and methods for interacting with Hyper-V.
package hyperv

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// DefaultOperationTimeout is the fallback timeout for VM operations when no deadline is set on ctx.
const DefaultOperationTimeout = 60 * time.Second

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
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("execution failed: %w", err)
	}
	return out, nil
}

// isKnownValuedParam returns true if the flag/parameter expects a following value argument.
func isKnownValuedParam(param string) bool {
	switch strings.ToLower(param) {
	case "-name", "-newname", "-snapshotname", "-vmname", "-path",
		"-vhddestinationpath", "-erroraction", "-expandproperty",
		"/t", "/c":
		return true
	default:
		return false
	}
}

// isKnownSwitchOrOperator returns true for known parameter switches or operators.
func isKnownSwitchOrOperator(arg string) bool {
	switch strings.ToLower(arg) {
	case "|", "select-object", "-force", "-passthru", "-copy", "-generatenewid", "-confirm:$false", "-confirm:$true", "/r":
		return true
	default:
		return false
	}
}

func isSafeIdentRune(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z':
		return true
	case r >= 'A' && r <= 'Z':
		return true
	case r >= '0' && r <= '9':
		return true
	case r == '_':
		return true
	default:
		return false
	}
}

func isSafeSlashParam(s string) bool {
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		default:
			return false
		}
	}
	return true
}

func isSafeFlagParam(s string) bool {
	for _, r := range s {
		if !isSafeIdentRune(r) {
			return false
		}
	}
	return true
}

// formatCmdletScript formats a PowerShell cmdlet invocation safely.
// Parameter values are strictly single-quoted with internal single-quotes escaped.
// Values following value-expecting flags are always treated as values even if prefixed with '-'.
func formatCmdletScript(cmdlet string, args ...string) string {
	parts := make([]string, 0, 1+len(args))
	parts = append(parts, cmdlet)
	expectingValue := false

	for _, arg := range args {
		if expectingValue {
			parts = append(parts, "'"+strings.ReplaceAll(arg, "'", "''")+"'")
			expectingValue = false
			continue
		}

		if isKnownValuedParam(arg) {
			parts = append(parts, arg)
			expectingValue = true
			continue
		}

		if isKnownSwitchOrOperator(arg) {
			parts = append(parts, arg)
			continue
		}

		if strings.HasPrefix(arg, "-") && len(arg) > 1 && isSafeFlagParam(arg[1:]) {
			parts = append(parts, arg)
			continue
		}
		if strings.HasPrefix(arg, "/") && len(arg) > 1 && isSafeSlashParam(arg[1:]) {
			parts = append(parts, arg)
			continue
		}

		parts = append(parts, "'"+strings.ReplaceAll(arg, "'", "''")+"'")
	}
	return strings.Join(parts, " ")
}

// RunCmdlet executes a PowerShell cmdlet safely, ensuring parameter values are properly single-quoted and escaped
func (p *PowerShellRunner) RunCmdlet(ctx context.Context, cmdlet string, args ...string) ([]byte, error) {
	script := formatCmdletScript(cmdlet, args...)
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("cmdlet execution failed: %w", err)
	}
	return out, nil
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
//
//nolint:funlen,gocyclo // Parsing logic is verbose and cyclomatic complexity is high
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
	switch {
	case strings.HasPrefix(outputStr, "{"):
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
	case strings.HasPrefix(outputStr, "["):
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
	default:
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
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("VM name cannot be empty")
	}

	// why: Fallback to DefaultOperationTimeout ensures external PowerShell processes
	// do not hang indefinitely if the caller supplied a context without a deadline.
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultOperationTimeout)
		defer cancel()
	}

	// why: Using RunCmdlet with separate args prevents shell injection attacks
	// where 'name' could contain malicious PowerShell commands.
	output, err := m.Exec.RunCmdlet(ctx, "Start-VM", "-Name", name)
	if err != nil {
		outputStr := string(output)
		// why: Hyper-V throws an exception when Start-VM is called on an already running VM.
		// Treating this as an idempotent success avoids reporting false errors during batch or UI operations.
		if strings.Contains(strings.ToLower(outputStr), "already running") ||
			strings.Contains(strings.ToLower(err.Error()), "already running") {
			return nil
		}
		return fmt.Errorf("failed to start VM '%s': %w\nOutput: %s", name, err, outputStr)
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
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("VM name cannot be empty")
	}

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
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("VM name cannot be empty")
	}

	// why: Enforce context timeout/cancellation and safe execution.
	output, err := m.Exec.RunCmdlet(ctx, "Restart-VM", "-Name", name, "-Force")
	if err != nil {
		return fmt.Errorf("failed to restart VM '%s': %v\nOutput: %s", name, err, string(output))
	}
	return nil
}

// GetVMStatus gets the status of a specific VM by name
func (m *Manager) GetVMStatus(ctx context.Context, name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("VM name cannot be empty")
	}

	// why: Use RunCmdlet to safely query properties without script injection risks.
	output, err := m.Exec.RunCmdlet(ctx, "Get-VM", "-Name", name, "|", "Select-Object", "-ExpandProperty", "State")
	if err != nil {
		return "", fmt.Errorf("failed to get VM status: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}
