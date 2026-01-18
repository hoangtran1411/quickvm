---
description: Standard development cycle for making changes
---

# Development Cycle

// turbo-all

Standard workflow for implementing changes:

## 1. Format Code

```powershell
go fmt ./...
```

## 2. Run Linter

```powershell
golangci-lint run
```

Or if not installed:
```powershell
go vet ./...
```

## 3. Run Tests

```powershell
go test -v ./...
```

## 4. Build

```powershell
go build -o quickvm.exe
```

## 5. Manual Test (if needed)

```powershell
./quickvm.exe --help
./quickvm.exe list
./quickvm.exe version
```

## Quick One-liner

For fast iteration:
```powershell
go fmt ./... && go test ./... && go build -o quickvm.exe
```

Or use Makefile:
```powershell
make dev
```
