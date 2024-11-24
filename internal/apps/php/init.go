package php

import (
	"fmt"

	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	php := []uint{74, 80, 81, 82, 83, 84}
	for _, version := range php {
		apploader.Register(&types.App{
			Slug: fmt.Sprintf("php%d", version),
			Route: func(r chi.Router) {
				service := NewService(version)
				r.Post("/setCli", service.SetCli)
				r.Get("/config", service.GetConfig)
				r.Post("/config", service.UpdateConfig)
				r.Get("/fpmConfig", service.GetFPMConfig)
				r.Post("/fpmConfig", service.UpdateFPMConfig)
				r.Get("/load", service.Load)
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
}
