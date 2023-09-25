package postgresql16

var (
	Name        = "PostgreSQL-16"
	Description = "PostgreSQL 是开源的对象 - 关系数据库数据库管理系统，在类似 BSD 许可与 MIT 许可的 PostgreSQL 许可下发行。"
	Slug        = "postgresql16"
	Version     = "16.0"
	Requires    = []string{}
	Excludes    = []string{"postgresql15"}
	Install     = `bash /www/panel/scripts/postgresql/install.sh 16`
	Uninstall   = `bash /www/panel/scripts/postgresql/uninstall.sh 16`
	Update      = `bash /www/panel/scripts/postgresql/update.sh 16`
)
