package cmd

import (
	"context"
	"fmt"
	"quickvm/internal/hyperv"
	"sync/atomic"
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
		wantErr  bool
		setup    func(*MockManager)
	}{
		{
			name:    "Restart single VM",
			args:    []string{"1"},
			wantErr: false,
			setup: func(m *MockManager) {
				m.RestartVMByNameFn = func(_ context.Context, name string) error {
					if name != "VM1" {
						return fmt.Errorf("wrong VM")
					}
					return nil
				}
			},
		},
		{
			name:    "Restart all VMs",
			all:     true,
			wantErr: false,
			setup: func(m *MockManager) {
				var count int32
				m.RestartVMByNameFn = func(_ context.Context, _ string) error {
					atomic.AddInt32(&count, 1)
					return nil
				}
			},
		},
		{
			name:    "Failed to get VMs",
			args:    []string{"1"},
			wantErr: true,
			setup: func(m *MockManager) {
				m.GetVMsFn = func(_ context.Context) ([]hyperv.VM, error) {
					return nil, fmt.Errorf("hyper-v error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origExit := osExit
			var exitCalled bool
			var exitCode int
			osExit = func(code int) {
				exitCalled = true
				exitCode = code
			}
			defer func() { osExit = origExit }()

			m := &MockManager{
				GetVMsFn: func(_ context.Context) ([]hyperv.VM, error) {
					return mockVMs, nil
				},
			}
			if tt.setup != nil {
				tt.setup(m)
			}

			runRestart(context.Background(), m, tt.args, tt.rangeStr, tt.all)

			if tt.wantErr {
				if !exitCalled {
					t.Errorf("Expected osExit(1) to be called, but it was not")
				} else if exitCode != 1 {
					t.Errorf("Expected exit code 1, got %d", exitCode)
				}
			} else if exitCalled {
				t.Errorf("Unexpected osExit(%d) called", exitCode)
			}
		})
	}
}
