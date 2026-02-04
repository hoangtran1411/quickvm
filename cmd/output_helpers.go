package cmd

import (
	"context"
	"fmt"

	"quickvm/internal/hyperv"
	"quickvm/internal/output"
)

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
		return
	}

	// Use shared getIndices logic
	indices, err := getIndices(args, rangeStr, all, len(vms))
	if err != nil {
		output.PrintError("INVALID_ARGS", "Invalid arguments", err.Error())
		if !output.IsJSON() {
			fmt.Printf("❌ Error: %v\n", err)
		}
		return
	}

	if !output.IsJSON() && len(indices) > 1 {
		fmt.Printf("%s %s %d VMs...\n\n", config.ActionEmoji, config.ActionVerb, len(indices))
	}

	results := make([]VMOperationResult, 0, len(indices))
	successCount := 0
	failCount := 0

	for _, index := range indices {
		vm := vms[index-1]
		result := VMOperationResult{
			Index: index,
			Name:  vm.Name,
		}

		if !output.IsJSON() {
			fmt.Printf("%s %s VM: %s (Index: %d)...\n", config.ActionEmoji, config.ActionVerb, vm.Name, index)
		}

		if err := config.OperationFunc(ctx, manager, vm); err != nil {
			result.Success = false
			result.Error = err.Error()
			failCount++
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to %s VM '%s': %v\n", config.Operation, vm.Name, err)
			}
		} else {
			result.Success = true
			result.Message = fmt.Sprintf("VM %s successfully", config.SuccessVerb)
			successCount++
			if !output.IsJSON() {
				fmt.Printf("✅ VM '%s' %s successfully!\n", vm.Name, config.SuccessVerb)
			}
		}
		results = append(results, result)
	}

	// JSON output for AI agents
	if output.IsJSON() {
		output.PrintData(VMBatchResult{
			Operation:    config.Operation,
			Results:      results,
			SuccessCount: successCount,
			FailCount:    failCount,
			TotalCount:   len(indices),
		})
		return
	}

	if len(indices) > 1 {
		fmt.Printf("\n📊 Summary: %d %s, %d failed\n", successCount, config.SuccessVerb, failCount)
	}
}
