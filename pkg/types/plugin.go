package types

import "github.com/go-chi/chi/v5"

// Plugin 插件元数据结构
type Plugin struct {
	Slug        string             `json:"slug"`        // 插件标识
	Name        string             `json:"name"`        // 插件名称
	Description string             `json:"description"` // 插件描述
	Version     string             `json:"version"`     // 插件版本
	Requires    []string           `json:"requires"`    // 依赖插件
	Excludes    []string           `json:"excludes"`    // 排除插件
	Install     string             `json:"-"`           // 安装命令
	Uninstall   string             `json:"-"`           // 卸载命令
	Update      string             `json:"-"`           // 更新命令
	Route       func(r chi.Router) `json:"-"`           // 路由
}
