# Hướng Dẫn Cài Đặt QuickVM

## 📋 Mục Lục

- [Tổng Quan](#tổng-quan)
- [Phương Pháp 1: Cài Đặt Tương Tác (Khuyến Nghị)](#phương-pháp-1-cài-đặt-tương-tác-khuyến-nghị)
- [Phương Pháp 2: Cài Đặt Tự Động](#phương-pháp-2-cài-đặt-tự-động)
- [Phương Pháp 3: Build Từ Source Code](#phương-pháp-3-build-từ-source-code)
- [So Sánh Các Vị Trí Cài Đặt](#so-sánh-các-vị-trí-cài-đặt)
- [Gỡ Cài Đặt](#gỡ-cài-đặt)
- [Xử Lý Lỗi](#xử-lý-lỗi)

---

## 🎯 Tổng Quan

QuickVM cung cấp **3 phương pháp cài đặt** khác nhau để phù hợp với mọi nhu cầu:

| Phương Pháp | Độ Khó | Khuyến Nghị | Yêu Cầu Admin |
|-------------|---------|-------------|---------------|
| **Cài Đặt Tương Tác** | ⭐ | ✅ Người dùng mới | Tùy chọn |
| **Cài Đặt Tự Động** | ⭐⭐ | 👨‍💻 Power User | Tùy chọn |
| **Build Từ Source** | ⭐⭐⭐ | 👨‍💼 Developer | Tùy chọn |

---

## 📦 Phương Pháp 1: Cài Đặt Tương Tác (Khuyến Nghị)

### Bước 1: Tải QuickVM

1. Truy cập [GitHub Releases](https://github.com/hoangtran1411/quickvm/releases)
2. Tải phiên bản phù hợp với hệ thống của bạn:
   - **Windows AMD64** (Intel/AMD 64-bit) - `quickvm-vX.X.X-windows-amd64.zip`
   - **Windows ARM64** (Surface X, ARM PC) - `quickvm-vX.X.X-windows-arm64.zip`

### Bước 2: Giải Nén File

1. Click phải vào file ZIP đã tải
2. Chọn **"Extract All..."** hoặc **"Giải nén tất cả..."**
3. Chọn thư mục đích (ví dụ: `C:\QuickVM`)

### Bước 3: Chạy Menu Cài Đặt

Có **2 cách** để chạy menu cài đặt:

#### Cách 1: Sử Dụng Batch File (Đơn Giản Nhất)
```
1. Mở thư mục đã giải nén
2. Double-click vào file "install-menu.bat"
```

#### Cách 2: Sử Dụng PowerShell Script
```
1. Mở thư mục đã giải nén
2. Click phải vào "install-menu.ps1"
3. Chọn "Run with PowerShell"
```

### Bước 4: Chọn Vị Trí Cài Đặt

Menu sẽ hiển thị:

```
╔══════════════════════════════════════════════════════╗
║   QuickVM - Fast Hyper-V Virtual Machine Manager    ║
║          Interactive Installation Menu              ║
╚══════════════════════════════════════════════════════╝

Chọn vị trí cài đặt QuickVM:

  1. System   - Cài vào C:\Windows\System32 (cần quyền Admin)
                Sử dụng được từ mọi nơi cho tất cả người dùng

  2. User     - Cài vào ~\bin (không cần quyền Admin)
                Sử dụng được từ mọi nơi cho người dùng hiện tại

  3. Current  - Giữ trong thư mục hiện tại
                Chạy bằng lệnh .\quickvm.exe

  0. Thoát

Nhập lựa chọn của bạn (0-3):
```

**Nhập số tương ứng** (1, 2, hoặc 3) và nhấn Enter.

### Bước 5: Hoàn Tất

- Sau khi cài đặt thành công, terminal sẽ hiển thị thông báo xác nhận
- **Khởi động lại terminal** để sử dụng lệnh `quickvm`
- Kiểm tra cài đặt: `quickvm version`

---

## ⚙️ Phương Pháp 2: Cài Đặt Tự Động

### Sử Dụng PowerShell Script

Thích hợp cho:
- Tự động hóa cài đặt
- Cài đặt hàng loạt trên nhiều máy
- CI/CD pipelines

### Cú Pháp Cơ Bản

```powershell
.\install.ps1 -InstallLocation <Location>
```

### Các Tham Số

| Tham Số | Giá Trị | Mô Tả |
|---------|---------|-------|
| `-InstallLocation` | `System` | Cài vào System32 (cần Admin) |
| | `User` | Cài vào %USERPROFILE%\bin (khuyến nghị) |
| | `Current` | Giữ trong thư mục hiện tại |
| `-SkipBuild` | Switch | Bỏ qua bước build (dùng binary có sẵn) |
| `-CreateAlias` | Switch | Tạo alias 'qvm' cho PowerShell |

### Ví Dụ

#### Cài đặt cho người dùng hiện tại
```powershell
.\install.ps1 -InstallLocation User
```

#### Cài đặt system-wide với alias
```powershell
# Chạy PowerShell as Administrator
.\install.ps1 -InstallLocation System -CreateAlias
```

#### Cài đặt portable mode
```powershell
.\install.ps1 -InstallLocation Current
```

#### Cài đặt từ binary có sẵn
```powershell
.\install.ps1 -InstallLocation User -SkipBuild
```

---

## 🔨 Phương Pháp 3: Build Từ Source Code

### Yêu Cầu

- **Git** - [Tải Git](https://git-scm.com/download/win)
- **Go 1.21+** - [Tải Go](https://golang.org/dl/)

### Các Bước

#### 1. Clone Repository

```bash
git clone https://github.com/hoangtran1411/quickvm.git
cd quickvm
```

#### 2. Download Dependencies

```bash
go mod download
go mod verify
```

#### 3. Build

```bash
# Build cơ bản
go build -o quickvm.exe

# Build tối ưu (giảm kích thước)
go build -ldflags="-s -w" -o quickvm.exe

# Build với thông tin version
go build -ldflags="-s -w -X 'main.Version=1.0.0'" -o quickvm.exe
```

#### 4. Cài Đặt

```bash
# Sử dụng menu tương tác
.\install-menu.bat

# Hoặc sử dụng script trực tiếp
.\install.ps1 -InstallLocation User
```

---

## 🔍 So Sánh Các Vị Trí Cài Đặt

### Option 1: System (`C:\Windows\System32`)

**✅ Ưu điểm:**
- Có thể sử dụng từ bất kỳ đâu trên hệ thống
- Tất cả người dùng đều có thể dùng
- Không cần cấu hình PATH

**❌ Nhược điểm:**
- **Yêu cầu quyền Administrator**
- Khó cập nhật (cần Admin mỗi lần)
- Có thể bị Windows Defender cảnh báo

**🎯 Phù hợp với:**
- Server/Workstation dùng chung
- IT Administrator

### Option 2: User (`%USERPROFILE%\bin`)

**✅ Ưu điểm:**
- **Không cần quyền Administrator**
- Dễ dàng cập nhật
- Tự động thêm vào PATH
- An toàn hơn

**❌ Nhược điểm:**
- Chỉ người dùng hiện tại sử dụng được
- Cần restart terminal lần đầu

**🎯 Phù hợp với:** (⭐ **KHUYẾN NGHỊ**)
- Máy cá nhân
- Hầu hết người dùng
- Development environment

### Option 3: Current (Thư mục hiện tại)

**✅ Ưu điểm:**
- Portable - di chuyển được
- Không cần cài đặt
- Không ảnh hưởng hệ thống

**❌ Nhược điểm:**
- Phải chạy từ thư mục chứa file
- Cần gõ đường dẫn đầy đủ: `.\quickvm.exe`

**🎯 Phù hợp với:**
- USB drive / Portable tools
- Testing
- Temporary use

---

## 🗑️ Gỡ Cài Đặt

### Nếu cài ở System

```powershell
# Chạy PowerShell as Administrator
Remove-Item C:\Windows\System32\quickvm.exe
```

### Nếu cài ở User

```powershell
Remove-Item $env:USERPROFILE\bin\quickvm.exe

# (Optional) Xóa thư mục bin nếu rỗng
Remove-Item $env:USERPROFILE\bin -Force
```

### Nếu cài ở Current

```powershell
# Chỉ cần xóa file trong thư mục
Remove-Item quickvm.exe
```

### Xóa Alias (nếu đã tạo)

```powershell
# Mở PowerShell profile
notepad $PROFILE

# Xóa các dòng có chứa "QuickVM Alias"
# Lưu và đóng file
```

---

## 🐛 Xử Lý Lỗi

### Lỗi: "Execution policy"

**Triệu chứng:**
```
install-menu.ps1 cannot be loaded because running scripts is disabled
```

**Giải pháp:**
```powershell
# Tạm thời cho phép chạy script
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process

# Sau đó chạy lại
.\install-menu.ps1
```

### Lỗi: "Go is not installed"

**Triệu chứng:**
```
❌ Error: Go is not installed or not in PATH
```

**Giải pháp:**
1. Tải Go từ https://golang.org/dl/
2. Cài đặt và restart terminal
3. Kiểm tra: `go version`

### Lỗi: "Access is denied" (khi cài System)

**Triệu chứng:**
```
❌ Error: Failed to copy to System32
Access is denied
```

**Giải pháp:**
```powershell
# Chạy PowerShell as Administrator
# Click phải PowerShell → "Run as Administrator"
# Sau đó chạy lại script
```

### Lỗi: "quickvm is not recognized"

**Triệu chứng:**
```
'quickvm' is not recognized as an internal or external command
```

**Giải pháp:**
1. **Restart terminal** (quan trọng!)
2. Kiểm tra PATH:
   ```powershell
   $env:Path -split ';' | Select-String "bin"
   ```
3. Nếu không thấy, chạy lại cài đặt với option User

### Script bị block bởi Windows Security

**Triệu chứng:**
- File bị gắn cờ "Unblock"

**Giải pháp:**
```powershell
# Click phải file → Properties → Unblock
# Hoặc dùng PowerShell:
Unblock-File .\install-menu.ps1
Unblock-File .\install-menu.bat
Unblock-File .\install.ps1
```

---

## 📞 Hỗ Trợ

Nếu gặp vấn đề không giải quyết được:

1. Kiểm tra [Issues](https://github.com/hoangtran1411/quickvm/issues) trên GitHub
2. Tạo Issue mới với thông tin:
   - Hệ điều hành (Windows version)
   - Phương pháp cài đặt đã sử dụng
   - Thông báo lỗi đầy đủ
   - Screenshot (nếu có)

---

## ✅ Checklist Sau Khi Cài Đặt

- [ ] Restart terminal
- [ ] Chạy `quickvm version` để kiểm tra
- [ ] Chạy `quickvm list` để xem danh sách VM
- [ ] (Optional) Tạo alias `qvm` bằng cách chạy lại với flag `-CreateAlias`

---

**🎉 Chúc mừng! Bạn đã cài đặt QuickVM thành công!**

Xem [Quick Reference](QUICK_REFERENCE.md) để biết cách sử dụng.
