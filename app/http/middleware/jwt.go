package middleware

import (
	"errors"

	"github.com/goravel/framework/auth"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/models"
)

// Jwt 确保通过 JWT 鉴权
func Jwt() http.Middleware {
	return func(ctx http.Context) {
		token := ctx.Request().Header("Authorization", "")
		if len(token) == 0 {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
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
					// Refresh time exceeded
					ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
						"code":    401,
						"message": "登录已过期",
					})
					return
				}

				token = "Bearer " + token
			} else {
				ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
					"code":    401,
					"message": "登录已过期",
				})
				return
			}
		}

		// 取出用户信息
		var user models.User
		if err := facades.Auth().User(ctx, &user); err != nil {
			ctx.Request().AbortWithStatusJson(http.StatusForbidden, http.Json{
				"code":    403,
				"message": "用户不存在",
			})
			return
		}

		ctx.WithValue("user", user)

		ctx.Response().Header("Authorization", token)
		ctx.Request().Next()
	}
}
