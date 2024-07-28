package plugins

import (
	"database/sql"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
	_ "github.com/lib/pq"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type PostgreSQLController struct {
	setting internal.Setting
	backup  internal.Backup
}

func NewPostgreSQLController() *PostgreSQLController {
	return &PostgreSQLController{
		setting: services.NewSettingImpl(),
		backup:  services.NewBackupImpl(),
	}
}

// GetConfig 获取配置
func (r *PostgreSQLController) GetConfig(ctx http.Context) http.Response {
	// 获取配置
	config, err := io.Read("/www/server/postgresql/data/postgresql.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL配置失败")
	}

	return h.Success(ctx, config)
}

// GetUserConfig 获取用户配置
func (r *PostgreSQLController) GetUserConfig(ctx http.Context) http.Response {
	// 获取配置
	config, err := io.Read("/www/server/postgresql/data/pg_hba.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL配置失败")
	}

	return h.Success(ctx, config)
}

// SaveConfig 保存配置
func (r *PostgreSQLController) SaveConfig(ctx http.Context) http.Response {
	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "配置不能为空")
	}

	if err := io.Write("/www/server/postgresql/data/postgresql.conf", config, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "写入PostgreSQL配置失败")
	}

	if err := systemctl.Reload("postgresql"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "重载服务失败")
	}

	return h.Success(ctx, nil)
}

// SaveUserConfig 保存用户配置
func (r *PostgreSQLController) SaveUserConfig(ctx http.Context) http.Response {
	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "配置不能为空")
	}

	if err := io.Write("/www/server/postgresql/data/pg_hba.conf", config, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "写入PostgreSQL配置失败")
	}

	if err := systemctl.Reload("postgresql"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "重载服务失败")
	}

	return h.Success(ctx, nil)
}

// Load 获取负载
func (r *PostgreSQLController) Load(ctx http.Context) http.Response {
	status, _ := systemctl.Status("postgresql")
	if !status {
		return h.Success(ctx, []types.NV{})
	}

	time, err := shell.Execf(`echo "select pg_postmaster_start_time();" | su - postgres -c "psql" | sed -n 3p | cut -d'.' -f1`)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL启动时间失败")
	}
	pid, err := shell.Execf(`echo "select pg_backend_pid();" | su - postgres -c "psql" | sed -n 3p`)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL进程PID失败")
	}
	process, err := shell.Execf(`ps aux | grep postgres | grep -v grep | wc -l`)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL进程数失败")
	}
	connections, err := shell.Execf(`echo "SELECT count(*) FROM pg_stat_activity WHERE NOT pid=pg_backend_pid();" | su - postgres -c "psql" | sed -n 3p`)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL连接数失败")
	}
	storage, err := shell.Execf(`echo "select pg_size_pretty(pg_database_size('postgres'));" | su - postgres -c "psql" | sed -n 3p`)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL空间占用失败")
	}

	data := []types.NV{
		{Name: "启动时间", Value: carbon.Parse(time).ToDateTimeString()},
		{Name: "进程 PID", Value: pid},
		{Name: "进程数", Value: process},
		{Name: "总连接数", Value: connections},
		{Name: "空间占用", Value: storage},
	}

	return h.Success(ctx, data)
}

// Log 获取日志
func (r *PostgreSQLController) Log(ctx http.Context) http.Response {
	log, err := shell.Execf("tail -n 100 /www/server/postgresql/logs/postgresql-" + carbon.Now().ToDateString() + ".log")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, log)
	}

	return h.Success(ctx, log)
}

// ClearLog 清空日志
func (r *PostgreSQLController) ClearLog(ctx http.Context) http.Response {
	if out, err := shell.Execf("echo '' > /www/server/postgresql/logs/postgresql-" + carbon.Now().ToDateString() + ".log"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}

// DatabaseList 获取数据库列表
func (r *PostgreSQLController) DatabaseList(ctx http.Context) http.Response {
	type database struct {
		Name     string `json:"name"`
		Owner    string `json:"owner"`
		Encoding string `json:"encoding"`
	}

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return h.Success(ctx, http.Json{
			"total": 0,
			"items": []database{},
		})
	}

	if err = db.Ping(); err != nil {
		return h.Success(ctx, http.Json{
			"total": 0,
			"items": []database{},
		})
	}

	query := `
        SELECT d.datname, pg_catalog.pg_get_userbyid(d.datdba), pg_catalog.pg_encoding_to_char(d.encoding)
        FROM pg_catalog.pg_database d
        WHERE datistemplate = false;
    `
	rows, err := db.Query(query)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var databases []database
	for rows.Next() {
		var db database
		if err := rows.Scan(&db.Name, &db.Owner, &db.Encoding); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		databases = append(databases, db)
	}
	if err = rows.Err(); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	paged, total := h.Paginate(ctx, databases)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// AddDatabase 添加数据库
func (r *PostgreSQLController) AddDatabase(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"database": "required|min_len:1|max_len:63|regex:^[a-zA-Z0-9_]+$",
		"user":     "required|min_len:1|max_len:30|regex:^[a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:40",
	}); sanitize != nil {
		return sanitize
	}

	database := ctx.Request().Input("database")
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")

	if out, err := shell.Execf(`echo "CREATE DATABASE ` + database + `;" | su - postgres -c "psql"`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf(`echo "CREATE USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf(`echo "ALTER DATABASE ` + database + ` OWNER TO ` + user + `;" | su - postgres -c "psql"`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf(`echo "GRANT ALL PRIVILEGES ON DATABASE ` + database + ` TO ` + user + `;" | su - postgres -c "psql"`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	userConfig := "host    " + database + "    " + user + "    127.0.0.1/32    scram-sha-256"
	if out, err := shell.Execf(`echo "` + userConfig + `" >> /www/server/postgresql/data/pg_hba.conf`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	if err := systemctl.Reload("postgresql"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "重载服务失败")
	}

	return h.Success(ctx, nil)
}

// DeleteDatabase 删除数据库
func (r *PostgreSQLController) DeleteDatabase(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"database": "required|min_len:1|max_len:63|regex:^[a-zA-Z0-9_]+$|not_in:postgres,template0,template1",
	}); sanitize != nil {
		return sanitize
	}

	database := ctx.Request().Input("database")
	if out, err := shell.Execf(`echo "DROP DATABASE ` + database + `;" | su - postgres -c "psql"`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}

// BackupList 获取备份列表
func (r *PostgreSQLController) BackupList(ctx http.Context) http.Response {
	backups, err := r.backup.PostgresqlList()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取备份列表失败")
	}

	paged, total := h.Paginate(ctx, backups)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// UploadBackup 上传备份
func (r *PostgreSQLController) UploadBackup(ctx http.Context) http.Response {
	file, err := ctx.Request().File("file")
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "上传文件失败")
	}

	backupPath := r.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
	if !io.Exists(backupPath) {
		if err = io.Mkdir(backupPath, 0644); err != nil {
			return nil
		}
	}

	name := file.GetClientOriginalName()
	_, err = file.StoreAs(backupPath, name)
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "上传文件失败")
	}

	return h.Success(ctx, nil)
}

// CreateBackup 创建备份
func (r *PostgreSQLController) CreateBackup(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"database": "required|min_len:1|max_len:63|regex:^[a-zA-Z0-9_]+$|not_in:postgres,template0,template1",
	}); sanitize != nil {
		return sanitize
	}

	database := ctx.Request().Input("database")
	if err := r.backup.PostgresqlBackup(database); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// DeleteBackup 删除备份
func (r *PostgreSQLController) DeleteBackup(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"name": "required|min_len:1|max_len:255",
	}); sanitize != nil {
		return sanitize
	}

	backupPath := r.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
	fileName := ctx.Request().Input("name")
	if err := io.Remove(backupPath + "/" + fileName); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// RestoreBackup 还原备份
func (r *PostgreSQLController) RestoreBackup(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"backup":   "required|min_len:1|max_len:255",
		"database": "required|min_len:1|max_len:63|regex:^[a-zA-Z0-9_]+$|not_in:postgres,template0,template1",
	}); sanitize != nil {
		return sanitize
	}

	if err := r.backup.PostgresqlRestore(ctx.Request().Input("database"), ctx.Request().Input("backup")); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "还原失败: "+err.Error())
	}

	return h.Success(ctx, nil)
}

// RoleList 角色列表
func (r *PostgreSQLController) RoleList(ctx http.Context) http.Response {
	type role struct {
		Role       string   `json:"role"`
		Attributes []string `json:"attributes"`
	}

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return h.Success(ctx, http.Json{
			"total": 0,
			"items": []role{},
		})
	}
	if err = db.Ping(); err != nil {
		return h.Success(ctx, http.Json{
			"total": 0,
			"items": []role{},
		})
	}

	query := `
        SELECT rolname,
               rolsuper,
               rolcreaterole,
               rolcreatedb,
               rolreplication,
               rolbypassrls
        FROM pg_roles
        WHERE rolcanlogin = true;
    `
	rows, err := db.Query(query)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var roles []role
	for rows.Next() {
		var r role
		var super, canCreateRole, canCreateDb, replication, bypassRls bool
		if err = rows.Scan(&r.Role, &super, &canCreateRole, &canCreateDb, &replication, &bypassRls); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}

		permissions := map[string]bool{
			"超级用户":   super,
			"创建角色":   canCreateRole,
			"创建数据库":  canCreateDb,
			"可以复制":   replication,
			"绕过行级安全": bypassRls,
		}
		for perm, enabled := range permissions {
			if enabled {
				r.Attributes = append(r.Attributes, perm)
			}
		}

		if len(r.Attributes) == 0 {
			r.Attributes = append(r.Attributes, "无")
		}

		roles = append(roles, r)
	}
	if err = rows.Err(); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	paged, total := h.Paginate(ctx, roles)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// AddRole 添加角色
func (r *PostgreSQLController) AddRole(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"database": "required|min_len:1|max_len:63|regex:^[a-zA-Z0-9_]+$",
		"user":     "required|min_len:1|max_len:30|regex:^[a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:40",
	}); sanitize != nil {
		return sanitize
	}

	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	database := ctx.Request().Input("database")
	if out, err := shell.Execf(`echo "CREATE USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf(`echo "GRANT ALL PRIVILEGES ON DATABASE ` + database + ` TO ` + user + `;" | su - postgres -c "psql"`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	userConfig := "host    " + database + "    " + user + "    127.0.0.1/32    scram-sha-256"
	if out, err := shell.Execf(`echo "` + userConfig + `" >> /www/server/postgresql/data/pg_hba.conf`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	if err := systemctl.Reload("postgresql"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "重载服务失败")
	}

	return h.Success(ctx, nil)
}

// DeleteRole 删除角色
func (r *PostgreSQLController) DeleteRole(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"user": "required|min_len:1|max_len:30|regex:^[a-zA-Z0-9_]+$",
	}); sanitize != nil {
		return sanitize
	}

	user := ctx.Request().Input("user")
	if out, err := shell.Execf(`echo "DROP USER ` + user + `;" | su - postgres -c "psql"`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf(`sed -i '/` + user + `/d' /www/server/postgresql/data/pg_hba.conf`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	if err := systemctl.Reload("postgresql"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "重载服务失败")
	}

	return h.Success(ctx, nil)
}

// SetRolePassword 设置用户密码
func (r *PostgreSQLController) SetRolePassword(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"user":     "required|min_len:1|max_len:30|regex:^[a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:40",
	}); sanitize != nil {
		return sanitize
	}

	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	if out, err := shell.Execf(`echo "ALTER USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}
