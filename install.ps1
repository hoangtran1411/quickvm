# QuickVM Installation Script
# This script builds and installs QuickVM to your system

param(
    [Parameter()]
    [ValidateSet('System', 'User', 'Current')]
    [string]$InstallLocation = 'User',
    
    [Parameter()]
    [switch]$SkipBuild,
    
    [Parameter()]
    [switch]$CreateAlias
)

$ErrorActionPreference = "Stop"

Write-Host "╔══════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   QuickVM - Fast Hyper-V Virtual Machine Manager     ║" -ForegroundColor Cyan
Write-Host "║              Installation Script v1.0                ║" -ForegroundColor Cyan
Write-Host "╚══════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Check for Administrator rights if installing to System
if ($InstallLocation -eq 'System') {
    $currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    $isAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    
    if (-not $isAdmin) {
        Write-Host "❌ Error: System installation requires Administrator privileges" -ForegroundColor Red
        Write-Host "   Please run this script as Administrator or use -InstallLocation User" -ForegroundColor Yellow
        exit 1
    }
}

# Check if Go is installed
Write-Host "🔍 Checking prerequisites..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "✅ Go found: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Error: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "   Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

# Build QuickVM
if (-not $SkipBuild) {
    Write-Host ""
    Write-Host "🔨 Building QuickVM..." -ForegroundColor Yellow
    
    try {
        # Download dependencies
        Write-Host "   📦 Downloading dependencies..." -ForegroundColor Cyan
        go mod download
        
        # Build with optimizations
        Write-Host "   🏗️  Compiling..." -ForegroundColor Cyan
        go build -ldflags="-s -w" -o quickvm.exe
        
        Write-Host "✅ Build completed successfully!" -ForegroundColor Green
    } catch {
        Write-Host "❌ Error: Build failed" -ForegroundColor Red
        Write-Host "   $_" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "⏭️  Skipping build (using existing quickvm.exe)" -ForegroundColor Yellow
}

# Check if binary exists
if (-not (Test-Path "quickvm.exe")) {
    Write-Host "❌ Error: quickvm.exe not found" -ForegroundColor Red
    Write-Host "   Please build the project first or remove -SkipBuild flag" -ForegroundColor Yellow
    exit 1
}

# Install to chosen location
Write-Host ""
Write-Host "📦 Installing QuickVM..." -ForegroundColor Yellow

switch ($InstallLocation) {
    'System' {
        $installPath = "C:\Windows\System32\quickvm.exe"
        Write-Host "   Installing to: $installPath" -ForegroundColor Cyan
        
        try {
            Copy-Item "quickvm.exe" $installPath -Force
            Write-Host "✅ Installed to System32" -ForegroundColor Green
            Write-Host "   You can now use 'quickvm' from anywhere!" -ForegroundColor Green
        } catch {
            Write-Host "❌ Error: Failed to copy to System32" -ForegroundColor Red
            Write-Host "   $_" -ForegroundColor Red
            exit 1
        }
    }
    
    'User' {
        $binDir = "$env:USERPROFILE\bin"
        $installPath = "$binDir\quickvm.exe"
        
        # Create bin directory if it doesn't exist
        if (-not (Test-Path $binDir)) {
            Write-Host "   Creating directory: $binDir" -ForegroundColor Cyan
            New-Item -ItemType Directory -Force -Path $binDir | Out-Null
        }
        
        try {
            Copy-Item "quickvm.exe" $installPath -Force
            Write-Host "✅ Installed to: $installPath" -ForegroundColor Green
            
            # Check if bin directory is in PATH
            $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
            if ($userPath -notlike "*$binDir*") {
                Write-Host ""
                Write-Host "⚠️  Adding $binDir to PATH..." -ForegroundColor Yellow
                
                $newPath = $userPath + ";$binDir"
                [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
                
                Write-Host "✅ Added to PATH" -ForegroundColor Green
                Write-Host "   Please restart your terminal for the changes to take effect" -ForegroundColor Yellow
            } else {
                Write-Host "   $binDir is already in PATH" -ForegroundColor Green
            }
        } catch {
            Write-Host "❌ Error: Installation failed" -ForegroundColor Red
            Write-Host "   $_" -ForegroundColor Red
            exit 1
        }
    }
    
    'Current' {
        Write-Host "   Keeping quickvm.exe in current directory" -ForegroundColor Cyan
        Write-Host "   You can run it with: .\quickvm.exe" -ForegroundColor Green
    }
}

# Create PowerShell alias
if ($CreateAlias) {
    Write-Host ""
    Write-Host "🔗 Creating PowerShell alias..." -ForegroundColor Yellow
    
    $profilePath = $PROFILE
    $aliasContent = @"

# QuickVM Alias
Set-Alias -Name qvm -Value quickvm
"@
    
    try {
        # Create profile if it doesn't exist
        if (-not (Test-Path $profilePath)) {
            $profileDir = Split-Path $profilePath -Parent
            if (-not (Test-Path $profileDir)) {
                New-Item -ItemType Directory -Force -Path $profileDir | Out-Null
            }
            New-Item -ItemType File -Force -Path $profilePath | Out-Null
        }
        
        # Check if alias already exists
        $profileContent = Get-Content $profilePath -Raw -ErrorAction SilentlyContinue
        if ($profileContent -notlike "*QuickVM Alias*") {
            Add-Content -Path $profilePath -Value $aliasContent
            Write-Host "✅ Added 'qvm' alias to PowerShell profile" -ForegroundColor Green
            Write-Host "   You can use 'qvm' instead of 'quickvm'" -ForegroundColor Green
            Write-Host "   Restart your terminal for the alias to take effect" -ForegroundColor Yellow
        } else {
            Write-Host "   Alias already exists in profile" -ForegroundColor Green
        }
    } catch {
        Write-Host "⚠️  Warning: Failed to create alias" -ForegroundColor Yellow
        Write-Host "   $_" -ForegroundColor Yellow
    }
}

# Test installation
Write-Host ""
Write-Host "🧪 Testing installation..." -ForegroundColor Yellow

try {
    $version = & "quickvm" version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ QuickVM is working correctly!" -ForegroundColor Green
    } else {
        throw "QuickVM execution failed"
    }
} catch {
    Write-Host "⚠️  Warning: Could not test installation" -ForegroundColor Yellow
    Write-Host "   Please restart your terminal and try: quickvm version" -ForegroundColor Yellow
}

# Summary
Write-Host ""
Write-Host "╔══════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║          Installation Completed Successfully!        ║" -ForegroundColor Green
Write-Host "╚══════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Host "📚 Quick Start Guide:" -ForegroundColor Cyan
Write-Host "   quickvm           - Launch interactive TUI" -ForegroundColor White
Write-Host "   quickvm list      - List all VMs" -ForegroundColor White
Write-Host "   quickvm start 1   - Start VM #1" -ForegroundColor White
Write-Host "   quickvm stop 1    - Stop VM #1" -ForegroundColor White
Write-Host "   quickvm version   - Show version" -ForegroundColor White
Write-Host ""
Write-Host "📖 For more information, see:" -ForegroundColor Cyan
Write-Host "   - README.md for full documentation" -ForegroundColor White
Write-Host "   - DEMO.md for examples and use cases" -ForegroundColor White
Write-Host ""
Write-Host "🚀 Happy VM managing!" -ForegroundColor Magenta
