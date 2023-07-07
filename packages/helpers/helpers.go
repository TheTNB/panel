// Package helpers 存放辅助方法
package helpers

import (
	"errors"
	"os"
	"os/exec"
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

type PanelInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	DownloadUrl string `json:"download_url"`
	Body        string `json:"body"`
	Date        string `json:"date"`
}

// GetLatestPanelVersion 获取最新版本
func GetLatestPanelVersion() (PanelInfo, error) {
	var info PanelInfo

	cmd := exec.Command("/bin/bash", "-c", "curl \"https://api.github.com/repos/HaoZi-Team/Panel/releases/latest\"")
	output, err := cmd.Output()
	if err != nil {
		return info, errors.New("获取最新版本失败")
	}

	file, err := os.CreateTemp("", "panel")
	if err != nil {
		return info, errors.New("创建临时文件失败")
	}
	defer os.Remove(file.Name())
	_, err = file.Write(output)
	if err != nil {
		return info, errors.New("写入临时文件失败")
	}
	err = file.Close()
	if err != nil {
		return info, errors.New("关闭临时文件失败")
	}
	fileName := file.Name()

	name := exec.Command("/bin/bash", "-c", "jq -r '.name' "+fileName)
	version := exec.Command("/bin/bash", "-c", "jq -r '.tag_name' "+fileName)
	body := exec.Command("/bin/bash", "-c", "jq -r '.body' "+fileName)
	date := exec.Command("/bin/bash", "-c", "jq -r '.published_at' "+fileName)
	nameOutput, _ := name.Output()
	versionOutput, _ := version.Output()
	bodyOutput, _ := body.Output()
	dateOutput, _ := date.Output()
	info.Name = strings.TrimSpace(string(nameOutput))
	info.Version = strings.TrimSpace(string(versionOutput))
	info.Body = strings.TrimSpace(string(bodyOutput))
	info.Date = strings.TrimSpace(string(dateOutput))
	if IsArm() {
		downloadUrl := exec.Command("/bin/bash", "-c", "jq -r '.assets[] | select(.name | contains(\"arm64\")) | .browser_download_url' "+fileName)
		downloadUrlOutput, _ := downloadUrl.Output()
		info.DownloadUrl = strings.TrimSpace(string(downloadUrlOutput))
	} else {
		downloadUrl := exec.Command("/bin/bash", "-c", "jq -r '.assets[] | select(.name | contains(\"amd64v3\")) | .browser_download_url' "+fileName)
		downloadUrlOutput, _ := downloadUrl.Output()
		info.DownloadUrl = strings.TrimSpace(string(downloadUrlOutput))
	}

	return info, nil
}

// UpdatePanel 更新面板
func UpdatePanel() error {
	panelInfo, err := GetLatestPanelVersion()
	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/bash", "-c", "wget -O panel.tar.gz "+panelInfo.DownloadUrl+" && tar -zxvf panel.tar.gz && rm -rf panel.tar.gz && chmod +x panel && ./panel artisan migrate")
	_, err = cmd.Output()
	if err != nil {
		return errors.New("更新面板失败")
	}

	return nil
}
