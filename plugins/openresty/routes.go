package openresty

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"panel/plugins/openresty/http/controllers"
)

func Route() {
	facades.Route().Prefix("api/plugins/openresty").Group(func(route route.Route) {
		route.Get("/openresty", func(ctx http.Context) {
			ctx.Response().Json(http.StatusOK, http.Json{
				"Hello": "Openresty",
			})
		})

		openRestyController := controllers.NewOpenrestyController()
		route.Get("/openresty/users/{id}", openRestyController.Show)
	})
}
