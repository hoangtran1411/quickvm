# QuickVM - Quick Reference Card ⚡

> **Fast cheat sheet and command reference for QuickVM v1.5.0**  
> Manage Hyper-V virtual machines with speed, concurrency, and style via TUI, CLI, and AI-agent JSON automation.

---

## 📑 Table of Contents

- [Quick Start & Installation](#quick-start--installation)
- [Interactive TUI Navigation & Shortcuts](#interactive-tui-navigation--shortcuts)
- [CLI Command Matrix](#cli-command-matrix)
- [Batch Operations & Concurrency](#batch-operations--concurrency)
- [Remote Desktop (RDP) Integration](#remote-desktop-rdp-integration)
- [Snapshot & Checkpoint Management](#snapshot--checkpoint-management)
- [VM Workspaces (Group Management)](#vm-workspaces-group-management)
- [VM Cloning, Export & Import](#vm-cloning-export--import)
- [GPU Partitioning (GPU-P Passthrough)](#gpu-partitioning-gpu-p-passthrough)
- [AI Agent & Automation (JSON API)](#ai-agent--automation-json-api)
- [System & Maintenance](#system--maintenance)
- [Developer Commands & Makefile](#developer-commands--makefile)
- [Troubleshooting & Error Codes](#troubleshooting--error-codes)

---

## <a id="quick-start--installation"></a>🚀 Quick Start & Installation

### Requirements

- **OS**: Windows 10 / 11 (Pro, Enterprise, Education) or Windows Server
- **Hyper-V**: Enabled (`quickvm enable -y` or via Windows Features)
- **Privilege**: Administrator privileges required for Hyper-V management

### Installation Options

```powershell
# Option 1: Interactive installer menu (Recommended)
.\install-menu.bat

# Option 2: Scripted install for current user with 'qvm' alias
.\install.ps1 -InstallLocation User -CreateAlias

# Option 3: System-wide install (requires Admin)
.\install.ps1 -InstallLocation System

# Option 4: Manual PowerShell alias (add to $PROFILE)
Set-Alias qvm quickvm
```

### Launch Modes

```powershell
quickvm         # Launch interactive full-screen TUI dashboard
quickvm list    # CLI mode: print formatted VM table and exit
quickvm --help  # Show global help and available commands
```

---

## <a id="interactive-tui-navigation--shortcuts"></a>⌨️ Interactive TUI Navigation & Shortcuts

Run `quickvm` without arguments to launch the Bubble Tea terminal dashboard.

### Keybindings

| Key | Action | Description |
| --- | --- | --- |
| `↑` / `k` | **Cursor Up** | Select previous VM in list |
| `↓` / `j` | **Cursor Down** | Select next VM in list |
| `Enter` | **Start VM** | Start selected VM (60s timeout, non-blocking UI) |
| `s` | **Stop VM** | Stop selected VM safely |
| `t` | **Restart VM** | Restart selected VM |
| `r` | **Refresh** | Query Hyper-V and update table immediately |
| `q` / `Esc` / `Ctrl+C` | **Quit** | Cleanly exit TUI |

### Visual Indicators

| Indicator | VM State | Description |
| --- | --- | --- |
| 🟢 Green | `Running` | VM is actively running |
| 🔴 Red | `Off` | VM is shut down |
| 🟡 Yellow | `Paused` / Other | VM is paused, saved, or transitioning |

### Display Columns

`Index` • `Name` • `State` • `IP Address` • `CPU%` • `Memory(MB)` • `Uptime` • `Status`

---

## <a id="cli-command-matrix"></a>📋 CLI Command Matrix

| Command | Syntax | Description | Example |
| --- | --- | --- | --- |
| `list` / `ls` | `quickvm list` | List all Hyper-V virtual machines | `quickvm list` |
| `start` | `quickvm start <indices...>` | Start one or more VMs by index | `quickvm start 1 3 5` |
| `stop` | `quickvm stop <indices...>` | Stop one or more VMs by index | `quickvm stop 2` |
| `restart` | `quickvm restart <indices...>` | Restart one or more VMs by index | `quickvm restart 1` |
| `rdp` | `quickvm rdp <index>` | Connect to running VM via RDP | `quickvm rdp 1 -u admin` |
| `snapshot` | `quickvm snapshot <subcmd>` | Checkpoint lifecycle management | `quickvm snapshot list 1` |
| `clone` | `quickvm clone <idx> <name>` | Full independent clone of a VM | `quickvm clone 1 "DevBox"` |
| `export` | `quickvm export <idx> <path>` | Export VM files and disks | `quickvm export 1 "D:\Backup"` |
| `import` | `quickvm import <path>` | Import VM into Hyper-V | `quickvm import "D:\Backup\VM"` |
| `workspace` / `ws` | `quickvm ws <subcmd>` | Group VM batch profiles | `quickvm ws start "DevStack"` |
| `gpu` | `quickvm gpu <subcmd>` | Manage GPU-P partitioning | `quickvm gpu status` |
| `info` | `quickvm info` | Show host CPU, RAM, disk, Hyper-V | `quickvm info -d` |
| `enable` | `quickvm enable` | Enable Hyper-V Windows feature | `quickvm enable -y` |
| `update` | `quickvm update` | Check & apply updates from GitHub | `quickvm update -y` |
| `version` | `quickvm version` | Print version, date, git commit | `quickvm version` |

### Global Flags

| Flag | Short | Description | Default |
| --- | --- | --- | --- |
| `--output` | `-o` | Output format: `json`, `table`, `text` | `table` |
| `--update` | | Check for updates before running command | `false` |
| `--help` | `-h` | Display help for quickvm or any subcommand | |

---

## <a id="batch-operations--concurrency"></a>⚡ Batch Operations & Concurrency

The `start`, `stop`, and `restart` commands support bounded concurrent execution (default 4 worker threads), graceful idempotency (ignores already running/stopped VMs), and strict 60s timeout safeguards.

```powershell
# 1. Arbitrary list of indices
quickvm start 1 3 5
quickvm stop 2 4

# 2. Sequential range syntax (-r / --range)
quickvm start -r 1-5
quickvm stop --range 3-8

# 3. Comma-separated range syntax
quickvm restart -r 1,3,5

# 4. All virtual machines (-a / --all)
quickvm start --all
quickvm stop -a
quickvm restart -a
```

---

## <a id="remote-desktop-rdp-integration"></a>🔌 Remote Desktop (RDP) Integration

QuickVM automatically detects the VM's IPv4 address from Hyper-V Integration Services and launches Windows `mstsc.exe`.

```powershell
# Basic connection (prompts for Windows credentials)
quickvm rdp 1

# Connect with pre-filled username
quickvm rdp 1 -u "Administrator"
quickvm rdp 1 --user "Administrator"

# Auto-login: Save password to Windows Credential Manager (cmdkey)
quickvm rdp 1 -u "Administrator@MySecretPassword"

# Domain auto-login
quickvm rdp 1 -u "CORP\user@MySecretPassword"

# Clear cached credentials from Windows Credential Manager
quickvm rdp 1 --clean-creds
```

> 💡 **Requirement**: VM must be running, have Hyper-V Guest Integration Services enabled, and have an assigned IPv4 address.

---

## <a id="snapshot--checkpoint-management"></a>📸 Snapshot & Checkpoint Management

```powershell
# List all snapshots for VM #1
quickvm snapshot list 1

# Create a snapshot before major updates or risky operations
quickvm snapshot create 1 "Before-Node20-Upgrade"

# Restore VM state from a snapshot (⚠️ Reverts disk changes!)
quickvm snapshot restore 1 "Before-Node20-Upgrade"

# Delete snapshot to reclaim disk storage (merge checkpoints)
quickvm snapshot delete 1 "Old-Snapshot"
```

---

## <a id="vm-workspaces-group-management"></a>📦 VM Workspaces (Group Management)

Workspaces allow defining logical groups of VMs that can be inspected, started, or stopped simultaneously. Workspace files are saved under `~/.quickvm/workspaces/<name>.yaml`.

```powershell
# Create workspace with specific VM names (-v / --vms is required)
quickvm ws create "DevStack" -v "ProxyVM,AppServer,DatabaseVM"

# List all saved workspaces
quickvm ws list

# Show detailed status of all VMs inside a workspace
quickvm ws show "DevStack"

# Start all VMs in the workspace in parallel
quickvm ws start "DevStack"

# Stop all VMs in the workspace in parallel
quickvm ws stop "DevStack"

# Delete a workspace configuration (does not delete actual VMs)
quickvm ws delete "DevStack"
```

---

## <a id="vm-cloning-export--import"></a>💾 VM Cloning, Export & Import

### Clone

Performs an automated full clone (independent VM copy with a new identity):

```powershell
quickvm clone 1 "Ubuntu-Worker-2"
```

### Export

Exports the complete VM configuration, virtual hard disks (.vhdx), and snapshots:

```powershell
# Export VM 1 into a backup directory
quickvm export 1 "D:\HyperV-Backups"

# Export to current directory
quickvm export 1 .
```

### Import

Imports previously exported VMs back into Hyper-V:

```powershell
# Register in-place (keeps existing files in directory)
quickvm import "D:\HyperV-Backups\Ubuntu-Worker-2"

# Copy files to default Hyper-V storage location (-c / --copy)
quickvm import "D:\HyperV-Backups\Ubuntu-Worker-2" -c
quickvm import "D:\HyperV-Backups\Ubuntu-Worker-2" --copy

# Generate a brand new unique ID (allows importing duplicate VMs) (-n / --new-id)
quickvm import "D:\HyperV-Backups\Ubuntu-Worker-2" -c -n
quickvm import "D:\HyperV-Backups\Ubuntu-Worker-2" --copy --new-id

# Custom destination path for virtual hard disks (-v / --vhd-path)
quickvm import "D:\HyperV-Backups\Ubuntu-Worker-2" -v "E:\FastNVMe\VHDs"
```

---

## <a id="gpu-partitioning-gpu-p-passthrough"></a>🎮 GPU Partitioning (GPU-P Passthrough)

Share your host physical GPU (NVIDIA / AMD / Intel) with Hyper-V virtual machines for CUDA, 3D acceleration, and rendering.

```powershell
# 1. Verify host GPU partitioning support and VRAM limits
quickvm gpu status

# 2. Stop target VM before modifying hardware settings
quickvm stop 1

# 3. Attach GPU partition to VM #1 (requires Administrator)
quickvm gpu add 1

# 4. Show host driver files to copy into guest VM
quickvm gpu drivers

# 5. Detach GPU partition when no longer needed (VM must be stopped)
quickvm gpu remove 1
```

---

## <a id="ai-agent--automation-json-api"></a>🤖 AI Agent & Automation (JSON API)

All QuickVM commands support the `-o json` / `--output json` flag for structured, machine-readable output.

### Standard Response Envelopes

#### Success Envelope

```json
{
  "success": true,
  "data": { ... }
}
```

#### Error Envelope

```json
{
  "success": false,
  "error": {
    "code": "VM_LIST_FAILED",
    "message": "Failed to get VMs",
    "details": "Hyper-V PowerShell module is not accessible."
  }
}
```

### JSON Response Mapping by Command

| Command | `--output json` Result Structure | Key Fields |
| --- | --- | --- |
| `list` | `VMListResponse` | `vms[]` (`name`, `state`, `cpuUsage`, `memoryMB`, `uptime`, `status`, `version`, `ipAddresses`), `total` |
| `start` / `stop` / `restart` | `VMBatchResult` | `operation`, `results[]` (`index`, `name`, `success`, `message`, `error`), `successCount`, `failCount`, `totalCount` |
| `info` | `SystemInfo` | `cpu` (`name`, `cores`), `memory` (`totalMB`, `totalGB`, `freeMB`, `freeGB`, `usedMB`, `usedGB`), `disks[]`, `hyperV` (`enabled`, `status`) |
| `snapshot list` | `SnapshotListResult` | `vmName`, `vmIndex`, `snapshots[]` (`name`, `creationTime`, `snapshotType`), `total` |
| `snapshot create/restore/del` | `SnapshotOpResult` | `operation`, `vmName`, `vmIndex`, `snapshotName`, `success`, `message`, `error` |
| `clone` | `CloneResult` | `sourceName`, `sourceIndex`, `newName`, `success`, `message`, `error` |
| `export` | `ExportResult` | `vmName`, `vmIndex`, `exportPath`, `success`, `message`, `error` |
| `import` | `ImportResult` | `vmName`, `importPath`, `success`, `message`, `error` |
| `rdp` | `RDPResult` | `vmName`, `vmIndex`, `ipAddress`, `success`, `message`, `error` |
| `ws list` | `WorkspaceListResult` | `workspaces[]`, `total` |
| `ws show` | `WorkspaceShowResult` | `workspace` (`Name`, `Description`, `VMs[]`) |
| `ws start` / `ws stop` | `WorkspaceBatchResult` | `workspace`, `operation`, `results[]`, `successCount`, `failCount`, `totalCount` |
| `gpu status` | `GPUStatusResult` | `gpus[]` (`name`, `instanceId`, `driverVersion`, `totalVRAMMB`, `freeVRAMMB`, `partitionable`), `supported`, `total` |
| `gpu add` / `gpu remove` | `GPUOpResult` | `operation`, `vmName`, `vmIndex`, `success`, `message`, `error` |
| `gpu drivers` | `GPUDriversResult` | `paths[]`, `total` |
| `enable` | `EnableResult` | `alreadyEnabled`, `needsRestart`, `restartScheduled`, `success`, `message` |
| `update` | `UpdateResult` | `hasUpdate`, `currentVersion`, `latestVersion`, `releaseNotes`, `installed`, `success`, `message` |
| `version` | `VersionInfo` | `name`, `version`, `buildDate`, `gitCommit` |

### Automation Examples

```powershell
# Parse running VMs with PowerShell
(quickvm list -o json | ConvertFrom-Json).data.vms | Where-Object { $_.state -eq "Running" }

# Extract first IP address of VM 1 for automation
$ip = (quickvm list -o json | ConvertFrom-Json).data.vms[0].ipAddresses[0]

# Pre-execution auto-update check in CI/CD pipelines
quickvm --update list -o json
```

---

## <a id="system--maintenance"></a>🛠️ System & Maintenance

### System Hardware Info

```powershell
# Overview: CPU, RAM, Hyper-V feature status
quickvm info

# Include detailed disk drives & capacity (slower)
quickvm info -d
quickvm info --disk
```

### Enable Hyper-V

```powershell
# Prompt for restart if installation requires reboot
quickvm enable

# Enable and immediately reboot system (-y / --yes)
quickvm enable -y
quickvm enable --yes

# Enable without restarting (manual reboot later)
quickvm enable --no-restart
```

### Updates

```powershell
# Check and install latest version from GitHub releases
quickvm update

# Check for updates without installing
quickvm update --check-only

# Unattended auto-update (skip confirmation prompt)
quickvm update -y
quickvm update --yes
```

---

## <a id="developer-commands--makefile"></a>🔧 Developer Commands & Makefile

```powershell
# Build application
go build -o quickvm.exe
make build

# Optimized production build (stripped symbols & dwarf)
go build -ldflags="-s -w" -o quickvm.exe
make build-optimized

# Build both AMD64 and ARM64 binaries
make build-all

# Run all unit tests
go test ./...
make test

# Run tests with HTML coverage report
go test -coverprofile=coverage.out ./...
make test-coverage

# Run benchmarks
go test -bench=. -benchmem ./...
make bench

# Format codebase (gofmt + goimports)
gofmt -w .
goimports -w .
make fmt

# Run linter (requires golangci-lint v2.8.0+)
golangci-lint run ./...
make lint

# Lint and format markdown files
make lint-md
make fmt-md

# Run all checks (format + lint + test)
make check

# Full development cycle (format + test + build + run)
make dev

# Clean build artifacts
make clean
```

---

## <a id="troubleshooting--error-codes"></a>🐛 Troubleshooting & Error Codes

### Common Issues & Solutions

| Symptom / Error | Root Cause | Solution |
| --- | --- | --- |
| `❌ Failed to get VMs: Access is denied` | Shell lacks Administrator rights | Right-click Terminal / PowerShell and select **Run as Administrator** |
| `Hyper-V is not enabled` | Windows Hyper-V feature is turned off | Run `quickvm enable -y` and reboot |
| `Invalid VM index: X` | Stale index number used | Run `quickvm list` to view current 1-based index numbers |
| `Failed to get VM IP address` | VM has no IP or Integration Services inactive | Ensure VM is running and "Guest Service Interface" is enabled in VM settings |
| `A VM with name 'X' already exists` | Target clone/import name collision | Choose a unique VM name or delete old VM |
| `VM is currently running` | Target VM must be stopped before GPU/Snapshot operations | Run `quickvm stop <index>` prior to modifying hardware or restoring checkpoints |
| `No GPUs with partitioning support found` | Host GPU lacks WDDM 2.5+ or partition drivers | Update host GPU drivers or verify vendor GPU-P support |

### Complete Error Code Catalog

All error envelopes returned under `--output json` use one of the standard codes below:

| Error Code | Meaning / Cause | Command Origin |
| --- | --- | --- |
| `ADMIN_REQUIRED` | Elevated Administrator privileges required | `enable`, `gpu add`, `gpu remove` |
| `CLONE_FAILED` | Temporary export, import, or rename step failed | `clone` |
| `DIR_CREATE_FAILED` | Target export directory could not be created | `export` |
| `ENABLE_FAILED` | Enabling Windows Hyper-V optional feature failed | `enable` |
| `EXPORT_FAILED` | VM export I/O error or insufficient disk space | `export` |
| `GPU_ADD_FAILED` | Failed to assign GPU partition adapter to VM | `gpu add` |
| `GPU_CHECK_FAILED` | WMI query failed for partitionable GPUs | `gpu status`, `gpu add` |
| `GPU_DRIVERS_FAILED` | Failed to locate GPU DriverStore folder | `gpu drivers` |
| `GPU_NOT_SUPPORTED` | No installed GPUs support GPU-P partitioning | `gpu add` |
| `GPU_REMOVE_FAILED` | Failed to detach GPU partition adapter from VM | `gpu remove` |
| `IMPORT_FAILED` | Invalid export directory, missing VHD, or import error | `import` |
| `INVALID_ARGS` | Missing, excess, or malformed arguments | `start`, `stop`, `restart`, `ws create` |
| `INVALID_INDEX` | Index argument is non-numeric or out of range | `start`, `stop`, `restart`, `snapshot`, `rdp`, `clone`, `export`, `gpu` |
| `INVALID_NAME` | New VM name argument is empty | `clone` |
| `IP_GET_FAILED` | Could not query IPv4 address from Integration Services | `rdp` |
| `PATH_ERROR` | Failed to resolve relative or working directory path | `export`, `import` |
| `PATH_NOT_FOUND` | Specified import source path does not exist on disk | `import` |
| `RDP_FAILED` | Failed to launch `mstsc.exe` or update Credential Manager | `rdp` |
| `SNAPSHOT_CREATE_FAILED` | Checkpoint creation failed | `snapshot create` |
| `SNAPSHOT_DELETE_FAILED` | Checkpoint deletion or disk merge failed | `snapshot delete` |
| `SNAPSHOT_LIST_FAILED` | Failed to query snapshots for VM | `snapshot list` |
| `SNAPSHOT_RESTORE_FAILED` | Reverting VM state to checkpoint failed | `snapshot restore` |
| `STATUS_CHECK_FAILED` | Failed to query Hyper-V service status | `enable` |
| `SYSTEM_INFO_FAILED` | WMI query failed for host CPU/RAM/Disk stats | `info` |
| `UPDATE_CHECK_FAILED` | Network failure reaching GitHub releases API | `update` |
| `UPDATE_INSTALL_FAILED` | Failed downloading or extracting new release binary | `update` |
| `VM_CHECK_FAILED` | Failed checking if destination VM name exists | `clone` |
| `VM_EXISTS` | Destination VM name already in use by another VM | `clone` |
| `VM_GET_FAILED` | Failed to retrieve VM info from Hyper-V WMI | `start`, `stop`, `restart`, `snapshot`, `export`, `clone`, `rdp`, `gpu` |
| `VM_LIST_FAILED` | Failed to query list of virtual machines | `list` |
| `VM_RUNNING` | Action cannot proceed because VM is currently running | `gpu add`, `gpu remove` |
| `WORKSPACE_CREATE_FAILED` | Failed to write workspace configuration file | `ws create` |
| `WORKSPACE_DELETE_FAILED` | Failed to delete workspace configuration file | `ws delete` |
| `WORKSPACE_GET_FAILED` | Workspace file does not exist or cannot be read | `ws show`, `ws start`, `ws stop` |
| `WORKSPACE_LIST_FAILED` | Could not read workspace storage directory | `ws list` |

---

<div align="center">

**QuickVM** • Built with Go, Bubble Tea & Cobra • [GitHub Repository](https://github.com/hoangtran1411/quickvm)

</div>
