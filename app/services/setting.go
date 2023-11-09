// Package services 设置服务
package services

import (
	"github.com/goravel/framework/facades"
	"panel/app/models"
)

type Setting interface {
	Get(key string, defaultValue ...string) string
	Set(key, value string) error
	Delete(key string) error
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
	if err := facades.Orm().Query().UpdateOrCreate(&setting, models.Setting{Key: key}, models.Setting{Value: value}); err != nil {
		return err
	}

	return nil
}

// Delete 删除设置
func (r *SettingImpl) Delete(key string) error {
	var setting models.Setting
	if _, err := facades.Orm().Query().Where("key", key).Delete(&setting); err != nil {
		return err
	}

	return nil
}
