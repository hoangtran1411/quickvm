# Lint Fixes - 2026-01-05

This document describes the linting errors that were fixed to make the CI/CD pipeline pass.

## Fixed Issues

### 1. Ineffectual Assignment in `ui/table.go` (Line 168)

**Error:**
```
ui\table.go:168:3: ineffectual assignment to state (ineffassign)
    state := vm.State
    ^
```

**Problem:**
The variable `state` was being assigned `vm.State` but was immediately reassigned in the switch statement below, making the initial assignment useless.

**Fix:**
Changed from:
```go
state := vm.State
switch strings.ToLower(vm.State) {
case "running":
    state = statusRunningStyle.Render(vm.State)
// ...
}
```

To:
```go
var state string
switch strings.ToLower(vm.State) {
case "running":
    state = statusRunningStyle.Render(vm.State)
// ...
}
```

### 2. Potential Nil Pointer Dereference in `updater/updater_test.go` (Line 14)

**Error:**
```
updater\updater_test.go:14:7: SA5011: possible nil pointer dereference (staticcheck)
    if u.currentVersion != "1.0.0" {
         ^
updater\updater_test.go:10:5: SA5011(related information): this check suggests that the pointer can be nil (staticcheck)
    if u == nil {
       ^
```

**Problem:**
The test checks if `u == nil` on line 10 but doesn't return early, meaning the code could continue to line 14 and dereference the nil pointer.

**Fix:**
Added early return after nil check:
```go
func TestNewUpdater(t *testing.T) {
    u := NewUpdater("1.0.0")
    if u == nil {
        t.Error("NewUpdater() returned nil")
        return  // <- Added this line
    }
    
    if u.currentVersion != "1.0.0" {
        t.Errorf("Expected version 1.0.0, got %s", u.currentVersion)
    }
}
```

## Summary

Both issues have been resolved:
- ✅ Removed ineffectual assignment in `ui/table.go`
- ✅ Added early return to prevent nil pointer dereference in `updater/updater_test.go`

The CI/CD lint job should now pass successfully.
