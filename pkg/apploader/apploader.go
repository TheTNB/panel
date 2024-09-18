// Package apploader 面板应用加载器
package apploader

import (
	"fmt"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/types"
)

var plugins sync.Map

func Register(plugin *types.App) {
	plugins.Store(plugin.Slug, plugin)
}

func Get(slug string) (*types.App, error) {
	if plugin, ok := plugins.Load(slug); ok {
		return plugin.(*types.App), nil
	}
	return nil, fmt.Errorf("plugin %s not found", slug)
}

func All() []*types.App {
	var list []*types.App
	plugins.Range(func(_, plugin any) bool {
		if p, ok := plugin.(*types.App); ok {
			list = append(list, p)
		}
		return true
	})

	// 排序
	/*slices.SortFunc(list, func(a, b *types.App) int {
		return cmp.Compare(a.Order, b.Order)
	})*/

	return list
}

func Boot(r chi.Router) {
	plugins.Range(func(_, plugin any) bool {
		if p, ok := plugin.(*types.App); ok {
			r.Route(fmt.Sprintf("/api/plugins/%s", p.Slug), p.Route)
		}
		return true
	})
}
