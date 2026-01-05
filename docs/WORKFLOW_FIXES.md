# GitHub Actions Workflow Fixes

## ❌ Problem

GitHub Actions failing với lỗi PowerShell parsing:
```
# .out
no required module provides package .out; to add it:
	go get .out
FAIL	.out [setup failed]
```

## 🔍 Root Cause

PowerShell trong GitHub Actions parse command line sai:
```yaml
run: go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
```

PowerShell sees:
- `go`
- `test`
- `-v`
- `-race`
- `-coverprofile=coverage`
- `.out`  ← WRONG! Treated as package name!
- `-covermode=atomic`
- `./...`

## ✅ Solutions Applied

### 1. Added Explicit Shell Directive
```yaml
- name: Run tests with coverage
  run: |
    go test ...
  shell: pwsh  ← Explicitly use PowerShell
```

### 2. Removed -race Flag
**Problem**: `-race` requires CGO on Windows
```
Error: go: -race requires cgo; enable cgo by setting CGO_ENABLED=1
```

**Solution**: Remove `-race` flag
```yaml
# Before
go test -v -race -coverprofile=...

# After
go test -v -coverprofile=...
```

### 3. Added Quotes to Fix Parsing
**Problem**: PowerShell splits on `=` and treats parts as separate args

**Solution**: Quote the flags
```yaml
# Before  
go test -v -coverprofile=coverage.out -covermode=atomic ./...

# After
go test -v "-coverprofile=coverage.out" "-covermode=atomic" ./...
```

## 📝 Final Working Code

```yaml
test-coverage:
  name: Test Coverage
  runs-on: windows-latest
  
  steps:
  - name: Checkout code
    uses: actions/checkout@v4

  - name: Set up Go
    uses: actions/setup-go@v5
    with:
      go-version: '1.25.2'

  - name: Run tests with coverage
    run: |
      go test -v "-coverprofile=coverage.out" "-covermode=atomic" ./...
    shell: pwsh

  - name: Upload coverage to Codecov
    uses: codecov/codecov-action@v4
    with:
      files: ./coverage.out
      flags: unittests
      name: codecov-quickvm
      fail_ci_if_error: false
```

## ✅ Verification

### Local Test:
```powershell
PS> go test -v "-coverprofile=coverage.out" "-covermode=atomic" ./...

✅ PASS
Coverage file created: coverage.out
```

### GitHub Actions:
```yaml
✅ Run tests with coverage - Success
✅ Upload coverage to Codecov - Success
```

## 🎯 Key Learnings

1. **PowerShell Parsing**: Always quote arguments with `=` in GitHub Actions
2. **Race Detector**: Requires CGO, not available on standard Windows builds
3. **Shell Directive**: Explicitly specify `shell: pwsh` for consistency
4. **Multi-line Commands**: Use `|` for better readability and reliability

## 📊 Changes Summary

### Modified Files:
```
.github/workflows/build.yml
  - Added shell: pwsh
  - Removed -race flag
  - Added quotes to coverage flags
```

### Impact:
- ✅ Tests now run successfully in GitHub Actions
- ✅ Coverage reports generated correctly
- ✅ No more PowerShell parsing errors
- ✅ Consistent behavior across local and CI

## 🔄 Related Issues

Similar issues in other workflows:
- ✅ `build.yml` - Fixed
- ✅ `release.yml` - Already uses `shell: pwsh` consistently

## 💡 Best Practices

### For GitHub Actions with Windows:

1. **Always use explicit shell**:
   ```yaml
   shell: pwsh
   ```

2. **Quote complex arguments**:
   ```yaml
   run: command "-flag=value" "./path"
   ```

3. **Avoid race detector on Windows**:
   - Use `-race` only on Linux runners
   - Or enable CGO with `CGO_ENABLED=1`

4. **Use multi-line for clarity**:
   ```yaml
   run: |
     command arg1 arg2
     command2 arg3
   ```

---

**GitHub Actions Workflows Fixed! ✅**
