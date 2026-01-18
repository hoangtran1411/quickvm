---
name: Go Testing Patterns
description: Comprehensive testing patterns for Go applications including table-driven tests, mocking interfaces, build tags, and CI integration.
---

# Go Testing Patterns

This skill provides professional testing patterns for Go applications, focusing on maintainability, coverage, and CI integration.

## When to Use

- Writing unit tests for Go packages
- Mocking external dependencies (databases, APIs, shells)
- Setting up CI/CD pipelines with test coverage
- Creating integration tests with build tags

## Core Patterns

### 1. Table-Driven Tests

The preferred pattern for testing multiple scenarios:

```go
func TestValidateInput(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    bool
        wantErr bool
    }{
        {
            name:  "valid input",
            input: "hello",
            want:  true,
        },
        {
            name:    "empty input",
            input:   "",
            wantErr: true,
        },
        {
            name:  "special characters",
            input: "hello-world_123",
            want:  true,
        },
        {
            name:    "invalid characters",
            input:   "hello;drop",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ValidateInput(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateInput() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if got != tt.want {
                t.Errorf("ValidateInput() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 2. Interface-Based Mocking

Define interfaces for external dependencies to enable mocking:

```go
// In production code
type DataStore interface {
    Get(id string) (*Item, error)
    Save(item *Item) error
}

type Service struct {
    store DataStore
}

// In test code
type MockDataStore struct {
    GetFunc  func(id string) (*Item, error)
    SaveFunc func(item *Item) error
}

func (m *MockDataStore) Get(id string) (*Item, error) {
    if m.GetFunc != nil {
        return m.GetFunc(id)
    }
    return nil, nil
}

func (m *MockDataStore) Save(item *Item) error {
    if m.SaveFunc != nil {
        return m.SaveFunc(item)
    }
    return nil
}

// Usage in tests
func TestService_GetItem(t *testing.T) {
    mock := &MockDataStore{
        GetFunc: func(id string) (*Item, error) {
            return &Item{ID: id, Name: "Test"}, nil
        },
    }

    service := &Service{store: mock}
    item, err := service.GetItem("123")

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if item.Name != "Test" {
        t.Errorf("expected Name 'Test', got '%s'", item.Name)
    }
}
```

### 3. Test Helpers

Create reusable test helpers:

```go
// testhelpers_test.go

// Helper to create test items
func makeTestItem(t *testing.T, name string) *Item {
    t.Helper()
    return &Item{
        ID:        uuid.New().String(),
        Name:      name,
        CreatedAt: time.Now(),
    }
}

// Helper for comparing errors
func assertError(t *testing.T, got, want error) {
    t.Helper()
    if got == nil && want != nil {
        t.Errorf("expected error %v, got nil", want)
        return
    }
    if got != nil && want == nil {
        t.Errorf("unexpected error: %v", got)
        return
    }
    if got != nil && want != nil && !errors.Is(got, want) {
        t.Errorf("expected error %v, got %v", want, got)
    }
}

// Helper for cleanup
func setupTestEnv(t *testing.T) func() {
    t.Helper()
    // Setup...
    return func() {
        // Cleanup...
    }
}
```

### 4. Build Tags for Integration Tests

Separate unit tests from integration tests:

```go
//go:build integration
// +build integration

package mypackage

import (
    "testing"
)

func TestIntegration_RealDatabase(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    // Real database tests here
}
```

Run with: `go test -tags=integration ./...`

### 5. Windows-Specific Tests

For Windows-only functionality:

```go
//go:build windows
// +build windows

package hyperv

import (
    "testing"
)

func TestHyperV_ListVMs(t *testing.T) {
    // This test only runs on Windows
}
```

### 6. Test Main for Setup/Teardown

```go
func TestMain(m *testing.M) {
    // Setup before all tests
    setup()
    
    // Run all tests
    code := m.Run()
    
    // Teardown after all tests
    teardown()
    
    os.Exit(code)
}

func setup() {
    // Initialize resources
}

func teardown() {
    // Clean up resources
}
```

## Coverage Commands

```bash
# Run tests with coverage
go test -v -cover ./...

# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Get coverage percentage
go test -cover ./... | grep coverage

# Coverage for specific package with detailed output
go test -v -coverprofile=coverage.out -covermode=atomic ./mypackage/...
```

## Makefile Integration

```makefile
.PHONY: test test-coverage test-integration

test: ## Run unit tests
	go test -v ./...

test-coverage: ## Run tests with coverage report
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-integration: ## Run integration tests
	go test -v -tags=integration ./...

test-short: ## Run quick tests only
	go test -v -short ./...

bench: ## Run benchmarks
	go test -bench=. -benchmem ./...
```

## GitHub Actions CI

```yaml
name: Tests

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: windows-latest  # or ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'

    - name: Cache Go modules
      uses: actions/cache@v4
      with:
        path: |
          ~/go/pkg/mod
          ~/.cache/go-build
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-

    - name: Run tests
      run: go test -v ./...

    - name: Run tests with coverage
      run: go test -v -coverprofile=coverage.out -covermode=atomic ./...

    - name: Upload coverage
      uses: codecov/codecov-action@v5
      with:
        files: ./coverage.out
        fail_ci_if_error: false
```

## Best Practices

1. **Name tests descriptively**: `TestFunction_Scenario_ExpectedResult`

2. **Use `t.Helper()`** in helper functions for better error reporting

3. **Avoid test pollution**: Each test should be independent

4. **Use `t.Parallel()`** for parallel execution when safe:
   ```go
   func TestParallel(t *testing.T) {
       t.Parallel()
       // Test code
   }
   ```

5. **Keep tests fast**: Unit tests should run in milliseconds

6. **Test error paths**: Don't just test happy paths

7. **Use meaningful assertions**:
   ```go
   // ❌ Bad
   if result != expected {
       t.Fail()
   }
   
   // ✅ Good
   if result != expected {
       t.Errorf("got %v, want %v", result, expected)
   }
   ```

8. **Test behavior, not implementation**: Focus on what the function does, not how

## Common Patterns for Mocking

### Mock with Call Tracking

```go
type MockService struct {
    Calls     []string
    ReturnVal interface{}
    ReturnErr error
}

func (m *MockService) DoSomething(id string) error {
    m.Calls = append(m.Calls, "DoSomething:"+id)
    return m.ReturnErr
}

// In test
func TestTrackCalls(t *testing.T) {
    mock := &MockService{}
    // ... run code ...
    
    if len(mock.Calls) != 2 {
        t.Errorf("expected 2 calls, got %d", len(mock.Calls))
    }
}
```

### Mock with Sequence of Returns

```go
type SequenceMock struct {
    Returns []interface{}
    Index   int
}

func (m *SequenceMock) Next() interface{} {
    if m.Index >= len(m.Returns) {
        return nil
    }
    val := m.Returns[m.Index]
    m.Index++
    return val
}
```
