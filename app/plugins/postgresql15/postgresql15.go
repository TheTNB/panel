package postgresql15

var (
	Name        = "PostgreSQL-15"
	Description = "PostgreSQL 是世界上最先进的开源关系数据库，在类似 BSD 与 MIT 许可的 PostgreSQL 许可下发行。"
	Slug        = "postgresql15"
	Version     = "15.5"
	Requires    = []string{}
	Excludes    = []string{"postgresql16"}
	Install     = `bash /www/panel/scripts/postgresql/install.sh 15`
	Uninstall   = `bash /www/panel/scripts/postgresql/uninstall.sh 15`
	Update      = `bash /www/panel/scripts/postgresql/update.sh 15`
)
