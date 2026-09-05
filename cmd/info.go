package cmd

import (
	"fmt"
	"os"
	"strings"

	"quickvm/internal/hyperv"
	"quickvm/internal/output"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display system information",
	Long: `Display system information including:
- CPU name and cores
- RAM (MB and GB)
- Disk drives with free/total space
- Hyper-V status`,
	Run: func(cmd *cobra.Command, _ []string) {
		manager := hyperv.NewManager()

		includeDisk, _ := cmd.Flags().GetBool("disk")

		info, err := manager.GetSystemInfo(cmd.Context(), includeDisk)
		if err != nil {
			output.PrintError("SYSTEM_INFO_FAILED", "Error getting system info", err.Error())
			if !output.IsJSON() {
				color.Red("❌ Error getting system info: %v", err)
			}
			os.Exit(1)
		}

		// JSON output for AI agents
		if output.IsJSON() {
			output.PrintData(info)
			return
		}

		printSystemInfo(info)
	},
}

func init() {
	infoCmd.Flags().BoolP("disk", "d", false, "Include disk usage information (slower)")
	rootCmd.AddCommand(infoCmd)
}

//nolint:funlen // UI function
func printSystemInfo(info *hyperv.SystemInfo) {
	// Header
	headerColor := color.New(color.FgCyan, color.Bold)
	labelColor := color.New(color.FgYellow)
	valueColor := color.New(color.FgWhite)
	successColor := color.New(color.FgGreen)
	errorColor := color.New(color.FgRed)

	fmt.Println()
	_, _ = headerColor.Println("╔══════════════════════════════════════════════════════════════╗")
	_, _ = headerColor.Println("║                       SYSTEM INFORMATION                     ║")
	_, _ = headerColor.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// CPU Section
	_, _ = headerColor.Println("🖥️  CPU")
	fmt.Println(strings.Repeat("─", 60))
	_, _ = labelColor.Print("   Name:  ")
	_, _ = valueColor.Println(info.CPU.Name)
	_, _ = labelColor.Print("   Cores: ")
	_, _ = valueColor.Printf("%d cores\n", info.CPU.Cores)
	fmt.Println()

	// Memory Section
	_, _ = headerColor.Println("💾 MEMORY (RAM)")
	fmt.Println(strings.Repeat("─", 60))
	_, _ = labelColor.Print("   Total: ")
	_, _ = valueColor.Printf("%d MB (%.2f GB)\n", info.Memory.TotalMB, info.Memory.TotalGB)
	_, _ = labelColor.Print("   Used:  ")
	_, _ = valueColor.Printf("%d MB (%.2f GB)\n", info.Memory.UsedMB, info.Memory.UsedGB)
	_, _ = labelColor.Print("   Free:  ")
	_, _ = successColor.Printf("%d MB (%.2f GB)\n", info.Memory.FreeMB, info.Memory.FreeGB)

	// Memory progress bar
	usedPercent := float64(info.Memory.UsedMB) / float64(info.Memory.TotalMB) * 100
	printProgressBar("   Usage", usedPercent, 40)
	fmt.Println()

	// Disk Section
	_, _ = headerColor.Println("💿 DISK DRIVES")
	fmt.Println(strings.Repeat("─", 60))
	for _, disk := range info.Disks {
		_, _ = labelColor.Printf("   Drive %s\n", disk.Name)
		_, _ = valueColor.Printf("      Total: %d MB (%.2f GB)\n", disk.TotalMB, disk.TotalGB)
		_, _ = valueColor.Printf("      Used:  %d MB (%.2f GB)\n", disk.UsedMB, disk.UsedGB)
		_, _ = successColor.Printf("      Free:  %d MB (%.2f GB)\n", disk.FreeMB, disk.FreeGB)

		diskUsedPercent := float64(disk.UsedMB) / float64(disk.TotalMB) * 100
		printProgressBar("      Usage", diskUsedPercent, 35)
		fmt.Println()
	}

	// Hyper-V Section
	_, _ = headerColor.Println("🔧 HYPER-V STATUS")
	fmt.Println(strings.Repeat("─", 60))
	_, _ = labelColor.Print("   Status:  ")
	if info.HyperV.Enabled {
		_, _ = successColor.Printf("✅ %s (Enabled)\n", info.HyperV.Status)
	} else {
		_, _ = errorColor.Printf("❌ %s (Disabled)\n", info.HyperV.Status)
	}
	fmt.Println()

	_, _ = headerColor.Println("════════════════════════════════════════════════════════════════")
}

// printProgressBar prints a visual progress bar
func printProgressBar(label string, percent float64, width int) {
	filled := int(percent / 100 * float64(width))
	empty := width - filled

	// Choose color based on usage
	// Choose color based on usage
	var barColor *color.Color
	switch {
	case percent >= 90:
		barColor = color.New(color.FgRed)
	case percent >= 70:
		barColor = color.New(color.FgYellow)
	default:
		barColor = color.New(color.FgGreen)
	}

	fmt.Printf("%s: [", label)
	_, _ = barColor.Print(strings.Repeat("█", filled))
	fmt.Print(strings.Repeat("░", empty))
	fmt.Printf("] %.1f%%\n", percent)
}
