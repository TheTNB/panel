package mysql57

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"

	"panel/app/http/controllers"
	"panel/app/models"
	"panel/app/services"
	"panel/pkg/tools"
)

type Mysql57Controller struct {
	setting services.Setting
	backup  services.Backup
}

func NewMysql57Controller() *Mysql57Controller {
	return &Mysql57Controller{
		setting: services.NewSettingImpl(),
		backup:  services.NewBackupImpl(),
	}
}

// Status 获取运行状态
func (c *Mysql57Controller) Status(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	status := tools.ExecShell("systemctl status mysqld | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Reload 重载配置
func (c *Mysql57Controller) Reload(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	tools.ExecShell("systemctl reload mysqld")
	status := tools.ExecShell("systemctl status mysqld | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Restart 重启服务
func (c *Mysql57Controller) Restart(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	tools.ExecShell("systemctl restart mysqld")
	status := tools.ExecShell("systemctl status mysqld | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Start 启动服务
func (c *Mysql57Controller) Start(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	tools.ExecShell("systemctl start mysqld")
	status := tools.ExecShell("systemctl status mysqld | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Stop 停止服务
func (c *Mysql57Controller) Stop(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	tools.ExecShell("systemctl stop mysqld")
	status := tools.ExecShell("systemctl status mysqld | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL状态失败")
		return
	}

	if status != "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// GetConfig 获取配置
func (c *Mysql57Controller) GetConfig(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	// 获取配置
	config := tools.ReadFile("/www/server/mysql/conf/my.cnf")
	if len(config) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL配置失败")
		return
	}

	controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (c *Mysql57Controller) SaveConfig(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	config := ctx.Request().Input("config")
	if len(config) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "配置不能为空")
		return
	}

	if !tools.WriteFile("/www/server/mysql/conf/my.cnf", config, 0644) {
		controllers.Error(ctx, http.StatusInternalServerError, "写入MySQL配置失败")
		return
	}

	controllers.Success(ctx, "保存MySQL配置成功")
}

// Load 获取负载
func (c *Mysql57Controller) Load(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	rootPassword := c.setting.Get(models.SettingKeyMysqlRootPassword)
	if len(rootPassword) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "MySQL root密码为空")
		return
	}

	status := tools.ExecShell("systemctl status mysqld | grep Active | grep -v grep | awk '{print $2}'")
	if status != "active" {
		controllers.Error(ctx, http.StatusInternalServerError, "MySQL 已停止运行")
		return
	}

	raw := tools.ExecShell("/www/server/mysql/bin/mysqladmin -uroot -p" + rootPassword + " extended-status 2>&1")
	if strings.Contains(raw, "Access denied for user") {
		controllers.Error(ctx, http.StatusBadRequest, "MySQL root密码错误")
		return
	}
	if !strings.Contains(raw, "Uptime") {
		controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL负载失败")
		return
	}

	data := make(map[int]map[string]string)
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

	for i, expression := range expressions {
		re := regexp.MustCompile(expression.regex)
		matches := re.FindStringSubmatch(raw)
		if len(matches) > 1 {
			data[i] = make(map[string]string)
			data[i] = map[string]string{"name": expression.name, "value": matches[1]}

			if expression.name == "发送" || expression.name == "接收" {
				data[i]["value"] = tools.FormatBytes(cast.ToFloat64(matches[1]))
			}
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

	controllers.Success(ctx, data)
}

// ErrorLog 获取错误日志
func (c *Mysql57Controller) ErrorLog(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	log := tools.ExecShell("tail -n 100 /www/server/mysql/mysql-error.log")
	controllers.Success(ctx, log)
}

// ClearErrorLog 清空错误日志
func (c *Mysql57Controller) ClearErrorLog(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	tools.ExecShell("echo '' > /www/server/mysql/mysql-error.log")
	controllers.Success(ctx, "清空错误日志成功")
}

// SlowLog 获取慢查询日志
func (c *Mysql57Controller) SlowLog(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	log := tools.ExecShell("tail -n 100 /www/server/mysql/mysql-slow.log")
	controllers.Success(ctx, log)
}

// ClearSlowLog 清空慢查询日志
func (c *Mysql57Controller) ClearSlowLog(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	tools.ExecShell("echo '' > /www/server/mysql/mysql-slow.log")
	controllers.Success(ctx, "清空慢查询日志成功")
}

// GetRootPassword 获取root密码
func (c *Mysql57Controller) GetRootPassword(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	rootPassword := c.setting.Get(models.SettingKeyMysqlRootPassword)
	if len(rootPassword) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "MySQL root密码为空")
		return
	}

	controllers.Success(ctx, rootPassword)
}

// SetRootPassword 设置root密码
func (c *Mysql57Controller) SetRootPassword(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	status := tools.ExecShell("systemctl status mysqld | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL状态失败")
		return
	}
	if status != "active" {
		controllers.Error(ctx, http.StatusInternalServerError, "MySQL 未运行")
		return
	}

	rootPassword := ctx.Request().Input(models.SettingKeyMysqlRootPassword)
	if len(rootPassword) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "MySQL root密码不能为空")
		return
	}

	oldRootPassword := c.setting.Get(models.SettingKeyMysqlRootPassword)
	if oldRootPassword != rootPassword {
		tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + oldRootPassword + " -e \"ALTER USER 'root'@'localhost' IDENTIFIED BY '" + rootPassword + "';\"")
		tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + oldRootPassword + " -e \"FLUSH PRIVILEGES;\"")
		err := c.setting.Set(models.SettingKeyMysqlRootPassword, rootPassword)
		if err != nil {
			tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"ALTER USER 'root'@'localhost' IDENTIFIED BY '" + oldRootPassword + "';\"")
			tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\"")
			controllers.Error(ctx, http.StatusInternalServerError, "设置root密码失败")
			return
		}
	}

	controllers.Success(ctx, "设置root密码成功")
}

// DatabaseList 获取数据库列表
func (c *Mysql57Controller) DatabaseList(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	rootPassword := c.setting.Get(models.SettingKeyMysqlRootPassword)
	type database struct {
		Name string `json:"name"`
	}

	db, err := sql.Open("mysql", "root:"+rootPassword+"@unix(/tmp/mysql.sock)/")
	if err != nil {
		facades.Log().Error("[MySQL57] 连接数据库失败" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "连接数据库失败")
		return
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		facades.Log().Error("[MySQL57] 获取数据库列表失败" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取数据库列表失败")
		return
	}
	defer rows.Close()

	var databases []database
	for rows.Next() {
		var d database
		err := rows.Scan(&d.Name)
		if err != nil {
			continue
		}

		databases = append(databases, d)
	}

	if err := rows.Err(); err != nil {
		facades.Log().Error("[MySQL57] 获取数据库列表失败" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取数据库列表失败")
		return
	}

	controllers.Success(ctx, databases)
}

// AddDatabase 添加数据库
func (c *Mysql57Controller) AddDatabase(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	rootPassword := c.setting.Get(models.SettingKeyMysqlRootPassword)
	database := ctx.Request().Input("database")
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")

	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"CREATE DATABASE IF NOT EXISTS " + database + " DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;\"")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"CREATE USER '" + user + "'@'localhost' IDENTIFIED BY '" + password + "';\"")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"GRANT ALL PRIVILEGES ON " + database + ".* TO '" + user + "'@'localhost';\"")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\"")

	controllers.Success(ctx, "添加数据库成功")
}

// DeleteDatabase 删除数据库
func (c *Mysql57Controller) DeleteDatabase(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$|not_in:information_schema,mysql,performance_schema,sys",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	rootPassword := c.setting.Get(models.SettingKeyMysqlRootPassword)
	database := ctx.Request().Input("database")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"DROP DATABASE IF EXISTS " + database + ";\"")

	controllers.Success(ctx, "删除数据库成功")
}

// BackupList 获取备份列表
func (c *Mysql57Controller) BackupList(ctx http.Context) {
	backupList, err := c.backup.MysqlList()
	if err != nil {
		facades.Log().Error("[MySQL57] 获取备份列表失败：" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取备份列表失败")
		return
	}

	controllers.Success(ctx, backupList)
}

// UploadBackup 上传备份
func (c *Mysql57Controller) UploadBackup(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	file, err := ctx.Request().File("file")
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "上传文件失败")
		return
	}

	backupPath := c.setting.Get(models.SettingKeyBackupPath) + "/mysql"
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}

	name := file.GetClientOriginalName()
	_, err = file.StoreAs(backupPath, name)
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "上传文件失败")
		return
	}

	controllers.Success(ctx, "上传文件成功")
}

// CreateBackup 创建备份
func (c *Mysql57Controller) CreateBackup(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$|not_in:information_schema,mysql,performance_schema,sys",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	database := ctx.Request().Input("database")
	err = c.backup.MysqlBackup(database)
	if err != nil {
		facades.Log().Error("[MYSQL57] 创建备份失败：" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "创建备份失败")
		return
	}

	controllers.Success(ctx, "备份成功")
}

// DeleteBackup 删除备份
func (c *Mysql57Controller) DeleteBackup(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"name": "required|min_len:1|max_len:255",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	backupPath := c.setting.Get(models.SettingKeyBackupPath) + "/mysql"
	fileName := ctx.Request().Input("name")
	tools.RemoveFile(backupPath + "/" + fileName)

	controllers.Success(ctx, "删除备份成功")
}

// RestoreBackup 还原备份
func (c *Mysql57Controller) RestoreBackup(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"name":     "required|min_len:1|max_len:255",
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$|not_in:information_schema,mysql,performance_schema,sys",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	err = c.backup.MysqlRestore(ctx.Request().Input("database"), ctx.Request().Input("name"))
	if err != nil {
		facades.Log().Error("[MYSQL57] 还原失败：" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "还原失败: "+err.Error())
		return
	}

	controllers.Success(ctx, "还原成功")
}

// UserList 用户列表
func (c *Mysql57Controller) UserList(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	type user struct {
		User   string   `json:"user"`
		Host   string   `json:"host"`
		Grants []string `json:"grants"`
	}

	rootPassword := c.setting.Get(models.SettingKeyMysqlRootPassword)
	db, err := sql.Open("mysql", "root:"+rootPassword+"@unix(/tmp/mysql.sock)/")
	if err != nil {
		facades.Log().Error("[MYSQL57] 连接数据库失败：" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "连接数据库失败")
	}
	defer db.Close()

	rows, err := db.Query("SELECT user, host FROM mysql.user")
	if err != nil {
		facades.Log().Error("[MYSQL57] 查询数据库失败：" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "查询数据库失败")
	}
	defer rows.Close()

	var userGrants []user

	for rows.Next() {
		var u user
		err := rows.Scan(&u.User, &u.Host)
		if err != nil {
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
			err := grantsRows.Scan(&grant)
			if err != nil {
				continue
			}

			u.Grants = append(u.Grants, grant)
		}

		if err := grantsRows.Err(); err != nil {
			continue
		}

		userGrants = append(userGrants, u)
	}

	if err := rows.Err(); err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	controllers.Success(ctx, userGrants)
}

// AddUser 添加用户
func (c *Mysql57Controller) AddUser(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	rootPassword := c.setting.Get(models.SettingKeyMysqlRootPassword)
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	database := ctx.Request().Input("database")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"CREATE USER '" + user + "'@'localhost' IDENTIFIED BY '" + password + ";'\"")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"GRANT ALL PRIVILEGES ON " + database + ".* TO '" + user + "'@'localhost';\"")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\"")

	controllers.Success(ctx, "添加成功")
}

// DeleteUser 删除用户
func (c *Mysql57Controller) DeleteUser(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"user": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	rootPassword := c.setting.Get(models.SettingKeyMysqlRootPassword)
	user := ctx.Request().Input("user")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"DROP USER '" + user + "'@'localhost';\"")

	controllers.Success(ctx, "删除成功")
}

// SetUserPassword 设置用户密码
func (c *Mysql57Controller) SetUserPassword(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	rootPassword := c.setting.Get(models.SettingKeyMysqlRootPassword)
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"ALTER USER '" + user + "'@'localhost' IDENTIFIED BY '" + password + "';\"")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\"")

	controllers.Success(ctx, "修改成功")
}

// SetUserPrivileges 设置用户权限
func (c *Mysql57Controller) SetUserPrivileges(ctx http.Context) {
	if !controllers.Check(ctx, "mysql57") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"database": "required|min_len:1|max_len:255",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	rootPassword := c.setting.Get(models.SettingKeyMysqlRootPassword)
	user := ctx.Request().Input("user")
	database := ctx.Request().Input("database")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"REVOKE ALL PRIVILEGES ON *.* FROM '" + user + "'@'localhost';\"")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"GRANT ALL PRIVILEGES ON " + database + ".* TO '" + user + "'@'localhost';\"")
	tools.ExecShell("/www/server/mysql/bin/mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\"")

	controllers.Success(ctx, "修改成功")
}
