package bootstrap

import "panel/plugins/openresty"

func Plugins() {
	openresty.Boot()
}
