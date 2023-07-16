package commands

import (
	"os"
	"regexp"

	"github.com/gookit/color"
	"github.com/spf13/cast"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"

	"panel/app/models"
	"panel/app/services"
	"panel/packages/helper"
)

type Panel struct {
}

// Signature The name and signature of the console command.
func (receiver *Panel) Signature() string {
	return "panel"
}

// Description The console command description.
func (receiver *Panel) Description() string {
	return "[面板] 命令行"
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

	switch action {
	case "init":
		var check models.User
		err := facades.Orm().Query().FirstOrFail(&check)
		if err == nil {
			color.Redln("面板已初始化")
			return nil
		}

		settings := []models.Setting{{Key: "name", Value: "耗子Linux面板"}, {Key: "monitor", Value: "1"}, {Key: "monitor_days", Value: "30"}, {Key: "backup_path", Value: "/www/backup"}, {Key: "website_path", Value: "/www/wwwroot"}, {Key: "panel_entrance", Value: "/"}}
		err = facades.Orm().Query().Create(&settings)
		if err != nil {
			color.Redln("初始化失败")
			return nil
		}

		hash, err := facades.Hash().Make(helper.RandomString(32))
		if err != nil {
			color.Redln("初始化失败")
			return nil
		}

		user := services.NewUserImpl()
		_, err = user.Create("admin", hash)
		if err != nil {
			color.Redln("创建管理员失败")
			return nil
		}

		color.Greenln("初始化成功")

	case "update":
		err := helper.UpdatePanel()
		if err != nil {
			color.Redln("更新失败: " + err.Error())
			return nil
		}

		color.Greenln("更新成功")

	case "getInfo":
		var user models.User
		err := facades.Orm().Query().Where("id", 1).FirstOrFail(&user)
		if err != nil {
			color.Redln("获取管理员信息失败")
			return nil
		}

		password := helper.RandomString(16)
		hash, err := facades.Hash().Make(password)
		if err != nil {
			color.Redln("生成密码失败")
			return nil
		}

		user.Username = helper.RandomString(8)
		user.Password = hash

		err = facades.Orm().Query().Save(&user)
		if err != nil {
			color.Redln("保存管理员信息失败")
			return nil
		}

		color.Greenln("用户名: " + user.Username)
		color.Greenln("密码: " + password)
		// color.Greenln("面板端口: " + port)
		color.Greenln("面板入口: " + services.NewSettingImpl().Get("panel_entrance", "/"))

	case "getPort":
		nginxConf, err := os.ReadFile("/www/server/nginx/conf/nginx.conf")
		if err != nil {
			color.Redln("获取面板端口失败，请检查Nginx主配置文件")
			return nil
		}

		match := regexp.MustCompile(`listen\s+(\d+)`).FindStringSubmatch(string(nginxConf))
		if len(match) < 2 {
			color.Redln("获取面板端口失败，请检查Nginx主配置文件")
			return nil
		}

		port := match[1]
		color.Greenln("面板端口: " + port)

	case "getEntrance":
		color.Greenln("面板入口: " + services.NewSettingImpl().Get("panel_entrance", "/"))

	case "writePlugin":
		slug := arg1
		version := arg2
		if len(slug) == 0 || len(version) == 0 {
			color.Redln("参数错误")
			return nil
		}

		var plugin models.Plugin
		err := facades.Orm().Query().UpdateOrCreate(&plugin, models.Plugin{
			Slug: slug,
		}, models.Plugin{
			Version: version,
		})

		if err != nil {
			color.Redln("写入插件安装状态失败")
			return nil
		}

		color.Greenln("写入插件安装状态成功")

	case "deletePlugin":
		slug := arg1
		if len(slug) == 0 {
			color.Redln("参数错误")
			return nil
		}

		_, err := facades.Orm().Query().Where("slug", slug).Delete(&models.Plugin{})
		if err != nil {
			color.Redln("移除插件安装状态失败")
			return nil
		}

		color.Greenln("移除插件安装状态成功")

	case "writeMysqlPassword":
		password := arg1
		if len(password) == 0 {
			color.Redln("参数错误")
			return nil
		}

		var setting models.Setting
		err := facades.Orm().Query().UpdateOrCreate(&setting, models.Setting{
			Key: "mysql_root_password",
		}, models.Setting{
			Value: password,
		})

		if err != nil {
			color.Redln("写入MySQL root密码失败")
			return nil
		}

		color.Greenln("写入MySQL root密码成功")

	case "cleanRunningTask":
		_, err := facades.Orm().Query().Model(&models.Task{}).Where("status", models.TaskStatusRunning).Update("status", models.TaskStatusFailed)
		if err != nil {
			color.Redln("清理正在运行的任务失败")
			return nil
		}

		color.Greenln("清理正在运行的任务成功")

	case "backup":

	case "writeSite":
		name := arg1
		status := cast.ToBool(arg2)
		path := ctx.Argument(3)
		php := cast.ToInt(ctx.Argument(4))
		ssl := cast.ToBool(ctx.Argument(5))
		if len(name) == 0 || len(path) == 0 {
			color.Redln("参数错误")
			return nil
		}

		var website models.Website
		if err := facades.Orm().Query().Where("name", name).FirstOrFail(&website); err == nil {
			color.Redln("网站已存在")
			return nil
		}

		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			color.Redln("网站目录不存在")
			return nil
		}

		err = facades.Orm().Query().Create(&models.Website{
			Name:   name,
			Status: status,
			Path:   path,
			Php:    php,
			Ssl:    ssl,
		})
		if err != nil {
			color.Redln("写入网站失败")
			return nil
		}

		color.Greenln("写入网站成功")

	case "deleteSite":
		name := arg1
		if len(name) == 0 {
			color.Redln("参数错误")
			return nil
		}

		_, err := facades.Orm().Query().Where("name", name).Delete(&models.Website{})
		if err != nil {
			color.Redln("删除网站失败")
			return nil
		}

		color.Greenln("删除网站成功")

	case "writeSetting":
		key := arg1
		value := arg2
		if len(key) == 0 || len(value) == 0 {
			color.Redln("参数错误")
			return nil
		}

		var setting models.Setting
		err := facades.Orm().Query().UpdateOrCreate(&setting, models.Setting{
			Key: key,
		}, models.Setting{
			Value: value,
		})
		if err != nil {
			color.Redln("写入设置失败")
			return nil
		}

		color.Greenln("写入设置成功")

	case "deleteSetting":
		key := arg1
		if len(key) == 0 {
			color.Redln("参数错误")
			return nil
		}

		_, err := facades.Orm().Query().Where("key", key).Delete(&models.Setting{})
		if err != nil {
			color.Redln("删除设置失败")
			return nil
		}

		color.Greenln("删除设置成功")

	default:
		color.Yellowln(facades.Config().GetString("panel.name") + "命令行工具 - " + facades.Config().GetString("panel.version"))
		color.Greenln("请使用以下命令：")
		color.Greenln("panel update 更新/修复面板到最新版本")
		color.Greenln("panel getInfo 重新初始化面板账号信息")
		color.Greenln("panel getPort 获取面板访问端口")
		color.Greenln("panel getEntrance 获取面板访问入口")
		color.Greenln("panel cleanRunningTask 强制清理面板正在运行的任务")
		color.Greenln("panel backup {website/mysql/postgresql} {name} {path} 备份网站/MySQL数据库/PostgreSQL数据库到指定目录")
		color.Redln("以下命令请在开发者指导下使用：")
		color.Yellowln("panel init 初始化面板")
		color.Yellowln("panel writePlugin {slug} 写入插件安装状态")
		color.Yellowln("panel deletePlugin {slug} 移除插件安装状态")
		color.Yellowln("panel writeMysqlPassword {password} 写入MySQL root密码")
		color.Yellowln("panel writeSite {name} {status} {path} {php} {ssl} 写入网站数据到面板")
		color.Yellowln("panel deleteSite {name} 删除面板网站数据")
		color.Yellowln("panel writeSetting {name} {value} 写入/更新面板设置数据")
		color.Yellowln("panel deleteSetting {name} 删除面板设置数据")
	}

	return nil
}
