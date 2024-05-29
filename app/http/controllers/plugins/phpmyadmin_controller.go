package plugins

import (
	"os"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/app/http/controllers"
	"github.com/TheTNB/panel/pkg/tools"
)

type PhpMyAdminController struct {
}

func NewPhpMyAdminController() *PhpMyAdminController {
	return &PhpMyAdminController{}
}

func (r *PhpMyAdminController) Info(ctx http.Context) http.Response {
	files, err := os.ReadDir("/www/server/phpmyadmin")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "找不到 phpMyAdmin 目录")
	}

	var phpmyadmin string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "phpmyadmin_") {
			phpmyadmin = f.Name()
		}
	}
	if len(phpmyadmin) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "找不到 phpMyAdmin 目录")
	}

	conf, err := tools.Read("/www/server/vhost/phpmyadmin.conf")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	match := regexp.MustCompile(`listen\s+(\d+);`).FindStringSubmatch(conf)
	if len(match) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "找不到 phpMyAdmin 端口")
	}

	return controllers.Success(ctx, http.Json{
		"path": phpmyadmin,
		"port": cast.ToInt(match[1]),
	})
}

func (r *PhpMyAdminController) SetPort(ctx http.Context) http.Response {
	port := ctx.Request().Input("port")
	if len(port) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "端口不能为空")
	}

	conf, err := tools.Read("/www/server/vhost/phpmyadmin.conf")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	conf = regexp.MustCompile(`listen\s+(\d+);`).ReplaceAllString(conf, "listen "+port+";")
	if err := tools.Write("/www/server/vhost/phpmyadmin.conf", conf, 0644); err != nil {
		facades.Log().Request(ctx.Request()).Tags("插件", "phpMyAdmin").With(map[string]any{
			"error": err.Error(),
		}).Info("修改 phpMyAdmin 端口失败")
		return controllers.ErrorSystem(ctx)
	}

	if tools.IsRHEL() {
		if out, err := tools.Exec("firewall-cmd --zone=public --add-port=" + port + "/tcp --permanent"); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := tools.Exec("firewall-cmd --reload"); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, out)
		}
	} else {
		if out, err := tools.Exec("ufw allow " + port + "/tcp"); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := tools.Exec("ufw reload"); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	if err := tools.ServiceReload("openresty"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "重载OpenResty失败")
	}

	return controllers.Success(ctx, nil)
}
