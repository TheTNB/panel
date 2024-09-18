package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	sessionmiddleware "github.com/go-rat/sessions/middleware"

	"github.com/TheTNB/panel/internal/panel"
)

// GlobalMiddleware is a collection of global middleware that will be applied to every request.
func GlobalMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		sessionmiddleware.StartSession(panel.Session),
		//middleware.SupressNotFound(app.Http),// bug https://github.com/go-chi/chi/pull/940
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Compress(5),
	}
}
