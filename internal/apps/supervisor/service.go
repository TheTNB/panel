package supervisor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-rat/chix"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/os"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type Service struct {
	name string
}

func NewService() *Service {
	var name string
	if os.IsRHEL() {
		name = "supervisord"
	} else {
		name = "supervisor"
	}

	return &Service{
		name: name,
	}
}

// Service 获取服务名称
func (s *Service) Service(w http.ResponseWriter, r *http.Request) {
	service.Success(w, s.name)
}

// Log 日志
func (s *Service) Log(w http.ResponseWriter, r *http.Request) {
	log, err := shell.Execf(`tail -n 200 /var/log/supervisor/supervisord.log`)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, log)
		return
	}

	service.Success(w, log)
}

// ClearLog 清空日志
func (s *Service) ClearLog(w http.ResponseWriter, r *http.Request) {
	if out, err := shell.Execf(`echo "" > /var/log/supervisor/supervisord.log`); err != nil {
		service.Error(w, http.StatusInternalServerError, out)
		return
	}

	service.Success(w, nil)
}

// GetConfig 获取配置
func (s *Service) GetConfig(w http.ResponseWriter, r *http.Request) {
	var config string
	var err error
	if os.IsRHEL() {
		config, err = io.Read(`/etc/supervisord.conf`)
	} else {
		config, err = io.Read(`/etc/supervisor/supervisord.conf`)
	}

	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, config)
}

// UpdateConfig 保存配置
func (s *Service) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if os.IsRHEL() {
		err = io.Write(`/etc/supervisord.conf`, req.Config, 0644)
	} else {
		err = io.Write(`/etc/supervisor/supervisord.conf`, req.Config, 0644)
	}

	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = systemctl.Restart(s.name); err != nil {
		service.Error(w, http.StatusInternalServerError, fmt.Sprintf("重启 %s 服务失败", s.name))
		return
	}

	service.Success(w, nil)
}

// Processes 进程列表
func (s *Service) Processes(w http.ResponseWriter, r *http.Request) {
	out, err := shell.Execf(`supervisorctl status | awk '{print $1}'`)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, out)
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
func (s *Service) StartProcess(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if out, err := shell.Execf(`supervisorctl start %s`, req.Process); err != nil {
		service.Error(w, http.StatusInternalServerError, out)
		return
	}

	service.Success(w, nil)
}

// StopProcess 停止进程
func (s *Service) StopProcess(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if out, err := shell.Execf(`supervisorctl stop %s`, req.Process); err != nil {
		service.Error(w, http.StatusInternalServerError, out)
		return
	}

	service.Success(w, nil)
}

// RestartProcess 重启进程
func (s *Service) RestartProcess(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if out, err := shell.Execf(`supervisorctl restart %s`, req.Process); err != nil {
		service.Error(w, http.StatusInternalServerError, out)
		return
	}

	service.Success(w, nil)
}

// ProcessLog 进程日志
func (s *Service) ProcessLog(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var logPath string
	if os.IsRHEL() {
		logPath, err = shell.Execf(`cat '/etc/supervisord.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
	} else {
		logPath, err = shell.Execf(`cat '/etc/supervisor/conf.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
	}

	if err != nil {
		service.Error(w, http.StatusInternalServerError, fmt.Sprintf("无法从进程 %s 的配置文件中获取日志路径", req.Process))
		return
	}

	log, err := shell.Execf(`tail -n 200 '%s'`, logPath)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, log)
		return
	}

	service.Success(w, log)
}

// ClearProcessLog 清空进程日志
func (s *Service) ClearProcessLog(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var logPath string
	if os.IsRHEL() {
		logPath, err = shell.Execf(`cat '/etc/supervisord.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
	} else {
		logPath, err = shell.Execf(`cat '/etc/supervisor/conf.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
	}

	if err != nil {
		service.Error(w, http.StatusInternalServerError, fmt.Sprintf("无法从进程 %s 的配置文件中获取日志路径", req.Process))
		return
	}

	if out, err := shell.Execf(`echo "" > '%s'`, logPath); err != nil {
		service.Error(w, http.StatusInternalServerError, out)
		return
	}

	service.Success(w, nil)
}

// ProcessConfig 获取进程配置
func (s *Service) ProcessConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var config string
	if os.IsRHEL() {
		config, err = io.Read(`/etc/supervisord.d/` + req.Process + `.conf`)
	} else {
		config, err = io.Read(`/etc/supervisor/conf.d/` + req.Process + `.conf`)
	}

	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, config)
}

// UpdateProcessConfig 保存进程配置
func (s *Service) UpdateProcessConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateProcessConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if os.IsRHEL() {
		err = io.Write(`/etc/supervisord.d/`+req.Process+`.conf`, req.Config, 0644)
	} else {
		err = io.Write(`/etc/supervisor/conf.d/`+req.Process+`.conf`, req.Config, 0644)
	}

	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, _ = shell.Execf(`supervisorctl reread`)
	_, _ = shell.Execf(`supervisorctl update`)
	_, _ = shell.Execf(`supervisorctl restart '%s'`, req.Process)

	service.Success(w, nil)
}

// CreateProcess 添加进程
func (s *Service) CreateProcess(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[CreateProcess](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
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
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, _ = shell.Execf(`supervisorctl reread`)
	_, _ = shell.Execf(`supervisorctl update`)
	_, _ = shell.Execf(`supervisorctl start '%s'`, req.Name)

	service.Success(w, nil)
}

// DeleteProcess 删除进程
func (s *Service) DeleteProcess(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ProcessName](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if out, err := shell.Execf(`supervisorctl stop '%s'`, req.Process); err != nil {
		service.Error(w, http.StatusInternalServerError, out)
		return
	}

	var logPath string
	if os.IsRHEL() {
		logPath, err = shell.Execf(`cat '/etc/supervisord.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
		if err := io.Remove(`/etc/supervisord.d/` + req.Process + `.conf`); err != nil {
			service.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		logPath, err = shell.Execf(`cat '/etc/supervisor/conf.d/%s.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`, req.Process)
		if err := io.Remove(`/etc/supervisor/conf.d/` + req.Process + `.conf`); err != nil {
			service.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err != nil {
		service.Error(w, http.StatusInternalServerError, fmt.Sprintf("无法从进程 %s 的配置文件中获取日志路径", req.Process))
		return
	}

	if err = io.Remove(logPath); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, _ = shell.Execf(`supervisorctl reread`)
	_, _ = shell.Execf(`supervisorctl update`)

	service.Success(w, nil)
}
