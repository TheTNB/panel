package models

import "github.com/goravel/framework/support/carbon"

const (
	SettingKeyName              = "name"
	SettingKeyVersion           = "version"
	SettingKeyMonitor           = "monitor"
	SettingKeyMonitorDays       = "monitor_days"
	SettingKeyBackupPath        = "backup_path"
	SettingKeyWebsitePath       = "website_path"
	SettingKeyEntrance          = "entrance"
	SettingKeyMysqlRootPassword = "mysql_root_password"
	SettingKeySshHost           = "ssh_host"
	SettingKeySshPort           = "ssh_port"
	SettingKeySshUser           = "ssh_user"
	SettingKeySshPassword       = "ssh_password"
)

type Setting struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Key       string          `gorm:"unique;not null" json:"key"`
	Value     string          `gorm:"default:''" json:"value"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
