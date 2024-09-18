package openresty

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.Plugin{
		Order:       -100,
		Slug:        "openresty",
		Name:        "OpenResty",
		Description: "OpenResty® 是一款基于 NGINX 和 LuaJIT 的 Web 平台",
		Version:     "1.25.3.1",
		Requires:    []string{},
		Excludes:    []string{},
		Install:     "bash /www/panel/scripts/openresty/install.sh",
		Uninstall:   "bash /www/panel/scripts/openresty/uninstall.sh",
		Update:      "bash /www/panel/scripts/openresty/install.sh",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/load", service.Load)
			r.Get("/config", service.GetConfig)
			r.Post("/config", service.SaveConfig)
			r.Get("/errorLog", service.ErrorLog)
			r.Post("/clearErrorLog", service.ClearErrorLog)
		},
	})
}
