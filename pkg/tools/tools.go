// Package tools 存放辅助方法
package tools

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/imroc/req/v3"
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

// VersionCompare 版本比较
func VersionCompare(ver1, ver2, operator string) bool {
	v1 := strings.TrimPrefix(ver1, "v")
	v2 := strings.TrimPrefix(ver2, "v")

	v1s := strings.Split(v1, ".")
	v2s := strings.Split(v2, ".")

	for len(v1s) < len(v2s) {
		v1s = append(v1s, "0")
	}

	for len(v2s) < len(v1s) {
		v2s = append(v2s, "0")
	}

	for i := 0; i < len(v1s); i++ {
		if v1s[i] > v2s[i] {
			return operator == ">" || operator == ">=" || operator == "!="
		} else if v1s[i] < v2s[i] {
			return operator == "<" || operator == "<=" || operator == "!="
		}
	}
	return operator == "==" || operator == ">=" || operator == "<="
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
	var output string
	isChina := IsChina()

	if isChina {
		output = Exec(`curl -sSL "https://jihulab.com/api/v4/projects/haozi-team%2Fpanel/releases/permalink/latest"`)
	} else {
		output = Exec(`curl -sSL "https://api.github.com/repos/haozi-team/panel/releases/latest"`)
	}

	if len(output) == 0 {
		return info, errors.New("获取最新版本失败")
	}

	file, err := os.CreateTemp("", "panel")
	if err != nil {
		return info, errors.New("创建临时文件失败")
	}
	defer os.Remove(file.Name())
	_, err = file.Write([]byte(output))
	if err != nil {
		return info, errors.New("写入临时文件失败")
	}
	err = file.Close()
	if err != nil {
		return info, errors.New("关闭临时文件失败")
	}
	fileName := file.Name()

	var name, version, body, date, downloadUrl string
	if isChina {
		name = Exec("jq -r '.name' " + fileName)
		version = Exec("jq -r '.tag_name' " + fileName)
		body = Exec("jq -r '.description' " + fileName)
		date = Exec("jq -r '.created_at' " + fileName)
		if IsArm() {
			downloadUrl = Exec("jq -r '.assets.links[] | select(.name | contains(\"arm64\")) | .direct_asset_url' " + fileName)
		} else {
			downloadUrl = Exec("jq -r '.assets.links[] | select(.name | contains(\"amd64v2\")) | .direct_asset_url' " + fileName)
		}
	} else {
		name = Exec("jq -r '.name' " + fileName)
		version = Exec("jq -r '.tag_name' " + fileName)
		body = Exec("jq -r '.body' " + fileName)
		date = Exec("jq -r '.published_at' " + fileName)
		if IsArm() {
			downloadUrl = Exec("jq -r '.assets[] | select(.name | contains(\"arm64\")) | .browser_download_url' " + fileName)
		} else {
			downloadUrl = Exec("jq -r '.assets[] | select(.name | contains(\"amd64v2\")) | .browser_download_url' " + fileName)
		}
	}

	info.Name = name
	info.Version = version
	info.Body = body
	info.Date = date
	info.DownloadUrl = downloadUrl

	return info, nil
}

// GetPanelVersion 获取指定面板版本
func GetPanelVersion(version string) (PanelInfo, error) {
	var info PanelInfo
	var output string
	isChina := IsChina()

	if isChina {
		output = Exec(`curl -sSL "https://jihulab.com/api/v4/projects/haozi-team%2Fpanel/releases/"` + version + `"`)
	} else {
		output = Exec(`curl -sSL "https://api.github.com/repos/haozi-team/panel/releases/tags/` + version + `"`)
	}

	if len(output) == 0 {
		return info, errors.New("获取面板版本失败")
	}

	file, err := os.CreateTemp("", "panel")
	if err != nil {
		return info, errors.New("创建临时文件失败")
	}
	defer os.Remove(file.Name())
	_, err = file.Write([]byte(output))
	if err != nil {
		return info, errors.New("写入临时文件失败")
	}
	err = file.Close()
	if err != nil {
		return info, errors.New("关闭临时文件失败")
	}
	fileName := file.Name()

	var name, version2, body, date, downloadUrl string
	if isChina {
		name = Exec("jq -r '.name' " + fileName)
		version2 = Exec("jq -r '.tag_name' " + fileName)
		body = Exec("jq -r '.description' " + fileName)
		date = Exec("jq -r '.created_at' " + fileName)
		if IsArm() {
			downloadUrl = Exec("jq -r '.assets.links[] | select(.name | contains(\"arm64\")) | .direct_asset_url' " + fileName)
		} else {
			downloadUrl = Exec("jq -r '.assets.links[] | select(.name | contains(\"amd64v2\")) | .direct_asset_url' " + fileName)
		}
	} else {
		name = Exec("jq -r '.name' " + fileName)
		version2 = Exec("jq -r '.tag_name' " + fileName)
		body = Exec("jq -r '.body' " + fileName)
		date = Exec("jq -r '.published_at' " + fileName)
		if IsArm() {
			downloadUrl = Exec("jq -r '.assets[] | select(.name | contains(\"arm64\")) | .browser_download_url' " + fileName)
		} else {
			downloadUrl = Exec("jq -r '.assets[] | select(.name | contains(\"amd64v2\")) | .browser_download_url' " + fileName)
		}
	}

	info.Name = name
	info.Version = version2
	info.Body = body
	info.Date = date
	info.DownloadUrl = downloadUrl

	return info, nil
}

// UpdatePanel 更新面板
func UpdatePanel() error {
	panelInfo, err := GetLatestPanelVersion()
	if err != nil {
		return err
	}

	color.Greenln("最新版本: " + panelInfo.Version)
	color.Greenln("下载链接: " + panelInfo.DownloadUrl)

	color.Greenln("备份面板配置...")
	Exec("cp -f /www/panel/database/panel.db /tmp/panel.db.bak")
	Exec("cp -f /www/panel/panel.conf /tmp/panel.conf.bak")
	if !Exists("/tmp/panel.db.bak") || !Exists("/tmp/panel.conf.bak") {
		return errors.New("备份面板配置失败")
	}
	color.Greenln("备份完成")

	color.Greenln("清理旧版本...")
	Exec("rm -rf /www/panel/*")
	color.Greenln("清理完成")

	color.Greenln("正在下载...")
	Exec("wget -T 120 -t 3 -O /www/panel/panel.zip " + panelInfo.DownloadUrl)

	if !Exists("/www/panel/panel.zip") {
		return errors.New("下载失败")
	}

	color.Greenln("下载完成")

	color.Greenln("更新新版本...")
	Exec("cd /www/panel && unzip -o panel.zip && rm -rf panel.zip && chmod 700 panel && bash scripts/update_panel.sh")

	if !Exists("/www/panel/panel") {
		return errors.New("更新失败，可能是下载过程中出现了问题")
	}

	color.Greenln("更新完成")

	color.Greenln("恢复面板配置...")
	Exec("cp -f /tmp/panel.db.bak /www/panel/database/panel.db")
	Exec("cp -f /tmp/panel.conf.bak /www/panel/panel.conf")
	if !Exists("/www/panel/database/panel.db") || !Exists("/www/panel/panel.conf") {
		return errors.New("恢复面板配置失败")
	}
	Exec("/www/panel/panel --env=panel.conf artisan migrate")
	color.Greenln("恢复完成")

	Exec("panel writeSetting version " + panelInfo.Version)

	Exec("rm -rf /tmp/panel.db.bak")
	Exec("rm -rf /tmp/panel.conf.bak")

	color.Greenln("重启面板...")
	Exec("systemctl restart panel")
	color.Greenln("重启完成")

	return nil
}

// IsChina 是否中国大陆
func IsChina() bool {
	client := req.C()
	client.SetTimeout(5 * time.Second)
	client.SetCommonRetryCount(2)
	client.ImpersonateSafari()

	resp, err := client.R().Get("https://www.cloudflare-cn.com/cdn-cgi/trace")
	if err != nil || !resp.IsSuccessState() {
		return false
	}

	if strings.Contains(resp.String(), "loc=CN") {
		return true
	}

	return false
}
