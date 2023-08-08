package redis

var (
	Name        = "Redis"
	Description = "Redis 是一个开源的使用ANSI C语言编写、支持网络、可基于内存亦可持久化的日志型、Key-Value数据库，并提供多种语言的API。"
	Slug        = "redis"
	Version     = "7.0.12"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/redis/install.sh`
	Uninstall   = `bash /www/panel/scripts/redis/uninstall.sh`
	Update      = `bash /www/panel/scripts/redis/update.sh`
)
