package biz

import "github.com/TheTNB/panel/pkg/types"

type BackupType string

const (
	BackupTypePath     BackupType = "path"
	BackupTypeWebsite  BackupType = "website"
	BackupTypeMySQL    BackupType = "mysql"
	BackupTypePostgres BackupType = "postgres"
	BackupTypeRedis    BackupType = "redis"
	BackupTypePanel    BackupType = "panel"
)

type BackupRepo interface {
	List(typ BackupType) ([]*types.BackupFile, error)
	Create(typ BackupType, target string, path ...string) error
	Delete(typ BackupType, name string) error
	Restore(typ BackupType, backup, target string) error
	ClearExpired(path, prefix string, save int) error
	CutoffLog(path, target string) error
	GetPath(typ BackupType) (string, error)
}
