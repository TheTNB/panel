package openresty

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/str"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type Controller struct {
	// Dependent services
}

func NewController() *Controller {
	return &Controller{}
}

// GetConfig
//
//	@Summary	获取配置
//	@Tags		插件-OpenResty
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	h.SuccessResponse
//	@Router		/plugins/openresty/config [get]
func (r *Controller) GetConfig(ctx http.Context) http.Response {
	config, err := io.Read("/www/server/openresty/conf/nginx.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取配置失败")
	}

	return h.Success(ctx, config)
}

// SaveConfig
//
//	@Summary	保存配置
//	@Tags		插件-OpenResty
//	@Produce	json
//	@Security	BearerToken
//	@Param		config	body		string	true	"配置"
//	@Success	200		{object}	h.SuccessResponse
//	@Router		/plugins/openresty/config [post]
func (r *Controller) SaveConfig(ctx http.Context) http.Response {
	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return h.Error(ctx, http.StatusInternalServerError, "配置不能为空")
	}

	if err := io.Write("/www/server/openresty/conf/nginx.conf", config, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "保存配置失败")
	}

	if err := systemctl.Reload("openresty"); err != nil {
		_, err = shell.Execf("openresty -t")
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("重载服务失败: %v", err))
	}

	return h.Success(ctx, nil)
}

// ErrorLog
//
//	@Summary	获取错误日志
//	@Tags		插件-OpenResty
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	h.SuccessResponse
//	@Router		/plugins/openresty/errorLog [get]
func (r *Controller) ErrorLog(ctx http.Context) http.Response {
	if !io.Exists("/www/wwwlogs/nginx_error.log") {
		return h.Success(ctx, "")
	}

	out, err := shell.Execf("tail -n 100 /www/wwwlogs/openresty_error.log")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, out)
}

// ClearErrorLog
//
//	@Summary	清空错误日志
//	@Tags		插件-OpenResty
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	h.SuccessResponse
//	@Router		/plugins/openresty/clearErrorLog [post]
func (r *Controller) ClearErrorLog(ctx http.Context) http.Response {
	if out, err := shell.Execf("echo '' > /www/wwwlogs/openresty_error.log"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}

// Load
//
//	@Summary	获取负载状态
//	@Tags		插件-OpenResty
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	h.SuccessResponse
//	@Router		/plugins/openresty/load [get]
func (r *Controller) Load(ctx http.Context) http.Response {
	client := resty.New().SetTimeout(10 * time.Second)
	resp, err := client.R().Get("http://127.0.0.1/nginx_status")
	if err != nil || !resp.IsSuccess() {
		return h.Success(ctx, []types.NV{})
	}

	raw := resp.String()
	var data []types.NV

	workers, err := shell.Execf("ps aux | grep nginx | grep 'worker process' | wc -l")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取负载失败")
	}
	data = append(data, types.NV{
		Name:  "工作进程",
		Value: workers,
	})

	out, err := shell.Execf("ps aux | grep nginx | grep 'worker process' | awk '{memsum+=$6};END {print memsum}'")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取负载失败")
	}
	mem := str.FormatBytes(cast.ToFloat64(out))
	data = append(data, types.NV{
		Name:  "内存占用",
		Value: mem,
	})

	match := regexp.MustCompile(`Active connections:\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) == 2 {
		data = append(data, types.NV{
			Name:  "活跃连接数",
			Value: match[1],
		})
	}

	match = regexp.MustCompile(`server accepts handled requests\s+(\d+)\s+(\d+)\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) == 4 {
		data = append(data, types.NV{
			Name:  "总连接次数",
			Value: match[1],
		})
		data = append(data, types.NV{
			Name:  "总握手次数",
			Value: match[2],
		})
		data = append(data, types.NV{
			Name:  "总请求次数",
			Value: match[3],
		})
	}

	match = regexp.MustCompile(`Reading:\s+(\d+)\s+Writing:\s+(\d+)\s+Waiting:\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) == 4 {
		data = append(data, types.NV{
			Name:  "请求数",
			Value: match[1],
		})
		data = append(data, types.NV{
			Name:  "响应数",
			Value: match[2],
		})
		data = append(data, types.NV{
			Name:  "驻留进程",
			Value: match[3],
		})
	}

	return h.Success(ctx, data)
}
