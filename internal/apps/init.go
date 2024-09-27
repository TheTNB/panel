package apps

import (
	"github.com/go-chi/chi/v5"

	_ "github.com/TheTNB/panel/internal/apps/fail2ban"
	_ "github.com/TheTNB/panel/internal/apps/frp"
	_ "github.com/TheTNB/panel/internal/apps/gitea"
	_ "github.com/TheTNB/panel/internal/apps/openresty"
	_ "github.com/TheTNB/panel/internal/apps/percona"
	_ "github.com/TheTNB/panel/internal/apps/php"
	_ "github.com/TheTNB/panel/internal/apps/phpmyadmin"
	"github.com/TheTNB/panel/pkg/apploader"
)

func Boot(r chi.Router) {
	apploader.Boot(r)
}
