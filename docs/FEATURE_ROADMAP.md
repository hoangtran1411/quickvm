# QuickVM Feature Roadmap 🗺️

> Tài liệu này mô tả các feature được đề xuất cho QuickVM, sắp xếp theo mức độ ưu tiên và độ phức tạp.

**Cập nhật lần cuối:** 2026-01-07

---

## 📊 Tổng quan

| Tier | Mô tả | Số lượng Features |
|------|-------|-------------------|
| Tier 1 | High Value, Medium Effort | 4 |
| Tier 2 | Quick Wins | 5 |
| Tier 3 | Advanced Features | 5 |
| Tier 4 | Nice to Have | 5 |

---

## 🚀 Tier 1: High Value, Medium Effort

> **Ưu tiên cao** - Các feature mang lại giá trị lớn, nên triển khai sớm.

### 1. VM Snapshots/Checkpoints ⭐ ✅ DONE

**Command:** `quickvm snapshot`

```bash
quickvm snapshot list <vm-index>              # Liệt kê snapshots của VM
quickvm snapshot create <vm-index> "name"     # Tạo snapshot mới
quickvm snapshot restore <vm-index> "name"    # Khôi phục snapshot
quickvm snapshot delete <vm-index> "name"     # Xóa snapshot
```

**Lý do:** Quản lý checkpoint là tính năng rất quan trọng khi làm việc với VM. Cho phép người dùng:
- Lưu trạng thái trước khi thực hiện thay đổi
- Khôi phục nhanh khi có lỗi
- Test safely với khả năng rollback

**Độ phức tạp:** ⭐⭐⭐ (Medium)

**PowerShell Commands:**
```powershell
Get-VMSnapshot -VMName "VMName"
Checkpoint-VM -Name "VMName" -SnapshotName "SnapshotName"
Restore-VMSnapshot -VMName "VMName" -Name "SnapshotName" -Confirm:$false
Remove-VMSnapshot -VMName "VMName" -Name "SnapshotName"
```

---

### 2. VM Clone

**Command:** `quickvm clone`

```bash
quickvm clone <vm-index> <new-name>           # Clone VM với tên mới
quickvm clone 1 "WebServer-Test"              # Ví dụ
```

**Lý do:** Clone VM nhanh để test/dev mà không ảnh hưởng VM gốc.

**Độ phức tạp:** ⭐⭐⭐ (Medium)

**PowerShell Commands:**
```powershell
Export-VM -Name "SourceVM" -Path "C:\VMs\Export"
Import-VM -Path "C:\VMs\Export\SourceVM\Virtual Machines\*.vmcx" -Copy -GenerateNewId
Rename-VM -Name "SourceVM" -NewName "NewVMName"
```

---

### 3. Export/Import VM

**Command:** `quickvm export` / `quickvm import`

```bash
quickvm export <vm-index> <path>              # Export VM ra thư mục
quickvm export 1 "D:\Backups\VMs"

quickvm import <path>                          # Import VM từ file
quickvm import "D:\Backups\VMs\WebServer"
```

**Lý do:** Backup và di chuyển VM giữa các máy.

**Độ phức tạp:** ⭐⭐⭐ (Medium)

**PowerShell Commands:**
```powershell
Export-VM -Name "VMName" -Path "D:\Backups"
Import-VM -Path "D:\Backups\VMName\Virtual Machines\*.vmcx"
```

---

### 4. VM Config

**Command:** `quickvm config`

```bash
quickvm config <vm-index> --memory 4GB        # Thay đổi RAM
quickvm config <vm-index> --cpu 2             # Thay đổi số CPU
quickvm config <vm-index> --memory 8GB --cpu 4  # Cả hai
quickvm config show <vm-index>                # Xem config hiện tại
```

**Lý do:** Thay đổi RAM/CPU của VM mà không cần mở Hyper-V Manager.

**Độ phức tạp:** ⭐⭐ (Low-Medium)

**PowerShell Commands:**
```powershell
Set-VM -Name "VMName" -MemoryStartupBytes 4GB
Set-VMProcessor -VMName "VMName" -Count 2
Get-VM -Name "VMName" | Select-Object *
```

---

## ⚡ Tier 2: Quick Wins

> **Làm nhanh, giá trị cao** - Các feature đơn giản nhưng hữu ích.

### 5. Connect to VM ⭐

**Command:** `quickvm connect`

```bash
quickvm connect <vm-index>                    # Mở VM Connect GUI
quickvm connect 1
```

**Lý do:** Mở VMConnect.exe trực tiếp từ terminal, không cần mở Hyper-V Manager.

**Độ phức tạp:** ⭐ (Low)

**Implementation:**
```go
// Đơn giản: gọi vmconnect.exe
exec.Command("vmconnect.exe", "localhost", vmName).Start()
```

---

### 6. SSH/RDP Quick Connect

**Command:** `quickvm ssh` / `quickvm rdp`

```bash
quickvm ssh <vm-index>                        # SSH vào VM (Linux)
quickvm ssh 1 -u admin                        # Với username

quickvm rdp <vm-index>                        # RDP vào VM (Windows)
quickvm rdp 1
```

**Lý do:** Kết nối nhanh vào VM nếu biết IP address.

**Độ phức tạp:** ⭐⭐ (Low-Medium)

**Yêu cầu:** Cần lấy IP của VM trước:
```powershell
(Get-VMNetworkAdapter -VMName "VMName").IPAddresses
```

---

### 7. VM Logs

**Command:** `quickvm logs`

```bash
quickvm logs <vm-index>                       # Xem event logs của VM
quickvm logs 1 --tail 50                      # Chỉ 50 dòng cuối
quickvm logs 1 --follow                       # Follow mode (real-time)
```

**Lý do:** Debug và troubleshoot VM issues.

**Độ phức tạp:** ⭐⭐ (Low-Medium)

---

### 8. Bulk Operations ⭐

**Command:** `quickvm start/stop/restart --all`

```bash
quickvm start --all                           # Start tất cả VMs
quickvm stop --all                            # Stop tất cả VMs
quickvm restart --all                         # Restart tất cả VMs

quickvm start --filter "Running"              # Start VMs đang Running
quickvm stop --filter "Web*"                  # Stop VMs có tên bắt đầu bằng "Web"
```

**Lý do:** Quản lý nhiều VMs cùng lúc.

**Độ phức tạp:** ⭐ (Low)

---

### 9. Watch Mode

**Command:** `quickvm watch`

```bash
quickvm watch                                 # Real-time monitoring TUI
quickvm watch --interval 5                    # Refresh mỗi 5 giây
quickvm list --watch                          # Watch mode cho list command
```

**Lý do:** Theo dõi trạng thái VMs real-time, đặc biệt hữu ích khi waiting cho VM start/stop.

**Độ phức tạp:** ⭐⭐ (Low-Medium)

---

## 🔧 Tier 3: Advanced Features

> **Nâng cao** - Các feature phức tạp hơn, dành cho power users.

### 10. VM Templates

**Command:** `quickvm template`

```bash
quickvm template create <vm-index> "TemplateName"   # Tạo template từ VM
quickvm template list                               # Liệt kê templates
quickvm template apply "TemplateName" "NewVMName"   # Tạo VM từ template
quickvm template delete "TemplateName"              # Xóa template
```

**Lý do:** Tạo VMs mới nhanh từ template đã chuẩn bị sẵn.

**Độ phức tạp:** ⭐⭐⭐⭐ (High)

---

### 11. Network Management

**Command:** `quickvm network`

```bash
quickvm network list                          # Liệt kê Virtual Switches
quickvm network create "SwitchName" --type internal
quickvm network attach <vm-index> "SwitchName"
quickvm network detach <vm-index>
```

**Lý do:** Quản lý Virtual Switch và network cho VMs.

**Độ phức tạp:** ⭐⭐⭐ (Medium)

---

### 12. Storage Management

**Command:** `quickvm disk`

```bash
quickvm disk list <vm-index>                  # Liệt kê disks của VM
quickvm disk create "disk.vhdx" --size 50GB   # Tạo VHD mới
quickvm disk resize "disk.vhdx" --size 100GB  # Resize VHD
quickvm disk attach <vm-index> "disk.vhdx"    # Attach disk vào VM
quickvm disk detach <vm-index> "disk.vhdx"    # Detach disk
```

**Lý do:** Quản lý VHD/VHDX files.

**Độ phức tạp:** ⭐⭐⭐ (Medium)

---

### 13. Resource Quotas

**Command:** `quickvm quota`

```bash
quickvm quota set <vm-index> --max-cpu 50%    # Giới hạn CPU
quickvm quota set <vm-index> --max-memory 4GB # Giới hạn RAM
quickvm quota show <vm-index>                 # Xem quotas hiện tại
```

**Lý do:** Set giới hạn resource, hữu ích cho lab environments.

**Độ phức tạp:** ⭐⭐⭐ (Medium)

---

### 14. Scheduled Tasks

**Command:** `quickvm schedule`

```bash
quickvm schedule start <vm-index> --at "08:00"        # Start VM lúc 8h sáng
quickvm schedule stop <vm-index> --at "18:00"         # Stop VM lúc 6h tối
quickvm schedule list                                  # Xem schedules
quickvm schedule delete <schedule-id>                  # Xóa schedule
```

**Lý do:** Tự động start/stop VMs theo lịch.

**Độ phức tạp:** ⭐⭐⭐⭐ (High)

---

## 🎯 Tier 4: Nice to Have

> **Tương lai** - Các feature bổ sung khi có thời gian.

### 15. Profile/Workspace

```bash
quickvm workspace create "Development"        # Tạo workspace
quickvm workspace add 1 2 3                   # Thêm VMs vào workspace
quickvm workspace start "Development"         # Start tất cả VMs trong workspace
```

**Lý do:** Nhóm VMs theo project/mục đích.

---

### 16. Remote Host Management

```bash
quickvm remote add "server1" --host 192.168.1.100
quickvm remote list
quickvm --host server1 list                   # Quản lý VMs trên máy khác
```

**Lý do:** Quản lý Hyper-V trên các máy khác (remote management).

---

### 17. Metrics Export

```bash
quickvm metrics export --format prometheus    # Export metrics
quickvm metrics serve --port 9090             # HTTP endpoint cho metrics
```

**Lý do:** Integration với monitoring tools (Prometheus/Grafana).

---

### 18. Configuration File

**File:** `~/.quickvmrc` hoặc `quickvm.yaml`

```yaml
# quickvm.yaml
defaults:
  memory: 4GB
  cpu: 2
  
aliases:
  web: 1
  db: 2
  
autostart:
  - web
  - db
```

**Lý do:** Lưu settings và preferences.

---

### 19. Plugin System

```bash
quickvm plugin install quickvm-docker         # Cài plugin
quickvm plugin list                           # Liệt kê plugins
quickvm docker ps                             # Command từ plugin
```

**Lý do:** Extensible architecture cho custom commands.

---

## 📋 Implementation Priority

### Phase 1 (Tuần 1-2)
- [ ] VM Connect (Tier 2, #5)
- [ ] Bulk Operations (Tier 2, #8)
- [x] VM Snapshots (Tier 1, #1) ✅ **Completed 2026-01-07**

### Phase 2 (Tuần 3-4)
- [ ] VM Config (Tier 1, #4)
- [ ] Watch Mode (Tier 2, #9)
- [ ] VM Logs (Tier 2, #7)

### Phase 3 (Tuần 5-6)
- [ ] VM Clone (Tier 1, #2)
- [ ] Export/Import (Tier 1, #3)
- [ ] SSH/RDP Connect (Tier 2, #6)

### Phase 4 (Tương lai)
- [ ] Tier 3 & 4 features

---

## 🔗 Related Documents

- [README.md](../README.md) - Project overview
- [DEVELOPER.md](DEVELOPER.md) - Developer guide
- [CONTRIBUTING.md](CONTRIBUTING.md) - How to contribute

---

> 💡 **Ghi chú:** Đây là roadmap dự kiến và có thể thay đổi dựa trên feedback từ users.
