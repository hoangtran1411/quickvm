// Package output provides machine-readable and human-readable output formatting.
// This package improves AX (Agent Experience) by allowing AI agents to parse
// structured JSON output instead of human-formatted text.
package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// Format represents the output format type
type Format string

const (
	// FormatTable is the default human-readable table format
	FormatTable Format = "table"
	// FormatJSON is machine-readable JSON format (AX-friendly)
	FormatJSON Format = "json"
	// FormatText is plain text format
	FormatText Format = "text"
)

// Response represents a standardized API response for JSON output
// This structure is designed for optimal AI agent parsing.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains structured error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// CurrentFormat holds the global output format setting
var CurrentFormat = FormatTable

// SetFormat sets the global output format
func SetFormat(format string) error {
	switch format {
	case "json":
		CurrentFormat = FormatJSON
	case "table":
		CurrentFormat = FormatTable
	case "text":
		CurrentFormat = FormatText
	default:
		return fmt.Errorf("invalid output format: %s (valid: json, table, text)", format)
	}
	return nil
}

// IsJSON returns true if current format is JSON
func IsJSON() bool {
	return CurrentFormat == FormatJSON
}

// PrintSuccess prints a success response
func PrintSuccess(message string, data interface{}) {
	if IsJSON() {
		printJSON(Response{
			Success: true,
			Message: message,
			Data:    data,
		})
		return
	}
	// Human-readable output is handled by the caller
}

// PrintError prints an error response
func PrintError(code, message string, details string) {
	if IsJSON() {
		printJSON(Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    code,
				Message: message,
				Details: details,
			},
		})
		return
	}
	// Human-readable: print to stderr
	fmt.Fprintf(os.Stderr, "❌ %s: %s\n", code, message)
	if details != "" {
		fmt.Fprintf(os.Stderr, "   Details: %s\n", details)
	}
}

// PrintData prints data in the appropriate format
func PrintData(data interface{}) {
	if IsJSON() {
		printJSON(Response{
			Success: true,
			Data:    data,
		})
		return
	}
	// For non-JSON, caller handles the formatting
}

// printJSON marshals and prints JSON to stdout
func printJSON(v interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(v)
}

// MustJSON marshals data to JSON and panics on error (for internal use)
func MustJSON(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	}
	return string(data)
}
