package db

import (
	"errors"
	"fmt"

	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
)

// MySQLResetRootPassword 重置 MySQL root密码
func MySQLResetRootPassword(password string) error {
	_ = systemctl.Stop("mysqld")
	if run, err := systemctl.Status("mysqld"); err != nil || run {
		return fmt.Errorf("停止MySQL失败: %w", err)
	}
	_, _ = shell.Execf(`systemctl set-environment MYSQLD_OPTS="--skip-grant-tables --skip-networking"`)
	if err := systemctl.Start("mysqld"); err != nil {
		return fmt.Errorf("以安全模式启动MySQL失败: %w", err)
	}
	if _, err := shell.Execf(`mysql -uroot -e "FLUSH PRIVILEGES;UPDATE mysql.user SET authentication_string=null WHERE user='root' AND host='localhost';ALTER USER 'root'@'localhost' IDENTIFIED BY '%s';FLUSH PRIVILEGES;"`, password); err != nil {
		return errors.New("设置root密码失败")
	}
	if err := systemctl.Stop("mysqld"); err != nil {
		return fmt.Errorf("停止MySQL失败: %w", err)
	}
	_, _ = shell.Execf(`systemctl unset-environment MYSQLD_OPTS`)
	if err := systemctl.Start("mysqld"); err != nil {
		return fmt.Errorf("启动MySQL失败: %w", err)
	}

	return nil
}
