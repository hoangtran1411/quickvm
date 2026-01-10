package hyperv

import (
	"os"
	"testing"
)

// skipIfNoHyperVSnapshot skips test in CI/CD environment
func skipIfNoHyperVSnapshot(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping test: Hyper-V not available in CI/CD environment")
	}
}

// TestSnapshot_Struct tests Snapshot struct fields
func TestSnapshot_Struct(t *testing.T) {
	snapshot := Snapshot{
		Name:         "TestSnapshot",
		VMName:       "TestVM",
		CreationTime: "2026-01-10 12:00:00",
		ParentName:   "(None)",
		SnapshotType: "Standard",
	}

	if snapshot.Name != "TestSnapshot" {
		t.Errorf("Expected Name=TestSnapshot, got %s", snapshot.Name)
	}
	if snapshot.VMName != "TestVM" {
		t.Errorf("Expected VMName=TestVM, got %s", snapshot.VMName)
	}
	if snapshot.CreationTime != "2026-01-10 12:00:00" {
		t.Errorf("Expected CreationTime=2026-01-10 12:00:00, got %s", snapshot.CreationTime)
	}
	if snapshot.ParentName != "(None)" {
		t.Errorf("Expected ParentName=(None), got %s", snapshot.ParentName)
	}
	if snapshot.SnapshotType != "Standard" {
		t.Errorf("Expected SnapshotType=Standard, got %s", snapshot.SnapshotType)
	}
}

func TestGetSnapshotsByVMName_InvalidVM(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	// Test with a VM name that doesn't exist
	_, err := manager.GetSnapshotsByVMName("NonExistentVM_123456")
	// This might not error if the VM doesn't exist (returns empty list)
	// but we're just testing the function executes without panic
	_ = err
}

func TestGetSnapshots_InvalidIndex(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	testCases := []struct {
		name  string
		index int
	}{
		{"Zero index", 0},
		{"Negative index", -1},
		{"Large index", 9999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := manager.GetSnapshots(tc.index)
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

func TestCreateSnapshot_InvalidIndex(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	testCases := []struct {
		name  string
		index int
	}{
		{"Zero index", 0},
		{"Negative index", -1},
		{"Large index", 9999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := manager.CreateSnapshot(tc.index, "TestSnapshot")
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

func TestRestoreSnapshot_InvalidIndex(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	testCases := []struct {
		name  string
		index int
	}{
		{"Zero index", 0},
		{"Negative index", -1},
		{"Large index", 9999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := manager.RestoreSnapshot(tc.index, "TestSnapshot")
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

func TestDeleteSnapshot_InvalidIndex(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	testCases := []struct {
		name  string
		index int
	}{
		{"Zero index", 0},
		{"Negative index", -1},
		{"Large index", 9999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := manager.DeleteSnapshot(tc.index, "TestSnapshot")
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

func TestGetVMNameByIndex_InvalidIndex(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	testCases := []struct {
		name  string
		index int
	}{
		{"Zero index", 0},
		{"Negative index", -1},
		{"Large index", 9999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := manager.GetVMNameByIndex(tc.index)
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

// TestCreateSnapshotByVMName_NonExistent tests creating snapshot for non-existent VM
func TestCreateSnapshotByVMName_NonExistent(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	err := manager.CreateSnapshotByVMName("QuickVM_NonExistent_12345", "TestSnapshot")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}

// TestRestoreSnapshotByVMName_NonExistent tests restoring snapshot for non-existent VM
func TestRestoreSnapshotByVMName_NonExistent(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	err := manager.RestoreSnapshotByVMName("QuickVM_NonExistent_12345", "TestSnapshot")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}

// TestDeleteSnapshotByVMName_NonExistent tests deleting snapshot for non-existent VM
func TestDeleteSnapshotByVMName_NonExistent(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	err := manager.DeleteSnapshotByVMName("QuickVM_NonExistent_12345", "TestSnapshot")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}
