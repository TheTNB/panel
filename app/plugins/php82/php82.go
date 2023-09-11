package php82

var (
	Name        = "PHP-8.2"
	Description = "PHP 是世界上最好的语言！"
	Slug        = "php82"
	Version     = "8.2.10"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/php/install.sh 82`
	Uninstall   = `bash /www/panel/scripts/php/uninstall.sh 82`
	Update      = `bash /www/panel/scripts/php/install.sh 82`
)
