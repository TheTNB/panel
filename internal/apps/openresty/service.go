package openresty

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/panel"
	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/str"
	"github.com/TheTNB/panel/pkg/systemctl"
	"github.com/TheTNB/panel/pkg/types"
)

type Service struct {
	// Dependent services
}

func NewService() *Service {
	return &Service{}
}

// GetConfig
//
//	@Summary	获取配置
//	@Tags		插件-OpenResty
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	h.SuccessResponse
//	@Router		/plugins/openresty/config [get]
func (s *Service) GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read(fmt.Sprintf("%s/server/openresty/conf/nginx.conf", panel.Root))
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取配置失败")
		return
	}

	service.Success(w, config)
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
func (s *Service) SaveConfig(w http.ResponseWriter, r *http.Request) {
	config := r.FormValue("config")
	if len(config) == 0 {
		service.Error(w, http.StatusInternalServerError, "配置不能为空")
	}

	if err := io.Write(fmt.Sprintf("%s/server/openresty/conf/nginx.conf", panel.Root), config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "保存配置失败")
	}

	if err := systemctl.Reload("openresty"); err != nil {
		_, err = shell.Execf("openresty -t")
		service.Error(w, http.StatusInternalServerError, fmt.Sprintf("重载服务失败: %v", err))
	}

	service.Success(w, nil)
}

// ErrorLog
//
//	@Summary	获取错误日志
//	@Tags		插件-OpenResty
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	h.SuccessResponse
//	@Router		/plugins/openresty/errorLog [get]
func (s *Service) ErrorLog(w http.ResponseWriter, r *http.Request) {
	if !io.Exists(fmt.Sprintf("%s/wwwlogs/nginx_error.log", panel.Root)) {
		service.Success(w, "")
	}

	out, err := shell.Execf("tail -n 100 %s/%s", panel.Root, "/wwwlogs/openresty_error.log")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
	}

	service.Success(w, out)
}

// ClearErrorLog
//
//	@Summary	清空错误日志
//	@Tags		插件-OpenResty
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	h.SuccessResponse
//	@Router		/plugins/openresty/clearErrorLog [post]
func (s *Service) ClearErrorLog(w http.ResponseWriter, r *http.Request) {
	if _, err := shell.Execf("echo '' > %s/%s", panel.Root, "/wwwlogs/openresty_error.log"); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
	}

	service.Success(w, nil)
}

// Load
//
//	@Summary	获取负载状态
//	@Tags		插件-OpenResty
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	h.SuccessResponse
//	@Router		/plugins/openresty/load [get]
func (s *Service) Load(w http.ResponseWriter, r *http.Request) {
	client := resty.New().SetTimeout(10 * time.Second)
	resp, err := client.R().Get("http://127.0.0.1/nginx_status")
	if err != nil || !resp.IsSuccess() {
		service.Success(w, []types.NV{})
	}

	raw := resp.String()
	var data []types.NV

	workers, err := shell.Execf("ps aux | grep nginx | grep 'worker process' | wc -l")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取负载失败")
	}
	data = append(data, types.NV{
		Name:  "工作进程",
		Value: workers,
	})

	out, err := shell.Execf("ps aux | grep nginx | grep 'worker process' | awk '{memsum+=$6};END {print memsum}'")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取负载失败")
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

	service.Success(w, data)
}
