package fail2ban

var (
	Name        = "Fail2ban"
	Description = "Fail2ban 扫描系统日志文件并从中找出多次尝试失败的IP地址，将该IP地址加入防火墙的拒绝访问列表中。"
	Slug        = "fail2ban"
	Version     = "1.0.0"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/fail2ban/install.sh`
	Uninstall   = `bash /www/panel/scripts/fail2ban/uninstall.sh`
	Update      = `bash /www/panel/scripts/fail2ban/update.sh`
)
