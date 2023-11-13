package plugins

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
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

func NewPostgresql15Controller() *Postgresql15Controller {
	return &Postgresql15Controller{
		setting: services.NewSettingImpl(),
		backup:  services.NewBackupImpl(),
	}
}

// Status 获取运行状态
func (r *Postgresql15Controller) Status(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	status, err := tools.ServiceStatus("postgresql")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
	}

	return controllers.Success(ctx, status)
}

// Reload 重载配置
func (r *Postgresql15Controller) Reload(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	if err := tools.ServiceReload("postgresql"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "重载PostgreSQL失败")
	}

	return controllers.Success(ctx, nil)
}

// Restart 重启服务
func (r *Postgresql15Controller) Restart(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	if err := tools.ServiceRestart("postgresql"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "重启PostgreSQL失败")
	}

	return controllers.Success(ctx, nil)
}

// Start 启动服务
func (r *Postgresql15Controller) Start(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	if err := tools.ServiceStart("postgresql"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "启动PostgreSQL失败")
	}

	return controllers.Success(ctx, nil)
}

// Stop 停止服务
func (r *Postgresql15Controller) Stop(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	if err := tools.ServiceStop("postgresql"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "停止PostgreSQL失败")
	}

	return controllers.Success(ctx, nil)
}

// GetConfig 获取配置
func (r *Postgresql15Controller) GetConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
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
func (r *Postgresql15Controller) GetUserConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
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
func (r *Postgresql15Controller) SaveConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "配置不能为空")
	}

	if err := tools.Write("/www/server/postgresql/data/postgresql.conf", config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "写入PostgreSQL配置失败")
	}

	return r.Restart(ctx)
}

// SaveUserConfig 保存用户配置
func (r *Postgresql15Controller) SaveUserConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	config := ctx.Request().Input("config")
	if len(config) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "配置不能为空")
	}

	if err := tools.Write("/www/server/postgresql/data/pg_hba.conf", config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "写入PostgreSQL配置失败")
	}

	return r.Restart(ctx)
}

// Load 获取负载
func (r *Postgresql15Controller) Load(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	status, err := tools.ServiceStatus("postgresql")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
	}
	if !status {
		return controllers.Error(ctx, http.StatusInternalServerError, "PostgreSQL已停止运行")
	}

	time, err := tools.Exec(`echo "select pg_postmaster_start_time();" | su - postgres -c "psql" | sed -n 3p | cut -d'.' -f1`)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL启动时间失败")
	}
	pid, err := tools.Exec(`echo "select pg_backend_pid();" | su - postgres -c "psql" | sed -n 3p`)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL进程PID失败")
	}
	process, err := tools.Exec(`ps aux | grep postgres | grep -v grep | wc -l`)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL进程数失败")
	}
	connections, err := tools.Exec(`echo "SELECT count(*) FROM pg_stat_activity WHERE NOT pid=pg_backend_pid();" | su - postgres -c "psql" | sed -n 3p`)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL连接数失败")
	}
	storage, err := tools.Exec(`echo "select pg_size_pretty(pg_database_size('postgres'));" | su - postgres -c "psql" | sed -n 3p`)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL空间占用失败")
	}

	data := []LoadInfo{
		{"启动时间", carbon.Parse(time).ToDateTimeString()},
		{"进程 PID", pid},
		{"进程数", process},
		{"总连接数", connections},
		{"空间占用", storage},
	}

	return controllers.Success(ctx, data)
}

// Log 获取日志
func (r *Postgresql15Controller) Log(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	log, err := tools.Exec("tail -n 100 /www/server/postgresql/logs/postgresql-" + carbon.Now().ToDateString() + ".log")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, log)
	}

	return controllers.Success(ctx, log)
}

// ClearLog 清空日志
func (r *Postgresql15Controller) ClearLog(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	if out, err := tools.Exec("echo '' > /www/server/postgresql/logs/postgresql-" + carbon.Now().ToDateString() + ".log"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// DatabaseList 获取数据库列表
func (r *Postgresql15Controller) DatabaseList(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	status, err := tools.ServiceStatus("postgresql")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PostgreSQL状态失败")
	}
	if !status {
		return controllers.Error(ctx, http.StatusInternalServerError, "PostgreSQL已停止运行")
	}

	type database struct {
		Name     string `json:"name"`
		Owner    string `json:"owner"`
		Encoding string `json:"encoding"`
	}

	raw, err := tools.Exec(`echo "\l" | su - postgres -c "psql"`)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, raw)
	}
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
func (r *Postgresql15Controller) AddDatabase(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
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

	if out, err := tools.Exec(`echo "CREATE DATABASE ` + database + `;" | su - postgres -c "psql"`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := tools.Exec(`echo "CREATE USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := tools.Exec(`echo "ALTER DATABASE ` + database + ` OWNER TO ` + user + `;" | su - postgres -c "psql"`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := tools.Exec(`echo "GRANT ALL PRIVILEGES ON DATABASE ` + database + ` TO ` + user + `;" | su - postgres -c "psql"`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	userConfig := "host    " + database + "    " + user + "    127.0.0.1/32    scram-sha-256"
	if out, err := tools.Exec(`echo "` + userConfig + `" >> /www/server/postgresql/data/pg_hba.conf`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return r.Reload(ctx)
}

// DeleteDatabase 删除数据库
func (r *Postgresql15Controller) DeleteDatabase(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
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
	if out, err := tools.Exec(`echo "DROP DATABASE ` + database + `;" | su - postgres -c "psql"`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// BackupList 获取备份列表
func (r *Postgresql15Controller) BackupList(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	backupList, err := r.backup.PostgresqlList()
	if err != nil {
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
func (r *Postgresql15Controller) UploadBackup(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	file, err := ctx.Request().File("file")
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "上传文件失败")
	}

	backupPath := r.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
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
func (r *Postgresql15Controller) CreateBackup(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
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
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

// DeleteBackup 删除备份
func (r *Postgresql15Controller) DeleteBackup(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
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
func (r *Postgresql15Controller) RestoreBackup(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
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
		return controllers.Error(ctx, http.StatusInternalServerError, "还原失败: "+err.Error())
	}

	return controllers.Success(ctx, nil)
}

// UserList 用户列表
func (r *Postgresql15Controller) UserList(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
	if check != nil {
		return check
	}

	type user struct {
		User string `json:"user"`
		Role string `json:"role"`
	}

	raw, err := tools.Exec(`echo "\du" | su - postgres -c "psql"`)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, raw)
	}
	users := strings.Split(raw, "\n")
	if len(users) < 4 {
		return controllers.Error(ctx, http.StatusInternalServerError, "用户列表为空")
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
func (r *Postgresql15Controller) AddUser(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
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
	if out, err := tools.Exec(`echo "CREATE USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := tools.Exec(`echo "GRANT ALL PRIVILEGES ON DATABASE ` + database + ` TO ` + user + `;" | su - postgres -c "psql"`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	userConfig := "host    " + database + "    " + user + "    127.0.0.1/32    scram-sha-256"
	if out, err := tools.Exec(`echo "` + userConfig + `" >> /www/server/postgresql/data/pg_hba.conf`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return r.Reload(ctx)
}

// DeleteUser 删除用户
func (r *Postgresql15Controller) DeleteUser(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
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
	if out, err := tools.Exec(`echo "DROP USER ` + user + `;" | su - postgres -c "psql"`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := tools.Exec(`sed -i '/` + user + `/d' /www/server/postgresql/data/pg_hba.conf`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return r.Reload(ctx)
}

// SetUserPassword 设置用户密码
func (r *Postgresql15Controller) SetUserPassword(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "postgresql15")
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
	if out, err := tools.Exec(`echo "ALTER USER ` + user + ` WITH PASSWORD '` + password + `';" | su - postgres -c "psql"`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}
