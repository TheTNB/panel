package openresty

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"panel/app/http/middleware"
	"panel/plugins/openresty/http/controllers"
)

func Route() {
	facades.Route().Prefix("api/plugins/openresty").Middleware(middleware.Jwt()).Group(func(route route.Route) {
		openRestyController := controllers.NewOpenrestyController()
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
