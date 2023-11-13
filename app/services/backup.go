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
	PostgresqlList() ([]BackupFile, error)
	PostgresqlBackup(database string) error
	PostgresqlRestore(database string, backupFile string) error
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
		if err := tools.Mkdir(backupPath, 0644); err != nil {
			return []BackupFile{}, err
		}
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

	if _, err := tools.Exec(`rm -rf '` + website.Path + `/*'`); err != nil {
		return err
	}
	if _, err := tools.Exec(`unzip -o '` + backupFile + `' -d '` + website.Path + `' 2>&1`); err != nil {
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
func (s *BackupImpl) MysqlList() ([]BackupFile, error) {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return []BackupFile{}, nil
	}

	backupPath += "/mysql"
	if !tools.Exists(backupPath) {
		if err := tools.Mkdir(backupPath, 0644); err != nil {
			return []BackupFile{}, err
		}
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
	tools.Remove(backupPath + "/" + backupFile)

	return os.Unsetenv("MYSQL_PWD")
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
		if _, err := tools.Exec("unzip -o " + backupFile + " -d " + backupPath); err != nil {
			return err
		}
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".gz":
		if strings.HasSuffix(backupFile, ".tar.gz") {
			// 解压.tar.gz文件
			if _, err := tools.Exec("tar -zxvf " + backupFile + " -C " + backupPath); err != nil {
				return err
			}
			backupFile = strings.TrimSuffix(backupFile, ".tar.gz")
		} else {
			// 解压.gz文件
			if _, err := tools.Exec("gzip -d " + backupFile); err != nil {
				return err
			}
			backupFile = strings.TrimSuffix(backupFile, ext)
		}
	case ".bz2":
		if _, err := tools.Exec("bzip2 -d " + backupFile); err != nil {
			return err
		}
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".tar":
		if _, err := tools.Exec("tar -xvf " + backupFile + " -C " + backupPath); err != nil {
			return err
		}
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".rar":
		if _, err := tools.Exec("unrar x " + backupFile + " " + backupPath); err != nil {
			return err
		}
		backupFile = strings.TrimSuffix(backupFile, ext)
	}

	if !tools.Exists(backupFile) {
		return errors.New("自动解压失败，请手动解压")
	}

	if _, err := tools.Exec("/www/server/mysql/bin/mysql -uroot " + database + " < " + backupFile); err != nil {
		return err
	}

	return os.Unsetenv("MYSQL_PWD")
}

// PostgresqlList PostgreSQL备份列表
func (s *BackupImpl) PostgresqlList() ([]BackupFile, error) {
	backupPath := s.setting.Get(models.SettingKeyBackupPath)
	if len(backupPath) == 0 {
		return []BackupFile{}, nil
	}

	backupPath += "/postgresql"
	if !tools.Exists(backupPath) {
		if err := tools.Mkdir(backupPath, 0644); err != nil {
			return []BackupFile{}, err
		}
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

	tools.Remove(backupPath + "/" + backupFile)
	return nil
}

// PostgresqlRestore PostgreSQL恢复
func (s *BackupImpl) PostgresqlRestore(database string, backupFile string) error {
	backupPath := s.setting.Get(models.SettingKeyBackupPath) + "/postgresql"
	ext := filepath.Ext(backupFile)
	backupFile = backupPath + "/" + backupFile
	if !tools.Exists(backupFile) {
		return errors.New("备份文件不存在")
	}

	switch ext {
	case ".zip":
		if _, err := tools.Exec("unzip -o " + backupFile + " -d " + backupPath); err != nil {
			return err
		}
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".gz":
		if strings.HasSuffix(backupFile, ".tar.gz") {
			// 解压.tar.gz文件
			if _, err := tools.Exec("tar -zxvf " + backupFile + " -C " + backupPath); err != nil {
				return err
			}
			backupFile = strings.TrimSuffix(backupFile, ".tar.gz")
		} else {
			// 解压.gz文件
			if _, err := tools.Exec("gzip -d " + backupFile); err != nil {
				return err
			}
			backupFile = strings.TrimSuffix(backupFile, ext)
		}
	case ".bz2":
		if _, err := tools.Exec("bzip2 -d " + backupFile); err != nil {
			return err
		}
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".tar":
		if _, err := tools.Exec("tar -xvf " + backupFile + " -C " + backupPath); err != nil {
			return err
		}
		backupFile = strings.TrimSuffix(backupFile, ext)
	case ".rar":
		if _, err := tools.Exec("unrar x " + backupFile + " " + backupPath); err != nil {
			return err
		}
		backupFile = strings.TrimSuffix(backupFile, ext)
	}

	if !tools.Exists(backupFile) {
		return errors.New("自动解压失败，请手动解压")
	}

	if _, err := tools.Exec(`su - postgres -c "psql ` + database + `" < ` + backupFile); err != nil {
		return err
	}

	return nil
}
