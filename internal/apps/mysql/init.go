package mysql

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "mysql",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/load", service.Load)
			r.Get("/config", service.GetConfig)
			r.Post("/config", service.UpdateConfig)
			r.Post("/clearErrorLog", service.ClearErrorLog)
			r.Get("/slowLog", service.SlowLog)
			r.Post("/clearSlowLog", service.ClearSlowLog)
			r.Get("/rootPassword", service.GetRootPassword)
			r.Post("/rootPassword", service.SetRootPassword)
		},
	})
}
