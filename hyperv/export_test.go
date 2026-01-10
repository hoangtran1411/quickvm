package hyperv

import (
	"os"
	"path/filepath"
	"testing"
)

// skipIfNoHyperVExport skips test in CI/CD environment
func skipIfNoHyperVExport(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping test: Hyper-V not available in CI/CD environment")
	}
}

func TestExportVM_InvalidIndex(t *testing.T) {
	skipIfNoHyperVExport(t)
	manager := NewManager()

	testCases := []struct {
		name  string
		index int
	}{
		{"Zero index", 0},
		{"Negative index", -1},
		{"Large index", 99999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := manager.ExportVM(tc.index, "C:\\Temp")
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

func TestFindVMCXFile_DirectVMCXPath(t *testing.T) {
	manager := NewManager()

	// Test with direct .vmcx path
	testPath := "C:\\Test\\VM\\Virtual Machines\\test.vmcx"
	result, err := manager.findVMCXFile(testPath)
	if err != nil {
		t.Errorf("Expected no error for direct .vmcx path, got: %v", err)
	}
	if result != testPath {
		t.Errorf("Expected path %s, got %s", testPath, result)
	}
}

func TestFindVMCXFile_NotFound(t *testing.T) {
	manager := NewManager()

	// Create a temp directory without any .vmcx files
	tempDir, err := os.MkdirTemp("", "quickvm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_, err = manager.findVMCXFile(tempDir)
	if err == nil {
		t.Error("Expected error for directory without .vmcx files, got nil")
	}
}

func TestFindVMCXFile_InVirtualMachinesSubdir(t *testing.T) {
	manager := NewManager()

	// Create a temp directory structure with a .vmcx file in "Virtual Machines" subdir
	tempDir, err := os.MkdirTemp("", "quickvm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	vmDir := filepath.Join(tempDir, "Virtual Machines")
	if err := os.MkdirAll(vmDir, 0755); err != nil {
		t.Fatalf("Failed to create Virtual Machines directory: %v", err)
	}

	vmcxFile := filepath.Join(vmDir, "test.vmcx")
	if err := os.WriteFile(vmcxFile, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create .vmcx file: %v", err)
	}

	result, err := manager.findVMCXFile(tempDir)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != vmcxFile {
		t.Errorf("Expected path %s, got %s", vmcxFile, result)
	}
}

func TestImportVMOptions_Defaults(t *testing.T) {
	opts := ImportVMOptions{}

	if opts.Copy != false {
		t.Error("Expected Copy to default to false")
	}
	if opts.GenerateNewID != false {
		t.Error("Expected GenerateNewID to default to false")
	}
	if opts.VHDPath != "" {
		t.Error("Expected VHDPath to default to empty")
	}
}

func TestImportVMOptions_WithFlags(t *testing.T) {
	testCases := []struct {
		name          string
		opts          ImportVMOptions
		expectCopy    bool
		expectNewID   bool
		expectVHDPath string
	}{
		{
			name: "Copy enabled",
			opts: ImportVMOptions{
				Path: "C:\\Test\\VM",
				Copy: true,
			},
			expectCopy:    true,
			expectNewID:   false,
			expectVHDPath: "",
		},
		{
			name: "GenerateNewID enabled",
			opts: ImportVMOptions{
				Path:          "C:\\Test\\VM",
				GenerateNewID: true,
			},
			expectCopy:    false,
			expectNewID:   true,
			expectVHDPath: "",
		},
		{
			name: "VHDPath specified",
			opts: ImportVMOptions{
				Path:    "C:\\Test\\VM",
				VHDPath: "E:\\VHDs",
			},
			expectCopy:    false,
			expectNewID:   false,
			expectVHDPath: "E:\\VHDs",
		},
		{
			name: "All options enabled",
			opts: ImportVMOptions{
				Path:          "C:\\Test\\VM",
				Copy:          true,
				GenerateNewID: true,
				VHDPath:       "E:\\VHDs",
			},
			expectCopy:    true,
			expectNewID:   true,
			expectVHDPath: "E:\\VHDs",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.opts.Copy != tc.expectCopy {
				t.Errorf("Expected Copy=%v, got %v", tc.expectCopy, tc.opts.Copy)
			}
			if tc.opts.GenerateNewID != tc.expectNewID {
				t.Errorf("Expected GenerateNewID=%v, got %v", tc.expectNewID, tc.opts.GenerateNewID)
			}
			if tc.opts.VHDPath != tc.expectVHDPath {
				t.Errorf("Expected VHDPath=%s, got %s", tc.expectVHDPath, tc.opts.VHDPath)
			}
		})
	}
}

func TestImportVM_PathNotExists(t *testing.T) {
	manager := NewManager()

	opts := ImportVMOptions{
		Path: "C:\\NonExistent\\Path\\That\\Does\\Not\\Exist\\12345",
	}

	_, err := manager.ImportVM(opts)
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}
}

func TestImportVM_NoVMCXFile(t *testing.T) {
	manager := NewManager()

	// Create a temp directory without any .vmcx files
	tempDir, err := os.MkdirTemp("", "quickvm-import-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	opts := ImportVMOptions{
		Path: tempDir,
	}

	_, err = manager.ImportVM(opts)
	if err == nil {
		t.Error("Expected error for directory without .vmcx files, got nil")
	}
}

func TestFindVMCXFile_CaseInsensitive(t *testing.T) {
	manager := NewManager()

	// Test uppercase .VMCX extension
	testPath := "C:\\Test\\VM\\Virtual Machines\\test.VMCX"
	result, err := manager.findVMCXFile(testPath)
	if err != nil {
		t.Errorf("Expected no error for .VMCX path, got: %v", err)
	}
	if result != testPath {
		t.Errorf("Expected path %s, got %s", testPath, result)
	}

	// Test mixed case .VmCx extension
	testPath2 := "C:\\Test\\VM\\Virtual Machines\\test.VmCx"
	result2, err := manager.findVMCXFile(testPath2)
	if err != nil {
		t.Errorf("Expected no error for .VmCx path, got: %v", err)
	}
	if result2 != testPath2 {
		t.Errorf("Expected path %s, got %s", testPath2, result2)
	}
}

func TestFindVMCXFile_DirectlyInBasePath(t *testing.T) {
	manager := NewManager()

	// Create a temp directory with a .vmcx file directly in basePath
	tempDir, err := os.MkdirTemp("", "quickvm-test-direct-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	vmcxFile := filepath.Join(tempDir, "test.vmcx")
	if err := os.WriteFile(vmcxFile, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create .vmcx file: %v", err)
	}

	result, err := manager.findVMCXFile(tempDir)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != vmcxFile {
		t.Errorf("Expected path %s, got %s", vmcxFile, result)
	}
}

func TestExportVMByName_EmptyName(t *testing.T) {
	skipIfNoHyperVExport(t)
	manager := NewManager()

	// Test with empty VM name - should fail at PowerShell level
	err := manager.ExportVMByName("", "C:\\Temp")
	if err == nil {
		t.Error("Expected error for empty VM name, got nil")
	}
}

func TestExportVMOptions_Struct(t *testing.T) {
	opts := ExportVMOptions{
		VMIndex: 1,
		Path:    "C:\\Backups\\VMs",
	}

	if opts.VMIndex != 1 {
		t.Errorf("Expected VMIndex=1, got %d", opts.VMIndex)
	}
	if opts.Path != "C:\\Backups\\VMs" {
		t.Errorf("Expected Path=C:\\Backups\\VMs, got %s", opts.Path)
	}
}
