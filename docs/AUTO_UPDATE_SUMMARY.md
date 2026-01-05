# QuickVM Auto-Update Feature - Summary

## ✅ Completed

Đã thêm thành công tính năng **Auto-Update** vào QuickVM!

## 🎯 Tính Năng Mới

### 1️⃣ Update Command
```powershell
# Check và install updates
quickvm update

# Chỉ check, không install
quickvm update --check-only

# Auto-install không hỏi
quickvm update -y
```

### 2️⃣ Global --update Flag
```powershell
# Check update trước khi chạy bất kỳ lệnh nào
quickvm --update list
quickvm --update start 1
quickvm --update
```

### 3️⃣ Tính Năng Nổi Bật

- ✅ **Auto-check** version từ GitHub Releases
- ✅ **Auto-download** phiên bản mới
- ✅ **Auto-backup** trước khi update
- ✅ **Auto-rollback** nếu update fail
- ✅ **Architecture detection** (AMD64/ARM64)
- ✅ **Silent fail** để không làm gián đoạn workflow
- ✅ **User-friendly messages** với emojis

## 📁 Files Created/Modified

### New Files
```
quickvm/
├── updater/
│   ├── updater.go          ← Update logic
│   └── updater_test.go     ← Unit tests
├── cmd/
│   └── update.go           ← Update command
└── docs/
    └── AUTO_UPDATE.md      ← Documentation
```

### Modified Files
```
- cmd/root.go               ← Added --update flag
- docs/QUICK_REFERENCE.md   ← Added update commands
- docs/README.md            ← Added AUTO_UPDATE link
- README.md                 ← Added update section
- CHANGELOG.md              ← Documented new feature
```

## 🔧 Technical Implementation

### 1. Updater Package (`updater/updater.go`)
**Responsibilities:**
- Check GitHub Releases API
- Compare versions
- Download new binary
- Create backup
- Replace executable
- Rollback on failure

**Key Functions:**
```go
NewUpdater(version) *Updater
CheckForUpdates() (release, hasUpdate, error)
DownloadAndInstall(release) error
```

### 2. Update Command (`cmd/update.go`)
**Flags:**
- `--yes, -y`: Auto-install without prompt
- `--check-only`: Only check, don't install

**Behavior:**
```
1. Check for updates
2. Display release info
3. Prompt user (unless -y)
4. Download and install
5. Show success message
```

### 3. Global --update Flag (`cmd/root.go`)
**Implementation:**
- `PersistentPreRun`: Runs before any command
- Checks for updates if `--update` flag is set
- Prompts to install if available
- Continues with original command or exits if updated

## 📊 Update Flow

```
User runs: quickvm update
         ↓
Check GitHub API
         ↓
Compare versions
         ↓
    New version?
    ├─ No ──> "Already latest!"
    └─ Yes
         ↓
Show release info
         ↓
Prompt user (unless -y)
         ↓
Download binary
         ↓
Create backup (.backup)
         ↓
Replace executable
         ↓
Verify + cleanup
         ↓
"Update complete!"
```

## 🎨 User Experience

### Example Output

#### When Update Available:
```
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

#### When Already Latest:
```
PS> quickvm update

🔍 Checking for updates...
✅ You are already using the latest version!
   Current version: 1.0.0
```

#### With --update Flag:
```
PS> quickvm --update list

✅ The current version is the latest!

📋 Fetching Hyper-V virtual machines...
...
```

## 🔒 Security Features

### 1. Automatic Backup
- Creates `.backup` before updating
- Automatically restored if update fails
- Deleted after successful update

### 2. Rollback on Failure
```go
if err := installUpdate(); err != nil {
    restoreBackup()  // Auto rollback
}
```

### 3. SHA256 Checksums
- GitHub provides SHA256 for releases
- Can be verified manually

### 4. Official Source Only
- Only downloads from `github.com/hoangtran1411/quickvm`
- Uses GitHub Releases API

## 📚 Documentation

### AUTO_UPDATE.md Includes:
- ✅ All commands and options
- ✅ Usage examples (5+ scenarios)
- ✅ Troubleshooting guide
- ✅ Security details
- ✅ Automation examples
- ✅ FAQ section
- ✅ Technical details

### Updated Docs:
- ✅ QUICK_REFERENCE.md - Command table
- ✅ README.md - Usage section
- ✅ docs/README.md - New link
- ✅ CHANGELOG.md - Feature log

## 🧪 Testing

### Manual Testing Commands:
```powershell
# Test update command
.\quickvm.exe update --help
.\quickvm.exe update --check-only

# Test global flag
.\quickvm.exe --update version

# Test version display
.\quickvm.exe version
```

### Unit Tests:
```powershell
# Run tests
go test ./updater

# With coverage
go test -cover ./updater
```

## 🚀 Usage Scenarios

### 1. Regular User
```powershell
# Check monthly
quickvm update --check-only

# Install when ready
quickvm update
```

### 2. Automated Scripts
```powershell
# Always use latest
quickvm update -y
quickvm start 1
```

### 3. CI/CD
```powershell
# Update in pipeline
quickvm update -y || true  # Don't fail pipeline
```

### 4. Safety-First
```powershell
# Check before important work
quickvm --update list
```

## 💡 Design Decisions

### Why GitHub Releases API?
- ✅ Official release mechanism
- ✅ Built-in version management
- ✅ Automatic asset hosting
- ✅ Free and reliable

### Why Prompt by Default?
- ✅ User control
- ✅ See release notes first
- ✅ Choose update timing
- ✅ Automation still possible ( -y flag)

### Why Backup?
- ✅ Safety net
- ✅ Auto-rollback capability
- ✅ Zero-risk updates
- ✅ No manual intervention needed

### Why Silent Fail on --update?
- ✅ Don't block workflow
- ✅ Offline scenarios
- ✅ Network issues gracefully handled
- ✅ Original command still runs

## 📈 Benefits

### For Users:
- ✅ Easy updates (one command)
- ✅ Always latest features
- ✅ Bug fixes automatically
- ✅ No manual download needed

### For Developers:
- ✅ Better adoption of new versions
- ✅ Less support for old bugs
- ✅ Easy to push critical fixes
- ✅ Clear update metrics

### For Project:
- ✅ Professional feature
- ✅ Better user experience
- ✅ Reduced documentation burden
- ✅ Competitive advantage

## 🔮 Future Enhancements

### Potential Additions:
- [ ] Auto-update on startup (opt-in)
- [ ] Update notifications in TUI
- [ ] Release notes in app
- [ ] Update changelog viewer
- [ ] Downgrade capability
- [ ] Beta channel support
- [ ] Update size preview
- [ ] Progress bar for download

## 📝 Commit Message

```
feat: add auto-update functionality

- Add updater package for GitHub releases integration
- Add 'quickvm update' command with --yes and --check-only flags
- Add --update global flag to check before any command
- Automatic backup and rollback on update failure
- Support for both AMD64 and ARM64 architectures
- Comprehensive documentation in docs/AUTO_UPDATE.md
- Updated QUICK_REFERENCE, README, and CHANGELOG

Features:
- Check latest version from GitHub releases
- Download and install updates automatically
- Safe update with automatic backup
- User-friendly prompts and messages
- Silent fail for non-critical errors
```

## ✅ Ready to Commit

All files ready for commit:
```powershell
git add .
git commit -m "feat: add auto-update functionality"
git push origin main
```

---

**Auto-Update Feature Complete! 🎉**

User có thể:
- ✅ Check updates: `quickvm update --check-only`
- ✅ Install updates: `quickvm update`
- ✅ Auto-check: `quickvm --update <command>`

Message khi không có update: **"The current version is the latest!"** ✅ (Đúng như yêu cầu!)

