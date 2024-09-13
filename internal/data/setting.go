package data

import (
	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
)

type settingRepo struct{}

func NewSettingRepo() biz.SettingRepo {
	return &settingRepo{}
}

func (r *settingRepo) Get(key string, defaultValue ...string) (string, error) {
	setting := new(biz.Setting)
	if err := app.Orm.Where("key = ?", key).First(setting).Error; err != nil {
		return "", err
	}

	if setting.Value == "" && len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return setting.Value, nil
}

func (r *settingRepo) Set(key, value string) error {
	setting := new(biz.Setting)
	if err := app.Orm.Where("key = ?", key).First(setting).Error; err != nil {
		return err
	}

	setting.Value = value
	return app.Orm.Save(setting).Error
}

func (r *settingRepo) Delete(key string) error {
	setting := new(biz.Setting)
	if err := app.Orm.Where("key = ?", key).Delete(setting).Error; err != nil {
		return err
	}

	return nil
}
