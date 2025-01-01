package route

import (
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/tnb-labs/panel/internal/apps"
	"github.com/tnb-labs/panel/internal/embed"
	"github.com/tnb-labs/panel/internal/http/middleware"
	"github.com/tnb-labs/panel/internal/service"
)

type Http struct {
	user             *service.UserService
	dashboard        *service.DashboardService
	task             *service.TaskService
	website          *service.WebsiteService
	database         *service.DatabaseService
	databaseServer   *service.DatabaseServerService
	databaseUser     *service.DatabaseUserService
	backup           *service.BackupService
	cert             *service.CertService
	certDNS          *service.CertDNSService
	certAccount      *service.CertAccountService
	app              *service.AppService
	cron             *service.CronService
	process          *service.ProcessService
	safe             *service.SafeService
	firewall         *service.FirewallService
	ssh              *service.SSHService
	container        *service.ContainerService
	containerNetwork *service.ContainerNetworkService
	containerImage   *service.ContainerImageService
	containerVolume  *service.ContainerVolumeService
	file             *service.FileService
	monitor          *service.MonitorService
	setting          *service.SettingService
	systemctl        *service.SystemctlService
	apps             *apps.Loader
}

func NewHttp(
	user *service.UserService,
	dashboard *service.DashboardService,
	task *service.TaskService,
	website *service.WebsiteService,
	database *service.DatabaseService,
	databaseServer *service.DatabaseServerService,
	databaseUser *service.DatabaseUserService,
	backup *service.BackupService,
	cert *service.CertService,
	certDNS *service.CertDNSService,
	certAccount *service.CertAccountService,
	app *service.AppService,
	cron *service.CronService,
	process *service.ProcessService,
	safe *service.SafeService,
	firewall *service.FirewallService,
	ssh *service.SSHService,
	container *service.ContainerService,
	containerNetwork *service.ContainerNetworkService,
	containerImage *service.ContainerImageService,
	containerVolume *service.ContainerVolumeService,
	file *service.FileService,
	monitor *service.MonitorService,
	setting *service.SettingService,
	systemctl *service.SystemctlService,
	apps *apps.Loader,
) *Http {
	return &Http{
		user:             user,
		dashboard:        dashboard,
		task:             task,
		website:          website,
		database:         database,
		databaseServer:   databaseServer,
		databaseUser:     databaseUser,
		backup:           backup,
		cert:             cert,
		certDNS:          certDNS,
		certAccount:      certAccount,
		app:              app,
		cron:             cron,
		process:          process,
		safe:             safe,
		firewall:         firewall,
		ssh:              ssh,
		container:        container,
		containerNetwork: containerNetwork,
		containerImage:   containerImage,
		containerVolume:  containerVolume,
		file:             file,
		monitor:          monitor,
		setting:          setting,
		systemctl:        systemctl,
		apps:             apps,
	}
}

func (route *Http) Register(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Get("/key", route.user.GetKey)
			r.With(middleware.Throttle(5, time.Minute)).Post("/login", route.user.Login)
			r.Post("/logout", route.user.Logout)
			r.Get("/isLogin", route.user.IsLogin)
			r.Get("/info", route.user.Info)
		})

		r.Route("/dashboard", func(r chi.Router) {
			r.Get("/panel", route.dashboard.Panel)
			r.Get("/homeApps", route.dashboard.HomeApps)
			r.Post("/current", route.dashboard.Current)
			r.Get("/systemInfo", route.dashboard.SystemInfo)
			r.Get("/countInfo", route.dashboard.CountInfo)
			r.Get("/installedDbAndPhp", route.dashboard.InstalledDbAndPhp)
			r.Get("/checkUpdate", route.dashboard.CheckUpdate)
			r.Get("/updateInfo", route.dashboard.UpdateInfo)
			r.Post("/update", route.dashboard.Update)
			r.Post("/restart", route.dashboard.Restart)
		})

		r.Route("/task", func(r chi.Router) {
			r.Get("/status", route.task.Status)
			r.Get("/", route.task.List)
			r.Get("/{id}", route.task.Get)
			r.Delete("/{id}", route.task.Delete)
		})

		r.Route("/website", func(r chi.Router) {
			r.Get("/rewrites", route.website.GetRewrites)
			r.Get("/defaultConfig", route.website.GetDefaultConfig)
			r.Post("/defaultConfig", route.website.UpdateDefaultConfig)
			r.Get("/", route.website.List)
			r.Post("/", route.website.Create)
			r.Get("/{id}", route.website.Get)
			r.Put("/{id}", route.website.Update)
			r.Delete("/{id}", route.website.Delete)
			r.Delete("/{id}/log", route.website.ClearLog)
			r.Post("/{id}/updateRemark", route.website.UpdateRemark)
			r.Post("/{id}/resetConfig", route.website.ResetConfig)
			r.Post("/{id}/status", route.website.UpdateStatus)
			r.Post("/{id}/obtainCert", route.website.ObtainCert)
		})

		r.Route("/database", func(r chi.Router) {
			r.Get("/", route.database.List)
			r.Post("/", route.database.Create)
			r.Delete("/", route.database.Delete)
			r.Post("/comment", route.database.Comment)
		})

		r.Route("/databaseServer", func(r chi.Router) {
			r.Get("/", route.databaseServer.List)
			r.Post("/", route.databaseServer.Create)
			r.Get("/{id}", route.databaseServer.Get)
			r.Put("/{id}", route.databaseServer.Update)
			r.Put("/{id}/remark", route.databaseServer.UpdateRemark)
			r.Delete("/{id}", route.databaseServer.Delete)
			r.Post("/{id}/sync", route.databaseServer.Sync)
		})

		r.Route("/databaseUser", func(r chi.Router) {
			r.Get("/", route.databaseUser.List)
			r.Post("/", route.databaseUser.Create)
			r.Get("/{id}", route.databaseUser.Get)
			r.Put("/{id}", route.databaseUser.Update)
			r.Put("/{id}/remark", route.databaseUser.UpdateRemark)
			r.Delete("/{id}", route.databaseUser.Delete)
		})

		r.Route("/backup", func(r chi.Router) {
			r.Get("/{type}", route.backup.List)
			r.Post("/{type}", route.backup.Create)
			r.Post("/{type}/upload", route.backup.Upload)
			r.Delete("/{type}/delete", route.backup.Delete)
			r.Post("/{type}/restore", route.backup.Restore)
		})

		r.Route("/cert", func(r chi.Router) {
			r.Get("/caProviders", route.cert.CAProviders)
			r.Get("/dnsProviders", route.cert.DNSProviders)
			r.Get("/algorithms", route.cert.Algorithms)
			r.Route("/cert", func(r chi.Router) {
				r.Get("/", route.cert.List)
				r.Post("/", route.cert.Create)
				r.Post("/upload", route.cert.Upload)
				r.Put("/{id}", route.cert.Update)
				r.Get("/{id}", route.cert.Get)
				r.Delete("/{id}", route.cert.Delete)
				r.Post("/{id}/obtainAuto", route.cert.ObtainAuto)
				r.Post("/{id}/obtainManual", route.cert.ObtainManual)
				r.Post("/{id}/obtainSelfSigned", route.cert.ObtainSelfSigned)
				r.Post("/{id}/renew", route.cert.Renew)
				r.Post("/{id}/manualDNS", route.cert.ManualDNS)
				r.Post("/{id}/deploy", route.cert.Deploy)
			})
			r.Route("/dns", func(r chi.Router) {
				r.Get("/", route.certDNS.List)
				r.Post("/", route.certDNS.Create)
				r.Put("/{id}", route.certDNS.Update)
				r.Get("/{id}", route.certDNS.Get)
				r.Delete("/{id}", route.certDNS.Delete)
			})
			r.Route("/account", func(r chi.Router) {
				r.Get("/", route.certAccount.List)
				r.Post("/", route.certAccount.Create)
				r.Put("/{id}", route.certAccount.Update)
				r.Get("/{id}", route.certAccount.Get)
				r.Delete("/{id}", route.certAccount.Delete)
			})
		})

		r.Route("/app", func(r chi.Router) {
			r.Get("/list", route.app.List)
			r.Post("/install", route.app.Install)
			r.Post("/uninstall", route.app.Uninstall)
			r.Post("/update", route.app.Update)
			r.Post("/updateShow", route.app.UpdateShow)
			r.Get("/isInstalled", route.app.IsInstalled)
			r.Get("/updateCache", route.app.UpdateCache)
		})

		r.Route("/cron", func(r chi.Router) {
			r.Get("/", route.cron.List)
			r.Post("/", route.cron.Create)
			r.Put("/{id}", route.cron.Update)
			r.Get("/{id}", route.cron.Get)
			r.Delete("/{id}", route.cron.Delete)
			r.Post("/{id}/status", route.cron.Status)
		})

		r.Route("/process", func(r chi.Router) {
			r.Get("/", route.process.List)
			r.Post("/kill", route.process.Kill)
		})

		r.Route("/safe", func(r chi.Router) {
			r.Get("/ssh", route.safe.GetSSH)
			r.Post("/ssh", route.safe.UpdateSSH)
			r.Get("/ping", route.safe.GetPingStatus)
			r.Post("/ping", route.safe.UpdatePingStatus)
		})

		r.Route("/firewall", func(r chi.Router) {
			r.Get("/status", route.firewall.GetStatus)
			r.Post("/status", route.firewall.UpdateStatus)
			r.Get("/rule", route.firewall.GetRules)
			r.Post("/rule", route.firewall.CreateRule)
			r.Delete("/rule", route.firewall.DeleteRule)
			r.Get("/ipRule", route.firewall.GetIPRules)
			r.Post("/ipRule", route.firewall.CreateIPRule)
			r.Delete("/ipRule", route.firewall.DeleteIPRule)
			r.Get("/forward", route.firewall.GetForwards)
			r.Post("/forward", route.firewall.CreateForward)
			r.Delete("/forward", route.firewall.DeleteForward)
		})

		r.Route("/ssh", func(r chi.Router) {
			r.Get("/", route.ssh.List)
			r.Post("/", route.ssh.Create)
			r.Put("/{id}", route.ssh.Update)
			r.Get("/{id}", route.ssh.Get)
			r.Delete("/{id}", route.ssh.Delete)
		})

		r.Route("/container", func(r chi.Router) {
			r.Route("/container", func(r chi.Router) {
				r.Get("/", route.container.List)
				r.Get("/search", route.container.Search)
				r.Post("/", route.container.Create)
				r.Delete("/{id}", route.container.Remove)
				r.Post("/{id}/start", route.container.Start)
				r.Post("/{id}/stop", route.container.Stop)
				r.Post("/{id}/restart", route.container.Restart)
				r.Post("/{id}/pause", route.container.Pause)
				r.Post("/{id}/unpause", route.container.Unpause)
				r.Post("/{id}/kill", route.container.Kill)
				r.Post("/{id}/rename", route.container.Rename)
				r.Get("/{id}/logs", route.container.Logs)
				r.Post("/prune", route.container.Prune)
			})
			r.Route("/network", func(r chi.Router) {
				r.Get("/", route.containerNetwork.List)
				r.Post("/", route.containerNetwork.Create)
				r.Delete("/{id}", route.containerNetwork.Remove)
				r.Post("/prune", route.containerNetwork.Prune)
			})
			r.Route("/image", func(r chi.Router) {
				r.Get("/", route.containerImage.List)
				r.Post("/", route.containerImage.Pull)
				r.Delete("/{id}", route.containerImage.Remove)
				r.Post("/prune", route.containerImage.Prune)
			})
			r.Route("/volume", func(r chi.Router) {
				r.Get("/", route.containerVolume.List)
				r.Post("/", route.containerVolume.Create)
				r.Delete("/{id}", route.containerVolume.Remove)
				r.Post("/prune", route.containerVolume.Prune)
			})
		})

		r.Route("/file", func(r chi.Router) {
			r.Post("/create", route.file.Create)
			r.Get("/content", route.file.Content)
			r.Post("/save", route.file.Save)
			r.Post("/delete", route.file.Delete)
			r.Post("/upload", route.file.Upload)
			r.Post("/exist", route.file.Exist)
			r.Post("/move", route.file.Move)
			r.Post("/copy", route.file.Copy)
			r.Get("/download", route.file.Download)
			r.Post("/remoteDownload", route.file.RemoteDownload)
			r.Get("/info", route.file.Info)
			r.Post("/permission", route.file.Permission)
			r.Post("/compress", route.file.Compress)
			r.Post("/unCompress", route.file.UnCompress)
			r.Get("/search", route.file.Search)
			r.Get("/list", route.file.List)
		})

		r.Route("/monitor", func(r chi.Router) {
			r.Get("/setting", route.monitor.GetSetting)
			r.Post("/setting", route.monitor.UpdateSetting)
			r.Post("/clear", route.monitor.Clear)
			r.Get("/list", route.monitor.List)
		})

		r.Route("/setting", func(r chi.Router) {
			r.Get("/", route.setting.Get)
			r.Post("/", route.setting.Update)
		})

		r.Route("/systemctl", func(r chi.Router) {
			r.Get("/status", route.systemctl.Status)
			r.Get("/isEnabled", route.systemctl.IsEnabled)
			r.Post("/enable", route.systemctl.Enable)
			r.Post("/disable", route.systemctl.Disable)
			r.Post("/restart", route.systemctl.Restart)
			r.Post("/reload", route.systemctl.Reload)
			r.Post("/start", route.systemctl.Start)
			r.Post("/stop", route.systemctl.Stop)
		})

		r.Route("/apps", func(r chi.Router) {
			route.apps.Register(r)
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
