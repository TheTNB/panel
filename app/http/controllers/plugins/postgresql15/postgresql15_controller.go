package postgresql15

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	"panel/app/http/controllers"
	"panel/app/models"
	"panel/app/services"
	"panel/pkg/tools"
)

type Postgresql15Controller struct {
	setting services.Setting
	backup  services.Backup
}

type Info struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewPostgresql15Controller() *Postgresql15Controller {
	return &Postgresql15Controller{
		setting: services.NewSettingImpl(),
		backup:  services.NewBackupImpl(),
	}
}

// Status 获取运行状态
func (c *Postgresql15Controller) Status(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Reload 重载配置
func (c *Postgresql15Controller) Reload(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	tools.Exec("systemctl reload postgresql")
	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Restart 重启服务
func (c *Postgresql15Controller) Restart(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	tools.Exec("systemctl restart postgresql")
	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Start 启动服务
func (c *Postgresql15Controller) Start(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	tools.Exec("systemctl start postgresql")
	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Stop 停止服务
func (c *Postgresql15Controller) Stop(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	tools.Exec("systemctl stop postgresql")
	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
		return
	}

	if status != "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// GetConfig 获取配置
func (c *Postgresql15Controller) GetConfig(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	// 获取配置
	config := tools.Read("/www/server/postgresql/data/postgresql.conf")
	if len(config) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL配置失败")
		return
	}

	controllers.Success(ctx, config)
}

// GetUserConfig 获取用户配置
func (c *Postgresql15Controller) GetUserConfig(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	// 获取配置
	config := tools.Read("/www/server/postgresql/data/pg_hba.conf")
	if len(config) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL配置失败")
		return
	}

	controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (c *Postgresql15Controller) SaveConfig(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	config := ctx.Request().Input("config")
	if len(config) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "配置不能为空")
		return
	}

	if !tools.Write("/www/server/postgresql/data/postgresql.conf", config, 0644) {
		controllers.Error(ctx, http.StatusInternalServerError, "写入PostgreSQL配置失败")
		return
	}

	c.Restart(ctx)
}

// SaveUserConfig 保存用户配置
func (c *Postgresql15Controller) SaveUserConfig(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	config := ctx.Request().Input("config")
	if len(config) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "配置不能为空")
		return
	}

	if !tools.Write("/www/server/postgresql/data/pg_hba.conf", config, 0644) {
		controllers.Error(ctx, http.StatusInternalServerError, "写入PostgreSQL配置失败")
		return
	}

	c.Restart(ctx)
}

// Load 获取负载
func (c *Postgresql15Controller) Load(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if status != "active" {
		controllers.Error(ctx, http.StatusInternalServerError, "PostgreSQL 已停止运行")
		return
	}

	data := []Info{
		{"启动时间", carbon.Parse(tools.Exec(`echo "select pg_postmaster_start_time();" | su - postgres -c "psql" | sed -n 3p | cut -d'.' -f1`)).ToDateTimeString()},
		{"进程 PID", tools.Exec(`echo "select pg_backend_pid();" | su - postgres -c "psql" | sed -n 3p`)},
		{"进程数", tools.Exec(`ps aux | grep postgres | grep -v grep | wc -l`)},
		{"总连接数", tools.Exec(`echo "SELECT count(*) FROM pg_stat_activity WHERE NOT pid=pg_backend_pid();" | su - postgres -c "psql" | sed -n 3p`)},
		{"空间占用", tools.Exec(`echo "select pg_size_pretty(pg_database_size('postgres'));" | su - postgres -c "psql" | sed -n 3p`)},
	}

	controllers.Success(ctx, data)
}

// Log 获取日志
func (c *Postgresql15Controller) Log(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	log := tools.Exec("tail -n 100 /www/server/postgresql/logs/postgresql-" + carbon.Now().ToDateString() + ".log")
	controllers.Success(ctx, log)
}

// ClearLog 清空日志
func (c *Postgresql15Controller) ClearLog(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	tools.Exec("echo '' > /www/server/postgresql/logs/postgresql-" + carbon.Now().ToDateString() + ".log")
	controllers.Success(ctx, nil)
}

// DatabaseList 获取数据库列表
func (c *Postgresql15Controller) DatabaseList(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if status != "active" {
		controllers.Error(ctx, http.StatusInternalServerError, "PostgreSQL 已停止运行")
		return
	}

	raw := tools.Exec(`echo "\l" | su - postgres -c "psql"`)
	databases := strings.Split(raw, "\n")
	databases = databases[3 : len(databases)-1]

	type database struct {
		Name     string `json:"name"`
		Owner    string `json:"owner"`
		Encoding string `json:"encoding"`
	}

	var databaseList []database
	for _, db := range databases {
		parts := strings.Split(db, "|")
		if len(parts) != 8 || len(strings.TrimSpace(parts[0])) == 0 {
			continue
		}

		databaseList = append(databaseList, database{
			Name:     strings.TrimSpace(parts[0]),
			Owner:    strings.TrimSpace(parts[1]),
			Encoding: strings.TrimSpace(parts[2]),
		})
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)
	startIndex := (page - 1) * limit
	endIndex := page * limit
	if startIndex > len(databaseList) {
		controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []database{},
		})
		return
	}
	if endIndex > len(databaseList) {
		endIndex = len(databaseList)
	}
	pagedDatabases := databaseList[startIndex:endIndex]

	controllers.Success(ctx, http.Json{
		"total": len(databaseList),
		"items": pagedDatabases,
	})
}

// AddDatabase 添加数据库
func (c *Postgresql15Controller) AddDatabase(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
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

	database := ctx.Request().Input("database")
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")

	tools.Exec(`echo "CREATE DATABASE ` + database + `;" | su - postgres -c "psql"`)
	tools.Exec(`echo "CREATE USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`)
	tools.Exec(`echo "GRANT ALL PRIVILEGES ON DATABASE ` + database + ` TO ` + user + `;" | su - postgres -c "psql"`)

	userConfig := "host    " + database + "    " + user + "    127.0.0.1/32    scram-sha-256"
	tools.Exec(`echo "` + userConfig + `" >> /www/server/postgresql/data/pg_hba.conf`)

	c.Reload(ctx)
}

// DeleteDatabase 删除数据库
func (c *Postgresql15Controller) DeleteDatabase(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$|not_in:postgres,template0,template1",
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
	tools.Exec(`echo "DROP DATABASE ` + database + `;" | su - postgres -c "psql"`)

	controllers.Success(ctx, nil)
}

// BackupList 获取备份列表
func (c *Postgresql15Controller) BackupList(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	backupList, err := c.backup.PostgresqlList()
	if err != nil {
		facades.Log().Error("[PostgreSQL] 获取备份列表失败：" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取备份列表失败")
		return
	}

	controllers.Success(ctx, backupList)
}

// UploadBackup 上传备份
func (c *Postgresql15Controller) UploadBackup(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	file, err := ctx.Request().File("file")
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "上传文件失败")
		return
	}

	backupPath := c.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}

	name := file.GetClientOriginalName()
	_, err = file.StoreAs(backupPath, name)
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "上传文件失败")
		return
	}

	controllers.Success(ctx, nil)
}

// CreateBackup 创建备份
func (c *Postgresql15Controller) CreateBackup(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
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
	err = c.backup.PostgresqlBackup(database)
	if err != nil {
		facades.Log().Error("[PostgreSQL] 创建备份失败：" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "创建备份失败")
		return
	}

	controllers.Success(ctx, nil)
}

// DeleteBackup 删除备份
func (c *Postgresql15Controller) DeleteBackup(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
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

	backupPath := c.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
	fileName := ctx.Request().Input("name")
	tools.Remove(backupPath + "/" + fileName)

	controllers.Success(ctx, nil)
}

// RestoreBackup 还原备份
func (c *Postgresql15Controller) RestoreBackup(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
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

	err = c.backup.PostgresqlRestore(ctx.Request().Input("database"), ctx.Request().Input("name"))
	if err != nil {
		facades.Log().Error("[PostgreSQL] 还原失败：" + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "还原失败: "+err.Error())
		return
	}

	controllers.Success(ctx, nil)
}

// UserList 用户列表
func (c *Postgresql15Controller) UserList(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
		return
	}

	type user struct {
		User string `json:"user"`
		Role string `json:"role"`
	}

	raw := tools.Exec(`echo "\du" | su - postgres -c "psql"`)
	users := strings.Split(raw, "\n")
	if len(users) < 4 {
		controllers.Error(ctx, http.StatusInternalServerError, "用户列表为空")
		return
	}
	users = users[3:]

	var userList []user
	for _, u := range users {
		userInfo := strings.Split(u, "|")
		if len(userInfo) != 3 {
			continue
		}

		userList = append(userList, user{
			User: strings.TrimSpace(userInfo[0]),
			Role: strings.TrimSpace(userInfo[1]),
		})
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)
	startIndex := (page - 1) * limit
	endIndex := page * limit
	if startIndex > len(userList) {
		controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []user{},
		})
		return
	}
	if endIndex > len(userList) {
		endIndex = len(userList)
	}
	pagedUsers := userList[startIndex:endIndex]

	controllers.Success(ctx, http.Json{
		"total": len(userList),
		"items": pagedUsers,
	})
}

// AddUser 添加用户
func (c *Postgresql15Controller) AddUser(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
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

	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	database := ctx.Request().Input("database")
	tools.Exec(`echo "CREATE USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`)
	tools.Exec(`echo "GRANT ALL PRIVILEGES ON DATABASE ` + database + ` TO ` + user + `;" | su - postgres -c "psql"`)

	userConfig := "host    " + database + "    " + user + "    127.0.0.1/32    scram-sha-256"
	tools.Exec(`echo "` + userConfig + `" >> /www/server/postgresql/data/pg_hba.conf`)

	c.Reload(ctx)
}

// DeleteUser 删除用户
func (c *Postgresql15Controller) DeleteUser(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
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

	user := ctx.Request().Input("user")
	tools.Exec(`echo "DROP USER ` + user + `;" | su - postgres -c "psql"`)
	tools.Exec(`sed -i '/` + user + `/d' /www/server/postgresql/data/pg_hba.conf`)

	c.Reload(ctx)
}

// SetUserPassword 设置用户密码
func (c *Postgresql15Controller) SetUserPassword(ctx http.Context) {
	if !controllers.Check(ctx, "postgresql15") {
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

	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	tools.Exec(`echo "ALTER USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`)

	controllers.Success(ctx, nil)
}
