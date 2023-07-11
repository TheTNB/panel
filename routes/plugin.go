package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"panel/app/http/controllers/plugins/mysql80"
	"panel/app/http/controllers/plugins/openresty"
	"panel/app/http/middleware"
)

// Plugin 加载插件路由
func Plugin() {
	facades.Route().Prefix("api/plugins/openresty").Middleware(middleware.Jwt()).Group(func(route route.Route) {
		openRestyController := openresty.NewOpenrestyController()
		route.Get("status", openRestyController.Status)
		route.Post("reload", openRestyController.Reload)
		route.Post("start", openRestyController.Start)
		route.Post("stop", openRestyController.Stop)
		route.Post("restart", openRestyController.Restart)
		route.Get("load", openRestyController.Load)
		route.Get("config", openRestyController.GetConfig)
		route.Post("config", openRestyController.SaveConfig)
		route.Get("errorLog", openRestyController.ErrorLog)
		route.Get("clearErrorLog", openRestyController.ClearErrorLog)
	})
	facades.Route().Prefix("api/plugins/mysql80").Middleware(middleware.Jwt()).Group(func(route route.Route) {
		mysql80Controller := mysql80.NewMysql80Controller()
		route.Get("status", mysql80Controller.Status)
		route.Post("reload", mysql80Controller.Reload)
		route.Post("start", mysql80Controller.Start)
		route.Post("stop", mysql80Controller.Stop)
		route.Post("restart", mysql80Controller.Restart)
		route.Get("load", mysql80Controller.Load)
		route.Get("config", mysql80Controller.GetConfig)
		route.Post("config", mysql80Controller.SaveConfig)
		route.Get("errorLog", mysql80Controller.ErrorLog)
		route.Get("clearErrorLog", mysql80Controller.ClearErrorLog)
		route.Get("slowLog", mysql80Controller.SlowLog)
		route.Get("clearSlowLog", mysql80Controller.ClearSlowLog)
		route.Get("rootPassword", mysql80Controller.GetRootPassword)
		route.Post("rootPassword", mysql80Controller.SetRootPassword)
		route.Get("database", mysql80Controller.DatabaseList)
		route.Post("addDatabase", mysql80Controller.AddDatabase)
		route.Post("deleteDatabase", mysql80Controller.DeleteDatabase)
		route.Get("backup", mysql80Controller.BackupList)
		route.Post("createBackup", mysql80Controller.CreateBackup)
		route.Post("deleteBackup", mysql80Controller.DeleteBackup)
		route.Post("restoreBackup", mysql80Controller.RestoreBackup)
		route.Get("user", mysql80Controller.UserList)
		route.Post("addUser", mysql80Controller.AddUser)
		route.Post("deleteUser", mysql80Controller.DeleteUser)
		route.Post("setUserPassword", mysql80Controller.SetUserPassword)
		route.Post("setUserPrivileges", mysql80Controller.SetUserPrivileges)
	})
}
