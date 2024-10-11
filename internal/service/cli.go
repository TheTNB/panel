package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/go-rat/utils/hash"
	"github.com/goccy/go-yaml"
	"github.com/gookit/color"
	"github.com/urfave/cli/v3"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/str"
	"github.com/TheTNB/panel/pkg/systemctl"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/pkg/types"
)

type CliService struct {
	app  biz.AppRepo
	user biz.UserRepo
	hash hash.Hasher
}

func NewCliService() *CliService {
	return &CliService{
		app:  data.NewAppRepo(),
		user: data.NewUserRepo(),
		hash: hash.NewArgon2id(),
	}
}

func (s *CliService) Restart(ctx context.Context, cmd *cli.Command) error {
	return systemctl.Restart("panel")
}

func (s *CliService) Stop(ctx context.Context, cmd *cli.Command) error {
	return systemctl.Stop("panel")
}

func (s *CliService) Start(ctx context.Context, cmd *cli.Command) error {
	return systemctl.Start("panel")
}

func (s *CliService) Update(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) Info(ctx context.Context, cmd *cli.Command) error {
	user := new(biz.User)
	if err := app.Orm.Where("id", 1).First(user).Error; err != nil {
		return fmt.Errorf("获取管理员信息失败：%v", err)
	}

	password := str.RandomString(16)
	hashed, err := s.hash.Make(password)
	if err != nil {
		return fmt.Errorf("密码生成失败：%v", err)
	}
	user.Username = str.RandomString(8)
	user.Password = hashed
	if user.Email == "" {
		user.Email = str.RandomString(8) + "@example.com"
	}

	if err = app.Orm.Save(user).Error; err != nil {
		return fmt.Errorf("管理员信息保存失败：%v", err)
	}

	protocol := "http"
	if app.Conf.Bool("http.tls") {
		protocol = "https"
	}
	ip, err := tools.GetPublicIP()
	if err != nil {
		ip = "127.0.0.1"
	}
	port := app.Conf.String("http.port")
	if port == "" {
		return fmt.Errorf("端口获取失败")
	}
	entrance := app.Conf.String("http.entrance")
	if entrance == "" {
		return fmt.Errorf("入口获取失败")
	}

	color.Greenln(fmt.Sprintf("用户名: %s", user.Username))
	color.Greenln(fmt.Sprintf("密码: %s", password))
	color.Greenln(fmt.Sprintf("端口: %s", port))
	color.Greenln(fmt.Sprintf("入口: %s", entrance))
	color.Greenln(fmt.Sprintf("地址: %s://%s:%s%s", protocol, ip, port, entrance))

	return nil
}

func (s *CliService) UserList(ctx context.Context, cmd *cli.Command) error {
	users := make([]biz.User, 0)
	if err := app.Orm.Find(&users).Error; err != nil {
		return fmt.Errorf("获取用户列表失败：%v", err)
	}

	for _, user := range users {
		color.Greenln(fmt.Sprintf("ID: %d, 用户名: %s, 邮箱: %s, 创建日期: %s", user.ID, user.Username, user.Email, user.CreatedAt.Format("2006-01-02 15:04:05")))
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

	return nil
}

func (s *CliService) HTTPSOn(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) HTTPSOff(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) EntranceOn(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) EntranceOff(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) Port(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) WebsiteCreate(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) WebsiteRemove(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) WebsiteDelete(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) WebsiteWrite(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) BackupWebsite(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) BackupDatabase(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) BackupPanel(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) CutoffWebsite(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) AppInstall(ctx context.Context, cmd *cli.Command) error {
	channel := cmd.Args().Get(0)
	slug := cmd.Args().Get(1)
	if channel == "" || slug == "" {
		return fmt.Errorf("参数不能为空")
	}

	if err := s.app.Install(channel, slug); err != nil {
		return fmt.Errorf("应用安装失败：%v", err)
	}

	color.Greenln(fmt.Sprintf("已创建应用 %s 安装任务", slug))

	return nil
}

func (s *CliService) AppUnInstall(ctx context.Context, cmd *cli.Command) error {
	slug := cmd.Args().First()
	if slug == "" {
		return fmt.Errorf("参数不能为空")
	}

	if err := s.app.UnInstall(slug); err != nil {
		return fmt.Errorf("应用卸载失败：%v", err)
	}

	color.Greenln(fmt.Sprintf("已创建应用 %s 卸载任务", slug))

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

func (s *CliService) ClearTask(ctx context.Context, cmd *cli.Command) error {
	if err := app.Orm.Model(&biz.Task{}).
		Where("status", biz.TaskStatusRunning).Or("status", biz.TaskStatusWaiting).
		Update("status", biz.TaskStatusFailed).
		Error; err != nil {
		return fmt.Errorf("任务清理失败：%v", err)
	}

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

	settings := []biz.Setting{{Key: biz.SettingKeyName, Value: "耗子面板"}, {Key: biz.SettingKeyMonitor, Value: "1"}, {Key: biz.SettingKeyMonitorDays, Value: "30"}, {Key: biz.SettingKeyBackupPath, Value: filepath.Join(app.Root, "backup")}, {Key: biz.SettingKeyWebsitePath, Value: filepath.Join(app.Root, "wwwroot")}, {Key: biz.SettingKeyVersion, Value: app.Conf.String("app.version")}}
	if err := app.Orm.Create(&settings).Error; err != nil {
		return fmt.Errorf("初始化失败：%v", err)
	}

	value, err := hash.NewArgon2id().Make(str.RandomString(32))
	if err != nil {
		return fmt.Errorf("初始化失败：%v", err)
	}

	user := data.NewUserRepo()
	_, err = user.Create("admin", value)
	if err != nil {
		return fmt.Errorf("初始化失败：%v", err)
	}

	config := new(types.PanelConfig)
	cm := yaml.CommentMap{}
	raw, err := io.Read(filepath.Join(app.Root, "panel/config/config.yml"))
	if err != nil {
		return err
	}
	if err = yaml.UnmarshalWithOptions([]byte(raw), config, yaml.CommentToMap(cm)); err != nil {
		return err
	}

	config.App.Key = str.RandomString(32)
	config.HTTP.Entrance = "/" + str.RandomString(6)

	encoded, err := yaml.MarshalWithOptions(config, yaml.WithComment(cm))
	if err != nil {
		return err
	}
	if err = io.Write(filepath.Join(app.Root, "panel/config/config.yml"), string(encoded), 0644); err != nil {
		return err
	}

	// 初始化应用中心缓存
	return s.app.UpdateCache()
}
