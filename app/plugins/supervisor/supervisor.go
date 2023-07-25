package supervisor

var (
	Name        = "Supervisor"
	Description = "Supervisor 是一个客户端/服务器系统，允许用户监视和控制类 UNIX 操作系统上的多个进程。"
	Slug        = "supervisor"
	Version     = "4.2.5"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/supervisor/install.sh`
	Uninstall   = `bash /www/panel/scripts/supervisor/uninstall.sh`
	Update      = `bash /www/panel/scripts/supervisor/install.sh`
)
