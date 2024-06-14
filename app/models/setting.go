package models

import "github.com/goravel/framework/database/orm"

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
	orm.Model
	Key   string `gorm:"not null;unique" json:"key"`
	Value string `gorm:"not null" json:"value"`
}
