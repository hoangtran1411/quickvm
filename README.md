# QuickVM - Fast Hyper-V Virtual Machine Manager 🚀

<div align="center">

![QuickVM Logo](https://img.shields.io/badge/QuickVM-Hyper--V%20Manager-blue?style=for-the-badge&logo=windows)
![Go Version](https://img.shields.io/badge/Go-1.25.2-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
[![Build Status](https://github.com/hoangtran1411/quickvm/actions/workflows/build.yml/badge.svg)](https://github.com/hoangtran1411/quickvm/actions/workflows/build.yml)

**A beautiful TUI-based command-line tool for managing Hyper-V virtual machines**

[Features](#features) • [Installation](#installation) • [Usage](#usage) • [Screenshots](#screenshots)

</div>

---

## ✨ Features

- 🎨 **Beautiful TUI Interface** - Interactive table view with color-coded VM states
- ⚡ **Quick Commands** - Start/stop/restart VMs by index number
- 📊 **Real-time Monitoring** - Live VM status, CPU usage, memory, and uptime
- 🎯 **Easy Navigation** - Keyboard shortcuts for efficient VM management
- 🔄 **Auto-refresh** - Keep your VM list up-to-date with a single keypress
- 💻 **Windows Native** - Direct integration with Hyper-V via PowerShell

## 📋 Prerequisites

- Windows 10/11 with Hyper-V enabled
- Administrator privileges (required for Hyper-V management)
- Go 1.21 or higher (for building from source)

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

#### Start a VM
```bash
quickvm start 1
```

#### Stop a VM
```bash
quickvm stop 1
```

#### Restart a VM
```bash
quickvm restart 1
```

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
│   ├── root.go     # Root command & TUI launcher
│   ├── start.go    # Start VM command
│   ├── stop.go     # Stop VM command
│   ├── restart.go  # Restart VM command
│   └── list.go     # List VMs command
├── hyperv/         # Hyper-V integration layer
│   └── hyperv.go   # VM management via PowerShell
├── ui/             # TUI components
│   └── table.go    # Interactive table view (Bubble Tea)
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
3. **Reliability** - Comprehensive error handling and validation
4. **Beauty** - Color-coded states and modern terminal aesthetics

## 📚 Documentation

For more detailed information, check out our comprehensive documentation:

### Getting Started
- **[Installation Guide](docs/INSTALLATION.md)** - Detailed installation instructions (English)
- **[Hướng Dẫn Cài Đặt](docs/CAI_DAT.md)** - Chi tiết cài đặt (Tiếng Việt)
- **[Quick Reference](docs/QUICK_REFERENCE.md)** - All commands and shortcuts at a glance

### User Guides
- **[Vietnamese Guide](docs/HUONG_DAN.md)** - Hướng dẫn chi tiết bằng tiếng Việt
- **[Demo & Examples](docs/DEMO.md)** - Real-world use cases and power user tips

### Developer Documentation
- **[Developer Guide](docs/DEVELOPER.md)** - Architecture and development notes
- **[Workflow Guide](docs/WORKFLOW.md)** - Development and deployment workflow
- **[Contributing Guide](docs/CONTRIBUTING.md)** - How to contribute to QuickVM
- **[Project Summary](docs/PROJECT_SUMMARY.md)** - Complete project overview

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

**Made with ❤️ by a Go enthusiast with 10 years of experience**

⭐ Star this repo if you find it useful!

</div>
