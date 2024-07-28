package plugins

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/os"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/str"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type Fail2banController struct {
	website internal.Website
}

func NewFail2banController() *Fail2banController {
	return &Fail2banController{
		website: services.NewWebsiteImpl(),
	}
}

// List 所有 Fail2ban 规则
func (r *Fail2banController) List(ctx http.Context) http.Response {
	raw, err := io.Read("/etc/fail2ban/jail.local")
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}

	jailList := regexp.MustCompile(`\[(.*?)]`).FindAllStringSubmatch(raw, -1)
	if len(jailList) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "Fail2ban 规则为空")
	}

	var jails []types.Fail2banJail
	for i, jail := range jailList {
		if i == 0 {
			continue
		}

		jailName := jail[1]
		jailRaw := str.Cut(raw, "# "+jailName+"-START", "# "+jailName+"-END")
		if len(jailRaw) == 0 {
			continue
		}
		jailEnabled := strings.Contains(jailRaw, "enabled = true")
		jailLogPath := regexp.MustCompile(`logpath = (.*)`).FindStringSubmatch(jailRaw)
		jailMaxRetry := regexp.MustCompile(`maxretry = (.*)`).FindStringSubmatch(jailRaw)
		jailFindTime := regexp.MustCompile(`findtime = (.*)`).FindStringSubmatch(jailRaw)
		jailBanTime := regexp.MustCompile(`bantime = (.*)`).FindStringSubmatch(jailRaw)

		jails = append(jails, types.Fail2banJail{
			Name:     jailName,
			Enabled:  jailEnabled,
			LogPath:  jailLogPath[1],
			MaxRetry: cast.ToInt(jailMaxRetry[1]),
			FindTime: cast.ToInt(jailFindTime[1]),
			BanTime:  cast.ToInt(jailBanTime[1]),
		})
	}

	paged, total := h.Paginate(ctx, jails)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// Add 添加 Fail2ban 规则
func (r *Fail2banController) Add(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"name":         "required",
		"type":         "required|in:website,service",
		"maxretry":     "required",
		"findtime":     "required",
		"bantime":      "required",
		"website_name": "required_if:type,website",
		"website_mode": "required_if:type,website",
		"website_path": "required_if:website_mode,path",
	}); sanitize != nil {
		return sanitize
	}

	jailName := ctx.Request().Input("name")
	jailType := ctx.Request().Input("type")
	jailMaxRetry := ctx.Request().Input("maxretry")
	jailFindTime := ctx.Request().Input("findtime")
	jailBanTime := ctx.Request().Input("bantime")
	jailWebsiteName := ctx.Request().Input("website_name")
	jailWebsiteMode := ctx.Request().Input("website_mode")
	jailWebsitePath := ctx.Request().Input("website_path")

	raw, err := io.Read("/etc/fail2ban/jail.local")
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if (strings.Contains(raw, "["+jailName+"]") && jailType == "service") || (strings.Contains(raw, "["+jailWebsiteName+"]"+"-cc") && jailType == "website" && jailWebsiteMode == "cc") || (strings.Contains(raw, "["+jailWebsiteName+"]"+"-path") && jailType == "website" && jailWebsiteMode == "path") {
		return h.Error(ctx, http.StatusUnprocessableEntity, "规则已存在")
	}

	switch jailType {
	case "website":
		var website models.Website
		err := facades.Orm().Query().Where("name", jailWebsiteName).FirstOrFail(&website)
		if err != nil {
			return h.Error(ctx, http.StatusUnprocessableEntity, "网站不存在")
		}
		config, err := r.website.GetConfig(website.ID)
		if err != nil {
			return h.Error(ctx, http.StatusUnprocessableEntity, "获取网站配置失败")
		}
		var ports string
		for _, port := range config.Ports {
			fields := strings.Fields(cast.ToString(port))
			ports += fields[0] + ","
		}

		rule := `
# ` + jailWebsiteName + `-` + jailWebsiteMode + `-START
[` + jailWebsiteName + `-` + jailWebsiteMode + `]
enabled = true
filter = haozi-` + jailWebsiteName + `-` + jailWebsiteMode + `
port = ` + ports + `
maxretry = ` + jailMaxRetry + `
findtime = ` + jailFindTime + `
bantime = ` + jailBanTime + `
action = %(action_mwl)s
logpath = /www/wwwlogs/` + website.Name + `.log
# ` + jailWebsiteName + `-` + jailWebsiteMode + `-END
`
		raw += rule
		if err = io.Write("/etc/fail2ban/jail.local", raw, 0644); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, "写入Fail2ban规则失败")
		}

		var filter string
		if jailWebsiteMode == "cc" {
			filter = `
[Definition]
failregex = ^<HOST>\s-.*HTTP/.*$
ignoreregex =
`
		} else {
			filter = `
[Definition]
failregex = ^<HOST>\s-.*\s` + jailWebsitePath + `.*HTTP/.*$
ignoreregex =
`
		}
		if err = io.Write("/etc/fail2ban/filter.d/haozi-"+jailWebsiteName+"-"+jailWebsiteMode+".conf", filter, 0644); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, "写入Fail2ban规则失败")
		}

	case "service":
		var logPath string
		var filter string
		var port string
		var err error
		switch jailName {
		case "ssh":
			if os.IsDebian() || os.IsUbuntu() {
				logPath = "/var/log/auth.log"
			} else {
				logPath = "/var/log/secure"
			}
			filter = "sshd"
			port, err = shell.Execf("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'")
		case "mysql":
			logPath = "/www/server/mysql/mysql-error.log"
			filter = "mysqld-auth"
			port, err = shell.Execf("cat /www/server/mysql/conf/my.cnf | grep 'port' | head -n 1 | awk '{print $3}'")
		case "pure-ftpd":
			logPath = "/var/log/messages"
			filter = "pure-ftpd"
			port, err = shell.Execf(`cat /www/server/pure-ftpd/etc/pure-ftpd.conf | grep "Bind" | awk '{print $2}' | awk -F "," '{print $2}'`)
		default:
			return h.Error(ctx, http.StatusUnprocessableEntity, "未知服务")
		}
		if len(port) == 0 || err != nil {
			return h.Error(ctx, http.StatusUnprocessableEntity, "获取服务端口失败，请检查是否安装")
		}

		rule := `
# ` + jailName + `-START
[` + jailName + `]
enabled = true
filter = ` + filter + `
port = ` + port + `
maxretry = ` + jailMaxRetry + `
findtime = ` + jailFindTime + `
bantime = ` + jailBanTime + `
action = %(action_mwl)s
logpath = ` + logPath + `
# ` + jailName + `-END
`
		raw += rule
		if err := io.Write("/etc/fail2ban/jail.local", raw, 0644); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, "写入Fail2ban规则失败")
		}
	}

	if _, err := shell.Execf("fail2ban-client reload"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "重载配置失败")
	}

	return h.Success(ctx, nil)
}

// Delete 删除规则
func (r *Fail2banController) Delete(ctx http.Context) http.Response {
	jailName := ctx.Request().Input("name")
	raw, err := io.Read("/etc/fail2ban/jail.local")
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if !strings.Contains(raw, "["+jailName+"]") {
		return h.Error(ctx, http.StatusUnprocessableEntity, "规则不存在")
	}

	rule := str.Cut(raw, "# "+jailName+"-START", "# "+jailName+"-END")
	raw = strings.Replace(raw, "\n# "+jailName+"-START"+rule+"# "+jailName+"-END", "", -1)
	raw = strings.TrimSpace(raw)
	if err := io.Write("/etc/fail2ban/jail.local", raw, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "写入Fail2ban规则失败")
	}

	if _, err := shell.Execf("fail2ban-client reload"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "重载配置失败")
	}

	return h.Success(ctx, nil)
}

// BanList 获取封禁列表
func (r *Fail2banController) BanList(ctx http.Context) http.Response {
	name := ctx.Request().Input("name")
	if len(name) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "缺少参数")
	}

	currentlyBan, err := shell.Execf(`fail2ban-client status %s | grep "Currently banned" | awk '{print $4}'`, name)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取封禁列表失败")
	}
	totalBan, err := shell.Execf(`fail2ban-client status %s | grep "Total banned" | awk '{print $4}'`, name)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取封禁列表失败")
	}
	bannedIp, err := shell.Execf(`fail2ban-client status %s | grep "Banned IP list" | awk -F ":" '{print $2}'`, name)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取封禁列表失败")
	}
	bannedIpList := strings.Split(bannedIp, " ")

	var list []map[string]string
	for _, ip := range bannedIpList {
		if len(ip) > 0 {
			list = append(list, map[string]string{
				"name": name,
				"ip":   ip,
			})
		}
	}
	if list == nil {
		list = []map[string]string{}
	}

	return h.Success(ctx, http.Json{
		"currently_ban": currentlyBan,
		"total_ban":     totalBan,
		"baned_list":    list,
	})
}

// Unban 解封
func (r *Fail2banController) Unban(ctx http.Context) http.Response {
	name := ctx.Request().Input("name")
	ip := ctx.Request().Input("ip")
	if len(name) == 0 || len(ip) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "缺少参数")
	}

	if _, err := shell.Execf("fail2ban-client set %s unbanip %s", name, ip); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "解封失败")
	}

	return h.Success(ctx, nil)
}

// SetWhiteList 设置白名单
func (r *Fail2banController) SetWhiteList(ctx http.Context) http.Response {
	ip := ctx.Request().Input("ip")
	if len(ip) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "缺少参数")
	}

	raw, err := io.Read("/etc/fail2ban/jail.local")
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	// 正则替换
	reg := regexp.MustCompile(`ignoreip\s*=\s*.*\n`)
	if reg.MatchString(raw) {
		raw = reg.ReplaceAllString(raw, "ignoreip = "+ip+"\n")
	} else {
		return h.Error(ctx, http.StatusInternalServerError, "解析Fail2ban规则失败，Fail2ban可能已损坏")
	}

	if err := io.Write("/etc/fail2ban/jail.local", raw, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "写入Fail2ban规则失败")
	}

	if _, err := shell.Execf("fail2ban-client reload"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "重载配置失败")
	}
	return h.Success(ctx, nil)
}

// GetWhiteList 获取白名单
func (r *Fail2banController) GetWhiteList(ctx http.Context) http.Response {
	raw, err := io.Read("/etc/fail2ban/jail.local")
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	reg := regexp.MustCompile(`ignoreip\s*=\s*(.*)\n`)
	if reg.MatchString(raw) {
		ignoreIp := reg.FindStringSubmatch(raw)[1]
		return h.Success(ctx, ignoreIp)
	} else {
		return h.Error(ctx, http.StatusInternalServerError, "解析Fail2ban规则失败，Fail2ban可能已损坏")
	}
}
