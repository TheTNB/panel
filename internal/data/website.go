package data

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/samber/lo"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/embed"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/cert"
	"github.com/TheTNB/panel/pkg/db"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/nginx"
	"github.com/TheTNB/panel/pkg/punycode"
	"github.com/TheTNB/panel/pkg/shell"
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
	// 解析nginx配置
	config, err := io.Read(filepath.Join(app.Root, "server/vhost", website.Name+".conf"))
	if err != nil {
		return nil, err
	}
	p, err := nginx.NewParser(config)
	if err != nil {
		return nil, err
	}

	setting := new(types.WebsiteSetting)
	setting.ID = website.ID
	setting.Name = website.Name
	setting.Path = website.Path
	setting.HTTPS = website.Https
	setting.PHP = p.GetPHP()
	setting.Raw = config
	// 监听地址
	listens, err := p.GetListen()
	if err != nil {
		return nil, err
	}
	setting.Listens = lo.Map(
		lo.UniqBy(listens, func(listen []string) string {
			if len(listen) == 0 {
				return ""
			}
			return listen[0]
		}),
		func(listen []string, _ int) types.WebsiteListen {
			addr := listen[0]
			grouped := lo.GroupBy(listens, func(listen []string) string {
				if len(listen) == 0 {
					return ""
				}
				return listen[0]
			})[addr]
			return types.WebsiteListen{
				Address: addr,
				HTTPS:   lo.SomeBy(grouped, func(listen []string) bool { return lo.Contains(listen, "ssl") }),
				QUIC:    lo.SomeBy(grouped, func(listen []string) bool { return lo.Contains(listen, "quic") }),
			}
		},
	)
	// 域名
	domains, err := p.GetServerName()
	if err != nil {
		return nil, err
	}
	domains, err = punycode.DecodeDomains(domains)
	if err != nil {
		return nil, err
	}
	setting.Domains = domains
	// 运行目录
	root, _ := p.GetRoot()
	setting.Root = root
	// 默认文档
	index, _ := p.GetIndex()
	setting.Index = index
	// 防跨站
	if io.Exists(filepath.Join(setting.Root, ".user.ini")) {
		userIni, _ := io.Read(filepath.Join(setting.Root, ".user.ini"))
		if strings.Contains(userIni, "open_basedir") {
			setting.OpenBasedir = true
		}
	}
	// HTTPS
	if setting.HTTPS {
		setting.HTTPRedirect = p.GetHTTPSRedirect()
		setting.HSTS = p.GetHSTS()
		setting.OCSP = p.GetOCSP()
	}
	// 证书
	crt, _ := io.Read(filepath.Join(app.Root, "server/vhost/ssl", website.Name+".pem"))
	setting.SSLCertificate = crt
	key, _ := io.Read(filepath.Join(app.Root, "server/vhost/ssl", website.Name+".key"))
	setting.SSLCertificateKey = key
	// 解析证书信息
	if decode, err := cert.ParseCert(crt); err == nil {
		setting.SSLNotBefore = decode.NotBefore.Format("2006-01-02 15:04:05")
		setting.SSLNotAfter = decode.NotAfter.Format("2006-01-02 15:04:05")
		setting.SSLIssuer = decode.Issuer.CommonName
		setting.SSLOCSPServer = decode.OCSPServer
		setting.SSLDNSNames = decode.DNSNames
	}
	// 伪静态
	rewrite, _ := io.Read(filepath.Join(app.Root, "server/vhost/rewrite", website.Name+".conf"))
	setting.Rewrite = rewrite
	// 访问日志
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
	// 初始化nginx配置
	p, err := nginx.NewParser()
	if err != nil {
		return nil, err
	}
	// 监听地址
	var listens [][]string
	for _, listen := range req.Listens {
		listens = append(listens, []string{listen})
	}
	if err = p.SetListen(listens); err != nil {
		return nil, err
	}
	// 域名
	domains, err := punycode.EncodeDomains(req.Domains)
	if err != nil {
		return nil, err
	}
	if err = p.SetServerName(domains); err != nil {
		return nil, err
	}
	// 运行目录
	if err = p.SetRoot(req.Path); err != nil {
		return nil, err
	}
	// PHP
	if err = p.SetPHP(req.PHP); err != nil {
		return nil, err
	}
	// 伪静态
	includes, comments, err := p.GetIncludes()
	if err != nil {
		return nil, err
	}
	includes = append(includes, filepath.Join(app.Root, "server/vhost/rewrite", req.Name+".conf"))
	comments = append(comments, []string{"# 伪静态规则"})
	if err = p.SetIncludes(includes, comments); err != nil {
		return nil, err
	}
	// 日志
	if err = p.SetAccessLog(filepath.Join(app.Root, "wwwlogs", req.Name+".log")); err != nil {
		return nil, err
	}
	if err = p.SetErrorLog(filepath.Join(app.Root, "wwwlogs", req.Name+".error.log")); err != nil {
		return nil, err
	}

	// 初始化网站目录
	if err = io.Mkdir(req.Path, 0755); err != nil {
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
	if err = io.Write(filepath.Join(req.Path, "404.html"), string(notFound), 0644); err != nil {
		return nil, err
	}

	// 写nginx配置
	if err = io.Write(filepath.Join(app.Root, "server/vhost", req.Name+".conf"), p.Dump(), 0644); err != nil {
		return nil, err
	}
	if err = io.Write(filepath.Join(app.Root, "server/vhost/rewrite", req.Name+".conf"), "", 0644); err != nil {
		return nil, err
	}
	if err = io.Write(filepath.Join(app.Root, "server/vhost/ssl", req.Name+".pem"), "", 0644); err != nil {
		return nil, err
	}
	if err = io.Write(filepath.Join(app.Root, "server/vhost/ssl", req.Name+".key"), "", 0644); err != nil {
		return nil, err
	}

	// 设置目录权限
	if err = io.Chmod(req.Path, 0755); err != nil {
		return nil, err
	}
	if err = io.Chown(req.Path, "www", "www"); err != nil {
		return nil, err
	}

	// 创建面板网站
	w := &biz.Website{
		Name:   req.Name,
		Status: true,
		Path:   req.Path,
		Https:  false,
	}
	if err = app.Orm.Create(w).Error; err != nil {
		return nil, err
	}

	if err = systemctl.Reload("nginx"); err != nil {
		_, err = shell.Execf("nginx -t")
		return nil, err
	}

	// 创建数据库
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

	// 解析nginx配置
	config, err := io.Read(filepath.Join(app.Root, "server/vhost", website.Name+".conf"))
	if err != nil {
		return err
	}
	// 如果修改了原文，直接写入返回
	if strings.TrimSpace(config) != strings.TrimSpace(req.Raw) {
		if err = io.Write(filepath.Join(app.Root, "server/vhost", website.Name+".conf"), req.Raw, 0644); err != nil {
			return err
		}
		if err = systemctl.Reload("nginx"); err != nil {
			_, err = shell.Execf("nginx -t")
			return err
		}
		return nil
	}

	// 初始化nginx配置
	p, err := nginx.NewParser(config)
	if err != nil {
		return err
	}
	// 监听地址
	var listens [][]string
	for _, listen := range req.Listens {
		if !listen.HTTPS && !listen.QUIC {
			listens = append(listens, []string{listen.Address})
		}
		if listen.HTTPS {
			listens = append(listens, []string{listen.Address, "ssl"})
		}
		if listen.QUIC {
			listens = append(listens, []string{listen.Address, "quic"})
		}
	}
	if err = p.SetListen(listens); err != nil {
		return err
	}
	// 域名
	domains, err := punycode.EncodeDomains(req.Domains)
	if err != nil {
		return err
	}
	if err = p.SetServerName(domains); err != nil {
		return err
	}
	// 首页文件
	if err = p.SetIndex(req.Index); err != nil {
		return err
	}
	// 运行目录
	if !io.Exists(req.Root) {
		return errors.New("运行目录不存在")
	}
	if err = p.SetRoot(req.Root); err != nil {
		return err
	}
	// 运行目录
	if !io.Exists(req.Path) {
		return errors.New("网站目录不存在")
	}
	website.Path = req.Path
	// PHP
	if err = p.SetPHP(req.PHP); err != nil {
		return err
	}
	// HTTPS
	certPath := filepath.Join(app.Root, "server/vhost/ssl", website.Name+".pem")
	keyPath := filepath.Join(app.Root, "server/vhost/ssl", website.Name+".key")
	if err = io.Write(certPath, req.SSLCertificate, 0644); err != nil {
		return err
	}
	if err = io.Write(keyPath, req.SSLCertificateKey, 0644); err != nil {
		return err
	}
	website.Https = req.HTTPS
	if req.HTTPS {
		if err = p.SetHTTPS(certPath, keyPath); err != nil {
			return err
		}
		if err = p.SetHTTPRedirect(req.HTTPRedirect); err != nil {
			return err
		}
		if err = p.SetHSTS(req.HSTS); err != nil {
			return err
		}
		if err = p.SetOCSP(req.OCSP); err != nil {
			return err
		}

	} else {
		if err = p.ClearSetHTTPS(); err != nil {
			return err
		}
		if err = p.SetHTTPRedirect(false); err != nil {
			return err
		}
		if err = p.SetHSTS(false); err != nil {
			return err
		}
		if err = p.SetOCSP(false); err != nil {
			return err
		}
	}
	// 防跨站
	if !strings.HasSuffix(req.Root, "/") {
		req.Root += "/"
	}
	userIni := filepath.Join(req.Root, ".user.ini")
	if req.OpenBasedir {
		_, _ = shell.Execf(`chattr -i '%s'`, userIni)
		if err = io.Write(userIni, fmt.Sprintf("open_basedir=%s:/tmp/", req.Root), 0644); err != nil {
			return err
		}
		_, _ = shell.Execf(`chattr +i '%s'`, userIni)
	} else {
		if io.Exists(userIni) {
			_, _ = shell.Execf(`chattr -i '%s'`, userIni)
			if err = io.Remove(userIni); err != nil {
				return err
			}
		}
	}

	if err = io.Write(filepath.Join(app.Root, "server/vhost", website.Name+".conf"), p.Dump(), 0644); err != nil {
		return err
	}
	if err = io.Write(filepath.Join(app.Root, "server/vhost/rewrite", website.Name+".conf"), req.Rewrite, 0644); err != nil {
		return err
	}

	if err = app.Orm.Save(website).Error; err != nil {
		return err
	}

	if err = systemctl.Reload("nginx"); err != nil {
		_, err = shell.Execf("nginx -t")
		return err
	}

	return nil
}

func (r *websiteRepo) Delete(req *request.WebsiteDelete) error {
	website := new(biz.Website)
	if err := app.Orm.Preload("Cert").Where("id", req.ID).First(website).Error; err != nil {
		return err
	}
	if website.Cert != nil {
		return errors.New("网站" + website.Name + "已绑定SSL证书，请先删除证书")
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

	if err := app.Orm.Delete(website).Error; err != nil {
		return err
	}

	if err := systemctl.Reload("nginx"); err != nil {
		_, err = shell.Execf("nginx -t")
		return err
	}

	return nil
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

	// 初始化nginx配置
	p, err := nginx.NewParser()
	if err != nil {
		return err
	}
	// 运行目录
	if err = p.SetRoot(website.Path); err != nil {
		return err
	}
	// 伪静态
	includes, comments, err := p.GetIncludes()
	if err != nil {
		return err
	}
	includes = append(includes, filepath.Join(app.Root, "server/vhost/rewrite", website.Name+".conf"))
	comments = append(comments, []string{"# 伪静态规则"})
	if err = p.SetIncludes(includes, comments); err != nil {
		return err
	}
	// 日志
	if err = p.SetAccessLog(filepath.Join(app.Root, "wwwlogs", website.Name+".log")); err != nil {
		return err
	}
	if err = p.SetErrorLog(filepath.Join(app.Root, "wwwlogs", website.Name+".error.log")); err != nil {
		return err
	}

	if err = io.Write(filepath.Join(app.Root, "server/vhost", website.Name+".conf"), p.Dump(), 0644); err != nil {
		return nil
	}
	if err = io.Write(filepath.Join(app.Root, "server/vhost/rewrite", website.Name+".conf"), "", 0644); err != nil {
		return nil
	}

	website.Status = true
	website.Https = false
	if err := app.Orm.Save(website).Error; err != nil {
		return err
	}

	if err = systemctl.Reload("nginx"); err != nil {
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

	// 解析nginx配置
	config, err := io.Read(filepath.Join(app.Root, "server/vhost", website.Name+".conf"))
	if err != nil {
		return err
	}
	p, err := nginx.NewParser(config)
	if err != nil {
		return err
	}

	// 取运行目录和默认文档
	root, rootComment, err := p.GetRootWithComment()
	if err != nil {
		return err
	}
	index, indexComment, err := p.GetIndexWithComment()
	if err != nil {
		return err
	}
	indexStr := strings.Join(index, " ")

	if status {
		if len(rootComment) == 0 {
			return fmt.Errorf("未找到运行目录注释")
		}
		if len(rootComment) != 1 {
			return fmt.Errorf("运行目录注释数量不正确，预期1个，实际%d个", len(rootComment))
		}
		rootComment[0] = strings.TrimPrefix(rootComment[0], "# ")
		if !io.Exists(rootComment[0]) {
			return fmt.Errorf("运行目录不存在")
		}
		if err = p.SetRoot(rootComment[0]); err != nil {
			return err
		}
		if len(indexComment) == 0 {
			return fmt.Errorf("未找到默认文档注释")
		}
		if len(indexComment) != 1 {
			return fmt.Errorf("默认文档注释数量不正确，预期1个，实际%d个", len(indexComment))
		}
		indexComment[0] = strings.TrimPrefix(indexComment[0], "# ")
		if err = p.SetIndex(strings.Fields(indexComment[0])); err != nil {
			return err
		}
	} else {
		if err = p.SetRootWithComment(filepath.Join(app.Root, "server/nginx/html"), []string{"# " + root}); err != nil {
			return err
		}
		if err = p.SetIndexWithComment([]string{"stop.html"}, []string{"# " + indexStr}); err != nil {
			return err
		}
	}

	if err = io.Write(filepath.Join(app.Root, "server/vhost", website.Name+".conf"), p.Dump(), 0644); err != nil {
		return err
	}

	website.Status = status
	if err = app.Orm.Save(website).Error; err != nil {
		return err
	}

	if err = systemctl.Reload("nginx"); err != nil {
		_, err = shell.Execf("nginx -t")
		return err
	}

	return nil
}
