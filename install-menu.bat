@echo off
:: QuickVM Installation Menu Launcher
:: This batch file launches the PowerShell installation menu

title QuickVM Installation

:: Check if install-menu.ps1 exists
if not exist "%~dp0install-menu.ps1" (
    echo Error: install-menu.ps1 not found!
    echo Please make sure install-menu.ps1 is in the same directory.
    pause
    exit /b 1
)

:: Run PowerShell script
powershell.exe -ExecutionPolicy Bypass -File "%~dp0install-menu.ps1"

:: Keep window open if there's an error
if errorlevel 1 (
    echo.
    echo An error occurred during installation.
    pause
)
