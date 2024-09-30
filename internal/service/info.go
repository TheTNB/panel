package service

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-rat/chix"
	"github.com/hashicorp/go-version"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/pkg/db"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/pkg/types"
)

type InfoService struct {
	taskRepo    biz.TaskRepo
	websiteRepo biz.WebsiteRepo
	appRepo     biz.AppRepo
	settingRepo biz.SettingRepo
	cronRepo    biz.CronRepo
}

func NewInfoService() *InfoService {
	return &InfoService{
		taskRepo:    data.NewTaskRepo(),
		websiteRepo: data.NewWebsiteRepo(),
		appRepo:     data.NewAppRepo(),
		settingRepo: data.NewSettingRepo(),
		cronRepo:    data.NewCronRepo(),
	}
}

func (s *InfoService) Panel(w http.ResponseWriter, r *http.Request) {
	name, _ := s.settingRepo.Get(biz.SettingKeyName)
	if name == "" {
		name = "耗子面板"
	}

	Success(w, chix.M{
		"name":     name,
		"language": app.Conf.MustString("app.locale"),
	})
}

func (s *InfoService) HomeApps(w http.ResponseWriter, r *http.Request) {
	apps, err := s.appRepo.GetHomeShow()
	if err != nil {
		Error(w, http.StatusInternalServerError, "获取首页应用失败")
		return
	}

	Success(w, apps)
}

func (s *InfoService) Realtime(w http.ResponseWriter, r *http.Request) {
	Success(w, tools.GetMonitoringInfo())
}

func (s *InfoService) SystemInfo(w http.ResponseWriter, r *http.Request) {
	monitorInfo := tools.GetMonitoringInfo()

	Success(w, chix.M{
		"os_name":       monitorInfo.Host.Platform + " " + monitorInfo.Host.PlatformVersion,
		"uptime":        fmt.Sprintf("%.2f", float64(monitorInfo.Host.Uptime)/86400),
		"panel_version": app.Conf.MustString("app.version"),
	})
}

func (s *InfoService) CountInfo(w http.ResponseWriter, r *http.Request) {
	websiteCount, err := s.websiteRepo.Count()
	if err != nil {
		Error(w, http.StatusInternalServerError, "获取网站数量失败")
		return
	}

	mysqlInstalled, _ := s.appRepo.IsInstalled("slug like ?", "mysql%")
	postgresqlInstalled, _ := s.appRepo.IsInstalled("slug like ?", "postgresql%")

	type database struct {
		Name string `json:"name"`
	}
	var databaseCount int64
	if mysqlInstalled {
		rootPassword, _ := s.settingRepo.Get(biz.SettingKeyPerconaRootPassword)
		mysql, err := db.NewMySQL("root", rootPassword, "/tmp/mysql.sock")
		if err == nil {
			defer mysql.Close()
			if err = mysql.Ping(); err != nil {
				databaseCount = -1
			} else {
				rows, err := mysql.Query("SHOW DATABASES")
				if err != nil {
					databaseCount = -1
				} else {
					defer rows.Close()
					var databases []database
					for rows.Next() {
						var d database
						if err := rows.Scan(&d.Name); err != nil {
							continue
						}
						if d.Name == "information_schema" || d.Name == "performance_schema" || d.Name == "mysql" || d.Name == "sys" {
							continue
						}

						databases = append(databases, d)
					}
					databaseCount = int64(len(databases))
				}
			}
		}
	}
	if postgresqlInstalled {
		postgres, err := db.NewPostgres("postgres", "", "127.0.0.1", fmt.Sprintf("%s/server/postgresql/data/pg_hba.conf", app.Root), 5432)
		if err == nil {
			defer postgres.Close()
			if err = postgres.Ping(); err != nil {
				databaseCount = -1
			} else {
				rows, err := postgres.Query("SELECT datname FROM pg_database WHERE datistemplate = false")
				if err != nil {
					databaseCount = -1
				} else {
					defer rows.Close()
					var databases []database
					for rows.Next() {
						var d database
						if err = rows.Scan(&d.Name); err != nil {
							continue
						}
						if d.Name == "postgres" || d.Name == "template0" || d.Name == "template1" {
							continue
						}
						databases = append(databases, d)
					}
					databaseCount = int64(len(databases))
				}
			}
		}
	}

	var ftpCount int64
	ftpInstalled, _ := s.appRepo.IsInstalled("slug = ?", "pureftpd")
	if ftpInstalled {
		listRaw, err := shell.Execf("pure-pw list")
		if len(listRaw) != 0 && err == nil {
			listArr := strings.Split(listRaw, "\n")
			ftpCount = int64(len(listArr))
		}
	}

	cronCount, err := s.cronRepo.Count()
	if err != nil {
		cronCount = -1
	}

	Success(w, chix.M{
		"website":  websiteCount,
		"database": databaseCount,
		"ftp":      ftpCount,
		"cron":     cronCount,
	})
}

func (s *InfoService) InstalledDbAndPhp(w http.ResponseWriter, r *http.Request) {
	mysqlInstalled, _ := s.appRepo.IsInstalled("slug like ?", "mysql%")
	postgresqlInstalled, _ := s.appRepo.IsInstalled("slug like ?", "postgresql%")
	php, _ := s.appRepo.GetInstalledAll("slug like ?", "php%")

	var phpData []types.LV
	var dbData []types.LV
	phpData = append(phpData, types.LV{Value: "0", Label: "不使用"})
	dbData = append(dbData, types.LV{Value: "0", Label: "不使用"})
	for _, p := range php {
		// 过滤 phpmyadmin
		match := regexp.MustCompile(`php(\d+)`).FindStringSubmatch(p.Slug)
		if len(match) == 0 {
			continue
		}

		app, _ := s.appRepo.Get(p.Slug)
		phpData = append(phpData, types.LV{Value: strings.ReplaceAll(p.Slug, "php", ""), Label: app.Name})
	}

	if mysqlInstalled {
		dbData = append(dbData, types.LV{Value: "mysql", Label: "MySQL"})
	}
	if postgresqlInstalled {
		dbData = append(dbData, types.LV{Value: "postgresql", Label: "PostgreSQL"})
	}

	Success(w, chix.M{
		"php": phpData,
		"db":  dbData,
	})
}

func (s *InfoService) CheckUpdate(w http.ResponseWriter, r *http.Request) {
	current := app.Conf.MustString("app.version")
	latest, err := tools.GetLatestPanelVersion()
	if err != nil {
		Error(w, http.StatusInternalServerError, "获取最新版本失败")
		return
	}

	v1, err := version.NewVersion(current)
	if err != nil {
		Error(w, http.StatusInternalServerError, "版本号解析失败")
		return
	}
	v2, err := version.NewVersion(latest.Version)
	if err != nil {
		Error(w, http.StatusInternalServerError, "版本号解析失败")
		return
	}
	if v1.GreaterThanOrEqual(v2) {
		Success(w, chix.M{
			"update": false,
		})
		return
	}

	Success(w, chix.M{
		"update": true,
	})
}

func (s *InfoService) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	current := app.Conf.MustString("app.version")
	latest, err := tools.GetLatestPanelVersion()
	if err != nil {
		Error(w, http.StatusInternalServerError, "获取最新版本失败")
		return
	}

	v1, err := version.NewVersion(current)
	if err != nil {
		Error(w, http.StatusInternalServerError, "版本号解析失败")
		return
	}
	v2, err := version.NewVersion(latest.Version)
	if err != nil {
		Error(w, http.StatusInternalServerError, "版本号解析失败")
		return
	}
	if v1.GreaterThanOrEqual(v2) {
		Error(w, http.StatusInternalServerError, "当前版本已是最新版本")
		return
	}

	versions, err := tools.GenerateVersions(current, latest.Version)
	if err != nil {
		Error(w, http.StatusInternalServerError, "获取更新信息失败")
		return
	}

	var versionInfo []tools.PanelInfo
	for _, v := range versions {
		info, err := tools.GetPanelVersion(v)
		if err != nil {
			continue
		}

		versionInfo = append(versionInfo, info)
	}

	Success(w, versionInfo)
}

func (s *InfoService) Update(w http.ResponseWriter, r *http.Request) {
	if s.taskRepo.HasRunningTask() {
		Error(w, http.StatusInternalServerError, "当前有任务正在执行，禁止更新")
		return
	}
	if err := app.Orm.Exec("PRAGMA wal_checkpoint(TRUNCATE)").Error; err != nil {
		types.Status = types.StatusFailed
		Error(w, http.StatusInternalServerError, fmt.Sprintf("面板数据库异常，已终止操作：%s", err.Error()))
		return
	}

	panel, err := tools.GetLatestPanelVersion()
	if err != nil {
		Error(w, http.StatusInternalServerError, "获取最新版本失败")
		return
	}

	types.Status = types.StatusUpgrade
	if err = tools.UpdatePanel(panel); err != nil {
		types.Status = types.StatusFailed
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	types.Status = types.StatusNormal
	tools.RestartPanel()
	Success(w, nil)
}

func (s *InfoService) Restart(w http.ResponseWriter, r *http.Request) {
	if s.taskRepo.HasRunningTask() {
		Error(w, http.StatusInternalServerError, "当前有任务正在执行，禁止重启")
		return
	}

	tools.RestartPanel()
	Success(w, nil)
}
