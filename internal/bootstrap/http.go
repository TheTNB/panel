package bootstrap

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/internal/apps"
	"github.com/TheTNB/panel/internal/http/middleware"
	"github.com/TheTNB/panel/internal/panel"
	"github.com/TheTNB/panel/internal/route"
)

func initHttp() {
	panel.Http = chi.NewRouter()

	// add middleware
	panel.Http.Use(middleware.GlobalMiddleware()...)

	// add route
	route.Http(panel.Http)
	apps.Boot(panel.Http)

	server := &http.Server{
		Addr:           panel.Conf.MustString("http.address"),
		Handler:        http.AllowQuerySemicolons(panel.Http),
		MaxHeaderBytes: 2048 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("failed to start http server: %v", err))
	}
}
