package phpmyadmin

var (
	Name        = "phpMyAdmin"
	Description = "phpMyAdmin 是一个以 PHP 为基础，以 Web-Base 方式架构在网站主机上的 MySQL 数据库管理工具。"
	Slug        = "phpmyadmin"
	Version     = "5.2.1"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/phpmyadmin/install.sh`
	Uninstall   = `bash /www/panel/scripts/phpmyadmin/uninstall.sh`
	Update      = `bash /www/panel/scripts/phpmyadmin/uninstall.sh && bash /www/panel/scripts/phpmyadmin/install.sh`
)
