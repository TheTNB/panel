package openresty

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/route"

	"github.com/TheTNB/panel/v2/app/http/middleware"
	"github.com/TheTNB/panel/v2/app/plugins/loader"
	"github.com/TheTNB/panel/v2/pkg/types"
)

func init() {
	loader.Register(&types.Plugin{
		Name:        "Fail2ban",
		Description: "Fail2ban 扫描系统日志文件并从中找出多次尝试失败的IP地址，将该IP地址加入防火墙的拒绝访问列表中",
		Slug:        "fail2ban",
		Version:     "1.0.2",
		Requires:    []string{},
		Excludes:    []string{},
		Install:     `bash /www/panel/scripts/fail2ban/install.sh`,
		Uninstall:   `bash /www/panel/scripts/fail2ban/uninstall.sh`,
		Update:      `bash /www/panel/scripts/fail2ban/update.sh`,
		Boot: func(app foundation.Application) {
			RouteFacade := app.MakeRoute()
			RouteFacade.Prefix("api/plugins/fail2ban").Middleware(middleware.Session(), middleware.MustInstall()).Group(func(r route.Router) {
				r.Prefix("openresty").Group(func(route route.Router) {
					controller := NewController()
					route.Get("jails", controller.List)
					route.Post("jails", controller.Add)
					route.Delete("jails", controller.Delete)
					route.Get("jails/{name}", controller.BanList)
					route.Post("unban", controller.Unban)
					route.Post("whiteList", controller.SetWhiteList)
					route.Get("whiteList", controller.GetWhiteList)
				})
			})
		},
	})
}
