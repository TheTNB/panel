package fail2ban

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/go-rat/chix"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/os"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/str"
)

type Service struct {
	websiteRepo biz.WebsiteRepo
}

func NewService() *Service {
	return &Service{
		websiteRepo: data.NewWebsiteRepo(),
	}
}

// List 所有规则
func (s *Service) List(w http.ResponseWriter, r *http.Request) {
	raw, err := io.Read("/etc/fail2ban/jail.local")
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	jailList := regexp.MustCompile(`\[(.*?)]`).FindAllStringSubmatch(raw, -1)
	if len(jailList) == 0 {
		service.Error(w, http.StatusUnprocessableEntity, "Fail2ban 规则为空")
		return
	}

	var jails []Jail
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

		jails = append(jails, Jail{
			Name:     jailName,
			Enabled:  jailEnabled,
			LogPath:  jailLogPath[1],
			MaxRetry: cast.ToInt(jailMaxRetry[1]),
			FindTime: cast.ToInt(jailFindTime[1]),
			BanTime:  cast.ToInt(jailBanTime[1]),
		})
	}

	paged, total := service.Paginate(r, jails)

	service.Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

// Create 添加规则
func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Add](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}
	jailName := req.Name
	jailType := req.Type
	jailMaxRetry := cast.ToString(req.MaxRetry)
	jailFindTime := cast.ToString(req.FindTime)
	jailBanTime := cast.ToString(req.BanTime)
	jailWebsiteName := req.WebsiteName
	jailWebsiteMode := req.WebsiteMode
	jailWebsitePath := req.WebsitePath

	raw, err := io.Read("/etc/fail2ban/jail.local")
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}
	if (strings.Contains(raw, "["+jailName+"]") && jailType == "service") || (strings.Contains(raw, "["+jailWebsiteName+"]"+"-cc") && jailType == "website" && jailWebsiteMode == "cc") || (strings.Contains(raw, "["+jailWebsiteName+"]"+"-path") && jailType == "website" && jailWebsiteMode == "path") {
		service.Error(w, http.StatusUnprocessableEntity, "规则已存在")
		return
	}

	switch jailType {
	case "website":
		website, err := s.websiteRepo.GetByName(jailWebsiteName)
		if err != nil {
			service.Error(w, http.StatusUnprocessableEntity, "获取网站配置失败：%v", err)
			return
		}
		var ports string
		for _, port := range website.Ports {
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
logpath = ` + app.Root + `/wwwlogs/` + website.Name + `.log
# ` + jailWebsiteName + `-` + jailWebsiteMode + `-END
`
		raw += rule
		if err = io.Write("/etc/fail2ban/jail.local", raw, 0644); err != nil {
			service.Error(w, http.StatusInternalServerError, "写入Fail2ban规则失败")
			return
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
			service.Error(w, http.StatusInternalServerError, "写入Fail2ban规则失败")
			return
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
			logPath = app.Root + "/server/mysql/mysql-error.log"
			filter = "mysqld-auth"
			port, err = shell.Execf("cat %s/server/mysql/conf/my.cnf | grep 'port' | head -n 1 | awk '{print $3}'", app.Root)
		case "pure-ftpd":
			logPath = "/var/log/messages"
			filter = "pure-ftpd"
			port, err = shell.Execf(`cat %s/server/pure-ftpd/etc/pure-ftpd.conf | grep "Bind" | awk '{print $2}' | awk -F "," '{print $2}'`, app.Root)
		default:
			service.Error(w, http.StatusUnprocessableEntity, "未知服务")
			return
		}
		if len(port) == 0 || err != nil {
			service.Error(w, http.StatusUnprocessableEntity, "获取服务端口失败，请检查是否安装")
			return
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
			service.Error(w, http.StatusInternalServerError, "写入Fail2ban规则失败")
			return
		}
	}

	if _, err := shell.Execf("fail2ban-client reload"); err != nil {
		service.Error(w, http.StatusInternalServerError, "重载配置失败")
		return
	}

	service.Success(w, nil)
}

// Delete 删除规则
func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Delete](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	raw, err := io.Read("/etc/fail2ban/jail.local")
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}
	if !strings.Contains(raw, "["+req.Name+"]") {
		service.Error(w, http.StatusUnprocessableEntity, "规则不存在")
		return
	}

	rule := str.Cut(raw, "# "+req.Name+"-START", "# "+req.Name+"-END")
	raw = strings.Replace(raw, "\n# "+req.Name+"-START"+rule+"# "+req.Name+"-END", "", -1)
	raw = strings.TrimSpace(raw)
	if err := io.Write("/etc/fail2ban/jail.local", raw, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "写入Fail2ban规则失败")
		return
	}

	if _, err := shell.Execf("fail2ban-client reload"); err != nil {
		service.Error(w, http.StatusInternalServerError, "重载配置失败")
		return
	}

	service.Success(w, nil)
}

// BanList 获取封禁列表
func (s *Service) BanList(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[BanList](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	currentlyBan, err := shell.Execf(`fail2ban-client status %s | grep "Currently banned" | awk '{print $4}'`, req.Name)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取封禁列表失败")
		return
	}
	totalBan, err := shell.Execf(`fail2ban-client status %s | grep "Total banned" | awk '{print $4}'`, req.Name)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取封禁列表失败")
		return
	}
	bannedIp, err := shell.Execf(`fail2ban-client status %s | grep "Banned IP list" | awk -F ":" '{print $2}'`, req.Name)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取封禁列表失败")
		return
	}
	bannedIpList := strings.Split(bannedIp, " ")

	var list []map[string]string
	for _, ip := range bannedIpList {
		if len(ip) > 0 {
			list = append(list, map[string]string{
				"name": req.Name,
				"ip":   ip,
			})
		}
	}
	if list == nil {
		list = []map[string]string{}
	}

	service.Success(w, chix.M{
		"currently_ban": currentlyBan,
		"total_ban":     totalBan,
		"baned_list":    list,
	})
}

// Unban 解封
func (s *Service) Unban(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Unban](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if _, err = shell.Execf("fail2ban-client set %s unbanip %s", req.Name, req.IP); err != nil {
		service.Error(w, http.StatusInternalServerError, "解封失败")
		return
	}

	service.Success(w, nil)
}

// SetWhiteList 设置白名单
func (s *Service) SetWhiteList(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[SetWhiteList](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	raw, err := io.Read("/etc/fail2ban/jail.local")
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}
	// 正则替换
	reg := regexp.MustCompile(`ignoreip\s*=\s*.*\n`)
	if reg.MatchString(raw) {
		raw = reg.ReplaceAllString(raw, "ignoreip = "+req.IP+"\n")
	} else {
		service.Error(w, http.StatusInternalServerError, "解析Fail2ban规则失败，Fail2ban可能已损坏")
		return
	}

	if err = io.Write("/etc/fail2ban/jail.local", raw, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "写入Fail2ban规则失败")
		return
	}

	if _, err = shell.Execf("fail2ban-client reload"); err != nil {
		service.Error(w, http.StatusInternalServerError, "重载配置失败")
		return
	}
	service.Success(w, nil)
}

// GetWhiteList 获取白名单
func (s *Service) GetWhiteList(w http.ResponseWriter, r *http.Request) {
	raw, err := io.Read("/etc/fail2ban/jail.local")
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}
	reg := regexp.MustCompile(`ignoreip\s*=\s*(.*)\n`)
	if reg.MatchString(raw) {
		ignoreIp := reg.FindStringSubmatch(raw)[1]
		service.Success(w, ignoreIp)
	} else {
		service.Error(w, http.StatusInternalServerError, "解析Fail2ban规则失败，Fail2ban可能已损坏")
		return
	}
}
