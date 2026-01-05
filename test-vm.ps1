Get-VM | Select-Object Name, State, CPUUsage, `
    @{Name='MemoryMB';Expression={[int]($_.MemoryAssigned/1MB)}}, `
    @{Name='Uptime';Expression={$_.Uptime.ToString()}}, `
    @{Name='Status';Expression={$_.Status.ToString()}}, `
    @{Name='Version';Expression={$_.Version.ToString()}} | ConvertTo-Json
