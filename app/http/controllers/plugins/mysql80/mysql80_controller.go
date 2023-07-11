package mysql80

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/spf13/cast"
	"golang.org/x/exp/slices"

	"panel/app/http/controllers"
	"panel/app/http/controllers/plugins"
	"panel/app/services"
	"panel/packages/helpers"
)

type Mysql80Controller struct {
	setting services.Setting
}

func NewMysql80Controller() *Mysql80Controller {
	return &Mysql80Controller{
		setting: services.NewSettingImpl(),
	}
}

// Status 获取运行状态
func (r *Mysql80Controller) Status(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	out := helpers.ExecShell("systemctl status mysql | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
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
func (r *Mysql80Controller) Reload(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	_ = helpers.ExecShell("systemctl reload mysql")
	out := helpers.ExecShell("systemctl status mysql | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
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
func (r *Mysql80Controller) Restart(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	_ = helpers.ExecShell("systemctl restart mysql")
	out := helpers.ExecShell("systemctl status mysql | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
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
func (r *Mysql80Controller) Start(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	_ = helpers.ExecShell("systemctl start mysql")
	out := helpers.ExecShell("systemctl status mysql | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
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
func (r *Mysql80Controller) Stop(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	_ = helpers.ExecShell("systemctl stop mysql")
	out := helpers.ExecShell("systemctl status mysql | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
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

// GetConfig 获取配置
func (r *Mysql80Controller) GetConfig(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	// 获取配置
	config := helpers.ReadFile("mysql80")
	if len(config) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL配置失败")
		return
	}

	controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (r *Mysql80Controller) SaveConfig(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	config := ctx.Request().Input("config")
	if len(config) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "配置不能为空")
		return
	}

	if !helpers.WriteFile("mysql80", config, 0644) {
		controllers.Error(ctx, http.StatusInternalServerError, "写入MySQL配置失败")
		return
	}

	controllers.Success(ctx, "保存MySQL配置成功")
}

// Load 获取负载
func (r *Mysql80Controller) Load(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	rootPassword := r.setting.Get("mysql_root_password")
	if len(rootPassword) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "MySQL root密码为空")
		return
	}

	status := helpers.ExecShell("systemctl status mysqld | grep Active | grep -v grep | awk '{print $2}'")
	if strings.TrimSpace(status) != "active" {
		controllers.Error(ctx, http.StatusInternalServerError, "MySQL 已停止运行")
		return
	}

	raw := helpers.ExecShell("mysqladmin -uroot -p" + rootPassword + " extended-status 2>&1")
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
		{`Uptime\s+\|\s+(\d+)\s+\|`, "总查询次数"},
		{`Queries\s+\|\s+(\d+)\s+\|`, "总连接次数"},
		{`Connections\s+\|\s+(\d+)\s+\|`, "每秒事务"},
		{`Com_commit\s+\|\s+(\d+)\s+\|`, "每秒回滚"},
		{`Com_rollback\s+\|\s+(\d+)\s+\|`, "发送"},
		{`Bytes_sent\s+\|\s+(\d+)\s+\|`, "接收"},
		{`Bytes_received\s+\|\s+(\d+)\s+\|`, "活动连接数"},
		{`Threads_connected\s+\|\s+(\d+)\s+\|`, "峰值连接数"},
		{`Max_used_connections\s+\|\s+(\d+)\s+\|`, "索引命中率"},
		{`Key_read_requests\s+\|\s+(\d+)\s+\|`, "Innodb索引命中率"},
		{`Innodb_buffer_pool_reads\s+\|\s+(\d+)\s+\|`, "创建临时表到磁盘"},
		{`Created_tmp_disk_tables\s+\|\s+(\d+)\s+\|`, "已打开的表"},
		{`Open_tables\s+\|\s+(\d+)\s+\|`, "没有使用索引的量"},
		{`Select_full_join\s+\|\s+(\d+)\s+\|`, "没有索引的JOIN量"},
		{`Select_full_range_join\s+\|\s+(\d+)\s+\|`, "没有索引的子查询量"},
		{`Select_range_check\s+\|\s+(\d+)\s+\|`, "排序后的合并次数"},
		{`Sort_merge_passes\s+\|\s+(\d+)\s+\|`, "锁表次数"},
		{`Table_locks_waited\s+\|\s+(\d+)\s+\|`, ""},
	}

	for i, expression := range expressions {
		re := regexp.MustCompile(expression.regex)
		matches := re.FindStringSubmatch(raw)
		if len(matches) > 1 {
			data[i] = make(map[string]string)
			data[i] = map[string]string{"name": expression.name, "value": matches[1]}

			if expression.name == "发送" || expression.name == "接收" {
				data[i]["value"] = helpers.FormatBytes(cast.ToFloat64(matches[1]))
			}
		}
	}

	// 索引命中率
	readRequests := cast.ToFloat64(data[9]["value"])
	reads := cast.ToFloat64(data[10]["value"])
	data[9]["value"] = fmt.Sprintf("%.2f%%", readRequests/(reads+readRequests)*100)
	// Innodb索引命中率
	bufferPoolReads := cast.ToFloat64(data[11]["value"])
	bufferPoolReadRequests := cast.ToFloat64(data[12]["value"])
	data[10]["value"] = fmt.Sprintf("%.2f%%", bufferPoolReadRequests/(bufferPoolReads+bufferPoolReadRequests)*100)

	controllers.Success(ctx, data)
}

// ErrorLog 获取错误日志
func (r *Mysql80Controller) ErrorLog(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	log := helpers.ExecShell("tail -n 100 /www/server/mysql/mysql-error.log")
	controllers.Success(ctx, log)
}

// ClearErrorLog 清空错误日志
func (r *Mysql80Controller) ClearErrorLog(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	helpers.ExecShell("echo '' > /www/server/mysql/mysql-error.log")
	controllers.Success(ctx, "清空错误日志成功")
}

// SlowLog 获取慢查询日志
func (r *Mysql80Controller) SlowLog(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	log := helpers.ExecShell("tail -n 100 /www/server/mysql/mysql-slow.log")
	controllers.Success(ctx, log)
}

// ClearSlowLog 清空慢查询日志
func (r *Mysql80Controller) ClearSlowLog(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	helpers.ExecShell("echo '' > /www/server/mysql/mysql-slow.log")
	controllers.Success(ctx, "清空慢查询日志成功")
}

// GetRootPassword 获取root密码
func (r *Mysql80Controller) GetRootPassword(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	rootPassword := r.setting.Get("mysql_root_password")
	if len(rootPassword) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "MySQL root密码为空")
		return
	}

	controllers.Success(ctx, rootPassword)
}

// SetRootPassword 设置root密码
func (r *Mysql80Controller) SetRootPassword(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	out := helpers.ExecShell("systemctl status mysql | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取MySQL状态失败")
		return
	}
	if status != "active" {
		controllers.Error(ctx, http.StatusInternalServerError, "MySQL 未运行")
		return
	}

	rootPassword := ctx.Request().Input("mysql_root_password")
	if len(rootPassword) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "MySQL root密码不能为空")
		return
	}

	oldRootPassword := r.setting.Get("mysql_root_password")
	if oldRootPassword != rootPassword {
		helpers.ExecShell("mysql -uroot -p" + oldRootPassword + " -e \"ALTER USER 'root'@'localhost' IDENTIFIED BY '" + rootPassword + "';\"")
		helpers.ExecShell("mysql -uroot -p" + oldRootPassword + " -e \"FLUSH PRIVILEGES;\"")
		err := r.setting.Set("mysql_root_password", rootPassword)
		if err != nil {
			helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"ALTER USER 'root'@'localhost' IDENTIFIED BY '" + oldRootPassword + "';\"")
			helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\"")
			controllers.Error(ctx, http.StatusInternalServerError, "设置root密码失败")
			return
		}
	}

	controllers.Success(ctx, "设置root密码成功")
}

// DatabaseList 获取数据库列表
func (r *Mysql80Controller) DatabaseList(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	out := helpers.ExecShell("mysql -uroot -p" + r.setting.Get("mysql_root_password") + " -e \"show databases;\"")
	databases := strings.Split(out, "\n")

	databases = databases[1 : len(databases)-1]
	systemDatabases := []string{"information_schema", "mysql", "performance_schema", "sys"}

	var userDatabases []string
	for _, db := range databases {
		if !slices.Contains(systemDatabases, db) {
			userDatabases = append(userDatabases, db)
		}
	}

	type Database struct {
		Name string
	}

	var dbStructs []Database
	for _, db := range userDatabases {
		dbStructs = append(dbStructs, Database{Name: db})
	}

	controllers.Success(ctx, dbStructs)
}

// AddDatabase 添加数据库
func (r *Mysql80Controller) AddDatabase(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255|regex:^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*(_|[^\\w])).+$",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	rootPassword := r.setting.Get("mysql_root_password")
	database := ctx.Request().Input("database")
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")

	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"CREATE DATABASE IF NOT EXISTS " + database + " DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;\"")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"CREATE USER '" + user + "'@'localhost' IDENTIFIED BY '" + password + "';\"")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"GRANT ALL PRIVILEGES ON " + database + ".* TO '" + user + "'@'localhost';\"")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\"")

	controllers.Success(ctx, "添加数据库成功")
}

// DeleteDatabase 删除数据库
func (r *Mysql80Controller) DeleteDatabase(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
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
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	rootPassword := r.setting.Get("mysql_root_password")
	database := ctx.Request().Input("database")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"DROP DATABASE IF EXISTS " + database + ";\"")

	controllers.Success(ctx, "删除数据库成功")
}

// BackupList 获取备份列表
func (r *Mysql80Controller) BackupList(ctx http.Context) {
	backupPath := "/www/backup/mysql"

	if !helpers.Exists(backupPath) {
		helpers.Mkdir(backupPath, 0644)
	}

	files, err := os.ReadDir(backupPath)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, "获取备份列表失败")
		return
	}

	var backupFiles []map[string]string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		backupFiles = append(backupFiles, map[string]string{
			"file": file.Name(),
			"size": helpers.FormatBytes(float64(info.Size())),
		})
	}

	controllers.Success(ctx, backupFiles)
}

// CreateBackup 创建备份
func (r *Mysql80Controller) CreateBackup(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
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
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	backupPath := "/www/backup/mysql"
	rootPassword := r.setting.Get("mysql_root_password")
	database := ctx.Request().Input("database")
	backupFile := backupPath + "/" + database + "_" + carbon.Now().ToShortDateTimeString() + ".sql"
	if !helpers.Exists(backupPath) {
		helpers.Mkdir(backupPath, 0644)
	}
	err = os.Setenv("MYSQL_PWD", rootPassword)
	if err != nil {
		facades.Log().Error("[MySQL80] 设置环境变量 MYSQL_PWD 失败：" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "备份失败")
		return
	}

	helpers.ExecShell("mysqldump -uroot " + database + " > " + backupFile)
	helpers.ExecShell("zip -r " + backupFile + ".zip " + backupFile)
	helpers.RemoveFile(backupFile)
	_ = os.Unsetenv("MYSQL_PWD")

	controllers.Success(ctx, "备份成功")
}

// DeleteBackup 删除备份
func (r *Mysql80Controller) DeleteBackup(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"file": "required|min_len:1|max_len:255",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	backupPath := "/www/backup/mysql"
	file := ctx.Request().Input("file")
	helpers.RemoveFile(backupPath + "/" + file)

	controllers.Success(ctx, "删除备份成功")
}

// RestoreBackup 还原备份
func (r *Mysql80Controller) RestoreBackup(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"file":     "required|min_len:1|max_len:255",
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$|not_in:information_schema,mysql,performance_schema,sys",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	backupPath := "/www/backup/mysql"
	rootPassword := r.setting.Get("mysql_root_password")
	file := ctx.Request().Input("file")
	backupFile := backupPath + "/" + file
	if !helpers.Exists(backupFile) {
		controllers.Error(ctx, http.StatusBadRequest, "备份文件不存在")
		return
	}

	err = os.Setenv("MYSQL_PWD", rootPassword)
	if err != nil {
		facades.Log().Error("[MySQL80] 设置环境变量 MYSQL_PWD 失败：" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "还原失败")
		return
	}

	// 获取文件拓展名
	ext := filepath.Ext(file)
	switch ext {
	case ".zip":
		helpers.ExecShell("unzip -o " + backupFile + " -d " + backupPath)
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".gz":
		if strings.HasSuffix(file, ".tar.gz") {
			// 解压.tar.gz文件
			helpers.ExecShell("tar -zxvf " + backupFile + " -C " + backupPath)
			backupFile = strings.TrimSuffix(backupFile, ".tar.gz")
		} else {
			// 解压.gz文件
			helpers.ExecShell("gzip -d " + backupFile)
			backupFile = strings.TrimSuffix(backupFile, ext)
		}
	case ".bz2":
		helpers.ExecShell("bzip2 -d " + backupFile)
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".tar":
		helpers.ExecShell("tar -xvf " + backupFile + " -C " + backupPath)
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".rar":
		helpers.ExecShell("unrar x " + backupFile + " " + backupPath)
		backupFile = strings.TrimSuffix(backupFile, ext)
	}

	if !helpers.Exists(backupFile) {
		controllers.Error(ctx, http.StatusBadRequest, "自动解压备份文件失败，请手动解压")
		return
	}

	helpers.ExecShell("mysql -uroot " + ctx.Request().Input("database") + " < " + backupFile)
	_ = os.Unsetenv("MYSQL_PWD")

	controllers.Success(ctx, "还原成功")
}

// UserList 用户列表
func (r *Mysql80Controller) UserList(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	type User struct {
		Username   string `json:"username"`
		Host       string `json:"host"`
		Privileges string `json:"privileges"`
	}

	rootPassword := r.setting.Get("mysql_root_password")
	out := helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e 'select user,host from mysql.user'")
	rawUsers := strings.Split(out, "\n")
	users := make([]User, 0)
	for _, rawUser := range rawUsers {
		user := strings.Split(rawUser, "\t")
		if user[0] == "root" || user[0] == "mysql.sys" || user[0] == "mysql.infoschema" || user[0] == "mysql.session" {
			continue
		}

		out := helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e 'show grants for " + user[0] + "@" + user[1] + "'")
		rawPrivileges := strings.Split(out, "\n")
		privileges := make([]string, 0)
		for _, rawPrivilege := range rawPrivileges {
			if rawPrivilege == "" {
				continue
			}
			privilege := rawPrivilege[6:strings.Index(rawPrivilege, " TO")]
			privileges = append(privileges, privilege)
		}
		users = append(users, User{Username: user[0], Host: user[1], Privileges: strings.Join(privileges, " | ")})
	}

	controllers.Success(ctx, users)
}

// AddUser 添加用户
func (r *Mysql80Controller) AddUser(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255|regex:^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*(_|[^\\w])).+$",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	rootPassword := r.setting.Get("mysql_root_password")
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	database := ctx.Request().Input("database")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"CREATE USER '" + user + "'@'localhost' IDENTIFIED BY '" + password + ";'\"")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"GRANT ALL PRIVILEGES ON " + database + ".* TO '" + user + "'@'localhost';\"")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\"")

	controllers.Success(ctx, "添加成功")
}

// DeleteUser 删除用户
func (r *Mysql80Controller) DeleteUser(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
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
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	rootPassword := r.setting.Get("mysql_root_password")
	user := ctx.Request().Input("user")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"DROP USER '" + user + "'@'localhost';\"")

	controllers.Success(ctx, "删除成功")
}

// SetUserPassword 设置用户密码
func (r *Mysql80Controller) SetUserPassword(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255|regex:^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*(_|[^\\w])).+$",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	rootPassword := r.setting.Get("mysql_root_password")
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"ALTER USER '" + user + "'@'localhost' IDENTIFIED BY '" + password + "';\"")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\"")

	controllers.Success(ctx, "修改成功")
}

// SetUserPrivileges 设置用户权限
func (r *Mysql80Controller) SetUserPrivileges(ctx http.Context) {
	if !plugins.Check(ctx, "mysql80") {
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
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	rootPassword := r.setting.Get("mysql_root_password")
	user := ctx.Request().Input("user")
	database := ctx.Request().Input("database")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"REVOKE ALL PRIVILEGES ON *.* FROM '" + user + "'@'localhost';\"")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"GRANT ALL PRIVILEGES ON " + database + ".* TO '" + user + "'@'localhost';\"")
	helpers.ExecShell("mysql -uroot -p" + rootPassword + " -e \"FLUSH PRIVILEGES;\"")

	controllers.Success(ctx, "修改成功")
}
