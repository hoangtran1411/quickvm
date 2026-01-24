package hyperv

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// GPUInfo contains information about a partitionable GPU
type GPUInfo struct {
	Name                    string `json:"name"`
	PartitionCount          int    `json:"partitionCount"`
	ValidPartitionCounts    []int  `json:"validPartitionCounts"`
	MinPartitionVRAM        int64  `json:"minPartitionVRAM"`
	MaxPartitionVRAM        int64  `json:"maxPartitionVRAM"`
	OptimalPartitionVRAM    int64  `json:"optimalPartitionVRAM"`
	MinPartitionEncode      int64  `json:"minPartitionEncode"`
	MaxPartitionEncode      int64  `json:"maxPartitionEncode"`
	OptimalPartitionEncode  int64  `json:"optimalPartitionEncode"`
	MinPartitionDecode      int64  `json:"minPartitionDecode"`
	MaxPartitionDecode      int64  `json:"maxPartitionDecode"`
	OptimalPartitionDecode  int64  `json:"optimalPartitionDecode"`
	MinPartitionCompute     int64  `json:"minPartitionCompute"`
	MaxPartitionCompute     int64  `json:"maxPartitionCompute"`
	OptimalPartitionCompute int64  `json:"optimalPartitionCompute"`
}

// VMGPUPartition contains GPU partition info for a VM
type VMGPUPartition struct {
	VMName         string `json:"vmName"`
	HasGPU         bool   `json:"hasGpu"`
	PartitionCount int    `json:"partitionCount"`
}

// GPUPartitionConfig contains configuration for GPU partitioning
type GPUPartitionConfig struct {
	MinVRAM        int64
	MaxVRAM        int64
	OptimalVRAM    int64
	MinEncode      int64
	MaxEncode      int64
	OptimalEncode  int64
	MinDecode      int64
	MaxDecode      int64
	OptimalDecode  int64
	MinCompute     int64
	MaxCompute     int64
	OptimalCompute int64
	LowMMIOSpace   string // e.g., "1Gb"
	HighMMIOSpace  string // e.g., "32GB"
}

// DefaultGPUPartitionConfig returns the default GPU partition configuration
func DefaultGPUPartitionConfig() *GPUPartitionConfig {
	return &GPUPartitionConfig{
		MinVRAM:        80000000,
		MaxVRAM:        100000000,
		OptimalVRAM:    100000000,
		MinEncode:      80000000,
		MaxEncode:      100000000,
		OptimalEncode:  100000000,
		MinDecode:      80000000,
		MaxDecode:      100000000,
		OptimalDecode:  100000000,
		MinCompute:     80000000,
		MaxCompute:     100000000,
		OptimalCompute: 100000000,
		LowMMIOSpace:   "1Gb",
		HighMMIOSpace:  "32GB",
	}
}

// CheckGPUPartitionable checks if the system has GPU(s) that support partitioning
func (m *Manager) CheckGPUPartitionable(ctx context.Context) ([]GPUInfo, error) {
	psScript := `
		$gpus = Get-VMPartitionableGpu -ErrorAction SilentlyContinue
		if ($gpus -eq $null -or $gpus.Count -eq 0) {
			Write-Output "[]"
		} else {
			$result = @()
			foreach ($gpu in $gpus) {
				$result += @{
					Name = $gpu.Name
					PartitionCount = $gpu.PartitionCount
					ValidPartitionCounts = $gpu.ValidPartitionCounts
					MinPartitionVRAM = $gpu.MinPartitionVRAM
					MaxPartitionVRAM = $gpu.MaxPartitionVRAM
					OptimalPartitionVRAM = $gpu.OptimalPartitionVRAM
					MinPartitionEncode = $gpu.MinPartitionEncode
					MaxPartitionEncode = $gpu.MaxPartitionEncode
					OptimalPartitionEncode = $gpu.OptimalPartitionEncode
					MinPartitionDecode = $gpu.MinPartitionDecode
					MaxPartitionDecode = $gpu.MaxPartitionDecode
					OptimalPartitionDecode = $gpu.OptimalPartitionDecode
					MinPartitionCompute = $gpu.MinPartitionCompute
					MaxPartitionCompute = $gpu.MaxPartitionCompute
					OptimalPartitionCompute = $gpu.OptimalPartitionCompute
				}
			}
			$result | ConvertTo-Json -Depth 3
		}
	`

	output, err := m.Exec.RunScript(ctx, psScript)
	if err != nil {
		return nil, fmt.Errorf("failed to check GPU partitioning support: %v\nOutput: %s", err, string(output))
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" || outputStr == "[]" {
		return []GPUInfo{}, nil
	}

	var gpus []GPUInfo
	// Handle single GPU case
	if strings.HasPrefix(outputStr, "{") {
		var gpu GPUInfo
		if err := json.Unmarshal([]byte(outputStr), &gpu); err != nil {
			return nil, fmt.Errorf("failed to parse GPU info: %v", err)
		}
		gpus = append(gpus, gpu)
	} else if strings.HasPrefix(outputStr, "[") {
		if err := json.Unmarshal([]byte(outputStr), &gpus); err != nil {
			return nil, fmt.Errorf("failed to parse GPU info: %v", err)
		}
	}

	return gpus, nil
}

// GetVMGPUPartition gets GPU partition info for a specific VM
func (m *Manager) GetVMGPUPartition(ctx context.Context, vmName string) (*VMGPUPartition, error) {
	psScript := fmt.Sprintf(`
		$adapter = Get-VMGpuPartitionAdapter -VMName "%s" -ErrorAction SilentlyContinue
		if ($adapter -eq $null) {
			@{
				VMName = "%s"
				HasGPU = $false
				PartitionCount = 0
			} | ConvertTo-Json
		} else {
			@{
				VMName = "%s"
				HasGPU = $true
				PartitionCount = 1
			} | ConvertTo-Json
		}
	`, vmName, vmName, vmName)

	output, err := m.Exec.RunScript(ctx, psScript)
	if err != nil {
		return nil, fmt.Errorf("failed to get VM GPU partition info: %v\nOutput: %s", err, string(output))
	}

	var result VMGPUPartition
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse VM GPU partition info: %v", err)
	}

	return &result, nil
}

// AddGPUPartition adds a GPU partition to a VM
func (m *Manager) AddGPUPartition(ctx context.Context, vmName string, config *GPUPartitionConfig) error {
	if config == nil {
		config = DefaultGPUPartitionConfig()
	}

	// Check if VM is running
	state, err := m.GetVMStatus(ctx, vmName)
	if err != nil {
		return fmt.Errorf("failed to get VM status: %v", err)
	}
	if state == "Running" {
		return fmt.Errorf("VM '%s' must be stopped before adding GPU partition", vmName)
	}

	// Check if VM already has GPU
	gpuInfo, err := m.GetVMGPUPartition(ctx, vmName)
	if err != nil {
		return fmt.Errorf("failed to check existing GPU partition: %v", err)
	}
	if gpuInfo.HasGPU {
		return fmt.Errorf("VM '%s' already has a GPU partition", vmName)
	}

	// Add GPU Partition Adapter
	psScript := fmt.Sprintf(`
		# Step 1: Add GPU Partition Adapter
		Add-VMGpuPartitionAdapter -VMName "%s"
		
		# Step 2: Configure GPU Partition parameters
		Set-VMGpuPartitionAdapter -VMName "%s" `+
		`-MinPartitionVRAM %d -MaxPartitionVRAM %d -OptimalPartitionVRAM %d `+
		`-MinPartitionEncode %d -MaxPartitionEncode %d -OptimalPartitionEncode %d `+
		`-MinPartitionDecode %d -MaxPartitionDecode %d -OptimalPartitionDecode %d `+
		`-MinPartitionCompute %d -MaxPartitionCompute %d -OptimalPartitionCompute %d
		
		# Step 3: Enable Guest Controlled Cache Types
		Set-VM -GuestControlledCacheTypes $true -VMName "%s"
		
		# Step 4: Set Memory Mapped IO Space
		Set-VM -LowMemoryMappedIoSpace %s -VMName "%s"
		Set-VM -HighMemoryMappedIoSpace %s -VMName "%s"
		
		Write-Output "SUCCESS"
	`,
		vmName,
		vmName,
		config.MinVRAM, config.MaxVRAM, config.OptimalVRAM,
		config.MinEncode, config.MaxEncode, config.OptimalEncode,
		config.MinDecode, config.MaxDecode, config.OptimalDecode,
		config.MinCompute, config.MaxCompute, config.OptimalCompute,
		vmName,
		config.LowMMIOSpace, vmName,
		config.HighMMIOSpace, vmName,
	)

	output, err := m.Exec.RunScript(ctx, psScript)
	if err != nil {
		return fmt.Errorf("failed to add GPU partition: %v\nOutput: %s", err, string(output))
	}

	if !strings.Contains(string(output), "SUCCESS") {
		return fmt.Errorf("GPU partition may not have been added correctly. Output: %s", string(output))
	}

	return nil
}

// RemoveGPUPartition removes a GPU partition from a VM
func (m *Manager) RemoveGPUPartition(ctx context.Context, vmName string) error {
	// Check if VM is running
	state, err := m.GetVMStatus(ctx, vmName)
	if err != nil {
		return fmt.Errorf("failed to get VM status: %v", err)
	}
	if state == "Running" {
		return fmt.Errorf("VM '%s' must be stopped before removing GPU partition", vmName)
	}

	// Check if VM has GPU
	gpuInfo, err := m.GetVMGPUPartition(ctx, vmName)
	if err != nil {
		return fmt.Errorf("failed to check GPU partition: %v", err)
	}
	if !gpuInfo.HasGPU {
		return fmt.Errorf("VM '%s' does not have a GPU partition", vmName)
	}

	psScript := fmt.Sprintf(`
		Remove-VMGpuPartitionAdapter -VMName "%s"
		Write-Output "SUCCESS"
	`, vmName)

	output, err := m.Exec.RunScript(ctx, psScript)
	if err != nil {
		return fmt.Errorf("failed to remove GPU partition: %v\nOutput: %s", err, string(output))
	}

	if !strings.Contains(string(output), "SUCCESS") {
		return fmt.Errorf("GPU partition may not have been removed correctly. Output: %s", string(output))
	}

	return nil
}

// GetGPUDriverPaths returns the paths where GPU drivers are located on the host
func (m *Manager) GetGPUDriverPaths(ctx context.Context) ([]string, error) {
	psScript := `
		$paths = @()
		$driverPath = "C:\Windows\System32\DriverStore\FileRepository"
		
		# Find NVIDIA driver folders
		$nvidiaPaths = Get-ChildItem -Path $driverPath -Filter "nv_dispi.inf_amd64_*" -Directory -ErrorAction SilentlyContinue
		foreach ($p in $nvidiaPaths) {
			$paths += $p.FullName
		}
		
		# Find AMD driver folders
		$amdPaths = Get-ChildItem -Path $driverPath -Filter "u0*" -Directory -ErrorAction SilentlyContinue
		foreach ($p in $amdPaths) {
			$paths += $p.FullName
		}
		
		$paths | ConvertTo-Json
	`

	output, err := m.Exec.RunScript(ctx, psScript)
	if err != nil {
		return nil, fmt.Errorf("failed to get GPU driver paths: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" || outputStr == "null" {
		return []string{}, nil
	}

	var paths []string
	// Handle single path case
	if strings.HasPrefix(outputStr, "\"") {
		var path string
		if err := json.Unmarshal([]byte(outputStr), &path); err != nil {
			return nil, fmt.Errorf("failed to parse driver path: %v", err)
		}
		paths = append(paths, path)
	} else if strings.HasPrefix(outputStr, "[") {
		if err := json.Unmarshal([]byte(outputStr), &paths); err != nil {
			return nil, fmt.Errorf("failed to parse driver paths: %v", err)
		}
	}

	return paths, nil
}
