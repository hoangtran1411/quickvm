# Contributing to QuickVM

Thank you for your interest in contributing to **QuickVM**! 🎉  
QuickVM is an open-source, fast, and lightweight Hyper-V Virtual Machine
manager CLI and TUI for Windows, built with Go.

Whether you are fixing a bug, suggesting a new feature, improving
documentation, or submitting a pull request, your contributions are warmly
welcomed.

---

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Improving Documentation](#improving-documentation)
  - [Code Contributions](#code-contributions)
- [Development Setup](#development-setup)
  - [Prerequisites](#prerequisites)
  - [First-Time Setup](#first-time-setup)
- [Development Workflow](#development-workflow)
  - [Building](#building)
  - [Running Tests](#running-tests)
  - [Code Quality & Linting](#code-quality--linting)
- [Coding & Architecture Guidelines](#coding--architecture-guidelines)
  - [Project Structure](#project-structure)
  - [PowerShell Security (CRITICAL)](#powershell-security-critical)
  - [Context & Process Lifecycle](#context--process-lifecycle)
  - [Error Handling](#error-handling)
  - [Testing Standards](#testing-standards)
  - [AI Agent & Machine-Readable Output](#ai-agent--machine-readable-output)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)
- [Questions & Support](#questions--support)

---

## Code of Conduct

This project and everyone participating in it is governed by our
[Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected
to uphold this code. Please report unacceptable behavior following the
guidelines in [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

---

## How Can I Contribute?

### Reporting Bugs

Before submitting a bug report:

1. Check the [existing issues](https://github.com/hoangtran1411/quickvm/issues)
   to avoid duplicates.
2. Verify that you are running the latest version of QuickVM
   (`quickvm version` or `quickvm update --check-only`).

When opening an issue, please include:

- **Environment**: Windows OS version (e.g., Windows 11 23H2), Go version,
  QuickVM version.
- **Clear description**: What happened vs. what you expected to happen.
- **Steps to reproduce**: Minimal, step-by-step commands.
- **Terminal output / logs**: Relevant output or screenshots (redacting any
  private or sensitive information).

### Suggesting Enhancements

Feature requests are welcome! QuickVM is designed around being
**fast, focused, and intuitive**. When proposing features:

- Explain the real-world use case and why it fits QuickVM's mission.
- Provide example CLI syntax or TUI interaction flows.
- Review [docs/FEATURE_ROADMAP.md](docs/FEATURE_ROADMAP.md) first to ensure
  the feature is not archived or out-of-scope.

### Improving Documentation

Documentation improvements (fixing typos, clarifying guides, adding examples,
or documenting edge cases) are always appreciated.

### Code Contributions

Looking for a place to contribute? Check open issues labeled `good first issue`
or review the feature roadmap. Please review the
[Development Setup](#development-setup) and
[Pull Request Process](#pull-request-process) sections before submitting.

---

## Development Setup

### Prerequisites

- **OS**: Windows 10 or Windows 11 with Hyper-V enabled.
- **Privileges**: Administrator privileges (required to interact with
  Hyper-V cmdlets and WMI).
- **Go**: Version 1.21 or higher (Go 1.27.0 is used in CI).
- **Git**: Latest version.
- **golangci-lint**: Version `v2.8.0` or higher (uses v2 schema).

### First-Time Setup

1. **Fork** the repository on GitHub:
   [https://github.com/hoangtran1411/quickvm](https://github.com/hoangtran1411/quickvm)
2. **Clone** your fork locally:

   ```powershell
   git clone https://github.com/<your-username>/quickvm.git
   cd quickvm
   ```

3. **Add upstream remote**:

   ```powershell
   git remote add upstream https://github.com/hoangtran1411/quickvm.git
   ```

4. **Download dependencies**:

   ```powershell
   go mod download
   go mod verify
   ```

---

## Development Workflow

### Building

Build the binary locally:

```powershell
# Standard development build
go build -o quickvm.exe

# Optimized production build (stripped symbols)
go build -ldflags="-s -w" -o quickvm.exe

# Cross-compiling for architectures
$env:GOOS="windows"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -o quickvm-windows-amd64.exe

$env:GOOS="windows"; $env:GOARCH="arm64"
go build -ldflags="-s -w" -o quickvm-windows-arm64.exe
```

### Running Tests

Unit tests are designed to run fast on any machine by using mocked
PowerShell executors:

```powershell
# Run all unit tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage profile
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality & Linting

QuickVM enforces strict code quality checks in CI. Ensure your code passes
all checks before opening a pull request:

```powershell
# Format code
gofmt -w .
goimports -w .

# Run static analysis
go vet ./...

# Run golangci-lint (v2.8.0+)
golangci-lint run ./...
```

> ⚠️ **Note**: Our `.golangci.yml` uses the **v2 schema** (`version: "2"`).
> Ensure your `golangci-lint` is updated to at least `v2.8.0`.

---

## Coding & Architecture Guidelines

### Project Structure

```text
quickvm/
├── cmd/                    # Cobra CLI commands & subcommands (entry points)
│   ├── root.go            # Root command & TUI launcher
│   ├── start.go           # quickvm start
│   ├── stop.go            # quickvm stop
│   ├── restart.go         # quickvm restart
│   ├── list.go            # quickvm list / ls
│   └── output_helpers.go  # Batch execution & worker pool logic
├── internal/
│   ├── hyperv/            # Hyper-V core domain logic & PowerShell wrappers
│   │   ├── executor.go    # ShellExecutor interface for mockable execution
│   │   ├── hyperv.go      # VM operations (start, stop, query, etc.)
│   │   ├── snapshot.go    # Checkpoint / snapshot operations
│   │   └── sysinfo.go     # Host system information gathering
│   └── output/            # Table & JSON formatters
├── ui/                    # Terminal User Interface (Bubble Tea & Lipgloss)
│   └── table.go           # Interactive dashboard
├── updater/               # Self-update logic and GitHub release checking
└── main.go                # Application main entry point
```

### PowerShell Security (CRITICAL)

Because QuickVM executes with Administrator privileges and interacts with
PowerShell, **security against command injection is vital**:

- ❌ **NEVER** concatenate unsanitized user input into script strings or
  command arguments:

  ```go
  // DANGEROUS - Vulnerable to command injection
  exec.Command("powershell", "-Command", "Get-VM -Name " + userInput)
  ```

- ✅ **ALWAYS** pass arguments as separate elements or use parameterization:

  ```go
  // SAFE - Separate arguments passed to executor
  exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive",
      "-Command", "Get-VM", "-Name", userInput)
  ```

### Context & Process Lifecycle

- Every function performing I/O or external command execution
  **must accept `context.Context` as its first parameter**.
- Always use `context.WithTimeout` for external operations to prevent
  orphaned or hanging PowerShell processes.
- Clean up resources with `defer cancel()` immediately after context creation.

### Error Handling

- Always wrap errors with descriptive context using `%w`:

  ```go
  if err != nil {
      return fmt.Errorf("failed to start VM %q: %w", vmName, err)
  }
  ```

- Fail fast using guard clauses.
- Do not silently discard error return values (e.g., in `os.WriteFile`).

### Testing Standards

- Use **table-driven tests** for all unit tests.
- Mock external PowerShell interactions via the `ShellExecutor` interface
  so tests execute reliably without requiring a live Hyper-V host.
- For integration tests requiring real Windows Hyper-V, isolate them with
  `//go:build windows`.

### AI Agent & Machine-Readable Output

QuickVM commands support `--output json` (`-o json`) for automated tooling
and AI agents. When introducing new commands or modifying output formats:

- Ensure output structures serialize cleanly to JSON.
- Maintain consistent response structures (`success`, `data`, and `error`).

---

## Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/)
specification:

```text
<type>(<optional scope>): <description>

[optional body]

[optional footer(s)]
```

### Allowed Types

- `feat`: A new feature or capability
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Formatting changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to build process, tooling, or dependencies

### Examples

```text
feat(cmd): add batch restart support for multiple VMs
fix(hyperv): resolve timeout error during snapshot creation
docs: update installation instructions for Windows ARM64
test(hyperv): add table-driven test for VM list parsing
```

---

## Pull Request Process

1. **Keep PRs focused**: Address a single feature or bug fix per pull request.
2. **Sync with main**: Rebase or merge upstream `main` into your feature branch.
3. **Verify locally**:
   - `go test ./...` passes.
   - `golangci-lint run ./...` reports no issues.
   - Code is formatted with `gofmt -w .`.
4. **Submit your PR**:
   - Fill out the PR template completely.
   - Link any related issues (`Fixes #123`).
5. **Code Review**: Address feedback promptly. Maintainers will review your
   submission and merge once all CI checks pass.

---

## Questions & Support

If you have questions or need guidance before starting work:

- Open a discussion or issue on
  [GitHub Issues](https://github.com/hoangtran1411/quickvm/issues).
- Check the [Developer Guide](docs/DEVELOPER.md) and
  [Getting Started](docs/GETTING_STARTED.md) for more technical details.

Thank you for helping make QuickVM better! 🚀
