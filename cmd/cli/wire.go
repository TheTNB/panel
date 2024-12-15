//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/bootstrap"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/route"
	"github.com/TheTNB/panel/internal/service"
)

// initCli init command line.
func initCli() (*app.Cli, error) {
	panic(wire.Build(bootstrap.ProviderSet, route.ProviderSet, service.ProviderSet, data.ProviderSet, app.NewCli))
}
