# GPU Passthrough (GPU-P) Feature

Share GPU between host and Hyper-V virtual machines through GPU Partitioning.

## System Requirements

- Windows 10/11 Pro with Hyper-V enabled
- GPU supporting GPU Partitioning (newer NVIDIA/AMD generations)
- GPU drivers installed on Host

## Check GPU Support

```powershell
# PowerShell as Admin
Get-VMPartitionableGPU
```

If there's output, your GPU supports partitioning.

## Commands

### Check GPU Status
```bash
quickvm gpu status
```

### Add GPU to VM
```bash
# VM must be OFF before adding GPU
quickvm stop <vm-index>
quickvm gpu add <vm-index>
```

### Remove GPU from VM
```bash
quickvm gpu remove <vm-index>
```

## Installing Drivers in Guest VM

After adding GPU, you need to copy drivers from Host to Guest:

### 1. Copy Driver Files

**From Host:**
```
C:\Windows\System32\DriverStore\FileRepository\nv_dispi.inf_amd64_[UUID]
```

**To Guest:**
```
C:\Windows\System32\HostDriverStore\FileRepository\nv_dispi.inf_amd64_[UUID]
```

### 2. Copy System Files

**From Host:**
```
C:\Windows\System32\nv*.*  (all files starting with "nv")
```

**To Guest:**
```
C:\Windows\System32\
```

## PowerShell Commands (Detailed)

The original script uses the following commands:

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

| Issue | Solution |
|-------|----------|
| GPU not visible in Guest | Check if drivers were copied correctly |
| Error when adding GPU | Ensure VM is OFF |
| `Get-VMPartitionableGPU` returns empty | GPU not supported or outdated driver |
| Admin required | Run PowerShell/Terminal as Administrator |

## References

- [GPU-P Tutorial.txt](../GPU-P%20Tutorial.txt)
- [GPU-P-Partition.ps1](../GPU-P-Partition.ps1)
