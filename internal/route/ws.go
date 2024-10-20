package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/internal/http/middleware"
	"github.com/TheTNB/panel/internal/service"
)

func Ws(r chi.Router) {
	r.Route("/api/ws", func(r chi.Router) {
		r.Use(middleware.MustLogin)
		ws := service.NewWsService()
		r.Get("/ssh", ws.Session)
		r.Get("/exec", ws.Exec)
	})
}
