package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"panel/app/http/controllers/plugins/openresty"
	"panel/app/http/middleware"
)

// Plugin 加载插件路由
func Plugin() {
	facades.Route().Prefix("api/plugins/openresty").Middleware(middleware.Jwt()).Group(func(route route.Route) {
		openRestyController := openresty.NewOpenrestyController()
		route.Get("status", openRestyController.Status)
		route.Post("reload", openRestyController.Reload)
		route.Post("start", openRestyController.Start)
		route.Post("stop", openRestyController.Stop)
		route.Post("restart", openRestyController.Restart)
		route.Get("load", openRestyController.Load)
		route.Get("config", openRestyController.GetConfig)
		route.Post("config", openRestyController.SaveConfig)
		route.Get("errorLog", openRestyController.ErrorLog)
		route.Get("cleanErrorLog", openRestyController.ClearErrorLog)
	})
}
