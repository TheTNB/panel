package php83

var (
	Name        = "PHP-8.3"
	Description = "PHP 是世界上最好的语言！"
	Slug        = "php83"
	Version     = "8.3.3"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/php/install.sh 83`
	Uninstall   = `bash /www/panel/scripts/php/uninstall.sh 83`
	Update      = `bash /www/panel/scripts/php/install.sh 83`
)
