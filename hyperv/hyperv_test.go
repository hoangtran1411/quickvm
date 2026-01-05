package hyperv

import (
	"os"
	"strings"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Error("NewManager() returned nil")
	}
}

func TestGetVMs(t *testing.T) {
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
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping VM index validation in CI/CD environment")
		return
	}

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

// Benchmark for GetVMs
func BenchmarkGetVMs(b *testing.B) {
	manager := NewManager()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.GetVMs()
	}
}
