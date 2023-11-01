package plugins

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"

	"panel/app/http/controllers"
	"panel/app/models"
	"panel/app/services"
	"panel/pkg/tools"
)

type Php82Controller struct {
	setting services.Setting
	task    services.Task
	version string
}

func NewPhp82Controller() *Php82Controller {
	return &Php82Controller{
		setting: services.NewSettingImpl(),
		task:    services.NewTaskImpl(),
		version: "82",
	}
}

func (r *Php82Controller) Status(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	status := tools.Exec("systemctl status php-fpm-" + r.version + " | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PHP-"+r.version+"运行状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

func (r *Php82Controller) Reload(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	tools.Exec("systemctl reload php-fpm-" + r.version)
	out := tools.Exec("systemctl status php-fpm-" + r.version + " | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PHP-"+r.version+"运行状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

func (r *Php82Controller) Start(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	tools.Exec("systemctl start php-fpm-" + r.version)
	out := tools.Exec("systemctl status php-fpm-" + r.version + " | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PHP-"+r.version+"运行状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

func (r *Php82Controller) Stop(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	tools.Exec("systemctl stop php-fpm-" + r.version)
	out := tools.Exec("systemctl status php-fpm-" + r.version + " | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PHP-"+r.version+"运行状态失败")
	}

	if status != "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

func (r *Php82Controller) Restart(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	tools.Exec("systemctl restart php-fpm-" + r.version)
	out := tools.Exec("systemctl status php-fpm-" + r.version + " | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PHP-"+r.version+"运行状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

func (r *Php82Controller) GetConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	config := tools.Read("/www/server/php/" + r.version + "/etc/php.ini")
	return controllers.Success(ctx, config)
}

func (r *Php82Controller) SaveConfig(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	config := ctx.Request().Input("config")
	tools.Write("/www/server/php/"+r.version+"/etc/php.ini", config, 0644)
	return r.Reload(ctx)
}

func (r *Php82Controller) Load(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	client := req.C().SetTimeout(10 * time.Second)
	resp, err := client.R().Get("http://127.0.0.1/phpfpm_status/" + r.version)
	if err != nil || !resp.IsSuccessState() {
		facades.Log().Error("获取PHP-" + r.version + "运行状态失败")
		return controllers.Error(ctx, http.StatusInternalServerError, "[PHP-"+r.version+"] 获取运行状态失败")
	}

	raw := resp.String()
	dataKeys := []string{"应用池", "工作模式", "启动时间", "接受连接", "监听队列", "最大监听队列", "监听队列长度", "空闲进程数量", "活动进程数量", "总进程数量", "最大活跃进程数量", "达到进程上限次数", "慢请求"}
	regexKeys := []string{"pool", "process manager", "start time", "accepted conn", "listen queue", "max listen queue", "listen queue len", "idle processes", "active processes", "total processes", "max active processes", "max children reached", "slow requests"}

	type Data struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	data := make([]Data, len(dataKeys))
	for i := range dataKeys {
		data[i].Name = dataKeys[i]

		r := regexp.MustCompile(fmt.Sprintf("%s:\\s+(.*)", regexKeys[i]))
		match := r.FindStringSubmatch(raw)

		if len(match) > 1 {
			data[i].Value = strings.TrimSpace(match[1])
		}
	}

	return controllers.Success(ctx, data)
}

func (r *Php82Controller) ErrorLog(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	log := tools.Escape(tools.Exec("tail -n 100 /www/server/php/" + r.version + "/var/log/php-fpm.log"))
	return controllers.Success(ctx, log)
}

func (r *Php82Controller) SlowLog(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	log := tools.Escape(tools.Exec("tail -n 100 /www/server/php/" + r.version + "/var/log/slow.log"))
	return controllers.Success(ctx, log)
}

func (r *Php82Controller) ClearErrorLog(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	tools.Exec("echo '' > /www/server/php/" + r.version + "/var/log/php-fpm.log")
	return controllers.Success(ctx, true)
}

func (r *Php82Controller) ClearSlowLog(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	tools.Exec("echo '' > /www/server/php/" + r.version + "/var/log/slow.log")
	return controllers.Success(ctx, true)
}

func (r *Php82Controller) GetExtensionList(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	extensions := r.GetExtensions()
	return controllers.Success(ctx, extensions)
}

func (r *Php82Controller) InstallExtension(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	slug := ctx.Request().Input("slug")
	if len(slug) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}

	extensions := r.GetExtensions()
	for _, item := range extensions {
		if item.Slug == slug {
			if item.Installed {
				return controllers.Error(ctx, http.StatusUnprocessableEntity, "扩展已安装")
			}

			var task models.Task
			task.Name = "安装PHP-" + r.version + "扩展-" + item.Name
			task.Status = models.TaskStatusWaiting
			task.Shell = `bash '/www/panel/scripts/php_extensions/` + item.Slug + `.sh' install ` + r.version + ` >> /tmp/` + item.Slug + `.log 2>&1`
			task.Log = "/tmp/" + item.Slug + ".log"
			if err := facades.Orm().Query().Create(&task); err != nil {
				facades.Log().Error("[PHP-" + r.version + "] 创建安装拓展任务失败：" + err.Error())
				return controllers.Error(ctx, http.StatusInternalServerError, "系统内部错误")
			}

			r.task.Process(task.ID)

			return controllers.Success(ctx, true)
		}
	}

	return controllers.Error(ctx, http.StatusUnprocessableEntity, "扩展不存在")
}

func (r *Php82Controller) UninstallExtension(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "php"+r.version)
	if check != nil {
		return check
	}

	slug := ctx.Request().Input("slug")
	if len(slug) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "参数错误")
	}

	extensions := r.GetExtensions()
	for _, item := range extensions {
		if item.Slug == slug {
			if !item.Installed {
				return controllers.Error(ctx, http.StatusUnprocessableEntity, "扩展未安装")
			}

			var task models.Task
			task.Name = "卸载PHP-" + r.version + "扩展-" + item.Name
			task.Status = models.TaskStatusWaiting
			task.Shell = `bash '/www/panel/scripts/php_extensions/` + item.Slug + `.sh' uninstall ` + r.version + ` >> /tmp/` + item.Slug + `.log 2>&1`
			task.Log = "/tmp/" + item.Slug + ".log"
			if err := facades.Orm().Query().Create(&task); err != nil {
				facades.Log().Error("[PHP-" + r.version + "] 创建卸载拓展任务失败：" + err.Error())
				return controllers.Error(ctx, http.StatusInternalServerError, "系统内部错误")
			}

			r.task.Process(task.ID)

			return controllers.Success(ctx, true)
		}
	}

	return controllers.Error(ctx, http.StatusUnprocessableEntity, "扩展不存在")
}

func (r *Php82Controller) GetExtensions() []PHPExtension {
	var extensions []PHPExtension
	extensions = append(extensions, PHPExtension{
		Name:        "OPcache",
		Slug:        "Zend OPcache",
		Description: "OPcache 通过将 PHP 脚本预编译的字节码存储到共享内存中来提升 PHP 的性能，存储预编译字节码可以省去每次加载和解析 PHP 脚本的开销。",
		Installed:   false,
	})
	extensions = append(extensions, PHPExtension{
		Name:        "PhpRedis",
		Slug:        "redis",
		Description: "PhpRedis 是一个用C语言编写的PHP模块，用来连接并操作 Redis 数据库上的数据。",
		Installed:   false,
	})
	extensions = append(extensions, PHPExtension{
		Name:        "ImageMagick",
		Slug:        "imagick",
		Description: "ImageMagick 是一个免费的创建、编辑、合成图片的软件。",
		Installed:   false,
	})
	extensions = append(extensions, PHPExtension{
		Name:        "Exif",
		Slug:        "exif",
		Description: "通过 exif 扩展，你可以操作图像元数据。",
		Installed:   false,
	})
	extensions = append(extensions, PHPExtension{
		Name:        "pdo_pgsql",
		Slug:        "pdo_pgsql",
		Description: "（需先安装PostgreSQL）pdo_pgsql 是一个驱动程序，它实现了 PHP 数据对象（PDO）接口以启用从 PHP 到 PostgreSQL 数据库的访问。",
		Installed:   false,
	})

	raw := tools.Exec("/www/server/php/" + r.version + "/bin/php -m")
	rawExtensionList := strings.Split(raw, "\n")

	for _, item := range rawExtensionList {
		if !strings.Contains(item, "[") && item != "" {
			for i := range extensions {
				if extensions[i].Slug == item {
					extensions[i].Installed = true
				}
			}
		}
	}

	return extensions
}
