package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// Log 记录请求日志
func Log() http.Middleware {
	return func(ctx http.Context) {
		facades.Log().Channel("http").With(map[string]any{
			"Method": ctx.Request().Method(),
			"URL":    ctx.Request().FullUrl(),
			"IP":     ctx.Request().Ip(),
			"UA":     ctx.Request().Header("User-Agent"),
			"Body":   ctx.Request().All(),
		}).Info("HTTP Request")
		ctx.Request().Next()
	}
}
