package biz

import "github.com/golang-module/carbon/v2"

const (
	SettingKeyName              = "name"
	SettingKeyVersion           = "version"
	SettingKeyMonitor           = "monitor"
	SettingKeyMonitorDays       = "monitor_days"
	SettingKeyBackupPath        = "backup_path"
	SettingKeyWebsitePath       = "website_path"
	SettingKeyMysqlRootPassword = "mysql_root_password"
	SettingKeySshHost           = "ssh_host"
	SettingKeySshPort           = "ssh_port"
	SettingKeySshUser           = "ssh_user"
	SettingKeySshPassword       = "ssh_password"
)

type Setting struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Key       string          `gorm:"not null;unique" json:"key"`
	Value     string          `gorm:"not null" json:"value"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

type SettingRepo interface {
	Get(key string, defaultValue ...string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}
