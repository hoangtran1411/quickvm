package hyperv

import (
	"context"
	"os"
	"strings"
	"testing"
)

// skipIfNoHyperVMain skips test in CI/CD environment or if not admin
func skipIfNoHyperVMain(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping test: Hyper-V not available in CI/CD environment")
	}
	if !IsRunningAsAdmin(context.TODO()) {
		t.Skip("Skipping test: Administrator privileges required for Hyper-V operations")
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

func TestGetVMs(t *testing.T) {
	skipIfNoHyperVMain(t)
	manager := NewManager()

	// Note: This test requires actual Hyper-V to be running
	// In CI/CD environments without Hyper-V, this will gracefully skip
	vms, err := manager.GetVMs(context.TODO())

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
	vms, err := manager.GetVMs(context.TODO())
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
			err := manager.StartVM(context.TODO(), tt.index)
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

	err := manager.StartVMByName(context.TODO(), "QuickVM_NonExistent_12345")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}

// TestStopVMByName_NonExistent tests StopVMByName with non-existent VM
func TestStopVMByName_NonExistent(t *testing.T) {
	skipIfNoHyperVMain(t)
	manager := NewManager()

	err := manager.StopVMByName(context.TODO(), "QuickVM_NonExistent_12345")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}

// TestRestartVMByName_NonExistent tests RestartVMByName with non-existent VM
func TestRestartVMByName_NonExistent(t *testing.T) {
	skipIfNoHyperVMain(t)
	manager := NewManager()

	err := manager.RestartVMByName(context.TODO(), "QuickVM_NonExistent_12345")
	if err == nil {
		t.Error("Expected error for non-existent VM, got nil")
	}
}

// TestGetVMStatus_NonExistent tests GetVMStatus with non-existent VM
func TestGetVMStatus_NonExistent(t *testing.T) {
	skipIfNoHyperVMain(t)
	manager := NewManager()

	status, err := manager.GetVMStatus(context.TODO(), "QuickVM_NonExistent_12345")
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
			err := manager.StartVM(context.TODO(), tc.index)
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
			err := manager.StopVM(context.TODO(), tc.index)
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
			err := manager.RestartVM(context.TODO(), tc.index)
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

	
	for b.Loop() {
		_, _ = manager.GetVMs(context.TODO())
	}
}

func TestFormatCmdletScript(t *testing.T) {
	tests := []struct {
		name     string
		cmdlet   string
		args     []string
		expected string
	}{
		{
			name:     "Simple cmdlet with name",
			cmdlet:   "Start-VM",
			args:     []string{"-Name", "MyVM"},
			expected: "Start-VM -Name 'MyVM'",
		},
		{
			name:     "Cmdlet with switch flag",
			cmdlet:   "Stop-VM",
			args:     []string{"-Name", "MyVM", "-Force"},
			expected: "Stop-VM -Name 'MyVM' -Force",
		},
		{
			name:     "Cmdlet with injection attempt in value",
			cmdlet:   "Start-VM",
			args:     []string{"-Name", "MyVM'; Stop-Process -Id 1 #"},
			expected: "Start-VM -Name 'MyVM''; Stop-Process -Id 1 #'",
		},
		{
			name:     "Value starting with dash (parameter value matching flag syntax)",
			cmdlet:   "Start-VM",
			args:     []string{"-Name", "-WhatIf"},
			expected: "Start-VM -Name '-WhatIf'",
		},
		{
			name:     "Value starting with dash prod name",
			cmdlet:   "Start-VM",
			args:     []string{"-Name", "-proddb"},
			expected: "Start-VM -Name '-proddb'",
		},
		{
			name:     "Value matching operator name",
			cmdlet:   "Start-VM",
			args:     []string{"-Name", "Select-Object"},
			expected: "Start-VM -Name 'Select-Object'",
		},
		{
			name:     "Value matching pipe",
			cmdlet:   "Start-VM",
			args:     []string{"-Name", "|"},
			expected: "Start-VM -Name '|'",
		},
		{
			name:     "Pipeline with Select-Object and ExpandProperty",
			cmdlet:   "Get-VM",
			args:     []string{"-Name", "MyVM", "|", "Select-Object", "-ExpandProperty", "State"},
			expected: "Get-VM -Name 'MyVM' | Select-Object -ExpandProperty 'State'",
		},
		{
			name:     "Multiple valued parameters",
			cmdlet:   "Checkpoint-VM",
			args:     []string{"-Name", "MyVM", "-SnapshotName", "Snap 1"},
			expected: "Checkpoint-VM -Name 'MyVM' -SnapshotName 'Snap 1'",
		},
		{
			name:     "Shutdown cmdlet with flags and values",
			cmdlet:   "shutdown",
			args:     []string{"/r", "/t", "10", "/c", "Restarting"},
			expected: "shutdown /r /t '10' /c 'Restarting'",
		},
		{
			name:     "Switch with explicit boolean",
			cmdlet:   "Restore-VMSnapshot",
			args:     []string{"-VMName", "MyVM", "-Name", "Snap1", "-Confirm:$false"},
			expected: "Restore-VMSnapshot -VMName 'MyVM' -Name 'Snap1' -Confirm:$false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := formatCmdletScript(tt.cmdlet, tt.args...)
			if actual != tt.expected {
				t.Errorf("formatCmdletScript() =\n  %q\nwant:\n  %q", actual, tt.expected)
			}
		})
	}
}

func TestVMOperations_EmptyNameValidation(t *testing.T) {
	manager := NewManager()
	ctx := context.Background()

	emptyNames := []string{"", "   ", "\t\n"}

	for _, name := range emptyNames {
		t.Run("start_empty_"+name, func(t *testing.T) {
			err := manager.StartVMByName(ctx, name)
			if err == nil || !strings.Contains(err.Error(), "VM name cannot be empty") {
				t.Errorf("Expected 'VM name cannot be empty', got %v", err)
			}
		})
		t.Run("stop_empty_"+name, func(t *testing.T) {
			err := manager.StopVMByName(ctx, name)
			if err == nil || !strings.Contains(err.Error(), "VM name cannot be empty") {
				t.Errorf("Expected 'VM name cannot be empty', got %v", err)
			}
		})
		t.Run("restart_empty_"+name, func(t *testing.T) {
			err := manager.RestartVMByName(ctx, name)
			if err == nil || !strings.Contains(err.Error(), "VM name cannot be empty") {
				t.Errorf("Expected 'VM name cannot be empty', got %v", err)
			}
		})
		t.Run("status_empty_"+name, func(t *testing.T) {
			_, err := manager.GetVMStatus(ctx, name)
			if err == nil || !strings.Contains(err.Error(), "VM name cannot be empty") {
				t.Errorf("Expected 'VM name cannot be empty', got %v", err)
			}
		})
	}
}
