package php82

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

func (c *Php82Controller) Status(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	status := tools.ExecShell("systemctl status php-fpm-" + c.version + " | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PHP-"+c.version+"运行状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

func (c *Php82Controller) Reload(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	tools.ExecShell("systemctl reload php-fpm-" + c.version)
	out := tools.ExecShell("systemctl status php-fpm-" + c.version + " | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PHP-"+c.version+"运行状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

func (c *Php82Controller) Start(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	tools.ExecShell("systemctl start php-fpm-" + c.version)
	out := tools.ExecShell("systemctl status php-fpm-" + c.version + " | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PHP-"+c.version+"运行状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

func (c *Php82Controller) Stop(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	tools.ExecShell("systemctl stop php-fpm-" + c.version)
	out := tools.ExecShell("systemctl status php-fpm-" + c.version + " | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PHP-"+c.version+"运行状态失败")
		return
	}

	if status != "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

func (c *Php82Controller) Restart(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	tools.ExecShell("systemctl restart php-fpm-" + c.version)
	out := tools.ExecShell("systemctl status php-fpm-" + c.version + " | grep Active | grep -v grep | awk '{print $2}'")
	status := strings.TrimSpace(out)
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PHP-"+c.version+"运行状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

func (c *Php82Controller) GetConfig(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	config := tools.ReadFile("/www/server/php/" + c.version + "/etc/php.ini")
	controllers.Success(ctx, config)
}

func (c *Php82Controller) SaveConfig(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	config := ctx.Request().Input("config")
	tools.WriteFile("/www/server/php/"+c.version+"/etc/php.ini", config, 0644)
	c.Reload(ctx)
}

func (c *Php82Controller) Load(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	client := req.C().SetTimeout(10 * time.Second)
	resp, err := client.R().Get("http://127.0.0.1/phpfpm_status/" + c.version)
	if err != nil || !resp.IsSuccessState() {
		facades.Log().Error("获取PHP-" + c.version + "运行状态失败")
		controllers.Error(ctx, http.StatusInternalServerError, "[PHP-"+c.version+"] 获取运行状态失败")
		return
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

	controllers.Success(ctx, data)
}

func (c *Php82Controller) ErrorLog(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	log := tools.ExecShell("tail -n 100 /www/server/php/" + c.version + "/var/log/php-fpm.log")
	controllers.Success(ctx, log)
}

func (c *Php82Controller) SlowLog(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	log := tools.ExecShell("tail -n 100 /www/server/php/" + c.version + "/var/log/slow.log")
	controllers.Success(ctx, log)
}

func (c *Php82Controller) ClearErrorLog(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	tools.ExecShell("echo '' > /www/server/php/" + c.version + "/var/log/php-fpm.log")
	controllers.Success(ctx, true)
}

func (c *Php82Controller) ClearSlowLog(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	tools.ExecShell("echo '' > /www/server/php/" + c.version + "/var/log/slow.log")
	controllers.Success(ctx, true)
}

type Extension struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Installed   bool   `json:"installed"`
}

func (c *Php82Controller) GetExtensionList(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	extensions := c.GetExtensions()
	controllers.Success(ctx, extensions)
}

func (c *Php82Controller) InstallExtension(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	slug := ctx.Request().Input("slug")
	if len(slug) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	extensions := c.GetExtensions()
	for _, item := range extensions {
		if item.Slug == slug {
			if item.Installed {
				controllers.Error(ctx, http.StatusBadRequest, "扩展已安装")
				return
			}

			var task models.Task
			task.Name = "安装PHP-" + c.version + "扩展-" + item.Name
			task.Status = models.TaskStatusWaiting
			task.Shell = `bash '/www/panel/scripts/php_extensions/` + item.Slug + `.sh' install ` + c.version + ` >> /tmp/` + item.Slug + `.log 2>&1`
			task.Log = "/tmp/" + item.Slug + ".log"
			if err := facades.Orm().Query().Create(&task); err != nil {
				facades.Log().Error("[PHP-" + c.version + "] 创建安装拓展任务失败：" + err.Error())
				controllers.Error(ctx, http.StatusInternalServerError, "系统内部错误")
				return
			}

			c.task.Process(task.ID)

			controllers.Success(ctx, true)
			return
		}
	}

	controllers.Error(ctx, http.StatusBadRequest, "扩展不存在")
}

func (c *Php82Controller) UninstallExtension(ctx http.Context) {
	if !controllers.Check(ctx, "php"+c.version) {
		return
	}

	slug := ctx.Request().Input("slug")
	if len(slug) == 0 {
		controllers.Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	extensions := c.GetExtensions()
	for _, item := range extensions {
		if item.Slug == slug {
			if !item.Installed {
				controllers.Error(ctx, http.StatusBadRequest, "扩展未安装")
				return
			}

			var task models.Task
			task.Name = "卸载PHP-" + c.version + "扩展-" + item.Name
			task.Status = models.TaskStatusWaiting
			task.Shell = `bash '/www/panel/scripts/php_extensions/` + item.Slug + `.sh' uninstall ` + c.version + ` >> /tmp/` + item.Slug + `.log 2>&1`
			task.Log = "/tmp/" + item.Slug + ".log"
			if err := facades.Orm().Query().Create(&task); err != nil {
				facades.Log().Error("[PHP-" + c.version + "] 创建卸载拓展任务失败：" + err.Error())
				controllers.Error(ctx, http.StatusInternalServerError, "系统内部错误")
				return
			}

			c.task.Process(task.ID)

			controllers.Success(ctx, true)
			return
		}
	}

	controllers.Error(ctx, http.StatusBadRequest, "扩展不存在")
}

func (c *Php82Controller) GetExtensions() []Extension {
	var extensions []Extension
	extensions = append(extensions, Extension{
		Name:        "OPcache",
		Slug:        "Zend OPcache",
		Description: "OPcache 通过将 PHP 脚本预编译的字节码存储到共享内存中来提升 PHP 的性能，存储预编译字节码可以省去每次加载和解析 PHP 脚本的开销。",
		Installed:   false,
	})
	extensions = append(extensions, Extension{
		Name:        "PhpRedis",
		Slug:        "redis",
		Description: "PhpRedis 是一个用C语言编写的PHP模块，用来连接并操作 Redis 数据库上的数据。",
		Installed:   false,
	})
	extensions = append(extensions, Extension{
		Name:        "ImageMagick",
		Slug:        "imagick",
		Description: "ImageMagick 是一个免费的创建、编辑、合成图片的软件。",
		Installed:   false,
	})
	extensions = append(extensions, Extension{
		Name:        "Exif",
		Slug:        "exif",
		Description: "通过 exif 扩展，你可以操作图像元数据。",
		Installed:   false,
	})
	extensions = append(extensions, Extension{
		Name:        "pdo_pgsql",
		Slug:        "pdo_pgsql",
		Description: "（需先安装PostgreSQL）pdo_pgsql 是一个驱动程序，它实现了 PHP 数据对象（PDO）接口以启用从 PHP 到 PostgreSQL 数据库的访问。",
		Installed:   false,
	})

	raw := tools.ExecShell("/www/server/php/" + c.version + "/bin/php -m")
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
