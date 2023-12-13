// Package services 备份服务
package services

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/support/carbon"

	"panel/app/internal"
	"panel/app/models"
	"panel/pkg/tools"
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
func (s *BackupImpl) WebsiteList() ([]internal.BackupFile, error) {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return []internal.BackupFile{}, nil
	}

	backupPath += "/website"
	if !tools.Exists(backupPath) {
		if err := tools.Mkdir(backupPath, 0644); err != nil {
			return []internal.BackupFile{}, err
		}
	}

	files, err := os.ReadDir(backupPath)
	if err != nil {
		return []internal.BackupFile{}, err
	}
	var backupList []internal.BackupFile
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		backupList = append(backupList, internal.BackupFile{
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
		if err := tools.Mkdir(backupPath, 0644); err != nil {
			return err
		}
	}

	backupFile := backupPath + "/" + website.Name + "_" + carbon.Now().ToShortDateTimeString() + ".zip"
	if _, err := tools.Exec(`cd '` + website.Path + `' && zip -r '` + backupFile + `' .`); err != nil {
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
	if !tools.Exists(backupPath) {
		if err := tools.Mkdir(backupPath, 0644); err != nil {
			return err
		}
	}

	backupFile = backupPath + "/" + backupFile
	if !tools.Exists(backupFile) {
		return errors.New("备份文件不存在")
	}

	if err := tools.Remove(website.Path); err != nil {
		return err
	}
	if err := tools.UnArchive(backupFile, website.Path); err != nil {
		return err
	}
	if err := tools.Chmod(website.Path, 0755); err != nil {
		return err
	}
	if err := tools.Chown(website.Path, "www", "www"); err != nil {
		return err
	}

	return nil
}

// MysqlList MySQL备份列表
func (s *BackupImpl) MysqlList() ([]internal.BackupFile, error) {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return []internal.BackupFile{}, nil
	}

	backupPath += "/mysql"
	if !tools.Exists(backupPath) {
		if err := tools.Mkdir(backupPath, 0644); err != nil {
			return []internal.BackupFile{}, err
		}
	}

	files, err := os.ReadDir(backupPath)
	if err != nil {
		return []internal.BackupFile{}, err
	}
	var backupList []internal.BackupFile
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		backupList = append(backupList, internal.BackupFile{
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
		if err := tools.Mkdir(backupPath, 0644); err != nil {
			return err
		}
	}
	err := os.Setenv("MYSQL_PWD", rootPassword)
	if err != nil {
		return err
	}

	if _, err := tools.Exec("/www/server/mysql/bin/mysqldump -uroot " + database + " > " + backupPath + "/" + backupFile); err != nil {
		return err
	}
	if _, err := tools.Exec("cd " + backupPath + " && zip -r " + backupPath + "/" + backupFile + ".zip " + backupFile); err != nil {
		return err
	}
	if err := tools.Remove(backupPath + "/" + backupFile); err != nil {
		return err
	}

	return os.Unsetenv("MYSQL_PWD")
}

// MysqlRestore MySQL恢复
func (s *BackupImpl) MysqlRestore(database string, backupFile string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath) + "/mysql"
	rootPassword := s.setting.Get(models.SettingKeyMysqlRootPassword)
	backupFullPath := filepath.Join(backupPath, backupFile)
	if !tools.Exists(backupFullPath) {
		return errors.New("备份文件不存在")
	}

	if err := os.Setenv("MYSQL_PWD", rootPassword); err != nil {
		return err
	}

	tempDir, err := tools.TempDir(backupFile)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(backupFile, ".sql") {
		backupFile = "" // 置空，防止干扰后续判断
		if err = tools.UnArchive(backupFullPath, tempDir); err != nil {
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
		if err = tools.Cp(backupFullPath, filepath.Join(tempDir, backupFile)); err != nil {
			return err
		}
	}

	if len(backupFile) == 0 {
		return errors.New("无法找到备份文件")
	}

	if _, err = tools.Exec("/www/server/mysql/bin/mysql -uroot " + database + " < " + filepath.Join(tempDir, backupFile)); err != nil {
		return err
	}

	if err = tools.Remove(tempDir); err != nil {
		return err
	}

	return os.Unsetenv("MYSQL_PWD")
}

// PostgresqlList PostgreSQL备份列表
func (s *BackupImpl) PostgresqlList() ([]internal.BackupFile, error) {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return []internal.BackupFile{}, nil
	}

	backupPath += "/postgresql"
	if !tools.Exists(backupPath) {
		if err := tools.Mkdir(backupPath, 0644); err != nil {
			return []internal.BackupFile{}, err
		}
	}

	files, err := os.ReadDir(backupPath)
	if err != nil {
		return []internal.BackupFile{}, err
	}
	var backupList []internal.BackupFile
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		backupList = append(backupList, internal.BackupFile{
			Name: file.Name(),
			Size: tools.FormatBytes(float64(info.Size())),
		})
	}

	return backupList, nil
}

// PostgresqlBackup PostgreSQL备份
func (s *BackupImpl) PostgresqlBackup(database string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
	backupFile := database + "_" + carbon.Now().ToShortDateTimeString() + ".sql"
	if !tools.Exists(backupPath) {
		if err := tools.Mkdir(backupPath, 0644); err != nil {
			return err
		}
	}

	if _, err := tools.Exec(`su - postgres -c "pg_dump ` + database + `" > ` + backupPath + "/" + backupFile); err != nil {
		return err
	}
	if _, err := tools.Exec("cd " + backupPath + " && zip -r " + backupPath + "/" + backupFile + ".zip " + backupFile); err != nil {
		return err
	}

	return tools.Remove(backupPath + "/" + backupFile)
}

// PostgresqlRestore PostgreSQL恢复
func (s *BackupImpl) PostgresqlRestore(database string, backupFile string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
	backupFullPath := filepath.Join(backupPath, backupFile)
	if !tools.Exists(backupFullPath) {
		return errors.New("备份文件不存在")
	}

	tempDir, err := tools.TempDir(backupFile)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(backupFile, ".sql") {
		backupFile = "" // 置空，防止干扰后续判断
		if err = tools.UnArchive(backupFullPath, tempDir); err != nil {
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
		if err = tools.Cp(backupFullPath, filepath.Join(tempDir, backupFile)); err != nil {
			return err
		}
	}

	if len(backupFile) == 0 {
		return errors.New("无法找到备份文件")
	}

	if _, err = tools.Exec(`su - postgres -c "psql ` + database + `" < ` + filepath.Join(tempDir, backupFile)); err != nil {
		return err
	}

	if err = tools.Remove(tempDir); err != nil {
		return err
	}

	return nil
}
