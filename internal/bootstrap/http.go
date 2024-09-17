package bootstrap

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/http/middleware"
	"github.com/TheTNB/panel/internal/plugin"
	"github.com/TheTNB/panel/internal/route"
)

func initHttp() {
	app.Http = chi.NewRouter()

	// add middleware
	app.Http.Use(middleware.GlobalMiddleware()...)

	// add route
	route.Http(app.Http)
	plugin.Boot(app.Http)

	server := &http.Server{
		Addr:           app.Conf.MustString("http.address"),
		Handler:        http.AllowQuerySemicolons(app.Http),
		MaxHeaderBytes: 2048 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("failed to start http server: %v", err))
	}
}
