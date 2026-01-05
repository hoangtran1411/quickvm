package hyperv

import (
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
	// In a real scenario, we would mock the PowerShell execution
	vms, err := manager.GetVMs()
	
	// We don't fail if there are no VMs
	if err != nil && len(vms) == 0 {
		t.Logf("No VMs found or Hyper-V not available: %v", err)
		return
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
	manager := NewManager()
	
	tests := []struct {
		name    string
		index   int
		wantErr bool
	}{
		{"Valid index", 1, false},
		{"Zero index", 0, true},
		{"Negative index", -1, true},
		{"Large index", 9999, true},
	}
	
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
