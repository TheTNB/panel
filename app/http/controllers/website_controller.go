package controllers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"panel/app/models"
	"panel/pkg/tools"

	"panel/app/services"
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

func (c *WebsiteController) Add(ctx http.Context) {
	validator, err := ctx.Request().Validate(map[string]string{
		"name":        "required|regex:^[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)*$",
		"domain":      "required",
		"php":         "required",
		"db":          "required",
		"db_type":     "required_if:db,1",
		"db_name":     "required_if:db,1",
		"db_user":     "required_if:db,1",
		"db_password": "required_if:db,1",
	})
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	var website services.PanelWebsite
	err = ctx.Request().Bind(&website)
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	newSite, err := c.website.Add(website)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 添加网站失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, newSite)
}

func (c *WebsiteController) Delete(ctx http.Context) {
	validator, err := ctx.Request().Validate(map[string]string{
		"id": "required|int",
	})
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	id := ctx.Request().InputInt("id")
	err = c.website.Delete(id)
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 删除网站失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}

func (c *WebsiteController) GetDefaultConfig(ctx http.Context) {
	index := tools.ReadFile("/www/server/openresty/html/index.html")
	stop := tools.ReadFile("/www/server/openresty/html/stop.html")

	Success(ctx, http.Json{
		"index": index,
		"stop":  stop,
	})
}

func (c *WebsiteController) SaveDefaultConfig(ctx http.Context) {
	index := ctx.Request().Input("index")
	stop := ctx.Request().Input("stop")

	if !tools.WriteFile("/www/server/openresty/html/index.html", index, 0644) {
		facades.Log().Error("[面板][WebsiteController] 保存默认配置失败")
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	if !tools.WriteFile("/www/server/openresty/html/stop.html", stop, 0644) {
		facades.Log().Error("[面板][WebsiteController] 保存默认配置失败")
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}

func (c *WebsiteController) GetConfig(ctx http.Context) {
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

func (c *WebsiteController) SaveConfig(ctx http.Context) {
	validator, err := ctx.Request().Validate(map[string]string{
		"id": "required|int",
	})
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	id := ctx.Request().InputInt("id")
	var website models.Website
	if facades.Orm().Query().Where("id", id).FirstOrFail(&website) != nil {
		Error(ctx, http.StatusBadRequest, "网站不存在")
		return
	}

	if !website.Status {
		Error(ctx, http.StatusBadRequest, "网站已停用，请先启用")
		return
	}

	// 原文
	raw := tools.ReadFile("/www/server/panel/vhost/openresty/" + website.Name + ".conf")
	if strings.TrimSpace(raw) != strings.TrimSpace(ctx.Request().Input("raw")) {
		tools.WriteFile("/www/server/panel/vhost/openresty/"+website.Name+".conf", ctx.Request().Input("raw"), 0644)
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
	domains := strings.Split(ctx.Request().Input("domain"), "\n")
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
	ports := strings.Split(ctx.Request().Input("port"), "\n")
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
	rootNew := strings.Replace(root, match[1], path, -1)
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
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if ctx.Request().InputBool("open_basedir") {
		tools.WriteFile(path+".user.ini", "open_basedir="+path+":/tmp/", 0644)
	} else {
		if tools.Exists(path + ".user.ini") {
			tools.RemoveFile(path + ".user.ini")
		}
	}

	// WAF
	waf := ctx.Request().Input("waf")
	wafMode := ctx.Request().Input("waf_mode", "DYNAMIC")
	wafCcDeny := ctx.Request().Input("waf_cc_deny", "rate=1000r/m duration=60m")
	wafCache := ctx.Request().Input("waf_cache", "capacity=50")
	wafConfig := `
# waf标记位开始
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
	if ssl {
		tools.WriteFile("/www/server/vhost/ssl/"+website.Name+".pem", ctx.Request().Input("ssl_certificate"), 0644)
		tools.WriteFile("/www/server/vhost/ssl/"+website.Name+".key", ctx.Request().Input("ssl_certificate_key"), 0644)
		sslConfig := `
# ssl标记位开始
    ssl_certificate /www/server/vhost/ssl/` + website.Name + `.pem;
    ssl_certificate_key /www/server/vhost/ssl/` + website.Name + `.key;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:10m;
    ssl_session_tickets off;
    ssl_dhparam /etc/ssl/certs/dhparam.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

`
		if ctx.Request().InputBool("http_redirect") {
			sslConfig += `

    # http重定向标记位开始
    if (\$server_port !~ 443){
        return 301 https://\$host\$request_uri;
    }
    error_page 497  https://\$host\$request_uri;
    # http重定向标记位结束

`
		}
		if ctx.Request().InputBool("hsts") {
			sslConfig += `

    # hsts标记位开始
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
			raw = strings.Replace(raw, sslConfigOld, "", -1)
		}
	}

	if website.Php != ctx.Request().InputInt("php") {
		website.Php = ctx.Request().InputInt("php")
		phpConfigOld := tools.Cut(raw, "# php标记位开始", "# php标记位结束")
		phpConfig := `
    include enable-php` + strconv.Itoa(website.Php) + `.conf;
    
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

	tools.WriteFile("/www/server/vhost/"+website.Name+".conf", raw, 0644)
	tools.WriteFile("/www/server/vhost/rewrite/"+website.Name+".conf", ctx.Request().Input("rewrite"), 0644)
	tools.ExecShell("systemctl reload openresty")

	Success(ctx, nil)
}

func (c *WebsiteController) ClearSiteLpg(ctx http.Context) {
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

	tools.RemoveFile("/www/wwwlogs/" + website.Name + ".log")

	Success(ctx, nil)
}

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

func (c *WebsiteController) BackupList(ctx http.Context) {

	backupList, err := c.backup.WebsiteList()
	if err != nil {
		facades.Log().Error("[面板][WebsiteController] 获取网站备份列表失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, backupList)
}

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

func (c *WebsiteController) ResetConfig(ctx http.Context) {
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
    waf_rule_path /www/server/nginx/ngx_waf/assets/rules/;
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

	tools.WriteFile("/www/server/vhost/"+website.Name+".conf", raw, 0644)
	tools.WriteFile("/www/server/vhost/rewrite"+website.Name+".conf", "", 0644)
	tools.ExecShell("systemctl reload openresty")

	Success(ctx, nil)
}

func (c *WebsiteController) SetStatus(ctx http.Context) {
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

	raw := tools.ReadFile("/www/server/vhost/" + website.Name + ".conf")

	// 运行目录
	rootConfig := tools.Cut(raw, "# root标记位开始", "# root标记位结束")
	match := regexp.MustCompile(`root\s+(.+);`).FindStringSubmatch(rootConfig)
	if len(match) == 2 {
		if website.Status {
			root := regexp.MustCompile(`# root\s+(.+);`).FindStringSubmatch(rootConfig)
			raw = strings.ReplaceAll(raw, rootConfig, "root "+root[1]+";")
		} else {
			raw = strings.ReplaceAll(raw, rootConfig, "root /www/server/openresty/html;\n# root "+match[1]+";\n")
		}
	}

	// 默认文件
	indexConfig := tools.Cut(raw, "# index标记位开始", "# index标记位结束")
	match = regexp.MustCompile(`index\s+(.+);`).FindStringSubmatch(indexConfig)
	if len(match) == 2 {
		if website.Status {
			index := regexp.MustCompile(`# index\s+(.+);`).FindStringSubmatch(indexConfig)
			raw = strings.ReplaceAll(raw, indexConfig, "index "+index[1]+";")
		} else {
			raw = strings.ReplaceAll(raw, indexConfig, "index stop.html;\n# index "+match[1]+";\n")
		}
	}

	tools.WriteFile("/www/server/vhost/"+website.Name+".conf", raw, 0644)
	tools.ExecShell("systemctl reload openresty")

	Success(ctx, nil)
}
