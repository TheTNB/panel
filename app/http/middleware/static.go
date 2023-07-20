package middleware

import (
	"github.com/gin-contrib/static"

	contractshttp "github.com/goravel/framework/contracts/http"
	frameworkhttp "github.com/goravel/framework/http"

	"panel/app/services"
)

func Static() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		static.Serve(services.NewSettingImpl().Get("entrance", "/"), static.LocalFile("/www/panel/public", false))(ctx.(*frameworkhttp.GinContext).Instance())

		ctx.Request().Next()
	}
}
