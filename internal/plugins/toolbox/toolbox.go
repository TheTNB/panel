package toolbox

var (
	Name        = "系统工具箱"
	Description = "可视化调整一些常用的配置项，如 DNS、SWAP、时区 等"
	Slug        = "toolbox"
	Version     = "1.0"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `panel writePlugin toolbox 1.0`
	Uninstall   = `panel deletePlugin toolbox`
	Update      = `panel writePlugin toolbox 1.0`
)
