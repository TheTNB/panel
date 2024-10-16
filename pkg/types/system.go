package types

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// CurrentInfo 监控信息
type CurrentInfo struct {
	Cpus      []cpu.InfoStat         `json:"cpus"`
	Percent   float64                `json:"percent"`  // 总使用率
	Percents  []float64              `json:"percents"` // 每个核心使用率
	Load      *load.AvgStat          `json:"load"`
	Host      *host.InfoStat         `json:"host"`
	Mem       *mem.VirtualMemoryStat `json:"mem"`
	Swap      *mem.SwapMemoryStat    `json:"swap"`
	Net       []net.IOCountersStat   `json:"net"`
	DiskIO    []disk.IOCountersStat  `json:"disk_io"`
	Disk      []disk.PartitionStat   `json:"disk"`
	DiskUsage []disk.UsageStat       `json:"disk_usage"`
	Time      time.Time              `json:"time"`
}
