# Changelog

All notable changes to QuickVM will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.5.0] - 2026-09-05

### Security & Hardening

- 🛡️ **PowerShell Injection Defense**: Introduced stateful parameter parsing (`formatCmdletScript`) and safe single-quote escaping in `RunCmdlet`, securely handling dash-prefixed values and arbitrary inputs without shell injection risks.
- 🔒 **Cryptographic Updater Verification**: Added SHA256 checksum verification (`.sha256`) before replacing binaries, secured cleanup script execution with escaped single quotes, and prioritized exact binary matches (`quickvm-<assetName>`) over checksum files.
- 🗂️ **Workspace Path Traversal Protection**: Added strict workspace name validation (`ValidateWorkspaceName`) rejecting path traversal sequences (`..`, slashes, colons, wildcards) across all workspace operations.
- 🚦 **Input Validation**: Added guard clauses rejecting empty or whitespace-only VM and snapshot names prior to issuing PowerShell calls.

### Added

- 🤖 **Complete Machine-Readable JSON Output**: Implemented structured JSON models (`--output json`) for `gpu`, `workspace`, `enable`, `update`, and snapshot commands (`restore`/`delete`).
- ⚡ **Native Windows Token Elevation**: Replaced slow `powershell.exe` execution in `IsRunningAsAdmin` with native Windows access token checks (`windows.OpenProcessToken`), cutting check latency from ~400ms to <1µs.
- 🚀 **Concurrent Workspace Execution**: Upgraded workspace batch VM start/stop operations to use bounded concurrency (`errgroup` with limit 4) matching core batch commands.
- 🖥️ **Multi-Socket Server Support**: Enhanced `getCPUInfo` to cleanly aggregate multi-socket processor arrays returned by `Win32_Processor` without unmarshaling crashes.

### Fixed

- 🛑 **Standardized Non-Zero Exit Codes**: Fixed CLI commands (`list`, `clone`, `export`, `rdp`, batch operations) to return exit code `1` on failure instead of exiting `0`.
- ⚠️ **Batch False-Success Resolution**: Batch VM operations now emit `"success": false` with error details in JSON mode whenever one or more operations fail.
- 🔇 **Headless Non-Blocking Mode**: Auto-update check in `root.go` skips interactive `fmt.Scanln` prompts when running in JSON mode.
- 🧪 **Test Suite Integrity**: Corrected snapshot unit tests that previously passed false-positively due to missing elevation checks.
- 📦 **Release Linker Flags**: Fixed `.github/workflows/release.yml` linker flags to correctly target `quickvm/cmd.Version`.
- 🧹 **Dangling Documentation References**: Cleaned up obsolete references to deleted `HUONG_DAN.md` in `Makefile` and `install.ps1`.

## [1.4.0] - 2026-09-03

### Added

- 🚀 **Concurrent Batch Operations**: Bounded parallel worker pool (4 workers via `errgroup`) for VM start, stop, and restart commands, cutting batch execution time by up to 75%.
- 🛡️ **Idempotent VM Starts**: Automatic detection and graceful handling of already-running Hyper-V VMs without reporting false failures.
- ⏱️ **Context & Timeout Safety**: Strict 60s operation timeouts across CLI batch operations and interactive TUI actions.
- ⚡ **Streamlined TUI VM Starts**: Calling `StartVMByName` directly from the table UI to bypass redundant `GetVMs` index lookups.
- 🤖 **AI-Agent & Automation Support**: Comprehensive `--output json` (`-o json`) support across commands for machine-readable automation.

### Added (Previous)

- 📸 **VM Snapshot Management** (2026-01-07)
  - `quickvm snapshot list <vm-index>` - List snapshots for a VM
  - `quickvm snapshot create <vm-index> <name>` - Create a new snapshot
  - `quickvm snapshot restore <vm-index> <name>` - Restore VM to a snapshot
  - `quickvm snapshot delete <vm-index> <name>` - Delete a snapshot

- 📦 **Export/Import VMs** (2026-01-07)
  - `quickvm export <vm-index> <path>` - Export VM to a directory
  - `quickvm import <path>` - Import VM from an export
  - `--copy` flag to copy VM files instead of registering in place
  - `--new-id` flag to generate a new unique ID
  - `--vhd-path` flag to specify custom VHD destination

- 🔄 **VM Clone** (2026-01-10)
  - `quickvm clone <vm-index> <new-name>` - Full clone VM with new identity
  - Export → Import with Copy → Generate new ID → Rename workflow
  - Completely independent clone (no linked disks)
  - Automatic temp cleanup

- 🔗 **RDP Quick Connect** (2026-01-10)
  - `quickvm rdp <vm-index>` - Open RDP connection to a VM
  - `quickvm rdp <vm-index> -u <username>` - With username hint
  - Automatically gets VM IP address
  - Opens Windows Remote Desktop client (mstsc.exe)

- 🔄 Auto-update functionality
  - `quickvm update` command to check and install updates
  - `--update` flag to check for updates before running any command
  - Download updates directly from GitHub releases
  - Auto-backup before updating
  - Support for both AMD64 and ARM64 architectures

### Planned Features

- Bulk operations (`--all` flag)
- Workspace/Profile system

## [1.0.0] - 2026-01-05

### Added

- 🎨 Beautiful TUI interface with Bubble Tea framework
- ⚡ CLI commands for quick VM operations
- 📊 Real-time VM monitoring (CPU, Memory, Uptime)
- 🎯 Index-based VM operations for speed
- 🔄 Interactive table view with keyboard navigation
- 📚 Comprehensive documentation in English and Vietnamese
- 🧪 Unit tests for core functionality
- 🚀 GitHub Actions CI/CD workflows
- 📦 Automated release builds for Windows AMD64 and ARM64

### Commands Implemented

- `quickvm` - Launch interactive TUI
- `quickvm list` - List all VMs
- `quickvm start <index>` - Start VM by index
- `quickvm stop <index>` - Stop VM by index
- `quickvm restart <index>` - Restart VM by index
- `quickvm version` - Show version information

### Documentation

- README.md - Main documentation
- QUICK_REFERENCE.md - Quick reference card
- HUONG_DAN.md - Vietnamese user guide
- DEMO.md - Examples and use cases
- DEVELOPER.md - Developer guide
- WORKFLOW.md - Development workflow
- CONTRIBUTING.md - Contributing guidelines
- PROJECT_SUMMARY.md - Complete overview

### Technical Details

- Built with Go 1.25.2
- Hyper-V integration via PowerShell
- Clean architecture with separation of concerns
- Color-coded VM states (🟢 Running, 🔴 Off, 🟡 Paused)
- Comprehensive error handling
- MIT License

### Performance

- Startup time: < 100ms
- Operation time: 1-2 seconds
- Memory usage: ~10-20MB
- Binary size: ~6-8MB

### Known Limitations

- Requires Administrator privileges
- Windows-only (Hyper-V specific)
- VM indices change when VMs are added/removed

---

## Version Guidelines

We use [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality in a backwards compatible manner
- **PATCH** version for backwards compatible bug fixes

## Links

- [GitHub Repository](https://github.com/hoangtran1411/quickvm)
- [Issue Tracker](https://github.com/hoangtran1411/quickvm/issues)
- [Releases](https://github.com/hoangtran1411/quickvm/releases)
