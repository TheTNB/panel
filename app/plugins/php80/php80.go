package php80

var (
	Name        = "PHP-8.0"
	Description = "PHP 是世界上最好的语言！"
	Slug        = "php80"
	Version     = "8.0.29"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/php/install.sh 80`
	Uninstall   = `bash /www/panel/scripts/php/uninstall.sh 80`
	Update      = `bash /www/panel/scripts/php/install.sh 80`
)
