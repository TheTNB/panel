// Package services 备份服务
package services

type Backup interface {
	BackupWebSite(string) error
	BackupDatabase(string, string) error
}
