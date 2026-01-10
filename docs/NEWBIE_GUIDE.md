# Newbie Guide: Understanding the QuickVM Codebase

Welcome! This guide is designed specifically for developers who are **new to Go** and want to understand how `quickvm` is structured. We'll break down the project flow so you can navigate and contribute effectively.

## 🗺️ High-Level Architecture

QuickVM follows a standard Go CLI (Command Line Interface) structure. Here is the mental model you should have:

```mermaid
graph TD
    User[User] --> Main[main.go]
    Main --> Cobra[cmd package (Cobra)]
    Cobra --> Logic[hyperv package (Business Logic)]
    Logic --> PowerShell[PowerShell / Hyper-V]
```

1.  **Entry Point**: The user runs the exe.
2.  **Routing**: The `cmd` package decides which command to run (e.g., `create`, `start`, `list`).
3.  **Execution**: The command calls functions in the `hyperv` package.
4.  **System Interaction**: The `hyperv` package executes PowerShell scripts to talk to Windows.

## 🚀 The Entry Point: `main.go`

If you open `main.go`, you'll see it's very small. This is a best practice in Go.

```go
package main

import "quickvm/cmd"

func main() {
    cmd.Execute()
}
```

**Key Takeaway**: `main.go` should only drive the application startup. All the real work happens in packages.

## 🛠️ The CLI Layer: `cmd/`

This project uses [Cobra](https://github.com/spf13/cobra), the most popular Go library for CLIs.

Navigate to `cmd/`. You will see files like `start.go`, `stop.go`, `list.go`. Each file typically corresponds to one command.

### Anatomy of a Command (`cmd/example.go`)

Most command files follow this pattern:

```go
// 1. Define the command variable
var startCmd = &cobra.Command{
    Use:   "start [vm-name]",
    Short: "Start a virtual machine",
    // 2. The main logic function
    RunE: func(cmd *cobra.Command, args []string) error {
        // Logic goes here
        return nil
    },
}

// 3. Register the command
func init() {
    rootCmd.AddCommand(startCmd)
}
```

**Go Concepts for Newbies:**
*   **`init()` function**: Go runs this automatically before `main()`. We use it to register our `startCmd` to the main `rootCmd`.
*   **`RunE`**: We prefer `RunE` over `Run` because it allows us to return an `error`. Go loves explicit error handling!

## 🧠 The Logic Layer: `hyperv/`

The `cmd` package should **not** contain specific Hyper-V logic. It should only parse flags and arguments. The real work happens in the `hyperv` folder.

This separation makes the code:
1.  **Testable**: We can test logic without running CLI commands.
2.  **Reusable**: Other Go programs could import our `hyperv` package.

### How we call PowerShell

Since Hyper-V is managed via PowerShell, this project wraps PowerShell commands. Look at `hyperv/manager.go` or specific feature files.

```go
func (m *Manager) StartVM(vmName string) error {
    // We construct a command string
    cmd := fmt.Sprintf("Start-VM -Name '%s'", vmName)
    
    // We execute it (simplified)
    return runPowerShell(cmd)
}
```

## 🎓 Key Go Idioms in Project

### 1. Error Handling (`if err != nil`)
You will see this everywhere. We don't use exceptions.
```go
if err := vm.Start(); err != nil {
    // Wrap errors to provide context!
    return fmt.Errorf("failed to start vm: %w", err)
}
```
*   `%w`: Wraps the error so we can unwrap it later to check the cause.

### 2. Context (`ctx`)
For long-running operations (like creating a 50GB VM), we use `context.Context`.
```go
func CreateVM(ctx context.Context, name string) error {
    // ...
}
```
This allows us to cancel the operation if the user hits Ctrl+C or if it times out.

## 👣 Walkthrough: Adding Your First Command

Let's say you want to add a command `quickvm hello`.

1.  **Create the file**: Create `cmd/hello.go`.
2.  **Define the command**:
    ```go
    package cmd

    import (
        "fmt"
        "github.com/spf13/cobra"
    )

    var helloCmd = &cobra.Command{
        Use:   "hello",
        Short: "Prints hello",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello, QuickVM Newbie!")
        },
    }

    func init() {
        rootCmd.AddCommand(helloCmd)
    }
    ```
3.  **Run it**: `go run main.go hello`

## 📚 Recommended Learning Path

1.  **Read `cmd/root.go`**: See how the base command is set up.
2.  **Read `hyperv/vm.go`** (or similar): See how a struct models a VM.
3.  **Check `DEVELOPER.md`**: For deeper architectural details.
4.  **Try to fix a "TODO"**: Search the codebase for `TODO` comments.

Happy Coding! 🚀
