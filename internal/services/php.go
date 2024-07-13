package services

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type PHPImpl struct {
	version string
}

func NewPHPImpl(version uint) *PHPImpl {
	return &PHPImpl{
		version: cast.ToString(version),
	}
}

func (r *PHPImpl) Reload() error {
	return systemctl.Reload("php-fpm-" + r.version)
}

func (r *PHPImpl) GetConfig() (string, error) {
	return io.Read("/www/server/php/" + r.version + "/etc/php.ini")
}

func (r *PHPImpl) SaveConfig(config string) error {
	if err := io.Write("/www/server/php/"+r.version+"/etc/php.ini", config, 0644); err != nil {
		return err
	}

	return r.Reload()
}

func (r *PHPImpl) GetFPMConfig() (string, error) {
	return io.Read("/www/server/php/" + r.version + "/etc/php-fpm.conf")
}

func (r *PHPImpl) SaveFPMConfig(config string) error {
	if err := io.Write("/www/server/php/"+r.version+"/etc/php-fpm.conf", config, 0644); err != nil {
		return err
	}

	return r.Reload()
}

func (r *PHPImpl) Load() ([]types.NV, error) {
	client := resty.New().SetTimeout(10 * time.Second)
	resp, err := client.R().Get("http://127.0.0.1/phpfpm_status/" + r.version)
	if err != nil || !resp.IsSuccess() {
		return []types.NV{}, nil
	}

	raw := resp.String()
	dataKeys := []string{"应用池", "工作模式", "启动时间", "接受连接", "监听队列", "最大监听队列", "监听队列长度", "空闲进程数量", "活动进程数量", "总进程数量", "最大活跃进程数量", "达到进程上限次数", "慢请求"}
	regexKeys := []string{"pool", "process manager", "start time", "accepted conn", "listen queue", "max listen queue", "listen queue len", "idle processes", "active processes", "total processes", "max active processes", "max children reached", "slow requests"}

	data := make([]types.NV, len(dataKeys))
	for i := range dataKeys {
		data[i].Name = dataKeys[i]

		r := regexp.MustCompile(fmt.Sprintf("%s:\\s+(.*)", regexKeys[i]))
		match := r.FindStringSubmatch(raw)

		if len(match) > 1 {
			data[i].Value = strings.TrimSpace(match[1])
		}
	}

	return data, nil
}

func (r *PHPImpl) GetErrorLog() (string, error) {
	return shell.Execf("tail -n 500 /www/server/php/%s/var/log/php-fpm.log", r.version)
}

func (r *PHPImpl) GetSlowLog() (string, error) {
	return shell.Execf("tail -n 500 /www/server/php/%s/var/log/slow.log", r.version)
}

func (r *PHPImpl) ClearErrorLog() error {
	if out, err := shell.Execf("echo '' > /www/server/php/%s/var/log/php-fpm.log", r.version); err != nil {
		return errors.New(out)
	}

	return r.Reload()
}

func (r *PHPImpl) ClearSlowLog() error {
	if out, err := shell.Execf("echo '' > /www/server/php/%s/var/log/slow.log", r.version); err != nil {
		return errors.New(out)
	}

	return nil
}

func (r *PHPImpl) GetExtensions() ([]types.PHPExtension, error) {
	extensions := []types.PHPExtension{
		{
			Name:        "fileinfo",
			Slug:        "fileinfo",
			Description: "Fileinfo 是一个用于识别文件类型的库。",
			Installed:   false,
		},
		{
			Name:        "OPcache",
			Slug:        "Zend OPcache",
			Description: "OPcache 通过将 PHP 脚本预编译的字节码存储到共享内存中来提升 PHP 的性能，存储预编译字节码可以省去每次加载和解析 PHP 脚本的开销。",
			Installed:   false,
		},
		{
			Name:        "PhpRedis",
			Slug:        "redis",
			Description: "PhpRedis 是一个用 C 语言编写的 PHP 模块，用来连接并操作 Redis 数据库上的数据。",
			Installed:   false,
		},
		{
			Name:        "ImageMagick",
			Slug:        "imagick",
			Description: "ImageMagick 是一个免费的创建、编辑、合成图片的软件。",
			Installed:   false,
		},
		{
			Name:        "exif",
			Slug:        "exif",
			Description: "通过 exif 扩展，你可以操作图像元数据。",
			Installed:   false,
		},
		{
			Name:        "pdo_pgsql",
			Slug:        "pdo_pgsql",
			Description: "（需先安装PostgreSQL）pdo_pgsql 是一个驱动程序，它实现了 PHP 数据对象（PDO）接口以启用从 PHP 到 PostgreSQL 数据库的访问。",
			Installed:   false,
		},
		{
			Name:        "imap",
			Slug:        "imap",
			Description: "IMAP 扩展允许 PHP 读取、搜索、删除、下载和管理邮件。",
			Installed:   false,
		},
		{
			Name:        "zip",
			Slug:        "zip",
			Description: "Zip 是一个用于处理 ZIP 文件的库。",
			Installed:   false,
		},
		{
			Name:        "bz2",
			Slug:        "bz2",
			Description: "Bzip2 是一个用于压缩和解压缩文件的库。",
			Installed:   false,
		},
		{
			Name:        "readline",
			Slug:        "readline",
			Description: "Readline 是一个库，它提供了一种用于处理文本的接口。",
			Installed:   false,
		},
		{
			Name:        "snmp",
			Slug:        "snmp",
			Description: "SNMP 是一种用于网络管理的协议。",
			Installed:   false,
		},
		{
			Name:        "ldap",
			Slug:        "ldap",
			Description: "LDAP 是一种用于访问目录服务的协议。",
		},
		{
			Name:        "enchant",
			Slug:        "enchant",
			Description: "Enchant 是一个拼写检查库。",
			Installed:   false,
		},
		{
			Name:        "pspell",
			Slug:        "pspell",
			Description: "Pspell 是一个拼写检查库。",
			Installed:   false,
		},
		{
			Name:        "calendar",
			Slug:        "calendar",
			Description: "Calendar 是一个用于处理日期的库。",
			Installed:   false,
		},
		{
			Name:        "gmp",
			Slug:        "gmp",
			Description: "GMP 是一个用于处理大整数的库。",
			Installed:   false,
		},
		{
			Name:        "sysvmsg",
			Slug:        "sysvmsg",
			Description: "Sysvmsg 是一个用于处理 System V 消息队列的库。",
			Installed:   false,
		},
		{
			Name:        "sysvsem",
			Slug:        "sysvsem",
			Description: "Sysvsem 是一个用于处理 System V 信号量的库。",
		},
		{
			Name:        "sysvshm",
			Slug:        "sysvshm",
			Description: "Sysvshm 是一个用于处理 System V 共享内存的库。",
			Installed:   false,
		},
		{
			Name:        "xsl",
			Slug:        "xsl",
			Description: "XSL 是一个用于处理 XML 文档的库。",
			Installed:   false,
		},
		{
			Name:        "intl",
			Slug:        "intl",
			Description: "Intl 是一个用于处理国际化和本地化的库。",
			Installed:   false,
		},
		{
			Name:        "gettext",
			Slug:        "gettext",
			Description: "Gettext 是一个用于处理多语言的库。",
			Installed:   false,
		},
		{
			Name:        "igbinary",
			Slug:        "igbinary",
			Description: "Igbinary 是一个用于序列化和反序列化数据的库。",
			Installed:   false,
		},
	}

	// ionCube 只支持 PHP 8.3 以下版本
	if cast.ToUint(r.version) < 83 {
		extensions = append(extensions, types.PHPExtension{
			Name:        "ionCube",
			Slug:        "ionCube Loader",
			Description: "ionCube 是一个专业级的 PHP 加密解密工具。",
			Installed:   false,
		})
	}
	// Swoole 和 Swow 不支持 PHP 8.0 以下版本
	if cast.ToUint(r.version) >= 80 {
		extensions = append(extensions, types.PHPExtension{
			Name:        "Swoole",
			Slug:        "swoole",
			Description: "Swoole 是一个用于构建高性能的异步并发服务器的 PHP 扩展。",
			Installed:   false,
		})
		extensions = append(extensions, types.PHPExtension{
			Name:        "Swow",
			Slug:        "Swow",
			Description: "Swow 是一个用于构建高性能的异步并发服务器的 PHP 扩展。",
			Installed:   false,
		})
	}

	raw, err := shell.Execf("/www/server/php/%s/bin/php -m", r.version)
	if err != nil {
		return extensions, err
	}

	extensionMap := make(map[string]*types.PHPExtension)
	for i := range extensions {
		extensionMap[extensions[i].Slug] = &extensions[i]
	}

	rawExtensionList := strings.Split(raw, "\n")
	for _, item := range rawExtensionList {
		if ext, exists := extensionMap[item]; exists && !strings.Contains(item, "[") && item != "" {
			ext.Installed = true
		}
	}

	return extensions, nil
}

func (r *PHPImpl) InstallExtension(slug string) error {
	if !r.checkExtension(slug) {
		return errors.New("扩展不存在")
	}

	shell := fmt.Sprintf(`bash '/www/panel/scripts/php_extensions/%s.sh' install %s >> '/tmp/%s.log' 2>&1`, slug, r.version, slug)

	officials := []string{"fileinfo", "exif", "imap", "pdo_pgsql", "zip", "bz2", "readline", "snmp", "ldap", "enchant", "pspell", "calendar", "gmp", "sysvmsg", "sysvsem", "sysvshm", "xsl", "intl", "gettext"}
	if slices.Contains(officials, slug) {
		shell = fmt.Sprintf(`bash '/www/panel/scripts/php_extensions/official.sh' install '%s' '%s' >> '/tmp/%s.log' 2>&1`, r.version, slug, slug)
	}

	var task models.Task
	task.Name = "安装PHP-" + r.version + "扩展-" + slug
	task.Status = models.TaskStatusWaiting
	task.Shell = shell
	task.Log = "/tmp/" + slug + ".log"
	if err := facades.Orm().Query().Create(&task); err != nil {
		return err
	}

	return NewTaskImpl().Process(task.ID)
}

func (r *PHPImpl) UninstallExtension(slug string) error {
	if !r.checkExtension(slug) {
		return errors.New("扩展不存在")
	}

	shell := fmt.Sprintf(`bash '/www/panel/scripts/php_extensions/%s.sh' uninstall %s >> '/tmp/%s.log' 2>&1`, slug, r.version, slug)

	officials := []string{"fileinfo", "exif", "imap", "pdo_pgsql", "zip", "bz2", "readline", "snmp", "ldap", "enchant", "pspell", "calendar", "gmp", "sysvmsg", "sysvsem", "sysvshm", "xsl", "intl", "gettext"}
	if slices.Contains(officials, slug) {
		shell = fmt.Sprintf(`bash '/www/panel/scripts/php_extensions/official.sh' uninstall '%s' '%s' >> '/tmp/%s.log' 2>&1`, r.version, slug, slug)
	}

	var task models.Task
	task.Name = "卸载PHP-" + r.version + "扩展-" + slug
	task.Status = models.TaskStatusWaiting
	task.Shell = shell
	task.Log = "/tmp/" + slug + ".log"
	if err := facades.Orm().Query().Create(&task); err != nil {
		return err
	}

	return NewTaskImpl().Process(task.ID)
}

func (r *PHPImpl) checkExtension(slug string) bool {
	extensions, err := r.GetExtensions()
	if err != nil {
		return false
	}

	for _, item := range extensions {
		if item.Slug == slug {
			return true
		}
	}

	return false
}
