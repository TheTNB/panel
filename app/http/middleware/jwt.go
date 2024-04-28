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
		translate := facades.Lang(ctx)
		token := ctx.Request().Header("Authorization", ctx.Request().Header("Sec-WebSocket-Protocol"))
		if len(token) == 0 {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
				"message": translate.Get("auth.token.missing"),
			})
			return
		}

		// JWT 鉴权
		if _, err := facades.Auth(ctx).Parse(token); err != nil {
			if errors.Is(err, auth.ErrorTokenExpired) {
				token, err = facades.Auth(ctx).Refresh()
				if err != nil {
					// 到达刷新时间上限
					ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
						"message": translate.Get("auth.token.expired"),
					})
					return
				}

				token = "Bearer " + token
			} else {
				ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
					"message": translate.Get("auth.token.expired"),
				})
				return
			}
		}

		ctx.Response().Header("Authorization", token)
		ctx.Request().Next()
	}
}
