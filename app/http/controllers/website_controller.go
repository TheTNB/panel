package controllers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/models"
	"panel/app/services"
	"panel/pkg/tools"
)

type WebsiteController struct {
	website services.Website
	setting services.Setting
	backup  services.Backup
}

func NewWebsiteController() *WebsiteController {
	return &WebsiteController{
		website: services.NewWebsiteImpl(),
		setting: services.NewSettingImpl(),
		backup:  services.NewBackupImpl(),
	}
}

// List 网站列表
func (c *WebsiteController) List(ctx http.Context) {
	limit := ctx.Request().QueryInt("limit")
	page := ctx.Request().QueryInt("page")

	total, websites, err := c.website.List(page, limit)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 获取网站列表失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, http.Json{
		"total": total,
		"items": websites,
	})
}

// Add 添加网站
func (c *WebsiteController) Add(ctx http.Context) {
	if !Check(ctx, "openresty") {
		return
	}
	validator, err := ctx.Request().Validate(map[string]string{
		"name":        "required|regex:^[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)*$|not_exists:websites,name",
		"domain":      "required",
		"php":         "required",
		"db":          "bool",
		"db_type":     "required_if:db,true",
		"db_name":     "required_if:db,true",
		"db_user":     "required_if:db,true",
		"db_password": "required_if:db,true",
	})
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	var website services.PanelWebsite
	website.Name = ctx.Request().Input("name")
	website.Domain = ctx.Request().Input("domain")
	website.Php = ctx.Request().InputInt("php")
	website.Db = ctx.Request().InputBool("db")
	website.DbType = ctx.Request().Input("db_type")
	website.DbName = ctx.Request().Input("db_name")
	website.DbUser = ctx.Request().Input("db_user")
	website.DbPassword = ctx.Request().Input("db_password")

	newSite, err := c.website.Add(website)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 添加网站失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, newSite)
}

// Delete 删除网站
func (c *WebsiteController) Delete(ctx http.Context) {
	if !Check(ctx, "openresty") {
		return
	}
	id := ctx.Request().InputInt("id")
	err := c.website.Delete(id)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 删除网站失败 ", err)
		Error(ctx, http.StatusInternalServerError, "删除网站失败: "+err.Error())
		return
	}

	Success(ctx, nil)
}

// GetDefaultConfig 获取默认配置
func (c *WebsiteController) GetDefaultConfig(ctx http.Context) {
	if !Check(ctx, "openresty") {
		return
	}
	index := tools.Read("/www/server/openresty/html/index.html")
	stop := tools.Read("/www/server/openresty/html/stop.html")

	Success(ctx, http.Json{
		"index": index,
		"stop":  stop,
	})
}

// SaveDefaultConfig 保存默认配置
func (c *WebsiteController) SaveDefaultConfig(ctx http.Context) {
	if !Check(ctx, "openresty") {
		return
	}
	index := ctx.Request().Input("index")
	stop := ctx.Request().Input("stop")

	if !tools.Write("/www/server/openresty/html/index.html", index, 0644) {
		facades.Log().Error("[面板][WebsiteController] 保存默认配置失败")
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	if !tools.Write("/www/server/openresty/html/stop.html", stop, 0644) {
		facades.Log().Error("[面板][WebsiteController] 保存默认配置失败")
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}

// GetConfig 获取配置
func (c *WebsiteController) GetConfig(ctx http.Context) {
	if !Check(ctx, "openresty") {
		return
	}
	id := ctx.Request().InputInt("id")
	if id == 0 {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	config, err := c.website.GetConfig(id)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 获取网站配置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, config)
}

// SaveConfig 保存配置
func (c *WebsiteController) SaveConfig(ctx http.Context) {
	if !Check(ctx, "openresty") {
		return
	}
	validator, err := ctx.Request().Validate(map[string]string{
		"id":                  "required",
		"domains":             "required",
		"ports":               "required",
		"hsts":                "bool",
		"ssl":                 "bool",
		"http_redirect":       "bool",
		"open_basedir":        "bool",
		"waf":                 "required",
		"waf_cache":           "required",
		"waf_mode":            "required",
		"waf_cc_deny":         "required",
		"index":               "required",
		"path":                "required",
		"root":                "required",
		"raw":                 "required",
		"php":                 "required",
		"ssl_certificate":     "required_if:ssl,true",
		"ssl_certificate_key": "required_if:ssl,true",
	})
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	var website models.Website
	if facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&website) != nil {
		Error(ctx, http.StatusBadRequest, "网站不存在")
		return
	}

	if !website.Status {
		Error(ctx, http.StatusBadRequest, "网站已停用，请先启用")
		return
	}

	// 原文
	raw := tools.Read("/www/server/vhost/" + website.Name + ".conf")
	if strings.TrimSpace(raw) != strings.TrimSpace(ctx.Request().Input("raw")) {
		tools.Write("/www/server/vhost/"+website.Name+".conf", ctx.Request().Input("raw"), 0644)
		tools.Exec("systemctl reload openresty")
		Success(ctx, nil)
		return
	}

	// 目录
	path := ctx.Request().Input("path")
	if !tools.Exists(path) {
		Error(ctx, http.StatusBadRequest, "网站目录不存在")
		return
	}
	website.Path = path

	// 域名
	domain := "server_name"
	domains := strings.Split(ctx.Request().Input("domains"), "\n")
	if len(domains) == 0 {
		Error(ctx, http.StatusBadRequest, "域名不能为空")
		return
	}
	for _, v := range domains {
		if v == "" {
			continue
		}
		domain += " " + v
	}
	domain += ";"
	domainConfigOld := tools.Cut(raw, "# server_name标记位开始", "# server_name标记位结束")
	if len(strings.TrimSpace(domainConfigOld)) == 0 {
		Error(ctx, http.StatusBadRequest, "配置文件中缺少server_name标记位")
		return
	}
	raw = strings.Replace(raw, domainConfigOld, "\n    "+domain+"\n    ", -1)

	// 端口
	var port strings.Builder
	ports := strings.Split(ctx.Request().Input("ports"), "\n")
	if len(ports) == 0 {
		Error(ctx, http.StatusBadRequest, "端口不能为空")
		return
	}
	for i, v := range ports {
		if _, err := strconv.Atoi(v); err != nil && v != "443 ssl http2" {
			Error(ctx, http.StatusBadRequest, "端口格式错误")
			return
		}
		if v == "443" && ctx.Request().InputBool("ssl") {
			v = "443 ssl http2"
		}
		if i != len(ports)-1 {
			port.WriteString("    listen " + v + ";\n")
		} else {
			port.WriteString("    listen " + v + ";")
		}
	}
	portConfigOld := tools.Cut(raw, "# port标记位开始", "# port标记位结束")
	if len(strings.TrimSpace(portConfigOld)) == 0 {
		Error(ctx, http.StatusBadRequest, "配置文件中缺少port标记位")
		return
	}
	raw = strings.Replace(raw, portConfigOld, "\n"+port.String()+"\n    ", -1)

	// 运行目录
	root := tools.Cut(raw, "# root标记位开始", "# root标记位结束")
	if len(strings.TrimSpace(root)) == 0 {
		Error(ctx, http.StatusBadRequest, "配置文件中缺少root标记位")
		return
	}
	match := regexp.MustCompile(`root\s+(.+);`).FindStringSubmatch(root)
	if len(match) != 2 {
		Error(ctx, http.StatusBadRequest, "配置文件中root标记位格式错误")
		return
	}
	rootNew := strings.Replace(root, match[1], ctx.Request().Input("root"), -1)
	raw = strings.Replace(raw, root, rootNew, -1)

	// 默认文件
	index := tools.Cut(raw, "# index标记位开始", "# index标记位结束")
	if len(strings.TrimSpace(index)) == 0 {
		Error(ctx, http.StatusBadRequest, "配置文件中缺少index标记位")
		return
	}
	match = regexp.MustCompile(`index\s+(.+);`).FindStringSubmatch(index)
	if len(match) != 2 {
		Error(ctx, http.StatusBadRequest, "配置文件中index标记位格式错误")
		return
	}
	indexNew := strings.Replace(index, match[1], ctx.Request().Input("index"), -1)
	raw = strings.Replace(raw, index, indexNew, -1)

	// 防跨站
	root = ctx.Request().Input("root")
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}
	if ctx.Request().InputBool("open_basedir") {
		tools.Write(root+".user.ini", "open_basedir="+path+":/tmp/", 0644)
	} else {
		if tools.Exists(root + ".user.ini") {
			tools.Remove(root + ".user.ini")
		}
	}

	// WAF
	waf := ctx.Request().Input("waf")
	wafMode := ctx.Request().Input("waf_mode", "DYNAMIC")
	wafCcDeny := ctx.Request().Input("waf_cc_deny", "rate=1000r/m duration=60m")
	wafCache := ctx.Request().Input("waf_cache", "capacity=50")
	wafConfig := `# waf标记位开始
    waf ` + waf + `;
    waf_rule_path /www/server/openresty/ngx_waf/assets/rules/;
    waf_mode ` + wafMode + `;
    waf_cc_deny ` + wafCcDeny + `;
    waf_cache ` + wafCache + `;
    `
	wafConfigOld := tools.Cut(raw, "# waf标记位开始", "# waf标记位结束")
	if len(strings.TrimSpace(wafConfigOld)) != 0 {
		raw = strings.Replace(raw, wafConfigOld, "", -1)
	}
	raw = strings.Replace(raw, "# waf标记位开始", wafConfig, -1)

	// SSL
	ssl := ctx.Request().InputBool("ssl")
	website.Ssl = ssl
	tools.Write("/www/server/vhost/ssl/"+website.Name+".pem", ctx.Request().Input("ssl_certificate"), 0644)
	tools.Write("/www/server/vhost/ssl/"+website.Name+".key", ctx.Request().Input("ssl_certificate_key"), 0644)
	if ssl {
		sslConfig := `# ssl标记位开始
    ssl_certificate /www/server/vhost/ssl/` + website.Name + `.pem;
    ssl_certificate_key /www/server/vhost/ssl/` + website.Name + `.key;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:10m;
    ssl_session_tickets off;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    `
		if ctx.Request().InputBool("http_redirect") {
			sslConfig += `# http重定向标记位开始
    if ($server_port !~ 443){
        return 301 https://$host$request_uri;
    }
    error_page 497  https://$host$request_uri;
    # http重定向标记位结束
    `
		}
		if ctx.Request().InputBool("hsts") {
			sslConfig += `# hsts标记位开始
    add_header Strict-Transport-Security "max-age=63072000" always;
    # hsts标记位结束
    `
		}
		sslConfigOld := tools.Cut(raw, "# ssl标记位开始", "# ssl标记位结束")
		if len(strings.TrimSpace(sslConfigOld)) != 0 {
			raw = strings.Replace(raw, sslConfigOld, "", -1)
		}
		raw = strings.Replace(raw, "# ssl标记位开始", sslConfig, -1)
	} else {
		sslConfigOld := tools.Cut(raw, "# ssl标记位开始", "# ssl标记位结束")
		if len(strings.TrimSpace(sslConfigOld)) != 0 {
			raw = strings.Replace(raw, sslConfigOld, "\n    ", -1)
		}
	}

	if website.Php != ctx.Request().InputInt("php") {
		website.Php = ctx.Request().InputInt("php")
		phpConfigOld := tools.Cut(raw, "# php标记位开始", "# php标记位结束")
		phpConfig := `
    include enable-php-` + strconv.Itoa(website.Php) + `.conf;
    `
		if len(strings.TrimSpace(phpConfigOld)) != 0 {
			raw = strings.Replace(raw, phpConfigOld, phpConfig, -1)
		}
	}

	err = facades.Orm().Query().Save(&website)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 保存网站配置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	tools.Write("/www/server/vhost/"+website.Name+".conf", raw, 0644)
	tools.Write("/www/server/vhost/rewrite/"+website.Name+".conf", ctx.Request().Input("rewrite"), 0644)
	tools.Exec("systemctl reload openresty")

	Success(ctx, nil)
}

// ClearLog 清空日志
func (c *WebsiteController) ClearLog(ctx http.Context) {
	if !Check(ctx, "openresty") {
		return
	}
	id := ctx.Request().InputInt("id")
	if id == 0 {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	website := models.Website{}
	err := facades.Orm().Query().Where("id", id).Get(&website)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 获取网站信息失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	tools.Remove("/www/wwwlogs/" + website.Name + ".log")

	Success(ctx, nil)
}

// UpdateRemark 更新备注
func (c *WebsiteController) UpdateRemark(ctx http.Context) {
	id := ctx.Request().InputInt("id")
	if id == 0 {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	website := models.Website{}
	err := facades.Orm().Query().Where("id", id).Get(&website)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 获取网站信息失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	website.Remark = ctx.Request().Input("remark")
	err = facades.Orm().Query().Save(&website)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 保存网站备注失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}

// BackupList 备份列表
func (c *WebsiteController) BackupList(ctx http.Context) {
	backupList, err := c.backup.WebsiteList()
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 获取网站备份列表失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, backupList)
}

// CreateBackup 创建备份
func (c *WebsiteController) CreateBackup(ctx http.Context) {
	id := ctx.Request().InputInt("id")
	if id == 0 {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	website := models.Website{}
	err := facades.Orm().Query().Where("id", id).Get(&website)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 获取网站信息失败 ", err)
		Error(ctx, http.StatusInternalServerError, "获取网站信息失败: "+err.Error())
		return
	}

	err = c.backup.WebSiteBackup(website)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 备份网站失败 ", err)
		Error(ctx, http.StatusInternalServerError, "备份网站失败: "+err.Error())
		return
	}

	Success(ctx, nil)
}

// UploadBackup 上传备份
func (c *WebsiteController) UploadBackup(ctx http.Context) {
	file, err := ctx.Request().File("file")
	if err != nil {
		Error(ctx, http.StatusBadRequest, "上传文件失败")
		return
	}

	backupPath := c.setting.Get(models.SettingKeyBackupPath) + "/website"
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}

	name := file.GetClientOriginalName()
	_, err = file.StoreAs(backupPath, name)
	if err != nil {
		Error(ctx, http.StatusBadRequest, "上传文件失败")
		return
	}

	Success(ctx, "上传文件成功")
}

// RestoreBackup 还原备份
func (c *WebsiteController) RestoreBackup(ctx http.Context) {
	id := ctx.Request().InputInt("id")
	if id == 0 {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}
	fileName := ctx.Request().Input("name")
	if len(fileName) == 0 {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	website := models.Website{}
	err := facades.Orm().Query().Where("id", id).Get(&website)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 获取网站信息失败 ", err)
		Error(ctx, http.StatusInternalServerError, "获取网站信息失败: "+err.Error())
		return
	}

	err = c.backup.WebsiteRestore(website, fileName)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 还原网站失败 ", err)
		Error(ctx, http.StatusInternalServerError, "还原网站失败: "+err.Error())
		return
	}

	Success(ctx, nil)
}

// DeleteBackup 删除备份
func (c *WebsiteController) DeleteBackup(ctx http.Context) {
	fileName := ctx.Request().Input("name")
	if len(fileName) == 0 {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	backupPath := c.setting.Get(models.SettingKeyBackupPath) + "/website"
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}

	if !tools.Remove(backupPath + "/" + fileName) {
		Error(ctx, http.StatusInternalServerError, "删除备份失败")
		return
	}

	Success(ctx, nil)
}

// ResetConfig 重置配置
func (c *WebsiteController) ResetConfig(ctx http.Context) {
	if !Check(ctx, "openresty") {
		return
	}
	id := ctx.Request().InputInt("id")
	if id == 0 {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	website := models.Website{}
	if err := facades.Orm().Query().Where("id", id).Get(&website); err != nil {
		facades.Log().Error("[面板][WebsiteController] 获取网站信息失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	website.Status = true
	website.Ssl = false
	if err := facades.Orm().Query().Save(&website); err != nil {
		facades.Log().Error("[面板][WebsiteController] 保存网站配置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	raw := fmt.Sprintf(`
# 配置文件中的标记位请勿随意修改，改错将导致面板无法识别！
# 有自定义配置需求的，请将自定义的配置写在各标记位下方。
server
{
    # port标记位开始
    listen 80;
    # port标记位结束
    # server_name标记位开始
    server_name localhost;
    # server_name标记位结束
    # index标记位开始
    index index.php index.html;
    # index标记位结束
    # root标记位开始
    root %s;
    # root标记位结束

    # ssl标记位开始
    # ssl标记位结束

    # php标记位开始
    include enable-php-%d.conf;
    # php标记位结束

    # waf标记位开始
    waf on;
    waf_rule_path /www/server/openresty/ngx_waf/assets/rules/;
    waf_mode DYNAMIC;
    waf_cc_deny rate=1000r/m duration=60m;
    waf_cache capacity=50;
    # waf标记位结束

    # 错误页配置，可自行设置
    #error_page 404 /404.html;
    #error_page 502 /502.html;

    # 伪静态规则引入，修改后将导致面板设置的伪静态规则失效
    include /www/server/vhost/rewrite/%s.conf;

    # 面板默认禁止访问部分敏感目录，可自行修改
    location ~ ^/(\.user.ini|\.htaccess|\.git|\.svn)
    {
        return 404;
    }
    # 面板默认不记录静态资源的访问日志并开启1小时浏览器缓存，可自行修改
    location ~ .*\.(js|css)$
    {
        expires 1h;
        error_log /dev/null;
        access_log /dev/null;
    }

    access_log /www/wwwlogs/%s.log;
    error_log /www/wwwlogs/%s.log;
}

`, website.Path, website.Php, website.Name, website.Name, website.Name)

	tools.Write("/www/server/vhost/"+website.Name+".conf", raw, 0644)
	tools.Write("/www/server/vhost/rewrite"+website.Name+".conf", "", 0644)
	tools.Exec("systemctl reload openresty")

	Success(ctx, nil)
}

// Status 网站状态
func (c *WebsiteController) Status(ctx http.Context) {
	if !Check(ctx, "openresty") {
		return
	}
	id := ctx.Request().InputInt("id")
	if id == 0 {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	website := models.Website{}
	if err := facades.Orm().Query().Where("id", id).Get(&website); err != nil {
		facades.Log().Error("[面板][WebsiteController] 获取网站信息失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	website.Status = ctx.Request().InputBool("status")
	if err := facades.Orm().Query().Save(&website); err != nil {
		facades.Log().Error("[面板][WebsiteController] 保存网站配置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	raw := tools.Read("/www/server/vhost/" + website.Name + ".conf")

	// 运行目录
	rootConfig := tools.Cut(raw, "# root标记位开始\n", "# root标记位结束")
	match := regexp.MustCompile(`root\s+(.+);`).FindStringSubmatch(rootConfig)
	if len(match) == 2 {
		if website.Status {
			root := regexp.MustCompile(`# root\s+(.+);`).FindStringSubmatch(rootConfig)
			raw = strings.ReplaceAll(raw, rootConfig, "    root "+root[1]+";\n    ")
		} else {
			raw = strings.ReplaceAll(raw, rootConfig, "    root /www/server/openresty/html;\n    # root "+match[1]+";\n    ")
		}
	}

	// 默认文件
	indexConfig := tools.Cut(raw, "# index标记位开始\n", "# index标记位结束")
	match = regexp.MustCompile(`index\s+(.+);`).FindStringSubmatch(indexConfig)
	if len(match) == 2 {
		if website.Status {
			index := regexp.MustCompile(`# index\s+(.+);`).FindStringSubmatch(indexConfig)
			raw = strings.ReplaceAll(raw, indexConfig, "    index "+index[1]+";\n    ")
		} else {
			raw = strings.ReplaceAll(raw, indexConfig, "    index stop.html;\n    # index "+match[1]+";\n    ")
		}
	}

	tools.Write("/www/server/vhost/"+website.Name+".conf", raw, 0644)
	tools.Exec("systemctl reload openresty")

	Success(ctx, nil)
}
