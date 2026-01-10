package hyperv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CloneVMOptions contains options for cloning a VM
type CloneVMOptions struct {
	VMIndex int    // Source VM index
	NewName string // Name for the cloned VM
}

// CloneVM clones a VM by index with a new name (full clone)
// This performs: Export -> Import with Copy and GenerateNewId -> Rename -> Cleanup
func (m *Manager) CloneVM(vmIndex int, newName string) error {
	// Validate new name
	if strings.TrimSpace(newName) == "" {
		return fmt.Errorf("new VM name cannot be empty")
	}

	// Get VM name by index
	vmName, err := m.GetVMNameByIndex(vmIndex)
	if err != nil {
		return err
	}

	// Check if new name already exists
	exists, err := m.VMExists(newName)
	if err != nil {
		return fmt.Errorf("failed to check if VM name exists: %v", err)
	}
	if exists {
		return fmt.Errorf("a VM with name '%s' already exists", newName)
	}

	// Create temp directory for export
	tempDir, err := os.MkdirTemp("", "quickvm-clone-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Cleanup on exit

	// Step 1: Export the source VM
	if err := m.ExportVMByName(vmName, tempDir); err != nil {
		return fmt.Errorf("failed to export source VM: %v", err)
	}

	// Step 2: Import with Copy and GenerateNewId
	exportedPath := filepath.Join(tempDir, vmName)
	importOpts := ImportVMOptions{
		Path:          exportedPath,
		Copy:          true, // Full clone - copy files
		GenerateNewID: true, // Generate new VM ID
	}

	importedName, err := m.ImportVM(importOpts)
	if err != nil {
		return fmt.Errorf("failed to import cloned VM: %v", err)
	}

	// Step 3: Rename to the new name if different
	if importedName != newName {
		if err := m.RenameVM(importedName, newName); err != nil {
			// Try to cleanup the imported VM if rename fails
			_ = m.DeleteVM(importedName)
			return fmt.Errorf("failed to rename cloned VM: %v", err)
		}
	}

	return nil
}

// CloneVMByName clones a VM by name with a new name (full clone)
func (m *Manager) CloneVMByName(sourceName, newName string) error {
	// Get index from name
	vms, err := m.GetVMs()
	if err != nil {
		return err
	}

	for _, vm := range vms {
		if vm.Name == sourceName {
			return m.CloneVM(vm.Index, newName)
		}
	}

	return fmt.Errorf("VM '%s' not found", sourceName)
}

// RenameVM renames a VM
func (m *Manager) RenameVM(oldName, newName string) error {
	psScript := fmt.Sprintf(`Rename-VM -Name "%s" -NewName "%s"`, oldName, newName)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to rename VM from '%s' to '%s': %v\nOutput: %s", oldName, newName, err, string(output))
	}
	return nil
}

// VMExists checks if a VM with the given name exists
func (m *Manager) VMExists(name string) (bool, error) {
	psScript := fmt.Sprintf(`(Get-VM -Name "%s" -ErrorAction SilentlyContinue) -ne $null`, name)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to check VM existence: %v\nOutput: %s", err, string(output))
	}

	result := strings.TrimSpace(string(output))
	return strings.EqualFold(result, "true"), nil
}

// DeleteVM deletes a VM by name (used for cleanup on error)
func (m *Manager) DeleteVM(name string) error {
	psScript := fmt.Sprintf(`Remove-VM -Name "%s" -Force`, name)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete VM '%s': %v\nOutput: %s", name, err, string(output))
	}
	return nil
}
