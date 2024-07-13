package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/color"
	"github.com/spf13/cast"

	requests "github.com/TheTNB/panel/v2/app/http/requests/website"
	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/str"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
	"github.com/TheTNB/panel/v2/pkg/tools"
	"github.com/TheTNB/panel/v2/pkg/types"
)

// Panel 面板命令行
type Panel struct {
}

// Signature The name and signature of the console command.
func (receiver *Panel) Signature() string {
	return "panel"
}

// Description The console command description.
func (receiver *Panel) Description() string {
	ctx := context.Background()
	return facades.Lang(ctx).Get("commands.panel.description")
}

// Extend The console command extend.
func (receiver *Panel) Extend() command.Extend {
	return command.Extend{
		Category: "panel",
	}
}

// Handle Execute the console command.
func (receiver *Panel) Handle(ctx console.Context) error {
	action := ctx.Argument(0)
	arg1 := ctx.Argument(1)
	arg2 := ctx.Argument(2)
	arg3 := ctx.Argument(3)
	arg4 := ctx.Argument(4)
	arg5 := ctx.Argument(5)

	translate := facades.Lang(context.Background())

	switch action {
	case "init":
		var check models.User
		err := facades.Orm().Query().FirstOrFail(&check)
		if err == nil {
			color.Red().Printfln(translate.Get("commands.panel.init.exist"))
			return nil
		}

		settings := []models.Setting{{Key: models.SettingKeyName, Value: "耗子面板"}, {Key: models.SettingKeyMonitor, Value: "1"}, {Key: models.SettingKeyMonitorDays, Value: "30"}, {Key: models.SettingKeyBackupPath, Value: "/www/backup"}, {Key: models.SettingKeyWebsitePath, Value: "/www/wwwroot"}, {Key: models.SettingKeyVersion, Value: facades.Config().GetString("panel.version")}}
		err = facades.Orm().Query().Create(&settings)
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.init.fail"))
			return nil
		}

		hash, err := facades.Hash().Make(str.RandomString(32))
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.init.fail"))
			return nil
		}

		user := services.NewUserImpl()
		_, err = user.Create("admin", hash)
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.init.adminFail"))
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.init.success"))

	case "update":
		var task models.Task
		if err := facades.Orm().Query().Where("status", models.TaskStatusRunning).OrWhere("status", models.TaskStatusWaiting).FirstOrFail(&task); err == nil {
			color.Red().Printfln(translate.Get("commands.panel.update.taskCheck"))
			return nil
		}
		if _, err := facades.Orm().Query().Exec("PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
			types.Status = types.StatusFailed
			color.Red().Printfln(translate.Get("commands.panel.update.dbFail"))
			return nil
		}

		panel, err := tools.GetLatestPanelVersion()
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.update.versionFail"))
			return err
		}

		// 停止面板服务，因为在shell中运行的和systemd的不同
		_ = systemctl.Stop("panel")

		types.Status = types.StatusUpgrade
		if err = tools.UpdatePanel(panel); err != nil {
			types.Status = types.StatusFailed
			color.Red().Printfln(translate.Get("commands.panel.update.fail") + ": " + err.Error())
			return nil
		}

		types.Status = types.StatusNormal
		color.Green().Printfln(translate.Get("commands.panel.update.success"))
		tools.RestartPanel()

	case "getInfo":
		var user models.User
		err := facades.Orm().Query().Where("id", 1).FirstOrFail(&user)
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.getInfo.adminGetFail"))
			return nil
		}

		password := str.RandomString(16)
		hash, err := facades.Hash().Make(password)
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.getInfo.passwordGenerationFail"))
			return nil
		}
		user.Username = str.RandomString(8)
		user.Password = hash
		if user.Email == "" {
			user.Email = str.RandomString(8) + "@example.com"
		}

		err = facades.Orm().Query().Save(&user)
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.getInfo.adminSaveFail"))
			return nil
		}

		port, err := shell.Execf(`cat /www/panel/panel.conf | grep APP_PORT | awk -F '=' '{print $2}' | tr -d '\n'`)
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.portFail"))
			return nil
		}
		ip, err := tools.GetPublicIP()
		if err != nil {
			ip = "127.0.0.1"
		}
		protocol := "http"
		if facades.Config().GetBool("panel.ssl") {
			protocol = "https"
		}

		color.Green().Printfln(translate.Get("commands.panel.getInfo.username") + ": " + user.Username)
		color.Green().Printfln(translate.Get("commands.panel.getInfo.password") + ": " + password)
		color.Green().Printfln(translate.Get("commands.panel.port") + ": " + port)
		color.Green().Printfln(translate.Get("commands.panel.entrance") + ": " + facades.Config().GetString("panel.entrance"))
		color.Green().Printfln(translate.Get("commands.panel.getInfo.address") + ": " + protocol + "://" + ip + ":" + port + facades.Config().GetString("panel.entrance"))

	case "getPort":
		port, err := shell.Execf(`cat /www/panel/panel.conf | grep APP_PORT | awk -F '=' '{print $2}' | tr -d '\n'`)
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.portFail"))
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.port") + ": " + port)

	case "getEntrance":
		color.Green().Printfln(translate.Get("commands.panel.entrance") + ": " + facades.Config().GetString("panel.entrance"))

	case "deleteEntrance":
		oldEntrance, err := shell.Execf(`cat /www/panel/panel.conf | grep APP_ENTRANCE | awk -F '=' '{print $2}' | tr -d '\n'`)
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.deleteEntrance.fail"))
			return nil
		}
		if _, err = shell.Execf("sed -i 's!APP_ENTRANCE=" + oldEntrance + "!APP_ENTRANCE=/!g' /www/panel/panel.conf"); err != nil {
			color.Red().Printfln(translate.Get("commands.panel.deleteEntrance.fail"))
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.deleteEntrance.success"))

	case "writePlugin":
		slug := arg1
		version := arg2
		if len(slug) == 0 || len(version) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.writePlugin.paramFail"))
			return nil
		}

		var plugin models.Plugin
		err := facades.Orm().Query().UpdateOrCreate(&plugin, models.Plugin{
			Slug: slug,
		}, models.Plugin{
			Version: version,
		})

		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.writePlugin.fail"))
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.writePlugin.success"))

	case "deletePlugin":
		slug := arg1
		if len(slug) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.deletePlugin.paramFail"))
			return nil
		}

		_, err := facades.Orm().Query().Where("slug", slug).Delete(&models.Plugin{})
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.deletePlugin.fail"))
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.deletePlugin.success"))

	case "writeMysqlPassword":
		password := arg1
		if len(password) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.writeMysqlPassword.paramFail"))
			return nil
		}

		var setting models.Setting
		err := facades.Orm().Query().UpdateOrCreate(&setting, models.Setting{
			Key: models.SettingKeyMysqlRootPassword,
		}, models.Setting{
			Value: password,
		})

		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.writeMysqlPassword.fail"))
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.writeMysqlPassword.success"))

	case "cleanTask":
		_, err := facades.Orm().Query().Model(&models.Task{}).Where("status", models.TaskStatusRunning).OrWhere("status", models.TaskStatusWaiting).Update("status", models.TaskStatusFailed)
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.cleanTask.fail"))
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.cleanTask.success"))

	case "backup":
		backupType := arg1
		name := arg2
		path := arg3
		save := arg4
		hr := `+----------------------------------------------------`
		if len(backupType) == 0 || len(name) == 0 || len(path) == 0 || len(save) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.backup.paramFail"))
			return nil
		}

		color.Green().Printfln(hr)
		color.Green().Printfln("★ " + translate.Get("commands.panel.backup.start") + " [" + carbon.Now().ToDateTimeString() + "]")
		color.Green().Printfln(hr)

		if !io.Exists(path) {
			if err := io.Mkdir(path, 0644); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.backupDirFail") + ": " + err.Error())
				return nil
			}
		}

		switch backupType {
		case "website":
			color.Yellow().Printfln("|-" + translate.Get("commands.panel.backup.targetSite") + ": " + name)
			var website models.Website
			if err := facades.Orm().Query().Where("name", name).FirstOrFail(&website); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.siteNotExist"))
				color.Green().Printfln(hr)
				return nil
			}

			backupFile := path + "/" + website.Name + "_" + carbon.Now().ToShortDateTimeString() + ".zip"
			if _, err := shell.Execf(`cd '` + website.Path + `' && zip -r '` + backupFile + `' .`); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.backupFail") + ": " + err.Error())
				return nil
			}
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.backupSuccess"))

		case "mysql":
			rootPassword := services.NewSettingImpl().Get(models.SettingKeyMysqlRootPassword)
			backupFile := name + "_" + carbon.Now().ToShortDateTimeString() + ".sql"

			err := os.Setenv("MYSQL_PWD", rootPassword)
			if err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.mysqlBackupFail") + ": " + err.Error())
				color.Green().Printfln(hr)
				return nil
			}

			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.targetMysql") + ": " + name)
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.startExport"))
			if _, err = shell.Execf(`mysqldump -uroot ` + name + ` > /tmp/` + backupFile + ` 2>&1`); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.exportFail") + ": " + err.Error())
				return nil
			}
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.exportSuccess"))
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.startCompress"))
			if _, err = shell.Execf("cd /tmp && zip -r " + backupFile + ".zip " + backupFile); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.compressFail") + ": " + err.Error())
				return nil
			}
			if err := io.Remove("/tmp/" + backupFile); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.deleteFail") + ": " + err.Error())
				return nil
			}
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.compressSuccess"))
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.startMove"))
			if err := io.Mv("/tmp/"+backupFile+".zip", path+"/"+backupFile+".zip"); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.moveFail") + ": " + err.Error())
				return nil
			}
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.moveSuccess"))
			_ = os.Unsetenv("MYSQL_PWD")
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.success"))

		case "postgresql":
			backupFile := name + "_" + carbon.Now().ToShortDateTimeString() + ".sql"
			check, err := shell.Execf(`su - postgres -c "psql -l" 2>&1`)
			if err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.databaseGetFail") + ": " + err.Error())
				color.Green().Printfln(hr)
				return nil
			}
			if !strings.Contains(check, name) {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.databaseNotExist"))
				color.Green().Printfln(hr)
				return nil
			}

			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.targetPostgres") + ": " + name)
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.startExport"))
			if _, err = shell.Execf(`su - postgres -c "pg_dump '` + name + `'" > /tmp/` + backupFile + ` 2>&1`); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.exportFail") + ": " + err.Error())
				return nil
			}
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.exportSuccess"))
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.startCompress"))
			if _, err = shell.Execf("cd /tmp && zip -r " + backupFile + ".zip " + backupFile); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.compressFail") + ": " + err.Error())
				return nil
			}
			if err := io.Remove("/tmp/" + backupFile); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.deleteFail") + ": " + err.Error())
				return nil
			}
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.compressSuccess"))
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.startMove"))
			if err := io.Mv("/tmp/"+backupFile+".zip", path+"/"+backupFile+".zip"); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.moveFail") + ": " + err.Error())
				return nil
			}
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.moveSuccess"))
			color.Green().Printfln("|-" + translate.Get("commands.panel.backup.success"))
		}

		color.Green().Printfln(hr)
		files, err := os.ReadDir(path)
		if err != nil {
			color.Red().Printfln("|-" + translate.Get("commands.panel.backup.cleanupFail") + ": " + err.Error())
			return nil
		}
		var filteredFiles []os.FileInfo
		for _, file := range files {
			if strings.HasPrefix(file.Name(), name) && strings.HasSuffix(file.Name(), ".zip") {
				fileInfo, err := os.Stat(filepath.Join(path, file.Name()))
				if err != nil {
					continue
				}
				filteredFiles = append(filteredFiles, fileInfo)
			}
		}
		sort.Slice(filteredFiles, func(i, j int) bool {
			return filteredFiles[i].ModTime().After(filteredFiles[j].ModTime())
		})
		for i := cast.ToInt(save); i < len(filteredFiles); i++ {
			fileToDelete := filepath.Join(path, filteredFiles[i].Name())
			color.Yellow().Printfln("|-" + translate.Get("commands.panel.backup.cleanBackup") + ": " + fileToDelete)
			if err := io.Remove(fileToDelete); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.backup.cleanupFail") + ": " + err.Error())
				return nil
			}
		}
		color.Green().Printfln("|-" + translate.Get("commands.panel.backup.cleanupSuccess"))
		color.Green().Printfln(hr)
		color.Green().Printfln("☆ " + translate.Get("commands.panel.backup.success") + " [" + carbon.Now().ToDateTimeString() + "]")
		color.Green().Printfln(hr)

	case "cutoff":
		name := arg1
		save := arg2
		hr := `+----------------------------------------------------`
		if len(name) == 0 || len(save) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.cutoff.paramFail"))
			return nil
		}

		color.Green().Printfln(hr)
		color.Green().Printfln("★ " + translate.Get("commands.panel.cutoff.start") + " [" + carbon.Now().ToDateTimeString() + "]")
		color.Green().Printfln(hr)

		color.Yellow().Printfln("|-" + translate.Get("commands.panel.cutoff.targetSite") + ": " + name)
		var website models.Website
		if err := facades.Orm().Query().Where("name", name).FirstOrFail(&website); err != nil {
			color.Red().Printfln("|-" + translate.Get("commands.panel.cutoff.siteNotExist"))
			color.Green().Printfln(hr)
			return nil
		}

		logPath := "/www/wwwlogs/" + website.Name + ".log"
		if !io.Exists(logPath) {
			color.Red().Printfln("|-" + translate.Get("commands.panel.cutoff.logNotExist"))
			color.Green().Printfln(hr)
			return nil
		}

		backupPath := "/www/wwwlogs/" + website.Name + "_" + carbon.Now().ToShortDateTimeString() + ".log.zip"
		if _, err := shell.Execf(`cd /www/wwwlogs && zip -r ` + backupPath + ` ` + website.Name + ".log"); err != nil {
			color.Red().Printfln("|-" + translate.Get("commands.panel.cutoff.backupFail") + ": " + err.Error())
			return nil
		}
		if _, err := shell.Execf(`echo "" > ` + logPath); err != nil {
			color.Red().Printfln("|-" + translate.Get("commands.panel.cutoff.clearFail") + ": " + err.Error())
			return nil
		}
		color.Green().Printfln("|-" + translate.Get("commands.panel.cutoff.cutSuccess"))

		color.Green().Printfln(hr)
		files, err := os.ReadDir("/www/wwwlogs")
		if err != nil {
			color.Red().Printfln("|-" + translate.Get("commands.panel.cutoff.cleanupFail") + ": " + err.Error())
			return nil
		}
		var filteredFiles []os.FileInfo
		for _, file := range files {
			if strings.HasPrefix(file.Name(), website.Name) && strings.HasSuffix(file.Name(), ".log.zip") {
				fileInfo, err := os.Stat(filepath.Join("/www/wwwlogs", file.Name()))
				if err != nil {
					continue
				}
				filteredFiles = append(filteredFiles, fileInfo)
			}
		}
		sort.Slice(filteredFiles, func(i, j int) bool {
			return filteredFiles[i].ModTime().After(filteredFiles[j].ModTime())
		})
		for i := cast.ToInt(save); i < len(filteredFiles); i++ {
			fileToDelete := filepath.Join("/www/wwwlogs", filteredFiles[i].Name())
			color.Yellow().Printfln("|-" + translate.Get("commands.panel.cutoff.clearLog") + ": " + fileToDelete)
			if err := io.Remove(fileToDelete); err != nil {
				color.Red().Printfln("|-" + translate.Get("commands.panel.cutoff.cleanupFail") + ": " + err.Error())
				return nil
			}
		}
		color.Green().Printfln("|-" + translate.Get("commands.panel.cutoff.cleanupSuccess"))
		color.Green().Printfln(hr)
		color.Green().Printfln("☆ " + translate.Get("commands.panel.cutoff.end") + " [" + carbon.Now().ToDateTimeString() + "]")
		color.Green().Printfln(hr)

	case "writeSite":
		name := arg1
		status := cast.ToBool(arg2)
		path := arg3
		php := cast.ToInt(arg4)
		ssl := cast.ToBool(ctx.Argument(5))
		if len(name) == 0 || len(path) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.writeSite.paramFail"))
			return nil
		}

		var website models.Website
		if err := facades.Orm().Query().Where("name", name).FirstOrFail(&website); err == nil {
			color.Red().Printfln(translate.Get("commands.panel.writeSite.siteExist"))
			return nil
		}

		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			color.Red().Printfln(translate.Get("commands.panel.writeSite.pathNotExist"))
			return nil
		}

		err = facades.Orm().Query().Create(&models.Website{
			Name:   name,
			Status: status,
			Path:   path,
			PHP:    php,
			SSL:    ssl,
		})
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.writeSite.fail"))
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.writeSite.success"))

	case "deleteSite":
		name := arg1
		if len(name) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.deleteSite.paramFail"))
			return nil
		}

		_, err := facades.Orm().Query().Where("name", name).Delete(&models.Website{})
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.deleteSite.fail"))
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.deleteSite.success"))

	case "writeSetting":
		key := arg1
		value := arg2
		if len(key) == 0 || len(value) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.writeSetting.paramFail"))
			return nil
		}

		var setting models.Setting
		err := facades.Orm().Query().UpdateOrCreate(&setting, models.Setting{
			Key: key,
		}, models.Setting{
			Value: value,
		})
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.writeSetting.fail"))
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.writeSetting.success"))

	case "getSetting":
		key := arg1
		if len(key) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.getSetting.paramFail"))
			return nil
		}

		var setting models.Setting
		if err := facades.Orm().Query().Where("key", key).FirstOrFail(&setting); err != nil {
			return nil
		}

		fmt.Printf("%s", setting.Value)

	case "deleteSetting":
		key := arg1
		if len(key) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.deleteSetting.paramFail"))
			return nil
		}

		_, err := facades.Orm().Query().Where("key", key).Delete(&models.Setting{})
		if err != nil {
			color.Red().Printfln(translate.Get("commands.panel.deleteSetting.fail"))
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.deleteSetting.success"))

	case "addSite":
		name := arg1
		domain := arg2
		port := arg3
		path := arg4
		php := arg5
		if len(name) == 0 || len(domain) == 0 || len(port) == 0 || len(path) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.addSite.paramFail"))
			return nil
		}

		domains := strings.Split(domain, ",")
		ports := strings.Split(port, ",")
		if len(domains) == 0 || len(ports) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.addSite.paramFail"))
			return nil
		}

		var uintPorts []uint
		for _, p := range ports {
			uintPorts = append(uintPorts, cast.ToUint(p))
		}

		website := services.NewWebsiteImpl()
		id, err := website.GetIDByName(name)
		if err != nil {
			color.Red().Printfln(err.Error())
			return nil
		}
		if id != 0 {
			color.Red().Printfln(translate.Get("commands.panel.addSite.siteExist"))
			return nil
		}

		_, err = website.Add(requests.Add{
			Name:    name,
			Domains: domains,
			Ports:   uintPorts,
			Path:    path,
			PHP:     php,
			DB:      false,
		})
		if err != nil {
			color.Red().Printfln(err.Error())
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.addSite.success"))

	case "removeSite":
		name := arg1
		if len(name) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.removeSite.paramFail"))
			return nil
		}

		website := services.NewWebsiteImpl()
		id, err := website.GetIDByName(name)
		if err != nil {
			color.Red().Printfln(err.Error())
			return nil
		}
		if id == 0 {
			color.Red().Printfln(translate.Get("commands.panel.removeSite.siteNotExist"))
			return nil
		}

		if err = website.Delete(requests.Delete{ID: id}); err != nil {
			color.Red().Printfln(err.Error())
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.removeSite.success"))

	case "installPlugin":
		slug := arg1
		if len(slug) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.installPlugin.paramFail"))
			return nil
		}

		plugin := services.NewPluginImpl()
		if err := plugin.Install(slug); err != nil {
			color.Red().Printfln(err.Error())
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.installPlugin.success"))

	case "uninstallPlugin":
		slug := arg1
		if len(slug) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.uninstallPlugin.paramFail"))
			return nil
		}

		plugin := services.NewPluginImpl()
		if err := plugin.Uninstall(slug); err != nil {
			color.Red().Printfln(err.Error())
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.uninstallPlugin.success"))

	case "updatePlugin":
		slug := arg1
		if len(slug) == 0 {
			color.Red().Printfln(translate.Get("commands.panel.updatePlugin.paramFail"))
			return nil
		}

		plugin := services.NewPluginImpl()
		if err := plugin.Update(slug); err != nil {
			color.Red().Printfln(err.Error())
			return nil
		}

		color.Green().Printfln(translate.Get("commands.panel.updatePlugin.success"))

	default:
		color.Yellow().Printfln(facades.Config().GetString("panel.name") + " - " + translate.Get("commands.panel.tool") + " - " + facades.Config().GetString("panel.version"))
		color.Green().Printfln(translate.Get("commands.panel.use") + "：")
		color.Green().Printfln("panel update " + translate.Get("commands.panel.update.description"))
		color.Green().Printfln("panel getInfo " + translate.Get("commands.panel.getInfo.description"))
		color.Green().Printfln("panel getPort " + translate.Get("commands.panel.getPort.description"))
		color.Green().Printfln("panel getEntrance " + translate.Get("commands.panel.getEntrance.description"))
		color.Green().Printfln("panel deleteEntrance " + translate.Get("commands.panel.deleteEntrance.description"))
		color.Green().Printfln("panel cleanTask " + translate.Get("commands.panel.cleanTask.description"))
		color.Green().Printfln("panel backup {website/mysql/postgresql} {name} {path} {save_copies} " + translate.Get("commands.panel.backup.description"))
		color.Green().Printfln("panel cutoff {website_name} {save_copies} " + translate.Get("commands.panel.cutoff.description"))
		color.Green().Printfln("panel installPlugin {slug} " + translate.Get("commands.panel.installPlugin.description"))
		color.Green().Printfln("panel uninstallPlugin {slug} " + translate.Get("commands.panel.uninstallPlugin.description"))
		color.Green().Printfln("panel updatePlugin {slug} " + translate.Get("commands.panel.updatePlugin.description"))
		color.Green().Printfln("panel addSite {name} {domain} {port} {path} {php} " + translate.Get("commands.panel.addSite.description"))
		color.Green().Printfln("panel removeSite {name} " + translate.Get("commands.panel.removeSite.description"))
		color.Red().Printfln(translate.Get("commands.panel.forDeveloper") + ":")
		color.Yellow().Printfln("panel init " + translate.Get("commands.panel.init.description"))
		color.Yellow().Printfln("panel writePlugin {slug} {version} " + translate.Get("commands.panel.writePlugin.description"))
		color.Yellow().Printfln("panel deletePlugin {slug} " + translate.Get("commands.panel.deletePlugin.description"))
		color.Yellow().Printfln("panel writeMysqlPassword {password} " + translate.Get("commands.panel.writeMysqlPassword.description"))
		color.Yellow().Printfln("panel writeSite {name} {status} {path} {php} {ssl} " + translate.Get("commands.panel.writeSite.description"))
		color.Yellow().Printfln("panel deleteSite {name} " + translate.Get("commands.panel.deleteSite.description"))
		color.Yellow().Printfln("panel getSetting {name} " + translate.Get("commands.panel.getSetting.description"))
		color.Yellow().Printfln("panel writeSetting {name} {value} " + translate.Get("commands.panel.writeSetting.description"))
		color.Yellow().Printfln("panel deleteSetting {name} " + translate.Get("commands.panel.deleteSetting.description"))
	}

	return nil
}
