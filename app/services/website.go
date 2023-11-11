// Package services 网站服务
package services

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/goravel/framework/facades"
	"golang.org/x/exp/slices"
	requests "panel/app/http/requests/website"

	"panel/app/models"
	"panel/pkg/tools"
)

type Website interface {
	List(page int, limit int) (int64, []models.Website, error)
	Add(website PanelWebsite) (models.Website, error)
	SaveConfig(config requests.SaveConfig) error
	Delete(id uint) error
	GetConfig(id uint) (WebsiteSetting, error)
	GetConfigByName(name string) (WebsiteSetting, error)
}

type PanelWebsite struct {
	Name       string   `json:"name"`
	Status     bool     `json:"status"`
	Domains    []string `json:"domains"`
	Ports      []string `json:"ports"`
	Path       string   `json:"path"`
	Php        int      `json:"php"`
	Ssl        bool     `json:"ssl"`
	Remark     string   `json:"remark"`
	Db         bool     `json:"db"`
	DbType     string   `json:"db_type"`
	DbName     string   `json:"db_name"`
	DbUser     string   `json:"db_user"`
	DbPassword string   `json:"db_password"`
}

// WebsiteSetting 网站设置
type WebsiteSetting struct {
	Name              string   `json:"name"`
	Domains           []string `json:"domains"`
	Ports             []string `json:"ports"`
	Root              string   `json:"root"`
	Path              string   `json:"path"`
	Index             string   `json:"index"`
	Php               string   `json:"php"`
	OpenBasedir       bool     `json:"open_basedir"`
	Ssl               bool     `json:"ssl"`
	SslCertificate    string   `json:"ssl_certificate"`
	SslCertificateKey string   `json:"ssl_certificate_key"`
	SslNotBefore      string   `json:"ssl_not_before"`
	SslNotAfter       string   `json:"ssl_not_after"`
	SSlDNSNames       []string `json:"ssl_dns_names"`
	SslIssuer         string   `json:"ssl_issuer"`
	SslOCSPServer     []string `json:"ssl_ocsp_server"`
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
	setting Setting
}

func NewWebsiteImpl() *WebsiteImpl {
	return &WebsiteImpl{
		setting: NewSettingImpl(),
	}
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
		website.Path = r.setting.Get(models.SettingKeyWebsitePath) + "/" + website.Name
	}
	// path不为/开头时，返回错误
	if website.Path[0] != '/' {
		return models.Website{}, errors.New("网站路径" + website.Path + "必须以/开头")
	}

	website.Ssl = false
	website.Status = true

	w := models.Website{
		Name:   website.Name,
		Status: website.Status,
		Path:   website.Path,
		Php:    website.Php,
		Ssl:    website.Ssl,
		Remark: website.Remark,
	}
	if err := facades.Orm().Query().Create(&w); err != nil {
		return models.Website{}, err
	}

	if err := tools.Mkdir(website.Path, 0755); err != nil {
		return models.Website{}, err
	}

	index := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>耗子Linux面板</title>
    <style>
        body {
            background-color: #f9f9f9;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 800px;
            margin: 2em auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 12px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        h1 {
            font-size: 2.5em;
            margin-top: 0;
            margin-bottom: 20px;
            text-align: center;
            color: #333;
            border-bottom: 2px solid #ddd;
            padding-bottom: 0.5em;
        }
        p {
            color: #555;
            line-height: 1.8;
        }
        @media screen and (max-width: 768px) {
            .container {
                padding: 15px;
                margin: 2em 15px;
            }
            h1 {
                font-size: 1.8em;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>耗子Linux面板</h1>
        <p>这是耗子Linux面板的网站默认页面！</p>
        <p>当您看到此页面，说明您的网站已创建成功。</p>
    </div>
</body>
</html>

`
	if err := tools.Write(website.Path+"/index.html", index, 0644); err != nil {
		return models.Website{}, err
	}

	portList := ""
	domainList := ""
	portUsed := make(map[string]bool)
	domainUsed := make(map[string]bool)

	for i, port := range website.Ports {
		if _, ok := portUsed[port]; !ok {
			if i == len(website.Ports)-1 {
				portList += "    listen " + port + ";"
			} else {
				portList += "    listen " + port + ";\n"
			}
			portUsed[port] = true
		}
	}
	if len(website.Ports) == 0 {
		portList += "    listen 80;\n"
	}
	for _, domain := range website.Domains {
		if _, ok := domainUsed[domain]; !ok {
			domainList += " " + domain
			domainUsed[domain] = true
		}
	}

	nginxConf := fmt.Sprintf(`# 配置文件中的标记位请勿随意修改，改错将导致面板无法识别！
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
`, portList, domainList, website.Path, website.Php, website.Name, website.Name, website.Name)

	if err := tools.Write("/www/server/vhost/"+website.Name+".conf", nginxConf, 0644); err != nil {
		return models.Website{}, err
	}
	if err := tools.Write("/www/server/vhost/rewrite/"+website.Name+".conf", "", 0644); err != nil {
		return models.Website{}, err
	}
	if err := tools.Write("/www/server/vhost/ssl/"+website.Name+".pem", "", 0644); err != nil {
		return models.Website{}, err
	}
	if err := tools.Write("/www/server/vhost/ssl/"+website.Name+".key", "", 0644); err != nil {
		return models.Website{}, err
	}

	if err := tools.Chmod(r.setting.Get(models.SettingKeyWebsitePath), 0755); err != nil {
		return models.Website{}, err
	}
	if err := tools.Chmod(website.Path, 0755); err != nil {
		return models.Website{}, err
	}
	if err := tools.Chown(r.setting.Get(models.SettingKeyWebsitePath), "www", "www"); err != nil {
		return models.Website{}, err
	}
	if err := tools.Chown(website.Path, "www", "www"); err != nil {
		return models.Website{}, err
	}

	tools.Exec("systemctl reload openresty")

	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	if website.Db && website.DbType == "mysql" {
		tools.Exec(`/www/server/mysql/bin/mysql -uroot -p` + rootPassword + ` -e "CREATE DATABASE IF NOT EXISTS ` + website.DbName + ` DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;"`)
		tools.Exec(`/www/server/mysql/bin/mysql -uroot -p` + rootPassword + ` -e "CREATE USER '` + website.DbUser + `'@'localhost' IDENTIFIED BY '` + website.DbPassword + `';"`)
		tools.Exec(`/www/server/mysql/bin/mysql -uroot -p` + rootPassword + ` -e "GRANT ALL PRIVILEGES ON ` + website.DbName + `.* TO '` + website.DbUser + `'@'localhost';"`)
		tools.Exec(`/www/server/mysql/bin/mysql -uroot -p` + rootPassword + ` -e "FLUSH PRIVILEGES;"`)
	}
	if website.Db && website.DbType == "postgresql" {
		tools.Exec(`echo "CREATE DATABASE ` + website.DbName + `;" | su - postgres -c "psql"`)
		tools.Exec(`echo "CREATE USER ` + website.DbUser + ` WITH PASSWORD '` + website.DbPassword + `';" | su - postgres -c "psql"`)
		tools.Exec(`echo "ALTER DATABASE ` + website.DbName + ` OWNER TO ` + website.DbUser + `;" | su - postgres -c "psql"`)
		tools.Exec(`echo "GRANT ALL PRIVILEGES ON DATABASE ` + website.DbName + ` TO ` + website.DbUser + `;" | su - postgres -c "psql"`)
	}

	return w, nil
}

// SaveConfig 保存网站配置
func (r *WebsiteImpl) SaveConfig(config requests.SaveConfig) error {
	var website models.Website
	if err := facades.Orm().Query().Where("id", config.ID).First(&website); err != nil {
		return err
	}

	if !website.Status {
		return errors.New("网站已停用，请先启用")
	}

	// 原文
	raw := tools.Read("/www/server/vhost/" + website.Name + ".conf")
	if strings.TrimSpace(raw) != strings.TrimSpace(config.Raw) {
		if err := tools.Write("/www/server/vhost/"+website.Name+".conf", config.Raw, 0644); err != nil {
			return err
		}
		tools.Exec("systemctl reload openresty")
		return nil
	}

	// 目录
	path := config.Path
	if !tools.Exists(path) {
		return errors.New("网站目录不存在")
	}
	website.Path = path

	// 域名
	domain := "server_name"
	domains := config.Domains
	for _, v := range domains {
		if v == "" {
			continue
		}
		domain += " " + v
	}
	domain += ";"
	domainConfigOld := tools.Cut(raw, "# server_name标记位开始", "# server_name标记位结束")
	if len(strings.TrimSpace(domainConfigOld)) == 0 {
		return errors.New("配置文件中缺少server_name标记位")
	}
	raw = strings.Replace(raw, domainConfigOld, "\n    "+domain+"\n    ", -1)

	// 端口
	var port strings.Builder
	ports := config.Ports
	for i, v := range ports {
		if _, err := strconv.Atoi(v); err != nil && v != "443 ssl http2" {
			return errors.New("端口格式错误")
		}
		if v == "443" && config.Ssl {
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
		return errors.New("配置文件中缺少port标记位")
	}
	raw = strings.Replace(raw, portConfigOld, "\n"+port.String()+"\n    ", -1)

	// 运行目录
	root := tools.Cut(raw, "# root标记位开始", "# root标记位结束")
	if len(strings.TrimSpace(root)) == 0 {
		return errors.New("配置文件中缺少root标记位")
	}
	match := regexp.MustCompile(`root\s+(.+);`).FindStringSubmatch(root)
	if len(match) != 2 {
		return errors.New("配置文件中root标记位格式错误")
	}
	rootNew := strings.Replace(root, match[1], config.Root, -1)
	raw = strings.Replace(raw, root, rootNew, -1)

	// 默认文件
	index := tools.Cut(raw, "# index标记位开始", "# index标记位结束")
	if len(strings.TrimSpace(index)) == 0 {
		return errors.New("配置文件中缺少index标记位")
	}
	match = regexp.MustCompile(`index\s+(.+);`).FindStringSubmatch(index)
	if len(match) != 2 {
		return errors.New("配置文件中index标记位格式错误")
	}
	indexNew := strings.Replace(index, match[1], config.Index, -1)
	raw = strings.Replace(raw, index, indexNew, -1)

	// 防跨站
	root = config.Root
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}
	if config.OpenBasedir {
		if err := tools.Write(root+".user.ini", "open_basedir="+path+":/tmp/", 0644); err != nil {
			return err
		}
	} else {
		if tools.Exists(root + ".user.ini") {
			tools.Remove(root + ".user.ini")
		}
	}

	// WAF
	waf := config.Waf
	wafStr := "off"
	if waf {
		wafStr = "on"
	}
	wafMode := config.WafMode
	wafCcDeny := config.WafCcDeny
	wafCache := config.WafCache
	wafConfig := `# waf标记位开始
    waf ` + wafStr + `;
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
	ssl := config.Ssl
	website.Ssl = ssl
	if err := tools.Write("/www/server/vhost/ssl/"+website.Name+".pem", config.SslCertificate, 0644); err != nil {
		return err
	}
	if err := tools.Write("/www/server/vhost/ssl/"+website.Name+".key", config.SslCertificateKey, 0644); err != nil {
		return err
	}
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
		if config.HttpRedirect {
			sslConfig += `# http重定向标记位开始
    if ($server_port !~ 443){
        return 301 https://$host$request_uri;
    }
    error_page 497  https://$host$request_uri;
    # http重定向标记位结束
    `
		}
		if config.Hsts {
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

	if website.Php != config.Php {
		website.Php = config.Php
		phpConfigOld := tools.Cut(raw, "# php标记位开始", "# php标记位结束")
		phpConfig := `
    include enable-php-` + strconv.Itoa(website.Php) + `.conf;
    `
		if len(strings.TrimSpace(phpConfigOld)) != 0 {
			raw = strings.Replace(raw, phpConfigOld, phpConfig, -1)
		}
	}

	if err := facades.Orm().Query().Save(&website); err != nil {
		return err
	}

	if err := tools.Write("/www/server/vhost/"+website.Name+".conf", raw, 0644); err != nil {
		return err
	}
	if err := tools.Write("/www/server/vhost/rewrite/"+website.Name+".conf", config.Rewrite, 0644); err != nil {
		return err
	}
	tools.Exec("systemctl reload openresty")
	return nil
}

// Delete 删除网站
func (r *WebsiteImpl) Delete(id uint) error {
	var website models.Website
	if err := facades.Orm().Query().With("Cert").Where("id", id).FirstOrFail(&website); err != nil {
		return err
	}

	if website.Cert != nil {
		return errors.New("网站" + website.Name + "已绑定SSL证书，请先删除证书")
	}

	if _, err := facades.Orm().Query().Delete(&website); err != nil {
		return err
	}

	tools.Remove("/www/server/vhost/" + website.Name + ".conf")
	tools.Remove("/www/server/vhost/rewrite/" + website.Name + ".conf")
	tools.Remove("/www/server/vhost/ssl/" + website.Name + ".pem")
	tools.Remove("/www/server/vhost/ssl/" + website.Name + ".key")
	tools.Remove(website.Path)

	tools.Exec("systemctl reload openresty")

	return nil
}

// GetConfig 获取网站配置
func (r *WebsiteImpl) GetConfig(id uint) (WebsiteSetting, error) {
	var website models.Website
	if err := facades.Orm().Query().Where("id", id).First(&website); err != nil {
		return WebsiteSetting{}, err
	}

	config := tools.Read("/www/server/vhost/" + website.Name + ".conf")

	var setting WebsiteSetting
	setting.Name = website.Name
	setting.Path = website.Path
	setting.Ssl = website.Ssl
	setting.Php = strconv.Itoa(website.Php)
	setting.Raw = config

	ports := tools.Cut(config, "# port标记位开始", "# port标记位结束")
	matches := regexp.MustCompile(`listen\s+(.*);`).FindAllStringSubmatch(ports, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		setting.Ports = append(setting.Ports, match[1])
	}
	serverName := tools.Cut(config, "# server_name标记位开始", "# server_name标记位结束")
	match := regexp.MustCompile(`server_name\s+(.*);`).FindStringSubmatch(serverName)
	if len(match) > 1 {
		setting.Domains = strings.Split(match[1], " ")
	}
	root := tools.Cut(config, "# root标记位开始", "# root标记位结束")
	match = regexp.MustCompile(`root\s+(.*);`).FindStringSubmatch(root)
	if len(match) > 1 {
		setting.Root = match[1]
	}
	index := tools.Cut(config, "# index标记位开始", "# index标记位结束")
	match = regexp.MustCompile(`index\s+(.*);`).FindStringSubmatch(index)
	if len(match) > 1 {
		setting.Index = match[1]
	}

	if tools.Exists(setting.Root + "/.user.ini") {
		userIni := tools.Read(setting.Path + "/.user.ini")
		if strings.Contains(userIni, "open_basedir") {
			setting.OpenBasedir = true
		} else {
			setting.OpenBasedir = false
		}
	} else {
		setting.OpenBasedir = false
	}

	setting.SslCertificate = tools.Read("/www/server/vhost/ssl/" + website.Name + ".pem")
	setting.SslCertificateKey = tools.Read("/www/server/vhost/ssl/" + website.Name + ".key")
	if setting.Ssl {
		ssl := tools.Cut(config, "# ssl标记位开始", "# ssl标记位结束")
		setting.HttpRedirect = strings.Contains(ssl, "# http重定向标记位")
		setting.Hsts = strings.Contains(ssl, "# hsts标记位")

		certData := tools.Read("/www/server/vhost/ssl/" + website.Name + ".pem")
		block, _ := pem.Decode([]byte(certData))
		if block != nil {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				setting.SslNotBefore = cert.NotBefore.Format("2006-01-02 15:04:05")
				setting.SslNotAfter = cert.NotAfter.Format("2006-01-02 15:04:05")
				setting.SslIssuer = cert.Issuer.CommonName
				setting.SslOCSPServer = cert.OCSPServer
				setting.SSlDNSNames = cert.DNSNames
			}
		}
	} else {
		setting.HttpRedirect = false
		setting.Hsts = false
	}

	waf := tools.Cut(config, "# waf标记位开始", "# waf标记位结束")
	setting.Waf = strings.Contains(waf, "waf on;")
	match = regexp.MustCompile(`waf_mode\s+(.+);`).FindStringSubmatch(waf)
	if len(match) > 1 {
		setting.WafMode = match[1]
	}
	match = regexp.MustCompile(`waf_cc_deny\s+(.+);`).FindStringSubmatch(waf)
	if len(match) > 1 {
		setting.WafCcDeny = match[1]
	}
	match = regexp.MustCompile(`waf_cache\s+(.+);`).FindStringSubmatch(waf)
	if len(match) > 1 {
		setting.WafCache = match[1]
	}

	setting.Rewrite = tools.Read("/www/server/vhost/rewrite/" + website.Name + ".conf")
	setting.Log = tools.Escape(tools.Exec(`tail -n 100 '/www/wwwlogs/` + website.Name + `.log'`))

	return setting, nil
}

// GetConfigByName 根据网站名称获取网站配置
func (r *WebsiteImpl) GetConfigByName(name string) (WebsiteSetting, error) {
	var website models.Website
	if err := facades.Orm().Query().Where("name", name).First(&website); err != nil {
		return WebsiteSetting{}, err
	}

	return r.GetConfig(website.ID)
}
