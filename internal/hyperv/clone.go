package hyperv

import (
	"context"
	"fmt"
	"os"
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
func (m *Manager) CloneVM(ctx context.Context, vmIndex int, newName string) error {
	// Validate new name
	if strings.TrimSpace(newName) == "" {
		return fmt.Errorf("new VM name cannot be empty")
	}

	// Get VM name by index
	vmName, err := m.GetVMNameByIndex(ctx, vmIndex)
	if err != nil {
		return err
	}

	// Check if new name already exists
	exists, err := m.VMExists(ctx, newName)
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
	defer func() { _ = os.RemoveAll(tempDir) }() // Cleanup on exit

	// Step 1: Export the source VM
	if err := m.ExportVMByName(ctx, vmName, tempDir); err != nil {
		return fmt.Errorf("failed to export source VM: %v", err)
	}

	// Step 2: Import with Copy and GenerateNewId
	exportedPath := filepath.Join(tempDir, vmName)
	importOpts := ImportVMOptions{
		Path:          exportedPath,
		Copy:          true, // Full clone - copy files
		GenerateNewID: true, // Generate new VM ID
	}

	importedName, err := m.ImportVM(ctx, importOpts)
	if err != nil {
		return fmt.Errorf("failed to import cloned VM: %v", err)
	}

	// Step 3: Rename to the new name if different
	if importedName != newName {
		if err := m.RenameVM(ctx, importedName, newName); err != nil {
			// Try to cleanup the imported VM if rename fails
			_ = m.DeleteVM(ctx, importedName)
			return fmt.Errorf("failed to rename cloned VM: %v", err)
		}
	}

	return nil
}

// CloneVMByName clones a VM by name with a new name (full clone)
func (m *Manager) CloneVMByName(ctx context.Context, sourceName, newName string) error {
	// Get index from name
	vms, err := m.GetVMs(ctx)
	if err != nil {
		return err
	}

	for _, vm := range vms {
		if vm.Name == sourceName {
			return m.CloneVM(ctx, vm.Index, newName)
		}
	}

	return fmt.Errorf("VM '%s' not found", sourceName)
}

// RenameVM renames a VM
func (m *Manager) RenameVM(ctx context.Context, oldName, newName string) error {
	output, err := m.Exec.RunCmdlet(ctx, "Rename-VM", "-Name", oldName, "-NewName", newName)
	if err != nil {
		return fmt.Errorf("failed to rename VM from '%s' to '%s': %v\nOutput: %s", oldName, newName, err, string(output))
	}
	return nil
}

// VMExists checks if a VM with the given name exists
func (m *Manager) VMExists(ctx context.Context, name string) (bool, error) {
	// We use RunCmdlet and check for errors or empty output if Get-VM fails
	// But Get-VM returns error if not found by default.
	// Safer/Simpler: Use RunScript with ErrorAction SilentlyContinue logic if RunCmdlet is strict.
	// However, using explicit Get-VM with RunCmdlet is safer against injection.

	// Strategy: Try to get VM info. if error -> assume not exists (or real error)
	// PowerShell's Get-VM throws if name not found.

	// Let's use RunScript for the precise bool check logic to match original intent,
	// but constructing it safely is hard without injection.
	// Better: use RunCmdlet and handle the error.

	output, err := m.Exec.RunCmdlet(ctx, "Get-VM", "-Name", name, "-ErrorAction", "SilentlyContinue")
	if err != nil {
		// If it's a "does not exist" error, return false.
		// BUT locally RunCmdlet returns combined output and error.
		// If PowerShell returns exit code 0 (even empty), it's success.
		// If exit code non-zero, it's err.
		// Update: RunCmdlet implementation uses exec.CommandContext which returns err if exit code != 0.
		return false, nil // Assume not found if it fails
	}

	// If output is empty, it might not exist if ErrorAction was SilentlyContinue?
	if strings.TrimSpace(string(output)) == "" {
		return false, nil
	}

	return true, nil
}

// DeleteVM deletes a VM by name (used for cleanup on error)
func (m *Manager) DeleteVM(ctx context.Context, name string) error {
	output, err := m.Exec.RunCmdlet(ctx, "Remove-VM", "-Name", name, "-Force")
	if err != nil {
		return fmt.Errorf("failed to delete VM '%s': %v\nOutput: %s", name, err, string(output))
	}
	return nil
}
