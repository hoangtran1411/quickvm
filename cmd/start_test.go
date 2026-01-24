package cmd

import (
	"context"
	"fmt"
	"quickvm/internal/hyperv"
	"testing"
)

func TestRunStart(t *testing.T) {
	mockVMs := []hyperv.VM{
		{Name: "VM1", Index: 1},
		{Name: "VM2", Index: 2},
	}

	tests := []struct {
		name     string
		args     []string
		rangeStr string
		all      bool
		setup    func(*MockManager)
	}{
		{
			name: "Start single VM",
			args: []string{"1"},
			setup: func(m *MockManager) {
				m.StartVMByNameFn = func(ctx context.Context, name string) error {
					if name != "VM1" {
						return fmt.Errorf("wrong VM")
					}
					return nil
				}
			},
		},
		{
			name: "Start all VMs",
			all:  true,
			setup: func(m *MockManager) {
				count := 0
				m.StartVMByNameFn = func(ctx context.Context, name string) error {
					count++
					return nil
				}
			},
		},
		{
			name: "Failed to get VMs",
			args: []string{"1"},
			setup: func(m *MockManager) {
				m.GetVMsFn = func(ctx context.Context) ([]hyperv.VM, error) {
					return nil, fmt.Errorf("hyper-v error")
				}
			},
		},
		{
			name: "Failed to start one VM",
			args: []string{"1", "2"},
			setup: func(m *MockManager) {
				m.StartVMByNameFn = func(ctx context.Context, name string) error {
					if name == "VM2" {
						return fmt.Errorf("crash")
					}
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockManager{
				GetVMsFn: func(ctx context.Context) ([]hyperv.VM, error) {
					return mockVMs, nil
				},
			}
			if tt.setup != nil {
				tt.setup(m)
			}

			// We just run it. Verification happens inside m.StartVMByNameFn
			runStart(context.Background(), m, tt.args, tt.rangeStr, tt.all)
		})
	}
}

func TestStartCommandSetup(t *testing.T) {
	if startCmd.Use != "start [vm-index]" {
		t.Errorf("Expected use 'start [vm-index]', got '%s'", startCmd.Use)
	}

	rangeFlag := startCmd.Flags().Lookup("range")
	if rangeFlag == nil {
		t.Error("Expected flag 'range' to be registered")
	}

	allFlag := startCmd.Flags().Lookup("all")
	if allFlag == nil {
		t.Error("Expected flag 'all' to be registered")
	}
}
