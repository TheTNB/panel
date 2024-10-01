package service

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/go-rat/utils/hash"
	"github.com/gookit/color"
	"github.com/urfave/cli/v3"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/pkg/str"
	"github.com/TheTNB/panel/pkg/systemctl"
	"github.com/TheTNB/panel/pkg/tools"
)

type CliService struct {
	repo biz.UserRepo
	hash hash.Hasher
}

func NewCliService() *CliService {
	return &CliService{
		repo: data.NewUserRepo(),
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

	err = app.Orm.Save(user).Error
	if err != nil {
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
	println("Hello, World!")
	return nil
}

func (s *CliService) UserName(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) UserPassword(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
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
	println("Hello, World!")
	return nil
}

func (s *CliService) AppUnInstall(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) AppWrite(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) AppRemove(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) ClearTask(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) WriteSetting(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) RemoveSetting(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) Init(ctx context.Context, cmd *cli.Command) error {
	var check biz.User
	if err := app.Orm.First(&check).Error; err == nil {
		return fmt.Errorf("已经初始化过了")
	}

	settings := []biz.Setting{{Key: biz.SettingKeyName, Value: "耗子面板"}, {Key: biz.SettingKeyMonitor, Value: "1"}, {Key: biz.SettingKeyMonitorDays, Value: "30"}, {Key: biz.SettingKeyBackupPath, Value: filepath.Join(app.Root, "backup")}, {Key: biz.SettingKeyWebsitePath, Value: filepath.Join(app.Root, "wwwroot")}, {Key: biz.SettingKeyVersion, Value: app.Conf.String("app.version")}}
	if err := app.Orm.Create(&settings).Error; err != nil {
		return fmt.Errorf("初始化失败: %v", err)
	}

	value, err := hash.NewArgon2id().Make(str.RandomString(32))
	if err != nil {
		return fmt.Errorf("初始化失败: %v", err)
	}

	user := data.NewUserRepo()
	_, err = user.Create("admin", value)
	if err != nil {
		return fmt.Errorf("初始化失败: %v", err)
	}

	return nil
}
