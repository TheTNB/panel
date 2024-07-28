package controllers

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/hashicorp/go-version"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
	"github.com/TheTNB/panel/v2/pkg/tools"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type MenuItem struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Jump  string `json:"jump"`
}

type InfoController struct {
	plugin  internal.Plugin
	setting internal.Setting
}

func NewInfoController() *InfoController {
	return &InfoController{
		plugin:  services.NewPluginImpl(),
		setting: services.NewSettingImpl(),
	}
}

// Panel 获取面板信息
func (r *InfoController) Panel(ctx http.Context) http.Response {
	return h.Success(ctx, http.Json{
		"name":     r.setting.Get(models.SettingKeyName),
		"language": facades.Config().GetString("app.locale"),
	})
}

// HomePlugins 获取首页插件
func (r *InfoController) HomePlugins(ctx http.Context) http.Response {
	var plugins []models.Plugin
	err := facades.Orm().Query().Where("show", 1).Find(&plugins)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "基础信息").With(map[string]any{
			"error": err.Error(),
		}).Info("获取首页插件失败")
		return h.ErrorSystem(ctx)
	}

	type pluginsData struct {
		models.Plugin
		Name string `json:"name"`
	}

	var pluginsJson []pluginsData
	for _, plugin := range plugins {
		pluginsJson = append(pluginsJson, pluginsData{
			Plugin: plugin,
			Name:   r.plugin.GetBySlug(plugin.Slug).Name,
		})
	}

	return h.Success(ctx, pluginsJson)
}

// NowMonitor 获取当前监控信息
func (r *InfoController) NowMonitor(ctx http.Context) http.Response {
	return h.Success(ctx, tools.GetMonitoringInfo())
}

// SystemInfo 获取系统信息
func (r *InfoController) SystemInfo(ctx http.Context) http.Response {
	monitorInfo := tools.GetMonitoringInfo()

	return h.Success(ctx, http.Json{
		"os_name":       monitorInfo.Host.Platform + " " + monitorInfo.Host.PlatformVersion,
		"uptime":        fmt.Sprintf("%.2f", float64(monitorInfo.Host.Uptime)/86400),
		"panel_version": facades.Config().GetString("panel.version"),
	})
}

// CountInfo 获取面板统计信息
func (r *InfoController) CountInfo(ctx http.Context) http.Response {
	var websiteCount int64
	err := facades.Orm().Query().Model(models.Website{}).Count(&websiteCount)
	if err != nil {
		websiteCount = -1
	}

	var mysql models.Plugin
	mysqlInstalled := true
	err = facades.Orm().Query().Where("slug like ?", "mysql%").FirstOrFail(&mysql)
	if err != nil {
		mysqlInstalled = false
	}
	var postgresql models.Plugin
	postgresqlInstalled := true
	err = facades.Orm().Query().Where("slug like ?", "postgresql%").FirstOrFail(&postgresql)
	if err != nil {
		postgresqlInstalled = false
	}
	var databaseCount int64
	if mysqlInstalled {
		status, err := systemctl.Status("mysqld")
		if status && err == nil {
			rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
			type database struct {
				Name string `json:"name"`
			}

			db, err := sql.Open("mysql", "root:"+rootPassword+"@unix(/tmp/mysql.sock)/")
			if err != nil {
				facades.Log().Request(ctx.Request()).Tags("面板", "基础信息").With(map[string]any{
					"error": err.Error(),
				}).Info("获取数据库列表失败")
				databaseCount = -1
			} else {
				defer db.Close()
				rows, err := db.Query("SHOW DATABASES")
				if err != nil {
					facades.Log().Request(ctx.Request()).Tags("面板", "基础信息").With(map[string]any{
						"error": err.Error(),
					}).Info("获取数据库列表失败")
					databaseCount = -1
				} else {
					defer rows.Close()
					var databases []database
					for rows.Next() {
						var d database
						err := rows.Scan(&d.Name)
						if err != nil {
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
		status, err := systemctl.Status("postgresql")
		if status && err == nil {
			raw, err := shell.Execf(`echo "\l" | su - postgres -c "psql"`)
			if err == nil {
				databases := strings.Split(raw, "\n")
				if len(databases) >= 4 {
					databases = databases[3 : len(databases)-1]
					for _, db := range databases {
						parts := strings.Split(db, "|")
						if len(parts) != 9 || len(strings.TrimSpace(parts[0])) == 0 || strings.TrimSpace(parts[0]) == "template0" || strings.TrimSpace(parts[0]) == "template1" || strings.TrimSpace(parts[0]) == "postgres" {
							continue
						}

						databaseCount++
					}
				}
			}
		}
	}

	var ftpCount int64
	var ftpPlugin = r.plugin.GetInstalledBySlug("pureftpd")
	if ftpPlugin.ID != 0 {
		listRaw, err := shell.Execf("pure-pw list")
		if len(listRaw) != 0 && err == nil {
			listArr := strings.Split(listRaw, "\n")
			ftpCount = int64(len(listArr))
		}
	}

	var cronCount int64
	err = facades.Orm().Query().Model(models.Cron{}).Count(&cronCount)
	if err != nil {
		cronCount = -1
	}

	return h.Success(ctx, http.Json{
		"website":  websiteCount,
		"database": databaseCount,
		"ftp":      ftpCount,
		"cron":     cronCount,
	})
}

// InstalledDbAndPhp 获取已安装的数据库和 PHP 版本
func (r *InfoController) InstalledDbAndPhp(ctx http.Context) http.Response {
	var php []models.Plugin
	err := facades.Orm().Query().Where("slug like ?", "php%").Find(&php)
	if err != nil {
		return h.ErrorSystem(ctx)
	}

	var mysql models.Plugin
	mysqlInstalled := true
	err = facades.Orm().Query().Where("slug like ?", "mysql%").FirstOrFail(&mysql)
	if err != nil {
		mysqlInstalled = false
	}

	var postgresql models.Plugin
	postgresqlInstalled := true
	err = facades.Orm().Query().Where("slug like ?", "postgresql%").FirstOrFail(&postgresql)
	if err != nil {
		postgresqlInstalled = false
	}

	type data struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}
	var phpData []data
	var dbData []data
	phpData = append(phpData, data{Value: "0", Label: "不使用"})
	dbData = append(dbData, data{Value: "0", Label: "不使用"})
	for _, p := range php {
		match := regexp.MustCompile(`php(\d+)`).FindStringSubmatch(p.Slug)
		if len(match) == 0 {
			continue
		}

		phpData = append(phpData, data{Value: strings.ReplaceAll(p.Slug, "php", ""), Label: r.plugin.GetBySlug(p.Slug).Name})
	}

	if mysqlInstalled {
		dbData = append(dbData, data{Value: "mysql", Label: "MySQL"})
	}
	if postgresqlInstalled {
		dbData = append(dbData, data{Value: "postgresql", Label: "PostgreSQL"})
	}

	return h.Success(ctx, http.Json{
		"php": phpData,
		"db":  dbData,
	})
}

// CheckUpdate 检查面板更新
func (r *InfoController) CheckUpdate(ctx http.Context) http.Response {
	current := facades.Config().GetString("panel.version")
	latest, err := tools.GetLatestPanelVersion()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取最新版本失败")
	}

	v1, err := version.NewVersion(current)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "版本号解析失败")
	}
	v2, err := version.NewVersion(latest.Version)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "版本号解析失败")
	}
	if v1.GreaterThanOrEqual(v2) {
		return h.Success(ctx, http.Json{
			"update": false,
		})
	}

	return h.Success(ctx, http.Json{
		"update": true,
	})
}

// UpdateInfo 获取更新信息
func (r *InfoController) UpdateInfo(ctx http.Context) http.Response {
	current := facades.Config().GetString("panel.version")
	latest, err := tools.GetLatestPanelVersion()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取最新版本失败")
	}

	v1, err := version.NewVersion(current)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "版本号解析失败")
	}
	v2, err := version.NewVersion(latest.Version)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "版本号解析失败")
	}
	if v1.GreaterThanOrEqual(v2) {
		return h.Error(ctx, http.StatusInternalServerError, "当前版本已是最新版本")
	}

	versions, err := tools.GenerateVersions(current, latest.Version)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取更新信息失败")
	}

	var versionInfo []tools.PanelInfo
	for _, v := range versions {
		info, err := tools.GetPanelVersion(v)
		if err != nil {
			continue
		}

		versionInfo = append(versionInfo, info)
	}

	return h.Success(ctx, versionInfo)
}

// Update 更新面板
func (r *InfoController) Update(ctx http.Context) http.Response {
	var task models.Task
	if err := facades.Orm().Query().Where("status", models.TaskStatusRunning).OrWhere("status", models.TaskStatusWaiting).FirstOrFail(&task); err == nil {
		return h.Error(ctx, http.StatusInternalServerError, "当前有任务正在执行，禁止更新")
	}
	if _, err := facades.Orm().Query().Exec("PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
		types.Status = types.StatusFailed
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("面板数据库异常，已终止操作：%s", err.Error()))
	}

	panel, err := tools.GetLatestPanelVersion()
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "基础信息").With(map[string]any{
			"error": err.Error(),
		}).Info("获取最新版本失败")
		return h.Error(ctx, http.StatusInternalServerError, "获取最新版本失败")
	}

	types.Status = types.StatusUpgrade
	if err = tools.UpdatePanel(panel); err != nil {
		types.Status = types.StatusFailed
		facades.Log().Request(ctx.Request()).Tags("面板", "基础信息").With(map[string]any{
			"error": err.Error(),
		}).Info("更新面板失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	types.Status = types.StatusNormal
	tools.RestartPanel()
	return h.Success(ctx, nil)
}

// Restart 重启面板
func (r *InfoController) Restart(ctx http.Context) http.Response {
	var task models.Task
	err := facades.Orm().Query().Where("status", models.TaskStatusRunning).OrWhere("status", models.TaskStatusWaiting).FirstOrFail(&task)
	if err == nil {
		return h.Error(ctx, http.StatusInternalServerError, "当前有任务正在执行，禁止重启")
	}

	tools.RestartPanel()
	return h.Success(ctx, nil)
}
