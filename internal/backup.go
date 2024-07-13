package internal

import (
	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type Backup interface {
	WebsiteList() ([]types.BackupFile, error)
	WebSiteBackup(website models.Website) error
	WebsiteRestore(website models.Website, backupFile string) error
	MysqlList() ([]types.BackupFile, error)
	MysqlBackup(database string) error
	MysqlRestore(database string, backupFile string) error
	PostgresqlList() ([]types.BackupFile, error)
	PostgresqlBackup(database string) error
	PostgresqlRestore(database string, backupFile string) error
}
