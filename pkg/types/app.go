package types

import "github.com/go-chi/chi/v5"

// App 应用元数据结构
type App struct {
	Slug  string             `json:"slug"` // 应用标识
	Route func(r chi.Router) `json:"-"`    // 路由
}

// StoreApp 商店应用结构
type StoreApp struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	Slug             string `json:"slug"`
	Version          string `json:"version"`
	Installed        bool   `json:"installed"`
	InstalledVersion string `json:"installed_version"`
	Show             bool   `json:"show"`
}
