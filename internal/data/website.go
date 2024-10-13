package data

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/embed"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/cert"
	"github.com/TheTNB/panel/pkg/db"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/str"
	"github.com/TheTNB/panel/pkg/systemctl"
	"github.com/TheTNB/panel/pkg/types"
)

type websiteRepo struct {
	settingRepo biz.SettingRepo
}

func NewWebsiteRepo() biz.WebsiteRepo {
	return &websiteRepo{
		settingRepo: NewSettingRepo(),
	}
}

func (r *websiteRepo) UpdateDefaultConfig(req *request.WebsiteDefaultConfig) error {
	if err := io.Write(filepath.Join(app.Root, "server/nginx/html/index.html"), req.Index, 0644); err != nil {
		return err
	}
	if err := io.Write(filepath.Join(app.Root, "server/nginx/html/stop.html"), req.Stop, 0644); err != nil {
		return err
	}

	return systemctl.Reload("nginx")
}

func (r *websiteRepo) Count() (int64, error) {
	var count int64
	if err := app.Orm.Model(&biz.Website{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *websiteRepo) Get(id uint) (*types.WebsiteSetting, error) {
	website := new(biz.Website)
	if err := app.Orm.Where("id", id).First(website).Error; err != nil {
		return nil, err
	}

	config, err := io.Read(filepath.Join(app.Root, "server/vhost", website.Name+".conf"))
	if err != nil {
		return nil, err
	}

	setting := new(types.WebsiteSetting)
	setting.ID = website.ID
	setting.Name = website.Name
	setting.Path = website.Path
	setting.SSL = website.SSL
	setting.PHP = website.PHP
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
		if !slices.Contains(setting.Ports, cast.ToUint(ports[0])) {
			setting.Ports = append(setting.Ports, cast.ToUint(ports[0]))
		}
		if len(ports) > 1 && ports[1] == "ssl" {
			setting.SSLPorts = append(setting.SSLPorts, cast.ToUint(ports[0]))
		} else if len(ports) > 1 && ports[1] == "quic" {
			setting.QUICPorts = append(setting.QUICPorts, cast.ToUint(ports[0]))
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

	crt, _ := io.Read(filepath.Join(app.Root, "server/vhost/ssl", website.Name+".pem"))
	setting.SSLCertificate = crt
	key, _ := io.Read(filepath.Join(app.Root, "server/vhost/ssl", website.Name+".key"))
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

	rewrite, _ := io.Read(filepath.Join(app.Root, "server/vhost/rewrite", website.Name+".conf"))
	setting.Rewrite = rewrite
	log, _ := shell.Execf(`tail -n 100 '%s/wwwlogs/%s.log'`, app.Root, website.Name)
	setting.Log = log

	return setting, err
}

func (r *websiteRepo) GetByName(name string) (*types.WebsiteSetting, error) {
	website := new(biz.Website)
	if err := app.Orm.Where("name", name).First(website).Error; err != nil {
		return nil, err
	}

	return r.Get(website.ID)

}

func (r *websiteRepo) List(page, limit uint) ([]*biz.Website, int64, error) {
	var websites []*biz.Website
	var total int64

	if err := app.Orm.Model(&biz.Website{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := app.Orm.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&websites).Error; err != nil {
		return nil, 0, err
	}

	return websites, total, nil
}

func (r *websiteRepo) Create(req *request.WebsiteCreate) (*biz.Website, error) {
	w := &biz.Website{
		Name:   req.Name,
		Status: true,
		Path:   req.Path,
		PHP:    req.PHP,
		SSL:    false,
	}
	if err := app.Orm.Create(w).Error; err != nil {
		return nil, err
	}

	if err := io.Mkdir(req.Path, 0755); err != nil {
		return nil, err
	}

	index, err := embed.WebsiteFS.ReadFile(filepath.Join("website", "index.html"))
	if err != nil {
		return nil, fmt.Errorf("获取index模板文件失败: %w", err)
	}
	if err = io.Write(filepath.Join(req.Path, "index.html"), string(index), 0644); err != nil {
		return nil, err
	}

	notFound, err := embed.WebsiteFS.ReadFile(filepath.Join("website", "404.html"))
	if err != nil {
		return nil, fmt.Errorf("获取404模板文件失败: %w", err)
	}
	if err = io.Write(req.Path+"/404.html", string(notFound), 0644); err != nil {
		return nil, err
	}

	portList := ""
	domainList := ""
	portUsed := make(map[uint]bool)
	domainUsed := make(map[string]bool)

	for i, port := range req.Ports {
		if _, ok := portUsed[port]; !ok {
			if i == len(req.Ports)-1 {
				portList += "    listen " + cast.ToString(port) + ";\n"
				portList += "    listen [::]:" + cast.ToString(port) + ";"
			} else {
				portList += "    listen " + cast.ToString(port) + ";\n"
				portList += "    listen [::]:" + cast.ToString(port) + ";\n"
			}
			portUsed[port] = true
		}
	}
	for _, domain := range req.Domains {
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

    # 错误页配置，可自行设置
    error_page 404 /404.html;
    #error_page 502 /502.html;

    # acme证书签发配置，不可修改
    include %s/server/vhost/acme/%s.conf;

    # 伪静态规则引入，修改后将导致面板设置的伪静态规则失效
    include %s/server/vhost/rewrite/%s.conf;

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

    access_log %s/wwwlogs/%s.log;
    error_log %s/wwwlogs/%s.log;
}
`, portList, domainList, req.Path, req.PHP, app.Root, req.Name, app.Root, req.Name, app.Root, req.Name, app.Root, req.Name)

	if err = io.Write(filepath.Join(app.Root, "server/vhost", req.Name+".conf"), nginxConf, 0644); err != nil {
		return nil, err
	}
	if err = io.Write(filepath.Join(app.Root, "server/vhost/rewrite", req.Name+".conf"), "", 0644); err != nil {
		return nil, err
	}
	if err = io.Write(filepath.Join(app.Root, "server/vhost/acme", req.Name+".conf"), "", 0644); err != nil {
		return nil, err
	}
	if err = io.Write(filepath.Join(app.Root, "server/vhost/ssl", req.Name+".pem"), "", 0644); err != nil {
		return nil, err
	}
	if err = io.Write(filepath.Join(app.Root, "server/vhost/ssl", req.Name+".key"), "", 0644); err != nil {
		return nil, err
	}

	if err = io.Chmod(req.Path, 0755); err != nil {
		return nil, err
	}
	if err = io.Chown(req.Path, "www", "www"); err != nil {
		return nil, err
	}

	if err = systemctl.Reload("nginx"); err != nil {
		_, err = shell.Execf("nginx -t")
		return nil, err
	}

	rootPassword, err := r.settingRepo.Get(biz.SettingKeyMySQLRootPassword)
	if err == nil && req.DB && req.DBType == "mysql" {
		mysql, err := db.NewMySQL("root", rootPassword, "/tmp/mysql.sock", "unix")
		if err != nil {
			return nil, err
		}
		if err = mysql.DatabaseCreate(req.DBName); err != nil {
			return nil, err
		}
		if err = mysql.UserCreate(req.DBUser, req.DBPassword); err != nil {
			return nil, err
		}
		if err = mysql.PrivilegesGrant(req.DBUser, req.DBName); err != nil {
			return nil, err
		}
	}
	if req.DB && req.DBType == "postgresql" {
		_, _ = shell.Execf(`echo "CREATE DATABASE '%s';" | su - postgres -c "psql"`, req.DBName)
		_, _ = shell.Execf(`echo "CREATE USER '%s' WITH PASSWORD '%s';" | su - postgres -c "psql"`, req.DBUser, req.DBPassword)
		_, _ = shell.Execf(`echo "ALTER DATABASE '%s' OWNER TO '%s';" | su - postgres -c "psql"`, req.DBName, req.DBUser)
		_, _ = shell.Execf(`echo "GRANT ALL PRIVILEGES ON DATABASE '%s' TO '%s';" | su - postgres -c "psql"`, req.DBName, req.DBUser)
		userConfig := "host    " + req.DBName + "    " + req.DBUser + "    127.0.0.1/32    scram-sha-256"
		_, _ = shell.Execf(`echo "`+userConfig+`" >> %s/server/postgresql/data/pg_hba.conf`, app.Root)
		_ = systemctl.Reload("postgresql")
	}

	return w, nil
}

func (r *websiteRepo) Update(req *request.WebsiteUpdate) error {
	website := new(biz.Website)
	if err := app.Orm.Where("id", req.ID).First(website).Error; err != nil {
		return err
	}

	if !website.Status {
		return errors.New("网站已停用，请先启用")
	}

	// 原文
	raw, err := io.Read(filepath.Join(app.Root, "server/vhost", website.Name+".conf"))
	if err != nil {
		return err
	}
	if strings.TrimSpace(raw) != strings.TrimSpace(req.Raw) {
		if err = io.Write(filepath.Join(app.Root, "server/vhost", website.Name+".conf"), req.Raw, 0644); err != nil {
			return err
		}
		if err = systemctl.Reload("nginx"); err != nil {
			_, err = shell.Execf("nginx -t")
			return err
		}

		return nil
	}

	// 目录
	path := req.Path
	if !io.Exists(path) {
		return errors.New("网站目录不存在")
	}
	website.Path = path

	// 域名
	domain := "server_name"
	domains := req.Domains
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
	ports := req.Ports
	for _, port := range ports {
		https := ""
		quic := false
		if slices.Contains(req.SSLPorts, port) {
			https = " ssl"
			if slices.Contains(req.QUICPorts, port) {
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
	rootNew := strings.Replace(root, match[1], req.Root, -1)
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
	indexNew := strings.Replace(index, match[1], req.Index, -1)
	raw = strings.Replace(raw, index, indexNew, -1)

	// 防跨站
	if !strings.HasSuffix(req.Root, "/") {
		req.Root += "/"
	}
	if req.OpenBasedir {
		if err = io.Write(req.Root+".user.ini", "open_basedir="+path+":/tmp/", 0644); err != nil {
			return err
		}
	} else {
		if io.Exists(req.Root + ".user.ini") {
			if err = io.Remove(req.Root + ".user.ini"); err != nil {
				return err
			}
		}
	}

	// SSL
	if err = io.Write(filepath.Join(app.Root, "server/vhost/ssl", website.Name+".pem"), req.SSLCertificate, 0644); err != nil {
		return err
	}
	if err = io.Write(filepath.Join(app.Root, "server/vhost/ssl", website.Name+".key"), req.SSLCertificateKey, 0644); err != nil {
		return err
	}
	website.SSL = req.SSL
	if req.SSL {
		if _, err = cert.ParseCert(req.SSLCertificate); err != nil {
			return errors.New("TLS证书格式错误")
		}
		if _, err = cert.ParseKey(req.SSLCertificateKey); err != nil {
			return errors.New("TLS私钥格式错误")
		}
		sslConfig := fmt.Sprintf(`# ssl标记位开始
    ssl_certificate %s/server/vhost/ssl/%s.pem;
    ssl_certificate_key %s/server/vhost/ssl/%s.key;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:10m;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_early_data on;
    `, app.Root, website.Name, app.Root, website.Name)
		if req.HTTPRedirect {
			sslConfig += `# http重定向标记位开始
    if ($server_port !~ 443){
        return 301 https://$host$request_uri;
    }
    error_page 497  https://$host$request_uri;
    # http重定向标记位结束
    `
		}
		if req.HSTS {
			sslConfig += `# hsts标记位开始
    add_header Strict-Transport-Security "max-age=63072000" always;
    # hsts标记位结束
    `
		}
		if req.OCSP {
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

	if website.PHP != req.PHP {
		website.PHP = req.PHP
		phpConfigOld := str.Cut(raw, "# php标记位开始", "# php标记位结束")
		phpConfig := `
    include enable-php-` + strconv.Itoa(website.PHP) + `.conf;
    `
		if len(strings.TrimSpace(phpConfigOld)) != 0 {
			raw = strings.Replace(raw, phpConfigOld, phpConfig, -1)
		}
	}

	if err = app.Orm.Save(website).Error; err != nil {
		return err
	}

	if err = io.Write(filepath.Join(app.Root, "server/vhost", website.Name+".conf"), raw, 0644); err != nil {
		return err
	}
	if err = io.Write(filepath.Join(app.Root, "server/vhost/rewrite", website.Name+".conf"), req.Rewrite, 0644); err != nil {
		return err
	}

	err = systemctl.Reload("nginx")
	if err != nil {
		_, err = shell.Execf("nginx -t")
	}

	return err
}

func (r *websiteRepo) Delete(req *request.WebsiteDelete) error {
	website := new(biz.Website)
	if err := app.Orm.Preload("Cert").Where("id", req.ID).First(website).Error; err != nil {
		return err
	}

	if website.Cert != nil {
		return errors.New("网站" + website.Name + "已绑定SSL证书，请先删除证书")
	}

	if err := app.Orm.Delete(website).Error; err != nil {
		return err
	}

	_ = io.Remove(filepath.Join(app.Root, "server/vhost", website.Name+".conf"))
	_ = io.Remove(filepath.Join(app.Root, "server/vhost/rewrite", website.Name+".conf"))
	_ = io.Remove(filepath.Join(app.Root, "server/vhost/acme", website.Name+".conf"))
	_ = io.Remove(filepath.Join(app.Root, "server/vhost/ssl", website.Name+".pem"))
	_ = io.Remove(filepath.Join(app.Root, "server/vhost/ssl", website.Name+".key"))

	if req.Path {
		_ = io.Remove(website.Path)
	}
	if req.DB {
		rootPassword, err := r.settingRepo.Get(biz.SettingKeyMySQLRootPassword)
		if err != nil {
			return err
		}
		mysql, err := db.NewMySQL("root", rootPassword, "/tmp/mysql.sock", "unix")
		if err != nil {
			return err
		}
		_ = mysql.DatabaseDrop(website.Name)
		_ = mysql.UserDrop(website.Name)
		_, _ = shell.Execf(`echo "DROP DATABASE IF EXISTS '%s';" | su - postgres -c "psql"`, website.Name)
		_, _ = shell.Execf(`echo "DROP USER IF EXISTS '%s';" | su - postgres -c "psql"`, website.Name)
	}

	err := systemctl.Reload("nginx")
	if err != nil {
		_, err = shell.Execf("nginx -t")
	}

	return err
}

func (r *websiteRepo) ClearLog(id uint) error {
	website := new(biz.Website)
	if err := app.Orm.Where("id", id).First(website).Error; err != nil {
		return err
	}

	_, err := shell.Execf(`echo "" > %s/wwwlogs/%s.log`, app.Root, website.Name)
	return err
}

func (r *websiteRepo) UpdateRemark(id uint, remark string) error {
	website := new(biz.Website)
	if err := app.Orm.Where("id", id).First(website).Error; err != nil {
		return err
	}

	website.Remark = remark
	return app.Orm.Save(website).Error
}

func (r *websiteRepo) ResetConfig(id uint) error {
	website := new(biz.Website)
	if err := app.Orm.Where("id", id).First(&website).Error; err != nil {
		return err
	}

	website.Status = true
	website.SSL = false
	if err := app.Orm.Save(website).Error; err != nil {
		return err
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

    # 错误页配置，可自行设置
    error_page 404 /404.html;
    #error_page 502 /502.html;

    # acme证书签发配置，不可修改
    include %s/server/vhost/acme/%s.conf;

    # 伪静态规则引入，修改后将导致面板设置的伪静态规则失效
    include %s/server/vhost/rewrite/%s.conf;

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

    access_log %s/wwwlogs/%s.log;
    error_log %s/wwwlogs/%s.log;
}

`, website.Path, website.PHP, app.Root, website.Name, app.Root, website.Name, app.Root, website.Name, app.Root, website.Name)
	if err := io.Write(filepath.Join(app.Root, "server/vhost", website.Name+".conf"), raw, 0644); err != nil {
		return nil
	}
	if err := io.Write(filepath.Join(app.Root, "server/vhost/rewrite", website.Name+".conf"), "", 0644); err != nil {
		return nil
	}
	if err := io.Write(filepath.Join(app.Root, "server/vhost/acme", website.Name+".conf"), "", 0644); err != nil {
		return nil
	}
	if err := systemctl.Reload("nginx"); err != nil {
		_, err = shell.Execf("nginx -t")
		return err
	}

	return nil
}

func (r *websiteRepo) UpdateStatus(id uint, status bool) error {
	website := new(biz.Website)
	if err := app.Orm.Where("id", id).First(&website).Error; err != nil {
		return err
	}

	website.Status = status
	if err := app.Orm.Save(website).Error; err != nil {
		return err
	}

	raw, err := io.Read(filepath.Join(app.Root, "server/vhost", website.Name+".conf"))
	if err != nil {
		return err
	}

	// 运行目录
	rootConfig := str.Cut(raw, "# root标记位开始\n", "# root标记位结束")
	match := regexp.MustCompile(`root\s+(.+);`).FindStringSubmatch(rootConfig)
	if len(match) == 2 {
		if website.Status {
			root := regexp.MustCompile(`# root\s+(.+);`).FindStringSubmatch(rootConfig)
			raw = strings.ReplaceAll(raw, rootConfig, fmt.Sprintf("    root %s;\n    ", root[1]))
		} else {
			raw = strings.ReplaceAll(raw, rootConfig, fmt.Sprintf("    root %s/server/nginx/html;\n    # root %s;\n    ", app.Root, match[1]))
		}
	}

	// 默认文件
	indexConfig := str.Cut(raw, "# index标记位开始\n", "# index标记位结束")
	match = regexp.MustCompile(`index\s+(.+);`).FindStringSubmatch(indexConfig)
	if len(match) == 2 {
		if website.Status {
			index := regexp.MustCompile(`# index\s+(.+);`).FindStringSubmatch(indexConfig)
			raw = strings.ReplaceAll(raw, indexConfig, "    index "+index[1]+";\n    ")
		} else {
			raw = strings.ReplaceAll(raw, indexConfig, "    index stop.html;\n    # index "+match[1]+";\n    ")
		}
	}

	if err = io.Write(filepath.Join(app.Root, "server/vhost", website.Name+".conf"), raw, 0644); err != nil {
		return err
	}
	if err = systemctl.Reload("nginx"); err != nil {
		_, err = shell.Execf("nginx -t")
		return err
	}

	return nil
}
