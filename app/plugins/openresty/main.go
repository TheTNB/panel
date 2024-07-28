package openresty

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/route"

	"github.com/TheTNB/panel/v2/app/http/middleware"
	"github.com/TheTNB/panel/v2/app/plugins/loader"
	"github.com/TheTNB/panel/v2/pkg/types"
)

func init() {
	loader.New(&types.Plugin{
		Name:        "OpenResty",
		Description: "OpenResty® 是一款基于 NGINX 和 LuaJIT 的 Web 平台",
		Slug:        "openresty",
		Version:     "1.25.3.1",
		Requires:    []string{},
		Excludes:    []string{},
		Install:     "bash /www/panel/scripts/openresty/install.sh",
		Uninstall:   "bash /www/panel/scripts/openresty/uninstall.sh",
		Update:      "bash /www/panel/scripts/openresty/install.sh",
		Boot: func(app foundation.Application) {
			RouteFacade := app.MakeRoute()
			RouteFacade.Prefix("api/plugins/openresty").Middleware(middleware.Session(), middleware.MustInstall()).Group(func(r route.Router) {
				r.Prefix("openresty").Group(func(route route.Router) {
					controller := NewController()
					route.Get("load", controller.Load)
					route.Get("config", controller.GetConfig)
					route.Post("config", controller.SaveConfig)
					route.Get("errorLog", controller.ErrorLog)
					route.Post("clearErrorLog", controller.ClearErrorLog)
				})
			})
		},
	})
}
