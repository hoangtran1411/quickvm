# QuickVM - Fast Hyper-V Virtual Machine Manager 🚀

<div align="center">

![QuickVM Logo](https://img.shields.io/badge/QuickVM-Hyper--V%20Manager-blue?style=for-the-badge&logo=windows)
![Release](https://img.shields.io/badge/Release-v1.4.0-blueviolet?style=for-the-badge)
![Go Version](https://img.shields.io/badge/Go-1.27.0-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
[![Build Status](https://github.com/hoangtran1411/quickvm/actions/workflows/build.yml/badge.svg)](https://github.com/hoangtran1411/quickvm/actions/workflows/build.yml)

**A beautiful TUI-based command-line tool for managing Hyper-V virtual machines**

[Features](#features) • [Installation](#installation) • [Usage](#usage) • [Screenshots](#screenshots)

</div>

---

## ✨ Features

- 🎨 **Beautiful TUI Interface** - Interactive table view with color-coded VM states and responsive actions
- ⚡ **Quick Commands** - Start/stop/restart VMs by index number
- 🚀 **Concurrent Batch Operations** - Start, stop, or restart multiple VMs simultaneously (`1 3 5`, `--range 1-5`, or `--all`) with bounded parallel worker pools
- 🛡️ **Reliable & Idempotent** - Graceful handling of already running VMs, deterministic result ordering, and strict 60s context timeout safety
- 🤖 **AI-Agent & Scripting Ready** - Machine-readable JSON output mode (`--output json` / `-o json`) across all commands
- 📊 **Real-time Monitoring** - Live VM status, CPU usage, memory, and uptime
- 🎯 **Easy Navigation** - Keyboard shortcuts for efficient VM management
- 🔄 **Auto-refresh** - Keep your VM list up-to-date with a single keypress
- 💻 **Windows Native** - Direct integration with Hyper-V via safe PowerShell execution

## 📋 Prerequisites

- Windows 10/11 with Hyper-V enabled
- Administrator privileges (required for Hyper-V management)
- Go 1.27.0 or higher (for building from source)

## 🚀 Installation

### Quick Install (Recommended)

For the easiest installation experience:

1. **Download** the latest release package for your architecture
   - [Windows AMD64 (64-bit Intel/AMD)](https://github.com/hoangtran1411/quickvm/releases)
   - [Windows ARM64](https://github.com/hoangtran1411/quickvm/releases)

2. **Extract** the ZIP file to a folder

3. **Run the installation menu**
   - Double-click `install-menu.bat`, or
   - Right-click on `install-menu.ps1` → Run with PowerShell

4. **Choose your installation location**:
   - **Option 1 (System)**: Install to `C:\Windows\System32` - available globally for all users (requires Admin)
   - **Option 2 (User)**: Install to `%USERPROFILE%\bin` - available for current user (recommended)
   - **Option 3 (Current)**: Keep in current directory - portable mode

### Automated Install

For scripted or custom installations:

```powershell
# Install for current user (recommended)
.\install.ps1 -InstallLocation User

# Install system-wide (requires Admin)
.\install.ps1 -InstallLocation System

# Keep in current directory
.\install.ps1 -InstallLocation Current

# Additional options
.\install.ps1 -InstallLocation User -CreateAlias  # Add 'qvm' alias
```

### Build from Source

For developers or those who want the latest code:

```bash
# Clone the repository
git clone https://github.com/hoangtran1411/quickvm.git
cd quickvm

# Download dependencies
go mod download

# Build the application
go build -ldflags="-s -w" -o quickvm.exe

# Lint the code (requires golangci-lint v2.8.0+)
golangci-lint run

# Install using the menu
.\install-menu.bat
```

## 📖 Usage

### Interactive TUI Mode

Launch the interactive interface by running:

```bash
quickvm
```

**Keyboard Shortcuts:**
- `↑/↓` - Navigate through VMs
- `Enter` - Start the selected VM
- `s` - Stop the selected VM
- `t` - Restart the selected VM
- `r` - Refresh VM list
- `q` or `Esc` - Quit

### Command Line Mode

#### List all VMs
```bash
quickvm list
# or
quickvm ls
```

#### Start VMs
```bash
# Start a single VM by index
quickvm start 1

# Start multiple VMs concurrently (parallel execution)
quickvm start 1 3 5

# Start a range of VMs
quickvm start --range 1-5

# Start all VMs
quickvm start --all
```

#### Stop VMs
```bash
# Stop a single VM
quickvm stop 1

# Stop multiple or all VMs
quickvm stop 1 2 3
quickvm stop --all
```

#### Restart VMs
```bash
# Restart a single VM or a range
quickvm restart 1
quickvm restart --range 1-3
```

#### AI Agent & Structured Output (JSON)
All commands support `--output json` (`-o json`) for machine-readable automation:
```bash
# List all VMs in JSON format
quickvm list -o json

# Batch start VMs with JSON status reporting
quickvm start 1 2 -o json

# Get system info as structured JSON
quickvm info -o json
```

#### View System Information
```bash
quickvm info
```

This will display:
- 🖥️ **CPU**: Name and number of cores
- 💾 **Memory**: Total, used, and free RAM (in MB and GB)
- 💿 **Disk Drives**: Name, free space, and total capacity for each drive
- 🔧 **Hyper-V Status**: Whether Hyper-V is enabled or disabled

#### Update QuickVM
```bash
# Check for updates and install
quickvm update

# Check for updates without installing
quickvm update --check-only

# Auto-install without prompting
quickvm update -y

# Check for updates before running any command
quickvm --update list
```

#### Enable Hyper-V
```bash
# Enable Hyper-V (will prompt for restart if needed)
quickvm enable

# Enable and restart immediately
quickvm enable -y

# Enable without restarting (manual restart required)
quickvm enable --no-restart
```

> ⚠️ **Note**: The `enable` command requires Administrator privileges.

#### Snapshot Management
```bash
# List all snapshots for a VM
quickvm snapshot list 1

# Create a new snapshot
quickvm snapshot create 1 "Before Update"

# Restore a VM to a snapshot
quickvm snapshot restore 1 "Before Update"

# Delete a snapshot
quickvm snapshot delete 1 "Old Snapshot"
```

#### Export/Import VMs
```bash
# Export a VM to a directory
quickvm export 1 "D:\Backups\VMs"

# Import a VM from an exported directory
quickvm import "D:\Backups\VMs\MyVM"

# Import with options
quickvm import "D:\Backups\VMs\MyVM" --copy        # Copy VM files
quickvm import "D:\Backups\VMs\MyVM" --new-id      # Generate new VM ID
```

#### GPU Passthrough (GPU-P)
```bash
# Check GPU partitioning support
quickvm gpu status

# Add GPU partition to a VM (requires Admin)
quickvm gpu add 1

# Remove GPU partition from a VM
quickvm gpu remove 1

# Show driver paths for manual copy to guest
quickvm gpu drivers
```

#### Remote Desktop (RDP)
```bash
# Connect to a running VM via RDP
quickvm rdp 1

# Connect with auto-login credentials
quickvm rdp 1 -u "admin@password123"
```

#### Workspace Management (VM Groups)
```bash
# Create a workspace with specific VMs
quickvm ws create "DevEnvironment" --vms "Proxy,WebApp,DB"

# List all workspaces
quickvm ws list

# Start all VMs in a workspace
quickvm ws start "DevEnvironment"

# Stop all VMs in a workspace
quickvm ws stop "DevEnvironment"
```

## 🎯 Quick Examples

```bash
# View all VMs in a formatted table
quickvm list

# Start the first VM in the list
quickvm start 1

# Stop the second VM
quickvm stop 2

# Restart the third VM
quickvm restart 3

# Launch interactive mode for visual management
quickvm
```

## 🏗️ Architecture

QuickVM is built with clean architecture principles:

```
quickvm/
├── cmd/            # CLI commands (Cobra)
│   ├── root.go            # Root command & TUI launcher
│   ├── list.go            # List VMs command
│   ├── start.go           # Start VM command
│   ├── stop.go            # Stop VM command
│   ├── restart.go         # Restart VM command
│   ├── output_helpers.go  # Concurrent batch operations & worker pool
│   ├── info.go            # System info command
│   ├── snapshot.go        # Snapshot management
│   ├── clone.go           # Clone VM command
│   ├── export.go          # Export VM command
│   ├── import.go          # Import VM command
│   ├── gpu.go             # GPU passthrough management
│   ├── rdp.go             # Remote Desktop connection
│   ├── workspace.go       # VM group management
│   ├── enable.go          # Enable Hyper-V command
│   ├── update.go          # Update command
│   └── version.go         # Version & build info
├── internal/       # Private application logic
│   ├── hyperv/            # Hyper-V integration layer
│   │   ├── hyperv.go      # Core VM management & idempotency
│   │   ├── snapshot.go    # Snapshot operations
│   │   ├── clone.go       # Clone operations
│   │   ├── export.go      # Export/Import operations
│   │   ├── gpu.go         # GPU passthrough logic
│   │   ├── rdp.go         # RDP & Credential logic
│   │   ├── sysinfo.go     # Hardware & System info
│   │   └── workspace.go   # Workspace profile logic
│   └── output/            # Formatter for console tables & AI-agent JSON
├── ui/             # TUI components (Bubble Tea)
│   └── table.go           # Interactive dashboard
├── updater/        # Auto-update functionality
├── main.go         # Application entry point
└── go.mod          # Go modules
```

## 🛠️ Technologies

- **[Cobra](https://github.com/spf13/cobra)** - CLI framework
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - TUI components
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **PowerShell** - Hyper-V integration

## 🎨 Design Principles

1. **User Experience First** - Intuitive keyboard navigation and clear visual feedback
2. **Performance** - Fast VM operations with minimal overhead
3. **Reliability** - Comprehensive error handling, automated testing, and proper process management via Context API
4. **Beauty** - Color-coded states and modern terminal aesthetics

## 📚 Documentation

For more detailed information, check out our comprehensive documentation:

### Getting Started
- **[Installation Guide](docs/INSTALLATION.md)** - Detailed installation instructions
- **[Quick Reference](docs/QUICK_REFERENCE.md)** - All commands and shortcuts at a glance

### User Guides
- **[Demo & Examples](docs/DEMO.md)** - Real-world use cases and power user tips

### Developer Documentation
- **[Developer Guide](docs/DEVELOPER.md)** - Architecture and development notes
- **[Feature Roadmap](docs/FEATURE_ROADMAP.md)** - Planned features and priorities
- **[Workflow Guide](docs/WORKFLOW.md)** - Development and deployment workflow
- **[Contributing Guide](docs/CONTRIBUTING.md)** - How to contribute to QuickVM
- **[Project Summary](docs/PROJECT_SUMMARY.md)** - Complete project overview
- **[AI Agent Setup](docs/AI_AGENT.md)** - Skills, workflows, and AI assistant configuration

## 🔒 Permissions

QuickVM requires administrator privileges because it manages Hyper-V virtual machines. Always run PowerShell or Command Prompt as Administrator when using QuickVM.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Charm.sh team for the amazing TUI libraries
- Cobra framework for CLI management
- The Go community for continued support

## 📧 Contact

For questions, suggestions, or issues, please open an issue on GitHub.

---

<div align="center">

**Made with ❤️ by a Go enthusiast <a href="https://github.com/hoangtran1411">Hoang Tran</a>**

⭐ Star this repo if you find it useful!

</div>

