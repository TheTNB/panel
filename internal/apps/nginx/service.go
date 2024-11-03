package nginx

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/systemctl"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/pkg/types"
)

type Service struct {
	// Dependent services
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read(fmt.Sprintf("%s/server/nginx/conf/nginx.conf", app.Root))
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取配置失败")
		return
	}

	service.Success(w, config)
}

func (s *Service) SaveConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = io.Write(fmt.Sprintf("%s/server/nginx/conf/nginx.conf", app.Root), req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "保存配置失败")
		return
	}

	if err = systemctl.Reload("nginx"); err != nil {
		_, err = shell.Execf("nginx -t")
		service.Error(w, http.StatusInternalServerError, "重载服务失败：%v", err)
		return
	}

	service.Success(w, nil)
}

func (s *Service) ErrorLog(w http.ResponseWriter, r *http.Request) {
	service.Success(w, fmt.Sprintf("%s/%s", app.Root, "wwwlogs/nginx-error.log"))
}

func (s *Service) ClearErrorLog(w http.ResponseWriter, r *http.Request) {
	if _, err := shell.Execf("echo '' > %s/%s", app.Root, "wwwlogs/nginx-error.log"); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}

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
		return
	}
	data = append(data, types.NV{
		Name:  "工作进程",
		Value: workers,
	})

	out, err := shell.Execf("ps aux | grep nginx | grep 'worker process' | awk '{memsum+=$6};END {print memsum}'")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取负载失败")
		return
	}
	mem := tools.FormatBytes(cast.ToFloat64(out))
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
