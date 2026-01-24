---
trigger: always_on
---

# Go Style Guide - quickvm

> **Core Rules** - For full idioms reference, see `go-idioms-reference.md`

This project is a **Hyper-V Virtual Machine Management CLI/TUI** built with:
- **Cobra** for CLI command structure
- **BubbleTea & Lipgloss** for Terminal User Interface
- **Internal packages** (`hyperv/`) for business logic wrapping PowerShell
- **Windows-focused** architecture

---

## Code Style

- Format with `gofmt`/`goimports`. Run `golangci-lint run ./...` before commit.
- Adhere to [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
- **Organization**:
  - `cmd/`: entry points and CLI flag parsing.
  - `internal/`: private application logic (Hyper-V, TUI models).
  - `pkg/`: public libraries (if any).

## Error Handling

- **Wrap Errors**: Always use `fmt.Errorf("...: %w", err)` for context, especially for PowerShell/Hyper-V failures.
- **Fail Fast**: Use guard clauses to check errors immediately.
- **No Silent Failures**: Do not ignore errors in `defer` (log them if unavoidable).

## CLI & TUI Guidelines

- **Cobra (CLI)**:
  - Follow the "Command Pattern".
  - Subcommands stay in `cmd/`.
  - Logic stays in `internal/`.
- **BubbleTea (TUI)**:
  - Models stay in `internal/tui`.
  - Styling with `lipgloss` (define consistent themes).
  - Use `tea.Cmd` for side effects (like VM operations).

## Hyper-V & Systems Integration

- **PowerShell Security**:
  - **NEVER** concatenate user input into command arguments.
  - Use `exec.CommandContext` with individual arguments.
- **Context Awareness**:
  - All I/O functions MUST accept `context.Context` as the first argument.
  - Use timeouts for all external process calls.
- **Platform Specifics**:
  - Use `//go:build windows` for code that strictly depends on Windows APIs.
  - Mock interfaces to allow tests to run on non-Windows systems.

## Testing & Linting

- **Table-Driven Tests**: Primary method for `internal` logic.
- **Mocking**: Abstract Hyper-V calls behind interfaces to unit test logic without Admin rights.
- **Integration**: Real Hyper-V tests should be separated (e.g., via build tags or manual triggers).

---

## AI Agent Rules (Critical)

### Enforcement

- **Prefer Clarity**: Clear, verbose variable names over single letters (except idiomatic `i`, `ctx`, `err`).
- **Idiomatic Go**: Avoid bringing patterns from other languages (like getters/setters unless necessary for interfaces).

### Context Accuracy

- **Verify APIs**: Check `go.mod` for library versions (especially `bubbletea`).
- **Assume "No Access"**: When mocking, assume we don't have Admin privileges unless confirmed.

### Context Engineering

- Check `internal/hyperv` existing patterns before adding new wrappers.
- Do not blindly copy WMI queries; verify they work on standard Hyper-V installations.

---

## Quick Reference Links

- [Effective Go](https://go.dev/doc/effective_go)
- [Cobra](https://github.com/spf13/cobra)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Go Hyper-V Reference](https://github.com/sheepla/go-hyperv)
- [golangci-lint](https://github.com/golangci/golangci-lint)