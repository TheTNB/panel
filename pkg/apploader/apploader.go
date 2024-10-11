// Package apploader 面板应用加载器
package apploader

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/types"
)

var apps sync.Map

func Register(app *types.App) {
	if _, ok := apps.Load(app.Slug); ok {
		log.Fatalf("app %s already exists", app.Slug)
	}
	apps.Store(app.Slug, app)
}

func Get(slug string) (*types.App, error) {
	if app, ok := apps.Load(slug); ok {
		return app.(*types.App), nil
	}
	return nil, fmt.Errorf("app %s not found", slug)
}

func All() []*types.App {
	var list []*types.App
	apps.Range(func(_, app any) bool {
		if p, ok := app.(*types.App); ok {
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
	apps.Range(func(_, app any) bool {
		if p, ok := app.(*types.App); ok {
			r.Route(fmt.Sprintf("/api/apps/%s", p.Slug), p.Route)
		}
		return true
	})
}
