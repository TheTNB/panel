package plugins

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

type Postgresql16Controller struct {
	setting services.Setting
	backup  services.Backup
}

func NewPostgresql16Controller() *Postgresql16Controller {
	return &Postgresql16Controller{
		setting: services.NewSettingImpl(),
		backup:  services.NewBackupImpl(),
	}
}

// Status 获取运行状态
func (r *Postgresql16Controller) Status(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Reload 重载配置
func (r *Postgresql16Controller) Reload(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	tools.Exec("systemctl reload postgresql")
	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Restart 重启服务
func (r *Postgresql16Controller) Restart(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	tools.Exec("systemctl restart postgresql")
	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Start 启动服务
func (r *Postgresql16Controller) Start(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	tools.Exec("systemctl start postgresql")
	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Stop 停止服务
func (r *Postgresql16Controller) Stop(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	tools.Exec("systemctl stop postgresql")
	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
	}

	if status != "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// GetConfig 获取配置
func (r *Postgresql16Controller) GetConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	// 获取配置
	config := tools.Read("/www/server/postgresql/data/postgresql.conf")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL配置失败")
	}

	return controllers.Success(ctx, config)
}

// GetUserConfig 获取用户配置
func (r *Postgresql16Controller) GetUserConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	// 获取配置
	config := tools.Read("/www/server/postgresql/data/pg_hba.conf")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL配置失败")
	}

	return controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (r *Postgresql16Controller) SaveConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "配置不能为空")
	}

	if !tools.Write("/www/server/postgresql/data/postgresql.conf", config, 0644) {
		return controllers.Error(ctx, http.StatusInternalServerError, "写入PostgreSQL配置失败")
	}

	return r.Restart(ctx)
}

// SaveUserConfig 保存用户配置
func (r *Postgresql16Controller) SaveUserConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "配置不能为空")
	}

	if !tools.Write("/www/server/postgresql/data/pg_hba.conf", config, 0644) {
		return controllers.Error(ctx, http.StatusInternalServerError, "写入PostgreSQL配置失败")
	}

	return r.Restart(ctx)
}

// Load 获取负载
func (r *Postgresql16Controller) Load(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if status != "active" {
		return controllers.Error(ctx, http.StatusInternalServerError, "PostgreSQL 已停止运行")
	}

	data := []LoadInfo{
		{"启动时间", carbon.Parse(tools.Exec(`echo "select pg_postmaster_start_time();" | su - postgres -c "psql" | sed -n 3p | cut -d'.' -f1`)).ToDateTimeString()},
		{"进程 PID", tools.Exec(`echo "select pg_backend_pid();" | su - postgres -c "psql" | sed -n 3p`)},
		{"进程数", tools.Exec(`ps aux | grep postgres | grep -v grep | wc -l`)},
		{"总连接数", tools.Exec(`echo "SELECT count(*) FROM pg_stat_activity WHERE NOT pid=pg_backend_pid();" | su - postgres -c "psql" | sed -n 3p`)},
		{"空间占用", tools.Exec(`echo "select pg_size_pretty(pg_database_size('postgres'));" | su - postgres -c "psql" | sed -n 3p`)},
	}

	return controllers.Success(ctx, data)
}

// Log 获取日志
func (r *Postgresql16Controller) Log(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	log := tools.Exec("tail -n 100 /www/server/postgresql/logs/postgresql-" + carbon.Now().ToDateString() + ".log")
	return controllers.Success(ctx, log)
}

// ClearLog 清空日志
func (r *Postgresql16Controller) ClearLog(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	tools.Exec("echo '' > /www/server/postgresql/logs/postgresql-" + carbon.Now().ToDateString() + ".log")
	return controllers.Success(ctx, nil)
}

// DatabaseList 获取数据库列表
func (r *Postgresql16Controller) DatabaseList(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	status := tools.Exec("systemctl status postgresql | grep Active | grep -v grep | awk '{print $2}'")
	if status != "active" {
		return controllers.Error(ctx, http.StatusInternalServerError, "PostgreSQL 已停止运行")
	}

	type database struct {
		Name     string `json:"name"`
		Owner    string `json:"owner"`
		Encoding string `json:"encoding"`
	}

	raw := tools.Exec(`echo "\l" | su - postgres -c "psql"`)
	databases := strings.Split(raw, "\n")
	if len(databases) >= 4 {
		databases = databases[3 : len(databases)-1]
	} else {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []database{},
		})
	}

	var databaseList []database
	for _, db := range databases {
		parts := strings.Split(db, "|")
		if len(parts) != 9 || len(strings.TrimSpace(parts[0])) == 0 {
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
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []database{},
		})
	}
	if endIndex > len(databaseList) {
		endIndex = len(databaseList)
	}
	pagedDatabases := databaseList[startIndex:endIndex]

	return controllers.Success(ctx, http.Json{
		"total": len(databaseList),
		"items": pagedDatabases,
	})
}

// AddDatabase 添加数据库
func (r *Postgresql16Controller) AddDatabase(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	database := ctx.Request().Input("database")
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")

	tools.Exec(`echo "CREATE DATABASE ` + database + `;" | su - postgres -c "psql"`)
	tools.Exec(`echo "CREATE USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`)
	tools.Exec(`echo "ALTER DATABASE ` + database + ` OWNER TO ` + user + `;" | su - postgres -c "psql"`)
	tools.Exec(`echo "GRANT ALL PRIVILEGES ON DATABASE ` + database + ` TO ` + user + `;" | su - postgres -c "psql"`)

	userConfig := "host    " + database + "    " + user + "    127.0.0.1/32    scram-sha-256"
	tools.Exec(`echo "` + userConfig + `" >> /www/server/postgresql/data/pg_hba.conf`)

	return r.Reload(ctx)
}

// DeleteDatabase 删除数据库
func (r *Postgresql16Controller) DeleteDatabase(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$|not_in:postgres,template0,template1",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	database := ctx.Request().Input("database")
	tools.Exec(`echo "DROP DATABASE ` + database + `;" | su - postgres -c "psql"`)

	return controllers.Success(ctx, nil)
}

// BackupList 获取备份列表
func (r *Postgresql16Controller) BackupList(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	backupList, err := r.backup.PostgresqlList()
	if err != nil {
		facades.Log().Error("[PostgreSQL] 获取备份列表失败：" + err.Error())
		return controllers.Error(ctx, http.StatusInternalServerError, "获取备份列表失败")
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)
	startIndex := (page - 1) * limit
	endIndex := page * limit
	if startIndex > len(backupList) {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []services.BackupFile{},
		})
	}
	if endIndex > len(backupList) {
		endIndex = len(backupList)
	}
	pagedBackupList := backupList[startIndex:endIndex]
	if pagedBackupList == nil {
		pagedBackupList = []services.BackupFile{}
	}

	return controllers.Success(ctx, http.Json{
		"total": len(backupList),
		"items": pagedBackupList,
	})
}

// UploadBackup 上传备份
func (r *Postgresql16Controller) UploadBackup(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	file, err := ctx.Request().File("file")
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "上传文件失败")
	}

	backupPath := r.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}

	name := file.GetClientOriginalName()
	_, err = file.StoreAs(backupPath, name)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "上传文件失败")
	}

	return controllers.Success(ctx, nil)
}

// CreateBackup 创建备份
func (r *Postgresql16Controller) CreateBackup(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$|not_in:information_schema,mysql,performance_schema,sys",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	database := ctx.Request().Input("database")
	err = r.backup.PostgresqlBackup(database)
	if err != nil {
		facades.Log().Error("[PostgreSQL] 创建备份失败：" + err.Error())
		return controllers.Error(ctx, http.StatusInternalServerError, "创建备份失败")
	}

	return controllers.Success(ctx, nil)
}

// DeleteBackup 删除备份
func (r *Postgresql16Controller) DeleteBackup(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"name": "required|min_len:1|max_len:255",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	backupPath := r.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
	fileName := ctx.Request().Input("name")
	tools.Remove(backupPath + "/" + fileName)

	return controllers.Success(ctx, nil)
}

// RestoreBackup 还原备份
func (r *Postgresql16Controller) RestoreBackup(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"backup":   "required|min_len:1|max_len:255",
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$|not_in:information_schema,mysql,performance_schema,sys",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	err = r.backup.PostgresqlRestore(ctx.Request().Input("database"), ctx.Request().Input("backup"))
	if err != nil {
		facades.Log().Error("[PostgreSQL] 还原失败：" + err.Error())
		return controllers.Error(ctx, http.StatusInternalServerError, "还原失败: "+err.Error())
	}

	return controllers.Success(ctx, nil)
}

// UserList 用户列表
func (r *Postgresql16Controller) UserList(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	type user struct {
		User string `json:"user"`
		Role string `json:"role"`
	}

	raw := tools.Exec(`echo "\du" | su - postgres -c "psql"`)
	users := strings.Split(raw, "\n")
	if len(users) < 4 {
		return controllers.Error(ctx, http.StatusInternalServerError, "用户列表为空")
	}
	users = users[3:]

	var userList []user
	for _, u := range users {
		userInfo := strings.Split(u, "|")
		if len(userInfo) != 2 {
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
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []user{},
		})
	}
	if endIndex > len(userList) {
		endIndex = len(userList)
	}
	pagedUsers := userList[startIndex:endIndex]

	return controllers.Success(ctx, http.Json{
		"total": len(userList),
		"items": pagedUsers,
	})
}

// AddUser 添加用户
func (r *Postgresql16Controller) AddUser(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"database": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	database := ctx.Request().Input("database")
	tools.Exec(`echo "CREATE USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`)
	tools.Exec(`echo "GRANT ALL PRIVILEGES ON DATABASE ` + database + ` TO ` + user + `;" | su - postgres -c "psql"`)

	userConfig := "host    " + database + "    " + user + "    127.0.0.1/32    scram-sha-256"
	tools.Exec(`echo "` + userConfig + `" >> /www/server/postgresql/data/pg_hba.conf`)

	return r.Reload(ctx)
}

// DeleteUser 删除用户
func (r *Postgresql16Controller) DeleteUser(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"user": "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	user := ctx.Request().Input("user")
	tools.Exec(`echo "DROP USER ` + user + `;" | su - postgres -c "psql"`)
	tools.Exec(`sed -i '/` + user + `/d' /www/server/postgresql/data/pg_hba.conf`)

	return r.Reload(ctx)
}

// SetUserPassword 设置用户密码
func (r *Postgresql16Controller) SetUserPassword(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql16")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"user":     "required|min_len:1|max_len:255|regex:^[a-zA-Z][a-zA-Z0-9_]+$",
		"password": "required|min_len:8|max_len:255",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	tools.Exec(`echo "ALTER USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`)

	return controllers.Success(ctx, nil)
}
