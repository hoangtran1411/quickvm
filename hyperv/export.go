package hyperv

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExportVMOptions contains options for exporting a VM
type ExportVMOptions struct {
	VMIndex int
	Path    string
}

// ImportVMOptions contains options for importing a VM
type ImportVMOptions struct {
	Path          string
	Copy          bool   // Copy the VM files instead of registering in place
	GenerateNewID bool   // Generate new VM ID
	VHDPath       string // Optional: custom path for VHD files
}

// ExportVM exports a VM by index to the specified path
func (m *Manager) ExportVM(vmIndex int, path string) error {
	vms, err := m.GetVMs()
	if err != nil {
		return err
	}

	if vmIndex < 1 || vmIndex > len(vms) {
		return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", vmIndex, len(vms))
	}

	vm := vms[vmIndex-1]
	return m.ExportVMByName(vm.Name, path)
}

// ExportVMByName exports a VM by name to the specified path
func (m *Manager) ExportVMByName(vmName, path string) error {
	// PowerShell script to export VM
	psScript := fmt.Sprintf(`Export-VM -Name "%s" -Path "%s"`, vmName, path)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to export VM '%s': %v\nOutput: %s", vmName, err, string(output))
	}

	return nil
}

// ImportVM imports a VM from the specified path
func (m *Manager) ImportVM(opts ImportVMOptions) (string, error) {
	// Find the .vmcx file in the path
	vmcxPath, err := m.findVMCXFile(opts.Path)
	if err != nil {
		return "", err
	}

	// Build the PowerShell command
	var psScriptBuilder strings.Builder
	psScriptBuilder.WriteString(fmt.Sprintf(`$vm = Import-VM -Path "%s"`, vmcxPath))

	if opts.Copy {
		psScriptBuilder.Reset()
		psScriptBuilder.WriteString(fmt.Sprintf(`$vm = Import-VM -Path "%s" -Copy`, vmcxPath))
	}

	if opts.GenerateNewID {
		psScriptBuilder.WriteString(" -GenerateNewId")
	}

	if opts.VHDPath != "" {
		psScriptBuilder.WriteString(fmt.Sprintf(` -VhdDestinationPath "%s"`, opts.VHDPath))
	}

	psScriptBuilder.WriteString("; $vm.Name")

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScriptBuilder.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to import VM from '%s': %v\nOutput: %s", opts.Path, err, string(output))
	}

	// Extract VM name from output
	vmName := strings.TrimSpace(string(output))
	return vmName, nil
}

// findVMCXFile finds the .vmcx file in the given path
func (m *Manager) findVMCXFile(basePath string) (string, error) {
	// First, check if the path is already a .vmcx file
	if strings.HasSuffix(strings.ToLower(basePath), ".vmcx") {
		return basePath, nil
	}

	// Search in "Virtual Machines" subdirectory (standard Hyper-V export structure)
	vmDir := filepath.Join(basePath, "Virtual Machines")
	pattern := filepath.Join(vmDir, "*.vmcx")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("error searching for .vmcx file: %v", err)
	}

	if len(matches) > 0 {
		return matches[0], nil
	}

	// Also try directly in the basePath
	pattern = filepath.Join(basePath, "*.vmcx")
	matches, err = filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("error searching for .vmcx file: %v", err)
	}

	if len(matches) > 0 {
		return matches[0], nil
	}

	return "", fmt.Errorf("no .vmcx file found in '%s' or '%s'", basePath, vmDir)
}

// GetExportedVMInfo gets information about an exported VM
func (m *Manager) GetExportedVMInfo(path string) (map[string]string, error) {
	vmcxPath, err := m.findVMCXFile(path)
	if err != nil {
		return nil, err
	}

	psScript := fmt.Sprintf(`
		$report = Compare-VM -Path "%s"
		@{
			VMName = $report.VM.Name
			State = $report.VM.State.ToString()
			MemoryMB = [int]($report.VM.MemoryStartup/1MB)
			ProcessorCount = $report.VM.ProcessorCount
			Incompatibilities = ($report.Incompatibilities | ForEach-Object { $_.Message }) -join "; "
		} | ConvertTo-Json
	`, vmcxPath)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get VM info from '%s': %v\nOutput: %s", path, err, string(output))
	}

	// Parse JSON output
	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		return nil, fmt.Errorf("empty output from PowerShell")
	}

	var info map[string]interface{}
	if err := json.Unmarshal([]byte(outputStr), &info); err != nil {
		return nil, fmt.Errorf("failed to parse VM info: %v", err)
	}

	// Convert to string map
	result := make(map[string]string)
	for k, v := range info {
		result[k] = fmt.Sprintf("%v", v)
	}

	return result, nil
}
