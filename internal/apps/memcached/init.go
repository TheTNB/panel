package memcached

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "memcached",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/load", service.Load)
			r.Get("/config", service.GetConfig)
			r.Post("/config", service.UpdateConfig)
		},
	})
}
