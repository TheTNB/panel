package data

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/io"
)

type settingRepo struct{}

func NewSettingRepo() biz.SettingRepo {
	return &settingRepo{}
}

func (r *settingRepo) Get(key biz.SettingKey, defaultValue ...string) (string, error) {
	setting := new(biz.Setting)
	if err := app.Orm.Where("key = ?", key).First(setting).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
	}

	if setting.Value == "" && len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return setting.Value, nil
}

func (r *settingRepo) Set(key biz.SettingKey, value string) error {
	setting := new(biz.Setting)
	if err := app.Orm.Where("key = ?", key).FirstOrInit(setting).Error; err != nil {
		return err
	}

	setting.Value = value
	return app.Orm.Save(setting).Error
}

func (r *settingRepo) Delete(key biz.SettingKey) error {
	setting := new(biz.Setting)
	if err := app.Orm.Where("key = ?", key).Delete(setting).Error; err != nil {
		return err
	}

	return nil
}

func (r *settingRepo) GetPanelSetting(ctx context.Context) (*request.PanelSetting, error) {
	name := new(biz.Setting)
	if err := app.Orm.Where("key = ?", biz.SettingKeyName).First(name).Error; err != nil {
		return nil, err
	}
	websitePath := new(biz.Setting)
	if err := app.Orm.Where("key = ?", biz.SettingKeyWebsitePath).First(websitePath).Error; err != nil {
		return nil, err
	}
	backupPath := new(biz.Setting)
	if err := app.Orm.Where("key = ?", biz.SettingKeyBackupPath).First(backupPath).Error; err != nil {
		return nil, err
	}

	userID := cast.ToUint(ctx.Value("user_id"))
	if userID == 0 {
		return nil, errors.New("获取用户 ID 失败")
	}
	user := new(biz.User)
	if err := app.Orm.Where("id = ?", userID).First(user).Error; err != nil {
		return nil, err
	}

	cert, err := io.Read(filepath.Join(app.Root, "panel/storage/cert.pem"))
	if err != nil {
		return nil, err
	}
	key, err := io.Read(filepath.Join(app.Root, "panel/storage/cert.key"))
	if err != nil {
		return nil, err
	}

	return &request.PanelSetting{
		Name:        name.Value,
		Locale:      app.Conf.String("app.locale"),
		Entrance:    app.Conf.String("http.entrance"),
		WebsitePath: websitePath.Value,
		BackupPath:  backupPath.Value,
		Username:    user.Username,
		Email:       user.Email,
		Port:        app.Conf.Int("http.port"),
		HTTPS:       app.Conf.Bool("http.tls"),
		Cert:        cert,
		Key:         key,
	}, nil
}

func (r *settingRepo) UpdatePanelSetting(setting *request.PanelSetting) error {
	if err := r.Set(biz.SettingKeyName, setting.Name); err != nil {
		return err
	}
	if err := r.Set(biz.SettingKeyWebsitePath, setting.WebsitePath); err != nil {
		return err
	}
	if err := r.Set(biz.SettingKeyBackupPath, setting.BackupPath); err != nil {
		return err
	}

	// TODO fix other settings

	return nil
}
