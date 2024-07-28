package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/http/limit"

	"github.com/TheTNB/panel/v2/app/http"
	"github.com/TheTNB/panel/v2/routes"
)

type RouteServiceProvider struct{}

func (receiver *RouteServiceProvider) Register(app foundation.Application) {
}

func (receiver *RouteServiceProvider) Boot(app foundation.Application) {
	// Add HTTP middlewares
	facades.Route().GlobalMiddleware(http.Kernel{}.Middleware()...)

	receiver.configureRateLimiting()

	routes.Plugin()
	routes.Api()
}

func (receiver *RouteServiceProvider) configureRateLimiting() {
	facades.RateLimiter().ForWithLimits("login", func(ctx contractshttp.Context) []contractshttp.Limit {
		return []contractshttp.Limit{
			limit.PerMinute(5).By(ctx.Request().Ip()).Response(func(ctx contractshttp.Context) {
				ctx.Request().AbortWithStatusJson(contractshttp.StatusTooManyRequests, contractshttp.Json{
					"message": "请求过于频繁，请等待一分钟后再试",
				})
			}),
			limit.PerHour(100).By(ctx.Request().Ip()).Response(func(ctx contractshttp.Context) {
				ctx.Request().AbortWithStatusJson(contractshttp.StatusTooManyRequests, contractshttp.Json{
					"message": "请求过于频繁，请等待一小时后再试",
				})
			}),
			limit.PerDay(1000).Response(func(ctx contractshttp.Context) {
				ctx.Request().AbortWithStatusJson(contractshttp.StatusTooManyRequests, contractshttp.Json{
					"message": "面板遭受登录爆破攻击过多，已暂时屏蔽登录，请立刻更换面板端口",
				})
			}),
		}
	})
}
