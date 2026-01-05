# QuickVM Development & Deployment Workflow

## Development

### 1. Setup Development Environment
```bash
# Navigate to project
cd d:\Workspace\Dev\Learning\Golang\quickvm

# Install dependencies
go mod download
go mod tidy
```

### 2. Build the Application
```bash
# Build for current platform (Windows)
go build -o quickvm.exe

# Build with optimizations
go build -ldflags="-s -w" -o quickvm.exe

# Build for different architectures
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o quickvm-amd64.exe
$env:GOOS="windows"; $env:GOARCH="arm64"; go build -o quickvm-arm64.exe
```

### 3. Test the Application
```bash
# Test listing VMs
.\quickvm.exe list

# Test starting a VM (replace 1 with actual VM index)
.\quickvm.exe start 1

# Test TUI mode
.\quickvm.exe
```

### 4. Code Quality
```bash
# Format code
go fmt ./...

# Run linter (install golangci-lint first)
golangci-lint run

# Run tests (when tests are added)
go test ./...

# Check for vulnerabilities
go mod verify
```

## Deployment

### 1. Local Installation
```powershell
# Option A: Copy to Windows System32 (requires Admin)
Copy-Item quickvm.exe C:\Windows\System32\

# Option B: Add to custom bin directory
$binDir = "$env:USERPROFILE\bin"
New-Item -ItemType Directory -Force -Path $binDir
Copy-Item quickvm.exe $binDir\
# Then add $binDir to PATH in System Environment Variables

# Option C: Add current directory to PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
$newPath = $currentPath + ";D:\Workspace\Dev\Learning\Golang\quickvm"
[Environment]::SetEnvironmentVariable("Path", $newPath, "User")
```

### 2. Create Installer (Optional)
```bash
# Using NSIS or similar tool to create Windows installer
# Or use chocolatey package
```

### 3. Distribution
```bash
# Create release archive
Compress-Archive -Path quickvm.exe, README.md, HUONG_DAN.md -DestinationPath quickvm-v1.0.0.zip

# Upload to GitHub Releases or other platforms
```

## Project Structure

```
quickvm/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command + TUI launcher
│   ├── start.go           # Start VM command
│   ├── stop.go            # Stop VM command
│   ├── restart.go         # Restart VM command
│   └── list.go            # List VMs command
├── hyperv/                # Hyper-V integration
│   └── hyperv.go          # VM management via PowerShell
├── ui/                    # TUI components
│   └── table.go           # Interactive table view
├── main.go                # Application entry point
├── go.mod                 # Go modules
├── go.sum                 # Dependencies checksums
├── README.md              # English documentation
├── HUONG_DAN.md          # Vietnamese guide
└── .gitignore            # Git ignore rules
```

## Best Practices

1. **Version Control**
   - Commit frequently with meaningful messages
   - Use semantic versioning (v1.0.0, v1.1.0, etc.)
   - Tag releases in Git

2. **Code Quality**
   - Follow Go conventions and best practices
   - Add comments for exported functions
   - Keep functions small and focused
   - Handle errors properly

3. **Testing**
   - Add unit tests for critical functions
   - Test with different VM configurations
   - Test error scenarios

4. **Documentation**
   - Keep README and guides up to date
   - Document new features
   - Add code examples

## Troubleshooting Development Issues

### Import Issues
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
go mod tidy
```

### Build Issues
```bash
# Clean build cache
go clean -cache

# Rebuild
go build -v -o quickvm.exe
```

### PowerShell Execution Policy
```powershell
# If you can't run test-vm.ps1
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## Future Enhancements

- [ ] Add config file support (.quickvm.yaml)
- [ ] Support for VM snapshots management
- [ ] Support for creating new VMs
- [ ] Export/import VM configurations
- [ ] Remote Hyper-V server support
- [ ] Performance metrics and monitoring
- [ ] Batch operations on multiple VMs
- [ ] Custom keybindings
- [ ] VM grouping/tagging
- [ ] Search/filter VMs by name

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - feel free to use and modify as needed.
