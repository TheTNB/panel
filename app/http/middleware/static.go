package middleware

import (
	"github.com/gin-contrib/static"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/gin"

	"github.com/TheTNB/panel/v2/embed"
)

func Static() http.Middleware {
	return func(ctx http.Context) {
		static.Serve("/", static.EmbedFolder(embed.PublicFS, "frontend"))(ctx.(*gin.Context).Instance())
		ctx.Request().Next()
	}
}
