package supervisor

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-rat/chix"
	"github.com/spf13/cast"

	"github.com/tnb-labs/panel/internal/service"
	"github.com/tnb-labs/panel/pkg/io"
	"github.com/tnb-labs/panel/pkg/os"
	"github.com/tnb-labs/panel/pkg/shell"
	"github.com/tnb-labs/panel/pkg/systemctl"
)

type App struct {
	name string
}

func NewApp() *App {
	var name string
	if os.IsRHEL() {
		name = "supervisord"
	} else {
		name = "supervisor"
	}

	return &App{
		name: name,
	}
}

func (s *App) Route(r chi.Router) {
	r.Get("/service", s.Service)
	r.Post("/clearLog", s.ClearLog)
	r.Get("/config", s.GetConfig)
	r.Post("/config", s.UpdateConfig)
	r.Get("/processes", s.Processes)
	r.Post("/processes/{process}/start", s.StartProcess)
	r.Post("/processes/{process}/stop", s.StopProcess)
	r.Post("/processes/{process}/restart", s.RestartProcess)
	r.Get("/processes/{process}/log", s.ProcessLog)
	r.Post("/processes/{process}/clearLog", s.ClearProcessLog)
	r.Get("/processes/{process}", s.ProcessConfig)
	r.Post("/processes/{process}", s.UpdateProcessConfig)
	r.Delete("/processes/{process}", s.DeleteProcess)
	r.Post("/processes", s.CreateProcess)
}

// Service 获取服务名称
func (s *App) Service(w http.ResponseWriter, r *http.Request) {
	service.Success(w, s.name)
}

// ClearLog 清空日志
func (s *App) ClearLog(w http.ResponseWriter, r *http.Request) {
	if _, err := shell.Execf(`echo "" > /var/log/supervisor/supervisord.log`); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}

// GetConfig 获取配置
func (s *App) GetConfig(w http.ResponseWriter, r *http.Request) {
	var config string
	var err error
	if os.IsRHEL() {
		config, err = io.Read(`/etc/supervisord.conf`)
	} else {
		config, err = io.Read(`/etc/supervisor/supervisord.conf`)
	}

	if err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, config)
}

// UpdateConfig 保存配置
func (s *App) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if os.IsRHEL() {
		err = io.Write(`/etc/supervisord.conf`, req.Config, 0644)
	} else {
		err = io.Write(`/etc/supervisor/supervisord.conf`, req.Config, 0644)
	}

	if err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = systemctl.Restart(s.name); err != nil {
		service.Error(w, http.StatusInternalServerError, "重启 %s 服务失败", s.name)
		return
	}

	service.Success(w, nil)
}

// Processes 进程列表
func (s *App) Processes(w http.ResponseWriter, r *http.Request) {
	out, err := shell.Execf(`supervisorctl status | awk '{print $1}'`)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	var processes []Process
	for _, line := range strings.Split(out, "\n") {
		if len(line) == 0 {
			continue
		}

		var p Process
		p.Name = line
		if status, err := shell.Execf(`supervisorctl status '%s' | awk '{print $2}'`, line); err == nil {
			p.Status = status
		}
		if p.Status == "RUNNING" {
			pid, _ := shell.Execf(`supervisorctl status '%s' | awk '{print $4}'`, line)
			p.Pid = strings.ReplaceAll(pid, ",", "")
			uptime, _ := shell.Execf(`supervisorctl status '%s' | awk '{print $6}'`, line)
			p.Uptime = uptime
		} else {
			p.Pid = "-"
			p.Uptime = "-"
		}
		processes = append(processes, p)
	}

	paged, total := service.Paginate(r, processes)

	service.Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

// StartProcess 启动进程
func (s *App) StartProcess(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if _, err = shell.Execf(`supervisorctl start %s`, req.Process); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}

// StopProcess 停止进程
func (s *App) StopProcess(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if _, err = shell.Execf(`supervisorctl stop %s`, req.Process); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}

// RestartProcess 重启进程
func (s *App) RestartProcess(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if _, err = shell.Execf(`supervisorctl restart %s`, req.Process); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}

// ProcessLog 进程日志
func (s *App) ProcessLog(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	var logPath string
	if os.IsRHEL() {
		logPath, err = shell.Execf(`cat '/etc/supervisord.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
	} else {
		logPath, err = shell.Execf(`cat '/etc/supervisor/conf.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
	}

	if err != nil {
		service.Error(w, http.StatusInternalServerError, "无法从进程 %s 的配置文件中获取日志路径", req.Process)
		return
	}

	service.Success(w, logPath)
}

// ClearProcessLog 清空进程日志
func (s *App) ClearProcessLog(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	var logPath string
	if os.IsRHEL() {
		logPath, err = shell.Execf(`cat '/etc/supervisord.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
	} else {
		logPath, err = shell.Execf(`cat '/etc/supervisor/conf.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
	}

	if err != nil {
		service.Error(w, http.StatusInternalServerError, "无法从进程 %s 的配置文件中获取日志路径", req.Process)
		return
	}

	if _, err = shell.Execf(`echo "" > '%s'`, logPath); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}

// ProcessConfig 获取进程配置
func (s *App) ProcessConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	var config string
	if os.IsRHEL() {
		config, err = io.Read(`/etc/supervisord.d/` + req.Process + `.conf`)
	} else {
		config, err = io.Read(`/etc/supervisor/conf.d/` + req.Process + `.conf`)
	}

	if err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, config)
}

// UpdateProcessConfig 保存进程配置
func (s *App) UpdateProcessConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateProcessConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if os.IsRHEL() {
		err = io.Write(`/etc/supervisord.d/`+req.Process+`.conf`, req.Config, 0644)
	} else {
		err = io.Write(`/etc/supervisor/conf.d/`+req.Process+`.conf`, req.Config, 0644)
	}

	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	_, _ = shell.Execf(`supervisorctl reread`)
	_, _ = shell.Execf(`supervisorctl update`)
	_, _ = shell.Execf(`supervisorctl restart '%s'`, req.Process)

	service.Success(w, nil)
}

// CreateProcess 添加进程
func (s *App) CreateProcess(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[CreateProcess](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	config := `[program:` + req.Name + `]
command=` + req.Command + `
process_name=%(program_name)s
directory=` + req.Path + `
autostart=true
autorestart=true
user=` + req.User + `
numprocs=` + cast.ToString(req.Num) + `
redirect_stderr=true
stdout_logfile=/var/log/supervisor/` + req.Name + `.log
stdout_logfile_maxbytes=2MB
`

	if os.IsRHEL() {
		err = io.Write(`/etc/supervisord.d/`+req.Name+`.conf`, config, 0644)
	} else {
		err = io.Write(`/etc/supervisor/conf.d/`+req.Name+`.conf`, config, 0644)
	}

	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	_, _ = shell.Execf(`supervisorctl reread`)
	_, _ = shell.Execf(`supervisorctl update`)
	_, _ = shell.Execf(`supervisorctl start '%s'`, req.Name)

	service.Success(w, nil)
}

// DeleteProcess 删除进程
func (s *App) DeleteProcess(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if _, err = shell.Execf(`supervisorctl stop '%s'`, req.Process); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	var logPath string
	if os.IsRHEL() {
		logPath, err = shell.Execf(`cat '/etc/supervisord.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
		if err != nil {
			service.Error(w, http.StatusInternalServerError, "无法从进程 %s 的配置文件中获取日志路径", req.Process)
			return
		}
		if err = io.Remove(`/etc/supervisord.d/` + req.Process + `.conf`); err != nil {
			service.Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
	} else {
		logPath, err = shell.Execf(`cat '/etc/supervisor/conf.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
		if err != nil {
			service.Error(w, http.StatusInternalServerError, "无法从进程 %s 的配置文件中获取日志路径", req.Process)
			return
		}
		if err = io.Remove(`/etc/supervisor/conf.d/` + req.Process + `.conf`); err != nil {
			service.Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
	}

	if err = io.Remove(logPath); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	_, _ = shell.Execf(`supervisorctl reread`)
	_, _ = shell.Execf(`supervisorctl update`)

	service.Success(w, nil)
}
