package plugin

import (
	"github.com/go-chi/chi/v5"

	_ "github.com/TheTNB/panel/internal/plugin/fail2ban"
	_ "github.com/TheTNB/panel/internal/plugin/openresty"
	"github.com/TheTNB/panel/pkg/pluginloader"
)

func Boot(r chi.Router) {
	pluginloader.Boot(r)
}
