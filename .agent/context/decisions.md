# Architecture Decision Log

This document records key architectural decisions for AI agent context.

## ADR-001: PowerShell over Windows API

**Status**: Accepted

**Context**: Need to interact with Hyper-V for VM management.

**Decision**: Use PowerShell cmdlets via `os/exec` instead of Windows API (WMI/WinRM).

**Rationale**:
- ✅ Simpler implementation (no CGO, no COM)
- ✅ Uses official Hyper-V cmdlets with full feature support
- ✅ Better error messages from PowerShell
- ✅ Easier debugging (test scripts in ISE)
- ❌ ~200-500ms overhead per command (acceptable for interactive use)

**Consequences**:
- All Hyper-V operations go through `ShellExecutor` interface
- JSON parsing required for complex data
- Must handle PowerShell encoding quirks

---

## ADR-002: Index-Based VM Reference

**Status**: Accepted

**Context**: Users need to specify VMs in CLI commands.

**Decision**: Use 1-based numeric indices instead of VM names.

**Rationale**:
- ✅ Faster typing: `quickvm start 1` vs `quickvm start "Dev VM"`
- ✅ No escaping needed for special characters
- ✅ Consistent with TUI row numbers
- ❌ Indices change when VMs are added/removed

**Mitigations**:
- Always show VM name in command output
- `list` command shows index → name mapping

---

## ADR-003: ShellExecutor Interface

**Status**: Accepted

**Context**: Need to test Hyper-V operations without Windows.

**Decision**: Abstract PowerShell execution behind `ShellExecutor` interface.

```go
type ShellExecutor interface {
    RunCommand(ctx context.Context, script string) ([]byte, error)
}
```

**Rationale**:
- ✅ Unit tests run on any OS with mock executor
- ✅ Easy to swap implementations
- ✅ Follows dependency injection pattern

---

## ADR-004: Bubble Tea for TUI

**Status**: Accepted

**Context**: Need interactive terminal UI.

**Decision**: Use Charmbracelet's Bubble Tea + Bubbles + Lipgloss.

**Rationale**:
- ✅ Elm Architecture (predictable state)
- ✅ Excellent documentation and examples
- ✅ Active community
- ✅ Beautiful default styling
- ✅ Works well with async operations

---

## ADR-005: Cobra for CLI

**Status**: Accepted

**Context**: Need CLI command structure.

**Decision**: Use Cobra framework.

**Rationale**:
- ✅ De facto standard for Go CLIs
- ✅ Automatic help generation
- ✅ Easy subcommand structure
- ✅ Flag parsing built-in
- ✅ Shell completion support

---

## ADR-006: golangci-lint v2

**Status**: Accepted

**Context**: Need consistent linting.

**Decision**: Use golangci-lint v2.8.0+ with v2 schema.

**Rationale**:
- ✅ Latest features and linters
- ✅ Better performance
- ✅ Cleaner config format

**Requirements**:
- Config MUST have `version: "2"` at top
- Use kebab-case for linter settings
- Exclusions under `linters: exclusions: rules`

---

## ADR-007: Conventional Commits

**Status**: Accepted

**Context**: Need consistent commit messages.

**Decision**: Follow Conventional Commits specification.

**Format**:
```
<type>(<scope>): <description>

[body]

[footer]
```

**Types**: feat, fix, docs, style, refactor, test, chore

---

## Future Decisions Pending

- [ ] Remote Hyper-V support (WinRM vs SSH)
- [ ] VM templates system
- [ ] Plugin architecture
- [ ] Web UI option
