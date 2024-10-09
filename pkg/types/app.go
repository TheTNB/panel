package types

import "github.com/go-chi/chi/v5"

// App 应用元数据结构
type App struct {
	Slug  string             `json:"slug"` // 应用标识
	Route func(r chi.Router) `json:"-"`    // 路由
}

// StoreApp 商店应用结构
type StoreApp struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	Versions    []struct {
		Slug      string `json:"slug"`
		Name      string `json:"name"`
		Panel     string `json:"panel"`
		Install   string `json:"install"`
		Uninstall string `json:"uninstall"`
		Update    string `json:"update"`
		Subs      []struct {
			Log     string `json:"log"`
			Version string `json:"version"`
		} `json:"versions"`
	} `json:"versions"`
	Installed            bool   `json:"installed"`
	InstalledVersion     string `json:"installed_version"`
	InstalledVersionSlug string `json:"installed_version_slug"`
	UpdateExist          bool   `json:"update_exist"`
	Show                 bool   `json:"show"`
}
