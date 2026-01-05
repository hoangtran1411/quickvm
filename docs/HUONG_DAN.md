# 🚀 QuickVM - Hướng Dẫn Nhanh

## Cài Đặt

### Yêu Cầu
- Windows 10/11 với Hyper-V đã bật
- Quyền Administrator
- Go 1.21+ (nếu build từ source)

### Build Từ Source

```powershell
# Clone repository
git clone <your-repo-url>
cd quickvm

# Tải dependencies
go mod download

# Build ứng dụng
go build -o quickvm.exe

# (Tùy chọn) Copy vào thư mục trong PATH
Copy-Item quickvm.exe C:\Windows\System32\
```

## 📚 Các Lệnh

### 1. Xem Danh Sách VM
```powershell
# Liệt kê tất cả máy ảo
quickvm list

# Hoặc dùng alias
quickvm ls
```

### 2. Khởi Động VM
```powershell
# Khởi động VM theo index (ví dụ: VM số 1)
quickvm start 1
```

### 3. Dừng VM
```powershell
# Dừng VM theo index
quickvm stop 1
```

### 4. Khởi Động Lại VM
```powershell
# Restart VM theo index
quickvm restart 1
```

### 5. Giao Diện TUI (Interactive)
```powershell
# Chạy chế độ TUI
quickvm

# Trong TUI:
# ↑/↓    - Di chuyển giữa các VM
# Enter  - Khởi động VM được chọn
# s      - Dừng VM được chọn
# t      - Restart VM được chọn
# r      - Refresh danh sách
# q/Esc  - Thoát
```

## 💡 Ví Dụ Sử Dụng

### Workflow Thông Thường

```powershell
# 1. Xem danh sách VM
PS> quickvm list

📋 Fetching Hyper-V virtual machines...

==============================================================================
Index   Name             State        CPU%    Memory(MB)  Uptime    Status
==============================================================================
1       Ubuntu-Dev       🔴 Off      0%      0           00:00:00  Operating normally
2       Windows-Test     🟢 Running  5%      4096        02:15:30  Operating normally
3       Docker-Host      🔴 Off      0%      0           00:00:00  Operating normally
==============================================================================

Total VMs: 3

💡 Tip: Use 'quickvm start <index>' to start a VM

# 2. Khởi động VM thứ 1
PS> quickvm start 1
🚀 Starting VM: Ubuntu-Dev (Index: 1)...
✅ VM 'Ubuntu-Dev' started successfully!

# 3. Dừng VM thứ 2
PS> quickvm stop 2
🛑 Stopping VM: Windows-Test (Index: 2)...
✅ VM 'Windows-Test' stopped successfully!
```

## 🎯 Tips & Tricks

### 1. Thêm Vào PATH
Để sử dụng `quickvm` từ bất kỳ đâu:
```powershell
# Copy executable vào thư mục System32
Copy-Item quickvm.exe C:\Windows\System32\
```

### 2. Tạo Alias PowerShell
Thêm vào PowerShell profile của bạn (`$PROFILE`):
```powershell
# Mở profile
notepad $PROFILE

# Thêm các alias
Set-Alias qvm "D:\path\to\quickvm.exe"
```

Sau đó bạn có thể dùng:
```powershell
qvm list
qvm start 1
```

### 3. Batch Operations
```powershell
# Khởi động nhiều VM
quickvm start 1
quickvm start 2
quickvm start 3

# Hoặc dùng loop
1..3 | ForEach-Object { quickvm start $_ }
```

## ⚠️ Lưu Ý

1. **Quyền Administrator**: Luôn chạy PowerShell/CMD với quyền Administrator
2. **Index VM**: Index của VM có thể thay đổi khi bạn thêm/xóa VM. Chạy `quickvm list` để xem index mới nhất
3. **Trạng Thái VM**: 
   - 🟢 Running - VM đang chạy
   - 🔴 Off - VM đã tắt
   - 🟡 Paused - VM đang tạm dừng

## 🐛 Troubleshooting

### Lỗi: "Failed to get VMs"
- Đảm bảo Hyper-V đã được bật
- Chạy với quyền Administrator
- Kiểm tra xem bạn có VM nào không: `Get-VM` trong PowerShell

### Lỗi: "Failed to start VM"
- VM có thể đã đang chạy
- Kiểm tra resource (CPU, RAM) còn đủ không
- Xem logs trong Event Viewer

### VM Không Hiển Thị
- Refresh lại: nhấn `r` trong TUI mode
- Hoặc chạy lại `quickvm list`

## 📞 Hỗ Trợ

Nếu gặp vấn đề, hãy:
1. Kiểm tra phần Troubleshooting ở trên
2. Xem logs PowerShell: `Get-VM -Name "VM-Name" | Format-List *`
3. Mở issue trên GitHub

## 🎓 Học Thêm

- [Hyper-V Documentation](https://docs.microsoft.com/en-us/virtualization/hyper-v-on-windows/)
- [PowerShell Get-VM Cmdlet](https://docs.microsoft.com/en-us/powershell/module/hyper-v/get-vm)

---

**Chúc bạn sử dụng QuickVM hiệu quả! 🚀**
