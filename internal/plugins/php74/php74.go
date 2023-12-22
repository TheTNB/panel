package php74

var (
	Name        = "PHP-7.4"
	Description = "PHP 是世界上最好的语言！（已停止维护，不建议使用！预计 2024 年 12 月移除）"
	Slug        = "php74"
	Version     = "7.4.33"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/php/install.sh 74`
	Uninstall   = `bash /www/panel/scripts/php/uninstall.sh 74`
	Update      = `bash /www/panel/scripts/php/install.sh 74`
)
