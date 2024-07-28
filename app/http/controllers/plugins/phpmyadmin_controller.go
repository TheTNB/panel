package plugins

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/os"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
)

type PhpMyAdminController struct {
}

func NewPhpMyAdminController() *PhpMyAdminController {
	return &PhpMyAdminController{}
}

func (r *PhpMyAdminController) Info(ctx http.Context) http.Response {
	files, err := io.ReadDir("/www/server/phpmyadmin")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "找不到 phpMyAdmin 目录")
	}

	var phpmyadmin string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "phpmyadmin_") {
			phpmyadmin = f.Name()
		}
	}
	if len(phpmyadmin) == 0 {
		return h.Error(ctx, http.StatusInternalServerError, "找不到 phpMyAdmin 目录")
	}

	conf, err := io.Read("/www/server/vhost/phpmyadmin.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	match := regexp.MustCompile(`listen\s+(\d+);`).FindStringSubmatch(conf)
	if len(match) == 0 {
		return h.Error(ctx, http.StatusInternalServerError, "找不到 phpMyAdmin 端口")
	}

	return h.Success(ctx, http.Json{
		"path": phpmyadmin,
		"port": cast.ToInt(match[1]),
	})
}

func (r *PhpMyAdminController) SetPort(ctx http.Context) http.Response {
	port := ctx.Request().InputInt("port")
	if port == 0 {
		return h.Error(ctx, http.StatusInternalServerError, "端口不能为空")
	}

	conf, err := io.Read("/www/server/vhost/phpmyadmin.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	conf = regexp.MustCompile(`listen\s+(\d+);`).ReplaceAllString(conf, "listen "+cast.ToString(port)+";")
	if err := io.Write("/www/server/vhost/phpmyadmin.conf", conf, 0644); err != nil {
		facades.Log().Request(ctx.Request()).Tags("插件", "phpMyAdmin").With(map[string]any{
			"error": err.Error(),
		}).Info("修改 phpMyAdmin 端口失败")
		return h.ErrorSystem(ctx)
	}

	if os.IsRHEL() {
		if out, err := shell.Execf("firewall-cmd --zone=public --add-port=%d/tcp --permanent", port); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := shell.Execf("firewall-cmd --reload"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	} else {
		if out, err := shell.Execf("ufw allow %d/tcp", port); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := shell.Execf("ufw reload"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	if err = systemctl.Reload("openresty"); err != nil {
		_, err = shell.Execf("openresty -t")
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("重载OpenResty失败: %v", err))
	}

	return h.Success(ctx, nil)
}

func (r *PhpMyAdminController) GetConfig(ctx http.Context) http.Response {
	config, err := io.Read("/www/server/vhost/phpmyadmin.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, config)
}

func (r *PhpMyAdminController) SaveConfig(ctx http.Context) http.Response {
	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return h.Error(ctx, http.StatusInternalServerError, "配置不能为空")
	}

	if err := io.Write("/www/server/vhost/phpmyadmin.conf", config, 0644); err != nil {
		facades.Log().Request(ctx.Request()).Tags("插件", "phpMyAdmin").With(map[string]any{
			"error": err.Error(),
		}).Info("修改 phpMyAdmin 配置失败")
		return h.ErrorSystem(ctx)
	}

	if err := systemctl.Reload("openresty"); err != nil {
		_, err = shell.Execf("openresty -t")
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("重载OpenResty失败: %v", err))
	}

	return h.Success(ctx, nil)
}
