package postgresql16

var (
	Name        = "PostgreSQL-16"
	Description = "PostgreSQL 是世界上最先进的开源关系数据库，在类似 BSD 与 MIT 许可的 PostgreSQL 许可下发行。"
	Slug        = "postgresql16"
	Version     = "16.2"
	Requires    = []string{}
	Excludes    = []string{"postgresql15"}
	Install     = `bash /www/panel/scripts/postgresql/install.sh 16`
	Uninstall   = `bash /www/panel/scripts/postgresql/uninstall.sh 16`
	Update      = `bash /www/panel/scripts/postgresql/update.sh 16`
)
