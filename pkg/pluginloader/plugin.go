// Package pluginloader 面板插件加载器
package pluginloader

import (
	"cmp"
	"fmt"
	"slices"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/types"
)

var plugins sync.Map

func Register(plugin *types.Plugin) {
	plugins.Store(plugin.Slug, plugin)
}

func Get(slug string) (*types.Plugin, error) {
	if plugin, ok := plugins.Load(slug); ok {
		return plugin.(*types.Plugin), nil
	}
	return nil, fmt.Errorf("plugin %s not found", slug)
}

func All() []*types.Plugin {
	var list []*types.Plugin
	plugins.Range(func(_, plugin any) bool {
		if p, ok := plugin.(*types.Plugin); ok {
			list = append(list, p)
		}
		return true
	})

	// 排序
	slices.SortFunc(list, func(a, b *types.Plugin) int {
		return cmp.Compare(a.Order, b.Order)
	})

	return list
}

func Boot(r chi.Router) {
	plugins.Range(func(_, plugin any) bool {
		if p, ok := plugin.(*types.Plugin); ok {
			r.Route(fmt.Sprintf("/api/plugins/%s", p.Slug), p.Route)
		}
		return true
	})
}
