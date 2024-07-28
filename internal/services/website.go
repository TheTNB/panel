// Package services 网站服务
package services

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"

	requests "github.com/TheTNB/panel/v2/app/http/requests/website"
	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/embed"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/pkg/cert"
	"github.com/TheTNB/panel/v2/pkg/db"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/str"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type WebsiteImpl struct {
	setting internal.Setting
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
func (r *WebsiteImpl) Add(website requests.Add) (models.Website, error) {
	w := models.Website{
		Name:   website.Name,
		Status: true,
		Path:   website.Path,
		PHP:    cast.ToInt(website.PHP),
		SSL:    false,
	}
	if err := facades.Orm().Query().Create(&w); err != nil {
		return models.Website{}, err
	}

	if err := io.Mkdir(website.Path, 0755); err != nil {
		return models.Website{}, err
	}

	index, err := embed.WebsiteFS.ReadFile(filepath.Join("website", "index.html"))
	if err != nil {
		return models.Website{}, fmt.Errorf("获取index模板文件失败: %w", err)
	}
	if err = io.Write(website.Path+"/index.html", string(index), 0644); err != nil {
		return models.Website{}, err
	}

	notFound, err := embed.WebsiteFS.ReadFile(filepath.Join("website", "404.html"))
	if err != nil {
		return models.Website{}, fmt.Errorf("获取404模板文件失败: %w", err)
	}
	if err = io.Write(website.Path+"/404.html", string(notFound), 0644); err != nil {
		return models.Website{}, err
	}

	portList := ""
	domainList := ""
	portUsed := make(map[uint]bool)
	domainUsed := make(map[string]bool)

	for i, port := range website.Ports {
		if _, ok := portUsed[port]; !ok {
			if i == len(website.Ports)-1 {
				portList += "    listen " + cast.ToString(port) + ";\n"
				portList += "    listen [::]:" + cast.ToString(port) + ";"
			} else {
				portList += "    listen " + cast.ToString(port) + ";\n"
				portList += "    listen [::]:" + cast.ToString(port) + ";\n"
			}
			portUsed[port] = true
		}
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
    include enable-php-%s.conf;
    # php标记位结束

    # waf标记位开始
    waf off;
    waf_rule_path /www/server/openresty/ngx_waf/assets/rules/;
    waf_mode DYNAMIC;
    waf_cc_deny rate=1000r/m duration=60m;
    waf_cache capacity=50;
    # waf标记位结束

    # 错误页配置，可自行设置
    error_page 404 /404.html;
    #error_page 502 /502.html;

    # acme证书签发配置，不可修改
    include /www/server/vhost/acme/%s.conf;

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
`, portList, domainList, website.Path, website.PHP, website.Name, website.Name, website.Name, website.Name)

	if err = io.Write("/www/server/vhost/"+website.Name+".conf", nginxConf, 0644); err != nil {
		return models.Website{}, err
	}
	if err = io.Write("/www/server/vhost/rewrite/"+website.Name+".conf", "", 0644); err != nil {
		return models.Website{}, err
	}
	if err = io.Write("/www/server/vhost/acme/"+website.Name+".conf", "", 0644); err != nil {
		return models.Website{}, err
	}
	if err = io.Write("/www/server/vhost/ssl/"+website.Name+".pem", "", 0644); err != nil {
		return models.Website{}, err
	}
	if err = io.Write("/www/server/vhost/ssl/"+website.Name+".key", "", 0644); err != nil {
		return models.Website{}, err
	}

	if err = io.Chmod(website.Path, 0755); err != nil {
		return models.Website{}, err
	}
	if err = io.Chown(website.Path, "www", "www"); err != nil {
		return models.Website{}, err
	}

	if err = systemctl.Reload("openresty"); err != nil {
		_, err = shell.Execf("openresty -t")
		return models.Website{}, err
	}

	rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
	if website.DB && website.DBType == "mysql" {
		mysql, err := db.NewMySQL("root", rootPassword, "/tmp/mysql.sock", "unix")
		if err != nil {
			return models.Website{}, err
		}
		if err = mysql.DatabaseCreate(website.DBName); err != nil {
			return models.Website{}, err
		}
		if err = mysql.UserCreate(website.DBUser, website.DBPassword); err != nil {
			return models.Website{}, err
		}
		if err = mysql.PrivilegesGrant(website.DBUser, website.DBName); err != nil {
			return models.Website{}, err
		}
	}
	if website.DB && website.DBType == "postgresql" {
		_, _ = shell.Execf(`echo "CREATE DATABASE '%s';" | su - postgres -c "psql"`, website.DBName)
		_, _ = shell.Execf(`echo "CREATE USER '%s' WITH PASSWORD '%s';" | su - postgres -c "psql"`, website.DBUser, website.DBPassword)
		_, _ = shell.Execf(`echo "ALTER DATABASE '%s' OWNER TO '%s';" | su - postgres -c "psql"`, website.DBName, website.DBUser)
		_, _ = shell.Execf(`echo "GRANT ALL PRIVILEGES ON DATABASE '%s' TO '%s';" | su - postgres -c "psql"`, website.DBName, website.DBUser)
		userConfig := "host    " + website.DBName + "    " + website.DBUser + "    127.0.0.1/32    scram-sha-256"
		_, _ = shell.Execf(`echo "` + userConfig + `" >> /www/server/postgresql/data/pg_hba.conf`)
		_ = systemctl.Reload("postgresql")
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
	raw, err := io.Read("/www/server/vhost/" + website.Name + ".conf")
	if err != nil {
		return err
	}
	if strings.TrimSpace(raw) != strings.TrimSpace(config.Raw) {
		if err = io.Write("/www/server/vhost/"+website.Name+".conf", config.Raw, 0644); err != nil {
			return err
		}
		if err = systemctl.Reload("openresty"); err != nil {
			_, err = shell.Execf("openresty -t")
			return err
		}

		return nil
	}

	// 目录
	path := config.Path
	if !io.Exists(path) {
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
	domainConfigOld := str.Cut(raw, "# server_name标记位开始", "# server_name标记位结束")
	if len(strings.TrimSpace(domainConfigOld)) == 0 {
		return errors.New("配置文件中缺少server_name标记位")
	}
	raw = strings.Replace(raw, domainConfigOld, "\n    "+domain+"\n    ", -1)

	// 端口
	var portConf strings.Builder
	ports := config.Ports
	for _, port := range ports {
		https := ""
		quic := false
		if slices.Contains(config.SSLPorts, port) {
			https = " ssl"
			if slices.Contains(config.QUICPorts, port) {
				quic = true
			}
		}

		portConf.WriteString(fmt.Sprintf("    listen %d%s;\n", port, https))
		portConf.WriteString(fmt.Sprintf("    listen [::]:%d%s;\n", port, https))
		if quic {
			portConf.WriteString(fmt.Sprintf("    listen %d%s;\n", port, " quic"))
			portConf.WriteString(fmt.Sprintf("    listen [::]:%d%s;\n", port, " quic"))
		}
	}
	portConf.WriteString("    ")
	portConfNew := portConf.String()
	portConfOld := str.Cut(raw, "# port标记位开始", "# port标记位结束")
	if len(strings.TrimSpace(portConfOld)) == 0 {
		return errors.New("配置文件中缺少port标记位")
	}
	raw = strings.Replace(raw, portConfOld, "\n"+portConfNew, -1)

	// 运行目录
	root := str.Cut(raw, "# root标记位开始", "# root标记位结束")
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
	index := str.Cut(raw, "# index标记位开始", "# index标记位结束")
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
		if err = io.Write(root+".user.ini", "open_basedir="+path+":/tmp/", 0644); err != nil {
			return err
		}
	} else {
		if io.Exists(root + ".user.ini") {
			if err = io.Remove(root + ".user.ini"); err != nil {
				return err
			}
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
	wafConfigOld := str.Cut(raw, "# waf标记位开始", "# waf标记位结束")
	if len(strings.TrimSpace(wafConfigOld)) != 0 {
		raw = strings.Replace(raw, wafConfigOld, "", -1)
	}
	raw = strings.Replace(raw, "# waf标记位开始", wafConfig, -1)

	// SSL
	ssl := config.SSL
	website.SSL = ssl
	if ssl {
		if _, err = cert.ParseCert(config.SSLCertificate); err != nil {
			return errors.New("TLS证书格式错误")
		}
		if _, err = cert.ParseKey(config.SSLCertificateKey); err != nil {
			return errors.New("TLS私钥格式错误")
		}
	}
	if err = io.Write("/www/server/vhost/ssl/"+website.Name+".pem", config.SSLCertificate, 0644); err != nil {
		return err
	}
	if err = io.Write("/www/server/vhost/ssl/"+website.Name+".key", config.SSLCertificateKey, 0644); err != nil {
		return err
	}
	if ssl {
		sslConfig := `# ssl标记位开始
    ssl_certificate /www/server/vhost/ssl/` + website.Name + `.pem;
    ssl_certificate_key /www/server/vhost/ssl/` + website.Name + `.key;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:10m;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_early_data on;
    `
		if config.HTTPRedirect {
			sslConfig += `# http重定向标记位开始
    if ($server_port !~ 443){
        return 301 https://$host$request_uri;
    }
    error_page 497  https://$host$request_uri;
    # http重定向标记位结束
    `
		}
		if config.HSTS {
			sslConfig += `# hsts标记位开始
    add_header Strict-Transport-Security "max-age=63072000" always;
    # hsts标记位结束
    `
		}
		if config.OCSP {
			sslConfig += `# ocsp标记位开始
    ssl_stapling on;
    ssl_stapling_verify on;
    # ocsp标记位结束
    `
		}
		sslConfigOld := str.Cut(raw, "# ssl标记位开始", "# ssl标记位结束")
		if len(strings.TrimSpace(sslConfigOld)) != 0 {
			raw = strings.Replace(raw, sslConfigOld, "", -1)
		}
		raw = strings.Replace(raw, "# ssl标记位开始", sslConfig, -1)
	} else {
		sslConfigOld := str.Cut(raw, "# ssl标记位开始", "# ssl标记位结束")
		if len(strings.TrimSpace(sslConfigOld)) != 0 {
			raw = strings.Replace(raw, sslConfigOld, "\n    ", -1)
		}
	}

	if website.PHP != config.PHP {
		website.PHP = config.PHP
		phpConfigOld := str.Cut(raw, "# php标记位开始", "# php标记位结束")
		phpConfig := `
    include enable-php-` + strconv.Itoa(website.PHP) + `.conf;
    `
		if len(strings.TrimSpace(phpConfigOld)) != 0 {
			raw = strings.Replace(raw, phpConfigOld, phpConfig, -1)
		}
	}

	if err = facades.Orm().Query().Save(&website); err != nil {
		return err
	}

	if err = io.Write("/www/server/vhost/"+website.Name+".conf", raw, 0644); err != nil {
		return err
	}
	if err = io.Write("/www/server/vhost/rewrite/"+website.Name+".conf", config.Rewrite, 0644); err != nil {
		return err
	}

	err = systemctl.Reload("openresty")
	if err != nil {
		_, err = shell.Execf("openresty -t")
	}

	return err
}

// Delete 删除网站
func (r *WebsiteImpl) Delete(request requests.Delete) error {
	var website models.Website
	if err := facades.Orm().Query().With("Cert").Where("id", request.ID).FirstOrFail(&website); err != nil {
		return err
	}

	if website.Cert != nil {
		return errors.New("网站" + website.Name + "已绑定SSL证书，请先删除证书")
	}

	if _, err := facades.Orm().Query().Delete(&website); err != nil {
		return err
	}

	_ = io.Remove("/www/server/vhost/" + website.Name + ".conf")
	_ = io.Remove("/www/server/vhost/rewrite/" + website.Name + ".conf")
	_ = io.Remove("/www/server/vhost/acme/" + website.Name + ".conf")
	_ = io.Remove("/www/server/vhost/ssl/" + website.Name + ".pem")
	_ = io.Remove("/www/server/vhost/ssl/" + website.Name + ".key")

	if request.Path {
		_ = io.Remove(website.Path)
	}
	if request.DB {
		rootPassword := r.setting.Get(models.SettingKeyMysqlRootPassword)
		mysql, err := db.NewMySQL("root", rootPassword, "/tmp/mysql.sock", "unix")
		if err != nil {
			return err
		}
		_ = mysql.DatabaseDrop(website.Name)
		_ = mysql.UserDrop(website.Name)
		_, _ = shell.Execf(`echo "DROP DATABASE IF EXISTS '%s';" | su - postgres -c "psql"`, website.Name)
		_, _ = shell.Execf(`echo "DROP USER IF EXISTS '%s';" | su - postgres -c "psql"`, website.Name)
	}

	err := systemctl.Reload("openresty")
	if err != nil {
		_, err = shell.Execf("openresty -t")
	}

	return err
}

// GetConfig 获取网站配置
func (r *WebsiteImpl) GetConfig(id uint) (types.WebsiteSetting, error) {
	var website models.Website
	if err := facades.Orm().Query().Where("id", id).First(&website); err != nil {
		return types.WebsiteSetting{}, err
	}

	config, err := io.Read("/www/server/vhost/" + website.Name + ".conf")
	if err != nil {
		return types.WebsiteSetting{}, err
	}

	var setting types.WebsiteSetting
	setting.Name = website.Name
	setting.Path = website.Path
	setting.SSL = website.SSL
	setting.PHP = strconv.Itoa(website.PHP)
	setting.Raw = config

	portStr := str.Cut(config, "# port标记位开始", "# port标记位结束")
	matches := regexp.MustCompile(`listen\s+([^;]*);?`).FindAllStringSubmatch(portStr, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		// 跳过 ipv6
		if strings.Contains(match[1], "[::]") {
			continue
		}

		// 处理 443 ssl 之类的情况
		ports := strings.Fields(match[1])
		if len(ports) == 0 {
			continue
		}
		if !slices.Contains(setting.Ports, ports[0]) {
			setting.Ports = append(setting.Ports, ports[0])
		}
		if len(ports) > 1 && ports[1] == "ssl" {
			setting.SSLPorts = append(setting.SSLPorts, ports[0])
		} else if len(ports) > 1 && ports[1] == "quic" {
			setting.QUICPorts = append(setting.QUICPorts, ports[0])
		}
	}
	serverName := str.Cut(config, "# server_name标记位开始", "# server_name标记位结束")
	match := regexp.MustCompile(`server_name\s+([^;]*);?`).FindStringSubmatch(serverName)
	if len(match) > 1 {
		setting.Domains = strings.Split(match[1], " ")
	}
	root := str.Cut(config, "# root标记位开始", "# root标记位结束")
	match = regexp.MustCompile(`root\s+([^;]*);?`).FindStringSubmatch(root)
	if len(match) > 1 {
		setting.Root = match[1]
	}
	index := str.Cut(config, "# index标记位开始", "# index标记位结束")
	match = regexp.MustCompile(`index\s+([^;]*);?`).FindStringSubmatch(index)
	if len(match) > 1 {
		setting.Index = match[1]
	}

	if io.Exists(filepath.Join(setting.Root, ".user.ini")) {
		userIni, _ := io.Read(filepath.Join(setting.Root, ".user.ini"))
		if strings.Contains(userIni, "open_basedir") {
			setting.OpenBasedir = true
		}
	}

	crt, _ := io.Read("/www/server/vhost/ssl/" + website.Name + ".pem")
	setting.SSLCertificate = crt
	key, _ := io.Read("/www/server/vhost/ssl/" + website.Name + ".key")
	setting.SSLCertificateKey = key
	if setting.SSL {
		ssl := str.Cut(config, "# ssl标记位开始", "# ssl标记位结束")
		setting.HTTPRedirect = strings.Contains(ssl, "# http重定向标记位")
		setting.HSTS = strings.Contains(ssl, "# hsts标记位")
		setting.OCSP = strings.Contains(ssl, "# ocsp标记位")
	}

	// 解析证书信息
	if decode, err := cert.ParseCert(crt); err == nil {
		setting.SSLNotBefore = decode.NotBefore.Format("2006-01-02 15:04:05")
		setting.SSLNotAfter = decode.NotAfter.Format("2006-01-02 15:04:05")
		setting.SSLIssuer = decode.Issuer.CommonName
		setting.SSLOCSPServer = decode.OCSPServer
		setting.SSLDNSNames = decode.DNSNames
	}

	waf := str.Cut(config, "# waf标记位开始", "# waf标记位结束")
	setting.Waf = strings.Contains(waf, "waf on;")
	match = regexp.MustCompile(`waf_mode\s+([^;]*);?`).FindStringSubmatch(waf)
	if len(match) > 1 {
		setting.WafMode = match[1]
	}
	match = regexp.MustCompile(`waf_cc_deny\s+([^;]*);?`).FindStringSubmatch(waf)
	if len(match) > 1 {
		setting.WafCcDeny = match[1]
	}
	match = regexp.MustCompile(`waf_cache\s+([^;]*);?`).FindStringSubmatch(waf)
	if len(match) > 1 {
		setting.WafCache = match[1]
	}

	rewrite, _ := io.Read("/www/server/vhost/rewrite/" + website.Name + ".conf")
	setting.Rewrite = rewrite
	log, _ := shell.Execf(`tail -n 100 '/www/wwwlogs/%s.log'`, website.Name)
	setting.Log = log

	return setting, err
}

// GetConfigByName 根据网站名称获取网站配置
func (r *WebsiteImpl) GetConfigByName(name string) (types.WebsiteSetting, error) {
	var website models.Website
	if err := facades.Orm().Query().Where("name", name).First(&website); err != nil {
		return types.WebsiteSetting{}, err
	}

	return r.GetConfig(website.ID)
}

// GetIDByName 根据网站名称获取网站ID
func (r *WebsiteImpl) GetIDByName(name string) (uint, error) {
	var website models.Website
	if err := facades.Orm().Query().Where("name", name).First(&website); err != nil {
		return 0, err
	}

	return website.ID, nil
}
