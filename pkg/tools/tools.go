// Package tools 存放辅助方法
package tools

import (
	"errors"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gookit/color"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"

	"github.com/TheTNB/panel/pkg/shell"
)

// MonitoringInfo 监控信息
type MonitoringInfo struct {
	Cpus      []cpu.InfoStat         `json:"cpus"`
	Percent   []float64              `json:"percent"`
	Load      *load.AvgStat          `json:"load"`
	Host      *host.InfoStat         `json:"host"`
	Mem       *mem.VirtualMemoryStat `json:"mem"`
	Swap      *mem.SwapMemoryStat    `json:"swap"`
	Net       []net.IOCountersStat   `json:"net"`
	DiskIO    []disk.IOCountersStat  `json:"disk_io"`
	Disk      []disk.PartitionStat   `json:"disk"`
	DiskUsage []disk.UsageStat       `json:"disk_usage"`
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
	res.Disk, _ = disk.Partitions(true)

	ioCounters, _ := disk.IOCounters()
	for _, info := range ioCounters {
		res.DiskIO = append(res.DiskIO, info)
	}

	for _, partition := range res.Disk {
		if strings.HasPrefix(partition.Mountpoint, "/dev") || strings.HasPrefix(partition.Mountpoint, "/sys") || strings.HasPrefix(partition.Mountpoint, "/proc") || strings.HasPrefix(partition.Mountpoint, "/run") || strings.HasPrefix(partition.Mountpoint, "/boot") || strings.HasPrefix(partition.Mountpoint, "/usr") || strings.HasPrefix(partition.Mountpoint, "/var") {
			continue
		}
		usage, _ := disk.Usage(partition.Mountpoint)
		res.DiskUsage = append(res.DiskUsage, *usage)
	}

	return res
}

func RestartPanel() {
	color.Greenln("重启面板...")
	_ = shell.ExecfAsync("sleep 1 && systemctl restart panel")
	color.Greenln("重启完成")
}

// IsChina 是否中国大陆
func IsChina() bool {
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	resp, err := client.R().Get("https://www.cloudflare-cn.com/cdn-cgi/trace")
	if err != nil || !resp.IsSuccess() {
		return false
	}

	if strings.Contains(resp.String(), "loc=CN") {
		return true
	}

	return false
}

// GetPublicIP 获取公网IP
func GetPublicIP() (string, error) {
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	resp, err := client.R().Get("https://www.cloudflare-cn.com/cdn-cgi/trace")
	if err != nil || !resp.IsSuccess() {
		return "", errors.New("获取公网IP失败")
	}

	return strings.TrimPrefix(strings.Split(resp.String(), "\n")[2], "ip="), nil
}
