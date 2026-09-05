//go:build windows

package hyperv

import (
	"context"

	"golang.org/x/sys/windows"
)

// IsRunningAsAdmin checks if the current process is running with administrator privileges
// using native Windows token elevation (< 1µs, zero subprocess overhead).
func IsRunningAsAdmin(_ context.Context) bool {
	var token windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return false
	}
	defer func() {
		_ = token.Close()
	}()

	return token.IsElevated()
}
