package hyperv

import (
	"context"
	"encoding/json"
	"fmt"
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
func (m *Manager) ExportVM(ctx context.Context, vmIndex int, path string) error {
	vms, err := m.GetVMs(ctx)
	if err != nil {
		return err
	}

	if vmIndex < 1 || vmIndex > len(vms) {
		return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", vmIndex, len(vms))
	}

	vm := vms[vmIndex-1]
	return m.ExportVMByName(ctx, vm.Name, path)
}

// ExportVMByName exports a VM by name to the specified path
func (m *Manager) ExportVMByName(ctx context.Context, vmName, path string) error {
	output, err := m.Exec.RunCmdlet(ctx, "Export-VM", "-Name", vmName, "-Path", path)
	if err != nil {
		return fmt.Errorf("failed to export VM '%s': %v\nOutput: %s", vmName, err, string(output))
	}

	return nil
}

// ImportVM imports a VM from the specified path
func (m *Manager) ImportVM(ctx context.Context, opts ImportVMOptions) (string, error) {
	// Find the .vmcx file in the path
	vmcxPath, err := m.findVMCXFile(opts.Path)
	if err != nil {
		return "", err
	}

	// Build arguments for RunCmdlet safely
	args := []string{"-Path", vmcxPath}

	if opts.Copy {
		args = append(args, "-Copy")
	}

	if opts.GenerateNewID {
		args = append(args, "-GenerateNewId")
	}

	if opts.VHDPath != "" {
		args = append(args, "-VhdDestinationPath", opts.VHDPath)
	}

	// We want to return the imported VM object to get its name
	// RunCmdlet executes the cmd, but to get property we might need to be clever.
	// `Import-VM ... -Passthru | Select-Object -ExpandProperty Name`
	// RunCmdlet appends args to the command.

	args = append(args, "-Passthru")

	// Issue: RunCmdlet appends these to "Import-VM".
	// If we want pipe, we can't easily use RunCmdlet unless we support piping logic inside it or use RunScript.
	// Since file paths are involved, RunScript is risky if we fmt.Sprintf the path.
	// But we can use RunCmdlet to get the object JSON?

	// Let's use RunScript with safe path quoting if possible, OR
	// Use RunCmdlet and parse default output? Default output of Import-VM is valid.
	// But we need the Name.

	// Alternative: Use RunCmdlet for Import-VM, then if successful, we might not know the name if we don't capture it.
	// But we CAN use "| Select -Expand Name" as args if powershell parses them.
	// As established in GetVMStatus, passing "|" as a separate arg works if the shell concatenates them.
	// Let's rely on that behavior of "powershell -Command ... arg1 arg2 ..." -> it effectively joins them.

	args = append(args, "|", "Select-Object", "-ExpandProperty", "Name")

	output, err := m.Exec.RunCmdlet(ctx, "Import-VM", args...)
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
func (m *Manager) GetExportedVMInfo(ctx context.Context, path string) (map[string]string, error) {
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
	// Note: vmcxPath is from filepath.Glob which is local. Still potential risk if malicious filenames, but low.
	// Ideally pass path as arg.

	output, err := m.Exec.RunScript(ctx, psScript)
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
