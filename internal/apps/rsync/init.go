package rsync

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "rsync",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/modules", service.List)
			r.Post("/modules", service.Create)
			r.Post("/modules/{name}", service.Update)
			r.Delete("/modules/{name}", service.Delete)
			r.Get("/config", service.GetConfig)
			r.Post("/config", service.UpdateConfig)
		},
	})
}
