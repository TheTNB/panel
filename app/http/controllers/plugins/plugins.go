package plugins

import (
	"sync"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/services"
)

// Check 检查插件是否可用
func Check(ctx http.Context, slug string) {
	plugin := services.NewPluginImpl().GetBySlug(slug)
	installedPlugin := services.NewPluginImpl().GetInstalledBySlug(slug)
	installedPlugins, err := services.NewPluginImpl().AllInstalled()
	if err != nil {
		facades.Log().Error("[面板][插件] 获取已安装插件失败")
		ctx.Request().AbortWithStatusJson(http.StatusInternalServerError, "系统内部错误")
	}

	if installedPlugin.Version != plugin.Version || installedPlugin.Slug != plugin.Slug {
		ctx.Request().AbortWithStatusJson(http.StatusForbidden, "插件 "+slug+" 需要更新至 "+plugin.Version+" 版本")
	}

	var lock sync.RWMutex
	pluginsMap := make(map[string]bool)

	for _, p := range installedPlugins {
		lock.Lock()
		pluginsMap[p.Slug] = true
		lock.Unlock()
	}

	for _, require := range plugin.Requires {
		lock.RLock()
		_, requireFound := pluginsMap[require]
		lock.RUnlock()
		if !requireFound {
			ctx.Request().AbortWithStatusJson(http.StatusForbidden, "插件 "+slug+" 需要依赖 "+require+" 插件")
		}
	}

	for _, exclude := range plugin.Excludes {
		lock.RLock()
		_, excludeFound := pluginsMap[exclude]
		lock.RUnlock()
		if excludeFound {
			ctx.Request().AbortWithStatusJson(http.StatusForbidden, "插件 "+slug+" 不兼容 "+exclude+" 插件")
		}
	}
}
