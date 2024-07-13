// Package tools 存放辅助方法
package tools

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/env"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
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

// GenerateVersions 获取版本列表
func GenerateVersions(start, end string) ([]string, error) {
	var versions []string
	start = strings.TrimPrefix(start, "v")
	end = strings.TrimPrefix(end, "v")
	startParts := strings.Split(start, ".")
	endParts := strings.Split(end, ".")

	if len(startParts) != 3 || len(endParts) != 3 {
		return nil, fmt.Errorf("版本格式错误")
	}

	startMajor := cast.ToInt(startParts[0])
	startMinor := cast.ToInt(startParts[1])
	startPatch := cast.ToInt(startParts[2])
	endMajor := cast.ToInt(endParts[0])
	endMinor := cast.ToInt(endParts[1])
	endPatch := cast.ToInt(endParts[2])

	for major := startMajor; major <= endMajor; major++ {
		for minor := 0; minor <= 99; minor++ {
			for patch := 0; patch <= 99; patch++ {
				if major == startMajor && minor < startMinor {
					continue
				}
				if major == startMajor && minor == startMinor && patch <= startPatch {
					continue
				}

				if major == endMajor && minor > endMinor {
					return versions, nil
				}
				if major == endMajor && minor == endMinor && patch > endPatch {
					return versions, nil
				}

				versions = append(versions, fmt.Sprintf("%d.%d.%d", major, minor, patch))
			}
		}
	}

	if len(versions) == 0 {
		return []string{}, nil
	}

	return versions, nil
}

type PanelInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	DownloadName string `json:"download_name"`
	DownloadUrl  string `json:"download_url"`
	Body         string `json:"body"`
	Date         string `json:"date"`
	Checksums    string `json:"checksums"`
	ChecksumsUrl string `json:"checksums_url"`
}

// GetLatestPanelVersion 获取最新版本
func GetLatestPanelVersion() (PanelInfo, error) {
	var info PanelInfo
	var output string
	var err error
	isChina := IsChina()

	if isChina {
		output, err = shell.Execf(`curl -sSL "https://git.haozi.net/api/v4/projects/opensource%%2Fpanel/releases/permalink/latest"`)
	} else {
		output, err = shell.Execf(`curl -sSL "https://api.github.com/repos/TheTNB/panel/releases/latest"`)
	}

	if len(output) == 0 || err != nil {
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

	var name, version, body, date, downloadName, downloadUrl, checksums, checksumsUrl string
	if isChina {
		if name, err = shell.Execf("jq -r '.name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if version, err = shell.Execf("jq -r '.tag_name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if body, err = shell.Execf("jq -r '.description' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if date, err = shell.Execf("jq -r '.created_at' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if checksums, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"checksums\")) | .name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if checksumsUrl, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"checksums\")) | .direct_asset_url' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if env.IsArm() {
			if downloadName, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"arm64\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
			if downloadUrl, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"arm64\")) | .direct_asset_url' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
		} else {
			if downloadName, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"amd64v2\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
			if downloadUrl, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"amd64v2\")) | .direct_asset_url' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
		}
	} else {
		if name, err = shell.Execf("jq -r '.name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if version, err = shell.Execf("jq -r '.tag_name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if body, err = shell.Execf("jq -r '.body' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if date, err = shell.Execf("jq -r '.published_at' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if checksums, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"checksums\")) | .name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if checksumsUrl, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"checksums\")) | .browser_download_url' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if env.IsArm() {
			if downloadName, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"arm64\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
			if downloadUrl, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"arm64\")) | .browser_download_url' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
		} else {
			if downloadName, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"amd64v2\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
			if downloadUrl, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"amd64v2\")) | .browser_download_url' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
		}
	}

	info.Name = name
	info.Version = version
	info.Body = body
	info.Date = date
	info.DownloadName = downloadName
	info.DownloadUrl = downloadUrl
	info.Checksums = checksums
	info.ChecksumsUrl = checksumsUrl

	return info, nil
}

// GetPanelVersion 获取指定面板版本
func GetPanelVersion(version string) (PanelInfo, error) {
	var info PanelInfo
	var output string
	var err error
	isChina := IsChina()

	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	if isChina {
		output, err = shell.Execf(`curl -sSL "https://git.haozi.net/api/v4/projects/opensource%%2Fpanel/releases/` + version + `"`)
	} else {
		output, err = shell.Execf(`curl -sSL "https://api.github.com/repos/TheTNB/panel/releases/tags/` + version + `"`)
	}

	if len(output) == 0 || err != nil {
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

	var name, version2, body, date, downloadName, downloadUrl, checksums, checksumsUrl string
	if isChina {
		if name, err = shell.Execf("jq -r '.name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if version2, err = shell.Execf("jq -r '.tag_name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if body, err = shell.Execf("jq -r '.description' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if date, err = shell.Execf("jq -r '.created_at' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if checksums, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"checksums\")) | .name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if checksumsUrl, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"checksums\")) | .direct_asset_url' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if env.IsArm() {
			if downloadName, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"arm64\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
			if downloadUrl, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"arm64\")) | .direct_asset_url' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
		} else {
			if downloadName, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"amd64v2\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
			if downloadUrl, err = shell.Execf("jq -r '.assets.links[] | select(.name | contains(\"amd64v2\")) | .direct_asset_url' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
		}
	} else {
		if name, err = shell.Execf("jq -r '.name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if version2, err = shell.Execf("jq -r '.tag_name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if body, err = shell.Execf("jq -r '.body' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if date, err = shell.Execf("jq -r '.published_at' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if checksums, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"checksums\")) | .name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if checksumsUrl, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"checksums\")) | .browser_download_url' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if env.IsArm() {
			if downloadName, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"arm64\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
			if downloadUrl, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"arm64\")) | .browser_download_url' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
		} else {
			if downloadName, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"amd64v2\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
			if downloadUrl, err = shell.Execf("jq -r '.assets[] | select(.name | contains(\"amd64v2\")) | .browser_download_url' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
		}
	}

	info.Name = name
	info.Version = version2
	info.Body = body
	info.Date = date
	info.DownloadName = downloadName
	info.DownloadUrl = downloadUrl
	info.Checksums = checksums
	info.ChecksumsUrl = checksumsUrl

	return info, nil
}

// UpdatePanel 更新面板
func UpdatePanel(panelInfo PanelInfo) error {
	color.Green().Printfln("目标版本: " + panelInfo.Version)
	color.Green().Printfln("下载链接: " + panelInfo.DownloadUrl)

	color.Green().Printfln("前置检查...")
	if io.Exists("/tmp/panel-storage.zip") || io.Exists("/tmp/panel.conf.bak") {
		return errors.New("检测到 /tmp 存在临时文件，可能是上次更新失败导致的，请谨慎排除后重试")
	}

	color.Green().Printfln("备份面板数据...")
	// 备份面板
	if err := io.Archive([]string{"/www/panel"}, "/www/backup/panel/panel-"+carbon.Now().ToShortDateTimeString()+".zip"); err != nil {
		color.Red().Printfln("备份面板失败")
		return err
	}
	if _, err := shell.Execf("cd /www/panel/storage && zip -r /tmp/panel-storage.zip *"); err != nil {
		color.Red().Printfln("备份面板数据失败")
		return err
	}
	if _, err := shell.Execf("cp -f /www/panel/panel.conf /tmp/panel.conf.bak"); err != nil {
		color.Red().Printfln("备份面板配置失败")
		return err
	}
	if !io.Exists("/tmp/panel-storage.zip") || !io.Exists("/tmp/panel.conf.bak") {
		return errors.New("备份面板数据失败")
	}
	color.Green().Printfln("备份完成")

	color.Green().Printfln("清理旧版本...")
	if _, err := shell.Execf("rm -rf /www/panel/*"); err != nil {
		color.Red().Printfln("清理旧版本失败")
		return err
	}
	color.Green().Printfln("清理完成")

	color.Green().Printfln("正在下载...")
	if _, err := shell.Execf("wget -T 120 -t 3 -O /www/panel/" + panelInfo.DownloadName + " " + panelInfo.DownloadUrl); err != nil {
		color.Red().Printfln("下载失败")
		return err
	}
	if _, err := shell.Execf("wget -T 20 -t 3 -O /www/panel/" + panelInfo.Checksums + " " + panelInfo.ChecksumsUrl); err != nil {
		color.Red().Printfln("下载失败")
		return err
	}
	if !io.Exists("/www/panel/"+panelInfo.DownloadName) || !io.Exists("/www/panel/"+panelInfo.Checksums) {
		return errors.New("下载失败")
	}
	color.Green().Printfln("下载完成")

	color.Green().Printfln("校验下载文件...")
	check, err := shell.Execf("cd /www/panel && sha256sum -c " + panelInfo.Checksums + " --ignore-missing")
	if check != panelInfo.DownloadName+": OK" || err != nil {
		return errors.New("下载文件校验失败")
	}
	if err = io.Remove("/www/panel/" + panelInfo.Checksums); err != nil {
		color.Red().Printfln("清理临时文件失败")
		return err
	}
	color.Green().Printfln("文件校验完成")

	color.Green().Printfln("更新新版本...")
	if _, err = shell.Execf("cd /www/panel && unzip -o " + panelInfo.DownloadName + " && rm -rf " + panelInfo.DownloadName); err != nil {
		color.Red().Printfln("更新失败")
		return err
	}
	if !io.Exists("/www/panel/panel") {
		return errors.New("更新失败，可能是下载过程中出现了问题")
	}
	color.Green().Printfln("更新完成")

	color.Green().Printfln("恢复面板数据...")
	if _, err = shell.Execf("cp -f /tmp/panel-storage.zip /www/panel/storage/panel-storage.zip && cd /www/panel/storage && unzip -o panel-storage.zip && rm -rf panel-storage.zip"); err != nil {
		color.Red().Printfln("恢复面板数据失败")
		return err
	}
	if _, err = shell.Execf("cp -f /tmp/panel.conf.bak /www/panel/panel.conf"); err != nil {
		color.Red().Printfln("恢复面板配置失败")
		return err
	}
	if _, err = shell.Execf("cp -f /www/panel/scripts/panel.sh /usr/bin/panel"); err != nil {
		color.Red().Printfln("恢复面板脚本失败")
		return err
	}
	if !io.Exists("/www/panel/storage/panel.db") || !io.Exists("/www/panel/panel.conf") {
		return errors.New("恢复面板数据失败")
	}
	color.Green().Printfln("恢复完成")

	color.Green().Printfln("设置面板文件权限...")
	_, _ = shell.Execf("chmod -R 700 /www/panel")
	_, _ = shell.Execf("chmod -R 700 /usr/bin/panel")
	color.Green().Printfln("设置完成")

	color.Green().Printfln("运行升级后脚本...")
	if _, err = shell.Execf("bash /www/panel/scripts/update_panel.sh"); err != nil {
		color.Red().Printfln("运行面板升级后脚本失败")
		return err
	}
	if _, err = shell.Execf("cp -f /www/panel/scripts/panel.service /etc/systemd/system/panel.service"); err != nil {
		color.Red().Printfln("写入面板服务文件失败")
		return err
	}
	_, _ = shell.Execf("systemctl daemon-reload")
	if _, err = shell.Execf("panel writeSetting version " + panelInfo.Version); err != nil {
		color.Red().Printfln("写入面板版本号失败")
		return err
	}
	color.Green().Printfln("升级完成")

	_, _ = shell.Execf("rm -rf /tmp/panel-storage.zip")
	_, _ = shell.Execf("rm -rf /tmp/panel.conf.bak")

	return nil
}

func RestartPanel() {
	color.Green().Printfln("重启面板...")
	err := shell.ExecfAsync("sleep 2 && systemctl restart panel")
	if err != nil {
		color.Red().Printfln("重启失败")
		return
	}

	color.Green().Printfln("重启完成")
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
