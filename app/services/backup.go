// Package services 备份服务
package services

import (
	"errors"
	"os"

	"github.com/goravel/framework/support/carbon"

	"panel/app/models"
	"panel/pkg/tools"
)

type Backup interface {
	WebsiteList() ([]BackupFile, error)
	WebSiteBackup(website models.Website) error
}

type BackupFile struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

type BackupImpl struct {
	setting Setting
}

func NewBackupImpl() *BackupImpl {
	return &BackupImpl{
		setting: NewSettingImpl(),
	}
}

func (s *BackupImpl) WebsiteList() ([]BackupFile, error) {
	path := s.setting.Get(models.SettingKeyBackupPath)
	if len(path) == 0 {
		return []BackupFile{}, nil
	}

	path += "/website"

	files, err := os.ReadDir(path)
	if err != nil {
		return []BackupFile{}, err
	}
	var backupList []BackupFile
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		backupList = append(backupList, BackupFile{
			Name: file.Name(),
			Size: tools.FormatBytes(float64(info.Size())),
		})
	}

	return backupList, nil
}

func (s *BackupImpl) WebSiteBackup(website models.Website) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return errors.New("未正确配置备份路径")
	}

	backupPath += "/website"
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}

	backupFile := backupPath + "/" + website.Name + carbon.Now().ToShortDateTimeString() + ".zip"
	tools.ExecShell("cd " + website.Path + " && zip -r " + backupFile + " .")

	return nil
}
