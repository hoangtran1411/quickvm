---
description: Review code changes before merging
---

# Code Review Workflow

Checklist for reviewing code changes:

## 1. Understand the Change

- What problem does this solve?
- What is the expected behavior?

## 2. Code Quality Checks

### Formatting
// turbo
```powershell
go fmt ./...
git diff  # Should be empty if properly formatted
```

### Linting
// turbo
```powershell
golangci-lint run
```

### Tests
// turbo
```powershell
go test -v ./...
```

## 3. Review Checklist

### Code Style
- [ ] Follows Go conventions (gofmt, effective go)
- [ ] Clear variable/function names
- [ ] No unused code/imports
- [ ] Proper error handling with `%w`

### Architecture
- [ ] Code in correct package (`cmd/`, `hyperv/`, `ui/`)
- [ ] No circular dependencies
- [ ] Uses interfaces for testability

### Security (for PowerShell integration)
- [ ] User input properly escaped
- [ ] No command injection vulnerabilities
- [ ] Uses context for timeouts

### Testing
- [ ] New code has tests
- [ ] Edge cases covered
- [ ] Mocks used appropriately

### Documentation
- [ ] Exported functions have comments
- [ ] README updated if needed
- [ ] CHANGELOG updated for user-facing changes

## 4. Manual Testing

// turbo
```powershell
go build -o quickvm.exe
./quickvm.exe <affected commands>
```

## 5. Final Verification

// turbo
```powershell
make check
```

## Review Rating Guide

- **Approve**: All checks pass, no issues
- **Request Changes**: Has issues that must be fixed
- **Comment**: Minor suggestions, can merge anyway
