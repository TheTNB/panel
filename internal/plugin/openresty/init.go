package openresty

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/pluginloader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	pluginloader.Register(&types.Plugin{
		Slug: "openresty",
		Name: "OpenResty",
		Route: func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("Hello, World!"))
			})
		},
	})
}
