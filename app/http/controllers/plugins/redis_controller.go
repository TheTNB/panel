package plugins

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"github.com/TheTNB/panel/app/http/controllers"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/types"
)

type RedisController struct {
}

func NewRedisController() *RedisController {
	return &RedisController{}
}

// Status 获取运行状态
func (r *RedisController) Status(ctx http.Context) http.Response {
	status, err := tools.ServiceStatus("redis")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
	}

	return controllers.Success(ctx, status)
}

// Restart 重启服务
func (r *RedisController) Restart(ctx http.Context) http.Response {
	if err := tools.ServiceRestart("redis"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "重启Redis失败")
	}

	return controllers.Success(ctx, nil)
}

// Start 启动服务
func (r *RedisController) Start(ctx http.Context) http.Response {
	if err := tools.ServiceStart("redis"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "启动Redis失败")
	}

	return controllers.Success(ctx, nil)
}

// Stop 停止服务
func (r *RedisController) Stop(ctx http.Context) http.Response {
	if err := tools.ServiceStop("redis"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "停止Redis失败")
	}

	return controllers.Success(ctx, nil)
}

// GetConfig 获取配置
func (r *RedisController) GetConfig(ctx http.Context) http.Response {
	// 获取配置
	config, err := tools.Read("/www/server/redis/redis.conf")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取Redis配置失败")
	}

	return controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (r *RedisController) SaveConfig(ctx http.Context) http.Response {
	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "配置不能为空")
	}

	if err := tools.Write("/www/server/redis/redis.conf", config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "写入Redis配置失败")
	}

	return r.Restart(ctx)
}

// Load 获取负载
func (r *RedisController) Load(ctx http.Context) http.Response {
	status, err := tools.ServiceStatus("redis")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
	}
	if !status {
		return controllers.Error(ctx, http.StatusInternalServerError, "Redis已停止运行")
	}

	raw, err := tools.Exec("redis-cli info")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取Redis负载失败")
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

	return controllers.Success(ctx, data)
}
