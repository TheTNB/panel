package biz

import "github.com/TheTNB/panel/pkg/types"

type BackupRepo interface {
	List(typ string) ([]*types.BackupFile, error)
	Create(typ, target string, path ...string) error
	Delete(typ, name string) error
	Restore(typ, backup, target string) error
	ClearExpired(path, prefix string, save int) error
	CutoffLog(path, target string) error
	GetPath(typ string) (string, error)
}
