package biz

import "github.com/TheTNB/panel/pkg/types"

type BackupRepo interface {
	List(typ string) ([]*types.BackupFile, error)
	Create(typ, target string, path ...string) error
	Delete(typ, name string) error
	CleanExpired(path, prefix string, save int) error
	CutoffLog(path, target string) error
	GetPath(typ string) (string, error)
}
