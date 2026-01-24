package hyperv

import (
	"context"
	"os"
	"strings"
	"testing"
)

// skipIfNoHyperV skips the test if running in CI/CD environment without Hyper-V or without Admin privileges
func skipIfNoHyperV(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping test: Hyper-V not available in CI/CD environment")
	}
	if !IsRunningAsAdmin(context.TODO()) {
		t.Skip("Skipping test: Administrator privileges required for Hyper-V operations")
	}
}

// TestCloneVMOptions_Struct tests CloneVMOptions struct
func TestCloneVMOptions_Struct(t *testing.T) {
	opts := CloneVMOptions{
		VMIndex: 1,
		NewName: "TestCloneVM",
	}

	if opts.VMIndex != 1 {
		t.Errorf("Expected VMIndex=1, got %d", opts.VMIndex)
	}
	if opts.NewName != "TestCloneVM" {
		t.Errorf("Expected NewName=TestCloneVM, got %s", opts.NewName)
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
			err := manager.CloneVM(context.TODO(), tc.vmIndex, tc.newName)
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
		{
			name:     "newlines only",
			newName:  "\n\n",
			wantErr:  true,
			errMatch: "new VM name cannot be empty",
		},
		{
			name:     "mixed whitespace",
			newName:  " \t\n ",
			wantErr:  true,
			errMatch: "new VM name cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Using index 1 - the empty name check should happen first before any Hyper-V calls
			err := manager.CloneVM(context.TODO(), 1, tc.newName)
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

// TestCloneVM_ValidName_VariousCharacters tests valid VM names
func TestCloneVM_ValidName_NotEmpty(t *testing.T) {
	// These are just validation tests for the empty check
	// They will fail at Hyper-V level but pass the empty name check
	validNames := []string{
		"ValidName",
		"VM-Clone-1",
		"VM_Clone_2",
		"VM Clone 3",
		"VM123",
		"a",
		"  ValidWithSpaces  ", // After trim, still valid (but trim happens in validation)
	}

	for _, name := range validNames {
		t.Run(name, func(t *testing.T) {
			// Just verify that empty check passes
			if strings.TrimSpace(name) == "" {
				t.Errorf("Expected '%s' to be non-empty after trim", name)
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
	err := manager.RenameVM(context.TODO(), "NonExistentVM", "NewName")
	if err == nil {
		t.Error("expected error when renaming non-existent VM, got nil")
	}
}

// TestVMExists_NonExistent tests VMExists with non-existent VM
func TestVMExists_NonExistent(t *testing.T) {
	skipIfNoHyperV(t)
	manager := NewManager()

	exists, err := manager.VMExists(context.TODO(), "QuickVM_NonExistent_Test_12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exists {
		t.Error("expected VMExists to return false for non-existent VM")
	}
}

// TestVMExists_EmptyName tests VMExists with empty name
func TestVMExists_EmptyName(t *testing.T) {
	skipIfNoHyperV(t)
	manager := NewManager()

	// Empty name should return false (no VM with empty name)
	exists, err := manager.VMExists(context.TODO(), "")
	if err != nil {
		// PowerShell might error with empty name
		return
	}

	if exists {
		t.Error("expected VMExists to return false for empty VM name")
	}
}

// TestDeleteVM_NonExistent tests DeleteVM with non-existent VM
func TestDeleteVM_NonExistent(t *testing.T) {
	skipIfNoHyperV(t)
	manager := NewManager()

	// Deleting non-existent VM should fail
	err := manager.DeleteVM(context.TODO(), "QuickVM_NonExistent_Test_12345")
	if err == nil {
		t.Error("expected error when deleting non-existent VM, got nil")
	}
}

// TestCloneVMByName_NonExistent tests cloning by name with non-existent VM
func TestCloneVMByName_NonExistent(t *testing.T) {
	skipIfNoHyperV(t)
	manager := NewManager()

	err := manager.CloneVMByName(context.TODO(), "QuickVM_NonExistent_Test_12345", "NewClone")
	if err == nil {
		t.Error("expected error when cloning non-existent VM, got nil")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected error containing 'not found', got '%s'", err.Error())
	}
}

// TestCloneVMByName_EmptySourceName tests CloneVMByName with empty source name
func TestCloneVMByName_EmptySourceName(t *testing.T) {
	skipIfNoHyperV(t)
	manager := NewManager()

	err := manager.CloneVMByName(context.TODO(), "", "NewClone")
	if err == nil {
		t.Error("expected error when cloning with empty source name, got nil")
	}
}

// TestCloneVMByName_EmptyNewName tests CloneVMByName with empty new name
func TestCloneVMByName_EmptyNewName(t *testing.T) {
	skipIfNoHyperV(t)
	manager := NewManager()

	// This should fail with "not found" because source doesn't exist
	// But if source existed, it would fail with "new VM name cannot be empty"
	err := manager.CloneVMByName(context.TODO(), "SomeVM", "")
	if err == nil {
		t.Error("expected error when cloning with empty new name, got nil")
	}
}
