package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-rat/utils/collect"
	"github.com/go-rat/utils/hash"
	"github.com/go-rat/utils/str"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/api"
	"github.com/TheTNB/panel/pkg/cert"
	"github.com/TheTNB/panel/pkg/firewall"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/ntp"
	"github.com/TheTNB/panel/pkg/systemctl"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/pkg/types"
)

type CliService struct {
	hr          string
	api         *api.API
	appRepo     biz.AppRepo
	userRepo    biz.UserRepo
	settingRepo biz.SettingRepo
	backupRepo  biz.BackupRepo
	websiteRepo biz.WebsiteRepo
	hash        hash.Hasher
}

func NewCliService() *CliService {
	return &CliService{
		hr:          `+----------------------------------------------------`,
		api:         api.NewAPI(app.Version),
		appRepo:     data.NewAppRepo(),
		userRepo:    data.NewUserRepo(),
		settingRepo: data.NewSettingRepo(),
		backupRepo:  data.NewBackupRepo(),
		websiteRepo: data.NewWebsiteRepo(),
		hash:        hash.NewArgon2id(),
	}
}

func (s *CliService) Restart(ctx context.Context, cmd *cli.Command) error {
	if err := systemctl.Restart("panel"); err != nil {
		return err
	}

	fmt.Println("面板服务已重启")
	return nil
}

func (s *CliService) Stop(ctx context.Context, cmd *cli.Command) error {
	if err := systemctl.Stop("panel"); err != nil {
		return err
	}

	fmt.Println("面板服务已停止")
	return nil
}

func (s *CliService) Start(ctx context.Context, cmd *cli.Command) error {
	if err := systemctl.Start("panel"); err != nil {
		return err
	}

	fmt.Println("面板服务已启动")
	return nil
}

func (s *CliService) Update(ctx context.Context, cmd *cli.Command) error {
	panel, err := s.api.LatestVersion()
	if err != nil {
		return fmt.Errorf("获取最新版本失败：%v", err)
	}

	download := collect.First(panel.Downloads)
	if download == nil {
		return fmt.Errorf("下载地址为空")
	}
	ver, url, checksum := panel.Version, download.URL, download.Checksum

	return s.settingRepo.UpdatePanel(ver, url, checksum)
}

func (s *CliService) Fix(ctx context.Context, cmd *cli.Command) error {
	return s.settingRepo.FixPanel()
}

func (s *CliService) Info(ctx context.Context, cmd *cli.Command) error {
	user := new(biz.User)
	if err := app.Orm.Where("id", 1).First(user).Error; err != nil {
		return fmt.Errorf("获取管理员信息失败：%v", err)
	}

	password := str.Random(16)
	hashed, err := s.hash.Make(password)
	if err != nil {
		return fmt.Errorf("密码生成失败：%v", err)
	}
	user.Username = str.Random(8)
	user.Password = hashed
	if user.Email == "" {
		user.Email = str.Random(8) + "@example.com"
	}

	if err = app.Orm.Save(user).Error; err != nil {
		return fmt.Errorf("管理员信息保存失败：%v", err)
	}

	protocol := "http"
	if app.Conf.Bool("http.tls") {
		protocol = "https"
	}

	port := app.Conf.String("http.port")
	if port == "" {
		return fmt.Errorf("端口获取失败")
	}
	entrance := app.Conf.String("http.entrance")
	if entrance == "" {
		return fmt.Errorf("入口获取失败")
	}

	fmt.Printf("用户名: %s\n", user.Username)
	fmt.Printf("密码: %s\n", password)
	fmt.Printf("端口: %s\n", port)
	fmt.Printf("入口: %s\n", entrance)

	lv4, err := tools.GetLocalIPv4()
	if err == nil {
		fmt.Printf("本地IPv4地址: %s://%s:%s%s\n", protocol, lv4, port, entrance)
	}
	lv6, err := tools.GetLocalIPv6()
	if err == nil {
		fmt.Printf("本地IPv6地址: %s://[%s]:%s%s\n", protocol, lv6, port, entrance)
	}
	rv4, err := tools.GetPublicIPv4()
	if err == nil {
		fmt.Printf("公网IPv4地址: %s://%s:%s%s\n", protocol, rv4, port, entrance)
	}
	rv6, err := tools.GetPublicIPv6()
	if err == nil {
		fmt.Printf("公网IPv6地址: %s://[%s]:%s%s\n", protocol, rv6, port, entrance)
	}

	fmt.Println("请根据自身网络情况自行选择合适的地址访问面板")
	fmt.Printf("如无法访问，请检查服务器运营商安全组和防火墙是否放行%s端口\n", port)
	fmt.Println("若仍无法访问，可尝试运行 panel-cli https off 关闭面板HTTPS")
	fmt.Println("警告：关闭面板HTTPS后，面板安全性将大大降低，请谨慎操作")

	return nil
}

func (s *CliService) UserList(ctx context.Context, cmd *cli.Command) error {
	users := make([]biz.User, 0)
	if err := app.Orm.Find(&users).Error; err != nil {
		return fmt.Errorf("获取用户列表失败：%v", err)
	}

	for _, user := range users {
		fmt.Printf("ID: %d, 用户名: %s, 邮箱: %s, 创建日期: %s\n", user.ID, user.Username, user.Email, user.CreatedAt.Format(time.DateTime))
	}

	return nil
}

func (s *CliService) UserName(ctx context.Context, cmd *cli.Command) error {
	user := new(biz.User)
	oldUsername := cmd.Args().Get(0)
	newUsername := cmd.Args().Get(1)
	if oldUsername == "" {
		return fmt.Errorf("旧用户名不能为空")
	}
	if newUsername == "" {
		return fmt.Errorf("新用户名不能为空")
	}

	if err := app.Orm.Where("username", oldUsername).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("用户不存在")
		} else {
			return fmt.Errorf("获取用户失败：%v", err)
		}
	}

	user.Username = newUsername
	if err := app.Orm.Save(user).Error; err != nil {
		return fmt.Errorf("用户名修改失败：%v", err)
	}

	fmt.Printf("用户 %s 修改为 %s 成功\n", oldUsername, newUsername)
	return nil
}

func (s *CliService) UserPassword(ctx context.Context, cmd *cli.Command) error {
	user := new(biz.User)
	username := cmd.Args().Get(0)
	password := cmd.Args().Get(1)
	if username == "" || password == "" {
		return fmt.Errorf("用户名和密码不能为空")
	}
	if len(password) < 6 {
		return fmt.Errorf("密码长度不能小于6")
	}

	if err := app.Orm.Where("username", username).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("用户不存在")
		} else {
			return fmt.Errorf("获取用户失败：%v", err)
		}
	}

	hashed, err := s.hash.Make(password)
	if err != nil {
		return fmt.Errorf("密码生成失败：%v", err)
	}
	user.Password = hashed
	if err = app.Orm.Save(user).Error; err != nil {
		return fmt.Errorf("密码修改失败：%v", err)
	}

	fmt.Printf("用户 %s 密码修改成功\n", username)
	return nil
}

func (s *CliService) HTTPSOn(ctx context.Context, cmd *cli.Command) error {
	config := new(types.PanelConfig)
	raw, err := io.Read("/usr/local/etc/panel/config.yml")
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal([]byte(raw), config); err != nil {
		return err
	}

	config.HTTP.TLS = true

	encoded, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err = io.Write("/usr/local/etc/panel/config.yml", string(encoded), 0700); err != nil {
		return err
	}

	fmt.Println("已开启HTTPS")
	return s.Restart(ctx, cmd)
}

func (s *CliService) HTTPSOff(ctx context.Context, cmd *cli.Command) error {
	config := new(types.PanelConfig)
	raw, err := io.Read("/usr/local/etc/panel/config.yml")
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal([]byte(raw), config); err != nil {
		return err
	}

	config.HTTP.TLS = false

	encoded, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err = io.Write("/usr/local/etc/panel/config.yml", string(encoded), 0700); err != nil {
		return err
	}

	fmt.Println("已关闭HTTPS")
	return s.Restart(ctx, cmd)
}

func (s *CliService) HTTPSGenerate(ctx context.Context, cmd *cli.Command) error {
	var names []string
	if lv4, err := tools.GetLocalIPv4(); err == nil {
		names = append(names, lv4)
	}
	if lv6, err := tools.GetLocalIPv6(); err == nil {
		names = append(names, lv6)
	}
	if rv4, err := tools.GetPublicIPv4(); err == nil {
		names = append(names, rv4)
	}
	if rv6, err := tools.GetPublicIPv6(); err == nil {
		names = append(names, rv6)
	}

	crt, key, err := cert.GenerateSelfSigned(names)
	if err != nil {
		return err
	}

	if err = io.Write(filepath.Join(app.Root, "panel/storage/cert.pem"), string(crt), 0644); err != nil {
		return err
	}
	if err = io.Write(filepath.Join(app.Root, "panel/storage/cert.key"), string(key), 0644); err != nil {
		return err
	}

	fmt.Println("已生成HTTPS证书")
	return s.Restart(ctx, cmd)
}

func (s *CliService) EntranceOn(ctx context.Context, cmd *cli.Command) error {
	config := new(types.PanelConfig)
	raw, err := io.Read("/usr/local/etc/panel/config.yml")
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal([]byte(raw), config); err != nil {
		return err
	}

	config.HTTP.Entrance = "/" + str.Random(6)

	encoded, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err = io.Write("/usr/local/etc/panel/config.yml", string(encoded), 0700); err != nil {
		return err
	}

	fmt.Println("已开启访问入口")
	fmt.Printf("访问入口：%s\n", config.HTTP.Entrance)
	return s.Restart(ctx, cmd)
}

func (s *CliService) EntranceOff(ctx context.Context, cmd *cli.Command) error {
	config := new(types.PanelConfig)
	raw, err := io.Read("/usr/local/etc/panel/config.yml")
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal([]byte(raw), config); err != nil {
		return err
	}

	config.HTTP.Entrance = "/"

	encoded, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err = io.Write("/usr/local/etc/panel/config.yml", string(encoded), 0700); err != nil {
		return err
	}

	fmt.Println("已关闭访问入口")
	return s.Restart(ctx, cmd)
}

func (s *CliService) Port(ctx context.Context, cmd *cli.Command) error {
	port := cast.ToInt(cmd.Args().First())
	if port < 1 || port > 65535 {
		return fmt.Errorf("端口范围错误")
	}

	config := new(types.PanelConfig)
	raw, err := io.Read("/usr/local/etc/panel/config.yml")
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal([]byte(raw), config); err != nil {
		return err
	}

	config.HTTP.Port = port

	encoded, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// 放行端口
	fw := firewall.NewFirewall()
	err = fw.Port(firewall.FireInfo{
		Type:      firewall.TypeNormal,
		PortStart: uint(config.HTTP.Port),
		PortEnd:   uint(config.HTTP.Port),
		Direction: firewall.DirectionIn,
		Strategy:  firewall.StrategyAccept,
	}, firewall.OperationAdd)
	if err != nil {
		return err
	}

	if err = io.Write("/usr/local/etc/panel/config.yml", string(encoded), 0700); err != nil {
		return err
	}

	fmt.Printf("已修改端口为 %d\n", port)
	return s.Restart(ctx, cmd)
}

func (s *CliService) WebsiteCreate(ctx context.Context, cmd *cli.Command) error {
	req := &request.WebsiteCreate{
		Name:    cmd.String("name"),
		Domains: cmd.StringSlice("domains"),
		Listens: cmd.StringSlice("listens"),
		Path:    cmd.String("path"),
		PHP:     int(cmd.Int("php")),
		DB:      false,
	}

	website, err := s.websiteRepo.Create(req)
	if err != nil {
		return err
	}

	fmt.Printf("网站 %s 创建成功\n", website.Name)
	return nil
}

func (s *CliService) WebsiteRemove(ctx context.Context, cmd *cli.Command) error {
	website, err := s.websiteRepo.GetByName(cmd.String("name"))
	if err != nil {
		return err
	}
	req := &request.WebsiteDelete{
		ID: website.ID,
	}

	if err = s.websiteRepo.Delete(req); err != nil {
		return err
	}

	fmt.Printf("网站 %s 移除成功\n", website.Name)
	return nil
}

func (s *CliService) WebsiteDelete(ctx context.Context, cmd *cli.Command) error {
	website, err := s.websiteRepo.GetByName(cmd.String("name"))
	if err != nil {
		return err
	}
	req := &request.WebsiteDelete{
		ID:   website.ID,
		Path: true,
		DB:   true,
	}

	if err = s.websiteRepo.Delete(req); err != nil {
		return err
	}

	fmt.Printf("网站 %s 删除成功\n", website.Name)
	return nil
}

func (s *CliService) WebsiteWrite(ctx context.Context, cmd *cli.Command) error {
	println("not support")
	return nil
}

func (s *CliService) BackupWebsite(ctx context.Context, cmd *cli.Command) error {
	fmt.Println(s.hr)
	fmt.Printf("★ 开始备份 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	fmt.Println("|-备份类型：网站")
	fmt.Printf("|-备份目标：%s\n", cmd.String("name"))
	if err := s.backupRepo.Create(biz.BackupTypeWebsite, cmd.String("name"), cmd.String("path")); err != nil {
		return fmt.Errorf("备份失败：%v", err)
	}
	fmt.Println(s.hr)
	fmt.Printf("☆ 备份成功 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	return nil
}

func (s *CliService) BackupDatabase(ctx context.Context, cmd *cli.Command) error {
	fmt.Println(s.hr)
	fmt.Printf("★ 开始备份 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	fmt.Println("|-备份类型：数据库")
	fmt.Printf("|-数据库：%s\n", cmd.String("type"))
	fmt.Printf("|-备份目标：%s\n", cmd.String("name"))
	if err := s.backupRepo.Create(biz.BackupType(cmd.String("type")), cmd.String("name"), cmd.String("path")); err != nil {
		return fmt.Errorf("备份失败：%v", err)
	}
	fmt.Println(s.hr)
	fmt.Printf("☆ 备份成功 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	return nil
}

func (s *CliService) BackupPanel(ctx context.Context, cmd *cli.Command) error {
	fmt.Println(s.hr)
	fmt.Printf("★ 开始备份 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	fmt.Println("|-备份类型：面板")
	if err := s.backupRepo.Create(biz.BackupTypePanel, "", cmd.String("path")); err != nil {
		return fmt.Errorf("备份失败：%v", err)
	}
	fmt.Println(s.hr)
	fmt.Printf("☆ 备份成功 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	return nil
}

func (s *CliService) BackupClear(ctx context.Context, cmd *cli.Command) error {
	path, err := s.backupRepo.GetPath(biz.BackupType(cmd.String("type")))
	if err != nil {
		return err
	}
	if cmd.String("path") != "" {
		path = cmd.String("path")
	}

	fmt.Println(s.hr)
	fmt.Printf("★ 开始清理 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	fmt.Printf("|-清理类型：%s\n", cmd.String("type"))
	fmt.Printf("|-清理目标：%s\n", cmd.String("file"))
	fmt.Printf("|-保留份数：%d\n", cmd.Int("save"))
	if err = s.backupRepo.ClearExpired(path, cmd.String("file"), int(cmd.Int("save"))); err != nil {
		return fmt.Errorf("清理失败：%v", err)
	}
	fmt.Println(s.hr)
	fmt.Printf("☆ 清理成功 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	return nil
}

func (s *CliService) CutoffWebsite(ctx context.Context, cmd *cli.Command) error {
	website, err := s.websiteRepo.GetByName(cmd.String("name"))
	if err != nil {
		return err
	}
	path := filepath.Join(app.Root, "wwwlogs")
	if cmd.String("path") != "" {
		path = cmd.String("path")
	}

	fmt.Println(s.hr)
	fmt.Printf("★ 开始切割日志 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	fmt.Println("|-切割类型：网站")
	fmt.Printf("|-切割目标：%s\n", website.Name)
	if err = s.backupRepo.CutoffLog(path, filepath.Join(app.Root, "wwwlogs", website.Name+".log")); err != nil {
		return err
	}
	fmt.Println(s.hr)
	fmt.Printf("☆ 切割成功 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	return nil
}

func (s *CliService) CutoffClear(ctx context.Context, cmd *cli.Command) error {
	if cmd.String("type") != "website" {
		return errors.New("当前仅支持网站日志切割")
	}
	path := filepath.Join(app.Root, "wwwlogs")
	if cmd.String("path") != "" {
		path = cmd.String("path")
	}

	fmt.Println(s.hr)
	fmt.Printf("★ 开始清理切割日志 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	fmt.Printf("|-清理类型：%s\n", cmd.String("type"))
	fmt.Printf("|-清理目标：%s\n", cmd.String("file"))
	fmt.Printf("|-保留份数：%d\n", cmd.Int("save"))
	if err := s.backupRepo.ClearExpired(path, cmd.String("file"), int(cmd.Int("save"))); err != nil {
		return err
	}
	fmt.Println(s.hr)
	fmt.Printf("☆ 清理成功 [%s]\n", time.Now().Format(time.DateTime))
	fmt.Println(s.hr)
	return nil
}

func (s *CliService) AppInstall(ctx context.Context, cmd *cli.Command) error {
	slug := cmd.Args().First()
	channel := cmd.Args().Get(1)
	if channel == "" || slug == "" {
		return fmt.Errorf("参数不能为空")
	}

	if err := s.appRepo.Install(channel, slug); err != nil {
		return fmt.Errorf("应用安装失败：%v", err)
	}

	fmt.Printf("应用 %s 安装完成\n", slug)

	return nil
}

func (s *CliService) AppUnInstall(ctx context.Context, cmd *cli.Command) error {
	slug := cmd.Args().First()
	if slug == "" {
		return fmt.Errorf("参数不能为空")
	}

	if err := s.appRepo.UnInstall(slug); err != nil {
		return fmt.Errorf("应用卸载失败：%v", err)
	}

	fmt.Printf("应用 %s 卸载完成\n", slug)

	return nil
}

func (s *CliService) AppUpdate(ctx context.Context, cmd *cli.Command) error {
	slug := cmd.Args().First()
	if slug == "" {
		return fmt.Errorf("参数不能为空")
	}

	if err := s.appRepo.Update(slug); err != nil {
		return fmt.Errorf("应用更新失败：%v", err)
	}

	fmt.Printf("应用 %s 更新完成\n", slug)

	return nil
}

func (s *CliService) AppWrite(ctx context.Context, cmd *cli.Command) error {
	slug := cmd.Args().Get(0)
	channel := cmd.Args().Get(1)
	version := cmd.Args().Get(2)
	if slug == "" || channel == "" || version == "" {
		return fmt.Errorf("参数不能为空")
	}

	newApp := new(biz.App)
	if err := app.Orm.Where("slug", slug).First(newApp).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("获取应用失败：%v", err)
		}
	}
	newApp.Slug = slug
	newApp.Channel = channel
	newApp.Version = version
	if err := app.Orm.Save(newApp).Error; err != nil {
		return fmt.Errorf("应用保存失败：%v", err)
	}

	return nil
}

func (s *CliService) AppRemove(ctx context.Context, cmd *cli.Command) error {
	slug := cmd.Args().First()
	if slug == "" {
		return fmt.Errorf("参数不能为空")
	}

	if err := app.Orm.Where("slug", slug).Delete(&biz.App{}).Error; err != nil {
		return fmt.Errorf("应用删除失败：%v", err)
	}

	return nil
}

func (s *CliService) SyncTime(ctx context.Context, cmd *cli.Command) error {
	now, err := ntp.Now()
	if err != nil {
		return err
	}

	if err = ntp.UpdateSystemTime(now); err != nil {
		return err
	}

	fmt.Println("时间同步成功")
	return nil
}

func (s *CliService) ClearTask(ctx context.Context, cmd *cli.Command) error {
	if err := app.Orm.Model(&biz.Task{}).
		Where("status", biz.TaskStatusRunning).Or("status", biz.TaskStatusWaiting).
		Update("status", biz.TaskStatusFailed).
		Error; err != nil {
		return fmt.Errorf("任务清理失败：%v", err)
	}

	fmt.Println("任务清理成功")
	return nil
}

func (s *CliService) GetSetting(ctx context.Context, cmd *cli.Command) error {
	key := cmd.Args().First()
	if key == "" {
		return fmt.Errorf("参数不能为空")
	}

	setting := new(biz.Setting)
	if err := app.Orm.Where("key", key).First(setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("设置不存在")
		}
		return fmt.Errorf("获取设置失败：%v", err)
	}

	fmt.Print(setting.Value)

	return nil
}

func (s *CliService) WriteSetting(ctx context.Context, cmd *cli.Command) error {
	key := cmd.Args().Get(0)
	value := cmd.Args().Get(1)
	if key == "" || value == "" {
		return fmt.Errorf("参数不能为空")
	}

	setting := new(biz.Setting)
	if err := app.Orm.Where("key", key).First(setting).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("获取设置失败：%v", err)
		}
	}
	setting.Key = biz.SettingKey(key)
	setting.Value = value
	if err := app.Orm.Save(setting).Error; err != nil {
		return fmt.Errorf("设置保存失败：%v", err)
	}

	return nil
}

func (s *CliService) RemoveSetting(ctx context.Context, cmd *cli.Command) error {
	key := cmd.Args().First()
	if key == "" {
		return fmt.Errorf("参数不能为空")
	}

	if err := app.Orm.Where("key", key).Delete(&biz.Setting{}).Error; err != nil {
		return fmt.Errorf("设置删除失败：%v", err)
	}

	return nil
}

func (s *CliService) Init(ctx context.Context, cmd *cli.Command) error {
	var check biz.User
	if err := app.Orm.First(&check).Error; err == nil {
		return fmt.Errorf("已经初始化过了")
	}

	settings := []biz.Setting{
		{Key: biz.SettingKeyName, Value: "耗子面板"},
		{Key: biz.SettingKeyMonitor, Value: "1"},
		{Key: biz.SettingKeyMonitorDays, Value: "30"},
		{Key: biz.SettingKeyBackupPath, Value: filepath.Join(app.Root, "backup")},
		{Key: biz.SettingKeyWebsitePath, Value: filepath.Join(app.Root, "wwwroot")},
		{Key: biz.SettingKeyVersion, Value: app.Version},
	}
	if err := app.Orm.Create(&settings).Error; err != nil {
		return fmt.Errorf("初始化失败：%v", err)
	}

	value, err := hash.NewArgon2id().Make(str.Random(32))
	if err != nil {
		return fmt.Errorf("初始化失败：%v", err)
	}

	user := data.NewUserRepo()
	_, err = user.Create("admin", value)
	if err != nil {
		return fmt.Errorf("初始化失败：%v", err)
	}

	if err = s.HTTPSGenerate(ctx, cmd); err != nil {
		return fmt.Errorf("初始化失败：%v", err)
	}

	config := new(types.PanelConfig)
	raw, err := io.Read("/usr/local/etc/panel/config.yml")
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal([]byte(raw), config); err != nil {
		return err
	}

	config.App.Key = str.Random(32)
	config.HTTP.Entrance = "/" + str.Random(6)

	encoded, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	if err = io.Write("/usr/local/etc/panel/config.yml", string(encoded), 0700); err != nil {
		return err
	}

	// 初始化应用中心缓存
	return s.appRepo.UpdateCache()
}
