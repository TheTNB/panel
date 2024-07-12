package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
)

func Entrance() http.Middleware {
	return func(ctx http.Context) {
		translate := facades.Lang(ctx)

		if !ctx.Request().HasSession() {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
				"message": translate.Get("auth.session.missing"),
			})
			return
		}

		entrance := facades.Config().GetString("panel.entrance")
		if ctx.Request().Path() == entrance {
			ctx.Request().Session().Put("verify_entrance", true)
			_ = ctx.Response().Redirect(http.StatusFound, "/login").Render()
			ctx.Request().AbortWithStatus(http.StatusFound)
			return
		}

		if !facades.Config().GetBool("app.debug") &&
			(ctx.Request().Session().Missing("verify_entrance") || !cast.ToBool(ctx.Request().Session().Get("verify_entrance"))) &&
			ctx.Request().Path() != "/robots.txt" {
			ctx.Request().AbortWithStatusJson(http.StatusTeapot, http.Json{
				"message": "请通过正确的入口访问",
			})
			return
		}

		ctx.Request().Next()
	}
}
