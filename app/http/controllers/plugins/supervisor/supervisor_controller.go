package supervisor

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
func (c *SupervisorController) Status(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	status := tools.Exec(`systemctl status ` + c.ServiceName + ` | grep Active | grep -v grep | awk '{print $2}'`)
	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Start 启动
func (c *SupervisorController) Start(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	tools.Exec(`systemctl start ` + c.ServiceName)
	status := tools.Exec(`systemctl status ` + c.ServiceName + ` | grep Active | grep -v grep | awk '{print $2}'`)
	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Stop 停止
func (c *SupervisorController) Stop(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	tools.Exec(`systemctl stop ` + c.ServiceName)
	status := tools.Exec(`systemctl status ` + c.ServiceName + ` | grep Active | grep -v grep | awk '{print $2}'`)
	if status != "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Restart 重启
func (c *SupervisorController) Restart(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	tools.Exec(`systemctl restart ` + c.ServiceName)
	status := tools.Exec(`systemctl status ` + c.ServiceName + ` | grep Active | grep -v grep | awk '{print $2}'`)
	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Reload 重载
func (c *SupervisorController) Reload(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	tools.Exec(`systemctl reload ` + c.ServiceName)
	status := tools.Exec(`systemctl status ` + c.ServiceName + ` | grep Active | grep -v grep | awk '{print $2}'`)
	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Log 日志
func (c *SupervisorController) Log(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	log := tools.Exec(`tail -n 200 /var/log/supervisor/supervisord.log`)
	controllers.Success(ctx, log)
}

// ClearLog 清空日志
func (c *SupervisorController) ClearLog(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	tools.Exec(`echo "" > /var/log/supervisor/supervisord.log`)
	controllers.Success(ctx, nil)
}

// Config 获取配置
func (c *SupervisorController) Config(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	var config string
	if tools.IsRHEL() {
		config = tools.Read(`/etc/supervisord.conf`)
	} else {
		config = tools.Read(`/etc/supervisor/supervisord.conf`)
	}
	controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (c *SupervisorController) SaveConfig(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	config := ctx.Request().Input("config")
	if tools.IsRHEL() {
		tools.Write(`/etc/supervisord.conf`, config, 0644)
	} else {
		tools.Write(`/etc/supervisor/supervisord.conf`, config, 0644)
	}

	c.Restart(ctx)
}

// Processes 进程列表
func (c *SupervisorController) Processes(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
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
		controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []process{},
		})
		return
	}
	if endIndex > len(processList) {
		endIndex = len(processList)
	}
	pagedProcessList := processList[startIndex:endIndex]

	controllers.Success(ctx, http.Json{
		"total": len(processList),
		"items": pagedProcessList,
	})
}

// StartProcess 启动进程
func (c *SupervisorController) StartProcess(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	process := ctx.Request().Input("process")
	tools.Exec(`supervisorctl start ` + process)
	controllers.Success(ctx, nil)
}

// StopProcess 停止进程
func (c *SupervisorController) StopProcess(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	process := ctx.Request().Input("process")
	tools.Exec(`supervisorctl stop ` + process)
	controllers.Success(ctx, nil)
}

// RestartProcess 重启进程
func (c *SupervisorController) RestartProcess(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	process := ctx.Request().Input("process")
	tools.Exec(`supervisorctl restart ` + process)
	controllers.Success(ctx, nil)
}

// ProcessLog 进程日志
func (c *SupervisorController) ProcessLog(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	process := ctx.Request().Input("process")
	var logPath string
	if tools.IsRHEL() {
		logPath = tools.Exec(`cat '/etc/supervisord.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	} else {
		logPath = tools.Exec(`cat '/etc/supervisor/conf.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	}

	log := tools.Exec(`tail -n 200 ` + logPath)
	controllers.Success(ctx, log)
}

// ClearProcessLog 清空进程日志
func (c *SupervisorController) ClearProcessLog(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	process := ctx.Request().Input("process")
	var logPath string
	if tools.IsRHEL() {
		logPath = tools.Exec(`cat '/etc/supervisord.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	} else {
		logPath = tools.Exec(`cat '/etc/supervisor/conf.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	}

	tools.Exec(`echo "" > ` + logPath)
	controllers.Success(ctx, nil)
}

// ProcessConfig 获取进程配置
func (c *SupervisorController) ProcessConfig(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	process := ctx.Request().Query("process")
	var config string
	if tools.IsRHEL() {
		config = tools.Read(`/etc/supervisord.d/` + process + `.conf`)
	} else {
		config = tools.Read(`/etc/supervisor/conf.d/` + process + `.conf`)
	}

	controllers.Success(ctx, config)
}

// SaveProcessConfig 保存进程配置
func (c *SupervisorController) SaveProcessConfig(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
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

	controllers.Success(ctx, nil)
}

// AddProcess 添加进程
func (c *SupervisorController) AddProcess(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"name":    "required|alpha_dash",
		"user":    "required|alpha_dash",
		"path":    "required",
		"command": "required",
		"num":     "required",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
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

	controllers.Success(ctx, nil)
}

// DeleteProcess 删除进程
func (c *SupervisorController) DeleteProcess(ctx http.Context) {
	if !controllers.Check(ctx, "supervisor") {
		return
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

	controllers.Success(ctx, nil)
}
