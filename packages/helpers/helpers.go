// Package helpers 存放辅助方法
package helpers

import (
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// MonitoringInfo 监控信息
type MonitoringInfo struct {
	Cpus      []cpu.InfoStat                 `json:"cpus"`
	Percent   []float64                      `json:"percent"`
	Load      *load.AvgStat                  `json:"load"`
	Host      *host.InfoStat                 `json:"host"`
	Mem       *mem.VirtualMemoryStat         `json:"mem"`
	Swap      *mem.SwapMemoryStat            `json:"swap"`
	Net       []net.IOCountersStat           `json:"net"`
	DiskIO    map[string]disk.IOCountersStat `json:"disk_io"`
	Disk      []disk.PartitionStat           `json:"disk"`
	DiskUsage map[string]*disk.UsageStat     `json:"disk_usage"`
}

// GetMonitoringInfo 获取监控数据
func GetMonitoringInfo() MonitoringInfo {
	var res MonitoringInfo
	res.Cpus, _ = cpu.Info()
	res.Percent, _ = cpu.Percent(time.Second, false)
	res.Load, _ = load.Avg()
	res.Host, _ = host.Info()
	res.Mem, _ = mem.VirtualMemory()
	res.Swap, _ = mem.SwapMemory()
	res.Net, _ = net.IOCounters(true)
	res.DiskIO, _ = disk.IOCounters()
	res.Disk, _ = disk.Partitions(true)

	res.DiskUsage = make(map[string]*disk.UsageStat)
	for _, partition := range res.Disk {
		if strings.HasPrefix(partition.Mountpoint, "/dev") || strings.HasPrefix(partition.Mountpoint, "/sys") || strings.HasPrefix(partition.Mountpoint, "/proc") || strings.HasPrefix(partition.Mountpoint, "/run") || strings.HasPrefix(partition.Mountpoint, "/boot") || strings.HasPrefix(partition.Mountpoint, "/usr") {
			continue
		}
		usage, _ := disk.Usage(partition.Mountpoint)
		res.DiskUsage[partition.Mountpoint] = usage
	}

	return res
}
