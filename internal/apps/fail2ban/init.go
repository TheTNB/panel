package fail2ban

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "fail2ban",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/jails", service.List)
			r.Post("/jails", service.Create)
			r.Delete("/jails", service.Delete)
			r.Get("/jails/{name}", service.BanList)
			r.Post("/unban", service.Unban)
			r.Post("/whiteList", service.SetWhiteList)
			r.Get("/whiteList", service.GetWhiteList)
		},
	})
}
