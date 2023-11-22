package rsync

var (
	Name        = "Rsync"
	Description = "Rsync 是一款提供快速增量文件传输的开源工具。"
	Slug        = "rsync"
	Version     = "3.2.7"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/rsync/install.sh`
	Uninstall   = `bash /www/panel/scripts/rsync/uninstall.sh`
	Update      = `bash /www/panel/scripts/rsync/install.sh`
)
