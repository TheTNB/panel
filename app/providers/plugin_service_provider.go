package providers

import (
	"github.com/goravel/framework/contracts/foundation"

	"github.com/TheTNB/panel/v2/app/plugins"
	"github.com/TheTNB/panel/v2/app/plugins/loader"
)

type PluginServiceProvider struct{}

func (receiver *PluginServiceProvider) Register(app foundation.Application) {
	plugins.Boot()
}

func (receiver *PluginServiceProvider) Boot(app foundation.Application) {
	for _, plugin := range loader.All() {
		plugin.Boot(app)
	}
}
