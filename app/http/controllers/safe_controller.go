package controllers

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/os"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
)

type SafeController struct {
	ssh string
}

func NewSafeController() *SafeController {
	var ssh string
	if os.IsRHEL() {
		ssh = "sshd"
	} else {
		ssh = "ssh"
	}

	return &SafeController{
		ssh: ssh,
	}
}

// GetFirewallStatus 获取防火墙状态
func (r *SafeController) GetFirewallStatus(ctx http.Context) http.Response {
	return h.Success(ctx, r.firewallStatus())
}

// SetFirewallStatus 设置防火墙状态
func (r *SafeController) SetFirewallStatus(ctx http.Context) http.Response {
	var err error
	if ctx.Request().InputBool("status") {
		if os.IsRHEL() {
			err = systemctl.Start("firewalld")
			if err == nil {
				err = systemctl.Enable("firewalld")
			}
		} else {
			_, err = shell.Execf("echo y | ufw enable")
			if err == nil {
				err = systemctl.Start("ufw")
			}
			if err == nil {
				err = systemctl.Enable("ufw")
			}
		}
	} else {
		if os.IsRHEL() {
			err = systemctl.Stop("firewalld")
			if err == nil {
				err = systemctl.Disable("firewalld")
			}
		} else {
			_, err = shell.Execf("ufw disable")
			if err == nil {
				err = systemctl.Stop("ufw")
			}
			if err == nil {
				err = systemctl.Disable("ufw")
			}
		}
	}

	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// GetFirewallRules 获取防火墙规则
func (r *SafeController) GetFirewallRules(ctx http.Context) http.Response {
	if !r.firewallStatus() {
		return h.Success(ctx, nil)
	}

	var rules []map[string]string
	if os.IsRHEL() {
		out, err := shell.Execf("firewall-cmd --list-all")
		if err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}

		match := regexp.MustCompile(`ports: (.*)`).FindStringSubmatch(out)
		if len(match) == 0 {
			return h.Success(ctx, http.Json{
				"total": 0,
				"items": []map[string]string{},
			})
		}
		ports := strings.Split(match[1], " ")
		for _, port := range ports {
			rule := strings.Split(port, "/")
			if len(rule) < 2 {
				rules = append(rules, map[string]string{
					"port":     rule[0],
					"protocol": "all",
				})
			} else {
				rules = append(rules, map[string]string{
					"port":     rule[0],
					"protocol": rule[1],
				})
			}
		}
	} else {
		out, err := shell.Execf("ufw status | grep -v '(v6)' | grep ALLOW | awk '{print $1}'")
		if err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}

		if len(out) == 0 {
			return h.Success(ctx, http.Json{
				"total": 0,
				"items": []map[string]string{},
			})
		}
		for _, port := range strings.Split(out, "\n") {
			rule := strings.Split(port, "/")
			if len(rule) < 2 {
				rules = append(rules, map[string]string{
					"port":     rule[0],
					"protocol": "all",
				})
			} else {
				rules = append(rules, map[string]string{
					"port":     rule[0],
					"protocol": rule[1],
				})
			}
		}
	}

	paged, total := h.Paginate(ctx, rules)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// AddFirewallRule 添加防火墙规则
func (r *SafeController) AddFirewallRule(ctx http.Context) http.Response {
	if !r.firewallStatus() {
		return h.Error(ctx, http.StatusInternalServerError, "防火墙未启动")
	}

	port := ctx.Request().Input("port")
	protocol := ctx.Request().Input("protocol")
	if port == "" || protocol == "" || (protocol != "tcp" && protocol != "udp") {
		return h.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}
	// 端口有 2 种写法，一种是 80-443，一种是 80
	if strings.Contains(port, "-") {
		ports := strings.Split(port, "-")
		startPort := cast.ToInt(ports[0])
		endPort := cast.ToInt(ports[1])
		if startPort < 1 || startPort > 65535 || endPort < 1 || endPort > 65535 || startPort > endPort {
			return h.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
		}
	} else {
		port := cast.ToInt(port)
		if port < 1 || port > 65535 {
			return h.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
		}
	}

	if os.IsRHEL() {
		if out, err := shell.Execf("firewall-cmd --remove-port=%s/%s --permanent", port, protocol); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := shell.Execf("firewall-cmd --add-port=%s/%s --permanent", port, protocol); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := shell.Execf("firewall-cmd --reload"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	} else {
		// ufw 需要替换 - 为 : 添加
		if strings.Contains(port, "-") {
			port = strings.ReplaceAll(port, "-", ":")
		}
		if out, err := shell.Execf("ufw delete allow %s/%s", port, protocol); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := shell.Execf("ufw allow %s/%s", port, protocol); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := shell.Execf("ufw reload"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	return h.Success(ctx, nil)
}

// DeleteFirewallRule 删除防火墙规则
func (r *SafeController) DeleteFirewallRule(ctx http.Context) http.Response {
	if !r.firewallStatus() {
		return h.Error(ctx, http.StatusUnprocessableEntity, "防火墙未启动")
	}

	port := ctx.Request().Input("port")
	protocol := ctx.Request().Input("protocol")
	if port == "" || protocol == "" {
		return h.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}
	if protocol == "all" {
		protocol = ""
	} else {
		protocol = "/" + protocol
	}

	if os.IsRHEL() {
		if out, err := shell.Execf("firewall-cmd --remove-port=%s%s --permanent", port, protocol); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := shell.Execf("firewall-cmd --reload"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	} else {
		if out, err := shell.Execf("ufw delete allow %s%s", port, protocol); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := shell.Execf("ufw reload"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	return h.Success(ctx, nil)
}

// firewallStatus 获取防火墙状态
func (r *SafeController) firewallStatus() bool {
	var running bool
	if os.IsRHEL() {
		running, _ = systemctl.Status("firewalld")
	} else {
		running, _ = systemctl.Status("ufw")
	}

	return running
}

// GetSshStatus 获取 SSH 状态
func (r *SafeController) GetSshStatus(ctx http.Context) http.Response {
	running, err := systemctl.Status(r.ssh)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, running)
}

// SetSshStatus 设置 SSH 状态
func (r *SafeController) SetSshStatus(ctx http.Context) http.Response {
	if ctx.Request().InputBool("status") {
		if err := systemctl.Enable(r.ssh); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		if err := systemctl.Start(r.ssh); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}
	} else {
		if err := systemctl.Stop(r.ssh); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		if err := systemctl.Disable(r.ssh); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}
	}

	return h.Success(ctx, nil)
}

// GetSshPort 获取 SSH 端口
func (r *SafeController) GetSshPort(ctx http.Context) http.Response {
	out, err := shell.Execf("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, cast.ToInt(out))
}

// SetSshPort 设置 SSH 端口
func (r *SafeController) SetSshPort(ctx http.Context) http.Response {
	port := ctx.Request().InputInt("port", 0)
	if port == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}

	oldPort, err := shell.Execf("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, oldPort)
	}
	_, _ = shell.Execf("sed -i 's/#Port %s/Port %d/g' /etc/ssh/sshd_config", oldPort, port)
	_, _ = shell.Execf("sed -i 's/Port %s/Port %d/g' /etc/ssh/sshd_config", oldPort, port)

	status, _ := systemctl.Status(r.ssh)
	if status {
		_ = systemctl.Restart(r.ssh)
	}

	return h.Success(ctx, nil)
}

// GetPingStatus 获取 Ping 状态
func (r *SafeController) GetPingStatus(ctx http.Context) http.Response {
	if os.IsRHEL() {
		out, err := shell.Execf(`firewall-cmd --list-all`)
		if err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}

		if !strings.Contains(out, `rule protocol value="icmp" drop`) {
			return h.Success(ctx, true)
		} else {
			return h.Success(ctx, false)
		}
	} else {
		config, err := io.Read("/etc/ufw/before.rules")
		if err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		if strings.Contains(config, "-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT") {
			return h.Success(ctx, true)
		} else {
			return h.Success(ctx, false)
		}
	}
}

// SetPingStatus 设置 Ping 状态
func (r *SafeController) SetPingStatus(ctx http.Context) http.Response {
	var out string
	var err error
	if os.IsRHEL() {
		if ctx.Request().InputBool("status") {
			out, err = shell.Execf(`firewall-cmd --permanent --remove-rich-rule='rule protocol value=icmp drop'`)
		} else {
			out, err = shell.Execf(`firewall-cmd --permanent --add-rich-rule='rule protocol value=icmp drop'`)
		}
	} else {
		if ctx.Request().InputBool("status") {
			out, err = shell.Execf(`sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/g' /etc/ufw/before.rules`)
		} else {
			out, err = shell.Execf(`sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/g' /etc/ufw/before.rules`)
		}
	}

	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	if os.IsRHEL() {
		out, err = shell.Execf(`firewall-cmd --reload`)
	} else {
		out, err = shell.Execf(`ufw reload`)
	}

	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}
