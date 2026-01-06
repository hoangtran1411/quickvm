# QuickVM - Demo & Examples

## 📸 Screenshots

### 1. TUI Mode (Interactive Interface)
```
╭───────────────────────────────────────────────────────────────────────────╮
│ 🖥️  QuickVM - Hyper-V Manager                                            │
╰───────────────────────────────────────────────────────────────────────────╯

╭───────────────────────────────────────────────────────────────────────────╮
│ Index │ Name              │ State        │ CPU%   │ Memory(MB) │ Uptime  │
│───────┼───────────────────┼──────────────┼────────┼────────────┼─────────│
│   1   │ Ubuntu-Dev        │ 🟢 Running   │   15%  │    4096    │ 05:23:41│
│   2   │ Windows-Test      │ 🔴 Off       │    0%  │       0    │ 00:00:00│
│   3   │ Docker-Host       │ 🟢 Running   │    8%  │    8192    │ 12:45:18│
│   4   │ Kali-Security     │ 🔴 Off       │    0%  │       0    │ 00:00:00│
│   5   │ CentOS-Server     │ 🟡 Paused    │    2%  │    2048    │ 02:11:33│
╰───────────────────────────────────────────────────────────────────────────╯

Status: VM list refreshed!

↑/↓: Navigate • Enter: Start VM • s: Stop VM • t: Restart VM • r: Refresh • q: Quit
```

### 2. List Command Output
```powershell
PS> quickvm list

📋 Fetching Hyper-V virtual machines...

==============================================================================
Index   Name              State        CPU%    Memory(MB)  Uptime      Status
==============================================================================
1       Ubuntu-Dev        🟢 Running   15%     4096        05:23:41    Operating normally
2       Windows-Test      🔴 Off       0%      0           00:00:00    Operating normally
3       Docker-Host       🟢 Running   8%      8192        12:45:18    Operating normally
4       Kali-Security     🔴 Off       0%      0           00:00:00    Operating normally
5       CentOS-Server     🟡 Paused    2%      2048        02:11:33    Operating normally
==============================================================================

Total VMs: 5

💡 Tip: Use 'quickvm start <index>' to start a VM
```

### 3. Start Command
```powershell
PS> quickvm start 2

🚀 Starting VM: Windows-Test (Index: 2)...
✅ VM 'Windows-Test' started successfully!
```

### 4. Start Range Command (New Feature!)
Start multiple VMs at once using the `--range` flag:

**Range format (start-end):**
```powershell
PS> quickvm start --range 1-3

🚀 Starting 3 VMs...

🚀 Starting VM: Ubuntu-Dev (Index: 1)...
✅ VM 'Ubuntu-Dev' started successfully!
🚀 Starting VM: Windows-Test (Index: 2)...
✅ VM 'Windows-Test' started successfully!
🚀 Starting VM: Docker-Host (Index: 3)...
✅ VM 'Docker-Host' started successfully!

📊 Summary: 3 started, 0 failed
```

**Comma-separated format:**
```powershell
PS> quickvm start --range 1,3,5

🚀 Starting 3 VMs...

🚀 Starting VM: Ubuntu-Dev (Index: 1)...
✅ VM 'Ubuntu-Dev' started successfully!
🚀 Starting VM: Docker-Host (Index: 3)...
✅ VM 'Docker-Host' started successfully!
🚀 Starting VM: CentOS-Server (Index: 5)...
✅ VM 'CentOS-Server' started successfully!

📊 Summary: 3 started, 0 failed
```

**Short flag syntax:**
```powershell
PS> quickvm start -r 1-5
```

### 5. Stop Command
```powershell
PS> quickvm stop 3

🛑 Stopping VM: Docker-Host (Index: 3)...
✅ VM 'Docker-Host' stopped successfully!
```

## 🎯 Use Cases & Examples

### Use Case 1: Daily Development Workflow
**Scenario**: You work with multiple development VMs and need to quickly start them each morning.

```powershell
# Morning routine - Start development environment (using --range flag)
quickvm start --range 1,3,5  # Start Ubuntu-Dev, Docker-Host, CentOS-Server

# Or use range format for consecutive VMs
quickvm start -r 1-3  # Start VMs 1, 2, 3

# Quick check
quickvm list
```

### Use Case 2: Testing Environment Management
**Scenario**: Running automated tests across different OS versions.

```powershell
# Start all test VMs at once (no loop needed!)
quickvm start --range 2,4,6,7

# Wait for VMs to boot
Start-Sleep -Seconds 60

# Run your tests...

# Clean up - Stop all test VMs (can also use --range when stop supports it)
foreach ($vm in @(2, 4, 6, 7)) {
    quickvm stop $vm
}
```

### Use Case 3: Resource Management
**Scenario**: You need to free up system resources quickly.

```powershell
# Use TUI mode for visual management
quickvm

# In TUI:
# - Navigate with ↑/↓ to see which VMs are consuming resources
# - Press 's' on running VMs you don't need
# - Press 'r' to refresh and see freed resources
```

### Use Case 4: Quick VM Status Check
**Scenario**: Check which VMs are running before starting work.

```powershell
# Quick list with color-coded status
quickvm list

# Or use PowerShell alias for faster access
Set-Alias qvm quickvm
qvm list
```

### Use Case 5: Automated Backup Workflow
**Scenario**: Stop VMs before running backups, then restart them.

```powershell
# Backup script
$criticalVMs = @(1, 3, 5)

# Stop VMs
Write-Host "Stopping VMs for backup..."
foreach ($vm in $criticalVMs) {
    quickvm stop $vm
}

# Wait for VMs to fully stop
Start-Sleep -Seconds 60

# Run your backup solution
# ... backup commands here ...

# Restart VMs
Write-Host "Restarting VMs..."
foreach ($vm in $criticalVMs) {
    quickvm start $vm
}
```

### Use Case 6: Development Lab Setup
**Scenario**: Quickly set up a multi-VM development lab.

```powershell
# Lab environment:
# VM 1: Database Server
# VM 2: Application Server  
# VM 3: Web Server
# VM 4: Load Balancer

# Start in correct order
quickvm start 1    # Database first
Start-Sleep -Seconds 45

quickvm start 2    # App server
Start-Sleep -Seconds 30

quickvm start 3    # Web server
Start-Sleep -Seconds 30

quickvm start 4    # Load balancer last

# Verify all running
quickvm list
```

## 🎨 UI/UX Features Showcase

### Color Coding
- 🟢 **Green (Running)**: VM is active and consuming resources
- 🔴 **Red (Off)**: VM is stopped and not consuming resources
- 🟡 **Yellow (Paused)**: VM is paused, can be quickly resumed

### Interactive Navigation
- **Arrow Keys**: Navigate through VM list smoothly
- **Enter Key**: Quick start for selected VM
- **Single Key Actions**: 
  - `s` for stop (red alert action)
  - `t` for restart (medium priority)
  - `r` for refresh (low impact)
- **Escape Hatch**: `q` or `Esc` for quick exit

### Real-time Information
- **CPU Usage**: Monitor VM resource consumption
- **Memory**: See allocated memory per VM
- **Uptime**: Track how long VMs have been running
- **Status**: Operating status at a glance

### Error Handling
```powershell
PS> quickvm start 99

❌ Invalid VM index: 99 (valid range: 1-5)

PS> quickvm start 1
# If VM is already running:
❌ Failed to start VM: The operation cannot be performed while 
   the virtual machine is in this state.
```

## 📊 Performance Comparison

### Traditional Method (PowerShell)
```powershell
# List VMs
Get-VM

# Start a VM
Start-VM -Name "Ubuntu-Dev"

# Check status
Get-VM -Name "Ubuntu-Dev"

# Total time: ~15-20 seconds with typing
```

### QuickVM Method
```powershell
# List VMs
quickvm list

# Start a VM
quickvm start 1

# Status already shown in list

# Total time: ~5-7 seconds
```

**Time Saved: ~65-70%**

## 🔥 Power User Tips

### 1. Create PowerShell Functions
Add to your `$PROFILE`:
```powershell
# Quick aliases
Set-Alias qvm quickvm

# Start multiple VMs
function Start-DevEnv {
    quickvm start 1
    quickvm start 3
    quickvm start 5
}

# Stop all running VMs
function Stop-AllVMs {
    $vms = @(1..10)  # Adjust based on your VM count
    foreach ($vm in $vms) {
        quickvm stop $vm 2>$null
    }
}
```

### 2. Scheduled Tasks
```powershell
# Create a scheduled task to start VMs at boot
$action = New-ScheduledTaskAction -Execute "quickvm.exe" -Argument "start 1"
$trigger = New-ScheduledTaskTrigger -AtStartup
Register-ScheduledTask -TaskName "StartVM1" -Action $action -Trigger $trigger
```

### 3. Context Menu Integration
Add QuickVM to Windows context menu (requires registry edit).

### 4. Notification Integration
Combine with Windows notifications:
```powershell
quickvm start 1
if ($LASTEXITCODE -eq 0) {
    New-BurntToastNotification -Text "VM Started", "Ubuntu-Dev is now running"
}
```

## 📈 Monitoring Dashboard Example

Create a simple monitoring loop:
```powershell
while ($true) {
    Clear-Host
    Write-Host "=== VM Dashboard ===" -ForegroundColor Cyan
    Write-Host "Last Update: $(Get-Date)" -ForegroundColor Yellow
    Write-Host ""
    
    quickvm list
    
    Start-Sleep -Seconds 30
}
```

## 🎓 Learning Path

1. **Beginner**: Use `quickvm list` and `quickvm start/stop`
2. **Intermediate**: Use TUI mode for visual management
3. **Advanced**: Create PowerShell functions and automation scripts
4. **Expert**: Integrate with monitoring tools and CI/CD pipelines

---

**Enjoy managing your Hyper-V VMs with QuickVM! 🚀**
