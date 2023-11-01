package plugins

import (
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"panel/pkg/tools"

	"panel/app/http/controllers"
)

type SupervisorController struct {
	ServiceName string
}

func NewSupervisorController() *SupervisorController {
	var serviceName string
	if tools.IsRHEL() {
		serviceName = "supervisord"
	} else {
		serviceName = "supervisor"
	}

	return &SupervisorController{
		ServiceName: serviceName,
	}
}

// Status 状态
func (r *SupervisorController) Status(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	status := tools.Exec(`systemctl status ` + r.ServiceName + ` | grep Active | grep -v grep | awk '{print $2}'`)
	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Start 启动
func (r *SupervisorController) Start(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	tools.Exec(`systemctl start ` + r.ServiceName)
	status := tools.Exec(`systemctl status ` + r.ServiceName + ` | grep Active | grep -v grep | awk '{print $2}'`)
	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Stop 停止
func (r *SupervisorController) Stop(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	tools.Exec(`systemctl stop ` + r.ServiceName)
	status := tools.Exec(`systemctl status ` + r.ServiceName + ` | grep Active | grep -v grep | awk '{print $2}'`)
	if status != "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Restart 重启
func (r *SupervisorController) Restart(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	tools.Exec(`systemctl restart ` + r.ServiceName)
	status := tools.Exec(`systemctl status ` + r.ServiceName + ` | grep Active | grep -v grep | awk '{print $2}'`)
	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Reload 重载
func (r *SupervisorController) Reload(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	tools.Exec(`systemctl reload ` + r.ServiceName)
	status := tools.Exec(`systemctl status ` + r.ServiceName + ` | grep Active | grep -v grep | awk '{print $2}'`)
	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Log 日志
func (r *SupervisorController) Log(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	log := tools.Exec(`tail -n 200 /var/log/supervisor/supervisord.log`)
	return controllers.Success(ctx, log)
}

// ClearLog 清空日志
func (r *SupervisorController) ClearLog(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	tools.Exec(`echo "" > /var/log/supervisor/supervisord.log`)
	return controllers.Success(ctx, nil)
}

// Config 获取配置
func (r *SupervisorController) Config(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	var config string
	if tools.IsRHEL() {
		config = tools.Read(`/etc/supervisord.conf`)
	} else {
		config = tools.Read(`/etc/supervisor/supervisord.conf`)
	}
	return controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (r *SupervisorController) SaveConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	config := ctx.Request().Input("config")
	if tools.IsRHEL() {
		tools.Write(`/etc/supervisord.conf`, config, 0644)
	} else {
		tools.Write(`/etc/supervisor/supervisord.conf`, config, 0644)
	}

	return r.Restart(ctx)
}

// Processes 进程列表
func (r *SupervisorController) Processes(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)

	type process struct {
		Name   string `json:"name"`
		Status string `json:"status"`
		Pid    string `json:"pid"`
		Uptime string `json:"uptime"`
	}

	out := tools.Exec(`supervisorctl status | awk '{print $1}'`)
	var processList []process
	for _, line := range strings.Split(out, "\n") {
		if len(line) == 0 {
			continue
		}

		var p process
		p.Name = line
		p.Status = tools.Exec(`supervisorctl status ` + line + ` | awk '{print $2}'`)
		if p.Status == "RUNNING" {
			p.Pid = strings.ReplaceAll(tools.Exec(`supervisorctl status `+line+` | awk '{print $4}'`), ",", "")
			p.Uptime = tools.Exec(`supervisorctl status ` + line + ` | awk '{print $6}'`)
		} else {
			p.Pid = "-"
			p.Uptime = "-"
		}
		processList = append(processList, p)
	}

	startIndex := (page - 1) * limit
	endIndex := page * limit
	if startIndex > len(processList) {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []process{},
		})
	}
	if endIndex > len(processList) {
		endIndex = len(processList)
	}
	pagedProcessList := processList[startIndex:endIndex]
	if pagedProcessList == nil {
		pagedProcessList = []process{}
	}

	return controllers.Success(ctx, http.Json{
		"total": len(processList),
		"items": pagedProcessList,
	})
}

// StartProcess 启动进程
func (r *SupervisorController) StartProcess(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	process := ctx.Request().Input("process")
	tools.Exec(`supervisorctl start ` + process)
	return controllers.Success(ctx, nil)
}

// StopProcess 停止进程
func (r *SupervisorController) StopProcess(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	process := ctx.Request().Input("process")
	tools.Exec(`supervisorctl stop ` + process)
	return controllers.Success(ctx, nil)
}

// RestartProcess 重启进程
func (r *SupervisorController) RestartProcess(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	process := ctx.Request().Input("process")
	tools.Exec(`supervisorctl restart ` + process)
	return controllers.Success(ctx, nil)
}

// ProcessLog 进程日志
func (r *SupervisorController) ProcessLog(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	process := ctx.Request().Input("process")
	var logPath string
	if tools.IsRHEL() {
		logPath = tools.Exec(`cat '/etc/supervisord.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	} else {
		logPath = tools.Exec(`cat '/etc/supervisor/conf.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	}

	log := tools.Exec(`tail -n 200 ` + logPath)
	return controllers.Success(ctx, log)
}

// ClearProcessLog 清空进程日志
func (r *SupervisorController) ClearProcessLog(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	process := ctx.Request().Input("process")
	var logPath string
	if tools.IsRHEL() {
		logPath = tools.Exec(`cat '/etc/supervisord.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	} else {
		logPath = tools.Exec(`cat '/etc/supervisor/conf.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	}

	tools.Exec(`echo "" > ` + logPath)
	return controllers.Success(ctx, nil)
}

// ProcessConfig 获取进程配置
func (r *SupervisorController) ProcessConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	process := ctx.Request().Query("process")
	var config string
	if tools.IsRHEL() {
		config = tools.Read(`/etc/supervisord.d/` + process + `.conf`)
	} else {
		config = tools.Read(`/etc/supervisor/conf.d/` + process + `.conf`)
	}

	return controllers.Success(ctx, config)
}

// SaveProcessConfig 保存进程配置
func (r *SupervisorController) SaveProcessConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	process := ctx.Request().Input("process")
	config := ctx.Request().Input("config")
	if tools.IsRHEL() {
		tools.Write(`/etc/supervisord.d/`+process+`.conf`, config, 0644)
	} else {
		tools.Write(`/etc/supervisor/conf.d/`+process+`.conf`, config, 0644)
	}
	tools.Exec(`supervisorctl reread`)
	tools.Exec(`supervisorctl update`)
	tools.Exec(`supervisorctl start ` + process)

	return controllers.Success(ctx, nil)
}

// AddProcess 添加进程
func (r *SupervisorController) AddProcess(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"name":    "required|alpha_dash",
		"user":    "required|alpha_dash",
		"path":    "required",
		"command": "required",
		"num":     "required",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
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
	if tools.IsRHEL() {
		tools.Write(`/etc/supervisord.d/`+name+`.conf`, config, 0644)
	} else {
		tools.Write(`/etc/supervisor/conf.d/`+name+`.conf`, config, 0644)
	}
	tools.Exec(`supervisorctl reread`)
	tools.Exec(`supervisorctl update`)
	tools.Exec(`supervisorctl start ` + name)

	return controllers.Success(ctx, nil)
}

// DeleteProcess 删除进程
func (r *SupervisorController) DeleteProcess(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "supervisor")
	if check != nil {
		return check
	}

	process := ctx.Request().Input("process")
	tools.Exec(`supervisorctl stop ` + process)
	var logPath string
	if tools.IsRHEL() {
		logPath = tools.Exec(`cat '/etc/supervisord.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
		tools.Remove(`/etc/supervisord.d/` + process + `.conf`)
	} else {
		logPath = tools.Exec(`cat '/etc/supervisor/conf.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
		tools.Remove(`/etc/supervisor/conf.d/` + process + `.conf`)
	}
	tools.Remove(logPath)
	tools.Exec(`supervisorctl reread`)
	tools.Exec(`supervisorctl update`)

	return controllers.Success(ctx, nil)
}
