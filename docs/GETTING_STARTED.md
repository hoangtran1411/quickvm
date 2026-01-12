# Getting Started with QuickVM

> Complete guide to install, understand, and use QuickVM effectively.

**Last Updated:** 2026-01-12

---

## 📋 Table of Contents

- [Quick Install](#-quick-install)
- [All Commands at a Glance](#-all-commands-at-a-glance)
- [TUI Keyboard Shortcuts](#-tui-keyboard-shortcuts)
- [Installation Methods](#-installation-methods)
- [Understanding the Codebase](#-understanding-the-codebase-for-developers)
- [Troubleshooting](#-troubleshooting)

---

## 🚀 Quick Install

### Fastest Way (2 minutes)

1. **Download** from [GitHub Releases](https://github.com/hoangtran1411/quickvm/releases)
   - Windows AMD64: `quickvm-vX.X.X-windows-amd64.zip`
   - Windows ARM64: `quickvm-vX.X.X-windows-arm64.zip`

2. **Extract** the ZIP file

3. **Run** `install-menu.bat` and choose option **2 (User)** *(recommended)*

4. **Restart terminal** and verify: `quickvm version`

---

## 📖 All Commands at a Glance

### Core Commands

| Command | Description | Example |
|---------|-------------|---------|
| `quickvm` | Launch interactive TUI | `quickvm` |
| `quickvm list` | List all VMs | `quickvm list` |
| `quickvm start <index>` | Start VM | `quickvm start 1` |
| `quickvm stop <index>` | Stop VM | `quickvm stop 1` |
| `quickvm restart <index>` | Restart VM | `quickvm restart 1` |
| `quickvm rdp <index>` | RDP connect to VM | `quickvm rdp 1` |
| `quickvm clone <index> <name>` | Clone VM | `quickvm clone 1 "MyVM-Copy"` |

### Snapshot Commands

| Command | Description |
|---------|-------------|
| `quickvm snapshot list <index>` | List snapshots |
| `quickvm snapshot create <index> "name"` | Create snapshot |
| `quickvm snapshot restore <index> "name"` | Restore snapshot |
| `quickvm snapshot delete <index> "name"` | Delete snapshot |

### Export/Import Commands

| Command | Description |
|---------|-------------|
| `quickvm export <index> <path>` | Export VM |
| `quickvm import <path>` | Import VM |

### System Commands

| Command | Description |
|---------|-------------|
| `quickvm info` | Show system info |
| `quickvm update` | Update QuickVM |
| `quickvm version` | Show version |

### Useful Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--range` | Start VMs in range | `quickvm start --range 1-5` |
| `--update` | Check updates first | `quickvm --update list` |
| `-y` | Auto-confirm | `quickvm update -y` |

---

## ⌨️ TUI Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate VMs |
| `Enter` | Start selected VM |
| `s` | Stop selected VM |
| `t` | Restart selected VM |
| `r` | Refresh VM list |
| `q` / `Esc` | Quit |

### Status Colors

| Color | State |
|-------|-------|
| 🟢 Green | Running |
| 🔴 Red | Off |
| 🟡 Yellow | Paused |

---

## 📦 Installation Methods

### Method 1: Interactive Menu (Recommended)

```
1. Extract ZIP
2. Double-click install-menu.bat
3. Choose:
   - 1. System (C:\Windows\System32) - requires Admin
   - 2. User (~\bin) - recommended, no Admin
   - 3. Current - portable mode
4. Restart terminal
```

### Method 2: PowerShell Script

```powershell
# User install (recommended)
.\install.ps1 -InstallLocation User

# System install (requires Admin)
.\install.ps1 -InstallLocation System

# With alias 'qvm'
.\install.ps1 -InstallLocation User -CreateAlias
```

### Method 3: Build from Source

```powershell
git clone https://github.com/hoangtran1411/quickvm.git
cd quickvm
go mod download
go build -ldflags="-s -w" -o quickvm.exe
.\install.ps1 -InstallLocation User
```

### Installation Location Comparison

| Location | Pros | Cons | Best For |
|----------|------|------|----------|
| **System** | Global access, no PATH config | Needs Admin | Shared workstations |
| **User** ⭐ | No Admin, easy update | Single user | Most users |
| **Current** | Portable, no install | Must use `.\quickvm.exe` | USB/Testing |

---

## 🧠 Understanding the Codebase (For Developers)

### High-Level Architecture

```
User → main.go → cmd/ (Cobra CLI) → hyperv/ (Business Logic) → PowerShell
```

### Project Structure

```
quickvm/
├── cmd/          # CLI commands (one file per command)
│   ├── root.go   # Root command + TUI launcher
│   ├── start.go  # quickvm start
│   ├── stop.go   # quickvm stop
│   └── ...
├── hyperv/       # Hyper-V integration via PowerShell
│   ├── hyperv.go # Core VM operations
│   ├── snapshot.go
│   └── ...
├── ui/           # TUI components (Bubble Tea)
└── main.go       # Entry point
```

### Key Concepts

**1. Entry Point (`main.go`)**
```go
func main() {
    cmd.Execute()
}
```
Just calls the CLI framework. All logic is in packages.

**2. Command Pattern (`cmd/start.go`)**
```go
var startCmd = &cobra.Command{
    Use:   "start <vm-index>",
    Short: "Start a virtual machine",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Parse args & call hyperv package
        return nil
    },
}

func init() {
    rootCmd.AddCommand(startCmd)
}
```

**3. Hyper-V Integration (`hyperv/hyperv.go`)**
```go
func (m *Manager) StartVM(index int) error {
    psScript := fmt.Sprintf(`Start-VM -Name "%s"`, vmName)
    output, err := m.Exec.RunCommand(psScript)
    // ...
}
```
Uses `ShellExecutor` interface for testability.

### Quick Tips for New Contributors

1. **Read `cmd/root.go`** - Understand CLI setup
2. **Check `hyperv/hyperv.go`** - See PowerShell integration
3. **Run tests**: `go test ./...`
4. **Format code**: `go fmt ./...`

---

## 🐛 Troubleshooting

### "Execution policy" error
```powershell
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process
.\install-menu.ps1
```

### "quickvm is not recognized"
1. Restart terminal
2. Check PATH: `$env:Path -split ';' | Select-String "bin"`
3. Re-install with User option

### "Access is denied" (System install)
Run PowerShell as Administrator

### "Failed to get VMs"
- Run as Administrator
- Check `Get-VM` works in PowerShell
- Verify Hyper-V is enabled

### Script blocked by Windows
```powershell
Unblock-File .\install-menu.ps1
Unblock-File .\install.ps1
```

---

## ⚡ Power User Tips

### Create alias
```powershell
Set-Alias qvm quickvm
```

### Start multiple VMs
```powershell
quickvm start --range 1-5
# or
1..3 | ForEach-Object { quickvm start $_ }
```

### Auto-start VMs at boot
```powershell
$action = New-ScheduledTaskAction -Execute "quickvm.exe" -Argument "start 1"
$trigger = New-ScheduledTaskTrigger -AtStartup
Register-ScheduledTask -TaskName "AutoStartVM" -Action $action -Trigger $trigger
```

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/hoangtran1411/quickvm/issues)
- **Releases**: [GitHub Releases](https://github.com/hoangtran1411/quickvm/releases)

---

**🎉 You're ready to use QuickVM!**

Next: Check out [Demo & Examples](DEMO.md) for real-world use cases.
