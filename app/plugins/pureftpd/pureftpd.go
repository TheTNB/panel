package pureftpd

var (
	Name        = "Pure-FTPd"
	Description = "Pure-Ftpd 是一个快速、高效、轻便、安全的 FTP 服务器，它以安全和配置简单为设计目标，支持虚拟主机，IPV6，PAM 等功能。"
	Slug        = "pureftpd"
	Version     = "1.0.50"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/pureftpd/install.sh`
	Uninstall   = `bash /www/panel/scripts/pureftpd/uninstall.sh`
	Update      = `bash /www/panel/scripts/pureftpd/update.sh`
)
