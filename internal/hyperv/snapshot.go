package hyperv

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// Snapshot represents a Hyper-V VM snapshot/checkpoint
type Snapshot struct {
	Name         string `json:"name"`
	VMName       string `json:"vmName"`
	CreationTime string `json:"creationTime"`
	ParentName   string `json:"parentName"`
	SnapshotType string `json:"snapshotType"`
}

// GetSnapshots retrieves all snapshots for a VM by index
func (m *Manager) GetSnapshots(ctx context.Context, vmIndex int) ([]Snapshot, error) {
	vms, err := m.GetVMs(ctx)
	if err != nil {
		return nil, err
	}

	if vmIndex < 1 || vmIndex > len(vms) {
		return nil, fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", vmIndex, len(vms))
	}

	vm := vms[vmIndex-1]
	return m.GetSnapshotsByVMName(ctx, vm.Name)
}

// GetSnapshotsByVMName retrieves all snapshots for a VM by name
func (m *Manager) GetSnapshotsByVMName(ctx context.Context, vmName string) ([]Snapshot, error) {
	psScript := fmt.Sprintf(`
		$snapshots = Get-VMSnapshot -VMName "%s" -ErrorAction SilentlyContinue
		if ($snapshots) {
			$snapshots | Select-Object @{Name='Name';Expression={$_.Name}},
				@{Name='VMName';Expression={$_.VMName}},
				@{Name='CreationTime';Expression={$_.CreationTime.ToString("yyyy-MM-dd HH:mm:ss")}},
				@{Name='ParentName';Expression={if($_.ParentSnapshotName){$_.ParentSnapshotName}else{"(None)"}}},
				@{Name='SnapshotType';Expression={$_.SnapshotType.ToString()}} | ConvertTo-Json
		} else {
			Write-Output "[]"
		}
	`, vmName)

	output, err := m.Exec.RunScript(ctx, psScript)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots for VM '%s': %v\nOutput: %s", vmName, err, string(output))
	}

	outputStr := strings.TrimSpace(string(output))

	// Handle empty case
	if outputStr == "" || outputStr == "[]" {
		return []Snapshot{}, nil
	}

	var snapshots []Snapshot

	// Handle single snapshot case (PowerShell returns object, not array)
	if strings.HasPrefix(outputStr, "{") {
		var snapshot Snapshot
		if err := json.Unmarshal([]byte(outputStr), &snapshot); err != nil {
			return nil, fmt.Errorf("failed to parse snapshot data: %v", err)
		}
		snapshots = append(snapshots, snapshot)
	} else if strings.HasPrefix(outputStr, "[") {
		if err := json.Unmarshal([]byte(outputStr), &snapshots); err != nil {
			return nil, fmt.Errorf("failed to parse snapshots data: %v", err)
		}
	}

	return snapshots, nil
}

// CreateSnapshot creates a new snapshot for a VM by index
func (m *Manager) CreateSnapshot(ctx context.Context, vmIndex int, snapshotName string) error {
	vms, err := m.GetVMs(ctx)
	if err != nil {
		return err
	}

	if vmIndex < 1 || vmIndex > len(vms) {
		return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", vmIndex, len(vms))
	}

	vm := vms[vmIndex-1]
	return m.CreateSnapshotByVMName(ctx, vm.Name, snapshotName)
}

// CreateSnapshotByVMName creates a new snapshot for a VM by name
func (m *Manager) CreateSnapshotByVMName(ctx context.Context, vmName, snapshotName string) error {
	output, err := m.Exec.RunCmdlet(ctx, "Checkpoint-VM", "-Name", vmName, "-SnapshotName", snapshotName)
	if err != nil {
		return fmt.Errorf("failed to create snapshot '%s' for VM '%s': %v\nOutput: %s", snapshotName, vmName, err, string(output))
	}
	return nil
}

// RestoreSnapshot restores a VM to a specific snapshot by index
func (m *Manager) RestoreSnapshot(ctx context.Context, vmIndex int, snapshotName string) error {
	vms, err := m.GetVMs(ctx)
	if err != nil {
		return err
	}

	if vmIndex < 1 || vmIndex > len(vms) {
		return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", vmIndex, len(vms))
	}

	vm := vms[vmIndex-1]
	return m.RestoreSnapshotByVMName(ctx, vm.Name, snapshotName)
}

// RestoreSnapshotByVMName restores a VM to a specific snapshot by name
func (m *Manager) RestoreSnapshotByVMName(ctx context.Context, vmName, snapshotName string) error {
	output, err := m.Exec.RunCmdlet(ctx, "Restore-VMSnapshot", "-VMName", vmName, "-Name", snapshotName, "-Confirm:$false")
	if err != nil {
		return fmt.Errorf("failed to restore snapshot '%s' for VM '%s': %v\nOutput: %s", snapshotName, vmName, err, string(output))
	}
	return nil
}

// DeleteSnapshot deletes a snapshot from a VM by index
func (m *Manager) DeleteSnapshot(ctx context.Context, vmIndex int, snapshotName string) error {
	vms, err := m.GetVMs(ctx)
	if err != nil {
		return err
	}

	if vmIndex < 1 || vmIndex > len(vms) {
		return fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", vmIndex, len(vms))
	}

	vm := vms[vmIndex-1]
	return m.DeleteSnapshotByVMName(ctx, vm.Name, snapshotName)
}

// DeleteSnapshotByVMName deletes a snapshot from a VM by name
func (m *Manager) DeleteSnapshotByVMName(ctx context.Context, vmName, snapshotName string) error {
	output, err := m.Exec.RunCmdlet(ctx, "Remove-VMSnapshot", "-VMName", vmName, "-Name", snapshotName, "-Confirm:$false")
	if err != nil {
		return fmt.Errorf("failed to delete snapshot '%s' from VM '%s': %v\nOutput: %s", snapshotName, vmName, err, string(output))
	}
	return nil
}

// GetVMNameByIndex returns the VM name for a given index
func (m *Manager) GetVMNameByIndex(ctx context.Context, vmIndex int) (string, error) {
	vms, err := m.GetVMs(ctx)
	if err != nil {
		return "", err
	}

	if vmIndex < 1 || vmIndex > len(vms) {
		return "", fmt.Errorf("invalid VM index: %d (valid range: 1-%d)", vmIndex, len(vms))
	}

	return vms[vmIndex-1].Name, nil
}
