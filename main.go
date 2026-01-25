// Package main is the entry point for the QuickVM application.
package main

import (
	"quickvm/cmd"
)

// QuickVM - Fast Hyper-V Virtual Machine Manager
//
// A TUI-based command-line tool for managing Hyper-V virtual machines
// with an intuitive interface and quick commands.
//
// Features:
// - Interactive TUI with table view of all VMs
// - Quick start/stop/restart commands by index
// - Real-time VM status monitoring
// - Beautiful color-coded interface
//
// Usage:
//   quickvm              - Launch interactive TUI
//   quickvm list         - List all VMs
//   quickvm start <idx>  - Start VM by index
//   quickvm stop <idx>   - Stop VM by index
//   quickvm restart <idx>- Restart VM by index

func main() {
	cmd.Execute()
}
