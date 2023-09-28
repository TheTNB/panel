package mysql80

var (
	Name        = "MySQL-8.0"
	Description = "MySQL 是最流行的关系型数据库管理系统之一，Oracle 旗下产品。(内存 < 4G 无法安装)"
	Slug        = "mysql80"
	Version     = "8.0.34"
	Requires    = []string{}
	Excludes    = []string{"mysql57"}
	Install     = `bash /www/panel/scripts/mysql/install.sh 80`
	Uninstall   = `bash /www/panel/scripts/mysql/uninstall.sh 80`
	Update      = `echo "not support now"`
)
