// Package services 备份服务
package services

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/support/carbon"

	"panel/app/models"
	"panel/pkg/tools"
)

type Backup interface {
	WebsiteList() ([]BackupFile, error)
	WebSiteBackup(website models.Website) error
	WebsiteRestore(website models.Website, backupFile string) error
	MysqlList() ([]BackupFile, error)
	MysqlBackup(database string) error
	MysqlRestore(database string, backupFile string) error
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

// WebsiteList 网站备份列表
func (s *BackupImpl) WebsiteList() ([]BackupFile, error) {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return []BackupFile{}, nil
	}

	backupPath += "/website"
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}

	files, err := os.ReadDir(backupPath)
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

// WebSiteBackup 网站备份
func (s *BackupImpl) WebSiteBackup(website models.Website) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return errors.New("未正确配置备份路径")
	}

	backupPath += "/website"
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}

	backupFile := backupPath + "/" + website.Name + "_" + carbon.Now().ToShortDateTimeString() + ".zip"
	tools.ExecShell(`cd '` + website.Path + `' && zip -r '` + backupFile + `' .`)

	return nil
}

// WebsiteRestore 网站恢复
func (s *BackupImpl) WebsiteRestore(website models.Website, backupFile string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return errors.New("未正确配置备份路径")
	}

	backupPath += "/website"
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}

	backupFile = backupPath + "/" + backupFile
	if !tools.Exists(backupFile) {
		return errors.New("备份文件不存在")
	}

	tools.ExecShell(`rm -rf '` + website.Path + `/*'`)
	tools.ExecShell(`unzip -o '` + backupFile + `' -d '` + website.Path + `' 2>&1`)
	tools.Chmod(website.Path, 0755)
	tools.Chown(website.Path, "www", "www")

	return nil
}

// MysqlList MySQL备份列表
func (s *BackupImpl) MysqlList() ([]BackupFile, error) {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return []BackupFile{}, nil
	}

	backupPath += "/mysql"
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}

	files, err := os.ReadDir(backupPath)
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

// MysqlBackup MySQL备份
func (s *BackupImpl) MysqlBackup(database string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath) + "/mysql"
	rootPassword := s.setting.Get(models.SettingKeyMysqlRootPassword)
	backupFile := database + "_" + carbon.Now().ToShortDateTimeString() + ".sql"
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}
	err := os.Setenv("MYSQL_PWD", rootPassword)
	if err != nil {
		return err
	}

	tools.ExecShell("/www/server/mysql/bin/mysqldump -uroot " + database + " > " + backupPath + "/" + backupFile)
	tools.ExecShell("cd " + backupPath + " && zip -r " + backupPath + "/" + backupFile + ".zip " + backupFile)
	tools.RemoveFile(backupPath + "/" + backupFile)
	_ = os.Unsetenv("MYSQL_PWD")

	return nil
}

// MysqlRestore MySQL恢复
func (s *BackupImpl) MysqlRestore(database string, backupFile string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath) + "/mysql"
	rootPassword := s.setting.Get(models.SettingKeyMysqlRootPassword)
	ext := filepath.Ext(backupFile)
	backupFile = backupPath + "/" + backupFile
	if !tools.Exists(backupFile) {
		return errors.New("备份文件不存在")
	}

	err := os.Setenv("MYSQL_PWD", rootPassword)
	if err != nil {
		return err
	}

	switch ext {
	case ".zip":
		tools.ExecShell("unzip -o " + backupFile + " -d " + backupPath)
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".gz":
		if strings.HasSuffix(backupFile, ".tar.gz") {
			// 解压.tar.gz文件
			tools.ExecShell("tar -zxvf " + backupFile + " -C " + backupPath)
			backupFile = strings.TrimSuffix(backupFile, ".tar.gz")
		} else {
			// 解压.gz文件
			tools.ExecShell("gzip -d " + backupFile)
			backupFile = strings.TrimSuffix(backupFile, ext)
		}
	case ".bz2":
		tools.ExecShell("bzip2 -d " + backupFile)
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".tar":
		tools.ExecShell("tar -xvf " + backupFile + " -C " + backupPath)
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".rar":
		tools.ExecShell("unrar x " + backupFile + " " + backupPath)
		backupFile = strings.TrimSuffix(backupFile, ext)
	}

	if !tools.Exists(backupFile) {
		return errors.New("自动解压失败，请手动解压")
	}

	tools.ExecShell("/www/server/mysql/bin/mysql -uroot " + database + " < " + backupFile)
	_ = os.Unsetenv("MYSQL_PWD")

	return nil
}
