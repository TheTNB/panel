package postgresql

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "postgresql",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/config", service.GetConfig)
			r.Post("/config", service.UpdateConfig)
			r.Get("/userConfig", service.GetUserConfig)
			r.Post("/userConfig", service.UpdateUserConfig)
			r.Get("/load", service.Load)
			r.Get("/log", service.Log)
			r.Post("/clearLog", service.ClearLog)
		},
	})
}
