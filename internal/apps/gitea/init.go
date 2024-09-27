package gitea

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "gitea",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/config", service.GetConfig)
			r.Post("/config", service.UpdateConfig)
		},
	})
}
