package plugins

import (
	"regexp"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/imroc/req/v3"
	"github.com/spf13/cast"

	"panel/app/http/controllers"
	"panel/pkg/tools"
)

type OpenRestyController struct {
	// Dependent services
}

func NewOpenrestyController() *OpenRestyController {
	return &OpenRestyController{
		// Inject services
	}
}

// Status 获取运行状态
func (r *OpenRestyController) Status(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "openresty")
	if check != nil {
		return check
	}

	status, err := tools.ServiceStatus("openresty")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty状态失败")
	}

	return controllers.Success(ctx, status)
}

// Reload 重载配置
func (r *OpenRestyController) Reload(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "openresty")
	if check != nil {
		return check
	}

	if err := tools.ServiceReload("openresty"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "重载OpenResty失败")
	}

	return controllers.Success(ctx, nil)
}

// Start 启动OpenResty
func (r *OpenRestyController) Start(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "openresty")
	if check != nil {
		return check
	}

	if err := tools.ServiceStart("openresty"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "启动OpenResty失败")
	}

	return controllers.Success(ctx, nil)
}

// Stop 停止OpenResty
func (r *OpenRestyController) Stop(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "openresty")
	if check != nil {
		return check
	}

	if err := tools.ServiceStop("openresty"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "停止OpenResty失败")
	}

	return controllers.Success(ctx, nil)
}

// Restart 重启OpenResty
func (r *OpenRestyController) Restart(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "openresty")
	if check != nil {
		return check
	}

	if err := tools.ServiceRestart("openresty"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "重启OpenResty失败")
	}

	return controllers.Success(ctx, nil)
}

// GetConfig 获取配置
func (r *OpenRestyController) GetConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "openresty")
	if check != nil {
		return check
	}

	config := tools.Read("/www/server/openresty/conf/nginx.conf")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty配置失败")
	}

	return controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (r *OpenRestyController) SaveConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "openresty")
	if check != nil {
		return check
	}

	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "配置不能为空")
	}

	if err := tools.Write("/www/server/openresty/conf/nginx.conf", config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "保存OpenResty配置失败")
	}

	return r.Reload(ctx)
}

// ErrorLog 获取错误日志
func (r *OpenRestyController) ErrorLog(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "openresty")
	if check != nil {
		return check
	}

	if !tools.Exists("/www/wwwlogs/nginx_error.log") {
		return controllers.Success(ctx, "")
	}

	out, err := tools.Exec("tail -n 100 /www/wwwlogs/nginx_error.log")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, out)
}

// ClearErrorLog 清空错误日志
func (r *OpenRestyController) ClearErrorLog(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "openresty")
	if check != nil {
		return check
	}

	if out, err := tools.Exec("echo '' > /www/wwwlogs/nginx_error.log"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// Load 获取负载
func (r *OpenRestyController) Load(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "openresty")
	if check != nil {
		return check
	}

	client := req.C().SetTimeout(10 * time.Second)
	resp, err := client.R().Get("http://127.0.0.1/nginx_status")
	if err != nil || !resp.IsSuccessState() {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty负载失败")
	}

	raw := resp.String()
	type nginxStatus struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	var data []nginxStatus

	workers, err := tools.Exec("ps aux | grep nginx | grep 'worker process' | wc -l")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty负载失败")
	}
	data = append(data, nginxStatus{
		Name:  "工作进程",
		Value: workers,
	})

	out, err := tools.Exec("ps aux | grep nginx | grep 'worker process' | awk '{memsum+=$6};END {print memsum}'")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty负载失败")
	}
	mem := tools.FormatBytes(cast.ToFloat64(out))
	data = append(data, nginxStatus{
		Name:  "内存占用",
		Value: mem,
	})

	match := regexp.MustCompile(`Active connections:\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) == 2 {
		data = append(data, nginxStatus{
			Name:  "活跃连接数",
			Value: match[1],
		})
	}

	match = regexp.MustCompile(`server accepts handled requests\s+(\d+)\s+(\d+)\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) == 4 {
		data = append(data, nginxStatus{
			Name:  "总连接次数",
			Value: match[1],
		})
		data = append(data, nginxStatus{
			Name:  "总握手次数",
			Value: match[2],
		})
		data = append(data, nginxStatus{
			Name:  "总请求次数",
			Value: match[3],
		})
	}

	match = regexp.MustCompile(`Reading:\s+(\d+)\s+Writing:\s+(\d+)\s+Waiting:\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) == 4 {
		data = append(data, nginxStatus{
			Name:  "请求数",
			Value: match[1],
		})
		data = append(data, nginxStatus{
			Name:  "响应数",
			Value: match[2],
		})
		data = append(data, nginxStatus{
			Name:  "驻留进程",
			Value: match[3],
		})
	}

	return controllers.Success(ctx, data)
}
