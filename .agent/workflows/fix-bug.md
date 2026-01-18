---
description: Debug and fix a bug in the project
---

# Bug Fix Workflow

Follow these steps to debug and fix bugs:

## 1. Reproduce the Bug

First, reproduce the issue:
```powershell
./quickvm.exe <command that causes bug>
```

Document:
- Exact command/input
- Expected behavior
- Actual behavior
- Error message (if any)

## 2. Locate the Issue

### Check error messages
Error messages often point to the file and function.

### Search for relevant code
```powershell
# Search for function/variable
grep -r "functionName" --include="*.go"

# Search for error message
grep -r "error text" --include="*.go"
```

### Common locations:
- CLI handling: `cmd/*.go`
- Business logic: `hyperv/*.go`
- TUI issues: `ui/*.go`

## 3. Write a Failing Test

// turbo
Create a test that reproduces the bug:

```go
func TestBugFix_IssueDescription(t *testing.T) {
    // Setup that triggers the bug
    // ...
    
    // This should fail before the fix
    result, err := functionThatHasBug()
    
    if err != nil {
        t.Fatalf("expected no error, got: %v", err)
    }
    
    expected := "expected value"
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

## 4. Fix the Code

Make the minimal change to fix the issue.

## 5. Verify the Fix

// turbo
Run the new test:
```powershell
go test -v -run "TestBugFix_IssueDescription" ./...
```

// turbo
Run all tests to ensure no regression:
```powershell
go test -v ./...
```

## 6. Manual Verification

// turbo
Build and test manually:
```powershell
go build -o quickvm.exe
./quickvm.exe <command that had bug>
```

## 7. Document the Fix

If significant:
- Update CHANGELOG.md in the `### Fixed` section
- Add comments explaining the fix if not obvious

## Debugging Tips

### Add Temporary Debug Output
```go
fmt.Printf("[DEBUG] value = %v\n", value)
```

### Check PowerShell Output
```go
output, err := m.Exec.RunCommand(script)
fmt.Printf("[DEBUG] PS Output: %s\n", string(output))
fmt.Printf("[DEBUG] PS Error: %v\n", err)
```

### Run Specific Package
```powershell
go test -v ./hyperv/...
go test -v ./cmd/...
```

### Run Specific Test
```powershell
go test -v -run "TestName" ./package/...
```

## Checklist

- [ ] Bug reproduced
- [ ] Root cause identified
- [ ] Failing test written
- [ ] Fix implemented
- [ ] New test passes
- [ ] All existing tests pass
- [ ] Manual verification done
- [ ] CHANGELOG updated (if needed)
