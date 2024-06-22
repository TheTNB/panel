package plugins

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/goravel/framework/contracts/http"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/app/http/controllers"
	"github.com/TheTNB/panel/app/models"
	"github.com/TheTNB/panel/internal"
	"github.com/TheTNB/panel/internal/services"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/types"
)

type MySQLController struct {
	setting internal.Setting
	backup  internal.Backup
}

func NewMySQLController() *MySQLController {
	return &MySQLController{
		setting: services.NewSettingImpl(),
		backup:  services.NewBackupImpl(),
	}
}

// GetConfig 获取配置
func (r *MySQLController) GetConfig(ctx http.Context) http.Response {
	config, err := tools.Read("/www/server/mysql/conf/my.cnf")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL配置失败")
	}

	return controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (r *MySQLController) SaveConfig(ctx http.Context) http.Response {
	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "配置不能为空")
	}

	if err := tools.Write("/www/server/mysql/conf/my.cnf", config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "写入MySQL配置失败")
	}

	if err := tools.ServiceReload("mysqld"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "重载MySQL失败")
	}

	return controllers.Success(ctx, nil)
}

// Load 获取负载
func (r *MySQLController) Load(ctx http.Context) http.Response {
	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	if len(rootPassword) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "MySQL root密码为空")
	}

	status, _ := tools.ServiceStatus("mysqld")
	if !status {
		return controllers.Success(ctx, []types.NV{})
	}

	raw, err := shell.Execf("/www/server/mysql/bin/mysqladmin -uroot -p" + rootPassword + " extended-status 2>&1")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL负载失败")
	}

	var data []map[string]string
	expressions := []struct {
		regex string
		name  string
	}{
		{`Uptime\s+\|\s+(\d+)\s+\|`, "运行时间"},
		{`Queries\s+\|\s+(\d+)\s+\|`, "总查询次数"},
		{`Connections\s+\|\s+(\d+)\s+\|`, "总连接次数"},
		{`Com_commit\s+\|\s+(\d+)\s+\|`, "每秒事务"},
		{`Com_rollback\s+\|\s+(\d+)\s+\|`, "每秒回滚"},
		{`Bytes_sent\s+\|\s+(\d+)\s+\|`, "发送"},
		{`Bytes_received\s+\|\s+(\d+)\s+\|`, "接收"},
		{`Threads_connected\s+\|\s+(\d+)\s+\|`, "活动连接数"},
		{`Max_used_connections\s+\|\s+(\d+)\s+\|`, "峰值连接数"},
		{`Key_read_requests\s+\|\s+(\d+)\s+\|`, "索引命中率"},
		{`Innodb_buffer_pool_reads\s+\|\s+(\d+)\s+\|`, "Innodb索引命中率"},
		{`Created_tmp_disk_tables\s+\|\s+(\d+)\s+\|`, "创建临时表到磁盘"},
		{`Open_tables\s+\|\s+(\d+)\s+\|`, "已打开的表"},
		{`Select_full_join\s+\|\s+(\d+)\s+\|`, "没有使用索引的量"},
		{`Select_full_range_join\s+\|\s+(\d+)\s+\|`, "没有索引的JOIN量"},
		{`Select_range_check\s+\|\s+(\d+)\s+\|`, "没有索引的子查询量"},
		{`Sort_merge_passes\s+\|\s+(\d+)\s+\|`, "排序后的合并次数"},
		{`Table_locks_waited\s+\|\s+(\d+)\s+\|`, "锁表次数"},
	}

	for _, expression := range expressions {
		re := regexp.MustCompile(expression.regex)
		matches := re.FindStringSubmatch(raw)
		if len(matches) > 1 {
			d := map[string]string{"name": expression.name, "value": matches[1]}
			if expression.name == "发送" || expression.name == "接收" {
				d["value"] = tools.FormatBytes(cast.ToFloat64(matches[1]))
			}

			data = append(data, d)
		}
	}

	// 索引命中率
	readRequests := cast.ToFloat64(data[9]["value"])
	reads := cast.ToFloat64(data[10]["value"])
	data[9]["value"] = fmt.Sprintf("%.2f%%", readRequests/(reads+readRequests)*100)
	// Innodb 索引命中率
	bufferPoolReads := cast.ToFloat64(data[11]["value"])
	bufferPoolReadRequests := cast.ToFloat64(data[12]["value"])
	data[10]["value"] = fmt.Sprintf("%.2f%%", bufferPoolReadRequests/(bufferPoolReads+bufferPoolReadRequests)*100)

	return controllers.Success(ctx, data)
}

// ErrorLog 获取错误日志
func (r *MySQLController) ErrorLog(ctx http.Context) http.Response {
	log, err := shell.Execf("tail -n 100 /www/server/mysql/mysql-error.log")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, log)
	}

	return controllers.Success(ctx, log)
}

// ClearErrorLog 清空错误日志
func (r *MySQLController) ClearErrorLog(ctx http.Context) http.Response {
	if out, err := shell.Execf("echo '' > /www/server/mysql/mysql-error.log"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// SlowLog 获取慢查询日志
func (r *MySQLController) SlowLog(ctx http.Context) http.Response {
	log, err := shell.Execf("tail -n 100 /www/server/mysql/mysql-slow.log")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, log)
	}

	return controllers.Success(ctx, log)
}

// ClearSlowLog 清空慢查询日志
func (r *MySQLController) ClearSlowLog(ctx http.Context) http.Response {
	if out, err := shell.Execf("echo '' > /www/server/mysql/mysql-slow.log"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	return controllers.Success(ctx, nil)
}

// GetRootPassword 获取root密码
func (r *MySQLController) GetRootPassword(ctx http.Context) http.Response {
	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	if len(rootPassword) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "MySQL root密码为空")
	}

	return controllers.Success(ctx, rootPassword)
}

// SetRootPassword 设置root密码
func (r *MySQLController) SetRootPassword(ctx http.Context) http.Response {
	status, err := tools.ServiceStatus("mysqld")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL状态失败")
	}
	if !status {
		return controllers.Error(ctx, http.StatusInternalServerError, "MySQL 未运行")
	}

	rootPassword := ctx.Request().Input("password")
	if len(rootPassword) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "MySQL root密码不能为空")
	}

	oldRootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	if oldRootPassword != rootPassword {
		if _, err = shell.Execf(fmt.Sprintf(`/www/server/mysql/bin/mysql -uroot -p%s -e "ALTER USER 'root'@'localhost' IDENTIFIED BY '%s';"`, oldRootPassword, rootPassword)); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("设置root密码失败: %v", err))
		}
		if _, err = shell.Execf(fmt.Sprintf(`/www/server/mysql/bin/mysql -uroot -p%s -e "FLUSH PRIVILEGES;"`, rootPassword)); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, "设置root密码失败")
		}
		if err = r.setting.Set(models.SettingKeyMysqlRootPassword, rootPassword); err != nil {
			_, _ = shell.Execf(fmt.Sprintf(`/www/server/mysql/bin/mysql -uroot -p%s -e "ALTER USER 'root'@'localhost' IDENTIFIED BY '%s';"`, rootPassword, oldRootPassword))
			_, _ = shell.Execf(fmt.Sprintf(`/www/server/mysql/bin/mysql -uroot -p%s -e "FLUSH PRIVILEGES;"`, oldRootPassword))
			return controllers.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("设置保存失败: %v", err))
		}
	}

	return controllers.Success(ctx, nil)
}

// DatabaseList 获取数据库列表
func (r *MySQLController) DatabaseList(ctx http.Context) http.Response {
	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	type database struct {
		Name string `json:"name"`
	}

	db, err := sql.Open("mysql", "root:"+rootPassword+"@unix(/tmp/mysql.sock)/")
	if err != nil {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []database{},
		})
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []database{},
		})
	}

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var databases []database
	for rows.Next() {
		var d database
		if err = rows.Scan(&d.Name); err != nil {
			continue
		}

		databases = append(databases, d)
	}

	if err = rows.Err(); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取数据库列表失败")
	}

	paged, total := controllers.Paginate(ctx, databases)

	return controllers.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// AddDatabase 添加数据库
func (r *MySQLController) AddDatabase(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255",
	}); sanitize != nil {
		return sanitize
	}

	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	database := ctx.Request().Input("database")
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")

	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"CREATE DATABASE IF NOT EXISTS " + database + " DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"CREATE USER '" + user + "'@'localhost' IDENTIFIED BY '" + password + "';\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"GRANT ALL PRIVILEGES ON " + database + ".* TO '" + user + "'@'localhost';\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// DeleteDatabase 删除数据库
func (r *MySQLController) DeleteDatabase(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$|not_in:information_schema,mysql,performance_schema,sys",
	}); sanitize != nil {
		return sanitize
	}

	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	database := ctx.Request().Input("database")
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"DROP DATABASE IF EXISTS " + database + ";\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// BackupList 获取备份列表
func (r *MySQLController) BackupList(ctx http.Context) http.Response {
	backups, err := r.backup.MysqlList()
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	paged, total := controllers.Paginate(ctx, backups)

	return controllers.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// UploadBackup 上传备份
func (r *MySQLController) UploadBackup(ctx http.Context) http.Response {
	file, err := ctx.Request().File("file")
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "上传文件失败")
	}

	backupPath := r.setting.Get(models.SettingKeyBackupPath) + "/mysql"
	if !tools.Exists(backupPath) {
		if err = tools.Mkdir(backupPath, 0644); err != nil {
			return nil
		}
	}

	name := file.GetClientOriginalName()
	_, err = file.StoreAs(backupPath, name)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "上传文件失败")
	}

	return controllers.Success(ctx, nil)
}

// CreateBackup 创建备份
func (r *MySQLController) CreateBackup(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$|not_in:information_schema,mysql,performance_schema,sys",
	}); sanitize != nil {
		return sanitize
	}

	database := ctx.Request().Input("database")
	if err := r.backup.MysqlBackup(database); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

// DeleteBackup 删除备份
func (r *MySQLController) DeleteBackup(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"name": "required|min_len:1|max_len:255",
	}); sanitize != nil {
		return sanitize
	}

	backupPath := r.setting.Get(models.SettingKeyBackupPath) + "/mysql"
	fileName := ctx.Request().Input("name")
	if err := tools.Remove(backupPath + "/" + fileName); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

// RestoreBackup 还原备份
func (r *MySQLController) RestoreBackup(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"backup":   "required|min_len:1|max_len:255",
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$|not_in:information_schema,mysql,performance_schema,sys",
	}); sanitize != nil {
		return sanitize
	}

	if err := r.backup.MysqlRestore(ctx.Request().Input("database"), ctx.Request().Input("backup")); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

// UserList 用户列表
func (r *MySQLController) UserList(ctx http.Context) http.Response {
	type user struct {
		User   string   `json:"user"`
		Host   string   `json:"host"`
		Grants []string `json:"grants"`
	}

	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	db, err := sql.Open("mysql", "root:"+rootPassword+"@unix(/tmp/mysql.sock)/")
	if err != nil {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []user{},
		})
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []user{},
		})
	}

	rows, err := db.Query("SELECT user, host FROM mysql.user")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var userGrants []user

	for rows.Next() {
		var u user
		if err = rows.Scan(&u.User, &u.Host); err != nil {
			continue
		}

		// 查询用户权限
		grantsRows, err := db.Query(fmt.Sprintf("SHOW GRANTS FOR '%s'@'%s'", u.User, u.Host))
		if err != nil {
			continue
		}
		defer grantsRows.Close()

		for grantsRows.Next() {
			var grant string
			if err = grantsRows.Scan(&grant); err != nil {
				continue
			}

			u.Grants = append(u.Grants, grant)
		}

		if err = grantsRows.Err(); err != nil {
			continue
		}

		userGrants = append(userGrants, u)
	}

	if err = rows.Err(); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取用户列表失败")
	}

	paged, total := controllers.Paginate(ctx, userGrants)

	return controllers.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// AddUser 添加用户
func (r *MySQLController) AddUser(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255",
	}); sanitize != nil {
		return sanitize
	}

	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	database := ctx.Request().Input("database")
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"CREATE USER '" + user + "'@'localhost' IDENTIFIED BY '" + password + ";'\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"GRANT ALL PRIVILEGES ON " + database + ".* TO '" + user + "'@'localhost';\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// DeleteUser 删除用户
func (r *MySQLController) DeleteUser(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"user": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
	}); sanitize != nil {
		return sanitize
	}

	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	user := ctx.Request().Input("user")
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"DROP USER '" + user + "'@'localhost';\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// SetUserPassword 设置用户密码
func (r *MySQLController) SetUserPassword(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255",
	}); sanitize != nil {
		return sanitize
	}

	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"ALTER USER '" + user + "'@'localhost' IDENTIFIED BY '" + password + "';\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// SetUserPrivileges 设置用户权限
func (r *MySQLController) SetUserPrivileges(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"database": "required|min_len:1|max_len:255",
	}); sanitize != nil {
		return sanitize
	}

	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	user := ctx.Request().Input("user")
	database := ctx.Request().Input("database")
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"REVOKE ALL PRIVILEGES ON *.* FROM '" + user + "'@'localhost';\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"GRANT ALL PRIVILEGES ON " + database + ".* TO '" + user + "'@'localhost';\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\""); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}
