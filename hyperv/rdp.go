package hyperv

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetVMIPAddress gets the IPv4 address of a VM by index
func (m *Manager) GetVMIPAddress(vmIndex int) (string, error) {
	vmName, err := m.GetVMNameByIndex(vmIndex)
	if err != nil {
		return "", err
	}

	return m.GetVMIPAddressByName(vmName)
}

// GetVMIPAddressByName gets the IPv4 address of a VM by name
func (m *Manager) GetVMIPAddressByName(vmName string) (string, error) {
	// First check if VM is running
	state, err := m.GetVMStatus(vmName)
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

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
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
func (m *Manager) ConnectRDP(vmIndex int, username string) error {
	ip, err := m.GetVMIPAddress(vmIndex)
	if err != nil {
		return err
	}

	return m.ConnectRDPByIP(ip, username)
}

// ConnectRDPByName opens an RDP connection to a VM by name
func (m *Manager) ConnectRDPByName(vmName, username string) error {
	ip, err := m.GetVMIPAddressByName(vmName)
	if err != nil {
		return err
	}

	return m.ConnectRDPByIP(ip, username)
}

// ConnectRDPByIP opens an RDP connection to a specific IP address
func (m *Manager) ConnectRDPByIP(ip, username string) error {
	var cmd *exec.Cmd

	if username != "" {
		// With username: mstsc /v:IP /admin (username will be prompted)
		// Note: For auto-login, Windows Credential Manager would be needed
		cmd = exec.Command("mstsc", "/v:"+ip)
	} else {
		cmd = exec.Command("mstsc", "/v:"+ip)
	}

	// Start mstsc without waiting for it to finish
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start RDP client: %v", err)
	}

	return nil
}
