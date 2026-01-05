# QuickVM Auto-Update Guide

## 🔄 Overview

QuickVM includes a built-in auto-update feature that allows you to easily keep your installation up-to-date with the latest releases from GitHub.

## ✨ Features

- ✅ Check for updates from GitHub releases
- ✅ Download and install updates automatically
- ✅ Automatic backup before updating
- ✅ Support for both AMD64 and ARM64 architectures
- ✅ Rollback capability if update fails
- ✅ Optional automatic update checks

## 📋 Commands

### Basic Update Command

```powershell
# Check for updates and install if available
quickvm update
```

The command will:
1. Check GitHub releases for the latest version
2. Compare with your current version
3. Prompt you to install if an update is available
4. Download and install the new version
5. Create a backup of the current version

### Update Options

#### Check Only (No Install)
```powershell
# Just check if updates are available
quickvm update --check-only
```

This will show you:
- Current version
- Latest available version
- Release notes
- But won't install anything

#### Auto-Install (No Prompt)
```powershell
# Automatically install without asking
quickvm update -y
# or
quickvm update --yes
```

Useful for automation scripts.

### Global --update Flag

Check for updates before running any command:

```powershell
# Check for updates, then list VMs
quickvm --update list

# Check for updates, then start VM
quickvm --update start 1

# Check for updates, then launch TUI
quickvm --update
```

If an update is available, you'll be prompted to install it before the command runs.

## 🎯 Usage Examples

### Example 1: Manual Update Check
```powershell
PS> quickvm update

🔍 Checking for updates...
🎉 New version available: v1.1.0 (current: v1.0.0)

❓ Do you want to install this update? [y/N]: y

📦 Downloading QuickVM v1.1.0 (8 MB)...
✅ Download complete!
📦 Creating backup...
🔄 Installing update...
✅ Successfully updated to version v1.1.0!
🔄 Please restart QuickVM to use the new version.
```

### Example 2: Check Only
```powershell
PS> quickvm update --check-only

🔍 Checking for updates...
🎉 New version available: v1.1.0 (current: v1.0.0)

📋 Release Notes:
- New feature: VM snapshots
- Bug fix: Memory leak in TUI
- Performance improvements

💡 Run 'quickvm update' without --check-only to install
```

### Example 3: Already Up-to-Date
```powershell
PS> quickvm update

🔍 Checking for updates...
✅ You are already using the latest version!
   Current version: 1.0.0
```

### Example 4: With Global Flag
```powershell
PS> quickvm --update list

✅ The current version is the latest!

📋 Fetching Hyper-V virtual machines...
# ...rest of list output
```

### Example 5: Auto-Install in Script
```powershell
# Update script
quickvm update -y

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ QuickVM updated successfully"
} else {
    Write-Host "❌ Update failed"
}
```

## 🔒 Security

### Checksum Verification
Updates are downloaded from official GitHub releases with SHA256 checksums for verification.

### Automatic Backup
Before installing an update, QuickVM creates a backup of the current executable:
- Backup location: `quickvm.exe.backup`
- Automatically restored if update fails
- Deleted after successful update

### Rollback
If an update fails, QuickVM automatically restores the backup:

```powershell
📦 Creating backup...
🔄 Installing update...
❌ Update failed: permission denied
🔙 Restoring backup...
✅ Rollback successful
```

## 🛠️ Troubleshooting

### Update Check Fails

**Problem**: Cannot check for updates

```powershell
❌ Failed to check for updates: connection timeout
💡 Tip: Check your internet connection and try again
```

**Solutions**:
- Check internet connection
- Verify GitHub is accessible
- Check firewall settings
- Try again later

### Download Fails

**Problem**: Update download interrupted

**Solutions**:
- Check internet connection
- Ensure enough disk space
- Try again with `quickvm update`

### Permission Denied

**Problem**: Cannot install update

```powershell
❌ Update failed: Access is denied
```

**Solutions**:
- Run PowerShell as Administrator
- Close other instances of QuickVM
- Temporarily disable antivirus

### Update Not Detected

**Problem**: QuickVM says already up-to-date but you know there's a new version

**Solutions**:
- Check your current version: `quickvm version`
- Check GitHub releases manually
- Version format must match (v1.0.0)

## 🚀 Automation

### Schedule Daily Update Checks

```powershell
# Create scheduled task to check for updates daily
$action = New-ScheduledTaskAction -Execute "quickvm.exe" -Argument "update --check-only"
$trigger = New-ScheduledTaskTrigger -Daily -At 9am
Register-ScheduledTask -TaskName "QuickVM Update Check" -Action $action -Trigger $trigger
```

### Script with Auto-Update

```powershell
# Always use latest version
function Start-QuickVM {
    # Update first
    quickvm update -y
    
    # Then run command
    quickvm $args
}

# Usage
Start-QuickVM list
Start-QuickVM start 1
```

### CI/CD Integration

```powershell
# In your CI/CD pipeline
# Download latest QuickVM
Invoke-WebRequest -Uri "https://github.com/hoangtran1411/quickvm/releases/latest/download/quickvm-windows-amd64.exe" -OutFile "quickvm.exe"

# Or use auto-update
quickvm update -y
```

## 📊 Update Process Flow

```
1. quickvm update
         ↓
2. Check GitHub API
         ↓
3. Compare versions
         ↓
4. New version? ──No──> "Already latest"
         ↓ Yes
5. Prompt user (unless -y)
         ↓
6. Download new version
         ↓
7. Create backup
         ↓
8. Replace executable
         ↓
9. Verify installation
         ↓
10. Delete backup
         ↓
11. Success message
```

## 🔧 Technical Details

### Update Source
- GitHub Releases API
- Repository: `hoangtran1411/quickvm`
- API: `https://api.github.com/repos/hoangtran1411/quickvm/releases/latest`

### Architecture Detection
QuickVM automatically detects your system architecture:
- AMD64 (64-bit Intel/AMD)
- ARM64 (ARM 64-bit)

And downloads the appropriate binary.

### Version Comparison
Versions are compared as strings after removing the "v" prefix:
- `v1.0.0` vs `v1.0.1` → Update available
- `v1.0.0` vs `v1.0.0` → Already latest

### File Locations
- Current executable: `quickvm.exe`
- During update: `quickvm-update-*.exe` (temp)
- Backup: `quickvm.exe.backup`

## 💡 Best Practices

1. **Check Regularly**: Run `quickvm update --check-only` periodically
2. **Read Release Notes**: Review changes before updating
3. **Backup Important Work**: Close all VMs before updating
4. **Use Administrator Rights**: Ensures smooth installation
5. **Test After Update**: Verify functionality with `quickvm version`

## 🎓 FAQ

**Q: How often should I update?**  
A: Check for updates before important tasks or weekly.

**Q: Will updates break my VMs?**  
A: No, QuickVM only manages VMs, it doesn't modify them.

**Q: Can I downgrade?**  
A: Manual download from GitHub releases required.

**Q: Do updates happen automatically?**  
A: No, you must run `quickvm update` or use `--update` flag.

**Q: What if I'm offline?**  
A: Update check will fail gracefully, QuickVM continues working.

**Q: Can I update from a specific version?**  
A: Always updates to latest. For specific versions, download from GitHub.

## 📞 Support

If you encounter issues with updates:

1. Check [CHANGELOG.md](../CHANGELOG.md) for known issues
2. Open an issue on GitHub
3. Include error message and `quickvm version` output

---

**QuickVM** - Fast Hyper-V Virtual Machine Manager  
Made with ❤️ using Go
