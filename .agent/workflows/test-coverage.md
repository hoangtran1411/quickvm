---
description: Generate and analyze test coverage reports
---

# Test Coverage Workflow

// turbo-all

Generate and analyze test coverage:

## 1. Run Tests with Coverage

```powershell
go test -v -cover ./...
```

## 2. Generate Coverage Profile

```powershell
go test -coverprofile=coverage.out ./...
```

## 3. View Coverage Summary

```powershell
go tool cover -func=coverage.out
```

## 4. Generate HTML Report

```powershell
go tool cover -html=coverage.out -o coverage.html
```

Then open `coverage.html` in browser.

## 5. Check Specific Package

```powershell
go test -v -coverprofile=hyperv_coverage -cover ./hyperv/...
go tool cover -func=hyperv_coverage
```

## Coverage Goals

| Package | Target | Current |
|---------|--------|---------|
| `hyperv/` | 70%+ | Check report |
| `cmd/` | 50%+ | Check report |
| `ui/` | Manual testing | N/A |

## Improving Coverage

### Find Uncovered Code
Open `coverage.html` - red lines are uncovered.

### Common Gaps
1. Error handling paths
2. Edge cases (empty input, invalid index)
3. JSON parsing variations

### Add Tests for Gaps
```go
func TestFunction_ErrorCase(t *testing.T) {
    mock := NewMockExecutor(nil, errors.New("simulated error"))
    manager := NewManagerWithExecutor(mock)
    
    _, err := manager.Function()
    
    if err == nil {
        t.Fatal("expected error")
    }
}
```

## Quick Commands

```powershell
# Full coverage workflow
make test-coverage

# Or manually
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html
```
