package pureftpd

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "pureftpd",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/users", service.List)
			r.Post("/users", service.Create)
			r.Delete("/users/{name}", service.Delete)
			r.Post("/users/{name}/password", service.ChangePassword)
			r.Get("/port", service.GetPort)
			r.Post("/port", service.UpdatePort)
		},
	})
}
