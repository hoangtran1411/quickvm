//go:build !windows

package hyperv

import "context"

// IsRunningAsAdmin returns false on non-Windows platforms.
func IsRunningAsAdmin(_ context.Context) bool {
	return false
}
