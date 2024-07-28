package types

import "github.com/goravel/framework/contracts/foundation"

// Plugin 插件元数据结构
type Plugin struct {
	Name        string                           // 插件名称
	Description string                           // 插件描述
	Slug        string                           // 插件标识
	Version     string                           // 插件版本
	Requires    []string                         // 依赖插件
	Excludes    []string                         // 排除插件
	Install     string                           // 安装命令
	Uninstall   string                           // 卸载命令
	Update      string                           // 更新命令
	Boot        func(app foundation.Application) // 启动时执行的命令
}
