package postgresql

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TheTNB/panel/internal/panel"
	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/systemctl"
	"github.com/TheTNB/panel/pkg/types"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// GetConfig 获取配置
func (s *Service) GetConfig(w http.ResponseWriter, r *http.Request) {
	// 获取配置
	config, err := io.Read(fmt.Sprintf("%s/server/postgresql/data/postgresql.conf", panel.Root))
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取PostgreSQL配置失败")
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

	if err := io.Write(fmt.Sprintf("%s/server/postgresql/data/postgresql.conf", panel.Root), req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "写入PostgreSQL配置失败")
		return
	}

	if err := systemctl.Reload("postgresql"); err != nil {
		service.Error(w, http.StatusInternalServerError, "重载服务失败")
		return
	}

	service.Success(w, nil)
}

// GetUserConfig 获取用户配置
func (s *Service) GetUserConfig(w http.ResponseWriter, r *http.Request) {
	// 获取配置
	config, err := io.Read(fmt.Sprintf("%s/server/postgresql/data/pg_hba.conf", panel.Root))
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取PostgreSQL配置失败")
		return
	}

	service.Success(w, config)
}

// UpdateUserConfig 保存用户配置
func (s *Service) UpdateUserConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := io.Write(fmt.Sprintf("%s/server/postgresql/data/pg_hba.conf", panel.Root), req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "写入PostgreSQL配置失败")
		return
	}

	if err := systemctl.Reload("postgresql"); err != nil {
		service.Error(w, http.StatusInternalServerError, "重载服务失败")
		return
	}

	service.Success(w, nil)
}

// Load 获取负载
func (s *Service) Load(w http.ResponseWriter, r *http.Request) {
	status, _ := systemctl.Status("postgresql")
	if !status {
		service.Success(w, []types.NV{})
		return
	}

	start, err := shell.Execf(`echo "select pg_postmaster_start_time();" | su - postgres -c "psql" | sed -n 3p | cut -d'.' -f1`)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取PostgreSQL启动时间失败")
		return
	}
	pid, err := shell.Execf(`echo "select pg_backend_pid();" | su - postgres -c "psql" | sed -n 3p`)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取PostgreSQL进程PID失败")
		return
	}
	process, err := shell.Execf(`ps aux | grep postgres | grep -v grep | wc -l`)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取PostgreSQL进程数失败")
		return
	}
	connections, err := shell.Execf(`echo "SELECT count(*) FROM pg_stat_activity WHERE NOT pid=pg_backend_pid();" | su - postgres -c "psql" | sed -n 3p`)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取PostgreSQL连接数失败")
		return
	}
	storage, err := shell.Execf(`echo "select pg_size_pretty(pg_database_size('postgres'));" | su - postgres -c "psql" | sed -n 3p`)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取PostgreSQL空间占用失败")
		return
	}

	data := []types.NV{
		{Name: "启动时间", Value: start},
		{Name: "进程 PID", Value: pid},
		{Name: "进程数", Value: process},
		{Name: "总连接数", Value: connections},
		{Name: "空间占用", Value: storage},
	}

	service.Success(w, data)
}

// Log 获取日志
func (s *Service) Log(w http.ResponseWriter, r *http.Request) {
	log, err := shell.Execf("tail -n 100 %s/server/postgresql/logs/postgresql-%s.log", panel.Root, time.Now().Format(time.DateOnly))
	if err != nil {
		service.Error(w, http.StatusInternalServerError, log)
		return
	}

	service.Success(w, log)
}

// ClearLog 清空日志
func (s *Service) ClearLog(w http.ResponseWriter, r *http.Request) {
	if _, err := shell.Execf("rm -rf %s/server/postgresql/logs/postgresql-*.log", panel.Root); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
}
