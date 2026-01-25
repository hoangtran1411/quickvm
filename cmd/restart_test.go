package cmd

import (
	"context"
	"fmt"
	"quickvm/internal/hyperv"
	"testing"
)

func TestRunRestart(t *testing.T) {
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
			name: "Restart single VM",
			args: []string{"1"},
			setup: func(m *MockManager) {
				m.RestartVMFn = func(_ context.Context, index int) error {
					if index != 1 {
						return fmt.Errorf("wrong index")
					}
					return nil
				}
			},
		},
		{
			name: "Restart all VMs",
			all:  true,
			setup: func(m *MockManager) {
				count := 0
				m.RestartVMFn = func(_ context.Context, _ int) error {
					count++
					return nil
				}
			},
		},
		{
			name: "Failed to get VMs",
			args: []string{"1"},
			setup: func(m *MockManager) {
				m.GetVMsFn = func(_ context.Context) ([]hyperv.VM, error) {
					return nil, fmt.Errorf("hyper-v error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			m := &MockManager{
				GetVMsFn: func(_ context.Context) ([]hyperv.VM, error) {
					return mockVMs, nil
				},
			}
			if tt.setup != nil {
				tt.setup(m)
			}

			runRestart(context.Background(), m, tt.args, tt.rangeStr, tt.all)
		})
	}
}
