# QuickVM Installation Menu
# Interactive menu for installing QuickVM with different locations

$ErrorActionPreference = "Stop"

# Display banner
Write-Host ""
Write-Host "╔══════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   QuickVM - Fast Hyper-V Virtual Machine Manager     ║" -ForegroundColor Cyan
Write-Host "║          Interactive Installation Menu               ║" -ForegroundColor Cyan
Write-Host "╚══════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Display menu
Write-Host "Chọn vị trí cài đặt QuickVM:" -ForegroundColor Yellow
Write-Host ""
Write-Host "  1. System   - Cài vào C:\Windows\System32 (cần quyền Admin)" -ForegroundColor White
Write-Host "                Sử dụng được từ mọi nơi cho tất cả người dùng" -ForegroundColor Gray
Write-Host ""
Write-Host "  2. User     - Cài vào ~\bin (không cần quyền Admin)" -ForegroundColor White
Write-Host "                Sử dụng được từ mọi nơi cho người dùng hiện tại" -ForegroundColor Gray
Write-Host ""
Write-Host "  3. Current  - Giữ trong thư mục hiện tại" -ForegroundColor White
Write-Host "                Chạy bằng lệnh .\quickvm.exe" -ForegroundColor Gray
Write-Host ""
Write-Host "  0. Thoát" -ForegroundColor Red
Write-Host ""

# Get user choice
do {
    $choice = Read-Host "Nhập lựa chọn của bạn (0-3)"
    
    switch ($choice) {
        "1" {
            Write-Host ""
            Write-Host "⚙️  Đang cài đặt vào System..." -ForegroundColor Cyan
            & ".\install.ps1" -InstallLocation System
            break
        }
        "2" {
            Write-Host ""
            Write-Host "⚙️  Đang cài đặt vào User..." -ForegroundColor Cyan
            & ".\install.ps1" -InstallLocation User
            break
        }
        "3" {
            Write-Host ""
            Write-Host "⚙️  Đang cài đặt vào Current..." -ForegroundColor Cyan
            & ".\install.ps1" -InstallLocation Current
            break
        }
        "0" {
            Write-Host ""
            Write-Host "👋 Đã hủy cài đặt!" -ForegroundColor Yellow
            exit 0
        }
        default {
            Write-Host ""
            Write-Host "❌ Lựa chọn không hợp lệ! Vui lòng chọn từ 0-3" -ForegroundColor Red
            Write-Host ""
            $choice = $null
        }
    }
} while ($null -eq $choice)

Write-Host ""
Write-Host "Nhấn phím bất kỳ để đóng..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
