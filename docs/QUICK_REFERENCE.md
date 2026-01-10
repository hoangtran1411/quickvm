# QuickVM - Quick Reference Card

## 🚀 Quick Start
```powershell
# Build
go build -o quickvm.exe

# Install
.\install.ps1 -InstallLocation User

# Run
quickvm
```

## 📋 CLI Commands

| Command | Description | Example |
|---------|-------------|---------|
| `quickvm` | Launch interactive TUI | `quickvm` |
| `quickvm list` | List all VMs | `quickvm list` |
| `quickvm start <index>` | Start VM by index | `quickvm start 1` |
| `quickvm stop <index>` | Stop VM by index | `quickvm stop 1` |
| `quickvm restart <index>` | Restart VM by index | `quickvm restart 1` |
| `quickvm update` | Check and install updates | `quickvm update` |
| `quickvm version` | Show version info | `quickvm version` |
| `quickvm help` | Show help | `quickvm help` |

### Update Flags
| Flag | Description | Example |
|------|-------------|---------|
| `--update` | Check for updates before running any command | `quickvm --update list` |
| `--yes, -y` | Auto-install without prompting | `quickvm update -y` |
| `--check-only` | Only check, don't install | `quickvm update --check-only` |


## ⌨️ TUI Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate VMs |
| `Enter` | Start selected VM |
| `s` | Stop selected VM |
| `t` | Restart selected VM |
| `r` | Refresh VM list |
| `q` / `Esc` | Quit |

## 🎨 Status Colors

| Color | State | Meaning |
|-------|-------|---------|
| 🟢 Green | Running | VM is active |
| 🔴 Red | Off | VM is stopped |
| 🟡 Yellow | Paused | VM is paused |

## 📁 Project Structure
```
quickvm/
├── cmd/          # CLI commands
├── hyperv/       # Hyper-V integration
├── ui/           # TUI interface
└── main.go       # Entry point
```

## 🔧 Development

### Build
```powershell
go build -o quickvm.exe                    # Normal build
go build -ldflags="-s -w" -o quickvm.exe   # Optimized
```

### Test
```powershell
go test ./...              # All tests
go test -cover ./...       # With coverage
go test -bench=. ./...     # Benchmarks
```

### Format & Lint
```powershell
go fmt ./...               # Format code
golangci-lint run          # Run linter
```

## 📖 Documentation Files

| File | Purpose |
|------|---------|
| `README.md` | Main documentation (English) |
| `DEMO.md` | Examples & use cases |
| `DEVELOPER.md` | Developer notes |
| `WORKFLOW.md` | Dev & deploy workflow |
| `CONTRIBUTING.md` | Contribution guide |
| `PROJECT_SUMMARY.md` | Project overview |

## 🛠️ Common Tasks

### Install Globally
```powershell
.\install.ps1 -InstallLocation User -CreateAlias
```

### Create Alias
```powershell
# Add to PowerShell $PROFILE
Set-Alias qvm quickvm
```

### Start Multiple VMs
```powershell
1..3 | ForEach-Object { quickvm start $_ }
```

### Stop All VMs
```powershell
quickvm list  # Get indices
1..5 | ForEach-Object { quickvm stop $_ }
```

## 🐛 Troubleshooting

### Issue: "Failed to get VMs"
**Solution**: 
- Run as Administrator
- Check Hyper-V is enabled
- Verify VMs exist: `Get-VM`

### Issue: "Invalid VM index"
**Solution**: 
- Run `quickvm list` to get current indices
- Indices change when VMs are added/removed

### Issue: Command not found
**Solution**:
- Add to PATH or use full path
- Run `.\quickvm.exe` in current directory

## ⚡ Power User Tips

### 1. Quick Status Check
```powershell
quickvm list | Select-String "Running"
```

### 2. Auto-start VMs at Boot
```powershell
$action = New-ScheduledTaskAction -Execute "quickvm.exe" -Argument "start 1"
$trigger = New-ScheduledTaskTrigger -AtStartup
Register-ScheduledTask -TaskName "AutoStartVM" -Action $action -Trigger $trigger
```

### 3. Monitoring Loop
```powershell
while ($true) {
    Clear-Host
    quickvm list
    Start-Sleep -Seconds 30
}
```

### 4. Notification on Start
```powershell
quickvm start 1
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ VM Started!" -ForegroundColor Green
}
```

## 🔗 Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/charmbracelet/bubbletea` | TUI framework |
| `github.com/charmbracelet/bubbles` | TUI components |
| `github.com/charmbracelet/lipgloss` | Terminal styling |

## 📊 Performance

- **Startup**: < 100ms
- **PowerShell exec**: 200-500ms
- **Total operation**: 1-2 seconds
- **Memory usage**: ~10-20MB

## 🎯 Use Cases

1. **Development**: Quick start/stop dev VMs
2. **Testing**: Manage test environments
3. **Demos**: Rapidly switch between environments
4. **Resource Management**: Monitor and control VM resources
5. **Automation**: Script VM operations

## 📞 Support

- **Documentation**: See `.md` files in project root
- **Issues**: GitHub Issues
- **Discussions**: GitHub Discussions

## 🌟 Quick Tips

💡 **Tip 1**: Use TUI mode for exploration, CLI for automation

💡 **Tip 2**: Create PowerShell functions for common workflows

💡 **Tip 3**: Use the alias `qvm` for faster typing

💡 **Tip 4**: Always run as Administrator

💡 **Tip 5**: Check `quickvm list` before operations

---

**QuickVM** - Fast Hyper-V Virtual Machine Manager
Version 1.0.0 | Made with ❤️ using Go
