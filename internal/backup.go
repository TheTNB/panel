package internal

import "github.com/TheTNB/panel/app/models"

type Backup interface {
	WebsiteList() ([]BackupFile, error)
	WebSiteBackup(website models.Website) error
	WebsiteRestore(website models.Website, backupFile string) error
	MysqlList() ([]BackupFile, error)
	MysqlBackup(database string) error
	MysqlRestore(database string, backupFile string) error
	PostgresqlList() ([]BackupFile, error)
	PostgresqlBackup(database string) error
	PostgresqlRestore(database string, backupFile string) error
}

type BackupFile struct {
	Name string `json:"name"`
	Size string `json:"size"`
}
