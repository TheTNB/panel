package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"panel/app/http/controllers"
	"panel/app/http/middleware"
)

func Web() {
	facades.Route().Prefix("api/panel").Group(func(r route.Route) {
		r.Prefix("info").Group(func(r route.Route) {
			infoController := controllers.NewInfoController()
			r.Get("name", infoController.Name)
			r.Middleware(middleware.Jwt()).Get("menu", infoController.Menu)
			r.Middleware(middleware.Jwt()).Get("homePlugins", infoController.HomePlugins)
			r.Middleware(middleware.Jwt()).Get("nowMonitor", infoController.NowMonitor)
			r.Middleware(middleware.Jwt()).Get("systemInfo", infoController.SystemInfo)
		})
		r.Prefix("user").Group(func(r route.Route) {
			userController := controllers.NewUserController()
			r.Post("login", userController.Login)
			r.Middleware(middleware.Jwt()).Get("info", userController.Info)
		})
		r.Prefix("task").Middleware(middleware.Jwt()).Group(func(r route.Route) {
			taskController := controllers.NewTaskController()
			r.Get("status", taskController.Status)
		})
	})

	facades.Route().Fallback(func(ctx http.Context) {
		ctx.Response().String(404, "not found")
	})
}
