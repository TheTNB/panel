package middleware

import (
	"fmt"

	"github.com/gin-contrib/static"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/gin"

	"panel/app/services"
)

func Static() http.Middleware {
	return func(ctx http.Context) {
		// 自动纠正 URL 格式
		if ctx.Request().Path() == services.NewSettingImpl().Get("entrance", "/") && ctx.Request().Path() != "/" {
			// ctx.Response().Redirect(http.StatusFound, ctx.Request().Path()+"/")
			ctx.Response().Writer().WriteHeader(http.StatusFound)
			_, err := ctx.Response().Writer().Write([]byte(fmt.Sprintf(`<html><head><meta http-equiv="refresh" content="0;url=%s/"></head></html>`, ctx.Request().Path())))
			if err != nil {
				return
			}
			ctx.Response().Flush()
			return
		}

		static.Serve(services.NewSettingImpl().Get("entrance", "/"), static.LocalFile("public", false))(ctx.(*gin.Context).Instance())
		ctx.Request().Next()
	}
}
