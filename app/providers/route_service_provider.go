package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"

	"panel/app/http"
	"panel/routes"
)

type RouteServiceProvider struct {
}

func (receiver *RouteServiceProvider) Register(app foundation.Application) {
	//Add HTTP middlewares
	facades.Route().GlobalMiddleware(http.Kernel{}.Middleware()...)
}

func (receiver *RouteServiceProvider) Boot(app foundation.Application) {
	receiver.configureRateLimiting()

	routes.Web()
	routes.Plugin()
}

func (receiver *RouteServiceProvider) configureRateLimiting() {

}
