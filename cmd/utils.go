package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// getIndices parses positional arguments, range strings, and the "all" flag
// to return a deduplicated, sorted list of 1-based VM indices.
func getIndices(args []string, rangeStr string, all bool, totalVMs int) ([]int, error) {
	indexMap := make(map[int]bool)

	// 1. Handle "all" flag
	if all {
		for i := 1; i <= totalVMs; i++ {
			indexMap[i] = true
		}
		return mapToSortedSlice(indexMap), nil
	}

	// 2. Handle range flag
	if rangeStr != "" {
		indices, err := parseRange(rangeStr, totalVMs)
		if err != nil {
			return nil, err
		}
		for _, idx := range indices {
			indexMap[idx] = true
		}
	}

	// 3. Handle positional arguments
	for _, arg := range args {
		idx, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid index '%s': must be a number", arg)
		}
		if idx < 1 || idx > totalVMs {
			return nil, fmt.Errorf("index %d out of range (1-%d)", idx, totalVMs)
		}
		indexMap[idx] = true
	}

	if len(indexMap) == 0 {
		return nil, fmt.Errorf("no VMs specified. Use index, --range, or --all")
	}

	return mapToSortedSlice(indexMap), nil
}

// mapToSortedSlice converts a map of indices to a sorted slice
func mapToSortedSlice(m map[int]bool) []int {
	slice := make([]int, 0, len(m))
	for k := range m {
		slice = append(slice, k)
	}
	sort.Ints(slice)
	return slice
}

// parseRange parses a range string like "1-5" or "1,3,5" into a slice of indices
func parseRange(rangeStr string, maxIndex int) ([]int, error) {
	var indices []int
	seen := make(map[int]bool)

	// Handle comma-separated format (e.g., "1,3,5" or "1-5,7,9-11")
	parts := strings.Split(rangeStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			// Handle range segment (e.g., "1-5")
			subParts := strings.Split(part, "-")
			if len(subParts) != 2 {
				return nil, fmt.Errorf("invalid range segment: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(subParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start index: %s", subParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(subParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end index: %s", subParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("start index (%d) must be <= end index (%d)", start, end)
			}

			if start < 1 || end > maxIndex {
				return nil, fmt.Errorf("range %d-%d out of bounds (1-%d)", start, end, maxIndex)
			}

			for i := start; i <= end; i++ {
				if !seen[i] {
					indices = append(indices, i)
					seen[i] = true
				}
			}
		} else {
			// Handle single index segment
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid index: %s", part)
			}

			if index < 1 || index > maxIndex {
				return nil, fmt.Errorf("index %d out of bounds (1-%d)", index, maxIndex)
			}

			if !seen[index] {
				indices = append(indices, index)
				seen[index] = true
			}
		}
	}

	if len(indices) == 0 {
		return nil, fmt.Errorf("no valid indices found in range '%s'", rangeStr)
	}

	return indices, nil
}
