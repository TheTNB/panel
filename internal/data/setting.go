package data

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/go-rat/utils/hash"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/cert"
	"github.com/TheTNB/panel/pkg/firewall"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/pkg/types"
)

type settingRepo struct {
	taskRepo biz.TaskRepo
}

func NewSettingRepo() biz.SettingRepo {
	return &settingRepo{
		taskRepo: NewTaskRepo(),
	}
}

func (r *settingRepo) Get(key biz.SettingKey, defaultValue ...string) (string, error) {
	setting := new(biz.Setting)
	if err := app.Orm.Where("key = ?", key).First(setting).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
	}

	if setting.Value == "" && len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return setting.Value, nil
}

func (r *settingRepo) GetBool(key biz.SettingKey, defaultValue ...bool) (bool, error) {
	setting := new(biz.Setting)
	if err := app.Orm.Where("key = ?", key).First(setting).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
	}

	if setting.Value == "" && len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return cast.ToBool(setting.Value), nil
}

func (r *settingRepo) Set(key biz.SettingKey, value string) error {
	setting := new(biz.Setting)
	if err := app.Orm.Where("key = ?", key).First(setting).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	setting.Key = key
	setting.Value = value
	return app.Orm.Save(setting).Error
}

func (r *settingRepo) Delete(key biz.SettingKey) error {
	setting := new(biz.Setting)
	if err := app.Orm.Where("key = ?", key).Delete(setting).Error; err != nil {
		return err
	}

	return nil
}

func (r *settingRepo) GetPanelSetting(ctx context.Context) (*request.PanelSetting, error) {
	name, err := r.Get(biz.SettingKeyName)
	if err != nil {
		return nil, err
	}
	offlineMode, err := r.Get(biz.SettingKeyOfflineMode)
	if err != nil {
		return nil, err
	}
	websitePath, err := r.Get(biz.SettingKeyWebsitePath)
	if err != nil {
		return nil, err
	}
	backupPath, err := r.Get(biz.SettingKeyBackupPath)
	if err != nil {
		return nil, err
	}

	userID := cast.ToUint(ctx.Value("user_id"))
	user := new(biz.User)
	if err := app.Orm.Where("id = ?", userID).First(user).Error; err != nil {
		return nil, err
	}

	crt, err := io.Read(filepath.Join(app.Root, "panel/storage/cert.pem"))
	if err != nil {
		return nil, err
	}
	key, err := io.Read(filepath.Join(app.Root, "panel/storage/cert.key"))
	if err != nil {
		return nil, err
	}

	return &request.PanelSetting{
		Name:        name,
		Locale:      app.Conf.String("app.locale"),
		Entrance:    app.Conf.String("http.entrance"),
		OfflineMode: cast.ToBool(offlineMode),
		WebsitePath: websitePath,
		BackupPath:  backupPath,
		Username:    user.Username,
		Email:       user.Email,
		Port:        app.Conf.Int("http.port"),
		HTTPS:       app.Conf.Bool("http.tls"),
		Cert:        crt,
		Key:         key,
	}, nil
}

func (r *settingRepo) UpdatePanelSetting(ctx context.Context, setting *request.PanelSetting) (bool, error) {
	if err := r.Set(biz.SettingKeyName, setting.Name); err != nil {
		return false, err
	}
	if err := r.Set(biz.SettingKeyOfflineMode, cast.ToString(setting.OfflineMode)); err != nil {
		return false, err
	}
	if err := r.Set(biz.SettingKeyWebsitePath, setting.WebsitePath); err != nil {
		return false, err
	}
	if err := r.Set(biz.SettingKeyBackupPath, setting.BackupPath); err != nil {
		return false, err
	}

	// 用户
	user := new(biz.User)
	userID := cast.ToUint(ctx.Value("user_id"))
	if err := app.Orm.Where("id = ?", userID).First(user).Error; err != nil {
		return false, err
	}

	user.Username = setting.Username
	user.Email = setting.Email
	if setting.Password != "" {
		value, err := hash.NewArgon2id().Make(setting.Password)
		if err != nil {
			return false, err
		}
		user.Password = value
	}
	if err := app.Orm.Save(user).Error; err != nil {
		return false, err
	}

	// 下面是需要需要重启的设置
	// 面板HTTPS
	restartFlag := false
	oldCert, _ := io.Read(filepath.Join(app.Root, "panel/storage/cert.pem"))
	oldKey, _ := io.Read(filepath.Join(app.Root, "panel/storage/cert.key"))
	if oldCert != setting.Cert || oldKey != setting.Key {
		if r.taskRepo.HasRunningTask() {
			return false, errors.New("后台任务正在运行，禁止修改部分设置，请稍后再试")
		}
		restartFlag = true
	}
	if _, err := cert.ParseCert(setting.Cert); err != nil {
		return false, fmt.Errorf("failed to parse certificate: %w", err)
	}
	if _, err := cert.ParseKey(setting.Key); err != nil {
		return false, fmt.Errorf("failed to parse private key: %w", err)
	}
	if err := io.Write(filepath.Join(app.Root, "panel/storage/cert.pem"), setting.Cert, 0644); err != nil {
		return false, err
	}
	if err := io.Write(filepath.Join(app.Root, "panel/storage/cert.key"), setting.Key, 0644); err != nil {
		return false, err
	}

	// 面板主配置
	config := new(types.PanelConfig)
	raw, err := io.Read("/usr/local/etc/panel/config.yml")
	if err != nil {
		return false, err
	}
	if err = yaml.Unmarshal([]byte(raw), config); err != nil {
		return false, err
	}

	config.App.Locale = setting.Locale
	config.HTTP.Port = setting.Port
	config.HTTP.Entrance = setting.Entrance
	config.HTTP.TLS = setting.HTTPS

	// 放行端口
	fw := firewall.NewFirewall()
	err = fw.Port(firewall.FireInfo{
		Type:      firewall.TypeNormal,
		PortStart: uint(config.HTTP.Port),
		PortEnd:   uint(config.HTTP.Port),
		Direction: firewall.DirectionIn,
		Strategy:  firewall.StrategyAccept,
	}, firewall.OperationAdd)
	if err != nil {
		return false, err
	}

	encoded, err := yaml.Marshal(config)
	if err != nil {
		return false, err
	}
	if raw != string(encoded) {
		if r.taskRepo.HasRunningTask() {
			return false, errors.New("后台任务正在运行，禁止修改部分设置，请稍后再试")
		}
		restartFlag = true
	}
	if err = io.Write("/usr/local/etc/panel/config.yml", string(encoded), 0644); err != nil {
		return false, err
	}

	return restartFlag, nil
}

func (r *settingRepo) UpdatePanel(version, url, checksum string) error {
	// 预先优化数据库
	if err := app.Orm.Exec("VACUUM").Error; err != nil {
		return err
	}
	if err := app.Orm.Exec("PRAGMA wal_checkpoint(TRUNCATE);").Error; err != nil {
		return err
	}

	name := filepath.Base(url)
	if app.IsCli {
		fmt.Printf("|-目标版本：%s\n", version)
		fmt.Printf("|-下载链接：%s\n", url)
		fmt.Printf("|-文件名：%s\n", name)
	}

	if app.IsCli {
		fmt.Println("|-正在下载...")
	}
	if _, err := shell.Execf("wget -T 120 -t 3 -O /tmp/%s %s", name, url); err != nil {
		return fmt.Errorf("下载失败：%w", err)
	}
	if _, err := shell.Execf("wget -T 20 -t 3 -O /tmp/%s %s", name+".sha256", checksum); err != nil {
		return fmt.Errorf("下载失败：%w", err)
	}
	if !io.Exists(filepath.Join("/tmp", name)) || !io.Exists(filepath.Join("/tmp", name+".sha256")) {
		return errors.New("下载文件检查失败")
	}

	if app.IsCli {
		fmt.Println("|-校验下载文件...")
	}
	if check, err := shell.Execf("cd /tmp && sha256sum -c %s --ignore-missing", name+".sha256"); check != name+": OK" || err != nil {
		return errors.New("下载文件校验失败")
	}
	if err := io.Remove(filepath.Join("/tmp", name+".sha256")); err != nil {
		if app.IsCli {
			fmt.Println("|-清理校验文件失败：", err)
		}
		return fmt.Errorf("清理校验文件失败：%w", err)
	}

	if app.IsCli {
		fmt.Println("|-前置检查...")
	}
	if io.Exists("/tmp/panel-storage.zip") {
		return errors.New("检测到 /tmp 存在临时文件，可能是上次更新失败所致，请运行 panel-cli fix 修复后重试")
	}

	if app.IsCli {
		fmt.Println("|-备份面板数据...")
	}
	// 备份面板
	backup := NewBackupRepo()
	if err := backup.Create(biz.BackupTypePanel, ""); err != nil {
		if app.IsCli {
			fmt.Println("|-备份面板失败：", err)
		}
		return fmt.Errorf("备份面板失败：%w", err)
	}
	if err := io.Compress(filepath.Join(app.Root, "panel/storage"), nil, "/tmp/panel-storage.zip"); err != nil {
		if app.IsCli {
			fmt.Println("|-备份面板数据失败：", err)
		}
		return fmt.Errorf("备份面板数据失败：%w", err)
	}
	if !io.Exists("/tmp/panel-storage.zip") {
		return errors.New("已备份面板数据检查失败")
	}

	if app.IsCli {
		fmt.Println("|-清理旧版本...")
	}
	if _, err := shell.Execf("rm -rf %s/panel/*", app.Root); err != nil {
		return fmt.Errorf("清理旧版本失败：%w", err)
	}

	if app.IsCli {
		fmt.Println("|-解压新版本...")
	}
	if err := io.UnCompress(filepath.Join("/tmp", name), filepath.Join(app.Root, "panel")); err != nil {
		return fmt.Errorf("解压失败：%w", err)
	}
	if !io.Exists(filepath.Join(app.Root, "panel", "web")) {
		return errors.New("解压失败，缺失文件")
	}

	if app.IsCli {
		fmt.Println("|-恢复面板数据...")
	}
	if err := io.UnCompress("/tmp/panel-storage.zip", filepath.Join(app.Root, "panel")); err != nil {
		return fmt.Errorf("恢复面板数据失败：%w", err)
	}
	if !io.Exists(filepath.Join(app.Root, "panel/storage/app.db")) {
		return errors.New("恢复面板数据失败")
	}

	if app.IsCli {
		fmt.Println("|-运行更新后脚本...")
	}
	if _, err := shell.Execf("curl -fsLm 10 https://dl.cdn.haozi.net/panel/auto_update.sh | bash"); err != nil {
		return fmt.Errorf("运行面板更新后脚本失败：%w", err)
	}
	if _, err := shell.Execf(`wget -O /etc/systemd/system/panel.service https://dl.cdn.haozi.net/panel/panel.service && sed -i "s|/www|%s|g" /etc/systemd/system/panel.service`, app.Root); err != nil {
		return fmt.Errorf("下载面板服务文件失败：%w", err)
	}
	if _, err := shell.Execf("panel-cli setting write version %s", version); err != nil {
		return fmt.Errorf("写入面板版本号失败：%w", err)
	}
	if err := io.Mv(filepath.Join(app.Root, "panel/cli"), "/usr/local/sbin/panel-cli"); err != nil {
		return fmt.Errorf("移动面板命令行工具失败：%w", err)
	}

	if app.IsCli {
		fmt.Println("|-设置关键文件权限...")
	}
	_ = io.Chmod("/usr/local/sbin/panel-cli", 0700)
	_ = io.Chmod("/etc/systemd/system/panel.service", 0700)
	_ = io.Chmod(filepath.Join(app.Root, "panel"), 0700)

	if app.IsCli {
		fmt.Println("|-更新完成")
	}

	_, _ = shell.Execf("systemctl daemon-reload")
	_ = io.Remove("/tmp/panel-storage.zip")
	_ = io.Remove(filepath.Join(app.Root, "panel/config.example.yml"))
	tools.RestartPanel()

	return nil
}

func (r *settingRepo) FixPanel() error {
	if app.IsCli {
		fmt.Println("|-开始修复面板...")
	}
	// 检查关键文件是否正常
	flag := false
	if !io.Exists(filepath.Join(app.Root, "panel", "web")) {
		flag = true
	}
	if !io.Exists(filepath.Join(app.Root, "panel", "storage", "app.db")) {
		flag = true
	}
	if io.Exists("/tmp/panel-storage.zip") {
		flag = true
	}
	if !flag {
		return fmt.Errorf("文件正常无需修复，请运行 panel-cli update 更新面板")
	}

	// 再次确认是否需要修复
	if io.Exists("/tmp/panel-storage.zip") {
		// 文件齐全情况下只移除临时文件
		if io.Exists(filepath.Join(app.Root, "panel", "web")) &&
			io.Exists(filepath.Join(app.Root, "panel", "storage", "app.db")) &&
			io.Exists("/usr/local/etc/panel/config.yml") {
			if err := io.Remove("/tmp/panel-storage.zip"); err != nil {
				return fmt.Errorf("清理临时文件失败：%w", err)
			}
			if app.IsCli {
				fmt.Println("已清理临时文件，请运行 panel-cli update 更新面板")
			}
			return nil
		}
	}

	// 从备份目录中找最新的备份文件
	backup := NewBackupRepo()
	list, err := backup.List(biz.BackupTypePanel)
	if err != nil {
		return err
	}
	slices.SortFunc(list, func(a *types.BackupFile, b *types.BackupFile) int {
		return int(b.Time.Unix() - a.Time.Unix())
	})
	if len(list) == 0 {
		return fmt.Errorf("未找到备份文件，无法自动修复")
	}
	latest := list[0]
	if app.IsCli {
		fmt.Printf("|-使用备份文件：%s\n", latest.Name)
	}

	// 解压备份文件
	if app.IsCli {
		fmt.Println("|-解压备份文件...")
	}
	if err = io.Remove("/tmp/panel-fix"); err != nil {
		return fmt.Errorf("清理临时目录失败：%w", err)
	}
	if err = io.UnCompress(latest.Path, "/tmp/panel-fix"); err != nil {
		return fmt.Errorf("解压备份文件失败：%w", err)
	}

	// 移动文件到对应位置
	if app.IsCli {
		fmt.Println("|-移动备份文件...")
	}
	if io.Exists(filepath.Join("/tmp/panel-fix", "panel")) && io.IsDir(filepath.Join("/tmp/panel-fix", "panel")) {
		if err = io.Remove(filepath.Join(app.Root, "panel")); err != nil {
			return fmt.Errorf("删除目录失败：%w", err)
		}
		if err = io.Mv(filepath.Join("/tmp/panel-fix", "panel"), filepath.Join(app.Root)); err != nil {
			return fmt.Errorf("移动目录失败：%w", err)
		}
	}
	if io.Exists(filepath.Join("/tmp/panel-fix", "config.yml")) {
		if err = io.Mv(filepath.Join("/tmp/panel-fix", "config.yml"), "/usr/local/etc/panel/config.yml"); err != nil {
			return fmt.Errorf("移动文件失败：%w", err)
		}
	}
	if io.Exists(filepath.Join("/tmp/panel-fix", "panel-cli")) {
		if err = io.Mv(filepath.Join("/tmp/panel-fix", "panel-cli"), "/usr/local/sbin/panel-cli"); err != nil {
			return fmt.Errorf("移动文件失败：%w", err)
		}
	}

	// tmp目录下如果有storage备份，则解压回去
	if app.IsCli {
		fmt.Println("|-恢复面板数据...")
	}
	if io.Exists("/tmp/panel-storage.zip") {
		if err = io.UnCompress("/tmp/panel-storage.zip", filepath.Join(app.Root, "panel")); err != nil {
			return fmt.Errorf("恢复面板数据失败：%w", err)
		}
		if err = io.Remove("/tmp/panel-storage.zip"); err != nil {
			return fmt.Errorf("清理临时文件失败：%w", err)
		}
	}

	// 下载服务文件
	if !io.Exists("/etc/systemd/system/panel.service") {
		if _, err = shell.Execf(`wget -O /etc/systemd/system/panel.service https://dl.cdn.haozi.net/panel/panel.service && sed -i "s|/www|%s|g" /etc/systemd/system/panel.service`, app.Root); err != nil {
			return err
		}
	}

	// 处理权限
	if app.IsCli {
		fmt.Println("|-设置关键文件权限...")
	}
	if err = io.Chmod("/usr/local/etc/panel/config.yml", 0600); err != nil {
		return err
	}
	if err = io.Chmod("/etc/systemd/system/panel.service", 0700); err != nil {
		return err
	}
	if err = io.Chmod("/usr/local/sbin/panel-cli", 0700); err != nil {
		return err
	}
	if err = io.Chmod(filepath.Join(app.Root, "panel"), 0700); err != nil {
		return err
	}

	if app.IsCli {
		fmt.Println("|-修复完成")
	}

	tools.RestartPanel()
	return nil
}
