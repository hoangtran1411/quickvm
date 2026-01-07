# QuickVM Feature Roadmap 🗺️

> This document describes proposed features for QuickVM, organized by priority and complexity.

**Last Updated:** 2026-01-07

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

### 1. VM Snapshots/Checkpoints ⭐ ✅ DONE

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

### 2. VM Clone

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

### 5. Connect to VM ⭐

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

### 6. SSH/RDP Quick Connect

**Command:** `quickvm ssh` / `quickvm rdp`

```bash
quickvm ssh <vm-index>                        # SSH into VM (Linux)
quickvm ssh 1 -u admin                        # With username

quickvm rdp <vm-index>                        # RDP into VM (Windows)
quickvm rdp 1
```

**Rationale:** Quick connection to VM if IP address is known.

**Complexity:** ⭐⭐ (Low-Medium)

**Requirement:** Need to get VM IP first:
```powershell
(Get-VMNetworkAdapter -VMName "VMName").IPAddresses
```

---

### 7. VM Logs

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

## 📋 Implementation Priority

### Phase 1 (Week 1-2)
- [ ] VM Connect (Tier 2, #5)
- [ ] Bulk Operations (Tier 2, #8)
- [x] VM Snapshots (Tier 1, #1) ✅ **Completed 2026-01-07**

### Phase 2 (Week 3-4)
- [ ] VM Config (Tier 1, #4)
- [ ] Watch Mode (Tier 2, #9)
- [ ] VM Logs (Tier 2, #7)

### Phase 3 (Week 5-6)
- [ ] VM Clone (Tier 1, #2)
- [x] Export/Import (Tier 1, #3) ✅ **Completed 2026-01-07**
- [ ] SSH/RDP Connect (Tier 2, #6)

### Phase 4 (Future)
- [ ] Tier 3 & 4 features

---

## 🔗 Related Documents

- [README.md](../README.md) - Project overview
- [DEVELOPER.md](DEVELOPER.md) - Developer guide
- [CONTRIBUTING.md](CONTRIBUTING.md) - How to contribute

---

> 💡 **Note:** This is a planned roadmap and may change based on user feedback.
