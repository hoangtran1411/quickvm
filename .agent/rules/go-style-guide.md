---
trigger: always_on
---

---
trigger: always_on
---

Code Style:
- "Ensure all Go code is formatted using `gofmt` or `goimports`."
- "Adhere to [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)."
- "Use `internal/` for core logic that should not be exposed, and `cmd/` for CLI entrypoints."

Error Handling:
- "Always wrap errors using `%w`: `fmt.Errorf(\"context: %w\", err)`. This is critical for tracing Hyper-V and Shell execution errors."
- "Implement 'fail fast' logic using guard clauses to minimize indentation."

Context & Concurrency:
- "Every function performing I/O or calling external commands (PowerShell) MUST accept `context.Context` as its first argument."
- "Use `context` to manage timeouts and cancellations for long-running Hyper-V operations."

Documentation:
- "Every exported function, variable, and type must have clear documentation comments explaining 'Why' rather than just 'What'."
- "Provide usage examples for complex packages in a `doc.go` file."

Testing:
- "Prioritize Table-driven tests combined with `t.Run` for comprehensive test coverage."
- "Mock Hyper-V interfaces to ensure unit tests can run without Administrator privileges or a Windows environment."

Efficiency & Tone:
- "Avoid greetings, apologies, or meta-commentary; focus strictly on code and execution logs."
- "Provide code as minimal diffs/blocks whenever possible."