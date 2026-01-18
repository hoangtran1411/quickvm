---
description: Safely refactor code with test coverage
---

# Refactoring Workflow

Follow these steps for safe refactoring:

## 1. Ensure Tests Exist

// turbo
Before any refactoring, verify current tests pass:
```powershell
go test -v ./...
```

If tests are missing for the code to refactor, write them first!

## 2. Check Current Coverage

// turbo
```powershell
go test -coverprofile=before_refactor.out ./...
go tool cover -func=before_refactor.out | grep "total:"
```

Note the coverage percentage.

## 3. Make Small Changes

Refactor in small increments:
- Extract function → test → commit
- Rename variable → test → commit
- Move code → test → commit

## 4. Run Tests After Each Change

// turbo
```powershell
go test -v ./...
```

## 5. Verify No Regression

// turbo
```powershell
go test -coverprofile=after_refactor.out ./...
go tool cover -func=after_refactor.out | grep "total:"
```

Coverage should be same or higher.

## Common Refactoring Patterns

### Extract Function
```go
// Before
func bigFunction() {
    // ... lots of code ...
    // Some logic here
    // ... more code ...
}

// After
func bigFunction() {
    // ... lots of code ...
    extractedLogic()
    // ... more code ...
}

func extractedLogic() {
    // Some logic here
}
```

### Extract Interface
```go
// Before: Concrete dependency
type Service struct {
    db *Database
}

// After: Interface dependency
type DataStore interface {
    Get(id string) (*Item, error)
    Save(item *Item) error
}

type Service struct {
    store DataStore
}
```

### Reduce Duplication
```go
// Before: Repeated validation
func StartVM(index int) error {
    vms, _ := m.GetVMs()
    if index < 1 || index > len(vms) {
        return fmt.Errorf("invalid index")
    }
    // ...
}

func StopVM(index int) error {
    vms, _ := m.GetVMs()
    if index < 1 || index > len(vms) {
        return fmt.Errorf("invalid index")
    }
    // ...
}

// After: Extracted validation
func (m *Manager) getVMByIndex(index int) (*VM, error) {
    vms, err := m.GetVMs()
    if err != nil {
        return nil, err
    }
    if index < 1 || index > len(vms) {
        return nil, fmt.Errorf("invalid VM index: %d", index)
    }
    return &vms[index-1], nil
}

func StartVM(index int) error {
    vm, err := m.getVMByIndex(index)
    if err != nil {
        return err
    }
    // ...
}
```

## 6. Format and Lint

// turbo
```powershell
go fmt ./...
go vet ./...
```

## 7. Final Verification

// turbo
```powershell
make check
```

## Checklist

- [ ] Existing tests pass before changes
- [ ] Coverage noted before refactoring
- [ ] Small incremental changes
- [ ] Tests run after each change
- [ ] No coverage regression
- [ ] Code formatted
- [ ] All checks pass
