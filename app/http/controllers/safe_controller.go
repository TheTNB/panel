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

func (r *SafeController) GetFirewallStatus(ctx http.Context) {
	Success(ctx, r.firewallStatus())
}

func (r *SafeController) SetFirewallStatus(ctx http.Context) {
	var out string
	if ctx.Request().InputBool("status") {
		if tools.IsRHEL() {
			out = tools.ExecShell("systemctl start firewalld")
		} else {
			out = tools.ExecShell("echo y | ufw enable")
		}
	} else {
		if tools.IsRHEL() {
			out = tools.ExecShell("systemctl stop firewalld")
		} else {
			out = tools.ExecShell("ufw disable")
		}
	}

	Success(ctx, out)
}

func (r *SafeController) GetFirewallRules(ctx http.Context) {
	if !r.firewallStatus() {
		Success(ctx, nil)
		return
	}
	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)

	if tools.IsRHEL() {
		out := tools.ExecShell("firewall-cmd --list-all 2>&1")
		match := regexp.MustCompile(`ports: (.*)`).FindStringSubmatch(out)
		if len(match) == 0 {
			Success(ctx, nil)
			return
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
			Success(ctx, http.Json{
				"total": 0,
				"items": []map[string]string{},
			})
			return
		}
		if endIndex > len(rules) {
			endIndex = len(rules)
		}
		pagedRules := rules[startIndex:endIndex]

		Success(ctx, http.Json{
			"total": len(rules),
			"items": pagedRules,
		})
	} else {
		out := tools.ExecShell("ufw status | grep -v '(v6)' | grep ALLOW | awk '{print $1}'")
		if len(out) == 0 {
			Success(ctx, nil)
			return
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
			Success(ctx, http.Json{
				"total": 0,
				"items": []map[string]string{},
			})
			return
		}
		if endIndex > len(rules) {
			endIndex = len(rules)
		}
		pagedRules := rules[startIndex:endIndex]

		Success(ctx, http.Json{
			"total": len(rules),
			"items": pagedRules,
		})
	}
}

func (r *SafeController) AddFirewallRule(ctx http.Context) {
	if !r.firewallStatus() {
		Error(ctx, http.StatusBadRequest, "防火墙未启动")
		return
	}

	port := ctx.Request().InputInt("port", 0)
	protocol := ctx.Request().Input("protocol", "")
	if port == 0 || protocol == "" {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	if tools.IsRHEL() {
		tools.ExecShell("firewall-cmd --remove-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1")
		tools.ExecShell("firewall-cmd --add-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1")
		tools.ExecShell("firewall-cmd --reload")
	} else {
		tools.ExecShell("ufw delete allow " + cast.ToString(port) + "/" + protocol)
		tools.ExecShell("ufw allow " + cast.ToString(port) + "/" + protocol)
		tools.ExecShell("ufw reload")
	}

	Success(ctx, nil)
}

func (r *SafeController) DeleteFirewallRule(ctx http.Context) {
	if !r.firewallStatus() {
		Error(ctx, http.StatusBadRequest, "防火墙未启动")
		return
	}

	port := ctx.Request().InputInt("port", 0)
	protocol := ctx.Request().Input("protocol", "")
	if port == 0 || protocol == "" {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	if tools.IsRHEL() {
		tools.ExecShell("firewall-cmd --remove-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1")
		tools.ExecShell("firewall-cmd --reload")
	} else {
		tools.ExecShell("ufw delete allow " + cast.ToString(port) + "/" + protocol)
		tools.ExecShell("ufw reload")
	}

	Success(ctx, nil)
}

func (r *SafeController) firewallStatus() bool {
	var out string
	var running bool
	if tools.IsRHEL() {
		out = tools.ExecShell("systemctl status firewalld | grep Active | awk '{print $3}'")
		if out == "(running)" {
			running = true
		} else {
			running = false
		}
	} else {
		out = tools.ExecShell("ufw status | grep Status | awk '{print $2}'")
		if out == "active" {
			running = true
		} else {
			running = false
		}
	}

	return running
}

func (r *SafeController) GetSshStatus(ctx http.Context) {
	var out string
	if tools.IsRHEL() {
		out = tools.ExecShell("systemctl status sshd | grep Active | awk '{print $3}'")
	} else {
		out = tools.ExecShell("systemctl status ssh | grep Active | awk '{print $3}'")
	}

	running := false
	if out == "(running)" {
		running = true
	}

	Success(ctx, running)
}

func (r *SafeController) SetSshStatus(ctx http.Context) {
	if ctx.Request().InputBool("status") {
		if tools.IsRHEL() {
			tools.ExecShell("systemctl enable sshd")
			tools.ExecShell("systemctl start sshd")
		} else {
			tools.ExecShell("systemctl enable ssh")
			tools.ExecShell("systemctl start ssh")
		}
	} else {
		if tools.IsRHEL() {
			tools.ExecShell("systemctl stop sshd")
			tools.ExecShell("systemctl disable sshd")
		} else {
			tools.ExecShell("systemctl stop ssh")
			tools.ExecShell("systemctl disable ssh")
		}
	}

	Success(ctx, nil)
}

func (r *SafeController) GetSshPort(ctx http.Context) {
	out := tools.ExecShell("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'")
	Success(ctx, out)
}

func (r *SafeController) SetSshPort(ctx http.Context) {
	port := ctx.Request().InputInt("port", 0)
	if port == 0 {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	oldPort := tools.ExecShell("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'")
	tools.ExecShell("sed -i 's/#Port " + oldPort + "/Port " + cast.ToString(port) + "/g' /etc/ssh/sshd_config")
	tools.ExecShell("sed -i 's/Port " + oldPort + "/Port " + cast.ToString(port) + "/g' /etc/ssh/sshd_config")

	if status := tools.ExecShell("systemctl status sshd | grep Active | awk '{print $3}'"); status == "(running)" {
		tools.ExecShell("systemctl restart sshd")
	}

	Success(ctx, nil)
}

func (r *SafeController) GetPingStatus(ctx http.Context) {
	if tools.IsRHEL() {
		out := tools.ExecShell(`firewall-cmd --list-all 2>&1`)
		if !strings.Contains(out, `rule protocol value="icmp" drop`) {
			Success(ctx, true)
		} else {
			Success(ctx, false)
		}
	} else {
		config := tools.ReadFile("/etc/ufw/before.rules")
		if strings.Contains(config, "-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT") {
			Success(ctx, true)
		} else {
			Success(ctx, false)
		}
	}
}

func (r *SafeController) SetPingStatus(ctx http.Context) {
	if tools.IsRHEL() {
		if ctx.Request().InputBool("status") {
			tools.ExecShell(`firewall-cmd --permanent --remove-rich-rule='rule protocol value=icmp drop'`)
		} else {
			tools.ExecShell(`firewall-cmd --permanent --add-rich-rule='rule protocol value=icmp drop'`)
		}
		tools.ExecShell(`firewall-cmd --reload`)
	} else {
		if ctx.Request().InputBool("status") {
			tools.ExecShell(`sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/g' /etc/ufw/before.rules`)
		} else {
			tools.ExecShell(`sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/g' /etc/ufw/before.rules`)
		}
		tools.ExecShell(`ufw reload`)
	}

	Success(ctx, nil)
}
