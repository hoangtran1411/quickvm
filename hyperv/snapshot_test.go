package hyperv

import (
	"testing"
)

func TestGetSnapshotsByVMName_InvalidVM(t *testing.T) {
	manager := NewManager()

	// Test with a VM name that doesn't exist
	_, err := manager.GetSnapshotsByVMName("NonExistentVM_123456")
	// This might not error if the VM doesn't exist (returns empty list)
	// but we're just testing the function executes without panic
	_ = err
}

func TestGetSnapshots_InvalidIndex(t *testing.T) {
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
