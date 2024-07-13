// Package services 备份服务
package services

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/support/carbon"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/str"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type BackupImpl struct {
	setting internal.Setting
}

func NewBackupImpl() *BackupImpl {
	return &BackupImpl{
		setting: NewSettingImpl(),
	}
}

// WebsiteList 网站备份列表
func (s *BackupImpl) WebsiteList() ([]types.BackupFile, error) {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return []types.BackupFile{}, nil
	}

	backupPath += "/website"
	if !io.Exists(backupPath) {
		if err := io.Mkdir(backupPath, 0644); err != nil {
			return []types.BackupFile{}, err
		}
	}

	files, err := os.ReadDir(backupPath)
	if err != nil {
		return []types.BackupFile{}, err
	}
	var backupList []types.BackupFile
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		backupList = append(backupList, types.BackupFile{
			Name: file.Name(),
			Size: str.FormatBytes(float64(info.Size())),
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
	if !io.Exists(backupPath) {
		if err := io.Mkdir(backupPath, 0644); err != nil {
			return err
		}
	}

	backupFile := backupPath + "/" + website.Name + "_" + carbon.Now().ToShortDateTimeString() + ".zip"
	if _, err := shell.Execf(`cd '` + website.Path + `' && zip -r '` + backupFile + `' .`); err != nil {
		return err
	}

	return nil
}

// WebsiteRestore 网站恢复
func (s *BackupImpl) WebsiteRestore(website models.Website, backupFile string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return errors.New("未正确配置备份路径")
	}

	backupPath += "/website"
	if !io.Exists(backupPath) {
		if err := io.Mkdir(backupPath, 0644); err != nil {
			return err
		}
	}

	backupFile = backupPath + "/" + backupFile
	if !io.Exists(backupFile) {
		return errors.New("备份文件不存在")
	}

	if err := io.Remove(website.Path); err != nil {
		return err
	}
	if err := io.UnArchive(backupFile, website.Path); err != nil {
		return err
	}
	if err := io.Chmod(website.Path, 0755); err != nil {
		return err
	}
	if err := io.Chown(website.Path, "www", "www"); err != nil {
		return err
	}

	return nil
}

// MysqlList MySQL备份列表
func (s *BackupImpl) MysqlList() ([]types.BackupFile, error) {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return []types.BackupFile{}, nil
	}

	backupPath += "/mysql"
	if !io.Exists(backupPath) {
		if err := io.Mkdir(backupPath, 0644); err != nil {
			return []types.BackupFile{}, err
		}
	}

	files, err := os.ReadDir(backupPath)
	if err != nil {
		return []types.BackupFile{}, err
	}
	var backupList []types.BackupFile
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		backupList = append(backupList, types.BackupFile{
			Name: file.Name(),
			Size: str.FormatBytes(float64(info.Size())),
		})
	}

	return backupList, nil
}

// MysqlBackup MySQL备份
func (s *BackupImpl) MysqlBackup(database string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath) + "/mysql"
	rootPassword := s.setting.Get(models.SettingKeyMysqlRootPassword)
	backupFile := database + "_" + carbon.Now().ToShortDateTimeString() + ".sql"
	if !io.Exists(backupPath) {
		if err := io.Mkdir(backupPath, 0644); err != nil {
			return err
		}
	}
	err := os.Setenv("MYSQL_PWD", rootPassword)
	if err != nil {
		return err
	}

	if _, err := shell.Execf("mysqldump -uroot " + database + " > " + backupPath + "/" + backupFile); err != nil {
		return err
	}
	if _, err := shell.Execf("cd " + backupPath + " && zip -r " + backupPath + "/" + backupFile + ".zip " + backupFile); err != nil {
		return err
	}
	if err := io.Remove(backupPath + "/" + backupFile); err != nil {
		return err
	}

	return os.Unsetenv("MYSQL_PWD")
}

// MysqlRestore MySQL恢复
func (s *BackupImpl) MysqlRestore(database string, backupFile string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath) + "/mysql"
	rootPassword := s.setting.Get(models.SettingKeyMysqlRootPassword)
	backupFullPath := filepath.Join(backupPath, backupFile)
	if !io.Exists(backupFullPath) {
		return errors.New("备份文件不存在")
	}

	if err := os.Setenv("MYSQL_PWD", rootPassword); err != nil {
		return err
	}

	tempDir, err := io.TempDir(backupFile)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(backupFile, ".sql") {
		backupFile = "" // 置空，防止干扰后续判断
		if err = io.UnArchive(backupFullPath, tempDir); err != nil {
			return err
		}
		if files, err := os.ReadDir(tempDir); err == nil {
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".sql") {
					backupFile = filepath.Base(file.Name())
					break
				}
			}
		}
	} else {
		if err = io.Cp(backupFullPath, filepath.Join(tempDir, backupFile)); err != nil {
			return err
		}
	}

	if len(backupFile) == 0 {
		return errors.New("无法找到备份文件")
	}

	if _, err = shell.Execf("mysql -uroot " + database + " < " + filepath.Join(tempDir, backupFile)); err != nil {
		return err
	}

	if err = io.Remove(tempDir); err != nil {
		return err
	}

	return os.Unsetenv("MYSQL_PWD")
}

// PostgresqlList PostgreSQL备份列表
func (s *BackupImpl) PostgresqlList() ([]types.BackupFile, error) {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return []types.BackupFile{}, nil
	}

	backupPath += "/postgresql"
	if !io.Exists(backupPath) {
		if err := io.Mkdir(backupPath, 0644); err != nil {
			return []types.BackupFile{}, err
		}
	}

	files, err := os.ReadDir(backupPath)
	if err != nil {
		return []types.BackupFile{}, err
	}
	var backupList []types.BackupFile
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		backupList = append(backupList, types.BackupFile{
			Name: file.Name(),
			Size: str.FormatBytes(float64(info.Size())),
		})
	}

	return backupList, nil
}

// PostgresqlBackup PostgreSQL备份
func (s *BackupImpl) PostgresqlBackup(database string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
	backupFile := database + "_" + carbon.Now().ToShortDateTimeString() + ".sql"
	if !io.Exists(backupPath) {
		if err := io.Mkdir(backupPath, 0644); err != nil {
			return err
		}
	}

	if _, err := shell.Execf(`su - postgres -c "pg_dump ` + database + `" > ` + backupPath + "/" + backupFile); err != nil {
		return err
	}
	if _, err := shell.Execf("cd " + backupPath + " && zip -r " + backupPath + "/" + backupFile + ".zip " + backupFile); err != nil {
		return err
	}

	return io.Remove(backupPath + "/" + backupFile)
}

// PostgresqlRestore PostgreSQL恢复
func (s *BackupImpl) PostgresqlRestore(database string, backupFile string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
	backupFullPath := filepath.Join(backupPath, backupFile)
	if !io.Exists(backupFullPath) {
		return errors.New("备份文件不存在")
	}

	tempDir, err := io.TempDir(backupFile)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(backupFile, ".sql") {
		backupFile = "" // 置空，防止干扰后续判断
		if err = io.UnArchive(backupFullPath, tempDir); err != nil {
			return err
		}
		if files, err := os.ReadDir(tempDir); err == nil {
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".sql") {
					backupFile = filepath.Base(file.Name())
					break
				}
			}
		}
	} else {
		if err = io.Cp(backupFullPath, filepath.Join(tempDir, backupFile)); err != nil {
			return err
		}
	}

	if len(backupFile) == 0 {
		return errors.New("无法找到备份文件")
	}

	if _, err = shell.Execf(`su - postgres -c "psql ` + database + `" < ` + filepath.Join(tempDir, backupFile)); err != nil {
		return err
	}

	if err = io.Remove(tempDir); err != nil {
		return err
	}

	return nil
}
