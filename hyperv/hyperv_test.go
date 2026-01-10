package hyperv

import (
	"os"
	"strings"
	"testing"
)

// skipIfNoHyperVMain skips test in CI/CD environment
func skipIfNoHyperVMain(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping test: Hyper-V not available in CI/CD environment")
	}
}

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Error("NewManager() returned nil")
	}
}

// TestVM_Struct tests VM struct fields
func TestVM_Struct(t *testing.T) {
	vm := VM{
		Index:    1,
		Name:     "TestVM",
		State:    "Running",
		CPUUsage: 25,
		MemoryMB: 4096,
		Uptime:   "1.02:30:45",
		Status:   "Operating normally",
		Version:  "9.0",
	}

	if vm.Index != 1 {
		t.Errorf("Expected Index=1, got %d", vm.Index)
	}
	if vm.Name != "TestVM" {
		t.Errorf("Expected Name=TestVM, got %s", vm.Name)
	}
	if vm.State != "Running" {
		t.Errorf("Expected State=Running, got %s", vm.State)
	}
	if vm.CPUUsage != 25 {
		t.Errorf("Expected CPUUsage=25, got %d", vm.CPUUsage)
	}
	if vm.MemoryMB != 4096 {
		t.Errorf("Expected MemoryMB=4096, got %d", vm.MemoryMB)
	}
	if vm.Uptime != "1.02:30:45" {
		t.Errorf("Expected Uptime=1.02:30:45, got %s", vm.Uptime)
	}
	if vm.Status != "Operating normally" {
		t.Errorf("Expected Status=Operating normally, got %s", vm.Status)
	}
	if vm.Version != "9.0" {
		t.Errorf("Expected Version=9.0, got %s", vm.Version)
	}
}

// TestManager_Struct tests Manager struct
func TestManager_Struct(t *testing.T) {
	manager := &Manager{}
	if manager == nil {
		t.Error("Manager struct is nil")
	}
}

func TestGetVMs(t *testing.T) {
	skipIfNoHyperVMain(t)
	manager := NewManager()

	// Note: This test requires actual Hyper-V to be running
	// In CI/CD environments without Hyper-V, this will gracefully skip
	vms, err := manager.GetVMs()

	// If no VMs or Hyper-V not available, skip the test
	if err != nil && len(vms) == 0 {
		if strings.Contains(err.Error(), "no VMs found") ||
			strings.Contains(err.Error(), "invalid output") ||
			strings.Contains(err.Error(), "failed to execute") {
			t.Skip("Skipping test: Hyper-V not available or no VMs configured")
			return
		}
		t.Fatalf("Unexpected error: %v", err)
	}

	// If we got VMs, verify structure
	if len(vms) > 0 {
		vm := vms[0]
		if vm.Name == "" {
			t.Error("VM Name should not be empty")
		}
		if vm.Index != 1 {
			t.Errorf("First VM Index should be 1, got %d", vm.Index)
		}
	}
}

func TestVMIndexValidation(t *testing.T) {
	// Skip this test if running in CI/CD without Hyper-V
	skipIfNoHyperVMain(t)

	manager := NewManager()

	// First check if VMs are available
	vms, err := manager.GetVMs()
	if err != nil || len(vms) == 0 {
		t.Skip("Skipping test: No VMs available for testing")
		return
	}

	tests := []struct {
		name    string
		index   int
		wantErr bool
	}{
		{"Zero index", 0, true},
		{"Negative index", -1, true},
		{"Large index", 9999, true},
	}

	// Only test invalid indices to avoid starting actual VMs
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.StartVM(tt.index)
			hasErr := err != nil

			if hasErr != tt.wantErr {
				t.Errorf("StartVM(%d) error = %v, wantErr %v", tt.index, err, tt.wantErr)
			}
		})
	}
}

// TestStartVMByName_NonExistent tests StartVMByName with non-existent VM
func TestStartVMByName_NonExistent(t *testing.T) {
	skipIfNoHyperVMain(t)
	manager := NewManager()

	err := manager.StartVMByName("QuickVM_NonExistent_12345")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}

// TestStopVMByName_NonExistent tests StopVMByName with non-existent VM
func TestStopVMByName_NonExistent(t *testing.T) {
	skipIfNoHyperVMain(t)
	manager := NewManager()

	err := manager.StopVMByName("QuickVM_NonExistent_12345")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}

// TestRestartVMByName_NonExistent tests RestartVMByName with non-existent VM
func TestRestartVMByName_NonExistent(t *testing.T) {
	skipIfNoHyperVMain(t)
	manager := NewManager()

	err := manager.RestartVMByName("QuickVM_NonExistent_12345")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}

// TestGetVMStatus_NonExistent tests GetVMStatus with non-existent VM
func TestGetVMStatus_NonExistent(t *testing.T) {
	skipIfNoHyperVMain(t)
	manager := NewManager()

	status, err := manager.GetVMStatus("QuickVM_NonExistent_12345")
	// PowerShell may return empty status without error for non-existent VM
	// We just verify the function doesn't panic and returns something expected
	if err != nil {
		// Error is acceptable for non-existent VM
		return
	}
	// If no error, status should be empty or contain unexpected output
	if status != "" {
		t.Logf("Status for non-existent VM: %q (expected empty or error)", status)
	}
}

// TestStartVM_InvalidIndex tests StartVM with invalid indices
func TestStartVM_InvalidIndex(t *testing.T) {
	skipIfNoHyperVMain(t)
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
			err := manager.StartVM(tc.index)
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

// TestStopVM_InvalidIndex tests StopVM with invalid indices
func TestStopVM_InvalidIndex(t *testing.T) {
	skipIfNoHyperVMain(t)
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
			err := manager.StopVM(tc.index)
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

// TestRestartVM_InvalidIndex tests RestartVM with invalid indices
func TestRestartVM_InvalidIndex(t *testing.T) {
	skipIfNoHyperVMain(t)
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
			err := manager.RestartVM(tc.index)
			if err == nil {
				t.Errorf("Expected error for index %d, got nil", tc.index)
			}
		})
	}
}

// Benchmark for GetVMs
func BenchmarkGetVMs(b *testing.B) {
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		b.Skip("Skipping benchmark in CI/CD environment")
	}

	manager := NewManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.GetVMs()
	}
}
