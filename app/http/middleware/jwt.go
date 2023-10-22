package middleware

import (
	"errors"

	"github.com/goravel/framework/auth"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// Jwt 确保通过 JWT 鉴权
func Jwt() http.Middleware {
	return func(ctx http.Context) {
		token := ctx.Request().Header("Authorization", "")
		if len(token) == 0 {
			ctx.Request().AbortWithStatusJson(http.StatusOK, http.Json{
				"code":    401,
				"message": "未登录",
			})
			return
		}

		// JWT 鉴权
		if _, err := facades.Auth().Parse(ctx, token); err != nil {
			if errors.Is(err, auth.ErrorTokenExpired) {
				token, err = facades.Auth().Refresh(ctx)
				if err != nil {
					// 到达刷新时间上限
					ctx.Request().AbortWithStatusJson(http.StatusOK, http.Json{
						"code":    401,
						"message": "登录已过期",
					})
					return
				}

				token = "Bearer " + token
			} else {
				ctx.Request().AbortWithStatusJson(http.StatusOK, http.Json{
					"code":    401,
					"message": "登录已过期",
				})
				return
			}
		}

		ctx.Response().Header("Authorization", token)
		ctx.Request().Next()
	}
}
