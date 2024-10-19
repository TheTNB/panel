package benchmark

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "benchmark",
		Route: func(r chi.Router) {
			service := NewService()
			r.Post("/test", service.Test)
		},
	})
}
