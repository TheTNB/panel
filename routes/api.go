package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"panel/app/http/controllers"
	"panel/app/http/middleware"
)

func Api() {
	facades.Route().Prefix("api/panel").Group(func(r route.Router) {
		r.Prefix("info").Group(func(r route.Router) {
			infoController := controllers.NewInfoController()
			r.Get("name", infoController.Name)
			r.Middleware(middleware.Jwt()).Get("homePlugins", infoController.HomePlugins)
			r.Middleware(middleware.Jwt()).Get("nowMonitor", infoController.NowMonitor)
			r.Middleware(middleware.Jwt()).Get("systemInfo", infoController.SystemInfo)
			r.Middleware(middleware.Jwt()).Get("countInfo", infoController.CountInfo)
			r.Middleware(middleware.Jwt()).Get("installedDbAndPhp", infoController.InstalledDbAndPhp)
			r.Middleware(middleware.Jwt()).Get("checkUpdate", infoController.CheckUpdate)
			r.Middleware(middleware.Jwt()).Get("updateInfo", infoController.UpdateInfo)
			r.Middleware(middleware.Jwt()).Post("update", infoController.Update)
			r.Middleware(middleware.Jwt()).Post("restart", infoController.Restart)
		})
		r.Prefix("user").Group(func(r route.Router) {
			userController := controllers.NewUserController()
			r.Post("login", userController.Login)
			r.Middleware(middleware.Jwt()).Get("info", userController.Info)
		})
		r.Prefix("task").Middleware(middleware.Jwt()).Group(func(r route.Router) {
			taskController := controllers.NewTaskController()
			r.Get("status", taskController.Status)
			r.Get("list", taskController.List)
			r.Get("log", taskController.Log)
			r.Post("delete", taskController.Delete)
		})
		r.Prefix("website").Middleware(middleware.Jwt()).Group(func(r route.Router) {
			websiteController := controllers.NewWebsiteController()
			r.Get("/", websiteController.List)
			r.Post("/", websiteController.Add)
			r.Delete("{id}", websiteController.Delete)
			r.Get("defaultConfig", websiteController.GetDefaultConfig)
			r.Post("defaultConfig", websiteController.SaveDefaultConfig)
			r.Get("config/{id}", websiteController.GetConfig)
			r.Post("config/{id}", websiteController.SaveConfig)
			r.Delete("log/{id}", websiteController.ClearLog)
			r.Post("updateRemark/{id}", websiteController.UpdateRemark)
			r.Get("backupList", websiteController.BackupList)
			r.Post("createBackup", websiteController.CreateBackup)
			r.Post("uploadBackup", websiteController.UploadBackup)
			r.Post("restoreBackup", websiteController.RestoreBackup)
			r.Post("deleteBackup", websiteController.DeleteBackup)
			r.Post("resetConfig", websiteController.ResetConfig)
			r.Post("status/{id}", websiteController.Status)
		})
		r.Prefix("cert").Middleware(middleware.Jwt()).Group(func(r route.Router) {
			certController := controllers.NewCertController()
			r.Get("caProviders", certController.CAProviders)
			r.Get("dnsProviders", certController.DNSProviders)
			r.Get("algorithms", certController.Algorithms)
			r.Get("users", certController.UserList)
			r.Post("users", certController.UserStore)
			r.Put("users/{id}", certController.UserUpdate)
			r.Get("users/{id}", certController.UserShow)
			r.Delete("users/{id}", certController.UserDestroy)
			r.Get("dns", certController.DNSList)
			r.Post("dns", certController.DNSStore)
			r.Put("dns/{id}", certController.DNSUpdate)
			r.Get("dns/{id}", certController.DNSShow)
			r.Delete("dns/{id}", certController.DNSDestroy)
			r.Get("certs", certController.CertList)
			r.Post("certs", certController.CertStore)
			r.Put("certs/{id}", certController.CertUpdate)
			r.Get("certs/{id}", certController.CertShow)
			r.Delete("certs/{id}", certController.CertDestroy)
			r.Post("obtain", certController.Obtain)
			r.Post("renew", certController.Renew)
			r.Post("manualDNS", certController.ManualDNS)
		})
		r.Prefix("plugin").Middleware(middleware.Jwt()).Group(func(r route.Router) {
			pluginController := controllers.NewPluginController()
			r.Get("list", pluginController.List)
			r.Post("install", pluginController.Install)
			r.Post("uninstall", pluginController.Uninstall)
			r.Post("update", pluginController.Update)
			r.Post("updateShow", pluginController.UpdateShow)
		})
		r.Prefix("cron").Middleware(middleware.Jwt()).Group(func(r route.Router) {
			cronController := controllers.NewCronController()
			r.Get("list", cronController.List)
			r.Get("{id}", cronController.Script)
			r.Post("add", cronController.Add)
			r.Put("{id}", cronController.Update)
			r.Delete("{id}", cronController.Delete)
			r.Post("status", cronController.Status)
			r.Get("log/{id}", cronController.Log)
		})
		r.Prefix("safe").Middleware(middleware.Jwt()).Group(func(r route.Router) {
			safeController := controllers.NewSafeController()
			r.Get("firewallStatus", safeController.GetFirewallStatus)
			r.Post("firewallStatus", safeController.SetFirewallStatus)
			r.Get("firewallRules", safeController.GetFirewallRules)
			r.Post("firewallRules", safeController.AddFirewallRule)
			r.Delete("firewallRules", safeController.DeleteFirewallRule)
			r.Get("sshStatus", safeController.GetSshStatus)
			r.Post("sshStatus", safeController.SetSshStatus)
			r.Get("sshPort", safeController.GetSshPort)
			r.Post("sshPort", safeController.SetSshPort)
			r.Get("pingStatus", safeController.GetPingStatus)
			r.Post("pingStatus", safeController.SetPingStatus)
		})
		r.Prefix("monitor").Middleware(middleware.Jwt()).Group(func(r route.Router) {
			monitorController := controllers.NewMonitorController()
			r.Post("switch", monitorController.Switch)
			r.Post("saveDays", monitorController.SaveDays)
			r.Post("clear", monitorController.Clear)
			r.Get("list", monitorController.List)
			r.Get("switchAndDays", monitorController.SwitchAndDays)
		})
		r.Prefix("ssh").Middleware(middleware.Jwt()).Group(func(r route.Router) {
			sshController := controllers.NewSshController()
			r.Get("info", sshController.GetInfo)
			r.Post("info", sshController.UpdateInfo)
			r.Get("session", sshController.Session)
		})
		r.Prefix("setting").Middleware(middleware.Jwt()).Group(func(r route.Router) {
			settingController := controllers.NewSettingController()
			r.Get("list", settingController.List)
			r.Post("update", settingController.Update)
		})
	})

	// 文档
	swaggerController := controllers.NewSwaggerController()
	facades.Route().Get("swagger", swaggerController.Index)
	facades.Route().Get("swagger/{any}", swaggerController.Index)

	// 静态文件
	entrance := facades.Config().GetString("http.entrance")
	if entrance == "/" {
		entrance = ""
	}
	assetController := controllers.NewAssetController()
	facades.Route().Get("favicon.png", assetController.Favicon)
	facades.Route().Get("robots.txt", assetController.Robots)
	facades.Route().Get(entrance+"/assets/{any}", assetController.Index)
	facades.Route().Get(entrance+"/loading/{any}", assetController.Index)
	facades.Route().Get(entrance+"/{any}", assetController.Index)
	facades.Route().Get(entrance+"/", assetController.Index)
	facades.Route().Fallback(assetController.NotFound)
}
