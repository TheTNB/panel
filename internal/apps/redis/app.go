package redis

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/tnb-labs/panel/internal/app"
	"github.com/tnb-labs/panel/internal/service"
	"github.com/tnb-labs/panel/pkg/io"
	"github.com/tnb-labs/panel/pkg/shell"
	"github.com/tnb-labs/panel/pkg/systemctl"
	"github.com/tnb-labs/panel/pkg/types"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (s *App) Route(r chi.Router) {
	r.Get("/load", s.Load)
	r.Get("/config", s.GetConfig)
	r.Post("/config", s.UpdateConfig)
}

func (s *App) Load(w http.ResponseWriter, r *http.Request) {
	status, err := systemctl.Status("redis")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取 Redis 状态失败")
		return
	}
	if !status {
		service.Success(w, []types.NV{})
		return
	}

	raw, err := shell.Execf("redis-cli info")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取 Redis 负载失败")
		return
	}

	infoLines := strings.Split(raw, "\n")
	dataRaw := make(map[string]string)

	for _, item := range infoLines {
		parts := strings.Split(item, ":")
		if len(parts) == 2 {
			dataRaw[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	data := []types.NV{
		{Name: "TCP 端口", Value: dataRaw["tcp_port"]},
		{Name: "已运行天数", Value: dataRaw["uptime_in_days"]},
		{Name: "连接的客户端数", Value: dataRaw["connected_clients"]},
		{Name: "已分配的内存总量", Value: dataRaw["used_memory_human"]},
		{Name: "占用内存总量", Value: dataRaw["used_memory_rss_human"]},
		{Name: "占用内存峰值", Value: dataRaw["used_memory_peak_human"]},
		{Name: "内存碎片比率", Value: dataRaw["mem_fragmentation_ratio"]},
		{Name: "运行以来连接过的客户端的总数", Value: dataRaw["total_connections_received"]},
		{Name: "运行以来执行过的命令的总数", Value: dataRaw["total_commands_processed"]},
		{Name: "每秒执行的命令数", Value: dataRaw["instantaneous_ops_per_sec"]},
		{Name: "查找数据库键成功次数", Value: dataRaw["keyspace_hits"]},
		{Name: "查找数据库键失败次数", Value: dataRaw["keyspace_misses"]},
		{Name: "最近一次 fork() 操作耗费的毫秒数", Value: dataRaw["latest_fork_usec"]},
	}

	service.Success(w, data)
}

func (s *App) GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read(fmt.Sprintf("%s/server/redis/redis.conf", app.Root))
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, config)
}

func (s *App) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = io.Write(fmt.Sprintf("%s/server/redis/redis.conf", app.Root), req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = systemctl.Restart("redis"); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}
