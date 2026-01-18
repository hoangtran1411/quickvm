---
name: Go Project Structure
description: Standard Go project structure for CLI applications with Makefile, GitHub Actions CI/CD, and professional documentation.
---

# Go Project Structure

This skill provides a battle-tested project structure for Go CLI applications with emphasis on maintainability, CI/CD, and professional documentation.

## When to Use

- Starting a new Go CLI project
- Restructuring an existing Go project
- Setting up CI/CD for Go projects
- Creating release workflows

## Standard Directory Structure

```
project/
├── .github/
│   └── workflows/
│       ├── build.yml           # Build and test on push/PR
│       └── release.yml         # Release workflow for tags
├── .agent/
│   ├── rules/                  # AI agent rules
│   │   └── go-style-guide.md
│   └── skills/                 # Reusable AI skills
├── cmd/                        # CLI commands (Cobra)
│   ├── root.go
│   ├── version.go
│   └── <command>.go
├── internal/                   # Private packages (not importable)
│   └── config/
├── pkg/                        # Public packages (importable)
│   └── utils/
├── <domain>/                   # Domain-specific packages
│   ├── <domain>.go
│   ├── <domain>_test.go
│   └── mock_test.go
├── ui/                         # TUI components (if applicable)
│   └── table.go
├── docs/                       # Documentation
│   ├── README.md
│   ├── DEVELOPER.md
│   ├── CONTRIBUTING.md
│   └── FEATURE_ROADMAP.md
├── .gitignore
├── CHANGELOG.md
├── LICENSE
├── Makefile
├── README.md
├── go.mod
├── go.sum
└── main.go
```

## Key Files

### 1. `main.go` (Entry Point)

Keep it minimal - just call the command executor:

```go
package main

import "yourproject/cmd"

func main() {
    cmd.Execute()
}
```

### 2. `.gitignore`

```gitignore
# Binaries
*.exe
*.dll
*.so
*.dylib

# Build output
/build/
/dist/

# Test outputs
coverage.out
coverage.html
*.coverprofile

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Go
vendor/
```

### 3. `Makefile`

```makefile
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=myapp.exe
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_DIR=build

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the application
	go build -o $(BINARY_NAME) .

.PHONY: build-optimized
build-optimized: ## Build with size optimization
	go build -ldflags="-s -w" -o $(BINARY_NAME) .

.PHONY: build-all
build-all: ## Build for multiple architectures
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) .
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/myapp-arm64.exe .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/myapp-linux .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/myapp-darwin .

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: clean
clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod tidy

.PHONY: check
check: fmt lint test ## Run all checks

.PHONY: dev
dev: fmt test build ## Development cycle

.PHONY: release
release: clean test build-all ## Create release
```

### 4. GitHub Actions - Build (`.github/workflows/build.yml`)

```yaml
name: Build

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: windows-latest
    
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'

    - name: Cache Go modules
      uses: actions/cache@v4
      with:
        path: |
          ~\go\pkg\mod
          ~\AppData\Local\go-build
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-

    - name: Download dependencies
      run: go mod download

    - name: Run tests
      run: go test -v ./...

    - name: Build
      run: go build -ldflags="-s -w" -o myapp.exe

    - name: Upload artifact
      uses: actions/upload-artifact@v4
      with:
        name: myapp-windows
        path: myapp.exe
        retention-days: 30

  lint:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4

    - uses: actions/setup-go@v5
      with:
        go-version: '1.21'

    - name: golangci-lint
      uses: golangci/golangci-lint-action@v6
      with:
        version: latest
```

### 5. GitHub Actions - Release (`.github/workflows/release.yml`)

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: windows-latest
    
    steps:
    - uses: actions/checkout@v4

    - uses: actions/setup-go@v5
      with:
        go-version: '1.21'

    - name: Build AMD64
      run: |
        $env:GOOS="windows"
        $env:GOARCH="amd64"
        go build -ldflags="-s -w" -o myapp-windows-amd64.exe
      shell: pwsh

    - name: Build ARM64
      run: |
        $env:GOOS="windows"
        $env:GOARCH="arm64"
        go build -ldflags="-s -w" -o myapp-windows-arm64.exe
      shell: pwsh

    - name: Create checksums
      run: |
        Get-FileHash *.exe -Algorithm SHA256 | ForEach-Object {
          "$($_.Hash)  $(Split-Path $_.Path -Leaf)" | Out-File -Append checksums.txt
        }
      shell: pwsh

    - name: Create Release
      uses: softprops/action-gh-release@v2
      with:
        files: |
          myapp-windows-amd64.exe
          myapp-windows-arm64.exe
          checksums.txt
        generate_release_notes: true
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 6. `CHANGELOG.md` Template

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
- New feature description

### Changed
- Change description

### Fixed
- Bug fix description

## [1.0.0] - 2024-01-01

### Added
- Initial release
- Feature 1
- Feature 2

[Unreleased]: https://github.com/user/project/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/user/project/releases/tag/v1.0.0
```

### 7. `README.md` Template

```markdown
# Project Name 🚀

<div align="center">

![Version](https://img.shields.io/badge/version-1.0.0-blue?style=for-the-badge)
![Go](https://img.shields.io/badge/Go-1.21-00ADD8?style=for-the-badge&logo=go)
[![Build](https://github.com/user/project/actions/workflows/build.yml/badge.svg)](https://github.com/user/project/actions)

**A brief, compelling description of your project**

[Features](#features) • [Installation](#installation) • [Usage](#usage) • [Documentation](#documentation)

</div>

---

## ✨ Features

- 🎯 Feature 1 - Description
- ⚡ Feature 2 - Description
- 🛡️ Feature 3 - Description

## 📋 Prerequisites

- Go 1.21+
- Additional requirements...

## 🚀 Installation

### Quick Install

\`\`\`bash
go install github.com/user/project@latest
\`\`\`

### Build from Source

\`\`\`bash
git clone https://github.com/user/project.git
cd project
go build -o myapp
\`\`\`

## 📖 Usage

\`\`\`bash
# Basic usage
myapp command

# With options
myapp command --flag value
\`\`\`

## 🏗️ Architecture

\`\`\`
project/
├── cmd/        # CLI commands
├── pkg/        # Public packages
└── internal/   # Private packages
\`\`\`

## 📚 Documentation

- [Developer Guide](docs/DEVELOPER.md)
- [Contributing](docs/CONTRIBUTING.md)

## 🤝 Contributing

Contributions welcome! Please read [CONTRIBUTING.md](docs/CONTRIBUTING.md).

## 📝 License

MIT License - see [LICENSE](LICENSE)
```

## Best Practices

1. **Keep `main.go` minimal** - Only call `cmd.Execute()`

2. **Use `internal/`** for truly private code that shouldn't be imported

3. **Domain packages at root** - For domain logic (e.g., `hyperv/`, `user/`)

4. **Consistent naming**:
   - Files: `lowercase_with_underscores.go`
   - Test files: `*_test.go`
   - Mock files: `mock_test.go`

5. **Document exported symbols** - Every exported function/type needs a comment

6. **Use semantic versioning** - `v1.0.0`, `v1.1.0`, `v2.0.0`

7. **Keep CHANGELOG updated** - Document all notable changes
