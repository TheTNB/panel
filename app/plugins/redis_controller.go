package plugins

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type RedisController struct {
}

func NewRedisController() *RedisController {
	return &RedisController{}
}

// GetConfig 获取配置
func (r *RedisController) GetConfig(ctx http.Context) http.Response {
	// 获取配置
	config, err := io.Read("/www/server/redis/redis.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取Redis配置失败")
	}

	return h.Success(ctx, config)
}

// SaveConfig 保存配置
func (r *RedisController) SaveConfig(ctx http.Context) http.Response {
	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "配置不能为空")
	}

	if err := io.Write("/www/server/redis/redis.conf", config, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "写入Redis配置失败")
	}

	if err := systemctl.Restart("redis"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "重启Redis失败")
	}

	return h.Success(ctx, nil)
}

// Load 获取负载
func (r *RedisController) Load(ctx http.Context) http.Response {
	status, err := systemctl.Status("redis")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
	}
	if !status {
		return h.Error(ctx, http.StatusInternalServerError, "Redis已停止运行")
	}

	raw, err := shell.Execf("redis-cli info")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取Redis负载失败")
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

	return h.Success(ctx, data)
}
