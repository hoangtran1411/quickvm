# AGENTS.md - QuickVM

QuickVM is a fast Hyper-V Virtual Machine management CLI/TUI for Windows, built with Go.

## Quick Commands

```powershell
# Install dependencies
go mod download

# Build
go build -o quickvm.exe

# Test (runs on any OS via mocked interfaces)
go test ./...

# Lint (MUST use golangci-lint v2.8.0+)
golangci-lint run ./...

# Format
gofmt -w .
goimports -w .

# Full dev cycle
make dev
```

## Project Structure

```
quickvm/
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go            # Entry point, launches TUI if no args
│   ├── start.go           # quickvm start <index>
│   ├── stop.go            # quickvm stop <index>
│   └── ...                # Other commands
├── internal/
│   └── hyperv/            # Hyper-V business logic
│       ├── hyperv.go      # Manager struct, VM operations
│       ├── executor.go    # ShellExecutor interface
│       └── *_test.go      # Unit tests with mocks
├── ui/                    # TUI (Bubble Tea + Lipgloss)
├── .agent/                # AI agent configuration
│   ├── rules/             # Always-on coding rules
│   ├── skills/            # Reusable code patterns
│   └── workflows/         # Step-by-step automation
└── docs/                  # Documentation
```

## Code Style

- **Linting**: Use `golangci-lint` v2.8.0+ with v2 schema configuration
- **Formatting**: Always run `gofmt` and `goimports` before committing
- **Error Handling**: Always wrap errors with context using `fmt.Errorf("...: %w", err)`
- **Context**: All I/O functions MUST accept `context.Context` as first argument
- **Naming**: Avoid stuttering (e.g., `hyperv.Manager` not `hyperv.HyperVManager`)

## PowerShell Security (CRITICAL)

**NEVER** concatenate user input into PowerShell commands:

```go
// ❌ DANGEROUS - Command injection vulnerability
exec.Command("powershell", "-Command", "Get-VM -Name " + userInput)

// ✅ SAFE - Use separate arguments
exec.CommandContext(ctx, "powershell", "-Command", "Get-VM", "-Name", userInput)
```

## Testing

- **Table-driven tests**: Primary pattern for all unit tests
- **Mocking**: Use `ShellExecutor` interface to mock PowerShell calls
- **Build tags**: Use `//go:build windows` for integration tests requiring real Hyper-V
- **Coverage**: Run `go test -coverprofile=coverage.out ./...`

```go
// Example test pattern
func TestStartVM(t *testing.T) {
    tests := []struct {
        name    string
        input   int
        wantErr bool
    }{
        {"valid index", 1, false},
        {"invalid index", 0, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Key Interfaces

```go
// ShellExecutor - Abstract PowerShell execution for testing
type ShellExecutor interface {
    RunCommand(ctx context.Context, script string) ([]byte, error)
}

// Manager - Main Hyper-V operations
type Manager struct {
    Exec ShellExecutor
}
```

## Common Tasks

| Task | Command/Location |
|------|------------------|
| Add CLI command | Create `cmd/<command>.go`, follow Cobra pattern |
| Add Hyper-V operation | Add method to `internal/hyperv/hyperv.go` |
| Add TUI feature | Modify `ui/table.go`, use Bubble Tea patterns |
| Run tests | `go test ./...` or `make test` |
| Build release | `make build` or `go build -ldflags="-s -w"` |

## AI Workflows

Use slash commands for automated tasks:
- `/dev-cycle` - Format, lint, test, build
- `/add-command` - Add new CLI command
- `/add-hyperv-feature` - Add Hyper-V functionality
- `/fix-bug` - Debug and fix issues

## Documentation

- [Getting Started](docs/GETTING_STARTED.md)
- [Developer Guide](docs/DEVELOPER.md)
- [AI Agent Setup](docs/AI_AGENT.md)
- [Contributing](docs/CONTRIBUTING.md)

## Pull Request Guidelines

1. Run `golangci-lint run ./...` - must pass
2. Run `go test ./...` - must pass
3. Update documentation if API changes
4. Follow [Conventional Commits](https://www.conventionalcommits.org/)
5. Keep PRs focused on single concerns
