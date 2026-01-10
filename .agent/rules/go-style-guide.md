---
trigger: always_on
---

---
trigger: always_on
---

Code Style:
- "Ensure all Go code is formatted using `gofmt` or `goimports`."
- "Adhere to [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)."
- "Organize code into domain-specific packages (e.g., `hyperv/`) for business logic. Use `internal/` only for code that should legally be hidden from external importers."

User Interface & UX (Charmbracelet):
- "Leverage `github.com/charmbracelet/lipgloss` for all CLI styling (colors, borders, margins). Avoid raw ANSI escape codes."
- "Use `github.com/charmbracelet/bubbletea` for interactive workflows (spinners, selection lists, forms)."
- "Ensure the CLI feels 'Premium': use consistent padding, clear headers, and intuitive navigation."

Error Handling:
- "Always wrap errors using `%w`: `fmt.Errorf(\"context: %w\", err)`. This is critical for tracing Hyper-V and Shell execution errors."
- "Implement 'fail fast' logic using guard clauses to minimize indentation."
- "Do not ignore errors in `defer` statements (e.g., closing files), log them if handling isn't possible."

Context & Concurrency:
- "Every function performing I/O or calling external commands (PowerShell) MUST accept `context.Context` as its first argument."
- "Use `context` to manage timeouts and cancellations for long-running Hyper-V operations."

Security & Shell Execution:
- "NEVER construct PowerShell commands using simple string concatenation with user input. Use proper argument quoting or sanitization helpers to prevent Command Injection."
- "Prefer `exec.CommandContext` over `exec.Command`."

Documentation:
- "Every exported function, variable, and type must have clear documentation comments explaining 'Why' rather than just 'What'."
- "Provide usage examples for complex packages in a `doc.go` file."

Testing:
- "Prioritize Table-driven tests combined with `t.Run` for comprehensive test coverage."
- "Mock Hyper-V interfaces to ensure unit tests can run without Administrator privileges or a Windows environment."
- "Use the 'TestMain' pattern or build tags (e.g., `//go:build windows`) for integration tests that require actual Hyper-V."

Efficiency & Tone:
- "Avoid greetings, apologies, or meta-commentary; focus strictly on code and execution logs."
- "Provide code as minimal diffs/blocks whenever possible."

Reference & Resource Mapping (GitHub Repositories)

### CLI Implementation:
- **Reference:** [spf13/cobra](https://github.com/spf13/cobra)
- **Guideline:** Use Cobra's structure for subcommands (e.g., `vm start`, `vm stop`). Follow the "Command Pattern" to decouple CLI logic from Hyper-V business logic.

### Hyper-V & WMI Integration:
- **Reference:** [sheepla/go-hyperv](https://github.com/sheepla/go-hyperv)
- **Guideline:** Model Hyper-V object structures (VM, VHD, Network Switch) based on this repo's WMI queries. Ensure type safety when parsing WMI outputs.

### Performance & Profiling:
- **Reference:** [google/pprof](https://github.com/google/pprof)
- **Guideline:** Integrate `net/http/pprof` in long-running management daemons to monitor memory leaks during heavy VM orchestration tasks.

### Automation Logic:
- **Reference:** [fdcastel/Hyper-V-Automation](https://github.com/fdcastel/Hyper-V-Automation)
- **Guideline:** Use the logic from these PowerShell scripts as a blueprint for the `os/exec` calls in Go. Ensure all scripts are executed with `context.Context` for strict timeout management.