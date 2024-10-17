package types

import "github.com/go-chi/chi/v5"

// App 应用元数据结构
type App struct {
	Slug  string             `json:"slug"` // 应用标识
	Route func(r chi.Router) `json:"-"`    // 路由
}

// AppCenter 应用中心结构
type AppCenter struct {
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	Channels    []struct {
		Slug      string `json:"slug"`
		Name      string `json:"name"`
		Panel     string `json:"panel"`
		Install   string `json:"-"`
		Uninstall string `json:"-"`
		Update    string `json:"-"`
		Subs      []struct {
			Log     string `json:"log"`
			Version string `json:"version"`
		} `json:"subs"`
	} `json:"channels"`
	Installed        bool   `json:"installed"`
	InstalledChannel string `json:"installed_channel"`
	InstalledVersion string `json:"installed_version"`
	UpdateExist      bool   `json:"update_exist"`
	Show             bool   `json:"show"`
}
