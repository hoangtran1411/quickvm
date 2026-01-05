# QuickVM Installation Guide

## 📋 Table of Contents

- [Overview](#overview)
- [Method 1: Interactive Installation (Recommended)](#method-1-interactive-installation-recommended)
- [Method 2: Automated Installation](#method-2-automated-installation)
- [Method 3: Build from Source](#method-3-build-from-source)
- [Installation Location Comparison](#installation-location-comparison)
- [Uninstallation](#uninstallation)
- [Troubleshooting](#troubleshooting)

---

## 🎯 Overview

QuickVM offers **3 installation methods** to suit all needs:

| Method | Difficulty | Recommended For | Requires Admin |
|--------|-----------|-----------------|----------------|
| **Interactive Installation** | ⭐ | New Users | Optional |
| **Automated Installation** | ⭐⭐ | Power Users | Optional |
| **Build from Source** | ⭐⭐⭐ | Developers | Optional |

---

## 📦 Method 1: Interactive Installation (Recommended)

### Step 1: Download QuickVM

1. Visit [GitHub Releases](https://github.com/hoangtran1411/quickvm/releases)
2. Download the package for your system:
   - **Windows AMD64** (Intel/AMD 64-bit) - `quickvm-vX.X.X-windows-amd64.zip`
   - **Windows ARM64** (Surface X, ARM PC) - `quickvm-vX.X.X-windows-arm64.zip`

### Step 2: Extract the Archive

1. Right-click the downloaded ZIP file
2. Select **"Extract All..."**
3. Choose destination folder (e.g., `C:\QuickVM`)

### Step 3: Run Installation Menu

There are **2 ways** to launch the installation menu:

#### Method A: Using Batch File (Easiest)
```
1. Open the extracted folder
2. Double-click "install-menu.bat"
```

#### Method B: Using PowerShell Script
```
1. Open the extracted folder
2. Right-click "install-menu.ps1"
3. Select "Run with PowerShell"
```

### Step 4: Choose Installation Location

The menu will display:

```
╔══════════════════════════════════════════════════════╗
║   QuickVM - Fast Hyper-V Virtual Machine Manager    ║
║          Interactive Installation Menu              ║
╚══════════════════════════════════════════════════════╝

Choose QuickVM installation location:

  1. System   - Install to C:\Windows\System32 (requires Admin)
                Available system-wide for all users

  2. User     - Install to ~\bin (no Admin required)
                Available for current user only

  3. Current  - Keep in current directory
                Run with .\quickvm.exe

  0. Exit

Enter your choice (0-3):
```

**Type the corresponding number** (1, 2, or 3) and press Enter.

### Step 5: Complete

- After successful installation, the terminal will show a confirmation message
- **Restart your terminal** to use the `quickvm` command
- Verify installation: `quickvm version`

---

## ⚙️ Method 2: Automated Installation

### Using PowerShell Script

Suitable for:
- Installation automation
- Batch deployment across multiple machines
- CI/CD pipelines

### Basic Syntax

```powershell
.\install.ps1 -InstallLocation <Location>
```

### Parameters

| Parameter | Value | Description |
|-----------|-------|-------------|
| `-InstallLocation` | `System` | Install to System32 (requires Admin) |
| | `User` | Install to %USERPROFILE%\bin (recommended) |
| | `Current` | Keep in current directory |
| `-SkipBuild` | Switch | Skip build step (use existing binary) |
| `-CreateAlias` | Switch | Create 'qvm' alias for PowerShell |

### Examples

#### Install for current user
```powershell
.\install.ps1 -InstallLocation User
```

#### System-wide installation with alias
```powershell
# Run PowerShell as Administrator
.\install.ps1 -InstallLocation System -CreateAlias
```

#### Portable mode installation
```powershell
.\install.ps1 -InstallLocation Current
```

#### Install from existing binary
```powershell
.\install.ps1 -InstallLocation User -SkipBuild
```

---

## 🔨 Method 3: Build from Source

### Requirements

- **Git** - [Download Git](https://git-scm.com/download/win)
- **Go 1.21+** - [Download Go](https://golang.org/dl/)

### Steps

#### 1. Clone Repository

```bash
git clone https://github.com/hoangtran1411/quickvm.git
cd quickvm
```

#### 2. Download Dependencies

```bash
go mod download
go mod verify
```

#### 3. Build

```bash
# Basic build
go build -o quickvm.exe

# Optimized build (smaller binary)
go build -ldflags="-s -w" -o quickvm.exe

# Build with version info
go build -ldflags="-s -w -X 'main.Version=1.0.0'" -o quickvm.exe
```

#### 4. Install

```bash
# Use interactive menu
.\install-menu.bat

# Or use script directly
.\install.ps1 -InstallLocation User
```

---

## 🔍 Installation Location Comparison

### Option 1: System (`C:\Windows\System32`)

**✅ Advantages:**
- Available from anywhere on the system
- All users can use it
- No PATH configuration needed

**❌ Disadvantages:**
- **Requires Administrator privileges**
- Harder to update (needs Admin each time)
- May trigger Windows Defender warnings

**🎯 Best for:**
- Shared servers/workstations
- IT Administrators

### Option 2: User (`%USERPROFILE%\bin`)

**✅ Advantages:**
- **No Administrator privileges required**
- Easy to update
- Automatically adds to PATH
- More secure

**❌ Disadvantages:**
- Only current user can use it
- Requires terminal restart on first install

**🎯 Best for:** (⭐ **RECOMMENDED**)
- Personal computers
- Most users
- Development environments

### Option 3: Current (Current directory)

**✅ Advantages:**
- Portable - can be moved around
- No installation required
- Doesn't affect system

**❌ Disadvantages:**
- Must run from containing directory
- Need to type full path: `.\quickvm.exe`

**🎯 Best for:**
- USB drives / Portable tools
- Testing
- Temporary use

---

## 🗑️ Uninstallation

### If installed to System

```powershell
# Run PowerShell as Administrator
Remove-Item C:\Windows\System32\quickvm.exe
```

### If installed to User

```powershell
Remove-Item $env:USERPROFILE\bin\quickvm.exe

# (Optional) Remove bin folder if empty
Remove-Item $env:USERPROFILE\bin -Force
```

### If installed to Current

```powershell
# Just delete the file from the folder
Remove-Item quickvm.exe
```

### Remove Alias (if created)

```powershell
# Open PowerShell profile
notepad $PROFILE

# Delete lines containing "QuickVM Alias"
# Save and close
```

---

## 🐛 Troubleshooting

### Error: "Execution policy"

**Symptoms:**
```
install-menu.ps1 cannot be loaded because running scripts is disabled
```

**Solution:**
```powershell
# Temporarily allow script execution
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process

# Then run again
.\install-menu.ps1
```

### Error: "Go is not installed"

**Symptoms:**
```
❌ Error: Go is not installed or not in PATH
```

**Solution:**
1. Download Go from https://golang.org/dl/
2. Install and restart terminal
3. Verify: `go version`

### Error: "Access is denied" (when installing to System)

**Symptoms:**
```
❌ Error: Failed to copy to System32
Access is denied
```

**Solution:**
```powershell
# Run PowerShell as Administrator
# Right-click PowerShell → "Run as Administrator"
# Then run the script again
```

### Error: "quickvm is not recognized"

**Symptoms:**
```
'quickvm' is not recognized as an internal or external command
```

**Solution:**
1. **Restart terminal** (important!)
2. Check PATH:
   ```powershell
   $env:Path -split ';' | Select-String "bin"
   ```
3. If not found, re-run installation with User option

### Script blocked by Windows Security

**Symptoms:**
- File has "Unblock" flag

**Solution:**
```powershell
# Right-click file → Properties → Unblock
# Or use PowerShell:
Unblock-File .\install-menu.ps1
Unblock-File .\install-menu.bat
Unblock-File .\install.ps1
```

---

## 📞 Support

If you encounter issues that can't be resolved:

1. Check [Issues](https://github.com/hoangtran1411/quickvm/issues) on GitHub
2. Create a new Issue with:
   - Operating system (Windows version)
   - Installation method used
   - Full error message
   - Screenshots (if applicable)

---

## ✅ Post-Installation Checklist

- [ ] Restart terminal
- [ ] Run `quickvm version` to verify
- [ ] Run `quickvm list` to see VM list
- [ ] (Optional) Create `qvm` alias by re-running with `-CreateAlias` flag

---

**🎉 Congratulations! You've successfully installed QuickVM!**

See [Quick Reference](QUICK_REFERENCE.md) for usage instructions.
