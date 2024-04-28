package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/internal"
)

// Status 检查程序状态
func Status() http.Middleware {
	return func(ctx http.Context) {
		translate := facades.Lang(ctx)
		switch internal.Status {
		case internal.StatusUpgrade:
			ctx.Request().AbortWithStatusJson(http.StatusServiceUnavailable, http.Json{
				"message": translate.Get("status.upgrade"),
			})
			return
		case internal.StatusMaintain:
			ctx.Request().AbortWithStatusJson(http.StatusServiceUnavailable, http.Json{
				"message": translate.Get("status.maintain"),
			})
			return
		case internal.StatusClosed:
			ctx.Request().AbortWithStatusJson(http.StatusForbidden, http.Json{
				"message": translate.Get("status.closed"),
			})
			return
		case internal.StatusFailed:
			ctx.Request().AbortWithStatusJson(http.StatusInternalServerError, http.Json{
				"message": translate.Get("status.failed"),
			})
			return
		default:
			ctx.Request().Next()
			return
		}
	}
}
