package phpmyadmin

import (
	"os"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"panel/app/http/controllers"
	"panel/pkg/tools"
)

type PhpMyAdminController struct {
}

func NewPhpMyAdminController() *PhpMyAdminController {
	return &PhpMyAdminController{}
}

func (c *PhpMyAdminController) Info(ctx http.Context) {
	if !controllers.Check(ctx, "phpmyadmin") {
		return
	}

	files, err := os.ReadDir("/www/server/phpmyadmin")
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "找不到 phpMyAdmin 目录")
		return
	}

	var phpmyadmin string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "phpmyadmin_") {
			phpmyadmin = f.Name()
		}
	}
	if len(phpmyadmin) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "找不到 phpMyAdmin 目录")
		return
	}

	conf := tools.ReadFile("/www/server/vhost/phpmyadmin.conf")
	match := regexp.MustCompile(`listen\s+(\d+);`).FindStringSubmatch(conf)
	if len(match) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "找不到 phpMyAdmin 端口")
		return
	}

	controllers.Success(ctx, http.Json{
		"phpmyadmin": phpmyadmin,
		"port":       match[1],
	})
}

func (c *PhpMyAdminController) SetPort(ctx http.Context) {
	if !controllers.Check(ctx, "phpmyadmin") {
		return
	}

	port := ctx.Request().Input("port")
	if len(port) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "端口不能为空")
		return
	}

	conf := tools.ReadFile("/www/server/vhost/phpmyadmin.conf")
	conf = regexp.MustCompile(`listen\s+(\d+);`).ReplaceAllString(conf, "listen "+port+";")
	tools.WriteFile("/www/server/vhost/phpmyadmin.conf", conf, 0644)

	if tools.IsRHEL() {
		tools.ExecShell("firewall-cmd --zone=public --add-port=" + port + "/tcp --permanent")
		tools.ExecShell("firewall-cmd --reload")
	} else {
		tools.ExecShell("ufw allow " + port + "/tcp")
		tools.ExecShell("ufw reload")
	}
	tools.ExecShell("systemctl reload openresty")

	controllers.Success(ctx, nil)
}
