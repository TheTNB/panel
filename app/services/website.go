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

	"panel/app/models"
	"panel/pkg/tools"
)

type Website interface {
	List(page int, limit int) (int64, []models.Website, error)
	Add(website PanelWebsite) (models.Website, error)
	Delete(id int) error
	GetConfig(id int) (WebsiteSetting, error)
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
		return w, err
	}

	tools.Mkdir(website.Path, 0755)

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
	tools.Write(website.Path+"/index.html", index, 0644)

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

	tools.Write("/www/server/vhost/"+website.Name+".conf", nginxConf, 0644)
	tools.Write("/www/server/vhost/rewrite/"+website.Name+".conf", "", 0644)
	tools.Write("/www/server/vhost/ssl/"+website.Name+".pem", "", 0644)
	tools.Write("/www/server/vhost/ssl/"+website.Name+".key", "", 0644)

	tools.Chmod(r.setting.Get(models.SettingKeyWebsitePath), 0755)
	tools.Chmod(website.Path, 0755)
	tools.Chown(r.setting.Get(models.SettingKeyWebsitePath), "www", "www")
	tools.Chown(website.Path, "www", "www")

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

// Delete 删除网站
func (r *WebsiteImpl) Delete(id int) error {
	var website models.Website
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&website); err != nil {
		return err
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
func (r *WebsiteImpl) GetConfig(id int) (WebsiteSetting, error) {
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

	return r.GetConfig(int(website.ID))
}
