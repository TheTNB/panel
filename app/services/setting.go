package services

import (
	"github.com/goravel/framework/facades"
	"panel/app/models"
)

type Setting interface {
	Get(key string, defaultValue ...string) string
	Set(key, value string) error
}

type SettingImpl struct {
}

func NewSettingImpl() *SettingImpl {
	return &SettingImpl{}
}

// Get 获取设置
func (r *SettingImpl) Get(key string, defaultValue ...string) string {
	var setting models.Setting
	if err := facades.Orm().Query().Where("key", key).FirstOrFail(&setting); err != nil {
		if len(defaultValue) == 0 {
			return ""
		}

		return defaultValue[0]
	}

	if len(setting.Value) == 0 {
		if len(defaultValue) == 0 {
			return ""
		}

		return defaultValue[0]
	}

	return setting.Value
}

// Set 更新或创建设置
func (r *SettingImpl) Set(key, value string) error {
	var setting models.Setting
	if err := facades.Orm().Query().Where("key", key).UpdateOrCreate(&setting, models.Setting{Key: key}, models.Setting{Value: value}); err != nil {
		return err
	}

	return nil
}
