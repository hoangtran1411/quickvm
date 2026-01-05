# QuickVM - Restructure & CI/CD Setup

## Changes Made

### 📁 Documentation Restructure
- ✅ Created `docs/` folder for better organization
- ✅ Moved all documentation files to `docs/` (except README.md)
  - CONTRIBUTING.md → docs/CONTRIBUTING.md
  - DEMO.md → docs/DEMO.md
  - DEVELOPER.md → docs/DEVELOPER.md
  - HUONG_DAN.md → docs/HUONG_DAN.md
  - PROJECT_SUMMARY.md → docs/PROJECT_SUMMARY.md
  - QUICK_REFERENCE.md → docs/QUICK_REFERENCE.md
  - WORKFLOW.md → docs/WORKFLOW.md
- ✅ Created docs/README.md as documentation index

### 🚀 GitHub Actions CI/CD
- ✅ Created `.github/workflows/build.yml`
  - Automated builds for Windows AMD64 and ARM64
  - Runs on every push to main
  - Runs on pull requests
  - Includes linting with golangci-lint
  - Test coverage with Codecov integration
  - Uploads build artifacts (30 days retention)

- ✅ Created `.github/workflows/release.yml`
  - Automated releases on version tags (v*.*.*)
  - Builds optimized binaries for both architectures
  - Creates complete release packages with docs
  - Generates SHA256 checksums
  - Auto-uploads to GitHub Releases
  - Professional release notes

### 📝 Updated Files
- ✅ Updated README.md
  - Added Go version 1.25.2 badge
  - Added GitHub Actions build status badge
  - Added Documentation section with links to docs/
  
- ✅ Updated .gitignore
  - Added release/ folder
  - Added coverage files (coverage.out, coverage.html)
  - Added release artifacts (*.zip, *.sha256)
  - Added quickvm-*.exe pattern

- ✅ Created CHANGELOG.md
  - Following Keep a Changelog format
  - Version 1.0.0 documented
  - Future roadmap included

## New File Structure

```
quickvm/
├── .github/
│   └── workflows/
│       ├── build.yml          # CI/CD build workflow
│       └── release.yml        # Release automation
│
├── docs/                      # Documentation folder
│   ├── README.md             # Documentation index
│   ├── CONTRIBUTING.md       # Contributing guide
│   ├── DEMO.md              # Examples
│   ├── DEVELOPER.md         # Developer guide
│   ├── HUONG_DAN.md         # Vietnamese guide
│   ├── PROJECT_SUMMARY.md   # Project summary
│   ├── QUICK_REFERENCE.md   # Quick reference
│   └── WORKFLOW.md          # Workflow guide
│
├── cmd/                      # CLI commands
├── hyperv/                   # Hyper-V integration
├── ui/                       # TUI interface
│
├── README.md                 # Main documentation
├── CHANGELOG.md             # Version history
├── LICENSE                  # MIT License
├── go.mod                   # Go modules
├── main.go                  # Entry point
└── ...other files
```

## CI/CD Features

### Build Workflow (build.yml)
**Triggers:**
- Push to main branch
- Pull requests to main
- Release creation

**Jobs:**
1. **build-windows**: Build AMD64 and ARM64 binaries
2. **lint**: Code quality checks with golangci-lint
3. **test-coverage**: Run tests and upload coverage

**Artifacts:**
- quickvm-windows-amd64.exe (30 days)
- quickvm-windows-arm64.exe (30 days)

### Release Workflow (release.yml)
**Triggers:**
- Git tags matching `v*.*.*` (e.g., v1.0.0)

**Process:**
1. Extract version from tag
2. Update version.go with build info
3. Run tests
4. Build optimized binaries
5. Create release packages with docs
6. Generate SHA256 checksums
7. Create GitHub Release with notes
8. Upload all assets

**Release Artifacts:**
- quickvm-windows-amd64.exe
- quickvm-windows-arm64.exe
- quickvm-v*.*.* -windows-amd64.zip (complete package)
- quickvm-v*.*.*-windows-arm64.zip (complete package)
- SHA256 checksums for all files

## How to Use

### Regular Development
```powershell
# Make changes and push
git add .
git commit -m "feat: add new feature"
git push

# GitHub Actions will automatically:
# - Build the project
# - Run tests
# - Run linter
# - Upload artifacts
```

### Creating a Release
```powershell
# Tag the version
git tag v1.0.1
git push origin v1.0.1

# GitHub Actions will automatically:
# - Build optimized binaries
# - Create release packages
# - Generate checksums
# - Create GitHub Release
# - Upload all assets
```

## Benefits

### For Users
- ✅ Organized documentation in `docs/` folder
- ✅ Easy to navigate with docs/README.md index
- ✅ Automated builds ensure quality
- ✅ Professional releases with checksums

### For Developers
- ✅ CI/CD catches issues early
- ✅ Automated testing on every push
- ✅ Code quality enforced with linting
- ✅ Easy release process (just tag)

### For the Project
- ✅ Professional development workflow
- ✅ Consistent build quality
- ✅ Better organization
- ✅ Automated documentation of changes (CHANGELOG)

## Testing the Workflows

The workflows will run automatically, but you can also:

1. **View workflow runs**: Go to Actions tab on GitHub
2. **Download artifacts**: Click on workflow run → Artifacts
3. **Check build status**: Badge on README shows current status

## Next Steps

1. ✅ Commit and push these changes
2. ✅ Workflows will run automatically
3. ✅ Create a v1.0.0 tag to test release workflow
4. ✅ Download and test the release artifacts

## Notes

- **Go Version**: Using Go 1.25.2 (latest)
- **Platforms**: Windows AMD64 and ARM64
- **Artifact Retention**: 30 days for regular builds, permanent for releases
- **Coverage**: Integrated with Codecov (optional)

---

**All set! Your project now has professional CI/CD! 🚀**
