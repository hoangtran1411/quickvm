# Contributing to QuickVM

First off, thank you for considering contributing to QuickVM! It's people like you that make QuickVM such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by respect and professionalism. By participating, you are expected to uphold this standard.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues as you might find that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps which reproduce the problem**
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed after following the steps**
* **Explain which behavior you expected to see instead and why**
* **Include screenshots if possible**

**Bug Report Template:**
```markdown
**Environment:**
- OS: Windows 10/11
- QuickVM Version: x.x.x
- Go Version: x.x.x
- Hyper-V Version: x.x

**Description:**
A clear description of the bug.

**Steps to Reproduce:**
1. Step 1
2. Step 2
3. ...

**Expected Behavior:**
What you expected to happen.

**Actual Behavior:**
What actually happened.

**Additional Context:**
Any other context about the problem.
```

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title**
* **Provide a step-by-step description of the suggested enhancement**
* **Provide specific examples to demonstrate the steps**
* **Describe the current behavior and why it's insufficient**
* **Explain why this enhancement would be useful**

### Pull Requests

* Fill in the required template
* Follow the Go coding style
* Include thoughtful comments in your code
* Write clear commit messages
* Include tests when adding new features
* Update documentation as needed
* End all files with a newline

## Development Setup

### Prerequisites

1. Go 1.21 or higher
2. Windows 10/11 with Hyper-V enabled
3. Git
4. (Optional) golangci-lint for code quality checks

### Setting Up Your Development Environment

```powershell
# 1. Fork and clone the repository
git clone https://github.com/YOUR-USERNAME/quickvm.git
cd quickvm

# 2. Install dependencies
go mod download

# 3. Create a new branch for your feature
git checkout -b feature/my-new-feature

# 4. Make your changes and test
go build -o quickvm.exe
.\quickvm.exe list

# 5. Run tests
go test ./...

# 6. Format your code
go fmt ./...

# 7. (Optional) Run linter
golangci-lint run
```

## Coding Guidelines

### Go Style Guide

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and:

1. **Use gofmt**: Always format your code with `go fmt`
2. **Write tests**: Add tests for new functionality
3. **Handle errors**: Never ignore errors, handle them appropriately
4. **Comment exports**: All exported functions, types, and constants should have comments
5. **Keep it simple**: Write clear, readable code over clever code

### Code Organization

```
quickvm/
├── cmd/         # CLI commands (one file per command)
├── hyperv/      # Hyper-V integration logic
├── ui/          # TUI components
└── main.go      # Entry point
```

### Commit Messages

Write clear, concise commit messages:

```
feat: add support for VM snapshots
fix: correct memory calculation in list command
docs: update README with new examples
test: add tests for hyperv package
refactor: simplify PowerShell script generation
```

Use prefixes:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Adding or updating tests
- `refactor:` - Code refactoring
- `style:` - Code style changes (formatting, etc.)
- `perf:` - Performance improvements
- `chore:` - Maintenance tasks

### PowerShell Scripts

When modifying PowerShell scripts:

1. Test scripts independently before integrating
2. Use explicit type conversions `[int]`, `[string]`
3. Handle errors with proper error messages
4. Keep scripts simple and readable
5. Comment complex PowerShell logic

### Testing

Run tests before submitting:

```powershell
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./hyperv

# Run benchmarks
go test -bench=. ./...
```

## Pull Request Process

1. **Update your fork**
   ```powershell
   git remote add upstream https://github.com/original/quickvm.git
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```

2. **Create a feature branch**
   ```powershell
   git checkout -b feature/my-feature
   ```

3. **Make your changes**
   - Write code
   - Add tests
   - Update documentation

4. **Test thoroughly**
   ```powershell
   go test ./...
   go build -o quickvm.exe
   # Manual testing
   ```

5. **Commit your changes**
   ```powershell
   git add .
   git commit -m "feat: add my amazing feature"
   ```

6. **Push to your fork**
   ```powershell
   git push origin feature/my-feature
   ```

7. **Create Pull Request**
   - Go to GitHub
   - Click "New Pull Request"
   - Provide clear description
   - Link related issues

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe testing performed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review performed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests pass
- [ ] No new warnings
```

## Project Structure

### Adding a New Command

1. Create `cmd/mycommand.go`:
```go
package cmd

import (
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Short description",
    Long:  `Long description`,
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

2. Add tests in `cmd/mycommand_test.go`

### Adding Hyper-V Functionality

1. Add methods to `hyperv/hyperv.go`
2. Follow existing patterns for PowerShell execution
3. Handle errors appropriately
4. Add tests in `hyperv/hyperv_test.go`

### Modifying TUI

1. Update `ui/table.go`
2. Follow Bubble Tea patterns
3. Test interactive features manually

## Documentation

Update documentation when:
- Adding new features
- Changing existing behavior
- Fixing bugs that affect usage
- Adding configuration options

Files to consider:
- `README.md` - Main documentation
- `HUONG_DAN.md` - Vietnamese guide
- `DEMO.md` - Examples and use cases
- `WORKFLOW.md` - Development workflow
- Code comments

## Release Process

Maintainers will handle releases, but contributors should:
1. Update version in `cmd/version.go` if needed
2. Update CHANGELOG.md (if exists)
3. Tag releases with semantic versioning

## Questions?

Feel free to:
- Open an issue for questions
- Join discussions in existing issues
- Reach out to maintainers

## Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes
- Project README (for significant contributions)

Thank you for contributing to QuickVM! 🚀
