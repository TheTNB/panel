package route

import (
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/internal/embed"
	"github.com/TheTNB/panel/internal/http/middleware"
	"github.com/TheTNB/panel/internal/service"
)

func Http(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			user := service.NewUserService()
			r.With(middleware.Throttle(5, time.Minute)).Post("/login", user.Login)
			r.Post("/logout", user.Logout)
			r.Get("/isLogin", user.IsLogin)
			r.With(middleware.MustLogin).Get("/info", user.Info)
		})

		r.Route("/info", func(r chi.Router) {
			info := service.NewInfoService()
			r.Get("/panel", info.Panel)
			r.With(middleware.MustLogin).Get("/homeApps", info.HomeApps)
			r.With(middleware.MustLogin).Get("/realtime", info.Realtime)
			r.With(middleware.MustLogin).Get("/systemInfo", info.SystemInfo)
			r.With(middleware.MustLogin).Get("/countInfo", info.CountInfo)
			r.With(middleware.MustLogin).Get("/installedDbAndPhp", info.InstalledDbAndPhp)
			r.With(middleware.MustLogin).Get("/checkUpdate", info.CheckUpdate)
			r.With(middleware.MustLogin).Get("/updateInfo", info.UpdateInfo)
			r.With(middleware.MustLogin).Post("/update", info.Update)
			r.With(middleware.MustLogin).Post("/restart", info.Restart)
		})

		r.Route("/task", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			task := service.NewTaskService()
			r.Get("/status", task.Status)
			r.Get("/", task.List)
			r.Get("/{id}", task.Get)
			r.Delete("/{id}", task.Delete)
		})

		r.Route("/website", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			website := service.NewWebsiteService()
			r.Get("/defaultConfig", website.GetDefaultConfig)
			r.Post("/defaultConfig", website.UpdateDefaultConfig)
			r.Get("/", website.List)
			r.Post("/", website.Create)
			r.Get("/{id}", website.Get)
			r.Put("/{id}", website.Update)
			r.Delete("/{id}", website.Delete)
			r.Delete("/{id}/log", website.ClearLog)
			r.Post("/{id}/updateRemark", website.UpdateRemark)
			r.Post("/{id}/resetConfig", website.ResetConfig)
			r.Post("/{id}/status", website.UpdateStatus)
		})

		// TODO
		r.Route("/backup", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			backup := service.NewBackupService()
			r.Get("/backup", backup.List)
			r.Post("/create", backup.Create)
			r.Post("/update", backup.Update)
			r.Get("/{id}", backup.Get)
			r.Delete("/{id}", backup.Delete)
			r.Delete("/{id}/restore", backup.Restore)
		})

		r.Route("/cert", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			cert := service.NewCertService()
			r.Get("/caProviders", cert.CAProviders)
			r.Get("/dnsProviders", cert.DNSProviders)
			r.Get("/algorithms", cert.Algorithms)
			r.Route("/cert", func(r chi.Router) {
				r.Get("/", cert.List)
				r.Post("/", cert.Create)
				r.Put("/{id}", cert.Update)
				r.Get("/{id}", cert.Get)
				r.Delete("/{id}", cert.Delete)
				r.Post("/{id}/obtain", cert.Obtain)
				r.Post("/{id}/renew", cert.Renew)
				r.Post("/{id}/manualDNS", cert.ManualDNS)
				r.Post("/{id}/deploy", cert.Deploy)
			})
			r.Route("/dns", func(r chi.Router) {
				certDNS := service.NewCertDNSService()
				r.Get("/", certDNS.List)
				r.Post("/", certDNS.Create)
				r.Put("/{id}", certDNS.Update)
				r.Get("/{id}", certDNS.Get)
				r.Delete("/{id}", certDNS.Delete)
			})
			r.Route("/account", func(r chi.Router) {
				certAccount := service.NewCertAccountService()
				r.Get("/", certAccount.List)
				r.Post("/", certAccount.Create)
				r.Put("/{id}", certAccount.Update)
				r.Get("/{id}", certAccount.Get)
				r.Delete("/{id}", certAccount.Delete)
			})
		})

		r.Route("/app", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			app := service.NewAppService()
			r.Get("/list", app.List)
			r.Post("/install", app.Install)
			r.Post("/uninstall", app.Uninstall)
			r.Post("/update", app.Update)
			r.Post("/updateShow", app.UpdateShow)
			r.Get("/isInstalled", app.IsInstalled)
			r.Get("/updateCache", app.UpdateCache)
		})

		r.Route("/cron", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			cron := service.NewCronService()
			r.Get("/", cron.List)
			r.Post("/", cron.Create)
			r.Put("/{id}", cron.Update)
			r.Get("/{id}", cron.Get)
			r.Delete("/{id}", cron.Delete)
			r.Post("/{id}/status", cron.Status)
			r.Get("/{id}/log", cron.Log)
		})

		r.Route("/safe", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			safe := service.NewSafeService()
			r.Get("/ssh", safe.GetSSH)
			r.Post("/ssh", safe.UpdateSSH)
			r.Get("/ping", safe.GetPingStatus)
			r.Post("/ping", safe.UpdatePingStatus)
		})

		r.Route("/firewall", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			firewall := service.NewFirewallService()
			r.Get("/status", firewall.GetStatus)
			r.Post("/status", firewall.UpdateStatus)
			r.Get("/rule", firewall.GetRules)
			r.Post("/rule", firewall.CreateRule)
			r.Delete("/rule", firewall.DeleteRule)
		})

		r.Route("/ssh", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			ssh := service.NewSSHService()
			r.Get("/info", ssh.GetInfo)
			r.Post("/info", ssh.UpdateInfo)
			r.Get("/session", ssh.Session)
		})

		r.Route("/container", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			r.Route("/container", func(r chi.Router) {
				container := service.NewContainerService()
				r.Get("/", container.List)
				r.Get("/search", container.Search)
				r.Post("/create", container.Create)
				r.Post("/remove", container.Remove)
				r.Post("/start", container.Start)
				r.Post("/stop", container.Stop)
				r.Post("/restart", container.Restart)
				r.Post("/pause", container.Pause)
				r.Post("/unpause", container.Unpause)
				r.Get("/inspect", container.Inspect)
				r.Post("/kill", container.Kill)
				r.Post("/rename", container.Rename)
				r.Get("/stats", container.Stats)
				r.Get("/exist", container.Exist)
				r.Get("/logs", container.Logs)
				r.Post("/prune", container.Prune)
			})
			r.Route("/network", func(r chi.Router) {
				containerNetwork := service.NewContainerNetworkService()
				r.Get("/list", containerNetwork.List)
				r.Post("/create", containerNetwork.Create)
				r.Post("/remove", containerNetwork.Remove)
				r.Get("/exist", containerNetwork.Exist)
				r.Get("/inspect", containerNetwork.Inspect)
				r.Post("/connect", containerNetwork.Connect)
				r.Post("/disconnect", containerNetwork.Disconnect)
				r.Post("/prune", containerNetwork.Prune)
			})
			r.Route("/image", func(r chi.Router) {
				containerImage := service.NewContainerImageService()
				r.Get("/list", containerImage.List)
				r.Get("/exist", containerImage.Exist)
				r.Post("/pull", containerImage.Pull)
				r.Post("/remove", containerImage.Remove)
				r.Get("/inspect", containerImage.Inspect)
				r.Post("/prune", containerImage.Prune)
			})
			r.Route("/volume", func(r chi.Router) {
				containerVolume := service.NewContainerVolumeService()
				r.Get("/list", containerVolume.List)
				r.Post("/create", containerVolume.Create)
				r.Get("/exist", containerVolume.Exist)
				r.Post("/remove", containerVolume.Remove)
				r.Get("/inspect", containerVolume.Inspect)
				r.Post("/prune", containerVolume.Prune)
			})
		})

		r.Route("/file", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			file := service.NewFileService()
			r.Post("/create", file.Create)
			r.Get("/content", file.Content)
			r.Post("/save", file.Save)
			r.Post("/delete", file.Delete)
			r.Post("/upload", file.Upload)
			r.Post("/move", file.Move)
			r.Post("/copy", file.Copy)
			r.Get("/download", file.Download)
			r.Post("/remoteDownload", file.RemoteDownload)
			r.Get("/info", file.Info)
			r.Post("/permission", file.Permission)
			r.Post("/compress", file.Compress)
			r.Post("/unCompress", file.UnCompress)
			r.Post("/search", file.Search)
			r.Get("/list", file.List)
		})

		r.Route("/monitor", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			monitor := service.NewMonitorService()
			r.Get("/setting", monitor.GetSetting)
			r.Post("/setting", monitor.UpdateSetting)
			r.Post("/clear", monitor.Clear)
			r.Get("/list", monitor.List)
		})

		r.Route("/setting", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			setting := service.NewSettingService()
			r.Get("/", setting.Get)
			r.Post("/", setting.Update)
		})

		r.Route("/systemctl", func(r chi.Router) {
			r.Use(middleware.MustLogin)
			systemctl := service.NewSystemctlService()
			r.Get("/status", systemctl.Status)
			r.Get("/isEnabled", systemctl.IsEnabled)
			r.Post("/enable", systemctl.Enable)
			r.Post("/disable", systemctl.Disable)
			r.Post("/restart", systemctl.Restart)
			r.Post("/reload", systemctl.Reload)
			r.Post("/start", systemctl.Start)
			r.Post("/stop", systemctl.Stop)
		})
	})

	r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		// /api 开头的返回 404
		if strings.HasPrefix(request.URL.Path, "/api") {
			http.NotFound(writer, request)
			return
		}
		// 其他返回前端页面
		frontend, _ := fs.Sub(embed.PublicFS, "frontend")
		spaHandler := func(fs http.FileSystem) http.HandlerFunc {
			fileServer := http.FileServer(fs)
			return func(w http.ResponseWriter, r *http.Request) {
				path := r.URL.Path
				f, err := fs.Open(path)
				if err != nil {
					indexFile, err := fs.Open("index.html")
					if err != nil {
						http.NotFound(w, r)
						return
					}
					defer indexFile.Close()

					fi, err := indexFile.Stat()
					if err != nil {
						http.NotFound(w, r)
						return
					}

					http.ServeContent(w, r, "index.html", fi.ModTime(), indexFile)
					return
				}
				defer f.Close()
				fileServer.ServeHTTP(w, r)
			}
		}
		spaHandler(http.FS(frontend)).ServeHTTP(writer, request)
	})
}
