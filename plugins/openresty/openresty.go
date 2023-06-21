package openresty

const (
	Name        = "OpenResty"
	Author      = "耗子"
	Description = "OpenResty® 是一款基于 NGINX 和 LuaJIT 的 Web 平台。"
	Slug        = "openresty"
	Version     = "1.21.4.1"
)

func Boot() {
	Route()
}
