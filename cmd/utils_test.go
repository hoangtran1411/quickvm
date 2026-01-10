package cmd

import (
	"testing"
)

func TestParseRange(t *testing.T) {
	maxIndex := 10

	tests := []struct {
		name     string
		rangeStr string
		want     []int
		wantErr  bool
	}{
		{"Single index", "1", []int{1}, false},
		{"Simple range", "1-3", []int{1, 2, 3}, false},
		{"Comma separated", "1,3,5", []int{1, 3, 5}, false},
		{"Mixed format", "1-2,4,6-7", []int{1, 2, 4, 6, 7}, false},
		{"With spaces", " 1 - 2 , 4 ", []int{1, 2, 4}, false},
		{"Duplicates", "1,1,1-2", []int{1, 2}, false},
		{"Out of bounds", "11", nil, true},
		{"Invalid segment", "1-2-3", nil, true},
		{"Start > End", "5-1", nil, true},
		{"Non-numeric", "abc", nil, true},
		{"Empty segment", "1,,3", []int{1, 3}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRange(tt.rangeStr, maxIndex)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !sliceEqual(got, tt.want) {
				t.Errorf("parseRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIndices(t *testing.T) {
	totalVMs := 5

	tests := []struct {
		name     string
		args     []string
		rangeStr string
		all      bool
		want     []int
		wantErr  bool
	}{
		{"All flag", nil, "", true, []int{1, 2, 3, 4, 5}, false},
		{"Positional args", []string{"1", "3"}, "", false, []int{1, 3}, false},
		{"Range flag", nil, "1-3", false, []int{1, 2, 3}, false},
		{"Combined", []string{"5"}, "1-2", false, []int{1, 2, 5}, false},
		{"Overlap", []string{"2"}, "1-3", false, []int{1, 2, 3}, false},
		{"No VMs specified", nil, "", false, nil, true},
		{"Invalid arg", []string{"abc"}, "", false, nil, true},
		{"Arg out of bounds", []string{"6"}, "", false, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIndices(tt.args, tt.rangeStr, tt.all, totalVMs)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIndices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !sliceEqual(got, tt.want) {
				t.Errorf("getIndices() = %v, want %v", got, tt.want)
			}
		})
	}
}

// sliceEqual checks if two int slices are equal
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
