package toolbox

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "toolbox",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/dns", service.GetDNS)
			r.Post("/dns", service.UpdateDNS)
			r.Get("/swap", service.GetSWAP)
			r.Post("/swap", service.UpdateSWAP)
			r.Get("/timezone", service.GetTimezone)
			r.Post("/timezone", service.UpdateTimezone)
			r.Post("/time", service.UpdateTime)
			r.Post("/syncTime", service.SyncTime)
			r.Get("/hostname", service.GetHostname)
			r.Post("/hostname", service.UpdateHostname)
			r.Get("/hosts", service.GetHosts)
			r.Post("/hosts", service.UpdateHosts)
			r.Post("/rootPassword", service.UpdateRootPassword)
		},
	})
}
