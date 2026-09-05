package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"quickvm/internal/hyperv"
	"quickvm/internal/output"
)

var osExit = os.Exit

const defaultBatchConcurrency = 4

// VMOperationResult represents the result of a single VM operation
type VMOperationResult struct {
	Index   int    `json:"index"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// VMBatchResult represents the result of a batch VM operation
type VMBatchResult struct {
	Operation    string              `json:"operation"`
	Results      []VMOperationResult `json:"results"`
	SuccessCount int                 `json:"successCount"`
	FailCount    int                 `json:"failCount"`
	TotalCount   int                 `json:"totalCount"`
}

// SnapshotListResult represents the result of listing snapshots
type SnapshotListResult struct {
	VMName    string            `json:"vmName"`
	VMIndex   int               `json:"vmIndex"`
	Snapshots []hyperv.Snapshot `json:"snapshots"`
	Total     int               `json:"total"`
}

// SnapshotOpResult represents the result of a snapshot operation
type SnapshotOpResult struct {
	Operation    string `json:"operation"`
	VMName       string `json:"vmName"`
	VMIndex      int    `json:"vmIndex"`
	SnapshotName string `json:"snapshotName"`
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
	Error        string `json:"error,omitempty"`
}

// ExportResult represents the result of an export operation
type ExportResult struct {
	VMName     string `json:"vmName"`
	VMIndex    int    `json:"vmIndex"`
	ExportPath string `json:"exportPath"`
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
}

// CloneResult represents the result of a clone operation
type CloneResult struct {
	SourceName  string `json:"sourceName"`
	SourceIndex int    `json:"sourceIndex"`
	NewName     string `json:"newName"`
	Success     bool   `json:"success"`
	Message     string `json:"message,omitempty"`
	Error       string `json:"error,omitempty"`
}

// RDPResult represents the result of an RDP connection attempt
type RDPResult struct {
	VMName    string `json:"vmName"`
	VMIndex   int    `json:"vmIndex"`
	IPAddress string `json:"ipAddress"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	VMName     string `json:"vmName"`
	ImportPath string `json:"importPath"`
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
}

// GPUStatusResult represents the result of checking GPU partitioning support
type GPUStatusResult struct {
	GPUs      []hyperv.GPUInfo `json:"gpus"`
	Supported bool             `json:"supported"`
	Total     int              `json:"total"`
}

// GPUOpResult represents the result of a GPU partition operation
type GPUOpResult struct {
	Operation string `json:"operation"`
	VMName    string `json:"vmName"`
	VMIndex   int    `json:"vmIndex"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
}

// GPUDriversResult represents the result of searching for GPU driver files
type GPUDriversResult struct {
	Paths []string `json:"paths"`
	Total int      `json:"total"`
}

// EnableResult represents the result of enabling Hyper-V
type EnableResult struct {
	AlreadyEnabled   bool   `json:"alreadyEnabled"`
	NeedsRestart     bool   `json:"needsRestart"`
	RestartScheduled bool   `json:"restartScheduled"`
	Success          bool   `json:"success"`
	Message          string `json:"message"`
}

// UpdateResult represents the result of checking or applying updates
type UpdateResult struct {
	HasUpdate      bool   `json:"hasUpdate"`
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	ReleaseNotes   string `json:"releaseNotes,omitempty"`
	Installed      bool   `json:"installed"`
	Success        bool   `json:"success"`
	Message        string `json:"message"`
}

// VMOperationFunc is a function that performs an operation on a VM
type VMOperationFunc func(ctx context.Context, manager hyperv.VMManager, vm hyperv.VM) error

// VMOperationConfig configures a batch VM operation
type VMOperationConfig struct {
	Operation     string          // e.g., "start", "stop", "restart"
	ActionVerb    string          // e.g., "Starting", "Stopping", "Restarting"
	ActionEmoji   string          // e.g., "🚀", "🛑", "🔄"
	SuccessVerb   string          // e.g., "started", "stopped", "restarted"
	OperationFunc VMOperationFunc // The actual operation to perform
}

// runVMBatchOperation executes a batch operation on VMs with JSON/table output support
//
//nolint:funlen // Dual output mode (JSON/table) requires handling both formats
func runVMBatchOperation(
	ctx context.Context,
	manager hyperv.VMManager,
	args []string,
	rangeStr string,
	all bool,
	config VMOperationConfig,
) {
	// Get VMs to validate index and get name
	vms, err := manager.GetVMs(ctx)
	if err != nil {
		output.PrintError("VM_GET_FAILED", "Failed to get VMs", err.Error())
		if !output.IsJSON() {
			fmt.Printf("❌ Failed to get VMs: %v\n", err)
		}
		osExit(1)
		return
	}

	// Use shared getIndices logic
	indices, err := getIndices(args, rangeStr, all, len(vms))
	if err != nil {
		output.PrintError("INVALID_ARGS", "Invalid arguments", err.Error())
		if !output.IsJSON() {
			fmt.Printf("❌ Error: %v\n", err)
		}
		osExit(1)
		return
	}

	if !output.IsJSON() && len(indices) > 1 {
		fmt.Printf("%s %s %d VMs...\n\n", config.ActionEmoji, config.ActionVerb, len(indices))
	}

	results := make([]VMOperationResult, len(indices))
	successCount := 0
	failCount := 0

	// why: Using a bounded worker pool (errgroup + semaphore) executes operations concurrently,
	// drastically reducing batch execution time (from N*10s to N/4*10s) while avoiding Hyper-V I/O storms.
	var (
		mu sync.Mutex
		g  errgroup.Group
	)
	sem := make(chan struct{}, defaultBatchConcurrency)

	for i, index := range indices {
		i, index := i, index
		vm := vms[index-1]

		sem <- struct{}{}
		g.Go(func() error {
			defer func() { <-sem }()

			// why: Enforce 60s per-VM timeout so a hung external cmdlet does not block the entire batch.
			opCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()

			res := VMOperationResult{
				Index: index,
				Name:  vm.Name,
			}

			if !output.IsJSON() {
				// why: Protect stdout to prevent interleaved lines when multiple workers log simultaneously.
				mu.Lock()
				fmt.Printf("%s %s VM: %s (Index: %d)...\n", config.ActionEmoji, config.ActionVerb, vm.Name, index)
				mu.Unlock()
			}

			opErr := config.OperationFunc(opCtx, manager, vm)

			// why: Synchronize updates to shared success/fail counters and results slice.
			mu.Lock()
			defer mu.Unlock()

			if opErr != nil {
				res.Success = false
				res.Error = opErr.Error()
				failCount++
				if !output.IsJSON() {
					fmt.Printf("❌ Failed to %s VM '%s': %v\n", config.Operation, vm.Name, opErr)
				}
			} else {
				res.Success = true
				res.Message = fmt.Sprintf("VM %s successfully", config.SuccessVerb)
				successCount++
				if !output.IsJSON() {
					fmt.Printf("✅ VM '%s' %s successfully!\n", vm.Name, config.SuccessVerb)
				}
			}
			// why: Direct index assignment preserves the original argument order regardless of completion order.
			results[i] = res
			return nil
		})
	}

	_ = g.Wait()

	outputBatchResults(config, results, successCount, failCount, len(indices))
}

func outputBatchResults(config VMOperationConfig, results []VMOperationResult, successCount, failCount, totalCount int) {
	if output.IsJSON() {
		topSuccess := failCount == 0
		var errInfo *output.ErrorInfo
		if failCount > 0 {
			errInfo = &output.ErrorInfo{
				Code:    "BATCH_OPERATION_FAILED",
				Message: fmt.Sprintf("%d of %d operations failed", failCount, totalCount),
			}
		}
		output.PrintResponse(output.Response{
			Success: topSuccess,
			Data: VMBatchResult{
				Operation:    config.Operation,
				Results:      results,
				SuccessCount: successCount,
				FailCount:    failCount,
				TotalCount:   totalCount,
			},
			Error: errInfo,
		})
		if failCount > 0 {
			osExit(1)
		}
		return
	}

	if totalCount > 1 {
		fmt.Printf("\n📊 Summary: %d %s, %d failed\n", successCount, config.SuccessVerb, failCount)
	}
	if failCount > 0 {
		osExit(1)
	}
}
