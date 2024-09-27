package php

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/panel"
	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/types"
)

type Service struct {
	version  uint
	taskRepo biz.TaskRepo
}

func NewService(version uint) *Service {
	return &Service{
		version:  version,
		taskRepo: data.NewTaskRepo(),
	}
}

// GetConfig
//
//	@Summary	获取配置
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/config [get]
func (s *Service) GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read(fmt.Sprintf("%s/server/php/%d/etc/php.ini", panel.Root, s.version))
	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, config)
}

// UpdateConfig
//
//	@Summary	保存配置
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int		true	"PHP 版本"
//	@Param		config	body		string	true	"配置"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/config [post]
func (s *Service) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = io.Write(fmt.Sprintf("%s/server/php/%d/etc/php.ini", panel.Root, s.version), req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
}

// GetFPMConfig
//
//	@Summary	获取 FPM 配置
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/fpmConfig [get]
func (s *Service) GetFPMConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read(fmt.Sprintf("%s/server/php/%d/etc/php-fpm.conf", panel.Root, s.version))
	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, config)
}

// UpdateFPMConfig
//
//	@Summary	保存 FPM 配置
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int		true	"PHP 版本"
//	@Param		config	body		string	true	"配置"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/fpmConfig [post]
func (s *Service) UpdateFPMConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = io.Write(fmt.Sprintf("%s/server/php/%d/etc/php-fpm.conf", panel.Root, s.version), req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
}

// Load
//
//	@Summary	获取负载状态
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/load [get]
func (s *Service) Load(w http.ResponseWriter, r *http.Request) {
	client := resty.New().SetTimeout(10 * time.Second)
	resp, err := client.R().Get(fmt.Sprintf("http://127.0.0.1/phpfpm_status/%d", s.version))
	if err != nil || !resp.IsSuccess() {
		service.Error(w, http.StatusInternalServerError, "获取负载状态失败")
		return
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

	service.Success(w, data)
}

// ErrorLog
//
//	@Summary	获取错误日志
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/errorLog [get]
func (s *Service) ErrorLog(w http.ResponseWriter, r *http.Request) {
	log, _ := shell.Execf("tail -n 500 %s/server/php/%d/var/log/php-fpm.log", panel.Root, s.version)
	service.Success(w, log)
}

// SlowLog
//
//	@Summary	获取慢日志
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/slowLog [get]
func (s *Service) SlowLog(w http.ResponseWriter, r *http.Request) {
	log, _ := shell.Execf("tail -n 500 %s/server/php/%d/var/log/slow.log", panel.Root, s.version)
	service.Success(w, log)
}

// ClearErrorLog
//
//	@Summary	清空错误日志
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/clearErrorLog [post]
func (s *Service) ClearErrorLog(w http.ResponseWriter, r *http.Request) {
	if out, err := shell.Execf("echo '' > %s/server/php/%d/var/log/php-fpm.log", panel.Root, s.version); err != nil {
		service.Error(w, http.StatusInternalServerError, out)
		return
	}

	service.Success(w, nil)
}

// ClearSlowLog
//
//	@Summary	清空慢日志
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/clearSlowLog [post]
func (s *Service) ClearSlowLog(w http.ResponseWriter, r *http.Request) {
	if out, err := shell.Execf("echo '' > %s/server/php/%d/var/log/slow.log", panel.Root, s.version); err != nil {
		service.Error(w, http.StatusInternalServerError, out)
		return
	}

	service.Success(w, nil)
}

// ExtensionList
//
//	@Summary	获取扩展列表
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int	true	"PHP 版本"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/extensions [get]
func (s *Service) ExtensionList(w http.ResponseWriter, r *http.Request) {
	extensions := s.getExtensions()

	// ionCube 只支持 PHP 8.3 以下版本
	if cast.ToUint(s.version) < 83 {
		extensions = append(extensions, types.PHPExtension{
			Name:        "ionCube",
			Slug:        "ionCube Loader",
			Description: "ionCube 是一个专业级的 PHP 加密解密工具。",
			Installed:   false,
		})
	}

	raw, err := shell.Execf("%s/server/php/%d/bin/php -m", panel.Root, s.version)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
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

	service.Success(w, extensions)
}

// InstallExtension
//
//	@Summary	安装扩展
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int		true	"PHP 版本"
//	@Param		slug	query		string	true	"slug"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/extensions [post]
func (s *Service) InstallExtension(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ExtensionSlug](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if !s.checkExtension(req.Slug) {
		service.Error(w, http.StatusUnprocessableEntity, "扩展不存在")
		return
	}

	cmd := fmt.Sprintf(`bash '%s/panel/scripts/php_extensions/%s.sh' install %d >> '/tmp/%s.log' 2>&1`, panel.Root, req.Slug, s.version, req.Slug)
	officials := []string{"fileinfo", "exif", "imap", "pdo_pgsql", "zip", "bz2", "readline", "snmp", "ldap", "enchant", "pspell", "calendar", "gmp", "sysvmsg", "sysvsem", "sysvshm", "xsl", "intl", "gettext"}
	if slices.Contains(officials, req.Slug) {
		cmd = fmt.Sprintf(`bash '%s/panel/scripts/php_extensions/official.sh' install '%d' '%s' >> '/tmp/%s.log' 2>&1`, panel.Root, s.version, req.Slug, req.Slug)
	}

	task := new(biz.Task)
	task.Name = fmt.Sprintf("安装 PHP-%d 扩展 %s", s.version, req.Slug)
	task.Status = biz.TaskStatusWaiting
	task.Shell = cmd
	task.Log = "/tmp/" + req.Slug + ".log"
	if err = s.taskRepo.Push(task); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
}

// UninstallExtension
//
//	@Summary	卸载扩展
//	@Tags		插件-PHP
//	@Produce	json
//	@Security	BearerToken
//	@Param		version	path		int		true	"PHP 版本"
//	@Param		slug	query		string	true	"slug"
//	@Success	200		{object}	controllers.SuccessResponse
//	@Router		/plugins/php/{version}/extensions [delete]
func (s *Service) UninstallExtension(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[ExtensionSlug](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if !s.checkExtension(req.Slug) {
		service.Error(w, http.StatusUnprocessableEntity, "扩展不存在")
		return
	}

	cmd := fmt.Sprintf(`bash '%s/panel/scripts/php_extensions/%s.sh' uninstall %d >> '/tmp/%s.log' 2>&1`, panel.Root, req.Slug, s.version, req.Slug)
	officials := []string{"fileinfo", "exif", "imap", "pdo_pgsql", "zip", "bz2", "readline", "snmp", "ldap", "enchant", "pspell", "calendar", "gmp", "sysvmsg", "sysvsem", "sysvshm", "xsl", "intl", "gettext"}
	if slices.Contains(officials, req.Slug) {
		cmd = fmt.Sprintf(`bash '%s/panel/scripts/php_extensions/official.sh' uninstall '%d' '%s' >> '/tmp/%s.log' 2>&1`, panel.Root, s.version, req.Slug, req.Slug)
	}

	task := new(biz.Task)
	task.Name = fmt.Sprintf("卸载 PHP-%d 扩展 %s", s.version, req.Slug)
	task.Status = biz.TaskStatusWaiting
	task.Shell = cmd
	task.Log = "/tmp/" + req.Slug + ".log"
	if err = s.taskRepo.Push(task); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
}

func (s *Service) getExtensions() []types.PHPExtension {
	extensions := []types.PHPExtension{
		{
			Name:        "fileinfo",
			Slug:        "fileinfo",
			Description: "Fileinfo 是一个用于识别文件类型的库。",
		},
		{
			Name:        "OPcache",
			Slug:        "Zend OPcache",
			Description: "OPcache 通过将 PHP 脚本预编译的字节码存储到共享内存中来提升 PHP 的性能，存储预编译字节码可以省去每次加载和解析 PHP 脚本的开销。",
		},
		{
			Name:        "PhpRedis",
			Slug:        "redis",
			Description: "PhpRedis 是一个用 C 语言编写的 PHP 模块，用来连接并操作 Redis 数据库上的数据。",
		},
		{
			Name:        "ImageMagick",
			Slug:        "imagick",
			Description: "ImageMagick 是一个免费的创建、编辑、合成图片的软件。",
		},
		{
			Name:        "exif",
			Slug:        "exif",
			Description: "通过 exif 扩展，你可以操作图像元数据。",
		},
		{
			Name:        "pdo_pgsql",
			Slug:        "pdo_pgsql",
			Description: "（需先安装PostgreSQL）pdo_pgsql 是一个驱动程序，它实现了 PHP 数据对象（PDO）接口以启用从 PHP 到 PostgreSQL 数据库的访问。",
		},
		{
			Name:        "imap",
			Slug:        "imap",
			Description: "IMAP 扩展允许 PHP 读取、搜索、删除、下载和管理邮件。",
		},
		{
			Name:        "zip",
			Slug:        "zip",
			Description: "Zip 是一个用于处理 ZIP 文件的库。",
		},
		{
			Name:        "bz2",
			Slug:        "bz2",
			Description: "Bzip2 是一个用于压缩和解压缩文件的库。",
		},
		{
			Name:        "readline",
			Slug:        "readline",
			Description: "Readline 是一个库，它提供了一种用于处理文本的接口。",
		},
		{
			Name:        "snmp",
			Slug:        "snmp",
			Description: "SNMP 是一种用于网络管理的协议。",
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
		},
		{
			Name:        "pspell",
			Slug:        "pspell",
			Description: "Pspell 是一个拼写检查库。",
		},
		{
			Name:        "calendar",
			Slug:        "calendar",
			Description: "Calendar 是一个用于处理日期的库。",
		},
		{
			Name:        "gmp",
			Slug:        "gmp",
			Description: "GMP 是一个用于处理大整数的库。",
		},
		{
			Name:        "sysvmsg",
			Slug:        "sysvmsg",
			Description: "Sysvmsg 是一个用于处理 System V 消息队列的库。",
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
		},
		{
			Name:        "xsl",
			Slug:        "xsl",
			Description: "XSL 是一个用于处理 XML 文档的库。",
		},
		{
			Name:        "intl",
			Slug:        "intl",
			Description: "Intl 是一个用于处理国际化和本地化的库。",
		},
		{
			Name:        "gettext",
			Slug:        "gettext",
			Description: "Gettext 是一个用于处理多语言的库。",
		},
		{
			Name:        "igbinary",
			Slug:        "igbinary",
			Description: "Igbinary 是一个用于序列化和反序列化数据的库。",
		},
		{
			Name:        "Swoole",
			Slug:        "swoole",
			Description: "Swoole 是一个用于构建高性能的异步并发服务器的 PHP 扩展。",
		},
		{
			Name:        "Swow",
			Slug:        "Swow",
			Description: "Swow 是一个用于构建高性能的异步并发服务器的 PHP 扩展。",
		},
	}

	// ionCube 只支持 PHP 8.3 以下版本
	if cast.ToUint(s.version) < 83 {
		extensions = append(extensions, types.PHPExtension{
			Name:        "ionCube",
			Slug:        "ionCube Loader",
			Description: "ionCube 是一个专业级的 PHP 加密解密工具。",
			Installed:   false,
		})
	}

	raw, _ := shell.Execf("%s/server/php/%d/bin/php -m", panel.Root, s.version)
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

	return extensions
}

func (s *Service) checkExtension(slug string) bool {
	extensions := s.getExtensions()

	for _, item := range extensions {
		if item.Slug == slug {
			return true
		}
	}

	return false
}
