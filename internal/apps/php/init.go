package php

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "php80",
		Route: func(r chi.Router) {
			service := NewService(80)
			r.Get("/load", service.Load)
			r.Get("/config", service.GetConfig)
			r.Post("/config", service.UpdateConfig)
			r.Get("/fpmConfig", service.GetFPMConfig)
			r.Post("/fpmConfig", service.UpdateFPMConfig)
			r.Get("/errorLog", service.ErrorLog)
			r.Get("/slowLog", service.SlowLog)
			r.Post("/clearErrorLog", service.ClearErrorLog)
			r.Post("/clearSlowLog", service.ClearSlowLog)
			r.Get("/extensions", service.ExtensionList)
			r.Post("/extensions", service.InstallExtension)
			r.Delete("/extensions", service.UninstallExtension)
		},
	})
	apploader.Register(&types.App{
		Slug: "php81",
		Route: func(r chi.Router) {
			service := NewService(81)
			r.Get("/load", service.Load)
			r.Get("/config", service.GetConfig)
			r.Post("/config", service.UpdateConfig)
			r.Get("/fpmConfig", service.GetFPMConfig)
			r.Post("/fpmConfig", service.UpdateFPMConfig)
			r.Get("/errorLog", service.ErrorLog)
			r.Get("/slowLog", service.SlowLog)
			r.Post("/clearErrorLog", service.ClearErrorLog)
			r.Post("/clearSlowLog", service.ClearSlowLog)
			r.Get("/extensions", service.ExtensionList)
			r.Post("/extensions", service.InstallExtension)
			r.Delete("/extensions", service.UninstallExtension)
		},
	})
	apploader.Register(&types.App{
		Slug: "php82",
		Route: func(r chi.Router) {
			service := NewService(82)
			r.Get("/load", service.Load)
			r.Get("/config", service.GetConfig)
			r.Post("/config", service.UpdateConfig)
			r.Get("/fpmConfig", service.GetFPMConfig)
			r.Post("/fpmConfig", service.UpdateFPMConfig)
			r.Get("/errorLog", service.ErrorLog)
			r.Get("/slowLog", service.SlowLog)
			r.Post("/clearErrorLog", service.ClearErrorLog)
			r.Post("/clearSlowLog", service.ClearSlowLog)
			r.Get("/extensions", service.ExtensionList)
			r.Post("/extensions", service.InstallExtension)
			r.Delete("/extensions", service.UninstallExtension)
		},
	})
	apploader.Register(&types.App{
		Slug: "php83",
		Route: func(r chi.Router) {
			service := NewService(83)
			r.Get("/load", service.Load)
			r.Get("/config", service.GetConfig)
			r.Post("/config", service.UpdateConfig)
			r.Get("/fpmConfig", service.GetFPMConfig)
			r.Post("/fpmConfig", service.UpdateFPMConfig)
			r.Get("/errorLog", service.ErrorLog)
			r.Get("/slowLog", service.SlowLog)
			r.Post("/clearErrorLog", service.ClearErrorLog)
			r.Post("/clearSlowLog", service.ClearSlowLog)
			r.Get("/extensions", service.ExtensionList)
			r.Post("/extensions", service.InstallExtension)
			r.Delete("/extensions", service.UninstallExtension)
		},
	})
}
