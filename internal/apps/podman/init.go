package podman

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "podman",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/registryConfig", service.GetRegistryConfig)
			r.Post("/registryConfig", service.UpdateRegistryConfig)
			r.Get("/storageConfig", service.GetStorageConfig)
			r.Post("/storageConfig", service.UpdateStorageConfig)
		},
	})
}
