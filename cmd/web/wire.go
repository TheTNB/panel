//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/tnb-labs/panel/internal/app"
	"github.com/tnb-labs/panel/internal/apps"
	"github.com/tnb-labs/panel/internal/bootstrap"
	"github.com/tnb-labs/panel/internal/data"
	"github.com/tnb-labs/panel/internal/http/middleware"
	"github.com/tnb-labs/panel/internal/job"
	"github.com/tnb-labs/panel/internal/route"
	"github.com/tnb-labs/panel/internal/service"
)

// initWeb init application.
func initWeb() (*app.Web, error) {
	panic(wire.Build(bootstrap.ProviderSet, middleware.ProviderSet, route.ProviderSet, service.ProviderSet, data.ProviderSet, apps.ProviderSet, job.ProviderSet, app.NewWeb))
}
