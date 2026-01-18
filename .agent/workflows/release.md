---
description: Create a new release version
---

# Release Workflow

Follow these steps to create a new release:

## 1. Update Version

Edit `cmd/version.go` and update the version:
```go
var (
    Version   = "1.x.x"  // Update this
    BuildDate = "2026-01-xx"
    GitCommit = "xxxxxxx"
)
```

## 2. Update CHANGELOG.md

Add new section at the top:
```markdown
## [1.x.x] - 2026-01-xx

### Added
- New feature description

### Changed
- Change description

### Fixed
- Bug fix description
```

## 3. Run All Checks

// turbo
```powershell
make check
```

Or manually:
```powershell
go fmt ./...
golangci-lint run
go test -v ./...
```

## 4. Build Release Binaries

// turbo
```powershell
make build-all
```

This creates:
- `build/quickvm.exe` (AMD64)
- `build/quickvm-arm64.exe` (ARM64)

## 5. Commit Changes

```powershell
git add -A
git commit -m "Release v1.x.x"
```

## 6. Create Tag

```powershell
git tag -a v1.x.x -m "Release v1.x.x"
```

## 7. Push to GitHub

```powershell
git push origin main
git push origin v1.x.x
```

## 8. Verify GitHub Release

GitHub Actions will automatically:
- Build binaries for all architectures
- Create checksums
- Publish release with artifacts

Check: https://github.com/hoangtran1411/quickvm/releases

## Checklist

- [ ] Version updated in `cmd/version.go`
- [ ] CHANGELOG.md updated
- [ ] All tests pass
- [ ] Code formatted and linted
- [ ] Tag created and pushed
- [ ] GitHub Release created automatically
