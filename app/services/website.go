// Package services 网站服务
package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goravel/framework/facades"
	"golang.org/x/exp/slices"

	"panel/app/models"
	"panel/packages/helper"
)

type Website interface {
	List() ([]models.Website, error)
}

type PanelWebsite struct {
	Name       string `json:"name"`
	Status     bool   `json:"status"`
	Domain     string `json:"domain"`
	Path       string `json:"path"`
	Php        int    `json:"php"`
	Ssl        bool   `json:"ssl"`
	Remark     string `json:"remark"`
	Db         bool   `json:"db"`
	DbType     string `json:"db_type"`
	DbName     string `json:"db_name"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
}

// WebsiteSetting 网站设置
type WebsiteSetting struct {
	Name              string   `json:"name"`
	Ports             []string `json:"ports"`
	Domains           []string `json:"domains"`
	Root              string   `json:"root"`
	Path              string   `json:"path"`
	Index             string   `json:"index"`
	OpenBasedir       bool     `json:"open_basedir"`
	Ssl               bool     `json:"ssl"`
	SslCertificate    string   `json:"ssl_certificate"`
	SslCertificateKey string   `json:"ssl_certificate_key"`
	HttpRedirect      bool     `json:"http_redirect"`
	Hsts              bool     `json:"hsts"`
	Waf               bool     `json:"waf"`
	WafMode           string   `json:"waf_mode"`
	WafCcDeny         string   `json:"waf_cc_deny"`
	WafCache          string   `json:"waf_cache"`
	Rewrite           string   `json:"rewrite"`
	Raw               string   `json:"raw"`
	Log               string   `json:"log"`
}

type WebsiteImpl struct {
}

func NewWebsiteImpl() *WebsiteImpl {
	return &WebsiteImpl{}
}

// List 列出网站
func (r *WebsiteImpl) List(page, limit int) (int64, []models.Website, error) {
	var websites []models.Website
	var total int64
	if err := facades.Orm().Query().Paginate(page, limit, &websites, &total); err != nil {
		return total, websites, err
	}

	return total, websites, nil
}

// Add 添加网站
func (r *WebsiteImpl) Add(website PanelWebsite) (models.Website, error) {
	// 禁止部分保留名称
	nameSlices := []string{"phpmyadmin", "mysql", "panel", "ssh"}
	if slices.Contains(nameSlices, website.Name) {
		return models.Website{}, errors.New("网站名称" + website.Name + "为保留名称，请更换")
	}

	// path为空时，设置默认值
	if len(website.Path) == 0 {
		website.Path = "/www/wwwroot/" + website.Name
	}
	// path不为/开头时，返回错误
	if website.Path[0] != '/' {
		return models.Website{}, errors.New("网站路径" + website.Path + "必须以/开头")
	}

	website.Ssl = false
	website.Status = true
	website.Domain = strings.TrimSpace(website.Domain)

	w := models.Website{
		Name:   website.Name,
		Status: website.Status,
		Path:   website.Path,
		Php:    website.Php,
		Ssl:    website.Ssl,
		Remark: website.Remark,
	}
	if err := facades.Orm().Query().Create(&w); err != nil {
		return w, err
	}

	helper.Mkdir(website.Path, 0755)

	index := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<title>耗子Linux面板</title>
</head>
<body>
<h1>耗子Linux面板</h1>
<p>这是耗子Linux面板的网站默认页面！</p>
<p>当您看到此页面，说明您的网站已创建成功。</p>
</body>
</html>

`
	helper.WriteFile(website.Path+"/index.html", index, 0644)

	domainArr := strings.Split(website.Domain, "\n")
	portList := ""
	portArr := make(map[string]bool)
	domainList := ""
	for key, value := range domainArr {
		temp := strings.Split(value, ":")
		domainList += " " + temp[0]

		if len(temp) < 2 {
			if _, ok := portArr["80"]; !ok {
				if key == len(domainArr)-1 {
					portList += "    listen 80;"
				} else {
					portList += "    listen 80;\n"
				}
				portArr["80"] = true
			}
		} else {
			if _, ok := portArr[temp[1]]; !ok {
				if key == len(domainArr)-1 {
					portList += "    listen " + temp[1] + ";"
				} else {
					portList += "    listen " + temp[1] + ";\n"
				}
				portArr[temp[1]] = true
			}
		}
	}

	nginxConf := fmt.Sprintf(`
# 配置文件中的标记位请勿随意修改，改错将导致面板无法识别！
# 有自定义配置需求的，请将自定义的配置写在各标记位下方。
server
{
    # port标记位开始
%s
    # port标记位结束
    # server_name标记位开始
    server_name%s;
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
    include /www/server/vhost/openresty/rewrite/%s.conf;

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

`, portList, domainList, website.Path, website.Php, website.Name, website.Name, website.Name)

	helper.WriteFile("/www/server/panel/vhost/openresty/"+website.Name+".conf", nginxConf, 0644)
	helper.WriteFile("/www/server/panel/vhost/openresty/rewrite/"+website.Name+".conf", "", 0644)
	helper.WriteFile("/www/server/panel/vhost/openresty/ssl/"+website.Name+".pem", "", 0644)
	helper.WriteFile("/www/server/panel/vhost/openresty/ssl/"+website.Name+".key", "", 0644)

	helper.ExecShellAsync("systemctl reload openresty")

	// TODO 创建数据库

	return w, nil
}

// Delete 删除网站
func (r *WebsiteImpl) Delete(name string) error {
	var website models.Website
	if err := facades.Orm().Query().Where("name", name).First(&website); err != nil {
		return err
	}

	if _, err := facades.Orm().Query().Delete(&website); err != nil {
		return err
	}

	helper.RemoveFile("/www/server/panel/vhost/openresty/" + website.Name + ".conf")
	helper.RemoveFile("/www/server/panel/vhost/openresty/rewrite/" + website.Name + ".conf")
	helper.RemoveFile("/www/server/panel/vhost/openresty/ssl/" + website.Name + ".pem")
	helper.RemoveFile("/www/server/panel/vhost/openresty/ssl/" + website.Name + ".key")
	helper.RemoveFile(website.Path)

	helper.ExecShellAsync("systemctl reload openresty")

	// TODO 删除数据库

	return nil
}
