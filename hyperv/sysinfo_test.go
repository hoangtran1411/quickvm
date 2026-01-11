package hyperv

import (
	"testing"
)

// MockShellExecutor for testing SysInfo
type MockSysInfoExecutor struct {
	MockOutput string
	MockError  error
}

func (m *MockSysInfoExecutor) RunCommand(script string) ([]byte, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return []byte(m.MockOutput), nil
}

func TestGetSystemInfo(t *testing.T) {
	// Values are mocked in SmartMockVerifyDisk below

	manager := NewManager()
	// Inject a mock executor that returns valid JSON for everything to avoid unmarshal errors
	// In a real scenario we'd need a smarter mock that responds based on input script.
	// Let's create a SmartMockExecutor.

	smartMock := &SmartMockVerifyDisk{
		DiskCheckTriggered: false,
	}
	manager.Exec = smartMock

	// Test 1: includeDisk = false
	_, _ = manager.GetSystemInfo(false)
	if smartMock.DiskCheckTriggered {
		t.Error("Expected GetSystemInfo(false) NOT to scan disks, but it did.")
	}

	// Test 2: includeDisk = true
	// Reset trigger
	smartMock.DiskCheckTriggered = false
	_, _ = manager.GetSystemInfo(true)
	if !smartMock.DiskCheckTriggered {
		// Note: This test might fail if the previous calls (CPU/Mem) fail first and return early.
		// So our mock needs to return valid JSON for CPU and Mem.
		t.Error("Expected GetSystemInfo(true) TO scan disks, but it didn't.")
	}
}

type SmartMockVerifyDisk struct {
	DiskCheckTriggered bool
}

func (m *SmartMockVerifyDisk) RunCommand(script string) ([]byte, error) {
	// Detect what kind of script is running based on keywords
	if contains(script, "Win32_Processor") {
		return []byte(`{"Name": "Mock CPU", "Cores": 4}`), nil
	}
	if contains(script, "Win32_OperatingSystem") {
		return []byte(`{"TotalMB": 8192, "FreeMB": 4096}`), nil
	}
	if contains(script, "Win32_LogicalDisk") {
		m.DiskCheckTriggered = true
		return []byte(`[{"Name": "C:", "TotalMB": 100, "FreeMB": 50}]`), nil
	}
	if contains(script, "Get-WindowsOptionalFeature") || contains(script, "Get-Service") {
		return []byte(`{"Enabled": true, "Status": "Running"}`), nil
	}

	return []byte("{}"), nil
}

func contains(s, substr string) bool {
	// Simple helper
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
