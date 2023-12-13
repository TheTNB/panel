package php81

var (
	Name        = "PHP-8.1"
	Description = "PHP 是世界上最好的语言！"
	Slug        = "php81"
	Version     = "8.1.26"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/php/install.sh 81`
	Uninstall   = `bash /www/panel/scripts/php/uninstall.sh 81`
	Update      = `bash /www/panel/scripts/php/install.sh 81`
)
