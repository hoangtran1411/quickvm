# QuickVM Makefile
# For Windows PowerShell

.DEFAULT_GOAL := help

# Variables
BINARY_NAME=quickvm.exe
MAIN_PATH=.
BUILD_DIR=build
VERSION?=1.0.0

# Colors for output
GREEN=\033[0;32m
NC=\033[0m # No Color

.PHONY: help
help: ## Show this help message
	@echo "QuickVM - Makefile Commands"
	@echo "=============================="
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the application
	@echo "Building QuickVM..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

.PHONY: build-optimized
build-optimized: ## Build with optimizations (smaller binary)
	@echo "Building optimized QuickVM..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Optimized build complete: $(BINARY_NAME)"

.PHONY: build-all
build-all: ## Build for all Windows architectures
	@echo "Building for all architectures..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/quickvm-arm64.exe $(MAIN_PATH)
	@echo "All builds complete in $(BUILD_DIR)/"

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: bench
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

.PHONY: lint
lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	golangci-lint run

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@if exist $(BINARY_NAME) del /Q $(BINARY_NAME)
	@if exist $(BUILD_DIR) rd /S /Q $(BUILD_DIR)
	@if exist coverage.out del /Q coverage.out
	@if exist coverage.html del /Q coverage.html
	@echo "Clean complete"

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated"

.PHONY: run
run: build ## Build and run the application
	@echo "Running QuickVM..."
	./$(BINARY_NAME)

.PHONY: run-list
run-list: build ## Build and run list command
	./$(BINARY_NAME) list

.PHONY: install
install: build-optimized ## Build and install to Windows System32 (requires admin)
	@echo "Installing QuickVM to System32..."
	copy $(BINARY_NAME) C:\Windows\System32\
	@echo "Installation complete. You can now use 'quickvm' from anywhere."

.PHONY: install-user
install-user: build-optimized ## Build and install to user bin directory
	@echo "Installing QuickVM to user bin..."
	@if not exist "%USERPROFILE%\bin" mkdir "%USERPROFILE%\bin"
	copy $(BINARY_NAME) "%USERPROFILE%\bin\"
	@echo "Installation complete."
	@echo "Add %USERPROFILE%\bin to your PATH if not already added."

.PHONY: release
release: clean test build-all ## Create release build
	@echo "Creating release package..."
	@mkdir -p $(BUILD_DIR)
	copy README.md $(BUILD_DIR)\
	copy HUONG_DAN.md $(BUILD_DIR)\
	copy DEMO.md $(BUILD_DIR)\
	@echo "Release package created in $(BUILD_DIR)/"

.PHONY: dev
dev: fmt test build run ## Full development cycle: format, test, build, and run

.PHONY: check
check: fmt lint test ## Run all checks: format, lint, and test
	@echo "All checks passed!"

.PHONY: version
version: ## Show version information
	@echo "QuickVM version: $(VERSION)"
	@echo "Go version:"
	@go version
