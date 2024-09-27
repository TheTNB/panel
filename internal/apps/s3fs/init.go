package s3fs

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "s3fs",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/mounts", service.List)
			r.Post("/mounts", service.Create)
			r.Delete("/mounts", service.Delete)
		},
	})
}
