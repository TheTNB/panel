package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/samber/do/v2"
	"github.com/shirou/gopsutil/disk"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/pkg/db"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/pkg/types"
)

type backupRepo struct{}

func NewBackupRepo() biz.BackupRepo {
	return do.MustInvoke[biz.BackupRepo](injector)
}

// List 备份列表
func (r *backupRepo) List(typ biz.BackupType) ([]*types.BackupFile, error) {
	path, err := r.GetPath(typ)
	if err != nil {
		return nil, err
	}

	files, err := io.ReadDir(path)
	if err != nil {
		return nil, err
	}

	list := make([]*types.BackupFile, 0)
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		list = append(list, &types.BackupFile{
			Name: file.Name(),
			Path: filepath.Join(path, file.Name()),
			Size: tools.FormatBytes(float64(info.Size())),
			Time: info.ModTime(),
		})
	}

	return list, nil
}

// Create 创建备份
// typ 备份类型
// target 目标名称
// path 可选备份保存路径
func (r *backupRepo) Create(typ biz.BackupType, target string, path ...string) error {
	defPath, err := r.GetPath(typ)
	if err != nil {
		return err
	}
	if len(path) > 0 && path[0] != "" {
		defPath = path[0]
	}

	switch typ {
	case biz.BackupTypeWebsite:
		return r.createWebsite(defPath, target)
	case biz.BackupTypeMySQL:
		return r.createMySQL(defPath, target)
	case biz.BackupTypePostgres:
		return r.createPostgres(defPath, target)
	case biz.BackupTypePanel:
		return r.createPanel(defPath)

	}

	return errors.New("未知备份类型")
}

// Delete 删除备份
func (r *backupRepo) Delete(typ biz.BackupType, name string) error {
	path, err := r.GetPath(typ)
	if err != nil {
		return err
	}

	file := filepath.Join(path, name)
	return io.Remove(file)
}

// Restore 恢复备份
// typ 备份类型
// backup 备份压缩包，可以是绝对路径或者相对路径
// target 目标名称
func (r *backupRepo) Restore(typ biz.BackupType, backup, target string) error {
	if !io.Exists(backup) {
		path, err := r.GetPath(typ)
		if err != nil {
			return err
		}
		backup = filepath.Join(path, backup)
	}

	switch typ {
	case biz.BackupTypeWebsite:
		return r.restoreWebsite(backup, target)
	case biz.BackupTypeMySQL:
		return r.restoreMySQL(backup, target)
	case biz.BackupTypePostgres:
		return r.restorePostgres(backup, target)
	}

	return errors.New("未知备份类型")
}

// CutoffLog 切割日志
// path 保存目录绝对路径
// target 待切割日志文件绝对路径
func (r *backupRepo) CutoffLog(path, target string) error {
	if !io.Exists(target) {
		return errors.New("日志文件不存在")
	}

	to := filepath.Join(path, fmt.Sprintf("%s_%s.zip", time.Now().Format("20060102150405"), filepath.Base(target)))
	if err := io.Compress(filepath.Dir(target), []string{filepath.Base(target)}, to); err != nil {
		return err
	}

	return io.Remove(target)
}

// ClearExpired 清理过期备份
// path 备份目录绝对路径
// prefix 目标文件前缀
// save 保存份数
func (r *backupRepo) ClearExpired(path, prefix string, save int) error {
	files, err := io.ReadDir(path)
	if err != nil {
		return err
	}

	var filtered []os.FileInfo
	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) && strings.HasSuffix(file.Name(), ".zip") {
			info, err := os.Stat(filepath.Join(path, file.Name()))
			if err != nil {
				continue
			}
			filtered = append(filtered, info)
		}
	}

	// 排序所有备份文件，从新到旧
	slices.SortFunc(filtered, func(a, b os.FileInfo) int {
		if a.ModTime().After(b.ModTime()) {
			return -1
		}
		if a.ModTime().Before(b.ModTime()) {
			return 1
		}
		return 0
	})
	if len(filtered) <= save {
		return nil
	}

	// 切片保留 save 份，删除剩余
	toDelete := filtered[save:]
	for _, file := range toDelete {
		filePath := filepath.Join(path, file.Name())
		if app.IsCli {
			fmt.Printf("|-清理过期文件：%s\n", filePath)
		}
		if err = os.Remove(filePath); err != nil {
			if app.IsCli {
				fmt.Printf("|-清理失败：%v\n", err)
			} else {
				return fmt.Errorf("清理失败：%v", err)
			}
		}
	}

	return nil
}

// GetPath 获取备份路径
func (r *backupRepo) GetPath(typ biz.BackupType) (string, error) {
	backupPath, err := NewSettingRepo().Get(biz.SettingKeyBackupPath)
	if err != nil {
		return "", err
	}
	if !slices.Contains([]biz.BackupType{biz.BackupTypePath, biz.BackupTypeWebsite, biz.BackupTypeMySQL, biz.BackupTypePostgres, biz.BackupTypeRedis, biz.BackupTypePanel}, typ) {
		return "", errors.New("未知备份类型")
	}

	backupPath = filepath.Join(backupPath, string(typ))
	if !io.Exists(backupPath) {
		if err = io.Mkdir(backupPath, 0644); err != nil {
			return "", err
		}
	}

	return backupPath, nil
}

// createWebsite 创建网站备份
func (r *backupRepo) createWebsite(to string, name string) error {
	website, err := NewWebsiteRepo().GetByName(name)
	if err != nil {
		return err
	}

	if err = r.preCheckPath(to, website.Path); err != nil {
		return err
	}

	start := time.Now()
	backup := filepath.Join(to, fmt.Sprintf("%s_%s.zip", website.Name, time.Now().Format("20060102150405")))
	if err = io.Compress(website.Path, nil, backup); err != nil {
		return err
	}

	if app.IsCli {
		fmt.Printf("|-备份耗时：%s\n", time.Since(start).String())
		fmt.Printf("|-已备份至文件：%s\n", filepath.Base(backup))
	}
	return nil
}

// createMySQL 创建 MySQL 备份
func (r *backupRepo) createMySQL(to string, name string) error {
	rootPassword, err := NewSettingRepo().Get(biz.SettingKeyMySQLRootPassword)
	if err != nil {
		return err
	}
	mysql, err := db.NewMySQL("root", rootPassword, "/tmp/mysql.sock", "unix")
	if err != nil {
		return err
	}
	defer mysql.Close()
	if exist, _ := mysql.DatabaseExists(name); !exist {
		return fmt.Errorf("数据库不存在：%s", name)
	}
	size, err := mysql.DatabaseSize(name)
	if err != nil {
		return err
	}
	if err = r.preCheckDB(to, size); err != nil {
		return err
	}

	if err = os.Setenv("MYSQL_PWD", rootPassword); err != nil {
		return err
	}
	start := time.Now()
	backup := filepath.Join(to, fmt.Sprintf("%s_%s.sql", name, time.Now().Format("20060102150405")))
	if _, err = shell.Execf(`mysqldump -u root '%s' > '%s'`, name, backup); err != nil {
		return err
	}
	if err = os.Unsetenv("MYSQL_PWD"); err != nil {
		return err
	}

	if err = io.Compress(filepath.Dir(backup), []string{filepath.Base(backup)}, backup+".zip"); err != nil {
		return err
	}
	if err = io.Remove(backup); err != nil {
		return err
	}

	if app.IsCli {
		fmt.Printf("|-备份耗时：%s\n", time.Since(start).String())
		fmt.Printf("|-已备份至文件：%s\n", filepath.Base(backup+".zip"))
	}
	return nil
}

// createPostgres 创建 PostgreSQL 备份
func (r *backupRepo) createPostgres(to string, name string) error {
	postgres, err := db.NewPostgres("postgres", "", "127.0.0.1", 5432)
	if err != nil {
		return err
	}
	defer postgres.Close()
	if exist, _ := postgres.DatabaseExist(name); !exist {
		return fmt.Errorf("数据库不存在：%s", name)
	}
	size, err := postgres.DatabaseSize(name)
	if err != nil {
		return err
	}
	if err = r.preCheckDB(to, size); err != nil {
		return err
	}

	start := time.Now()
	backup := filepath.Join(to, fmt.Sprintf("%s_%s.sql", name, time.Now().Format("20060102150405")))
	if _, err = shell.Execf(`su - postgres -c "pg_dump '%s'" > '%s'`, name, backup); err != nil {
		return err
	}

	if err = io.Compress(filepath.Dir(backup), []string{filepath.Base(backup)}, backup+".zip"); err != nil {
		return err
	}
	if err = io.Remove(backup); err != nil {
		return err
	}

	if app.IsCli {
		fmt.Printf("|-备份耗时：%s\n", time.Since(start).String())
		fmt.Printf("|-已备份至文件：%s\n", filepath.Base(backup+".zip"))
	}
	return nil
}

// createPanel 创建面板备份
func (r *backupRepo) createPanel(to string) error {
	backup := filepath.Join(to, fmt.Sprintf("panel_%s.zip", time.Now().Format("20060102150405")))

	if err := r.preCheckPath(to, filepath.Join(app.Root, "panel")); err != nil {
		return err
	}

	start := time.Now()

	temp, err := io.TempDir("panel-backup")
	if err != nil {
		return err
	}

	if err = io.Cp(filepath.Join(app.Root, "panel"), temp); err != nil {
		return err
	}
	if err = io.Cp("/usr/local/sbin/panel-cli", temp); err != nil {
		return err
	}
	if err = io.Cp("/usr/local/etc/panel/config.yml", temp); err != nil {
		return err
	}

	_ = io.Chmod(temp, 0600)
	if err = io.Compress(temp, nil, backup); err != nil {
		return err
	}
	if err = io.Chmod(backup, 0600); err != nil {
		return err
	}

	if app.IsCli {
		fmt.Printf("|-备份耗时：%s\n", time.Since(start).String())
		fmt.Printf("|-已备份至文件：%s\n", filepath.Base(backup))
	}

	return io.Remove(temp)
}

// restoreWebsite 恢复网站备份
func (r *backupRepo) restoreWebsite(backup, target string) error {
	if !io.Exists(backup) {
		return errors.New("备份文件不存在")
	}

	website, err := NewWebsiteRepo().GetByName(target)
	if err != nil {
		return err
	}

	if err = io.Remove(website.Path); err != nil {
		return err
	}
	if err = io.UnCompress(backup, website.Path); err != nil {
		return err
	}
	if err = io.Chmod(website.Path, 0755); err != nil {
		return err
	}
	if err = io.Chown(website.Path, "www", "www"); err != nil {
		return err
	}

	return nil
}

// restoreMySQL 恢复 MySQL 备份
func (r *backupRepo) restoreMySQL(backup, target string) error {
	if !io.Exists(backup) {
		return errors.New("备份文件不存在")
	}

	rootPassword, err := NewSettingRepo().Get(biz.SettingKeyMySQLRootPassword)
	if err != nil {
		return err
	}
	mysql, err := db.NewMySQL("root", rootPassword, "/tmp/mysql.sock", "unix")
	if err != nil {
		return err
	}
	defer mysql.Close()
	if exist, _ := mysql.DatabaseExists(target); !exist {
		return fmt.Errorf("数据库不存在：%s", target)
	}
	if err = os.Setenv("MYSQL_PWD", rootPassword); err != nil {
		return err
	}

	if !strings.HasSuffix(backup, ".sql") {
		backup, err = r.autoUnCompressSQL(backup)
		if err != nil {
			return err
		}
	}

	if _, err = shell.Execf(`mysql -u root '%s' < '%s'`, target, backup); err != nil {
		return err
	}
	if err = os.Unsetenv("MYSQL_PWD"); err != nil {
		return err
	}

	_ = io.Remove(filepath.Dir(backup))
	return nil
}

// restorePostgres 恢复 PostgreSQL 备份
func (r *backupRepo) restorePostgres(backup, target string) error {
	if !io.Exists(backup) {
		return errors.New("备份文件不存在")
	}

	postgres, err := db.NewPostgres("postgres", "", "127.0.0.1", 5432)
	if err != nil {
		return err
	}
	defer postgres.Close()
	if exist, _ := postgres.DatabaseExist(target); !exist {
		return fmt.Errorf("数据库不存在：%s", target)
	}

	if !strings.HasSuffix(backup, ".sql") {
		backup, err = r.autoUnCompressSQL(backup)
		if err != nil {
			return err
		}
	}

	if _, err = shell.Execf(`su - postgres -c "psql '%s'" < '%s'`, target, backup); err != nil {
		return err
	}

	_ = io.Remove(filepath.Dir(backup))
	return nil
}

// preCheckPath 预检空间和 inode 是否足够
// to 备份保存目录
// path 待备份目录
func (r *backupRepo) preCheckPath(to, path string) error {
	size, err := io.SizeX(path)
	if err != nil {
		return err
	}
	files, err := io.CountX(path)
	if err != nil {
		return err
	}

	usage, err := disk.Usage(to)
	if err != nil {
		return err
	}

	if app.IsCli {
		fmt.Printf("|-目标大小：%s\n", tools.FormatBytes(float64(size)))
		fmt.Printf("|-目标文件数：%d\n", files)
		fmt.Printf("|-备份目录可用空间：%s\n", tools.FormatBytes(float64(usage.Free)))
		fmt.Printf("|-备份目录可用Inode：%d\n", usage.InodesFree)
	}

	if uint64(size) > usage.Free {
		return errors.New("备份目录空间不足")
	}
	if uint64(files) > usage.InodesFree {
		return errors.New("备份目录Inode不足")
	}

	return nil
}

// preCheckDB 预检空间和 inode 是否足够
// to 备份保存目录
// size 数据库大小
func (r *backupRepo) preCheckDB(to string, size int64) error {
	usage, err := disk.Usage(to)
	if err != nil {
		return err
	}

	if app.IsCli {
		fmt.Printf("|-目标大小：%s\n", tools.FormatBytes(float64(size)))
		fmt.Printf("|-备份目录可用空间：%s\n", tools.FormatBytes(float64(usage.Free)))
		fmt.Printf("|-备份目录可用Inode：%d\n", usage.InodesFree)
	}

	if uint64(size) > usage.Free {
		return errors.New("备份目录空间不足")
	}

	return nil
}

// autoUnCompressSQL 自动处理压缩文件
func (r *backupRepo) autoUnCompressSQL(backup string) (string, error) {
	temp, err := io.TempDir(backup)
	if err != nil {
		return "", err
	}

	if err = io.UnCompress(backup, temp); err != nil {
		return "", err
	}

	backup = "" // 置空，防止干扰后续判断
	if files, err := os.ReadDir(temp); err == nil {
		if len(files) != 1 {
			return "", fmt.Errorf("压缩文件中包含的文件数量不为1，实际为%d", len(files))
		}
		if strings.HasSuffix(files[0].Name(), ".sql") {
			backup = filepath.Join(temp, files[0].Name())
		}
	}

	if backup == "" {
		return "", errors.New("无法找到.sql备份文件")
	}

	return backup, nil
}
