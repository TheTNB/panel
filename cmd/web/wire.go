//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/bootstrap"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/middleware"
	"github.com/TheTNB/panel/internal/job"
	"github.com/TheTNB/panel/internal/route"
	"github.com/TheTNB/panel/internal/service"
)

// initWeb init application.
func initWeb() (*app.Web, error) {
	panic(wire.Build(bootstrap.ProviderSet, middleware.ProviderSet, route.ProviderSet, service.ProviderSet, data.ProviderSet, job.ProviderSet, app.NewWeb))
}
