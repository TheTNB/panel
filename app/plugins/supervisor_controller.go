package plugins

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/os"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
)

type SupervisorController struct {
	service string
}

func NewSupervisorController() *SupervisorController {
	var service string
	if os.IsRHEL() {
		service = "supervisord"
	} else {
		service = "supervisor"
	}

	return &SupervisorController{
		service: service,
	}
}

// Service 获取服务名称
func (r *SupervisorController) Service(ctx http.Context) http.Response {
	return h.Success(ctx, r.service)
}

// Log 日志
func (r *SupervisorController) Log(ctx http.Context) http.Response {
	log, err := shell.Execf(`tail -n 200 /var/log/supervisor/supervisord.log`)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, log)
	}

	return h.Success(ctx, log)
}

// ClearLog 清空日志
func (r *SupervisorController) ClearLog(ctx http.Context) http.Response {
	if out, err := shell.Execf(`echo "" > /var/log/supervisor/supervisord.log`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}

// Config 获取配置
func (r *SupervisorController) Config(ctx http.Context) http.Response {
	var config string
	var err error
	if os.IsRHEL() {
		config, err = io.Read(`/etc/supervisord.conf`)
	} else {
		config, err = io.Read(`/etc/supervisor/supervisord.conf`)
	}

	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, config)
}

// SaveConfig 保存配置
func (r *SupervisorController) SaveConfig(ctx http.Context) http.Response {
	config := ctx.Request().Input("config")
	var err error
	if os.IsRHEL() {
		err = io.Write(`/etc/supervisord.conf`, config, 0644)
	} else {
		err = io.Write(`/etc/supervisor/supervisord.conf`, config, 0644)
	}

	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err = systemctl.Restart(r.service); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("重启 %s 服务失败", r.service))
	}

	return h.Success(ctx, nil)
}

// Processes 进程列表
func (r *SupervisorController) Processes(ctx http.Context) http.Response {
	type process struct {
		Name   string `json:"name"`
		Status string `json:"status"`
		Pid    string `json:"pid"`
		Uptime string `json:"uptime"`
	}

	out, err := shell.Execf(`supervisorctl status | awk '{print $1}'`)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	var processes []process
	for _, line := range strings.Split(out, "\n") {
		if len(line) == 0 {
			continue
		}

		var p process
		p.Name = line
		if status, err := shell.Execf(`supervisorctl status ` + line + ` | awk '{print $2}'`); err == nil {
			p.Status = status
		}
		if p.Status == "RUNNING" {
			pid, _ := shell.Execf(`supervisorctl status ` + line + ` | awk '{print $4}'`)
			p.Pid = strings.ReplaceAll(pid, ",", "")
			uptime, _ := shell.Execf(`supervisorctl status ` + line + ` | awk '{print $6}'`)
			p.Uptime = uptime
		} else {
			p.Pid = "-"
			p.Uptime = "-"
		}
		processes = append(processes, p)
	}

	paged, total := h.Paginate(ctx, processes)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// StartProcess 启动进程
func (r *SupervisorController) StartProcess(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	if out, err := shell.Execf(`supervisorctl start %s`, process); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}

// StopProcess 停止进程
func (r *SupervisorController) StopProcess(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	if out, err := shell.Execf(`supervisorctl stop %s`, process); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}

// RestartProcess 重启进程
func (r *SupervisorController) RestartProcess(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	if out, err := shell.Execf(`supervisorctl restart %s`, process); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}

// ProcessLog 进程日志
func (r *SupervisorController) ProcessLog(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	var logPath string
	var err error
	if os.IsRHEL() {
		logPath, err = shell.Execf(`cat '/etc/supervisord.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, process)
	} else {
		logPath, err = shell.Execf(`cat '/etc/supervisor/conf.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, process)
	}

	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "无法从进程"+process+"的配置文件中获取日志路径")
	}

	log, err := shell.Execf(`tail -n 200 ` + logPath)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, log)
	}

	return h.Success(ctx, log)
}

// ClearProcessLog 清空进程日志
func (r *SupervisorController) ClearProcessLog(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	var logPath string
	var err error
	if os.IsRHEL() {
		logPath, err = shell.Execf(`cat '/etc/supervisord.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, process)
	} else {
		logPath, err = shell.Execf(`cat '/etc/supervisor/conf.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, process)
	}

	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("无法从进程%s的配置文件中获取日志路径", process))
	}

	if out, err := shell.Execf(`echo "" > ` + logPath); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}

// ProcessConfig 获取进程配置
func (r *SupervisorController) ProcessConfig(ctx http.Context) http.Response {
	process := ctx.Request().Query("process")
	var config string
	var err error
	if os.IsRHEL() {
		config, err = io.Read(`/etc/supervisord.d/` + process + `.conf`)
	} else {
		config, err = io.Read(`/etc/supervisor/conf.d/` + process + `.conf`)
	}

	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, config)
}

// SaveProcessConfig 保存进程配置
func (r *SupervisorController) SaveProcessConfig(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	config := ctx.Request().Input("config")
	var err error
	if os.IsRHEL() {
		err = io.Write(`/etc/supervisord.d/`+process+`.conf`, config, 0644)
	} else {
		err = io.Write(`/etc/supervisor/conf.d/`+process+`.conf`, config, 0644)
	}

	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}

	_, _ = shell.Execf(`supervisorctl reread`)
	_, _ = shell.Execf(`supervisorctl update`)
	_, _ = shell.Execf(`supervisorctl restart %s`, process)

	return h.Success(ctx, nil)
}

// AddProcess 添加进程
func (r *SupervisorController) AddProcess(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"name":    "required|alpha_dash",
		"user":    "required|alpha_dash",
		"path":    "required",
		"command": "required",
		"num":     "required",
	}); sanitize != nil {
		return sanitize
	}

	name := ctx.Request().Input("name")
	user := ctx.Request().Input("user")
	path := ctx.Request().Input("path")
	command := ctx.Request().Input("command")
	num := ctx.Request().InputInt("num", 1)
	config := `[program:` + name + `]
command=` + command + `
process_name=%(program_name)s
directory=` + path + `
autostart=true
autorestart=true
user=` + user + `
numprocs=` + strconv.Itoa(num) + `
redirect_stderr=true
stdout_logfile=/var/log/supervisor/` + name + `.log
stdout_logfile_maxbytes=2MB
`

	var err error
	if os.IsRHEL() {
		err = io.Write(`/etc/supervisord.d/`+name+`.conf`, config, 0644)
	} else {
		err = io.Write(`/etc/supervisor/conf.d/`+name+`.conf`, config, 0644)
	}

	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}

	_, _ = shell.Execf(`supervisorctl reread`)
	_, _ = shell.Execf(`supervisorctl update`)
	_, _ = shell.Execf(`supervisorctl start %s`, name)

	return h.Success(ctx, nil)
}

// DeleteProcess 删除进程
func (r *SupervisorController) DeleteProcess(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	if out, err := shell.Execf(`supervisorctl stop %s`, process); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	var logPath string
	var err error
	if os.IsRHEL() {
		logPath, err = shell.Execf(`cat '/etc/supervisord.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, process)
		if err := io.Remove(`/etc/supervisord.d/` + process + `.conf`); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}
	} else {
		logPath, err = shell.Execf(`cat '/etc/supervisor/conf.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, process)
		if err := io.Remove(`/etc/supervisor/conf.d/` + process + `.conf`); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}
	}

	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "无法从进程"+process+"的配置文件中获取日志路径")
	}

	if err := io.Remove(logPath); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	_, _ = shell.Execf(`supervisorctl reread`)
	_, _ = shell.Execf(`supervisorctl update`)

	return h.Success(ctx, nil)
}
