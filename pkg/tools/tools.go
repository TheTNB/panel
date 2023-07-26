// Package tools 存放辅助方法
package tools

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
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
		if strings.HasPrefix(partition.Mountpoint, "/dev") || strings.HasPrefix(partition.Mountpoint, "/sys") || strings.HasPrefix(partition.Mountpoint, "/proc") || strings.HasPrefix(partition.Mountpoint, "/run") || strings.HasPrefix(partition.Mountpoint, "/boot") || strings.HasPrefix(partition.Mountpoint, "/usr") || strings.HasPrefix(partition.Mountpoint, "/var") {
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
		downloadUrl := exec.Command("/bin/bash", "-c", "jq -r '.assets[] | select(.name | contains(\"amd64v2\")) | .browser_download_url' "+fileName)
		downloadUrlOutput, _ := downloadUrl.Output()
		info.DownloadUrl = strings.TrimSpace(string(downloadUrlOutput))
	}

	return info, nil
}

// UpdatePanel 更新面板
func UpdatePanel(proxy bool) error {
	panelInfo, err := GetLatestPanelVersion()
	if err != nil {
		return err
	}

	color.Greenln("最新版本: " + panelInfo.Version)
	color.Greenln("下载链接: " + panelInfo.DownloadUrl)
	color.Greenln("使用代理: " + strconv.FormatBool(proxy))

	color.Greenln("备份面板配置...")
	ExecShell("cp -f /www/panel/database/panel.db /tmp/panel.db.bak")
	ExecShell("cp -f /www/panel/panel.conf /tmp/panel.conf.bak")
	if !Exists("/tmp/panel.db.bak") || !Exists("/tmp/panel.conf.bak") {
		return errors.New("备份面板配置失败")
	}
	color.Greenln("备份完成")

	color.Greenln("清理旧版本...")
	ExecShell("rm -rf /www/panel/*")
	color.Greenln("清理完成")

	color.Greenln("正在下载...")
	if proxy {
		ExecShell("wget -O /www/panel/panel.zip https://ghproxy.com/" + panelInfo.DownloadUrl)
	} else {
		ExecShell("wget -O /www/panel/panel.zip " + panelInfo.DownloadUrl)
	}
	color.Greenln("下载完成")

	color.Greenln("更新新版本...")
	ExecShell("cd /www/panel && unzip -o panel.zip && rm -rf panel.zip && chmod 700 panel")
	color.Greenln("更新完成")

	color.Greenln("恢复面板配置...")
	ExecShell("cp -f /tmp/panel.db.bak /www/panel/database/panel.db")
	ExecShell("cp -f /tmp/panel.conf.bak /www/panel/panel.conf")
	if !Exists("/www/panel/database/panel.db") || !Exists("/www/panel/panel.conf") {
		return errors.New("恢复面板配置失败")
	}
	ExecShell("/www/panel/panel --env=panel.conf artisan migrate")
	color.Greenln("恢复完成")

	color.Greenln("重启面板...")
	ExecShell("systemctl restart panel")
	color.Greenln("重启完成")

	return nil
}
