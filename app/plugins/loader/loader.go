package loader

import (
	"github.com/TheTNB/panel/v2/pkg/types"
)

var data []*types.Plugin

// All 获取所有插件
func All() []*types.Plugin {
	return data
}

// New 新注册插件
func New(plugin *types.Plugin) {
	data = append(data, plugin)
}
