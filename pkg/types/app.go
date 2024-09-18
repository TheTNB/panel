package types

import "github.com/go-chi/chi/v5"

// App 应用元数据结构
type App struct {
	Slug  string             `json:"slug"` // 插件标识
	Route func(r chi.Router) `json:"-"`    // 路由
}
