package middleware

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	sessionmiddleware "github.com/go-rat/sessions/middleware"
	"github.com/golang-cz/httplog"

	"github.com/TheTNB/panel/internal/app"
)

// GlobalMiddleware is a collection of global middleware that will be applied to every request.
func GlobalMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		sessionmiddleware.StartSession(app.Session),
		//middleware.SupressNotFound(app.Http),// bug https://github.com/go-chi/chi/pull/940
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.Compress(5),
		httplog.RequestLogger(app.Logger, &httplog.Options{
			Level:             slog.LevelInfo,
			LogRequestHeaders: []string{"User-Agent"},
		}),
		middleware.Recoverer,
		Entrance,
		Status,
		MustInstall,
	}
}
