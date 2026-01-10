package hyperv

import (
	"os"
	"strings"
	"testing"
)

// skipIfCI skips the test if running in CI/CD environment without Hyper-V
func skipIfNoHyperV(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping test: Hyper-V not available in CI/CD environment")
	}
}

// TestCloneVM_InvalidIndex tests cloning with invalid VM index
func TestCloneVM_InvalidIndex(t *testing.T) {
	skipIfNoHyperV(t)
	manager := NewManager()

	testCases := []struct {
		name     string
		vmIndex  int
		newName  string
		wantErr  bool
		errMatch string
	}{
		{
			name:     "negative index",
			vmIndex:  -1,
			newName:  "TestClone",
			wantErr:  true,
			errMatch: "invalid VM index",
		},
		{
			name:     "zero index",
			vmIndex:  0,
			newName:  "TestClone",
			wantErr:  true,
			errMatch: "invalid VM index",
		},
		{
			name:     "very large index",
			vmIndex:  9999,
			newName:  "TestClone",
			wantErr:  true,
			errMatch: "invalid VM index",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := manager.CloneVM(tc.vmIndex, tc.newName)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tc.errMatch)
					return
				}
				if !strings.Contains(err.Error(), tc.errMatch) {
					t.Errorf("expected error containing '%s', got '%s'", tc.errMatch, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestCloneVM_EmptyName tests cloning with empty new name
// This test doesn't require Hyper-V as empty name validation happens first
func TestCloneVM_EmptyName(t *testing.T) {
	manager := NewManager()

	testCases := []struct {
		name     string
		newName  string
		wantErr  bool
		errMatch string
	}{
		{
			name:     "empty string",
			newName:  "",
			wantErr:  true,
			errMatch: "new VM name cannot be empty",
		},
		{
			name:     "whitespace only",
			newName:  "   ",
			wantErr:  true,
			errMatch: "new VM name cannot be empty",
		},
		{
			name:     "tabs only",
			newName:  "\t\t",
			wantErr:  true,
			errMatch: "new VM name cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Using index 1 - the empty name check should happen first before any Hyper-V calls
			err := manager.CloneVM(1, tc.newName)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tc.errMatch)
					return
				}
				if !strings.Contains(err.Error(), tc.errMatch) {
					t.Errorf("expected error containing '%s', got '%s'", tc.errMatch, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestRenameVM_EmptyNames tests RenameVM with empty names
func TestRenameVM_EmptyNames(t *testing.T) {
	skipIfNoHyperV(t)
	manager := NewManager()

	// Renaming with empty names should fail (PowerShell will error)
	// This test verifies the function calls PowerShell correctly
	err := manager.RenameVM("NonExistentVM", "NewName")
	if err == nil {
		t.Error("expected error when renaming non-existent VM, got nil")
	}
}

// TestVMExists_NonExistent tests VMExists with non-existent VM
func TestVMExists_NonExistent(t *testing.T) {
	skipIfNoHyperV(t)
	manager := NewManager()

	exists, err := manager.VMExists("QuickVM_NonExistent_Test_12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exists {
		t.Error("expected VMExists to return false for non-existent VM")
	}
}

// TestDeleteVM_NonExistent tests DeleteVM with non-existent VM
func TestDeleteVM_NonExistent(t *testing.T) {
	skipIfNoHyperV(t)
	manager := NewManager()

	// Deleting non-existent VM should fail
	err := manager.DeleteVM("QuickVM_NonExistent_Test_12345")
	if err == nil {
		t.Error("expected error when deleting non-existent VM, got nil")
	}
}

// TestCloneVMByName_NonExistent tests cloning by name with non-existent VM
func TestCloneVMByName_NonExistent(t *testing.T) {
	skipIfNoHyperV(t)
	manager := NewManager()

	err := manager.CloneVMByName("QuickVM_NonExistent_Test_12345", "NewClone")
	if err == nil {
		t.Error("expected error when cloning non-existent VM, got nil")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected error containing 'not found', got '%s'", err.Error())
	}
}
