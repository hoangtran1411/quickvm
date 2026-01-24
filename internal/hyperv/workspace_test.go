package hyperv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorkspaceStorage(t *testing.T) {
	// Setup: use a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "quickvm-test-ws")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Since GetWorkspaceDir uses os.UserHomeDir, we might want to mock it.
	// But let's check if we can just test the logic via manual Save/Load to specific paths if possible.
	// The current implementation is hardcoded to ~/.quickvm.

	// Quick hack: override HOME for this test
	oldHome := os.Getenv("USERPROFILE") // Windows
	if oldHome == "" {
		oldHome = os.Getenv("HOME")
	}

	os.Setenv("USERPROFILE", tempDir)
	os.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("USERPROFILE", oldHome)
		os.Setenv("HOME", oldHome)
	}()

	ws := &Workspace{
		Name:        "TestDev",
		Description: "Dev lab",
		VMs:         []string{"VM1", "VM2"},
	}

	// Test Save
	if err := SaveWorkspace(ws); err != nil {
		t.Fatalf("SaveWorkspace failed: %v", err)
	}

	// Test Load
	loaded, err := LoadWorkspace("TestDev")
	if err != nil {
		t.Fatalf("LoadWorkspace failed: %v", err)
	}

	if loaded.Name != ws.Name || len(loaded.VMs) != 2 {
		t.Errorf("Loaded workspace mismatch: %+v", loaded)
	}

	// Test List
	names, err := ListWorkspaces()
	if err != nil {
		t.Fatalf("ListWorkspaces failed: %v", err)
	}
	if len(names) != 1 || names[0] != "TestDev" {
		t.Errorf("ListWorkspaces mismatch: %v", names)
	}

	// Test Delete
	if err := DeleteWorkspace("TestDev"); err != nil {
		t.Fatalf("DeleteWorkspace failed: %v", err)
	}

	_, err = LoadWorkspace("TestDev")
	if err == nil {
		t.Error("Expected error loading deleted workspace, got nil")
	}
}

func TestLoadWorkspace_Errors(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "quickvm-test-load-err")
	defer os.RemoveAll(tempDir)

	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", tempDir)
	defer os.Setenv("USERPROFILE", oldHome)

	// Test non-existent
	_, err := LoadWorkspace("NoSuchWorkspace")
	if err == nil {
		t.Error("Expected error for non-existent workspace, got nil")
	}

	// Test invalid YAML
	wsDir, _ := GetWorkspaceDir()
	invalidFile := filepath.Join(wsDir, "invalid.yaml")
	_ = os.WriteFile(invalidFile, []byte("invalid: yaml: : :"), 0644)

	_, err = LoadWorkspace("invalid")
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestDeleteWorkspace_NonExistent(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "quickvm-test-del-err")
	defer os.RemoveAll(tempDir)

	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", tempDir)
	defer os.Setenv("USERPROFILE", oldHome)

	err := DeleteWorkspace("NoSuchWs")
	if err == nil {
		t.Error("Expected error deleting non-existent workspace, got nil")
	}
}

func TestSaveWorkspace_Error(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "quickvm-test-save-err")
	defer os.RemoveAll(tempDir)

	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", tempDir)
	defer os.Setenv("USERPROFILE", oldHome)

	wsDir, _ := GetWorkspaceDir()
	// Make workspace directory read-only (not reliably working on all OS, but try)
	_ = os.Chmod(wsDir, 0100)
	defer func() { _ = os.Chmod(wsDir, 0755) }()

	ws := &Workspace{Name: "Fail"}
	err := SaveWorkspace(ws)
	if err == nil {
		// Note: On Windows, os.Chmod(dir, 0100) might not prevent file creation
		// So we only log if it actually failed as expected.
		t.Log("SaveWorkspace did not fail on read-only dir (common on Windows)")
	}
}

func TestGetWorkspaceDir_HomeError(t *testing.T) {
	oldHome := os.Getenv("USERPROFILE")
	os.Unsetenv("USERPROFILE")
	os.Unsetenv("HOME")
	defer func() {
		os.Setenv("USERPROFILE", oldHome)
	}()

	_, err := GetWorkspaceDir()
	if err == nil {
		t.Error("Expected error when home dir is unknown, got nil")
	}
}

func TestGetWorkspaceDir(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "quickvm-test-dir")
	defer os.RemoveAll(tempDir)

	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", tempDir)
	defer os.Setenv("USERPROFILE", oldHome)

	dir, err := GetWorkspaceDir()
	if err != nil {
		t.Fatalf("GetWorkspaceDir failed: %v", err)
	}

	expected := filepath.Join(tempDir, ".quickvm", "workspaces")
	if dir != expected {
		t.Errorf("Expected %s, got %s", expected, dir)
	}

	// Verify directory was created
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("Workspace directory was not created")
	}
}
