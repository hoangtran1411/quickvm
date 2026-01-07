package cmd

import (
	"testing"
)

// ===== FAILURE TESTS FIRST =====
// Ưu tiên test các trường hợp lỗi trước vì hệ thống VM luôn có thể gặp sự cố

func TestParseRange_FailureCases(t *testing.T) {
	// Test với maxIndex = 10 để có baseline
	maxIndex := 10

	tests := []struct {
		name        string
		rangeStr    string
		wantErr     bool
		errContains string
	}{
		// === Invalid Format Errors ===
		{
			name:        "Empty string",
			rangeStr:    "",
			wantErr:     true,
			errContains: "invalid", // parseRange treats empty as invalid index
		},
		{
			name:        "Whitespace only",
			rangeStr:    "   ",
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name:        "Invalid range format - multiple dashes",
			rangeStr:    "1-2-3",
			wantErr:     true,
			errContains: "invalid range format",
		},
		{
			name:        "Non-numeric start",
			rangeStr:    "abc-5",
			wantErr:     true,
			errContains: "invalid start index",
		},
		{
			name:        "Non-numeric end",
			rangeStr:    "1-xyz",
			wantErr:     true,
			errContains: "invalid end index",
		},
		{
			name:        "Non-numeric comma-separated value",
			rangeStr:    "1,abc,3",
			wantErr:     true,
			errContains: "invalid index",
		},
		{
			name:        "Mixed invalid characters",
			rangeStr:    "1; 2; 3",
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name:        "Special characters",
			rangeStr:    "1@3",
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name:        "Float numbers",
			rangeStr:    "1.5-3.5",
			wantErr:     true,
			errContains: "invalid",
		},

		// === Boundary Violations ===
		{
			name:        "Zero index in range",
			rangeStr:    "0-5",
			wantErr:     true,
			errContains: "start index must be at least 1",
		},
		{
			name:        "Negative start index",
			rangeStr:    "-1-5",
			wantErr:     true,
			errContains: "invalid range format", // -1-5 has 3 parts when split by "-"
		},
		{
			name:        "End exceeds max index",
			rangeStr:    "1-100",
			wantErr:     true,
			errContains: "exceeds maximum",
		},
		{
			name:        "Start greater than end",
			rangeStr:    "5-1",
			wantErr:     true,
			errContains: "start index must be less than or equal",
		},
		{
			name:        "Single zero index in comma list",
			rangeStr:    "0",
			wantErr:     true,
			errContains: "out of range",
		},
		{
			name:        "Negative index in comma list",
			rangeStr:    "1,-2,3",
			wantErr:     true,
			errContains: "out of range",
		},
		{
			name:        "Index exceeds max in comma list",
			rangeStr:    "1,50,3",
			wantErr:     true,
			errContains: "out of range",
		},
		{
			name:        "All indices exceed max",
			rangeStr:    "100,200,300",
			wantErr:     true,
			errContains: "out of range",
		},

		// === Edge Cases for VM System Errors ===
		{
			name:        "Only commas, no numbers",
			rangeStr:    ",,,",
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name:        "Trailing comma",
			rangeStr:    "1,2,",
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name:        "Leading comma",
			rangeStr:    ",1,2",
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name:        "Only dash",
			rangeStr:    "-",
			wantErr:     true,
			errContains: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseRange(tt.rangeStr, maxIndex)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseRange(%q) expected error containing %q, but got nil", tt.rangeStr, tt.errContains)
				} else if tt.errContains != "" && !containsSubstring(err.Error(), tt.errContains) {
					t.Errorf("parseRange(%q) error = %q, want error containing %q", tt.rangeStr, err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("parseRange(%q) unexpected error: %v", tt.rangeStr, err)
				}
			}
		})
	}
}

func TestParseRange_ZeroMaxIndex(t *testing.T) {
	// Test khi không có VM nào (hệ thống lỗi hoặc Hyper-V không khả dụng)
	tests := []struct {
		name     string
		rangeStr string
	}{
		{"Single index when no VMs", "1"},
		{"Range when no VMs", "1-5"},
		{"Comma list when no VMs", "1,2,3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseRange(tt.rangeStr, 0)
			if err == nil {
				t.Error("Expected error when maxIndex is 0, got nil")
			}
		})
	}
}

func TestParseRange_NegativeMaxIndex(t *testing.T) {
	// Test edge case với maxIndex âm (không nên xảy ra nhưng phòng hờ)
	_, err := parseRange("1", -1)
	if err == nil {
		t.Error("Expected error when maxIndex is negative, got nil")
	}
}

// ===== SUCCESS TESTS (sau khi đã test failure) =====
func TestParseRange_SuccessCases(t *testing.T) {
	maxIndex := 10

	tests := []struct {
		name     string
		rangeStr string
		want     []int
	}{
		// Cơ bản
		{
			name:     "Single index",
			rangeStr: "1",
			want:     []int{1},
		},
		{
			name:     "Simple range",
			rangeStr: "1-3",
			want:     []int{1, 2, 3},
		},
		{
			name:     "Comma separated",
			rangeStr: "1,3,5",
			want:     []int{1, 3, 5},
		},

		// Với whitespace
		{
			name:     "Range with spaces",
			rangeStr: " 1 - 3 ",
			want:     []int{1, 2, 3},
		},
		{
			name:     "Comma with spaces",
			rangeStr: " 1 , 2 , 3 ",
			want:     []int{1, 2, 3},
		},

		// Boundary valid
		{
			name:     "Range at max boundary",
			rangeStr: "8-10",
			want:     []int{8, 9, 10},
		},
		{
			name:     "Full range",
			rangeStr: "1-10",
			want:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:     "Same start and end",
			rangeStr: "5-5",
			want:     []int{5},
		},

		// Duplicate handling
		{
			name:     "Duplicates in comma list should be removed",
			rangeStr: "1,1,2,2,3",
			want:     []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRange(tt.rangeStr, maxIndex)
			if err != nil {
				t.Errorf("parseRange(%q) unexpected error: %v", tt.rangeStr, err)
				return
			}
			if !sliceEqual(got, tt.want) {
				t.Errorf("parseRange(%q) = %v, want %v", tt.rangeStr, got, tt.want)
			}
		})
	}
}

// Helper functions
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func sliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
