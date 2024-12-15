package middleware

import (
	"github.com/TheTNB/panel/internal/biz"
	"github.com/go-chi/chi/v5"
	"github.com/go-rat/sessions"
	"github.com/google/wire"
	"github.com/knadh/koanf/v2"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	sessionmiddleware "github.com/go-rat/sessions/middleware"
	"github.com/golang-cz/httplog"
)

var ProviderSet = wire.NewSet(NewMiddlewares)

type Middlewares struct {
	conf    *koanf.Koanf
	log     *slog.Logger
	session *sessions.Manager
	app     biz.AppRepo
}

func NewMiddlewares(conf *koanf.Koanf, log *slog.Logger, session *sessions.Manager, app biz.AppRepo) *Middlewares {
	return &Middlewares{
		conf:    conf,
		log:     log,
		session: session,
		app:     app,
	}
}

// Globals is a collection of global middleware that will be applied to every request.
func (r *Middlewares) Globals(mux *chi.Mux) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		sessionmiddleware.StartSession(r.session),
		//middleware.SupressNotFound(mux),// bug https://github.com/go-chi/chi/pull/940
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.Compress(5),
		httplog.RequestLogger(r.log, &httplog.Options{
			Level:             slog.LevelInfo,
			LogRequestHeaders: []string{"User-Agent"},
		}),
		middleware.Recoverer,
		Status,
		Entrance(r.conf, r.session),
		MustLogin(r.session),
		MustInstall(r.app),
	}
}
