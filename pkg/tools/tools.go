// Package tools 存放辅助方法
package tools

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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

	startMajor, err := strconv.Atoi(startParts[0])
	if err != nil {
		return nil, fmt.Errorf("无效的起始主版本号: %v", err)
	}
	startMinor, err := strconv.Atoi(startParts[1])
	if err != nil {
		return nil, fmt.Errorf("无效的起始次版本号: %v", err)
	}
	startPatch, err := strconv.Atoi(startParts[2])
	if err != nil {
		return nil, fmt.Errorf("无效的起始修订号: %v", err)
	}
	endMajor, err := strconv.Atoi(endParts[0])
	if err != nil {
		return nil, fmt.Errorf("无效的结束主版本号: %v", err)
	}
	endMinor, err := strconv.Atoi(endParts[1])
	if err != nil {
		return nil, fmt.Errorf("无效的结束次版本号: %v", err)
	}
	endPatch, err := strconv.Atoi(endParts[2])
	if err != nil {
		return nil, fmt.Errorf("无效的结束修订号: %v", err)
	}

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
		output, err = Exec(`curl -sSL "https://jihulab.com/api/v4/projects/haozi-team%2Fpanel/releases/permalink/latest"`)
	} else {
		output, err = Exec(`curl -sSL "https://api.github.com/repos/haozi-team/panel/releases/latest"`)
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
		if name, err = Exec("jq -r '.name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if version, err = Exec("jq -r '.tag_name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if body, err = Exec("jq -r '.description' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if date, err = Exec("jq -r '.created_at' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if checksums, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"checksums\")) | .name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if checksumsUrl, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"checksums\")) | .direct_asset_url' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if IsArm() {
			if downloadName, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"arm64\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
			if downloadUrl, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"arm64\")) | .direct_asset_url' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
		} else {
			if downloadName, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"amd64v2\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
			if downloadUrl, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"amd64v2\")) | .direct_asset_url' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
		}
	} else {
		if name, err = Exec("jq -r '.name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if version, err = Exec("jq -r '.tag_name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if body, err = Exec("jq -r '.body' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if date, err = Exec("jq -r '.published_at' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if checksums, err = Exec("jq -r '.assets[] | select(.name | contains(\"checksums\")) | .name' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if checksumsUrl, err = Exec("jq -r '.assets[] | select(.name | contains(\"checksums\")) | .browser_download_url' " + fileName); err != nil {
			return info, errors.New("获取最新版本失败")
		}
		if IsArm() {
			if downloadName, err = Exec("jq -r '.assets[] | select(.name | contains(\"arm64\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
			if downloadUrl, err = Exec("jq -r '.assets[] | select(.name | contains(\"arm64\")) | .browser_download_url' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
		} else {
			if downloadName, err = Exec("jq -r '.assets[] | select(.name | contains(\"amd64v2\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取最新版本失败")
			}
			if downloadUrl, err = Exec("jq -r '.assets[] | select(.name | contains(\"amd64v2\")) | .browser_download_url' " + fileName); err != nil {
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
		output, err = Exec(`curl -sSL "https://jihulab.com/api/v4/projects/haozi-team%2Fpanel/releases/` + version + `"`)
	} else {
		output, err = Exec(`curl -sSL "https://api.github.com/repos/haozi-team/panel/releases/tags/` + version + `"`)
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
		if name, err = Exec("jq -r '.name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if version2, err = Exec("jq -r '.tag_name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if body, err = Exec("jq -r '.description' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if date, err = Exec("jq -r '.created_at' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if checksums, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"checksums\")) | .name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if checksumsUrl, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"checksums\")) | .direct_asset_url' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if IsArm() {
			if downloadName, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"arm64\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
			if downloadUrl, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"arm64\")) | .direct_asset_url' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
		} else {
			if downloadName, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"amd64v2\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
			if downloadUrl, err = Exec("jq -r '.assets.links[] | select(.name | contains(\"amd64v2\")) | .direct_asset_url' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
		}
	} else {
		if name, err = Exec("jq -r '.name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if version2, err = Exec("jq -r '.tag_name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if body, err = Exec("jq -r '.body' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if date, err = Exec("jq -r '.published_at' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if checksums, err = Exec("jq -r '.assets[] | select(.name | contains(\"checksums\")) | .name' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if checksumsUrl, err = Exec("jq -r '.assets[] | select(.name | contains(\"checksums\")) | .browser_download_url' " + fileName); err != nil {
			return info, errors.New("获取面板版本失败")
		}
		if IsArm() {
			if downloadName, err = Exec("jq -r '.assets[] | select(.name | contains(\"arm64\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
			if downloadUrl, err = Exec("jq -r '.assets[] | select(.name | contains(\"arm64\")) | .browser_download_url' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
		} else {
			if downloadName, err = Exec("jq -r '.assets[] | select(.name | contains(\"amd64v2\")) | .name' " + fileName); err != nil {
				return info, errors.New("获取面板版本失败")
			}
			if downloadUrl, err = Exec("jq -r '.assets[] | select(.name | contains(\"amd64v2\")) | .browser_download_url' " + fileName); err != nil {
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
	color.Greenln("目标版本: " + panelInfo.Version)
	color.Greenln("下载链接: " + panelInfo.DownloadUrl)

	color.Greenln("备份面板配置...")
	if _, err := Exec("cp -f /www/panel/database/panel.db /tmp/panel.db.bak"); err != nil {
		color.Redln("备份面板数据库失败")
		return err
	}
	if _, err := Exec("cp -f /www/panel/panel.conf /tmp/panel.conf.bak"); err != nil {
		color.Redln("备份面板配置失败")
		return err
	}
	if !Exists("/tmp/panel.db.bak") || !Exists("/tmp/panel.conf.bak") {
		return errors.New("备份面板配置失败")
	}
	color.Greenln("备份完成")

	color.Greenln("清理旧版本...")
	if _, err := Exec("rm -rf /www/panel/*"); err != nil {
		color.Redln("清理旧版本失败")
		return err
	}
	color.Greenln("清理完成")

	color.Greenln("正在下载...")
	if _, err := Exec("wget -T 120 -t 3 -O /www/panel/" + panelInfo.DownloadName + " " + panelInfo.DownloadUrl); err != nil {
		color.Redln("下载失败")
		return err
	}
	if _, err := Exec("wget -T 20 -t 3 -O /www/panel/" + panelInfo.Checksums + " " + panelInfo.ChecksumsUrl); err != nil {
		color.Redln("下载失败")
		return err
	}
	if !Exists("/www/panel/"+panelInfo.DownloadName) || !Exists("/www/panel/"+panelInfo.Checksums) {
		return errors.New("下载失败")
	}
	color.Greenln("下载完成")

	color.Greenln("校验下载文件...")
	check, err := Exec("cd /www/panel && sha256sum -c " + panelInfo.Checksums + " --ignore-missing")
	if check != panelInfo.DownloadName+": OK" || err != nil {
		return errors.New("下载文件校验失败")
	}
	if err := Remove("/www/panel/" + panelInfo.Checksums); err != nil {
		color.Redln("清理临时文件失败")
		return err
	}
	color.Greenln("文件校验完成")

	color.Greenln("更新新版本...")
	if _, err = Exec("cd /www/panel && unzip -o " + panelInfo.DownloadName + " && rm -rf " + panelInfo.DownloadName); err != nil {
		color.Redln("更新失败")
		return err
	}
	if !Exists("/www/panel/panel") {
		return errors.New("更新失败，可能是下载过程中出现了问题")
	}
	color.Greenln("更新完成")

	color.Greenln("恢复面板配置...")
	if _, err = Exec("cp -f /tmp/panel.db.bak /www/panel/database/panel.db"); err != nil {
		color.Redln("恢复面板数据库失败")
		return err
	}
	if _, err = Exec("cp -f /tmp/panel.conf.bak /www/panel/panel.conf"); err != nil {
		color.Redln("恢复面板配置失败")
		return err
	}
	if !Exists("/www/panel/database/panel.db") || !Exists("/www/panel/panel.conf") {
		return errors.New("恢复面板配置失败")
	}
	if _, err = Exec("/www/panel/panel --env=panel.conf artisan migrate"); err != nil {
		color.Redln("运行面板数据库迁移失败")
		return err
	}
	color.Greenln("恢复完成")

	color.Greenln("设置面板文件权限...")
	if _, err = Exec("chmod -R 700 /www/panel"); err != nil {
		color.Redln("设置面板文件权限失败")
		return err
	}
	color.Greenln("设置完成")

	if _, err = Exec("bash /www/panel/scripts/update_panel.sh"); err != nil {
		color.Redln("执行面板升级后脚本失败")
		return err
	}
	if _, err = Exec("panel writeSetting version " + panelInfo.Version); err != nil {
		color.Redln("写入面板版本号失败")
		return err
	}

	if _, err = Exec("rm -rf /tmp/panel.db.bak"); err != nil {
		color.Redln("清理临时文件失败")
		return err
	}
	if _, err = Exec("rm -rf /tmp/panel.conf.bak"); err != nil {
		color.Redln("清理临时文件失败")
		return err
	}

	return nil
}

func RestartPanel() {
	color.Greenln("重启面板...")
	err := ExecAsync("sleep 2 && systemctl restart panel")
	if err != nil {
		color.Redln("重启失败")
		return
	}

	color.Greenln("重启完成")
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
