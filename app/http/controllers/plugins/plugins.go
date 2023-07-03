package plugins

import (
	"sync"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/http/controllers"
	"panel/app/services"
)

// Check 检查插件是否可用
func Check(ctx http.Context, slug string) bool {
	plugin := services.NewPluginImpl().GetBySlug(slug)
	installedPlugin := services.NewPluginImpl().GetInstalledBySlug(slug)
	installedPlugins, err := services.NewPluginImpl().AllInstalled()
	if err != nil {
		facades.Log().Error("[面板][插件] 获取已安装插件失败")
		controllers.Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return false
	}

	if installedPlugin.Version != plugin.Version || installedPlugin.Slug != plugin.Slug {
		controllers.Error(ctx, http.StatusForbidden, "插件 "+slug+" 需要更新至 "+plugin.Version+" 版本")
		return false
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
			controllers.Error(ctx, http.StatusForbidden, "插件 "+slug+" 需要依赖 "+require+" 插件")
			return false
		}
	}

	for _, exclude := range plugin.Excludes {
		lock.RLock()
		_, excludeFound := pluginsMap[exclude]
		lock.RUnlock()
		if excludeFound {
			controllers.Error(ctx, http.StatusForbidden, "插件 "+slug+" 不兼容 "+exclude+" 插件")
			return false
		}
	}

	return true
}
