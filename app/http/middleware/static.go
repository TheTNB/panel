package middleware

import (
	"github.com/gin-contrib/static"

	contractshttp "github.com/goravel/framework/contracts/http"
	frameworkhttp "github.com/goravel/framework/http"
)

func Static() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		static.Serve("/", static.LocalFile("./public", false))(ctx.(*frameworkhttp.GinContext).Instance())

		ctx.Request().Next()
	}
}
