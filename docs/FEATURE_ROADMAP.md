# QuickVM Feature Roadmap 🗺️

> This document describes proposed features for QuickVM, organized by priority and complexity.

**Last Updated:** 2026-01-10

---

## 📊 Overview

| Tier | Description | Feature Count |
|------|-------------|---------------|
| Tier 1 | High Value, Medium Effort | 4 |
| Tier 2 | Quick Wins | 5 |
| Tier 3 | Advanced Features | 5 |
| Tier 4 | Nice to Have | 5 |

---

## 🚀 Tier 1: High Value, Medium Effort

> **High Priority** - Features that provide great value, should be implemented early.

### 1. VM Snapshots/Checkpoints ⭐ ✅ DONE (2026-01-07)

**Command:** `quickvm snapshot`

```bash
quickvm snapshot list <vm-index>              # List snapshots of a VM
quickvm snapshot create <vm-index> "name"     # Create a new snapshot
quickvm snapshot restore <vm-index> "name"    # Restore a snapshot
quickvm snapshot delete <vm-index> "name"     # Delete a snapshot
```

**Rationale:** Checkpoint management is a crucial feature when working with VMs. It allows users to:
- Save state before making changes
- Quickly restore when errors occur
- Test safely with rollback capability

**Complexity:** ⭐⭐⭐ (Medium)

**PowerShell Commands:**
```powershell
Get-VMSnapshot -VMName "VMName"
Checkpoint-VM -Name "VMName" -SnapshotName "SnapshotName"
Restore-VMSnapshot -VMName "VMName" -Name "SnapshotName" -Confirm:$false
Remove-VMSnapshot -VMName "VMName" -Name "SnapshotName"
```

---

### 2. VM Clone ⭐ ✅ DONE (2026-01-10)

**Command:** `quickvm clone`

```bash
quickvm clone <vm-index> <new-name>           # Clone VM with a new name
quickvm clone 1 "WebServer-Test"              # Example
```

**Rationale:** Quickly clone VMs for testing/development without affecting the original VM.

**Complexity:** ⭐⭐⭐ (Medium)

**PowerShell Commands:**
```powershell
Export-VM -Name "SourceVM" -Path "C:\VMs\Export"
Import-VM -Path "C:\VMs\Export\SourceVM\Virtual Machines\*.vmcx" -Copy -GenerateNewId
Rename-VM -Name "SourceVM" -NewName "NewVMName"
```

---

### 3. Export/Import VM ✅ DONE

**Command:** `quickvm export` / `quickvm import`

```bash
quickvm export <vm-index> <path>              # Export VM to a directory
quickvm export 1 "D:\Backups\VMs"

quickvm import <path>                          # Import VM from file
quickvm import "D:\Backups\VMs\WebServer"
```

**Rationale:** Backup and move VMs between machines.

**Complexity:** ⭐⭐⭐ (Medium)

**PowerShell Commands:**
```powershell
Export-VM -Name "VMName" -Path "D:\Backups"
Import-VM -Path "D:\Backups\VMName\Virtual Machines\*.vmcx"
```

---

### 4. VM Config

**Command:** `quickvm config`

```bash
quickvm config <vm-index> --memory 4GB        # Change RAM
quickvm config <vm-index> --cpu 2             # Change CPU count
quickvm config <vm-index> --memory 8GB --cpu 4  # Both
quickvm config show <vm-index>                # View current config
```

**Rationale:** Change VM RAM/CPU without opening Hyper-V Manager.

**Complexity:** ⭐⭐ (Low-Medium)

**PowerShell Commands:**
```powershell
Set-VM -Name "VMName" -MemoryStartupBytes 4GB
Set-VMProcessor -VMName "VMName" -Count 2
Get-VM -Name "VMName" | Select-Object *
```

---

## ⚡ Tier 2: Quick Wins

> **Fast to implement, high value** - Simple but useful features.

### 5. Connect to VM ⭐ ✅ DONE (Implemented via RDP & VMConnect)

**Command:** `quickvm connect`

```bash
quickvm connect <vm-index>                    # Open VM Connect GUI
quickvm connect 1
```

**Rationale:** Open VMConnect.exe directly from terminal, no need to open Hyper-V Manager.

**Complexity:** ⭐ (Low)

**Implementation:**
```go
// Simple: call vmconnect.exe
exec.Command("vmconnect.exe", "localhost", vmName).Start()
```

---

### 6. RDP Quick Connect ⭐ ✅ DONE (2026-01-10)

**Command:** `quickvm rdp`

```bash
quickvm rdp <vm-index>                        # RDP into VM
quickvm rdp 1                                 # Example
quickvm rdp 1 -u admin                        # With username
```

**Rationale:** Quick connection to VM via Windows RDP client.

**Complexity:** ⭐⭐ (Low-Medium)

---

### 7. GPU Partitioning (GPU-P) ⭐ ✅ DONE (2026-01-08)

**Command:** `quickvm gpu`

```bash
quickvm gpu status                           # Check GPU-P support
quickvm gpu add <vm-index>                   # Add GPU to VM
quickvm gpu remove <vm-index>                # Remove GPU from VM
quickvm gpu drivers                          # Show driver paths for guest
```

**Rationale:** High-performance graphics acceleration for VMs. Essential for bypass anti-VM and heavy workloads.

**Complexity:** ⭐⭐⭐⭐ (High)

---

### 8. VM Logs

**Command:** `quickvm logs`

```bash
quickvm logs <vm-index>                       # View VM event logs
quickvm logs 1 --tail 50                      # Only last 50 lines
quickvm logs 1 --follow                       # Follow mode (real-time)
```

**Rationale:** Debug and troubleshoot VM issues.

**Complexity:** ⭐⭐ (Low-Medium)

---

### 8. Bulk Operations ⭐

**Command:** `quickvm start/stop/restart --all`

```bash
quickvm start --all                           # Start all VMs
quickvm stop --all                            # Stop all VMs
quickvm restart --all                         # Restart all VMs

quickvm start --filter "Running"              # Start VMs that are Running
quickvm stop --filter "Web*"                  # Stop VMs with names starting with "Web"
```

**Rationale:** Manage multiple VMs at once.

**Complexity:** ⭐ (Low)

---

### 9. Watch Mode

**Command:** `quickvm watch`

```bash
quickvm watch                                 # Real-time monitoring TUI
quickvm watch --interval 5                    # Refresh every 5 seconds
quickvm list --watch                          # Watch mode for list command
```

**Rationale:** Monitor VM status in real-time, especially useful when waiting for VM start/stop.

**Complexity:** ⭐⭐ (Low-Medium)

---

## 🔧 Tier 3: Advanced Features

> **Advanced** - More complex features for power users.

### 10. VM Templates

**Command:** `quickvm template`

```bash
quickvm template create <vm-index> "TemplateName"   # Create template from VM
quickvm template list                               # List templates
quickvm template apply "TemplateName" "NewVMName"   # Create VM from template
quickvm template delete "TemplateName"              # Delete template
```

**Rationale:** Quickly create new VMs from pre-configured templates.

**Complexity:** ⭐⭐⭐⭐ (High)

---

### 11. Network Management

**Command:** `quickvm network`

```bash
quickvm network list                          # List Virtual Switches
quickvm network create "SwitchName" --type internal
quickvm network attach <vm-index> "SwitchName"
quickvm network detach <vm-index>
```

**Rationale:** Manage Virtual Switches and networking for VMs.

**Complexity:** ⭐⭐⭐ (Medium)

---

### 12. Storage Management

**Command:** `quickvm disk`

```bash
quickvm disk list <vm-index>                  # List disks of a VM
quickvm disk create "disk.vhdx" --size 50GB   # Create new VHD
quickvm disk resize "disk.vhdx" --size 100GB  # Resize VHD
quickvm disk attach <vm-index> "disk.vhdx"    # Attach disk to VM
quickvm disk detach <vm-index> "disk.vhdx"    # Detach disk
```

**Rationale:** Manage VHD/VHDX files.

**Complexity:** ⭐⭐⭐ (Medium)

---

### 13. Resource Quotas

**Command:** `quickvm quota`

```bash
quickvm quota set <vm-index> --max-cpu 50%    # Limit CPU
quickvm quota set <vm-index> --max-memory 4GB # Limit RAM
quickvm quota show <vm-index>                 # View current quotas
```

**Rationale:** Set resource limits, useful for lab environments.

**Complexity:** ⭐⭐⭐ (Medium)

---

### 14. Scheduled Tasks

**Command:** `quickvm schedule`

```bash
quickvm schedule start <vm-index> --at "08:00"        # Start VM at 8 AM
quickvm schedule stop <vm-index> --at "18:00"         # Stop VM at 6 PM
quickvm schedule list                                  # View schedules
quickvm schedule delete <schedule-id>                  # Delete schedule
```

**Rationale:** Automatically start/stop VMs on schedule.

**Complexity:** ⭐⭐⭐⭐ (High)

---

## 🎯 Tier 4: Nice to Have

> **Future** - Additional features when time permits.

### 15. Profile/Workspace

```bash
quickvm workspace create "Development"        # Create workspace
quickvm workspace add 1 2 3                   # Add VMs to workspace
quickvm workspace start "Development"         # Start all VMs in workspace
```

**Rationale:** Group VMs by project/purpose.

---

### 16. Remote Host Management

```bash
quickvm remote add "server1" --host 192.168.1.100
quickvm remote list
quickvm --host server1 list                   # Manage VMs on another machine
```

**Rationale:** Manage Hyper-V on other machines (remote management).

---

### 17. Metrics Export

```bash
quickvm metrics export --format prometheus    # Export metrics
quickvm metrics serve --port 9090             # HTTP endpoint for metrics
```

**Rationale:** Integration with monitoring tools (Prometheus/Grafana).

---

### 18. Configuration File

**File:** `~/.quickvmrc` or `quickvm.yaml`

```yaml
# quickvm.yaml
defaults:
  memory: 4GB
  cpu: 2
  
aliases:
  web: 1
  db: 2
  
autostart:
  - web
  - db
```

**Rationale:** Save settings and preferences.

---

### 19. Plugin System

```bash
quickvm plugin install quickvm-docker         # Install plugin
quickvm plugin list                           # List plugins
quickvm docker ps                             # Command from plugin
```

**Rationale:** Extensible architecture for custom commands.

---

## 📋 Implementation Priority (Original)

> ⚠️ **Note:** See [Refined Priority Roadmap](#-refined-priority-roadmap-brainstorming-2026-01-10) below for updated priorities based on user workflow analysis.

### Phase 1 (Week 1-2)
- [x] VM Snapshots (Tier 1, #1) ✅
- [x] VM Clone (P0) ✅
- [x] RDP Connect (P0) ✅
- [x] GPU Partitioning (Hidden/Advanced) ✅
- [x] Lazy Loading Disk Info (Performance) ✅ 2026-01-11
- [x] IP Address in TUI ✅ 2026-01-11

### Phase 2 (Week 3-4)
- [ ] Bulk Operations Enhancement (Multi-index, --all)
- [ ] VM Config (Tier 1, #4)
- [ ] Watch Mode (Tier 2, #9)

### Phase 3 (Week 5-6)
- [ ] VM Clone (Tier 1, #2)
- [x] Export/Import (Tier 1, #3) ✅ **Completed 2026-01-07**
- [ ] SSH/RDP Connect (Tier 2, #6)

### Phase 4 (Future)
- [ ] Tier 3 & 4 features

---

## 🎯 Refined Priority Roadmap (Brainstorming 2026-01-10)

> **Context:** Based on user workflow analysis - managing 3-5 VMs daily, RDP for monitoring, full clone for anti-VM detection bypass, emphasis on test coverage and error handling.

### P0: Immediate Priority

#### VM Clone (Full Clone Only) ✅ DONE

**Command:** `quickvm clone`

```bash
quickvm clone <vm-index> <new-name>           # Full clone with new name
quickvm clone 1 "WebServer-Copy"              # Example
```

**Implementation Notes:**
- Full clone only (Export → Import → Rename)
- Linked clone NOT supported (anti-VM detection concerns)
- Generate new VM ID
- Estimated time: Several minutes depending on disk size

**PowerShell Flow:**
```powershell
# 1. Export VM
Export-VM -Name "SourceVM" -Path "$env:TEMP\quickvm-clone"

# 2. Import with new ID
Import-VM -Path "...\*.vmcx" -Copy -GenerateNewId

# 3. Rename
Rename-VM -Name "SourceVM" -NewName "NewVMName"

# 4. Cleanup temp export
Remove-Item -Recurse "$env:TEMP\quickvm-clone"
```

---

#### RDP Quick Connect ✅ DONE

**Command:** `quickvm rdp`

```bash
quickvm rdp <vm-index>                        # RDP into VM
quickvm rdp 1                                 # Example
quickvm rdp 1 -u admin                        # With username (optional)
```

**Implementation Notes:**
- Manual trigger only (no auto-connect after start)
- Get VM IP first, then call mstsc.exe
- Error if VM not running or no IP available

**PowerShell to get IP:**
```powershell
(Get-VMNetworkAdapter -VMName "VMName").IPAddresses | Where-Object { $_ -match '^\d+\.\d+\.\d+\.\d+$' }
```

---

### P1: High Priority

#### Bulk Operations Enhancement

**Command:** Enhanced `quickvm start/stop/restart`

```bash
# Current (already implemented)
quickvm start 1
quickvm start --range 1-5

# New additions
quickvm start 1 2 3                           # Multiple indexes
quickvm start --all                           # All VMs
quickvm stop --all
quickvm restart --all

quickvm start --filter "Web*"                 # By name pattern (optional)
```

---

#### Workspace/Profile System

**Command:** `quickvm workspace` (alias: `quickvm ws`)

```bash
# Workspace management
quickvm workspace list                        # List all workspaces
quickvm workspace create <name>               # Create (interactive wizard)
quickvm workspace show <name>                 # View workspace details
quickvm workspace edit <name>                 # Open YAML in editor
quickvm workspace delete <name>               # Delete workspace

# Daily usage
quickvm workspace start <name>                # Start all VMs in workspace
quickvm workspace stop <name>                 # Stop all VMs
quickvm ws start dev                          # Short alias
```

**Config Location:** `~/.quickvm/workspaces/<name>.yaml`

**YAML Structure:**
```yaml
# ~/.quickvm/workspaces/dev.yaml
name: "Development Environment"
description: "Daily development VMs"

vms:
  - name: "SQL-Server"        # VM name in Hyper-V
    alias: db                 # Short alias for reference
    
  - name: "API-Server"
    alias: api
    
  - name: "Web-Frontend"
    alias: web

# Optional settings (future)
# settings:
#   start_delay: 5            # Delay between VM starts
#   stop_order: reverse       # Stop in reverse order
```

**Implementation Notes:**
- Support both wizard mode and direct YAML editing
- VMs identified by Hyper-V name (not index, as index may change)
- Alias is optional, for user convenience
- Start/stop all VMs in workspace with single command

---

### P2: Medium Priority

#### Snapshot Improvements

**Command:** Enhanced `quickvm snapshot`

```bash
# Current (keep backward compatible)
quickvm snapshot list <vm-index>
quickvm snapshot create <vm-index> "name"
quickvm snapshot restore <vm-index> "name"
quickvm snapshot delete <vm-index> "name"

# New additions
quickvm snapshot create <vm-index>            # Auto-name: "quickvm-2026-01-10-120000"
quickvm snapshot quick <vm-index>             # Alias for auto-named snapshot
quickvm snapshot restore <vm-index> --latest  # Restore most recent snapshot
```

**Implementation Notes:**
- Must maintain backward compatibility
- Auto-naming format: `quickvm-YYYY-MM-DD-HHMMSS`
- `--latest` finds snapshot with most recent creation time

---

#### VM Config

Already defined in Tier 1, #4. No changes needed.

---

### Optional/Future Enhancements

> Features to consider when time permits. Not critical for core workflow.

| Feature | Description | Complexity |
|---------|-------------|------------|
| Auto-snapshot | Snapshot before workspace start | Medium |
| Health check | Verify VM running before RDP | Low |
| Workspace import/export | Share config with team | Low |
| Default workspace | Set one workspace as default | Low |
| Global aliases | Use VM alias without workspace context | Medium |
| Workspace RDP | `quickvm ws rdp dev db` - RDP to specific VM | Low |

---

## 🔗 Related Documents

- [README.md](../README.md) - Project overview
- [DEVELOPER.md](DEVELOPER.md) - Developer guide
- [CONTRIBUTING.md](CONTRIBUTING.md) - How to contribute

---

## 🛠️ Technical Improvements & Refactoring

### 2026-01-10: Dependency Injection for PowerShell
- **Change**: Refactored `hyperv` package to use `ShellExecutor` interface.
- **Reason**: Improved unit test coverage and testability without requiring a Windows environment.
- **Impact**: Developers can now mock PowerShell output in tests using `MockRunner`.

### 2026-01-11: Pivot to Simplicity
- **Decision**: Remove/Hide complex features (GPU, detailed Disk mgmt) from main view.
- **Change**: `GetSystemInfo` now accepts `includeDisk` flag to speed up default loading.
- **Change**: Added `IP Address` column to TUI for better utility.
- **Goal**: Focus on being a "Quick" VM manager, not a comprehensive infrastructure tool.

> 💡 **Note:** This is a planned roadmap and may change based on user feedback.
