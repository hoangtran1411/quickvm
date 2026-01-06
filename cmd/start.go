package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"quickvm/hyperv"

	"github.com/spf13/cobra"
)

var startRange string

var startCmd = &cobra.Command{
	Use:   "start [vm-index]",
	Short: "Start a Hyper-V virtual machine",
	Long: `Start a Hyper-V virtual machine by its index.

Examples:
  quickvm start 1           # Start VM at index 1
  quickvm start --range 1-5 # Start VMs from index 1 to 5
  quickvm start --range 1,3,5 # Start VMs at index 1, 3, and 5`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := hyperv.NewManager()

		// Get VMs to validate index and get name
		vms, err := manager.GetVMs()
		if err != nil {
			fmt.Printf("❌ Failed to get VMs: %v\n", err)
			return
		}

		// Handle range flag
		if startRange != "" {
			indices, err := parseRange(startRange, len(vms))
			if err != nil {
				fmt.Printf("❌ Invalid range: %v\n", err)
				return
			}

			fmt.Printf("🚀 Starting %d VMs...\n\n", len(indices))
			successCount := 0
			failCount := 0

			for _, index := range indices {
				vm := vms[index-1]
				fmt.Printf("🚀 Starting VM: %s (Index: %d)...\n", vm.Name, index)

				if err := manager.StartVMByName(vm.Name); err != nil {
					fmt.Printf("❌ Failed to start VM '%s': %v\n", vm.Name, err)
					failCount++
				} else {
					fmt.Printf("✅ VM '%s' started successfully!\n", vm.Name)
					successCount++
				}
			}

			fmt.Printf("\n📊 Summary: %d started, %d failed\n", successCount, failCount)
			return
		}

		// Handle single VM index
		if len(args) == 0 {
			fmt.Println("❌ Please provide a VM index or use --range flag")
			fmt.Println("Example: quickvm start 1")
			fmt.Println("Example: quickvm start --range 1-5")
			return
		}

		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("❌ Invalid VM index: %s\n", args[0])
			return
		}

		if index < 1 || index > len(vms) {
			fmt.Printf("❌ Invalid VM index: %d (valid range: 1-%d)\n", index, len(vms))
			return
		}

		vm := vms[index-1]
		fmt.Printf("🚀 Starting VM: %s (Index: %d)...\n", vm.Name, index)

		if err := manager.StartVM(index); err != nil {
			fmt.Printf("❌ Failed to start VM: %v\n", err)
			return
		}

		fmt.Printf("✅ VM '%s' started successfully!\n", vm.Name)
	},
}

// parseRange parses a range string like "1-5" or "1,3,5" into a slice of indices
func parseRange(rangeStr string, maxIndex int) ([]int, error) {
	var indices []int
	seen := make(map[int]bool)

	// Check if it's a range format (e.g., "1-5")
	if strings.Contains(rangeStr, "-") && !strings.Contains(rangeStr, ",") {
		parts := strings.Split(rangeStr, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format, use 'start-end' (e.g., 1-5)")
		}

		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid start index: %s", parts[0])
		}

		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid end index: %s", parts[1])
		}

		if start > end {
			return nil, fmt.Errorf("start index must be less than or equal to end index")
		}

		if start < 1 {
			return nil, fmt.Errorf("start index must be at least 1")
		}

		if end > maxIndex {
			return nil, fmt.Errorf("end index %d exceeds maximum VM index %d", end, maxIndex)
		}

		for i := start; i <= end; i++ {
			indices = append(indices, i)
		}
	} else {
		// Handle comma-separated format (e.g., "1,3,5")
		parts := strings.Split(rangeStr, ",")
		for _, part := range parts {
			index, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				return nil, fmt.Errorf("invalid index: %s", part)
			}

			if index < 1 || index > maxIndex {
				return nil, fmt.Errorf("index %d out of range (valid: 1-%d)", index, maxIndex)
			}

			if !seen[index] {
				indices = append(indices, index)
				seen[index] = true
			}
		}
	}

	if len(indices) == 0 {
		return nil, fmt.Errorf("no valid indices found")
	}

	return indices, nil
}

func init() {
	startCmd.Flags().StringVarP(&startRange, "range", "r", "", "Range of VM indices to start (e.g., '1-5' or '1,3,5')")
	rootCmd.AddCommand(startCmd)
}
