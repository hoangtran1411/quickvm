package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"quickvm/hyperv"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display system information",
	Long: `Display system information including:
- CPU name and cores
- RAM (MB and GB)
- Disk drives with free/total space
- Hyper-V status`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := hyperv.NewManager()
		info, err := manager.GetSystemInfo()
		if err != nil {
			color.Red("❌ Error getting system info: %v", err)
			return
		}

		printSystemInfo(info)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func printSystemInfo(info *hyperv.SystemInfo) {
	// Header
	headerColor := color.New(color.FgCyan, color.Bold)
	labelColor := color.New(color.FgYellow)
	valueColor := color.New(color.FgWhite)
	successColor := color.New(color.FgGreen)
	errorColor := color.New(color.FgRed)

	fmt.Println()
	headerColor.Println("╔══════════════════════════════════════════════════════════════╗")
	headerColor.Println("║                       SYSTEM INFORMATION                     ║")
	headerColor.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// CPU Section
	headerColor.Println("🖥️  CPU")
	fmt.Println(strings.Repeat("─", 60))
	labelColor.Print("   Name:  ")
	valueColor.Println(info.CPU.Name)
	labelColor.Print("   Cores: ")
	valueColor.Printf("%d cores\n", info.CPU.Cores)
	fmt.Println()

	// Memory Section
	headerColor.Println("💾 MEMORY (RAM)")
	fmt.Println(strings.Repeat("─", 60))
	labelColor.Print("   Total: ")
	valueColor.Printf("%d MB (%.2f GB)\n", info.Memory.TotalMB, info.Memory.TotalGB)
	labelColor.Print("   Used:  ")
	valueColor.Printf("%d MB (%.2f GB)\n", info.Memory.UsedMB, info.Memory.UsedGB)
	labelColor.Print("   Free:  ")
	successColor.Printf("%d MB (%.2f GB)\n", info.Memory.FreeMB, info.Memory.FreeGB)

	// Memory progress bar
	usedPercent := float64(info.Memory.UsedMB) / float64(info.Memory.TotalMB) * 100
	printProgressBar("   Usage", usedPercent, 40)
	fmt.Println()

	// Disk Section
	headerColor.Println("💿 DISK DRIVES")
	fmt.Println(strings.Repeat("─", 60))
	for _, disk := range info.Disks {
		labelColor.Printf("   Drive %s\n", disk.Name)
		valueColor.Printf("      Total: %d MB (%.2f GB)\n", disk.TotalMB, disk.TotalGB)
		valueColor.Printf("      Used:  %d MB (%.2f GB)\n", disk.UsedMB, disk.UsedGB)
		successColor.Printf("      Free:  %d MB (%.2f GB)\n", disk.FreeMB, disk.FreeGB)
		
		diskUsedPercent := float64(disk.UsedMB) / float64(disk.TotalMB) * 100
		printProgressBar("      Usage", diskUsedPercent, 35)
		fmt.Println()
	}

	// Hyper-V Section
	headerColor.Println("🔧 HYPER-V STATUS")
	fmt.Println(strings.Repeat("─", 60))
	labelColor.Print("   Status:  ")
	if info.HyperV.Enabled {
		successColor.Printf("✅ %s (Enabled)\n", info.HyperV.Status)
	} else {
		errorColor.Printf("❌ %s (Disabled)\n", info.HyperV.Status)
	}
	fmt.Println()

	headerColor.Println("════════════════════════════════════════════════════════════════")
}

// printProgressBar prints a visual progress bar
func printProgressBar(label string, percent float64, width int) {
	filled := int(percent / 100 * float64(width))
	empty := width - filled

	// Choose color based on usage
	var barColor *color.Color
	if percent >= 90 {
		barColor = color.New(color.FgRed)
	} else if percent >= 70 {
		barColor = color.New(color.FgYellow)
	} else {
		barColor = color.New(color.FgGreen)
	}

	fmt.Printf("%s: [", label)
	barColor.Print(strings.Repeat("█", filled))
	fmt.Print(strings.Repeat("░", empty))
	fmt.Printf("] %.1f%%\n", percent)
}
