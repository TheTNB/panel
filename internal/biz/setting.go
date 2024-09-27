package biz

import (
	"github.com/golang-module/carbon/v2"

	"github.com/TheTNB/panel/internal/http/request"
)

type SettingKey string

const (
	SettingKeyName                SettingKey = "name"
	SettingKeyVersion             SettingKey = "version"
	SettingKeyMonitor             SettingKey = "monitor"
	SettingKeyMonitorDays         SettingKey = "monitor_days"
	SettingKeyBackupPath          SettingKey = "backup_path"
	SettingKeyWebsitePath         SettingKey = "website_path"
	SettingKeyPerconaRootPassword SettingKey = "percona_root_password"
	SettingKeySshHost             SettingKey = "ssh_host"
	SettingKeySshPort             SettingKey = "ssh_port"
	SettingKeySshUser             SettingKey = "ssh_user"
	SettingKeySshPassword         SettingKey = "ssh_password"
)

type Setting struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Key       SettingKey      `gorm:"not null;unique" json:"key"`
	Value     string          `gorm:"not null" json:"value"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

type SettingRepo interface {
	Get(key SettingKey, defaultValue ...string) (string, error)
	Set(key SettingKey, value string) error
	Delete(key SettingKey) error
	GetPanelSetting() (*request.PanelSetting, error)
	UpdatePanelSetting(setting *request.PanelSetting) error
}
