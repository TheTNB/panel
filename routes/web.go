package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"panel/app/http/controllers"
	"panel/app/http/middleware"
)

func Web() {
	facades.Route().StaticFile("favicon.ico", "/www/panel/public/favicon.ico")
	facades.Route().Prefix("api/panel").Group(func(r route.Route) {
		r.Prefix("info").Group(func(r route.Route) {
			infoController := controllers.NewInfoController()
			r.Get("name", infoController.Name)
			r.Middleware(middleware.Jwt()).Get("menu", infoController.Menu)
			r.Middleware(middleware.Jwt()).Get("homePlugins", infoController.HomePlugins)
			r.Middleware(middleware.Jwt()).Get("nowMonitor", infoController.NowMonitor)
			r.Middleware(middleware.Jwt()).Get("systemInfo", infoController.SystemInfo)
			r.Middleware(middleware.Jwt()).Get("installedDbAndPhp", infoController.InstalledDbAndPhp)
			r.Middleware(middleware.Jwt()).Get("checkUpdate", infoController.CheckUpdate)
			r.Middleware(middleware.Jwt()).Post("update", infoController.Update)
		})
		r.Prefix("user").Group(func(r route.Route) {
			userController := controllers.NewUserController()
			r.Post("login", userController.Login)
			r.Middleware(middleware.Jwt()).Get("info", userController.Info)
		})
		r.Prefix("task").Middleware(middleware.Jwt()).Group(func(r route.Route) {
			taskController := controllers.NewTaskController()
			r.Get("status", taskController.Status)
			r.Get("list", taskController.List)
			r.Get("log", taskController.Log)
			r.Post("delete", taskController.Delete)
		})
		r.Prefix("website").Middleware(middleware.Jwt()).Group(func(r route.Route) {
			websiteController := controllers.NewWebsiteController()
			r.Get("list", websiteController.List)
			r.Post("add", websiteController.Add)
			r.Post("delete", websiteController.Delete)
			r.Get("defaultConfig", websiteController.GetDefaultConfig)
			r.Post("defaultConfig", websiteController.SaveDefaultConfig)
			r.Get("config", websiteController.GetConfig)
			r.Post("config", websiteController.SaveConfig)
			r.Get("clearLog", websiteController.ClearLog)
			r.Post("updateRemark", websiteController.UpdateRemark)
			r.Get("backupList", websiteController.BackupList)
			r.Post("createBackup", websiteController.CreateBackup)
			r.Post("resetConfig", websiteController.ResetConfig)
			r.Post("status", websiteController.Status)
		})
		r.Prefix("plugin").Middleware(middleware.Jwt()).Group(func(r route.Route) {
			pluginController := controllers.NewPluginController()
			r.Get("list", pluginController.List)
			r.Post("install", pluginController.Install)
			r.Post("uninstall", pluginController.Uninstall)
			r.Post("update", pluginController.Update)
			r.Post("updateShow", pluginController.UpdateShow)
		})
		r.Prefix("cron").Middleware(middleware.Jwt()).Group(func(r route.Route) {
			cronController := controllers.NewCronController()
			r.Get("list", cronController.List)
			r.Get("script", cronController.Script)
			r.Post("add", cronController.Add)
			r.Post("update", cronController.Update)
			r.Post("delete", cronController.Delete)
			r.Post("status", cronController.Status)
			r.Get("log", cronController.Log)
		})
		r.Prefix("safe").Middleware(middleware.Jwt()).Group(func(r route.Route) {
			safeController := controllers.NewSafeController()
			r.Get("firewallStatus", safeController.GetFirewallStatus)
			r.Post("firewallStatus", safeController.SetFirewallStatus)
			r.Get("firewallRules", safeController.GetFirewallRules)
			r.Post("addFirewallRule", safeController.AddFirewallRule)
			r.Post("deleteFirewallRule", safeController.DeleteFirewallRule)
			r.Get("sshStatus", safeController.GetSshStatus)
			r.Post("sshStatus", safeController.SetSshStatus)
			r.Get("sshPort", safeController.GetSshPort)
			r.Post("sshPort", safeController.SetSshPort)
			r.Get("pingStatus", safeController.GetPingStatus)
			r.Post("pingStatus", safeController.SetPingStatus)
		})
		r.Prefix("monitor").Middleware(middleware.Jwt()).Group(func(r route.Route) {
			monitorController := controllers.NewMonitorController()
			r.Post("switch", monitorController.Switch)
			r.Post("saveDays", monitorController.SaveDays)
			r.Post("clear", monitorController.Clear)
			r.Get("list", monitorController.List)
			r.Get("switchAndDays", monitorController.SwitchAndDays)
		})
		r.Prefix("setting").Middleware(middleware.Jwt()).Group(func(r route.Route) {
			settingController := controllers.NewSettingController()
			r.Get("list", settingController.List)
			r.Post("save", settingController.Save)
		})
	})

	facades.Route().Fallback(func(ctx http.Context) {
		ctx.Response().String(404, "not found")
	})
}
