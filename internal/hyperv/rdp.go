package hyperv

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// RDPCredentials contains parsed RDP credentials
type RDPCredentials struct {
	Username string
	Password string
}

// ParseCredentials parses "user@pass" or "user" format into RDPCredentials
// Format: "username" or "username@password"
func ParseCredentials(input string) RDPCredentials {
	if input == "" {
		return RDPCredentials{}
	}

	// Find the last @ to handle usernames like "domain\user@password" or "user@domain@password"
	lastAtIndex := strings.LastIndex(input, "@")
	if lastAtIndex == -1 {
		// No @ found, treat as username only
		return RDPCredentials{Username: input}
	}

	// Split at last @
	username := input[:lastAtIndex]
	password := input[lastAtIndex+1:]

	// If password is empty after @, treat entire input as username
	if password == "" {
		return RDPCredentials{Username: input}
	}

	return RDPCredentials{
		Username: username,
		Password: password,
	}
}

// GetVMIPAddress gets the IPv4 address of a VM by index
func (m *Manager) GetVMIPAddress(ctx context.Context, vmIndex int) (string, error) {
	vmName, err := m.GetVMNameByIndex(ctx, vmIndex)
	if err != nil {
		return "", err
	}

	return m.GetVMIPAddressByName(ctx, vmName)
}

// GetVMIPAddressByName gets the IPv4 address of a VM by name
func (m *Manager) GetVMIPAddressByName(ctx context.Context, vmName string) (string, error) {
	// First check if VM is running
	state, err := m.GetVMStatus(ctx, vmName)
	if err != nil {
		return "", fmt.Errorf("failed to get VM state: %v", err)
	}

	if state != "Running" {
		return "", fmt.Errorf("VM '%s' is not running (state: %s)", vmName, state)
	}

	// Get IPv4 address from VM network adapter
	psScript := fmt.Sprintf(`
		$ips = (Get-VMNetworkAdapter -VMName "%s").IPAddresses
		$ipv4 = $ips | Where-Object { $_ -match '^\d+\.\d+\.\d+\.\d+$' } | Select-Object -First 1
		if ($ipv4) { $ipv4 } else { "" }
	`, vmName)

	output, err := m.Exec.RunScript(ctx, psScript)
	if err != nil {
		return "", fmt.Errorf("failed to get VM IP address: %v\nOutput: %s", err, string(output))
	}

	ip := strings.TrimSpace(string(output))
	if ip == "" {
		return "", fmt.Errorf("VM '%s' has no IPv4 address assigned. Ensure:\n"+
			"  - VM has integration services installed\n"+
			"  - VM has a network adapter connected\n"+
			"  - VM has obtained an IP address", vmName)
	}

	return ip, nil
}

// ConnectRDP opens an RDP connection to a VM by index
func (m *Manager) ConnectRDP(ctx context.Context, vmIndex int, credentials string) error {
	ip, err := m.GetVMIPAddress(ctx, vmIndex)
	if err != nil {
		return err
	}

	return m.ConnectRDPByIP(ctx, ip, credentials)
}

// ConnectRDPByName opens an RDP connection to a VM by name
func (m *Manager) ConnectRDPByName(ctx context.Context, vmName, credentials string) error {
	ip, err := m.GetVMIPAddressByName(ctx, vmName)
	if err != nil {
		return err
	}

	return m.ConnectRDPByIP(ctx, ip, credentials)
}

// ConnectRDPByIP opens an RDP connection to a specific IP address
// credentials can be "username" or "username@password"
func (m *Manager) ConnectRDPByIP(ctx context.Context, ip, credentials string) error {
	creds := ParseCredentials(credentials)

	// If password is provided, save to Windows Credential Manager first
	if creds.Password != "" {
		if err := m.SaveRDPCredentials(ctx, ip, creds.Username, creds.Password); err != nil {
			return fmt.Errorf("failed to save RDP credentials: %v", err)
		}
	}

	// Build mstsc command
	// mstsc is an external interactive GUI program.
	// We use Command (not CommandContext) because we usually want it to detach or run independently if possible,
	// but CommandContext is safer if we want to kill it when CLI exits.
	// However, RDP usually stays open.
	// For now, let's use standard exec.Command but we can't use ctx easily without it being killed on context cancel.
	// If context is short-lived (CLI command), we might want it to survive?
	// But idiomatic Go says respect context. Use CommandContext.
	// If user hits Ctrl-C, we kill RDP window? Probably acceptable.

	cmd := exec.CommandContext(ctx, "mstsc", "/v:"+ip)

	// Start mstsc without waiting for it to finish
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start RDP client: %v", err)
	}

	return nil
}

// SaveRDPCredentials saves RDP credentials to Windows Credential Manager
func (m *Manager) SaveRDPCredentials(ctx context.Context, target, username, password string) error {
	// cmdkey
	cmd := exec.CommandContext(ctx, "cmdkey", fmt.Sprintf("/generic:TERMSRV/%s", target), fmt.Sprintf("/user:%s", username), fmt.Sprintf("/pass:%s", password))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to save credentials: %v\nOutput: %s", err, string(output))
	}

	return nil
}

// DeleteRDPCredentials removes RDP credentials from Windows Credential Manager
func (m *Manager) DeleteRDPCredentials(ctx context.Context, target string) error {
	cmd := exec.CommandContext(ctx, "cmdkey", fmt.Sprintf("/delete:TERMSRV/%s", target))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete credentials: %v\nOutput: %s", err, string(output))
	}

	return nil
}
