package cmd

import (
	"context"
	"quickvm/internal/hyperv"
)

// MockManager is a mock implementation of VMManager for testing
type MockManager struct {
	GetVMsFn          func(ctx context.Context) ([]hyperv.VM, error)
	StartVMFn         func(ctx context.Context, index int) error
	StartVMByNameFn   func(ctx context.Context, name string) error
	StopVMFn          func(ctx context.Context, index int) error
	StopVMByNameFn    func(ctx context.Context, name string) error
	RestartVMFn       func(ctx context.Context, index int) error
	RestartVMByNameFn func(ctx context.Context, name string) error
	GetVMStatusFn     func(ctx context.Context, name string) (string, error)
}

func (m *MockManager) GetVMs(ctx context.Context) ([]hyperv.VM, error) {
	if m.GetVMsFn != nil {
		return m.GetVMsFn(ctx)
	}
	return []hyperv.VM{}, nil
}

func (m *MockManager) StartVM(ctx context.Context, index int) error {
	if m.StartVMFn != nil {
		return m.StartVMFn(ctx, index)
	}
	return nil
}

func (m *MockManager) StartVMByName(ctx context.Context, name string) error {
	if m.StartVMByNameFn != nil {
		return m.StartVMByNameFn(ctx, name)
	}
	return nil
}

func (m *MockManager) StopVM(ctx context.Context, index int) error {
	if m.StopVMFn != nil {
		return m.StopVMFn(ctx, index)
	}
	return nil
}

func (m *MockManager) StopVMByName(ctx context.Context, name string) error {
	if m.StopVMByNameFn != nil {
		return m.StopVMByNameFn(ctx, name)
	}
	return nil
}

func (m *MockManager) RestartVM(ctx context.Context, index int) error {
	if m.RestartVMFn != nil {
		return m.RestartVMFn(ctx, index)
	}
	return nil
}

func (m *MockManager) RestartVMByName(ctx context.Context, name string) error {
	if m.RestartVMByNameFn != nil {
		return m.RestartVMByNameFn(ctx, name)
	}
	return nil
}

func (m *MockManager) GetVMStatus(ctx context.Context, name string) (string, error) {
	if m.GetVMStatusFn != nil {
		return m.GetVMStatusFn(ctx, name)
	}
	return "Running", nil
}
