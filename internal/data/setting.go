package data

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-rat/utils/hash"
	"github.com/goccy/go-yaml"
	"github.com/gookit/color"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/types"
)

type settingRepo struct{}

func NewSettingRepo() biz.SettingRepo {
	return &settingRepo{}
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
	name := new(biz.Setting)
	if err := app.Orm.Where("key = ?", biz.SettingKeyName).First(name).Error; err != nil {
		return nil, err
	}
	websitePath := new(biz.Setting)
	if err := app.Orm.Where("key = ?", biz.SettingKeyWebsitePath).First(websitePath).Error; err != nil {
		return nil, err
	}
	backupPath := new(biz.Setting)
	if err := app.Orm.Where("key = ?", biz.SettingKeyBackupPath).First(backupPath).Error; err != nil {
		return nil, err
	}

	userID := cast.ToUint(ctx.Value("user_id"))
	user := new(biz.User)
	if err := app.Orm.Where("id = ?", userID).First(user).Error; err != nil {
		return nil, err
	}

	cert, err := io.Read(filepath.Join(app.Root, "panel/storage/cert.pem"))
	if err != nil {
		return nil, err
	}
	key, err := io.Read(filepath.Join(app.Root, "panel/storage/cert.key"))
	if err != nil {
		return nil, err
	}

	return &request.PanelSetting{
		Name:        name.Value,
		Locale:      app.Conf.String("app.locale"),
		Entrance:    app.Conf.String("http.entrance"),
		WebsitePath: websitePath.Value,
		BackupPath:  backupPath.Value,
		Username:    user.Username,
		Email:       user.Email,
		Port:        app.Conf.Int("http.port"),
		HTTPS:       app.Conf.Bool("http.tls"),
		Cert:        cert,
		Key:         key,
	}, nil
}

func (r *settingRepo) UpdatePanelSetting(ctx context.Context, setting *request.PanelSetting) (bool, error) {
	if err := r.Set(biz.SettingKeyName, setting.Name); err != nil {
		return false, err
	}
	if err := r.Set(biz.SettingKeyWebsitePath, setting.WebsitePath); err != nil {
		return false, err
	}
	if err := r.Set(biz.SettingKeyBackupPath, setting.BackupPath); err != nil {
		return false, err
	}

	if err := io.Write(filepath.Join(app.Root, "panel/storage/cert.pem"), setting.Cert, 0644); err != nil {
		return false, err
	}
	if err := io.Write(filepath.Join(app.Root, "panel/storage/cert.key"), setting.Key, 0644); err != nil {
		return false, err
	}

	restartFlag := false
	config := new(types.PanelConfig)
	cm := yaml.CommentMap{}
	raw, err := io.Read("/usr/local/etc/panel/config.yml")
	if err != nil {
		return false, err
	}
	if err = yaml.UnmarshalWithOptions([]byte(raw), config, yaml.CommentToMap(cm)); err != nil {
		return false, err
	}

	config.App.Locale = setting.Locale
	config.HTTP.Port = setting.Port
	config.HTTP.Entrance = setting.Entrance
	config.HTTP.TLS = setting.HTTPS

	encoded, err := yaml.MarshalWithOptions(config, yaml.WithComment(cm))
	if err != nil {
		return false, err
	}
	if err = io.Write("/usr/local/etc/panel/config.yml", string(encoded), 0644); err != nil {
		return false, err
	}
	if raw != string(encoded) {
		restartFlag = true
	}

	user := new(biz.User)
	userID := cast.ToUint(ctx.Value("user_id"))
	if err = app.Orm.Where("id = ?", userID).First(user).Error; err != nil {
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
	if err = app.Orm.Save(user).Error; err != nil {
		return false, err
	}

	return restartFlag, nil
}

func (r *settingRepo) UpdatePanel(version, url, checksum string) error {
	name := filepath.Base(url)
	color.Greenln("目标版本: ", version)
	color.Greenln("下载链接: ", url)
	color.Greenln("文件名: ", name)

	color.Greenln("前置检查...")
	if io.Exists("/tmp/panel-storage.zip") {
		return errors.New("检测到 /tmp 存在临时文件，可能是上次更新失败导致的，请排除后重试")
	}

	color.Greenln("备份面板数据...")
	// 备份面板
	if err := io.Compress([]string{filepath.Join(app.Root, "panel")}, filepath.Join(app.Root, fmt.Sprintf("backup/panel/panel-%s.zip", time.Now().Format("20060102150405"))), io.Zip); err != nil {
		color.Redln("备份面板失败：", err)
		return err
	}
	if err := io.Compress([]string{filepath.Join(app.Root, "panel/storage")}, "/tmp/panel-storage.zip", io.Zip); err != nil {
		color.Redln("备份面板数据失败：", err)
		return err
	}
	if !io.Exists("/tmp/panel-storage.zip") {
		return errors.New("已备份面板数据检查失败")
	}
	color.Greenln("备份完成")

	color.Greenln("清理旧版本...")
	if _, err := shell.Execf("rm -rf %s/panel/*", app.Root); err != nil {
		color.Redln("清理旧版本失败：", err)
		return err
	}
	color.Greenln("清理完成")

	color.Greenln("正在下载...")
	if _, err := shell.Execf("wget -T 120 -t 3 -O %s/panel/%s %s", app.Root, name, url); err != nil {
		color.Redln("下载失败：", err)
		return err
	}
	if _, err := shell.Execf("wget -T 20 -t 3 -O %s/panel/%s %s", app.Root, name+".sha256", checksum); err != nil {
		color.Redln("下载失败：", err)
		return err
	}
	if !io.Exists(filepath.Join(app.Root, "panel", name)) || !io.Exists(filepath.Join(app.Root, "panel", name+".sha256")) {
		return errors.New("下载文件检查失败")
	}
	color.Greenln("下载完成")

	color.Greenln("校验下载文件...")
	check, err := shell.Execf("cd %s/panel && sha256sum -c %s --ignore-missing", app.Root, name+".sha256")
	if check != name+": OK" || err != nil {
		return errors.New("下载文件校验失败")
	}
	if err = io.Remove(filepath.Join(app.Root, "panel", name+".sha256")); err != nil {
		color.Redln("清理校验文件失败：", err)
		return err
	}
	color.Greenln("文件校验完成")

	color.Greenln("更新新版本...")
	if _, err = shell.Execf("cd %s/panel && unzip -o %s && rm -rf %s", app.Root, name, name); err != nil {
		color.Redln("更新失败：", err)
		return err
	}
	if !io.Exists(filepath.Join(app.Root, "panel", "web")) {
		return errors.New("更新失败，可能是下载过程中出现了问题")
	}
	color.Greenln("更新完成")

	color.Greenln("恢复面板数据...")
	if err = io.UnCompress("/tmp/panel-storage.zip", filepath.Join(app.Root, "panel"), io.Zip); err != nil {
		color.Redln("恢复面板数据失败：", err)
		return err
	}
	if !io.Exists(filepath.Join(app.Root, "panel/storage/app.db")) {
		return errors.New("恢复面板数据失败")
	}
	color.Greenln("恢复完成")

	color.Greenln("运行升级后脚本...")
	if _, err = shell.Execf("curl -fsLm 10 https://dl.cdn.haozi.net/panel/auto_update.sh | bash"); err != nil {
		color.Redln("运行面板升级后脚本失败：", err)
		return err
	}
	if _, err = shell.Execf(`wget -O /etc/systemd/system/panel.service https://dl.cdn.haozi.net/panel/panel.service && sed -i "s|/www|%s|g" /etc/systemd/system/panel.service`, app.Root); err != nil {
		color.Redln("下载面板服务文件失败：", err)
		return err
	}
	if _, err = shell.Execf("panel-cli setting write version %s", version); err != nil {
		color.Redln("写入面板版本号失败：", err)
		return err
	}
	if err = io.Mv(filepath.Join(app.Root, "panel/cli"), "/usr/local/sbin/panel-cli"); err != nil {
		color.Redln("移动面板命令行工具失败：", err)
		return err
	}

	color.Greenln("设置面板文件权限...")
	_ = io.Chmod("/usr/local/sbin/panel-cli", 0700)
	_ = io.Chmod("/etc/systemd/system/panel.service", 0700)
	_ = io.Chmod(filepath.Join(app.Root, "panel"), 0700)
	color.Greenln("设置完成")

	color.Greenln("升级完成")

	_, _ = shell.Execf("systemctl daemon-reload")
	_, _ = shell.Execf("systemctl restart panel")
	_ = io.Remove("/tmp/panel-storage.zip")
	_ = io.Remove(filepath.Join(app.Root, "panel/config.example.yml"))

	return nil
}
