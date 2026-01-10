package hyperv

import (
	"strings"
	"testing"
)

// TestGetVMIPAddress_InvalidIndex tests GetVMIPAddress with invalid indices
func TestGetVMIPAddress_InvalidIndex(t *testing.T) {
	manager := NewManager()

	testCases := []struct {
		name     string
		vmIndex  int
		wantErr  bool
		errMatch string
	}{
		{
			name:     "negative index",
			vmIndex:  -1,
			wantErr:  true,
			errMatch: "invalid VM index",
		},
		{
			name:     "zero index",
			vmIndex:  0,
			wantErr:  true,
			errMatch: "invalid VM index",
		},
		{
			name:     "very large index",
			vmIndex:  9999,
			wantErr:  true,
			errMatch: "invalid VM index",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := manager.GetVMIPAddress(tc.vmIndex)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tc.errMatch)
					return
				}
				if !strings.Contains(err.Error(), tc.errMatch) {
					t.Errorf("expected error containing '%s', got '%s'", tc.errMatch, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestGetVMIPAddressByName_NonExistent tests GetVMIPAddressByName with non-existent VM
func TestGetVMIPAddressByName_NonExistent(t *testing.T) {
	manager := NewManager()

	_, err := manager.GetVMIPAddressByName("QuickVM_NonExistent_Test_12345")
	if err == nil {
		t.Error("expected error for non-existent VM, got nil")
	}
}

// TestConnectRDP_InvalidIndex tests ConnectRDP with invalid indices
func TestConnectRDP_InvalidIndex(t *testing.T) {
	manager := NewManager()

	testCases := []struct {
		name     string
		vmIndex  int
		wantErr  bool
		errMatch string
	}{
		{
			name:     "negative index",
			vmIndex:  -1,
			wantErr:  true,
			errMatch: "invalid VM index",
		},
		{
			name:     "zero index",
			vmIndex:  0,
			wantErr:  true,
			errMatch: "invalid VM index",
		},
		{
			name:     "very large index",
			vmIndex:  9999,
			wantErr:  true,
			errMatch: "invalid VM index",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := manager.ConnectRDP(tc.vmIndex, "")
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tc.errMatch)
					return
				}
				if !strings.Contains(err.Error(), tc.errMatch) {
					t.Errorf("expected error containing '%s', got '%s'", tc.errMatch, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestConnectRDPByName_NonExistent tests ConnectRDPByName with non-existent VM
func TestConnectRDPByName_NonExistent(t *testing.T) {
	manager := NewManager()

	err := manager.ConnectRDPByName("QuickVM_NonExistent_Test_12345", "")
	if err == nil {
		t.Error("expected error for non-existent VM, got nil")
	}
}
