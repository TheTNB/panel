package plugins

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/spf13/cast"

	"panel/app/http/controllers"
	"panel/pkg/tools"
)

type ToolBoxController struct {
}

func NewToolBoxController() *ToolBoxController {
	return &ToolBoxController{}
}

// GetDNS 获取 DNS 信息
func (r *ToolBoxController) GetDNS(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "toolbox")
	if check != nil {
		return check
	}

	raw := tools.Read("/etc/resolv.conf")
	match := regexp.MustCompile(`nameserver\s+(\S+)`).FindAllStringSubmatch(raw, -1)
	if len(match) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "找不到 DNS 信息")
	}

	var dns []string
	for _, m := range match {
		dns = append(dns, m[1])
	}

	return controllers.Success(ctx, dns)
}

// SetDNS 设置 DNS 信息
func (r *ToolBoxController) SetDNS(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "toolbox")
	if check != nil {
		return check
	}

	dns1 := ctx.Request().Input("dns1")
	dns2 := ctx.Request().Input("dns2")

	if len(dns1) == 0 || len(dns2) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "DNS 信息不能为空")
	}

	var dns string
	dns += "nameserver " + dns1 + "\n"
	dns += "nameserver " + dns2 + "\n"

	tools.Write("/etc/resolv.conf", dns, 0644)
	return controllers.Success(ctx, nil)
}

// GetSWAP 获取 SWAP 信息
func (r *ToolBoxController) GetSWAP(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "toolbox")
	if check != nil {
		return check
	}

	var total, used, free string
	var size int64
	if tools.Exists("/www/swap") {
		s, _ := tools.FileSize("/www/swap")
		size = s / 1024 / 1024
		total = tools.FormatBytes(float64(s))
	} else {
		size = 0
		total = "0.00 B"
	}

	raw := tools.Exec("LC_ALL=C free | grep Swap")
	match := regexp.MustCompile(`Swap:\s+(\d+)\s+(\d+)\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) > 0 {
		used = tools.FormatBytes(cast.ToFloat64(match[2]) * 1024)
		free = tools.FormatBytes(cast.ToFloat64(match[3]) * 1024)
	}

	return controllers.Success(ctx, http.Json{
		"total": total,
		"size":  size,
		"used":  used,
		"free":  free,
	})
}

// SetSWAP 设置 SWAP 信息
func (r *ToolBoxController) SetSWAP(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "toolbox")
	if check != nil {
		return check
	}

	size := ctx.Request().InputInt("size")

	if tools.Exists("/www/swap") {
		tools.Exec("swapoff /www/swap")
		tools.Exec("rm -f /www/swap")
		tools.Exec("sed -i '/www\\/swap/d' /etc/fstab")
	}

	if size > 1 {
		free := tools.Exec("df -k /www | awk '{print $4}' | tail -n 1")
		if cast.ToInt64(free)*1024 < int64(size)*1024*1024 {
			return controllers.Error(ctx, http.StatusUnprocessableEntity, "磁盘空间不足，当前剩余 "+tools.FormatBytes(cast.ToFloat64(free)))
		}

		if strings.Contains(tools.Exec("df -T /www | awk '{print $2}' | tail -n 1"), "btrfs") {
			tools.Exec("btrfs filesystem mkswapfile --size " + cast.ToString(size) + "M --uuid clear /www/swap")
		} else {
			tools.Exec("dd if=/dev/zero of=/www/swap bs=1M count=" + cast.ToString(size))
			tools.Exec("mkswap -f /www/swap")
			tools.Chmod("/www/swap", 0600)
		}
		tools.Exec("swapon /www/swap")
		tools.Exec("echo '/www/swap    swap    swap    defaults    0 0' >> /etc/fstab")
	}

	return controllers.Success(ctx, nil)
}

// GetTimezone 获取时区
func (r *ToolBoxController) GetTimezone(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "toolbox")
	if check != nil {
		return check
	}

	raw := tools.Exec("LC_ALL=C timedatectl | grep zone")
	match := regexp.MustCompile(`zone:\s+(\S+)`).FindStringSubmatch(raw)
	if len(match) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "找不到时区信息")
	}

	type zone struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}

	zonesRaw := tools.Exec("LC_ALL=C timedatectl list-timezones")
	zones := strings.Split(zonesRaw, "\n")

	var zonesList []zone
	for _, z := range zones {
		zonesList = append(zonesList, zone{
			Label: z,
			Value: z,
		})
	}

	return controllers.Success(ctx, http.Json{
		"timezone":  match[1],
		"timezones": zonesList,
	})
}

// SetTimezone 设置时区
func (r *ToolBoxController) SetTimezone(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "toolbox")
	if check != nil {
		return check
	}

	timezone := ctx.Request().Input("timezone")
	if len(timezone) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "时区不能为空")
	}

	tools.Exec("timedatectl set-timezone " + timezone)

	return controllers.Success(ctx, nil)
}

// GetHosts 获取 hosts 信息
func (r *ToolBoxController) GetHosts(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "toolbox")
	if check != nil {
		return check
	}

	return controllers.Success(ctx, tools.Read("/etc/hosts"))
}

// SetHosts 设置 hosts 信息
func (r *ToolBoxController) SetHosts(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "toolbox")
	if check != nil {
		return check
	}

	hosts := ctx.Request().Input("hosts")
	if len(hosts) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "hosts 信息不能为空")
	}

	tools.Write("/etc/hosts", hosts, 0644)

	return controllers.Success(ctx, nil)
}

// SetRootPassword 设置 root 密码
func (r *ToolBoxController) SetRootPassword(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "toolbox")
	if check != nil {
		return check
	}

	password := ctx.Request().Input("password")
	if len(password) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "密码不能为空")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9·~!@#$%^&*()_+-=\[\]{};:'",./<>?]{6,20}$`).MatchString(password) {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "密码必须为 6-20 位字母、数字或特殊字符")
	}

	password = strings.ReplaceAll(password, `'`, `\'`)
	tools.Exec(`yes '` + password + `' | passwd root`)

	return controllers.Success(ctx, nil)
}
