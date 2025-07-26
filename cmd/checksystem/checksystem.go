package checksystem

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

// Экспортируемая команда
var Cmd = &cobra.Command{
	Use:   "checksystem",
	Short: "Show system metrics like CPU, memory, disk and uptime",
	Run: func(cmd *cobra.Command, args []string) {
		printSystemInfo()
	},
}

func printSystemInfo() {
	// CPU usage
	cpuPercents, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println("Error getting CPU info:", err)
		return
	}
	fmt.Printf("CPU Usage: %.2f%%\n", cpuPercents[0])

	// Memory info
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Error getting memory info:", err)
		return
	}
	fmt.Printf("Memory Usage: %.2f%% (Used: %v MB, Total: %v MB)\n",
		vmStat.UsedPercent,
		vmStat.Used/1024/1024,
		vmStat.Total/1024/1024,
	)

	// Disk usage
	diskStat, err := disk.Usage("/")
	if err != nil {
		fmt.Println("Error getting disk info:", err)
		return
	}
	fmt.Printf("Disk Usage (/): %.2f%% (Used: %v GB, Total: %v GB)\n",
		diskStat.UsedPercent,
		diskStat.Used/1024/1024/1024,
		diskStat.Total/1024/1024/1024,
	)

	// System uptime
	uptime, err := host.Uptime()
	if err != nil {
		fmt.Println("Error getting uptime:", err)
		return
	}
	fmt.Printf("System Uptime: %v hours\n", uptime/3600)
}
