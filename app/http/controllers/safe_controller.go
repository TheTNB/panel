package controllers

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/spf13/cast"

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
	if ctx.Request().InputBool("status") {
		if tools.IsRHEL() {
			out = tools.Exec("systemctl start firewalld")
		} else {
			out = tools.Exec("echo y | ufw enable")
		}
	} else {
		if tools.IsRHEL() {
			out = tools.Exec("systemctl stop firewalld")
		} else {
			out = tools.Exec("ufw disable")
		}
	}

	return Success(ctx, out)
}

// GetFirewallRules 获取防火墙规则
func (r *SafeController) GetFirewallRules(ctx http.Context) http.Response {
	if !r.firewallStatus() {
		return Success(ctx, nil)
	}
	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)

	if tools.IsRHEL() {
		out := tools.Exec("firewall-cmd --list-all 2>&1")
		match := regexp.MustCompile(`ports: (.*)`).FindStringSubmatch(out)
		if len(match) == 0 {
			return Success(ctx, http.Json{
				"total": 0,
				"items": []map[string]string{},
			})
		}
		ports := strings.Split(match[1], " ")
		var rules []map[string]string
		for _, port := range ports {
			rule := strings.Split(port, "/")
			rules = append(rules, map[string]string{
				"port":     rule[0],
				"protocol": rule[1],
			})
		}

		startIndex := (page - 1) * limit
		endIndex := page * limit
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
	} else {
		out := tools.Exec("ufw status | grep -v '(v6)' | grep ALLOW | awk '{print $1}'")
		if len(out) == 0 {
			return Success(ctx, http.Json{
				"total": 0,
				"items": []map[string]string{},
			})
		}
		var rules []map[string]string
		for _, port := range strings.Split(out, "\n") {
			rule := strings.Split(port, "/")
			rules = append(rules, map[string]string{
				"port":     rule[0],
				"protocol": rule[1],
			})
		}

		startIndex := (page - 1) * limit
		endIndex := page * limit
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
}

// AddFirewallRule 添加防火墙规则
func (r *SafeController) AddFirewallRule(ctx http.Context) http.Response {
	if !r.firewallStatus() {
		return Error(ctx, http.StatusBadRequest, "防火墙未启动")
	}

	port := ctx.Request().InputInt("port", 0)
	protocol := ctx.Request().Input("protocol", "")
	if port == 0 || protocol == "" {
		return Error(ctx, http.StatusBadRequest, "参数错误")
	}

	if tools.IsRHEL() {
		tools.Exec("firewall-cmd --remove-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1")
		tools.Exec("firewall-cmd --add-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1")
		tools.Exec("firewall-cmd --reload")
	} else {
		tools.Exec("ufw delete allow " + cast.ToString(port) + "/" + protocol)
		tools.Exec("ufw allow " + cast.ToString(port) + "/" + protocol)
		tools.Exec("ufw reload")
	}

	return Success(ctx, nil)
}

// DeleteFirewallRule 删除防火墙规则
func (r *SafeController) DeleteFirewallRule(ctx http.Context) http.Response {
	if !r.firewallStatus() {
		return Error(ctx, http.StatusBadRequest, "防火墙未启动")
	}

	port := ctx.Request().InputInt("port", 0)
	protocol := ctx.Request().Input("protocol", "")
	if port == 0 || protocol == "" {
		return Error(ctx, http.StatusBadRequest, "参数错误")
	}

	if tools.IsRHEL() {
		tools.Exec("firewall-cmd --remove-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1")
		tools.Exec("firewall-cmd --reload")
	} else {
		tools.Exec("ufw delete allow " + cast.ToString(port) + "/" + protocol)
		tools.Exec("ufw reload")
	}

	return Success(ctx, nil)
}

// firewallStatus 获取防火墙状态
func (r *SafeController) firewallStatus() bool {
	var out string
	var running bool
	if tools.IsRHEL() {
		out = tools.Exec("systemctl status firewalld | grep Active | awk '{print $3}'")
		if out == "(running)" {
			running = true
		} else {
			running = false
		}
	} else {
		out = tools.Exec("ufw status | grep Status | awk '{print $2}'")
		if out == "active" {
			running = true
		} else {
			running = false
		}
	}

	return running
}

// GetSshStatus 获取 SSH 状态
func (r *SafeController) GetSshStatus(ctx http.Context) http.Response {
	var out string
	if tools.IsRHEL() {
		out = tools.Exec("systemctl status sshd | grep Active | awk '{print $3}'")
	} else {
		out = tools.Exec("systemctl status ssh | grep Active | awk '{print $3}'")
	}

	running := false
	if out == "(running)" {
		running = true
	}

	return Success(ctx, running)
}

// SetSshStatus 设置 SSH 状态
func (r *SafeController) SetSshStatus(ctx http.Context) http.Response {
	if ctx.Request().InputBool("status") {
		if tools.IsRHEL() {
			tools.Exec("systemctl enable sshd")
			tools.Exec("systemctl start sshd")
		} else {
			tools.Exec("systemctl enable ssh")
			tools.Exec("systemctl start ssh")
		}
	} else {
		if tools.IsRHEL() {
			tools.Exec("systemctl stop sshd")
			tools.Exec("systemctl disable sshd")
		} else {
			tools.Exec("systemctl stop ssh")
			tools.Exec("systemctl disable ssh")
		}
	}

	return Success(ctx, nil)
}

// GetSshPort 获取 SSH 端口
func (r *SafeController) GetSshPort(ctx http.Context) http.Response {
	out := tools.Exec("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'")
	return Success(ctx, out)
}

// SetSshPort 设置 SSH 端口
func (r *SafeController) SetSshPort(ctx http.Context) http.Response {
	port := ctx.Request().InputInt("port", 0)
	if port == 0 {
		return Error(ctx, http.StatusBadRequest, "参数错误")
	}

	oldPort := tools.Exec("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'")
	tools.Exec("sed -i 's/#Port " + oldPort + "/Port " + cast.ToString(port) + "/g' /etc/ssh/sshd_config")
	tools.Exec("sed -i 's/Port " + oldPort + "/Port " + cast.ToString(port) + "/g' /etc/ssh/sshd_config")

	if status := tools.Exec("systemctl status sshd | grep Active | awk '{print $3}'"); status == "(running)" {
		tools.Exec("systemctl restart sshd")
	}

	return Success(ctx, nil)
}

// GetPingStatus 获取 Ping 状态
func (r *SafeController) GetPingStatus(ctx http.Context) http.Response {
	if tools.IsRHEL() {
		out := tools.Exec(`firewall-cmd --list-all 2>&1`)
		if !strings.Contains(out, `rule protocol value="icmp" drop`) {
			return Success(ctx, true)
		} else {
			return Success(ctx, false)
		}
	} else {
		config := tools.Read("/etc/ufw/before.rules")
		if strings.Contains(config, "-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT") {
			return Success(ctx, true)
		} else {
			return Success(ctx, false)
		}
	}
}

// SetPingStatus 设置 Ping 状态
func (r *SafeController) SetPingStatus(ctx http.Context) http.Response {
	if tools.IsRHEL() {
		if ctx.Request().InputBool("status") {
			tools.Exec(`firewall-cmd --permanent --remove-rich-rule='rule protocol value=icmp drop'`)
		} else {
			tools.Exec(`firewall-cmd --permanent --add-rich-rule='rule protocol value=icmp drop'`)
		}
		tools.Exec(`firewall-cmd --reload`)
	} else {
		if ctx.Request().InputBool("status") {
			tools.Exec(`sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/g' /etc/ufw/before.rules`)
		} else {
			tools.Exec(`sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/g' /etc/ufw/before.rules`)
		}
		tools.Exec(`ufw reload`)
	}

	return Success(ctx, nil)
}
