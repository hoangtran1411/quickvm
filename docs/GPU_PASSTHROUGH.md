# GPU Passthrough (GPU-P) Feature

Tính năng share GPU giữa máy thật và máy ảo Hyper-V thông qua GPU Partitioning.

## Yêu Cầu Hệ Thống

- Windows 10/11 Pro với Hyper-V enabled
- GPU hỗ trợ GPU Partitioning (NVIDIA/AMD thế hệ mới)
- GPU drivers đã cài đặt trên Host

## Kiểm Tra GPU Hỗ Trợ

```powershell
# PowerShell as Admin
Get-VMPartitionableGPU
```

Nếu có output, GPU của bạn hỗ trợ partitioning.

## Commands

### Kiểm tra GPU status
```bash
quickvm gpu status
```

### Add GPU cho VM
```bash
# VM phải OFF trước khi add GPU
quickvm stop <vm-index>
quickvm gpu add <vm-index>
```

### Remove GPU khỏi VM
```bash
quickvm gpu remove <vm-index>
```

## Cài Đặt Driver Trong Guest VM

Sau khi add GPU, cần copy drivers từ Host sang Guest:

### 1. Copy Driver Files

**Từ Host:**
```
C:\Windows\System32\DriverStore\FileRepository\nv_dispi.inf_amd64_[UUID]
```

**Sang Guest:**
```
C:\Windows\System32\HostDriverStore\FileRepository\nv_dispi.inf_amd64_[UUID]
```

### 2. Copy System Files

**Từ Host:**
```
C:\Windows\System32\nv*.*  (tất cả file bắt đầu bằng "nv")
```

**Sang Guest:**
```
C:\Windows\System32\
```

## PowerShell Commands (Chi tiết)

Script gốc sử dụng các commands sau:

```powershell
$vm = "VMName"

# 1. Add GPU Partition Adapter
Add-VMGpuPartitionAdapter -VMName $vm

# 2. Configure GPU Partition (VRAM, Encode, Decode, Compute)
Set-VMGpuPartitionAdapter -VMName $vm `
  -MinPartitionVRAM 80000000 -MaxPartitionVRAM 100000000 -OptimalPartitionVRAM 100000000 `
  -MinPartitionEncode 80000000 -MaxPartitionEncode 100000000 -OptimalPartitionEncode 100000000 `
  -MinPartitionDecode 80000000 -MaxPartitionDecode 100000000 -OptimalPartitionDecode 100000000 `
  -MinPartitionCompute 80000000 -MaxPartitionCompute 100000000 -OptimalPartitionCompute 100000000

# 3. Enable Guest Controlled Cache Types
Set-VM -GuestControlledCacheTypes $true -VMName $vm

# 4. Set Memory Mapped IO Space
Set-VM -LowMemoryMappedIoSpace 1Gb -VMName $vm
Set-VM -HighMemoryMappedIoSpace 32GB -VMName $vm
```

## Troubleshooting

| Vấn đề | Giải pháp |
|--------|-----------|
| GPU không hiện trong Guest | Kiểm tra driver đã copy đúng chưa |
| Lỗi khi add GPU | Đảm bảo VM đang OFF |
| `Get-VMPartitionableGPU` trống | GPU không hỗ trợ hoặc driver cũ |
| Cần Admin | Chạy PowerShell/Terminal as Administrator |

## Tham Khảo

- [GPU-P Tutorial.txt](../GPU-P%20Tutorial.txt)
- [GPU-P-Partition.ps1](../GPU-P-Partition.ps1)
