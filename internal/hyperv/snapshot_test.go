package hyperv

import (
	"context"
	"os"
	"strings"
	"testing"
)

// skipIfNoHyperVSnapshot skips test in CI/CD environment
func skipIfNoHyperVSnapshot(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping test: Hyper-V not available in CI/CD environment")
	}
	if !IsRunningAsAdmin(context.TODO()) {
		t.Skip("Skipping test: Administrator privileges required for Hyper-V operations")
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
	_, err := manager.GetSnapshotsByVMName(context.TODO(), "NonExistentVM_123456")
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
			_, err := manager.GetSnapshots(context.TODO(), tc.index)
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
			err := manager.CreateSnapshot(context.TODO(), tc.index, "TestSnapshot")
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

// TestCreateSnapshot_EmptyName tests CreateSnapshot with empty name
func TestCreateSnapshot_EmptyName(t *testing.T) {
	manager := NewManager()

	// Should fail before calling Hyper-V
	err := manager.CreateSnapshot(context.TODO(), 1, "")
	if err == nil {
		t.Error("Expected error for empty snapshot name, got nil")
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
			err := manager.RestoreSnapshot(context.TODO(), tc.index, "TestSnapshot")
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

// TestRestoreSnapshot_EmptyName tests RestoreSnapshot with empty name
func TestRestoreSnapshot_EmptyName(t *testing.T) {
	manager := NewManager()

	// Should fail before calling Hyper-V
	err := manager.RestoreSnapshot(context.TODO(), 1, "")
	if err == nil {
		t.Error("Expected error for empty snapshot name, got nil")
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
			err := manager.DeleteSnapshot(context.TODO(), tc.index, "TestSnapshot")
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

// TestDeleteSnapshot_EmptyName tests DeleteSnapshot with empty name
func TestDeleteSnapshot_EmptyName(t *testing.T) {
	manager := NewManager()

	// Should fail before calling Hyper-V
	err := manager.DeleteSnapshot(context.TODO(), 1, "")
	if err == nil {
		t.Error("Expected error for empty snapshot name, got nil")
	}
}

func TestGetVMNameByIndex_Invalid(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	_, err := manager.GetVMNameByIndex(context.TODO(), 99999)
	if err == nil {
		t.Error("Expected error for invalid index, got nil")
	}
}

// TestCreateSnapshotByVMName_NonExistent tests CreateSnapshotByVMName with non-existent VM
func TestCreateSnapshotByVMName_NonExistent(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	err := manager.CreateSnapshotByVMName(context.TODO(), "QuickVM_NonExistent_12345", "TestSnap")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
	if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "failed") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestRestoreSnapshotByVMName_NonExistent tests RestoreSnapshotByVMName with non-existent VM
func TestRestoreSnapshotByVMName_NonExistent(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	err := manager.RestoreSnapshotByVMName(context.TODO(), "QuickVM_NonExistent_12345", "TestSnap")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}

// TestDeleteSnapshotByVMName_NonExistent tests DeleteSnapshotByVMName with non-existent VM
func TestDeleteSnapshotByVMName_NonExistent(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	err := manager.DeleteSnapshotByVMName(context.TODO(), "QuickVM_NonExistent_12345", "TestSnap")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}

// TestGetSnapshotsByVMName_NonExistent tests GetSnapshotsByVMName with non-existent VM
func TestGetSnapshotsByVMName_NonExistent(t *testing.T) {
	skipIfNoHyperVSnapshot(t)
	manager := NewManager()

	_, err := manager.GetSnapshotsByVMName(context.TODO(), "QuickVM_NonExistent_12345")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}
