package controllers

import (
	"sync"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/services"
)

func Success(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

func Error(ctx http.Context, code int, message any) http.Response {
	return ctx.Response().Json(http.StatusOK, http.Json{
		"code":    code,
		"message": message,
	})
}

func SystemError(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusOK, http.Json{
		"code":    http.StatusInternalServerError,
		"message": "系统内部错误",
	})
}

// Check 检查插件是否可用
func Check(ctx http.Context, slug string) http.Response {
	plugin := services.NewPluginImpl().GetBySlug(slug)
	installedPlugin := services.NewPluginImpl().GetInstalledBySlug(slug)
	installedPlugins, err := services.NewPluginImpl().AllInstalled()
	if err != nil {
		facades.Log().Error("[面板][插件] 获取已安装插件失败")
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if installedPlugin.Version != plugin.Version || installedPlugin.Slug != plugin.Slug {
		return Error(ctx, http.StatusForbidden, "插件 "+slug+" 需要更新至 "+plugin.Version+" 版本")
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
			return Error(ctx, http.StatusForbidden, "插件 "+slug+" 需要依赖 "+require+" 插件")
		}
	}

	for _, exclude := range plugin.Excludes {
		lock.RLock()
		_, excludeFound := pluginsMap[exclude]
		lock.RUnlock()
		if excludeFound {
			return Error(ctx, http.StatusForbidden, "插件 "+slug+" 不兼容 "+exclude+" 插件")
		}
	}

	return nil
}
