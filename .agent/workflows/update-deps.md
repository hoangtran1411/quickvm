---
description: Update Go dependencies safely
---

# Update Dependencies Workflow

// turbo-all

Follow these steps to update dependencies:

## 1. Check Current Dependencies

```powershell
go list -m all
```

## 2. Check for Updates

```powershell
go list -m -u all
```

## 3. Update All Dependencies

```powershell
go get -u ./...
go mod tidy
```

## 4. Or Update Specific Package

```powershell
go get -u github.com/spf13/cobra@latest
go get -u github.com/charmbracelet/bubbletea@latest
go mod tidy
```

## 5. Run Tests

```powershell
go test -v ./...
```

## 6. Build and Verify

```powershell
go build -o quickvm.exe
./quickvm.exe version
./quickvm.exe list
```

## 7. Check for Breaking Changes

If tests fail after update:
1. Check the package's CHANGELOG
2. Update code to match new API
3. Re-run tests

## Common Dependencies & Update Notes

| Package | Notes |
|---------|-------|
| `cobra` | Usually backward compatible |
| `bubbletea` | Check for API changes |
| `lipgloss` | Check style API changes |
| `bubbles` | Check table API changes |

## Rollback if Needed

```powershell
git checkout go.mod go.sum
go mod download
```

## Checklist

- [ ] Current deps listed
- [ ] Updates checked
- [ ] Dependencies updated
- [ ] go.mod tidied
- [ ] All tests pass
- [ ] Manual verification done
- [ ] Committed changes
