package cmd

import (
	"quickvm/hyperv"
)

// MockManager is a mock implementation of VMManager for testing
type MockManager struct {
	GetVMsFn          func() ([]hyperv.VM, error)
	StartVMFn         func(index int) error
	StartVMByNameFn   func(name string) error
	StopVMFn          func(index int) error
	StopVMByNameFn    func(name string) error
	RestartVMFn       func(index int) error
	RestartVMByNameFn func(name string) error
	GetVMStatusFn     func(name string) (string, error)
}

func (m *MockManager) GetVMs() ([]hyperv.VM, error) {
	if m.GetVMsFn != nil {
		return m.GetVMsFn()
	}
	return []hyperv.VM{}, nil
}

func (m *MockManager) StartVM(index int) error {
	if m.StartVMFn != nil {
		return m.StartVMFn(index)
	}
	return nil
}

func (m *MockManager) StartVMByName(name string) error {
	if m.StartVMByNameFn != nil {
		return m.StartVMByNameFn(name)
	}
	return nil
}

func (m *MockManager) StopVM(index int) error {
	if m.StopVMFn != nil {
		return m.StopVMFn(index)
	}
	return nil
}

func (m *MockManager) StopVMByName(name string) error {
	if m.StopVMByNameFn != nil {
		return m.StopVMByNameFn(name)
	}
	return nil
}

func (m *MockManager) RestartVM(index int) error {
	if m.RestartVMFn != nil {
		return m.RestartVMFn(index)
	}
	return nil
}

func (m *MockManager) RestartVMByName(name string) error {
	if m.RestartVMByNameFn != nil {
		return m.RestartVMByNameFn(name)
	}
	return nil
}

func (m *MockManager) GetVMStatus(name string) (string, error) {
	if m.GetVMStatusFn != nil {
		return m.GetVMStatusFn(name)
	}
	return "Running", nil
}
