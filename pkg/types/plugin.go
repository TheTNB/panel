package types

// Plugin 插件元数据结构
type Plugin struct {
	Name        string
	Description string
	Slug        string
	Version     string
	Requires    []string
	Excludes    []string
	Install     string
	Uninstall   string
	Update      string
}
