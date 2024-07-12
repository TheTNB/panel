package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
)

// Session 确保通过 JWT 鉴权
func Session() http.Middleware {
	return func(ctx http.Context) {
		translate := facades.Lang(ctx)

		if !ctx.Request().HasSession() {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
				"message": translate.Get("auth.session.missing"),
			})
			return
		}

		if ctx.Request().Session().Missing("user_id") {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
				"message": translate.Get("auth.session.expired"),
			})
			return
		}

		userID := cast.ToUint(ctx.Request().Session().Get("user_id"))
		if userID == 0 {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
				"message": translate.Get("auth.session.invalid"),
			})
			return
		}

		// 刷新会话
		/*if err := ctx.Request().Session().Regenerate(); err == nil {
			ctx.Response().Cookie(http.Cookie{
				Name:     ctx.Request().Session().GetName(),
				Value:    ctx.Request().Session().GetID(),
				MaxAge:   facades.Config().GetInt("session.lifetime") * 60,
				Path:     facades.Config().GetString("session.path"),
				Domain:   facades.Config().GetString("session.domain"),
				Secure:   facades.Config().GetBool("session.secure"),
				HttpOnly: facades.Config().GetBool("session.http_only"),
				SameSite: facades.Config().GetString("session.same_site"),
			})
		}*/

		ctx.WithValue("user_id", userID)
		ctx.Request().Next()
	}
}
