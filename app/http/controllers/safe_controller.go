package controllers

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/spf13/cast"
	commonrequests "panel/app/http/requests/common"

	"panel/pkg/tools"
)

type SafeController struct {
	// Dependent services
}

func NewSafeController() *SafeController {
	return &SafeController{
		// Inject services
	}
}

// GetFirewallStatus 获取防火墙状态
func (r *SafeController) GetFirewallStatus(ctx http.Context) http.Response {
	return Success(ctx, r.firewallStatus())
}

// SetFirewallStatus 设置防火墙状态
func (r *SafeController) SetFirewallStatus(ctx http.Context) http.Response {
	var out string
	var err error
	if ctx.Request().InputBool("status") {
		if tools.IsRHEL() {
			out, err = tools.Exec("systemctl start firewalld")
		} else {
			out, err = tools.Exec("echo y | ufw enable")
		}
	} else {
		if tools.IsRHEL() {
			out, err = tools.Exec("systemctl stop firewalld")
		} else {
			out, err = tools.Exec("ufw disable")
		}
	}

	if err != nil {
		return Error(ctx, http.StatusInternalServerError, out)
	}

	return Success(ctx, nil)
}

// GetFirewallRules 获取防火墙规则
func (r *SafeController) GetFirewallRules(ctx http.Context) http.Response {
	var paginateRequest commonrequests.Paginate
	sanitize := Sanitize(ctx, &paginateRequest)
	if sanitize != nil {
		return sanitize
	}

	if !r.firewallStatus() {
		return Success(ctx, nil)
	}

	var rules []map[string]string
	if tools.IsRHEL() {
		out, err := tools.Exec("firewall-cmd --list-all 2>&1")
		if err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}

		match := regexp.MustCompile(`ports: (.*)`).FindStringSubmatch(out)
		if len(match) == 0 {
			return Success(ctx, http.Json{
				"total": 0,
				"items": []map[string]string{},
			})
		}
		ports := strings.Split(match[1], " ")
		for _, port := range ports {
			rule := strings.Split(port, "/")
			rules = append(rules, map[string]string{
				"port":     rule[0],
				"protocol": rule[1],
			})
		}
	} else {
		out, err := tools.Exec("ufw status | grep -v '(v6)' | grep ALLOW | awk '{print $1}'")
		if err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}

		if len(out) == 0 {
			return Success(ctx, http.Json{
				"total": 0,
				"items": []map[string]string{},
			})
		}
		for _, port := range strings.Split(out, "\n") {
			rule := strings.Split(port, "/")
			rules = append(rules, map[string]string{
				"port":     rule[0],
				"protocol": rule[1],
			})
		}
	}

	startIndex := (paginateRequest.Page - 1) * paginateRequest.Limit
	endIndex := paginateRequest.Page * paginateRequest.Limit
	if startIndex > len(rules) {
		return Success(ctx, http.Json{
			"total": 0,
			"items": []map[string]string{},
		})
	}
	if endIndex > len(rules) {
		endIndex = len(rules)
	}
	pagedRules := rules[startIndex:endIndex]

	return Success(ctx, http.Json{
		"total": len(rules),
		"items": pagedRules,
	})
}

// AddFirewallRule 添加防火墙规则
func (r *SafeController) AddFirewallRule(ctx http.Context) http.Response {
	if !r.firewallStatus() {
		return Error(ctx, http.StatusInternalServerError, "防火墙未启动")
	}

	port := ctx.Request().Input("port")
	protocol := ctx.Request().Input("protocol")
	if port == "" || protocol == "" || (protocol != "tcp" && protocol != "udp") {
		return Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}
	// 端口有 2 种写法，一种是 80-443，一种是 80
	if strings.Contains(port, "-") {
		ports := strings.Split(port, "-")
		startPort := cast.ToInt(ports[0])
		endPort := cast.ToInt(ports[1])
		if startPort < 1 || startPort > 65535 || endPort < 1 || endPort > 65535 || startPort > endPort {
			return Error(ctx, http.StatusUnprocessableEntity, "参数错误")
		}
	} else {
		port := cast.ToInt(port)
		if port < 1 || port > 65535 {
			return Error(ctx, http.StatusUnprocessableEntity, "参数错误")
		}
	}

	if tools.IsRHEL() {
		if out, err := tools.Exec("firewall-cmd --remove-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1"); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := tools.Exec("firewall-cmd --add-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1"); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := tools.Exec("firewall-cmd --reload"); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
	} else {
		// ufw 需要替换 - 为 : 添加
		if strings.Contains(port, "-") {
			port = strings.ReplaceAll(port, "-", ":")
		}
		if out, err := tools.Exec("ufw delete allow " + cast.ToString(port) + "/" + protocol); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := tools.Exec("ufw allow " + cast.ToString(port) + "/" + protocol); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := tools.Exec("ufw reload"); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
	}

	return Success(ctx, nil)
}

// DeleteFirewallRule 删除防火墙规则
func (r *SafeController) DeleteFirewallRule(ctx http.Context) http.Response {
	if !r.firewallStatus() {
		return Error(ctx, http.StatusUnprocessableEntity, "防火墙未启动")
	}

	port := ctx.Request().InputInt("port", 0)
	protocol := ctx.Request().Input("protocol", "")
	if port == 0 || protocol == "" {
		return Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}

	if tools.IsRHEL() {
		if out, err := tools.Exec("firewall-cmd --remove-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1"); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := tools.Exec("firewall-cmd --reload"); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
	} else {
		if out, err := tools.Exec("ufw delete allow " + cast.ToString(port) + "/" + protocol); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := tools.Exec("ufw reload"); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
	}

	return Success(ctx, nil)
}

// firewallStatus 获取防火墙状态
func (r *SafeController) firewallStatus() bool {
	var out string
	var err error
	var running bool
	if tools.IsRHEL() {
		out, err = tools.Exec("systemctl status firewalld | grep Active | awk '{print $3}'")
		if out == "(running)" {
			running = true
		} else {
			running = false
		}
	} else {
		out, err = tools.Exec("ufw status | grep Status | awk '{print $2}'")
		if out == "active" {
			running = true
		} else {
			running = false
		}
	}

	if err != nil {
		return false
	}

	return running
}

// GetSshStatus 获取 SSH 状态
func (r *SafeController) GetSshStatus(ctx http.Context) http.Response {
	var out string
	var err error
	if tools.IsRHEL() {
		out, err = tools.Exec("systemctl status sshd | grep Active | awk '{print $3}'")
	} else {
		out, err = tools.Exec("systemctl status ssh | grep Active | awk '{print $3}'")
	}

	if err != nil {
		return Error(ctx, http.StatusInternalServerError, out)
	}

	return Success(ctx, out == "(running)")
}

// SetSshStatus 设置 SSH 状态
func (r *SafeController) SetSshStatus(ctx http.Context) http.Response {
	if ctx.Request().InputBool("status") {
		if tools.IsRHEL() {
			if out, err := tools.Exec("systemctl enable sshd"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := tools.Exec("systemctl start sshd"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
		} else {
			if out, err := tools.Exec("systemctl enable ssh"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := tools.Exec("systemctl start ssh"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
		}
	} else {
		if tools.IsRHEL() {
			if out, err := tools.Exec("systemctl stop sshd"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := tools.Exec("systemctl disable sshd"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
		} else {
			if out, err := tools.Exec("systemctl stop ssh"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := tools.Exec("systemctl disable ssh"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
		}
	}

	return Success(ctx, nil)
}

// GetSshPort 获取 SSH 端口
func (r *SafeController) GetSshPort(ctx http.Context) http.Response {
	out, err := tools.Exec("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'")
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, out)
	}

	return Success(ctx, out)
}

// SetSshPort 设置 SSH 端口
func (r *SafeController) SetSshPort(ctx http.Context) http.Response {
	port := ctx.Request().InputInt("port", 0)
	if port == 0 {
		return Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}

	oldPort, err := tools.Exec("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'")
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, oldPort)
	}
	_, _ = tools.Exec("sed -i 's/#Port " + oldPort + "/Port " + cast.ToString(port) + "/g' /etc/ssh/sshd_config")
	_, _ = tools.Exec("sed -i 's/Port " + oldPort + "/Port " + cast.ToString(port) + "/g' /etc/ssh/sshd_config")

	out, err := tools.Exec("systemctl status sshd | grep Active | awk '{print $3}'")
	if err != nil || out != "(running)" {
		Error(ctx, http.StatusInternalServerError, out)
	}

	_, _ = tools.Exec("systemctl restart sshd")
	return Success(ctx, nil)
}

// GetPingStatus 获取 Ping 状态
func (r *SafeController) GetPingStatus(ctx http.Context) http.Response {
	if tools.IsRHEL() {
		out, err := tools.Exec(`firewall-cmd --list-all 2>&1`)
		if err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}

		if !strings.Contains(out, `rule protocol value="icmp" drop`) {
			return Success(ctx, true)
		} else {
			return Success(ctx, false)
		}
	} else {
		config, err := tools.Read("/etc/ufw/before.rules")
		if err != nil {
			return Error(ctx, http.StatusInternalServerError, err.Error())
		}
		if strings.Contains(config, "-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT") {
			return Success(ctx, true)
		} else {
			return Success(ctx, false)
		}
	}
}

// SetPingStatus 设置 Ping 状态
func (r *SafeController) SetPingStatus(ctx http.Context) http.Response {
	var out string
	var err error
	if tools.IsRHEL() {
		if ctx.Request().InputBool("status") {
			out, err = tools.Exec(`firewall-cmd --permanent --remove-rich-rule='rule protocol value=icmp drop'`)
		} else {
			out, err = tools.Exec(`firewall-cmd --permanent --add-rich-rule='rule protocol value=icmp drop'`)
		}
	} else {
		if ctx.Request().InputBool("status") {
			out, err = tools.Exec(`sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/g' /etc/ufw/before.rules`)
		} else {
			out, err = tools.Exec(`sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/g' /etc/ufw/before.rules`)
		}
	}

	if err != nil {
		return Error(ctx, http.StatusInternalServerError, out)
	}

	if tools.IsRHEL() {
		out, err = tools.Exec(`firewall-cmd --reload`)
	} else {
		out, err = tools.Exec(`ufw reload`)
	}

	if err != nil {
		return Error(ctx, http.StatusInternalServerError, out)
	}

	return Success(ctx, nil)
}
