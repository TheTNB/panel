package middleware

import (
	"github.com/gin-contrib/static"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/gin"

	"panel/app/services"
)

func Static() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		static.Serve(services.NewSettingImpl().Get("entrance", "/"), static.LocalFile("/www/panel/public", false))(ctx.(*gin.Context).Instance())

		ctx.Request().Next()
	}
}
