package middleware

import (
	"github.com/goravel/framework/contracts/http"

	"panel/internal"
)

// Status 检查程序状态
func Status() http.Middleware {
	return func(ctx http.Context) {
		switch internal.Status {
		case internal.StatusUpgrade:
			ctx.Request().AbortWithStatusJson(http.StatusOK, http.Json{
				"code":    503,
				"message": "面板升级中，请稍后",
			})
			return
		case internal.StatusMaintain:
			ctx.Request().AbortWithStatusJson(http.StatusOK, http.Json{
				"code":    503,
				"message": "面板正在运行维护，请稍后",
			})
			return
		case internal.StatusClosed:
			ctx.Request().AbortWithStatusJson(http.StatusOK, http.Json{
				"code":    403,
				"message": "面板已关闭",
			})
			return
		case internal.StatusFailed:
			ctx.Request().AbortWithStatusJson(http.StatusOK, http.Json{
				"code":    500,
				"message": "面板运行出错，请检查排除或联系支持",
			})
			return
		default:
			ctx.Request().Next()
			return
		}
	}
}
