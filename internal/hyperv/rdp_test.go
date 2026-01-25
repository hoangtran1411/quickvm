package hyperv

import (
	"context"
	"os"
	"strings"
	"testing"
)

// skipIfNoHyperVEnv skips the test if running in CI/CD environment without Hyper-V or without Admin privileges
func skipIfNoHyperVEnv(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping test: Hyper-V not available in CI/CD environment")
	}
	if !IsRunningAsAdmin(context.TODO()) {
		t.Skip("Skipping test: Administrator privileges required for Hyper-V operations")
	}
}

// TestParseCredentials tests the ParseCredentials function
//
//nolint:funlen // Table-driven tests are naturally long
func TestParseCredentials(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		wantUsername string
		wantPassword string
	}{
		{
			name:         "empty string",
			input:        "",
			wantUsername: "",
			wantPassword: "",
		},
		{
			name:         "username only",
			input:        "admin",
			wantUsername: "admin",
			wantPassword: "",
		},
		{
			name:         "username with password",
			input:        "admin@password123",
			wantUsername: "admin",
			wantPassword: "password123",
		},
		{
			name:         "domain user with password",
			input:        "domain\\user@password",
			wantUsername: "domain\\user",
			wantPassword: "password",
		},
		{
			name:         "email-like username with password",
			input:        "user@domain.com@password",
			wantUsername: "user@domain.com",
			wantPassword: "password",
		},
		{
			name:         "password with @ symbol",
			input:        "user@pass@word",
			wantUsername: "user@pass",
			wantPassword: "word",
		},
		{
			name:         "username with @ at end (no password)",
			input:        "user@",
			wantUsername: "user@",
			wantPassword: "",
		},
		{
			name:         "only @ symbol",
			input:        "@",
			wantUsername: "@",
			wantPassword: "",
		},
		{
			name:         "password only (starts with @)",
			input:        "@password",
			wantUsername: "",
			wantPassword: "password",
		},
		{
			name:         "complex password with special chars",
			input:        "admin@P@ss!w0rd#123",
			wantUsername: "admin@P",
			wantPassword: "ss!w0rd#123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			creds := ParseCredentials(tc.input)
			if creds.Username != tc.wantUsername {
				t.Errorf("Username: got %q, want %q", creds.Username, tc.wantUsername)
			}
			if creds.Password != tc.wantPassword {
				t.Errorf("Password: got %q, want %q", creds.Password, tc.wantPassword)
			}
		})
	}
}

// TestRDPCredentials_Struct tests RDPCredentials struct
func TestRDPCredentials_Struct(t *testing.T) {
	creds := RDPCredentials{
		Username: "testuser",
		Password: "testpass",
	}

	if creds.Username != "testuser" {
		t.Errorf("Username: got %q, want %q", creds.Username, "testuser")
	}
	if creds.Password != "testpass" {
		t.Errorf("Password: got %q, want %q", creds.Password, "testpass")
	}
}

// TestGetVMIPAddress_InvalidIndex tests GetVMIPAddress with invalid indices
func TestGetVMIPAddress_InvalidIndex(t *testing.T) {
	skipIfNoHyperVEnv(t)
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
			_, err := manager.GetVMIPAddress(context.TODO(), tc.vmIndex)
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
	skipIfNoHyperVEnv(t)
	manager := NewManager()

	_, err := manager.GetVMIPAddressByName(context.TODO(), "QuickVM_NonExistent_Test_12345")
	if err == nil {
		t.Error("expected error for non-existent VM, got nil")
	}
}

// TestConnectRDP_InvalidIndex tests ConnectRDP with invalid indices
func TestConnectRDP_InvalidIndex(t *testing.T) {
	skipIfNoHyperVEnv(t)
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
			err := manager.ConnectRDP(context.TODO(), tc.vmIndex, "")
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
	skipIfNoHyperVEnv(t)
	manager := NewManager()

	err := manager.ConnectRDPByName(context.TODO(), "QuickVM_NonExistent_Test_12345", "")
	if err == nil {
		t.Error("expected error for non-existent VM, got nil")
	}
}
