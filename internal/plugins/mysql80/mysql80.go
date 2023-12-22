package mysql80

var (
	Name        = "MySQL-8.0"
	Description = "MySQL 是最流行的关系型数据库管理系统之一，Oracle 旗下产品。（建议内存 > 2G 安装）"
	Slug        = "mysql80"
	Version     = "8.0.35"
	Requires    = []string{}
	Excludes    = []string{"mysql57"}
	Install     = `bash /www/panel/scripts/mysql/install.sh 80`
	Uninstall   = `bash /www/panel/scripts/mysql/uninstall.sh 80`
	Update      = `bash /www/panel/scripts/mysql/update.sh 80`
)
