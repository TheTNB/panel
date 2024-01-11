package openresty

var (
	Name        = "OpenResty"
	Description = "OpenResty® 是一款基于 NGINX 和 LuaJIT 的 Web 平台。"
	Slug        = "openresty"
	Version     = "1.25.3.1"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = "bash /www/panel/scripts/openresty/install.sh"
	Uninstall   = "bash /www/panel/scripts/openresty/uninstall.sh"
	Update      = "bash /www/panel/scripts/openresty/install.sh"
)
