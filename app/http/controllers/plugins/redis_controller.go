package plugins

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"panel/app/http/controllers"
	"panel/pkg/tools"
)

type RedisController struct {
}

func NewRedisController() *RedisController {
	return &RedisController{}
}

// Status 获取运行状态
func (r *RedisController) Status(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "redis")
	if check != nil {
		return check
	}

	status := tools.Exec("systemctl status redis | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Restart 重启服务
func (r *RedisController) Restart(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "redis")
	if check != nil {
		return check
	}

	tools.Exec("systemctl restart redis")
	status := tools.Exec("systemctl status redis | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Start 启动服务
func (r *RedisController) Start(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "redis")
	if check != nil {
		return check
	}

	tools.Exec("systemctl start redis")
	status := tools.Exec("systemctl status redis | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Stop 停止服务
func (r *RedisController) Stop(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "redis")
	if check != nil {
		return check
	}

	tools.Exec("systemctl stop redis")
	status := tools.Exec("systemctl status redis | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取Redis状态失败")
	}

	if status != "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// GetConfig 获取配置
func (r *RedisController) GetConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "redis")
	if check != nil {
		return check
	}

	// 获取配置
	config := tools.Read("/www/server/redis/redis.conf")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取Redis配置失败")
	}

	return controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (r *RedisController) SaveConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "redis")
	if check != nil {
		return check
	}

	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "配置不能为空")
	}

	if !tools.Write("/www/server/redis/redis.conf", config, 0644) {
		return controllers.Error(ctx, http.StatusInternalServerError, "写入Redis配置失败")
	}

	return r.Restart(ctx)
}

// Load 获取负载
func (r *RedisController) Load(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "redis")
	if check != nil {
		return check
	}

	status := tools.Exec("systemctl status redis | grep Active | grep -v grep | awk '{print $2}'")
	if status != "active" {
		return controllers.Error(ctx, http.StatusInternalServerError, "Redis 已停止运行")
	}

	raw := tools.Exec("redis-cli info")
	if len(raw) == 0 {
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

	data := []LoadInfo{
		{"TCP 端口", dataRaw["tcp_port"]},
		{"已运行天数", dataRaw["uptime_in_days"]},
		{"连接的客户端数", dataRaw["connected_clients"]},
		{"已分配的内存总量", dataRaw["used_memory_human"]},
		{"占用内存总量", dataRaw["used_memory_rss_human"]},
		{"占用内存峰值", dataRaw["used_memory_peak_human"]},
		{"内存碎片比率", dataRaw["mem_fragmentation_ratio"]},
		{"运行以来连接过的客户端的总数", dataRaw["total_connections_received"]},
		{"运行以来执行过的命令的总数", dataRaw["total_commands_processed"]},
		{"每秒执行的命令数", dataRaw["instantaneous_ops_per_sec"]},
		{"查找数据库键成功次数", dataRaw["keyspace_hits"]},
		{"查找数据库键失败次数", dataRaw["keyspace_misses"]},
		{"最近一次 fork() 操作耗费的毫秒数", dataRaw["latest_fork_usec"]},
	}

	return controllers.Success(ctx, data)
}
