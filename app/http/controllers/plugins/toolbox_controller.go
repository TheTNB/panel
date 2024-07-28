package plugins

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/str"
)

type ToolBoxController struct {
}

func NewToolBoxController() *ToolBoxController {
	return &ToolBoxController{}
}

// GetDNS 获取 DNS 信息
func (r *ToolBoxController) GetDNS(ctx http.Context) http.Response {
	raw, err := io.Read("/etc/resolv.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	match := regexp.MustCompile(`nameserver\s+(\S+)`).FindAllStringSubmatch(raw, -1)
	if len(match) == 0 {
		return h.Error(ctx, http.StatusInternalServerError, "找不到 DNS 信息")
	}

	var dns []string
	for _, m := range match {
		dns = append(dns, m[1])
	}

	return h.Success(ctx, dns)
}

// SetDNS 设置 DNS 信息
func (r *ToolBoxController) SetDNS(ctx http.Context) http.Response {
	dns1 := ctx.Request().Input("dns1")
	dns2 := ctx.Request().Input("dns2")

	if len(dns1) == 0 || len(dns2) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "DNS 信息不能为空")
	}

	var dns string
	dns += "nameserver " + dns1 + "\n"
	dns += "nameserver " + dns2 + "\n"

	if err := io.Write("/etc/resolv.conf", dns, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "写入 DNS 信息失败")
	}

	return h.Success(ctx, nil)
}

// GetSWAP 获取 SWAP 信息
func (r *ToolBoxController) GetSWAP(ctx http.Context) http.Response {
	var total, used, free string
	var size int64
	if io.Exists("/www/swap") {
		file, err := io.FileInfo("/www/swap")
		if err != nil {
			return h.Error(ctx, http.StatusUnprocessableEntity, "获取 SWAP 信息失败")
		}

		size = file.Size() / 1024 / 1024
		total = str.FormatBytes(float64(file.Size()))
	} else {
		size = 0
		total = "0.00 B"
	}

	raw, err := shell.Execf("free | grep Swap")
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "获取 SWAP 信息失败")
	}

	match := regexp.MustCompile(`Swap:\s+(\d+)\s+(\d+)\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) > 0 {
		used = str.FormatBytes(cast.ToFloat64(match[2]) * 1024)
		free = str.FormatBytes(cast.ToFloat64(match[3]) * 1024)
	}

	return h.Success(ctx, http.Json{
		"total": total,
		"size":  size,
		"used":  used,
		"free":  free,
	})
}

// SetSWAP 设置 SWAP 信息
func (r *ToolBoxController) SetSWAP(ctx http.Context) http.Response {
	size := ctx.Request().InputInt("size")

	if io.Exists("/www/swap") {
		if out, err := shell.Execf("swapoff /www/swap"); err != nil {
			return h.Error(ctx, http.StatusUnprocessableEntity, out)
		}
		if out, err := shell.Execf("rm -f /www/swap"); err != nil {
			return h.Error(ctx, http.StatusUnprocessableEntity, out)
		}
		if out, err := shell.Execf("sed -i '/www\\/swap/d' /etc/fstab"); err != nil {
			return h.Error(ctx, http.StatusUnprocessableEntity, out)
		}
	}

	if size > 1 {
		free, err := shell.Execf("df -k /www | awk '{print $4}' | tail -n 1")
		if err != nil {
			return h.Error(ctx, http.StatusUnprocessableEntity, "获取磁盘空间失败")
		}
		if cast.ToInt64(free)*1024 < int64(size)*1024*1024 {
			return h.Error(ctx, http.StatusUnprocessableEntity, "磁盘空间不足，当前剩余 "+str.FormatBytes(cast.ToFloat64(free)))
		}

		btrfsCheck, _ := shell.Execf("df -T /www | awk '{print $2}' | tail -n 1")
		if strings.Contains(btrfsCheck, "btrfs") {
			if out, err := shell.Execf("btrfs filesystem mkswapfile --size " + cast.ToString(size) + "M --uuid clear /www/swap"); err != nil {
				return h.Error(ctx, http.StatusUnprocessableEntity, out)
			}
		} else {
			if out, err := shell.Execf("dd if=/dev/zero of=/www/swap bs=1M count=" + cast.ToString(size)); err != nil {
				return h.Error(ctx, http.StatusUnprocessableEntity, out)
			}
			if out, err := shell.Execf("mkswap -f /www/swap"); err != nil {
				return h.Error(ctx, http.StatusUnprocessableEntity, out)
			}
			if err := io.Chmod("/www/swap", 0600); err != nil {
				return h.Error(ctx, http.StatusUnprocessableEntity, "设置 SWAP 权限失败")
			}
		}
		if out, err := shell.Execf("swapon /www/swap"); err != nil {
			return h.Error(ctx, http.StatusUnprocessableEntity, out)
		}
		if out, err := shell.Execf("echo '/www/swap    swap    swap    defaults    0 0' >> /etc/fstab"); err != nil {
			return h.Error(ctx, http.StatusUnprocessableEntity, out)
		}
	}

	return h.Success(ctx, nil)
}

// GetTimezone 获取时区
func (r *ToolBoxController) GetTimezone(ctx http.Context) http.Response {
	raw, err := shell.Execf("timedatectl | grep zone")
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "获取时区信息失败")
	}

	match := regexp.MustCompile(`zone:\s+(\S+)`).FindStringSubmatch(raw)
	if len(match) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "找不到时区信息")
	}

	type zone struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}

	zonesRaw, err := shell.Execf("timedatectl list-timezones")
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "获取时区列表失败")
	}
	zones := strings.Split(zonesRaw, "\n")

	var zonesList []zone
	for _, z := range zones {
		zonesList = append(zonesList, zone{
			Label: z,
			Value: z,
		})
	}

	return h.Success(ctx, http.Json{
		"timezone":  match[1],
		"timezones": zonesList,
	})
}

// SetTimezone 设置时区
func (r *ToolBoxController) SetTimezone(ctx http.Context) http.Response {
	timezone := ctx.Request().Input("timezone")
	if len(timezone) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "时区不能为空")
	}

	if out, err := shell.Execf("timedatectl set-timezone %s", timezone); err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, out)
	}

	return h.Success(ctx, nil)
}

// GetHosts 获取 hosts 信息
func (r *ToolBoxController) GetHosts(ctx http.Context) http.Response {
	hosts, err := io.Read("/etc/hosts")
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}

	return h.Success(ctx, hosts)
}

// SetHosts 设置 hosts 信息
func (r *ToolBoxController) SetHosts(ctx http.Context) http.Response {
	hosts := ctx.Request().Input("hosts")
	if len(hosts) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "hosts 信息不能为空")
	}

	if err := io.Write("/etc/hosts", hosts, 0644); err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "写入 hosts 信息失败")
	}

	return h.Success(ctx, nil)
}

// SetRootPassword 设置 root 密码
func (r *ToolBoxController) SetRootPassword(ctx http.Context) http.Response {
	password := ctx.Request().Input("password")
	if len(password) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "密码不能为空")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9·~!@#$%^&*()_+-=\[\]{};:'",./<>?]{6,20}$`).MatchString(password) {
		return h.Error(ctx, http.StatusUnprocessableEntity, "密码必须为 6-20 位字母、数字或特殊字符")
	}

	password = strings.ReplaceAll(password, `'`, `\'`)
	if out, err := shell.Execf(`yes '` + password + `' | passwd root`); err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, out)
	}

	return h.Success(ctx, nil)
}
