# QuickVM Architecture

This document describes the high-level architecture of QuickVM. It is intended for AI Agents and developers to understand the system design, components, and data flow.

## High-Level Overview

QuickVM follows a layered architecture designed for testability, security, and a superior Agent Experience (AX).

```text
┌─────────────────────────────────────────────────────────┐
│                    User / AI Agent                      │
└──────────────────────────┬──────────────────────────────┘
                           │
           ┌───────────────┴───────────────┐
           │        CLI (Cobra)            │
           │  (cmd/, main.go, root.go)     │
           └───────┬───────────────┬───────┘
                   │               │
       ┌───────────▼───┐       ┌───▼───────────┐
       │   TUI Layer   │       │  Output Layer │
       │ (Bubble Tea)  │       │ (int/output)  │
       │     (ui/)     │       └───────▲───────┘
       └───────────┬───┘               │
                   │                   │
       ┌───────────▼───────────────────┴───────┐
       │         Logic Layer (Core)            │
       │          (internal/hyperv)            │
       └───────────────────┬───────────────────┘
                           │
       ┌───────────────────▼───────────────────┐
       │       Infrastructure Layer            │
       │    (ShellExecutor / PowerShell)       │
       └───────────────────┬───────────────────┘
                           │
       ┌───────────────────▼───────────────────┐
       │         Windows Hyper-V API           │
       └───────────────────────────────────────┘
```

## Architectural Layers

### 1. Presentation Layer (CLI & TUI)
- **`cmd/`**: Uses `spf13/cobra`. Handles flag parsing, argument validation, and command orchestration. It acts as the primary entry point for both humans (CLI) and AI agents (CLI with `--output json`).
- **`ui/`**: Uses `charmbracelet/bubbletea` and `lipgloss`. Provides a rich interactive experience for humans. It interacts with the Logic Layer via `tea.Cmd`.

### 2. Logic Layer (Core Service)
- **`internal/hyperv`**: Contains the business logic for VM management.
- **`Manager`**: The central orchestrator. It doesn't execute PowerShell directly but uses the `ShellExecutor` interface.
- **`VMManager` Interface**: Defines the contract for all VM operations, allowing for easy mocking.

### 3. Output Layer (AX Foundation)
- **`internal/output`**: Standardizes responses across all commands. 
- **JSON Serialization**: Ensures predictable, machine-readable output for AI agents.
- **Error Handling**: Provides structured error codes (`VM_NOT_FOUND`, `HYPERV_ERROR`) instead of just raw strings.

### 4. Infrastructure Layer (OS Integration)
- **`ShellExecutor` Interface**: Abstracts the execution of shell commands.
- **`PowerShellRunner`**: The concrete implementation that calls `powershell.exe`. It enforces security by using `exec.CommandContext` with separate arguments to prevent command injection.

## Data Flow

### Command Execution (CLI)
1. `root rootCmd` receives command and flags.
2. `PersistentPreRunE` configures the `output` package (format: json/table).
3. The specific command (`list`, `start`, etc.) calls `internal/hyperv`.
4. `internal/hyperv` executes logic and requests data via `ShellExecutor`.
5. Results are returned to the command.
6. The command uses `internal/output` to print results to stdout.

### Interactive Loop (TUI)
1. `main.go` launches `ui.NewModel()`.
2. `bubbletea` manages state and rendering.
3. User actions trigger `tea.Cmd`.
4. `tea.Cmd` calls methods on `internal/hyperv`.
5. Results are returned as `tea.Msg` to update the UI model.

## Key Design Patterns

- **Dependency Inversion**: High-level modules (`cmd`, `hyperv.Manager`) depend on abstractions (`VMManager`, `ShellExecutor`).
- **Command Pattern**: Encapsulating VM operations as distinct CLI commands.
- **Response Wrapping**: All JSON outputs wrap data/errors in a consistent `Response` struct.
- **Secure Shell Execution**: Strict separation of commands and arguments to prevent injection.

## Technology Stack

- **Language**: Go
- **CLI Framework**: Cobra
- **TUI Framework**: Bubble Tea / Lipgloss / Bubbles
- **Serialization**: Standard `encoding/json`
- **Testing**: Table-driven tests with `testify` (optional) and native `testing`.
- **Target OS**: Windows (Hyper-V required)
