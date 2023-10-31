package certbot

var (
	Name        = "证书管理器"
	Description = "证书管理器使用 ACME 协议为服务器从证书颁发机构自动获取 HTTPS 证书。"
	Slug        = "certbot"
	Version     = "1.0.0"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/certbot/install.sh`
	Uninstall   = `bash /www/panel/scripts/certbot/uninstall.sh`
	Update      = `bash /www/panel/scripts/certbot/update.sh`
)
