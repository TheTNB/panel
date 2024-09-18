package apps

import (
	"github.com/go-chi/chi/v5"

	_ "github.com/TheTNB/panel/internal/apps/fail2ban"
	_ "github.com/TheTNB/panel/internal/apps/openresty"
	"github.com/TheTNB/panel/pkg/apploader"
)

func Boot(r chi.Router) {
	apploader.Boot(r)
}
