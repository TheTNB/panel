package data

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/go-rat/utils/hash"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/types"
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

func (r *settingRepo) UpdatePanelSetting(ctx context.Context, setting *request.PanelSetting) (bool, error) {
	if err := r.Set(biz.SettingKeyName, setting.Name); err != nil {
		return false, err
	}
	if err := r.Set(biz.SettingKeyWebsitePath, setting.WebsitePath); err != nil {
		return false, err
	}
	if err := r.Set(biz.SettingKeyBackupPath, setting.BackupPath); err != nil {
		return false, err
	}

	if err := io.Write(filepath.Join(app.Root, "panel/storage/cert.pem"), setting.Cert, 0644); err != nil {
		return false, err
	}
	if err := io.Write(filepath.Join(app.Root, "panel/storage/cert.key"), setting.Key, 0644); err != nil {
		return false, err
	}

	restartFlag := false
	config := new(types.PanelConfig)
	cm := yaml.CommentMap{}
	raw, err := io.Read("config/config.yml")
	if err != nil {
		return false, err
	}
	if err = yaml.UnmarshalWithOptions([]byte(raw), &config, yaml.CommentToMap(cm)); err != nil {
		return false, err
	}

	config.App.Locale = setting.Locale
	config.HTTP.Port = setting.Port
	config.HTTP.Entrance = setting.Entrance
	config.HTTP.TLS = setting.HTTPS

	encoded, err := yaml.MarshalWithOptions(config, yaml.WithComment(cm))
	if err != nil {
		return false, err
	}
	if err = io.Write("config/config.yml", string(encoded), 0644); err != nil {
		return false, err
	}
	if raw != string(encoded) {
		restartFlag = true
	}

	user := new(biz.User)
	userID := cast.ToUint(ctx.Value("user_id"))
	if err = app.Orm.Where("id = ?", userID).First(user).Error; err != nil {
		return false, err
	}

	user.Username = setting.Username
	user.Email = setting.Email
	if setting.Password != "" {
		value, err := hash.NewArgon2id().Make(setting.Password)
		if err != nil {
			return false, err
		}
		user.Password = value
	}
	if err = app.Orm.Save(user).Error; err != nil {
		return false, err
	}

	return restartFlag, nil
}
