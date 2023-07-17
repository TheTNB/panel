package controllers

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/spf13/cast"

	"panel/packages/helper"
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
	if ctx.Request().QueryBool("status") {
		if helper.IsRHEL() {
			out = helper.ExecShell("systemctl start firewalld")
		} else {
			out = helper.ExecShell("echo y | ufw enable")
		}
	} else {
		if helper.IsRHEL() {
			out = helper.ExecShell("systemctl stop firewalld")
		} else {
			out = helper.ExecShell("ufw disable")
		}
	}

	Success(ctx, out)
}

func (r *SafeController) GetFirewallRules(ctx http.Context) {
	if !r.firewallStatus() {
		Error(ctx, http.StatusBadRequest, "防火墙未启动")
		return
	}

	if helper.IsRHEL() {
		out := helper.ExecShell("firewall-cmd --list-all 2>&1")
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

		Success(ctx, rules)
	} else {
		out := helper.ExecShell("ufw status numbered | grep ALLOW | awk '{print $2}'")
		if len(out) == 0 {
			Success(ctx, nil)
			return
		}
		var rules []map[string]string
		for _, port := range strings.Split(out, "\n") {
			if strings.Contains(port, "]") {
				continue
			}
			rule := strings.Split(port, "/")
			rules = append(rules, map[string]string{
				"port":     rule[0],
				"protocol": rule[1],
			})
		}

		Success(ctx, rules)
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

	if helper.IsRHEL() {
		helper.ExecShell("firewall-cmd --remove-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1")
		helper.ExecShell("firewall-cmd --add-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1")
		helper.ExecShell("firewall-cmd --reload")
	} else {
		helper.ExecShell("ufw delete allow " + cast.ToString(port) + "/" + protocol)
		helper.ExecShell("ufw allow " + cast.ToString(port) + "/" + protocol)
		helper.ExecShell("ufw reload")
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

	if helper.IsRHEL() {
		helper.ExecShell("firewall-cmd --remove-port=" + cast.ToString(port) + "/" + protocol + " --permanent 2>&1")
		helper.ExecShell("firewall-cmd --reload")
	} else {
		helper.ExecShell("ufw delete allow " + cast.ToString(port) + "/" + protocol)
		helper.ExecShell("ufw reload")
	}

	Success(ctx, nil)
}

func (r *SafeController) firewallStatus() bool {
	var out string
	var running bool
	if helper.IsRHEL() {
		out = helper.ExecShell("systemctl status firewalld | grep Active | awk '{print $3}'")
		if out == "(running)" {
			running = true
		} else {
			running = false
		}
	} else {
		out = helper.ExecShell("ufw status | grep Status | awk '{print $2}'")
		if out == "active" {
			running = true
		} else {
			running = false
		}
	}

	return running
}

func (r *SafeController) GetSshStatus(ctx http.Context) {
	out := helper.ExecShell("systemctl status sshd | grep Active | awk '{print $3}'")
	running := false
	if out == "(running)" {
		running = true
	}

	Success(ctx, running)
}

func (r *SafeController) SetSshStatus(ctx http.Context) {
	if ctx.Request().QueryBool("status") {
		helper.ExecShell("systemctl enable sshd")
		helper.ExecShell("systemctl start sshd")
	} else {
		helper.ExecShell("systemctl stop sshd")
		helper.ExecShell("systemctl disable sshd")
	}

	Success(ctx, nil)
}

func (r *SafeController) GetSshPort(ctx http.Context) {
	out := helper.ExecShell("cat /etc/ssh/sshd_config | grep Port | awk '{print $2}'")
	Success(ctx, out)
}

func (r *SafeController) SetSshPort(ctx http.Context) {
	port := ctx.Request().InputInt("port", 0)
	if port == 0 {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	oldPort := helper.ExecShell("cat /etc/ssh/sshd_config | grep Port | awk '{print $2}'")
	helper.ExecShell("sed -i 's/#Port " + oldPort + "/Port " + cast.ToString(port) + "/g' /etc/ssh/sshd_config")
	helper.ExecShell("sed -i 's/Port " + oldPort + "/Port " + cast.ToString(port) + "/g' /etc/ssh/sshd_config")

	if status := helper.ExecShell("systemctl status sshd | grep Active | awk '{print $3}'"); status == "(running)" {
		helper.ExecShell("systemctl restart sshd")
	}

	Success(ctx, nil)
}

func (r *SafeController) GetPingStatus(ctx http.Context) {
	if helper.IsRHEL() {
		out := helper.ExecShell("firewall-cmd --query-rich-rule='rule protocol value=icmp drop' 2>&1")
		if out == "no" {
			Success(ctx, true)
		} else {
			Success(ctx, false)
		}
	} else {
		config := helper.ReadFile("/etc/ufw/before.rules")
		if strings.Contains(config, "-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT") {
			Success(ctx, true)
		} else {
			Success(ctx, false)
		}
	}
}

func (r *SafeController) SetPingStatus(ctx http.Context) {
	if helper.IsRHEL() {
		if ctx.Request().QueryBool("status") {
			helper.ExecShell("firewall-cmd --permanent --add-rich-rule='rule protocol value=icmp drop'")
		} else {
			helper.ExecShell("firewall-cmd --permanent --remove-rich-rule='rule protocol value=icmp drop'")
		}
		helper.ExecShell("firewall-cmd --reload")
	} else {
		if ctx.Request().QueryBool("status") {
			helper.ExecShell("sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/g' /etc/ufw/before.rules")
		} else {
			helper.ExecShell("sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/g' /etc/ufw/before.rules")
		}
		helper.ExecShell("ufw reload")
	}

	Success(ctx, nil)
}
