# Contributing to QuickVM

> How to contribute to QuickVM - from bug reports to pull requests.

**Last Updated:** 2026-01-12

---

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Quick Start for Contributors](#-quick-start-for-contributors)
- [Reporting Bugs](#-reporting-bugs)
- [Suggesting Features](#-suggesting-features)
- [Development Setup](#-development-setup)
- [Development Workflow](#-development-workflow)
- [Coding Guidelines](#-coding-guidelines)
- [Pull Request Process](#-pull-request-process)

---

## Code of Conduct

This project is governed by respect and professionalism. Be kind and constructive.

---

## 🚀 Quick Start for Contributors

```powershell
# 1. Fork & clone
git clone https://github.com/YOUR-USERNAME/quickvm.git
cd quickvm

# 2. Install dependencies
go mod download

# 3. Create feature branch
git checkout -b feature/my-feature

# 4. Make changes & test
go build -o quickvm.exe
go test ./...
go fmt ./...

# 5. Commit & push
git commit -m "feat: add my feature"
git push origin feature/my-feature

# 6. Create Pull Request on GitHub
```

---

## 🐛 Reporting Bugs

Before creating a bug report, check existing issues first.

### Bug Report Template

```markdown
**Environment:**
- OS: Windows 10/11
- QuickVM Version: x.x.x
- Go Version: x.x.x

**Description:**
A clear description of the bug.

**Steps to Reproduce:**
1. Step 1
2. Step 2
3. ...

**Expected Behavior:**
What you expected.

**Actual Behavior:**
What actually happened.

**Screenshots:**
If applicable.
```

---

## 💡 Suggesting Features

> ⚠️ **Important**: Check [FEATURE_ROADMAP.md](FEATURE_ROADMAP.md) first!
> 
> Features in **Tier 3 & 4 are ARCHIVED** and won't be implemented.
> QuickVM is intentionally a "Quick" VM manager, not a Hyper-V replacement.

When suggesting features:
- Explain the use case
- Describe how it fits the "Quick" philosophy
- Provide examples of usage

---

## 🛠️ Development Setup

### Prerequisites

- **Go 1.21+**: [Download Go](https://golang.org/dl/)
- **Windows 10/11** with Hyper-V enabled
- **Git**: [Download Git](https://git-scm.com/)
- **(Optional)** golangci-lint for linting

### First Time Setup

```powershell
# Clone your fork
git clone https://github.com/YOUR-USERNAME/quickvm.git
cd quickvm

# Install dependencies
go mod download
go mod tidy

# Build
go build -o quickvm.exe

# Test
go test ./...
```

---

## 📝 Development Workflow

### Build Commands

```powershell
# Normal build
go build -o quickvm.exe

# Optimized build (smaller binary)
go build -ldflags="-s -w" -o quickvm.exe

# Cross-compile
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o quickvm-amd64.exe
$env:GOOS="windows"; $env:GOARCH="arm64"; go build -o quickvm-arm64.exe
```

### Test Commands

```powershell
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Specific package
go test ./hyperv
```

### Code Quality

```powershell
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Verify dependencies
go mod verify
```

### Manual Testing Checklist

- [ ] `quickvm list`
- [ ] `quickvm start 1`
- [ ] `quickvm stop 1`
- [ ] `quickvm` (TUI mode)
- [ ] TUI navigation works
- [ ] `quickvm version`

---

## 📏 Coding Guidelines

### Go Style

Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments):

1. **Use `gofmt`**: Always format code
2. **Handle errors**: Never ignore errors
3. **Comment exports**: Document all exported functions/types
4. **Keep simple**: Readable code over clever code

### Error Handling

```go
if err := vm.Start(); err != nil {
    return fmt.Errorf("failed to start vm: %w", err)  // Wrap with %w
}
```

### PowerShell Scripts

1. Test scripts independently first
2. Use explicit type conversions: `[int]`, `[string]`
3. Handle errors with clear messages
4. Keep scripts simple

### ⚠️ Linting Rules (IMPORTANT)

This project uses `golangci-lint` in CI. **All PRs must pass linting.**

#### Common Lint Errors to Avoid

| Rule | Error | Fix |
|------|-------|-----|
| `errcheck` | Error return value not checked | Always check or explicitly ignore with `_ =` |
| `ineffassign` | Ineffectual assignment | Remove unused assignments |
| `staticcheck` | Various static analysis | Follow suggestions |

#### Error Return Values (`errcheck`)

**❌ Wrong** - Will fail CI:
```go
os.WriteFile(path, data, 0644)
w.Write([]byte(response))
copyFile(src, dst)
```

**✅ Correct** - Production code:
```go
if err := os.WriteFile(path, data, 0644); err != nil {
    return fmt.Errorf("failed to write file: %w", err)
}
```

**✅ Correct** - Test/Benchmark code (when error handling not needed):
```go
_ = os.WriteFile(path, data, 0644)
_, _ = w.Write([]byte(response))
_ = copyFile(src, dst)
```

#### Run Lint Locally Before Push

```powershell
# Install golangci-lint (one time)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run lint
golangci-lint run ./...

# Run with timeout (for large codebases)
golangci-lint run --timeout=5m ./...
```

#### If You Don't Have golangci-lint

At minimum, run these before pushing:
```powershell
go vet ./...      # Basic static analysis
go build ./...    # Ensure code compiles
go test ./...     # Ensure tests pass
```

### Commit Messages

Use conventional commits:

```
feat: add support for VM snapshots
fix: correct memory calculation
docs: update README
test: add tests for hyperv package
refactor: simplify PowerShell generation
```

Prefixes:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `test:` - Tests
- `refactor:` - Code refactoring
- `chore:` - Maintenance

---

## 🔄 Pull Request Process

### 1. Update your fork

```powershell
git remote add upstream https://github.com/hoangtran1411/quickvm.git
git fetch upstream
git checkout main
git merge upstream/main
```

### 2. Create feature branch

```powershell
git checkout -b feature/my-feature
```

### 3. Make changes

- Write code
- Add tests
- Update docs

### 4. Test thoroughly

```powershell
go test ./...
go build -o quickvm.exe
# Manual testing
```

### 5. Commit & push

```powershell
git add .
git commit -m "feat: add my feature"
git push origin feature/my-feature
```

### 6. Create Pull Request

Go to GitHub and create a PR with clear description.

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe testing performed

## Checklist
- [ ] Code follows project style
- [ ] Self-review performed
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] All tests pass
```

---

## 🏗️ Adding New Features

### Adding a New Command

1. Create `cmd/mycommand.go`:

```go
package cmd

import "github.com/spf13/cobra"

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Short description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

2. Add tests in `cmd/mycommand_test.go`
3. Update documentation

### Adding Hyper-V Functionality

1. Add method to `hyperv/hyperv.go`
2. Use `m.Exec.RunCommand()` for PowerShell
3. Add tests with mock executor
4. Handle errors properly

---

## 🐞 Troubleshooting Development

### Import issues
```powershell
go clean -modcache
go mod download
go mod tidy
```

### Build issues
```powershell
go clean -cache
go build -v -o quickvm.exe
```

### PowerShell execution policy
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

---

## 📞 Questions?

- Open an issue for questions
- Join discussions in existing issues
- Check existing documentation

---

## 🌟 Recognition

Contributors are recognized in:
- GitHub contributors list
- Release notes
- Project README (significant contributions)

---

**Thank you for contributing to QuickVM! 🚀**
