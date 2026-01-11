package hyperv

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// SystemInfo contains system information
type SystemInfo struct {
	CPU    CPUInfo      `json:"cpu"`
	Memory MemoryInfo   `json:"memory"`
	Disks  []DiskInfo   `json:"disks"`
	HyperV HyperVStatus `json:"hyperV"`
}

// CPUInfo contains CPU information
type CPUInfo struct {
	Name  string `json:"name"`
	Cores int    `json:"cores"`
}

// MemoryInfo contains memory information
type MemoryInfo struct {
	TotalMB int64   `json:"totalMB"`
	TotalGB float64 `json:"totalGB"`
	FreeMB  int64   `json:"freeMB"`
	FreeGB  float64 `json:"freeGB"`
	UsedMB  int64   `json:"usedMB"`
	UsedGB  float64 `json:"usedGB"`
}

// DiskInfo contains disk information
type DiskInfo struct {
	Name    string  `json:"name"`
	FreeMB  int64   `json:"freeMB"`
	FreeGB  float64 `json:"freeGB"`
	TotalMB int64   `json:"totalMB"`
	TotalGB float64 `json:"totalGB"`
	UsedMB  int64   `json:"usedMB"`
	UsedGB  float64 `json:"usedGB"`
}

// HyperVStatus contains Hyper-V status information
type HyperVStatus struct {
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
}

// GetSystemInfo retrieves system information including CPU, RAM, and Hyper-V status.
// diskInfo is optional and only retrieved if includeDisk is true.
func (m *Manager) GetSystemInfo(includeDisk bool) (*SystemInfo, error) {
	info := &SystemInfo{}

	// Get CPU info
	cpuInfo, err := m.getCPUInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %v", err)
	}
	info.CPU = *cpuInfo

	// Get Memory info
	memInfo, err := m.getMemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %v", err)
	}
	info.Memory = *memInfo

	// Get Disk info only if requested
	if includeDisk {
		diskInfo, err := m.getDiskInfo()
		if err != nil {
			return nil, fmt.Errorf("failed to get disk info: %v", err)
		}
		info.Disks = diskInfo
	} else {
		info.Disks = []DiskInfo{}
	}

	// Get Hyper-V status
	hyperVStatus, err := m.getHyperVStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get Hyper-V status: %v", err)
	}
	info.HyperV = *hyperVStatus

	return info, nil
}

// getCPUInfo retrieves CPU information
func (m *Manager) getCPUInfo() (*CPUInfo, error) {
	psScript := `
		$cpu = Get-WmiObject -Class Win32_Processor
		@{
			Name = $cpu.Name
			Cores = $cpu.NumberOfCores
		} | ConvertTo-Json
	`

	output, err := m.Exec.RunCommand(psScript)
	if err != nil {
		return nil, fmt.Errorf("failed to execute PowerShell command: %v\nOutput: %s", err, string(output))
	}

	var result struct {
		Name  string `json:"Name"`
		Cores int    `json:"Cores"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse CPU info: %v", err)
	}

	return &CPUInfo{
		Name:  strings.TrimSpace(result.Name),
		Cores: result.Cores,
	}, nil
}

// getMemoryInfo retrieves memory information
func (m *Manager) getMemoryInfo() (*MemoryInfo, error) {
	psScript := `
		$os = Get-WmiObject -Class Win32_OperatingSystem
		$totalMB = [math]::Round($os.TotalVisibleMemorySize / 1024)
		$freeMB = [math]::Round($os.FreePhysicalMemory / 1024)
		$usedMB = $totalMB - $freeMB
		@{
			TotalMB = $totalMB
			TotalGB = [math]::Round($totalMB / 1024, 2)
			FreeMB = $freeMB
			FreeGB = [math]::Round($freeMB / 1024, 2)
			UsedMB = $usedMB
			UsedGB = [math]::Round($usedMB / 1024, 2)
		} | ConvertTo-Json
	`

	output, err := m.Exec.RunCommand(psScript)
	if err != nil {
		return nil, fmt.Errorf("failed to execute PowerShell command: %v\nOutput: %s", err, string(output))
	}

	var result MemoryInfo
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse memory info: %v", err)
	}

	return &result, nil
}

// getDiskInfo retrieves disk information
func (m *Manager) getDiskInfo() ([]DiskInfo, error) {
	psScript := `
		Get-WmiObject -Class Win32_LogicalDisk -Filter "DriveType=3" | ForEach-Object {
			$totalMB = [math]::Round($_.Size / 1MB)
			$freeMB = [math]::Round($_.FreeSpace / 1MB)
			$usedMB = $totalMB - $freeMB
			@{
				Name = $_.DeviceID
				TotalMB = $totalMB
				TotalGB = [math]::Round($totalMB / 1024, 2)
				FreeMB = $freeMB
				FreeGB = [math]::Round($freeMB / 1024, 2)
				UsedMB = $usedMB
				UsedGB = [math]::Round($usedMB / 1024, 2)
			}
		} | ConvertTo-Json
	`

	output, err := m.Exec.RunCommand(psScript)
	if err != nil {
		return nil, fmt.Errorf("failed to execute PowerShell command: %v\nOutput: %s", err, string(output))
	}

	outputStr := strings.TrimSpace(string(output))
	var disks []DiskInfo

	// Handle single disk case (PowerShell returns object, not array)
	if strings.HasPrefix(outputStr, "{") {
		var disk DiskInfo
		if err := json.Unmarshal(output, &disk); err != nil {
			return nil, fmt.Errorf("failed to parse disk info: %v", err)
		}
		disks = append(disks, disk)
	} else if strings.HasPrefix(outputStr, "[") {
		if err := json.Unmarshal(output, &disks); err != nil {
			return nil, fmt.Errorf("failed to parse disks info: %v", err)
		}
	}

	return disks, nil
}

// getHyperVStatus retrieves Hyper-V status
func (m *Manager) getHyperVStatus() (*HyperVStatus, error) {
	psScript := `
		$hyperv = Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V
		if ($hyperv -eq $null) {
			@{
				Enabled = $false
				Status = "Not Installed"
			} | ConvertTo-Json
		} else {
			@{
				Enabled = ($hyperv.State -eq "Enabled")
				Status = $hyperv.State.ToString()
			} | ConvertTo-Json
		}
	`

	output, err := m.Exec.RunCommand(psScript)
	if err != nil {
		// If error, try alternative method using Get-Service
		return m.getHyperVStatusAlternative()
	}

	var result HyperVStatus
	if err := json.Unmarshal(output, &result); err != nil {
		// If parsing fails, try alternative method
		return m.getHyperVStatusAlternative()
	}

	return &result, nil
}

// getHyperVStatusAlternative uses Get-Service as fallback
func (m *Manager) getHyperVStatusAlternative() (*HyperVStatus, error) {
	psScript := `
		$vmms = Get-Service -Name vmms -ErrorAction SilentlyContinue
		if ($vmms -eq $null) {
			@{
				Enabled = $false
				Status = "Not Installed"
			} | ConvertTo-Json
		} else {
			@{
				Enabled = ($vmms.Status -eq "Running")
				Status = $vmms.Status.ToString()
			} | ConvertTo-Json
		}
	`

	output, err := m.Exec.RunCommand(psScript)
	if err != nil {
		return &HyperVStatus{
			Enabled: false,
			Status:  "Unknown",
		}, nil
	}

	var result HyperVStatus
	if err := json.Unmarshal(output, &result); err != nil {
		return &HyperVStatus{
			Enabled: false,
			Status:  "Unknown",
		}, nil
	}

	return &result, nil
}

// EnableHyperV enables Hyper-V on the system
// Returns true if a restart is needed, false otherwise
func (m *Manager) EnableHyperV() (bool, error) {
	// Try to enable all Hyper-V features
	psScript := `
		$features = @(
			"Microsoft-Hyper-V-All",
			"Microsoft-Hyper-V",
			"Microsoft-Hyper-V-Tools-All",
			"Microsoft-Hyper-V-Management-PowerShell",
			"Microsoft-Hyper-V-Hypervisor",
			"Microsoft-Hyper-V-Services",
			"Microsoft-Hyper-V-Management-Clients"
		)
		
		$needsRestart = $false
		$enabled = $false
		
		foreach ($feature in $features) {
			try {
				$result = Enable-WindowsOptionalFeature -Online -FeatureName $feature -All -NoRestart -ErrorAction SilentlyContinue
				if ($result -ne $null) {
					$enabled = $true
					if ($result.RestartNeeded) {
						$needsRestart = $true
					}
				}
			} catch {
				# Continue to next feature
			}
		}
		
		if (-not $enabled) {
			# Try DISM as fallback
			$dismResult = dism /online /enable-feature /featurename:Microsoft-Hyper-V-All /all /norestart 2>&1
			if ($LASTEXITCODE -eq 0 -or $LASTEXITCODE -eq 3010) {
				$enabled = $true
				if ($LASTEXITCODE -eq 3010) {
					$needsRestart = $true
				}
			}
		}
		
		@{
			Enabled = $enabled
			NeedsRestart = $needsRestart
		} | ConvertTo-Json
	`

	output, err := m.Exec.RunCommand(psScript)
	if err != nil {
		return false, fmt.Errorf("failed to enable Hyper-V: %v\nOutput: %s", err, string(output))
	}

	var result struct {
		Enabled      bool `json:"Enabled"`
		NeedsRestart bool `json:"NeedsRestart"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		// If we can't parse the result, assume it worked but needs restart
		return true, nil
	}

	if !result.Enabled {
		return false, fmt.Errorf("failed to enable Hyper-V. Please enable it manually through Windows Features")
	}

	return result.NeedsRestart, nil
}

// IsRunningAsAdmin checks if the current process is running with administrator privileges
func IsRunningAsAdmin() bool {
	psScript := `
		$identity = [Security.Principal.WindowsIdentity]::GetCurrent()
		$principal = New-Object Security.Principal.WindowsPrincipal($identity)
		$principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
	`

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(output)) == "True"
}

// ScheduleRestart schedules a system restart after the specified number of seconds
func (m *Manager) ScheduleRestart(seconds int) error {
	psScript := fmt.Sprintf(`shutdown /r /t %d /c "Restarting to complete Hyper-V installation"`, seconds)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to schedule restart: %v\nOutput: %s", err, string(output))
	}

	return nil
}
