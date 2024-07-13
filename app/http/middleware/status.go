package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/pkg/types"
)

// Status 检查程序状态
func Status() http.Middleware {
	return func(ctx http.Context) {
		translate := facades.Lang(ctx)
		switch types.Status {
		case types.StatusUpgrade:
			ctx.Request().AbortWithStatusJson(http.StatusServiceUnavailable, http.Json{
				"message": translate.Get("status.upgrade"),
			})
			return
		case types.StatusMaintain:
			ctx.Request().AbortWithStatusJson(http.StatusServiceUnavailable, http.Json{
				"message": translate.Get("status.maintain"),
			})
			return
		case types.StatusClosed:
			ctx.Request().AbortWithStatusJson(http.StatusForbidden, http.Json{
				"message": translate.Get("status.closed"),
			})
			return
		case types.StatusFailed:
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
