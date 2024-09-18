package data

import (
	"errors"

	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/internal/panel"
)

type settingRepo struct{}

func NewSettingRepo() biz.SettingRepo {
	return &settingRepo{}
}

func (r *settingRepo) Get(key biz.SettingKey, defaultValue ...string) (string, error) {
	setting := new(biz.Setting)
	if err := panel.Orm.Where("key = ?", key).First(setting).Error; err != nil {
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
	if err := panel.Orm.Where("key = ?", key).FirstOrInit(setting).Error; err != nil {
		return err
	}

	setting.Value = value
	return panel.Orm.Save(setting).Error
}

func (r *settingRepo) Delete(key biz.SettingKey) error {
	setting := new(biz.Setting)
	if err := panel.Orm.Where("key = ?", key).Delete(setting).Error; err != nil {
		return err
	}

	return nil
}

func (r *settingRepo) GetPanelSetting() (*request.PanelSetting, error) {
	setting := new(biz.Setting)
	if err := panel.Orm.Where("key = ?", biz.SettingKeyName).First(setting).Error; err != nil {
		return nil, err
	}

	// TODO fix

	return &request.PanelSetting{
		Name: setting.Value,
	}, nil
}

func (r *settingRepo) UpdatePanelSetting(setting *request.PanelSetting) error {
	if err := r.Set(biz.SettingKeyName, setting.Name); err != nil {
		return err
	}

	// TODO fix

	return nil
}
