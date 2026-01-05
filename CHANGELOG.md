# Changelog

All notable changes to QuickVM will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- 🔄 Auto-update functionality
  - `quickvm update` command to check and install updates
  - `--update` flag to check for updates before running any command
  - Download updates directly from GitHub releases
  - Auto-backup before updating
  - Support for both AMD64 and ARM64 architectures

### Planned Features
- VM snapshot management
- Create new VMs
- Configuration file support
- Remote Hyper-V server support

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
