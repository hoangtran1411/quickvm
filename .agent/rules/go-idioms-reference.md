---
trigger: always_on
---

# Go Idioms Reference - quickvm

> **Full Reference Document** - Contains detailed idioms, code examples, and best practices.
> For compact core rules, see `go-style-guide.md`

---

## Naming Conventions (Idiomatic Go)

- Use short but meaningful names, scoped by context:
  - `r`, `w`, `ctx`, `db`, `tx`, `cfg` are acceptable in small scopes.
  - Avoid `data`, `info`, `obj`, `temp`, `value` unless unavoidable.

- Prefer noun-based names for structs, verb-based names for functions:
  - `vm.Start()`, `manager.GetVM()`, `parser.Parse()`

- Boolean names should read naturally:
  - `isRunning`, `hasNetwork`, `enableLogging`

- Avoid stuttering:
  - ❌ `hyperv.HyperVManager`
  - ✅ `hyperv.Manager`

---

## Function Design

- Prefer small functions (≤ 40 lines).
- One function = one responsibility.
- Avoid flags that change behavior dramatically:

```go
// bad
func StartVM(name string, force bool)

// good
func StartVM(name string)
func ForceStartVM(name string)
```

- Return early (guard clauses):

```go
if err != nil {
    return nil, err
}
```

---

## Error Handling Idioms

- **Traceability**: Always wrap errors at package boundaries to maintain the call stack context.

```go
return nil, fmt.Errorf("failed to start VM: %w", err)
```

- **Fail Fast**: Check errors immediately.
- **Typed Errors**: Use typed errors for specific Hyper-V failure states if callers need to handle them differently (e.g., `ErrVMNotFound`).

```go
if errors.Is(err, hyperv.ErrVMNotFound) { ... }
```

---

## Package Design & Boundaries

- `internal` packages must be:
  - **Logic-Focused**: `internal/hyperv` should know nothing about Cobra or BubbleTea.
  - **Decoupled**: Use interfaces to interact between layers.

- `cmd` package (Cobra):
  - handles flags, arguments, and UI output.
  - calls `internal` logic.

- Avoid circular dependencies.

---

## Context Usage Idioms

- **Mandatory**: Every function performing I/O or calling external commands (PowerShell) MUST accept `context.Context` as its first argument.

```go
func (m *Manager) GetVM(ctx context.Context, name string) (*VM, error)
```

- Do not pass nil context in production code. Use `context.TODO()` if unsure during dev, but `context.Background()` only at entry points.
- **Timeouts**: Use `context.WithTimeout` for all Hyper-V operations to prevent hanging processes.

---

## Shell Execution & Security (PowerShell)

- **Sanitization**: NEVER construct PowerShell commands using string concatenation with user input.

```go
// ❌ Dangerous
exec.Command("powershell", "-Command", "Get-VM -Name " + input)

// ✅ Safe
exec.CommandContext(ctx, "powershell", "-Command", "Get-VM", "-Name", input)
```

- **Execution**: Prefer `exec.CommandContext` over `exec.Command`.

---

## Concurrency Patterns

- **Worker Pools**: Use for batch operations (e.g., starting multiple VMs).
- **Ownership**: Clearly define who starts and stops goroutines.
- **ErrGroup**: Use `errgroup.Group` for running multiple tasks that might fail.

---

## Struct & Interface Idioms

- Accept interfaces, return concrete types.
- **Mocking**: Define interfaces for Hyper-V or OS interactions to allow testing without actual Windows calls.

```go
type VMManager interface {
    Start(ctx context.Context, name string) error
    Stop(ctx context.Context, name string) error
}
```

---

## Zero Value Philosophy

- Design structs so zero value is usable where possible.
- Prefer empty slices over nil slices for JSON/List output.

---

## Testing Idioms

- **Table-Driven**: Primary pattern for logic tests.
- **Build Tags**: Use `//go:build windows` for integration tests that require real Hyper-V.
- **Mocking**: Use mocks for unit tests in `internal/hyperv` to run on any OS.
- **Linting**: `.golangci.yml` MUST use version 2 schema (`version: "2"`).
  - Use kebab-case for linter settings (e.g., `ignore-sigs`, `ignore-package-globs`).
  - Exclusions must be configured under `linters: exclusions: rules` instead of `issues: exclude-rules`.
  - Prefer global exclusions in config over redundant `//nolint` comments in test files.

---

## Comments & Documentation

- Comments explain **why**, not **what**.
- Exported comments must start with identifier name.
- Use `// TODO(username):` for tracking technical debt.

---

## Reference Links

### Official Go Documentation
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Modules](https://go.dev/ref/mod)

### Project Specific
- [Cobra (CLI)](https://github.com/spf13/cobra)
- [Bubble Tea (TUI)](https://github.com/charmbracelet/bubbletea)
- [Go Hyper-V (Reference)](https://github.com/sheepla/go-hyperv)
- [golangci-lint (v2.8.0)](https://golangci-lint.run/)
