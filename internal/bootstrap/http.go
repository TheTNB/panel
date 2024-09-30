package bootstrap

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/apps"
	"github.com/TheTNB/panel/internal/http/middleware"
	"github.com/TheTNB/panel/internal/route"
)

func initHttp() {
	app.Http = chi.NewRouter()

	// add middleware
	app.Http.Use(middleware.GlobalMiddleware()...)

	// add route
	route.Http(app.Http)
	apps.Boot(app.Http)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", app.Conf.MustInt("http.port")),
		Handler:        http.AllowQuerySemicolons(app.Http),
		MaxHeaderBytes: 2048 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("failed to start http server: %v", err))
	}
}
