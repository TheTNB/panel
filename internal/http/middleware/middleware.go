package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-rat/sessions"
	"github.com/knadh/koanf/v2"
	"gorm.io/gorm"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	sessionmiddleware "github.com/go-rat/sessions/middleware"
	"github.com/golang-cz/httplog"
)

// GlobalMiddleware is a collection of global middleware that will be applied to every request.
func GlobalMiddleware(r *chi.Mux, conf *koanf.Koanf, db *gorm.DB, log *slog.Logger, session *sessions.Manager) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		sessionmiddleware.StartSession(session),
		//middleware.SupressNotFound(r),// bug https://github.com/go-chi/chi/pull/940
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.Compress(5),
		httplog.RequestLogger(log, &httplog.Options{
			Level:             slog.LevelInfo,
			LogRequestHeaders: []string{"User-Agent"},
		}),
		middleware.Recoverer,
		Status,
		Entrance(conf, session),
		MustLogin(session),
		MustInstall(db),
	}
}
