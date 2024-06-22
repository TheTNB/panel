package types

var PluginOpenResty = Plugin{
	Name:        "OpenResty",
	Description: "OpenResty® 是一款基于 NGINX 和 LuaJIT 的 Web 平台",
	Slug:        "openresty",
	Version:     "1.25.3.1",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     "bash /www/panel/scripts/openresty/install.sh",
	Uninstall:   "bash /www/panel/scripts/openresty/uninstall.sh",
	Update:      "bash /www/panel/scripts/openresty/install.sh",
}

var PluginMySQL57 = Plugin{
	Name:        "MySQL-5.7",
	Description: "MySQL 是最流行的关系型数据库管理系统之一，Oracle 旗下产品（已停止维护，不建议使用！预计 2025 年 12 月移除）",
	Slug:        "mysql57",
	Version:     "5.7.44",
	Requires:    []string{},
	Excludes:    []string{"mysql80", "mysql84"},
	Install:     `bash /www/panel/scripts/mysql/install.sh 57`,
	Uninstall:   `bash /www/panel/scripts/mysql/uninstall.sh 57`,
	Update:      `bash /www/panel/scripts/mysql/update.sh 57`,
}

var PluginMySQL80 = Plugin{
	Name:        "MySQL-8.0",
	Description: "MySQL 是最流行的关系型数据库管理系统之一，Oracle 旗下产品（建议内存 > 2G 安装）",
	Slug:        "mysql80",
	Version:     "8.0.37",
	Requires:    []string{},
	Excludes:    []string{"mysql57", "mysql84"},
	Install:     `bash /www/panel/scripts/mysql/install.sh 80`,
	Uninstall:   `bash /www/panel/scripts/mysql/uninstall.sh 80`,
	Update:      `bash /www/panel/scripts/mysql/update.sh 80`,
}

var PluginMySQL84 = Plugin{
	Name:        "MySQL-8.4",
	Description: "MySQL 是最流行的关系型数据库管理系统之一，Oracle 旗下产品（建议内存 > 2G 安装）",
	Slug:        "mysql84",
	Version:     "8.4.0",
	Requires:    []string{},
	Excludes:    []string{"mysql57", "mysql80"},
	Install:     `bash /www/panel/scripts/mysql/install.sh 84`,
	Uninstall:   `bash /www/panel/scripts/mysql/uninstall.sh 84`,
	Update:      `bash /www/panel/scripts/mysql/update.sh 84`,
}

var PluginPostgreSQL15 = Plugin{
	Name:        "PostgreSQL-15",
	Description: "PostgreSQL 是世界上最先进的开源关系数据库，在类似 BSD 与 MIT 许可的 PostgreSQL 许可下发行",
	Slug:        "postgresql15",
	Version:     "15.7",
	Requires:    []string{},
	Excludes:    []string{"postgresql16"},
	Install:     `bash /www/panel/scripts/postgresql/install.sh 15`,
	Uninstall:   `bash /www/panel/scripts/postgresql/uninstall.sh 15`,
	Update:      `bash /www/panel/scripts/postgresql/update.sh 15`,
}

var PluginPostgreSQL16 = Plugin{
	Name:        "PostgreSQL-16",
	Description: "PostgreSQL 是世界上最先进的开源关系数据库，在类似 BSD 与 MIT 许可的 PostgreSQL 许可下发行",
	Slug:        "postgresql16",
	Version:     "16.3",
	Requires:    []string{},
	Excludes:    []string{"postgresql15"},
	Install:     `bash /www/panel/scripts/postgresql/install.sh 16`,
	Uninstall:   `bash /www/panel/scripts/postgresql/uninstall.sh 16`,
	Update:      `bash /www/panel/scripts/postgresql/update.sh 16`,
}

var PluginPHP74 = Plugin{
	Name:        "PHP-7.4",
	Description: "PHP 是一种创建动态交互性站点的强有力的服务器端脚本语言（已停止维护，不建议使用！预计 2024 年 12 月移除）",
	Slug:        "php74",
	Version:     "7.4.33",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/php/install.sh 74`,
	Uninstall:   `bash /www/panel/scripts/php/uninstall.sh 74`,
	Update:      `bash /www/panel/scripts/php/install.sh 74`,
}

var PluginPHP80 = Plugin{
	Name:        "PHP-8.0",
	Description: "PHP 是一种创建动态交互性站点的强有力的服务器端脚本语言（已停止维护，不建议使用！预计 2025 年 12 月移除）",
	Slug:        "php80",
	Version:     "8.0.30",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/php/install.sh 80`,
	Uninstall:   `bash /www/panel/scripts/php/uninstall.sh 80`,
	Update:      `bash /www/panel/scripts/php/install.sh 80`,
}

var PluginPHP81 = Plugin{
	Name:        "PHP-8.1",
	Description: "PHP 是一种创建动态交互性站点的强有力的服务器端脚本语言",
	Slug:        "php81",
	Version:     "8.1.29",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/php/install.sh 81`,
	Uninstall:   `bash /www/panel/scripts/php/uninstall.sh 81`,
	Update:      `bash /www/panel/scripts/php/install.sh 81`,
}

var PluginPHP82 = Plugin{
	Name:        "PHP-8.2",
	Description: "PHP 是一种创建动态交互性站点的强有力的服务器端脚本语言",
	Slug:        "php82",
	Version:     "8.2.20",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/php/install.sh 82`,
	Uninstall:   `bash /www/panel/scripts/php/uninstall.sh 82`,
	Update:      `bash /www/panel/scripts/php/install.sh 82`,
}

var PluginPHP83 = Plugin{
	Name:        "PHP-8.3",
	Description: "PHP 是一种创建动态交互性站点的强有力的服务器端脚本语言",
	Slug:        "php83",
	Version:     "8.3.8",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/php/install.sh 83`,
	Uninstall:   `bash /www/panel/scripts/php/uninstall.sh 83`,
	Update:      `bash /www/panel/scripts/php/install.sh 83`,
}

var PluginPHPMyAdmin = Plugin{
	Name:        "phpMyAdmin",
	Description: "phpMyAdmin 是一个以 PHP 为基础，以 Web-Base 方式架构在网站主机上的 MySQL 数据库管理工具",
	Slug:        "phpmyadmin",
	Version:     "5.2.1",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/phpmyadmin/install.sh`,
	Uninstall:   `bash /www/panel/scripts/phpmyadmin/uninstall.sh`,
	Update:      `bash /www/panel/scripts/phpmyadmin/uninstall.sh && bash /www/panel/scripts/phpmyadmin/install.sh`,
}

var PluginPureFTPd = Plugin{
	Name:        "Pure-FTPd",
	Description: "Pure-Ftpd 是一个快速、高效、轻便、安全的 FTP 服务器，它以安全和配置简单为设计目标，支持虚拟主机，IPV6，PAM 等功能",
	Slug:        "pureftpd",
	Version:     "1.0.50",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/pureftpd/install.sh`,
	Uninstall:   `bash /www/panel/scripts/pureftpd/uninstall.sh`,
	Update:      `bash /www/panel/scripts/pureftpd/update.sh`,
}

var PluginRedis = Plugin{
	Name:        "Redis",
	Description: "Redis 是一个开源的使用 ANSI C 语言编写、支持网络、可基于内存亦可持久化的日志型、Key-Value 数据库，并提供多种语言的 API",
	Slug:        "redis",
	Version:     "7.2.5",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/redis/install.sh`,
	Uninstall:   `bash /www/panel/scripts/redis/uninstall.sh`,
	Update:      `bash /www/panel/scripts/redis/update.sh`,
}

var PluginS3fs = Plugin{
	Name:        "S3fs",
	Description: "S3fs 通过 FUSE 挂载兼容 S3 标准的存储桶，例如 Amazon S3、阿里云 OSS、腾讯云 COS、七牛云 Kodo 等",
	Slug:        "s3fs",
	Version:     "1.9",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/s3fs/install.sh`,
	Uninstall:   `bash /www/panel/scripts/s3fs/uninstall.sh`,
	Update:      `bash /www/panel/scripts/s3fs/update.sh`,
}

var PluginRsync = Plugin{
	Name:        "Rsync",
	Description: "Rsync 是一款提供快速增量文件传输的开源工具",
	Slug:        "rsync",
	Version:     "3.2.7",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/rsync/install.sh`,
	Uninstall:   `bash /www/panel/scripts/rsync/uninstall.sh`,
	Update:      `bash /www/panel/scripts/rsync/install.sh`,
}

var PluginSupervisor = Plugin{
	Name:        "Supervisor",
	Description: "Supervisor 是一个客户端/服务器系统，允许用户监视和控制类 UNIX 操作系统上的多个进程",
	Slug:        "supervisor",
	Version:     "4.2.5",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/supervisor/install.sh`,
	Uninstall:   `bash /www/panel/scripts/supervisor/uninstall.sh`,
	Update:      `bash /www/panel/scripts/supervisor/update.sh`,
}

var PluginFail2ban = Plugin{
	Name:        "Fail2ban",
	Description: "Fail2ban 扫描系统日志文件并从中找出多次尝试失败的IP地址，将该IP地址加入防火墙的拒绝访问列表中",
	Slug:        "fail2ban",
	Version:     "1.0.2",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/fail2ban/install.sh`,
	Uninstall:   `bash /www/panel/scripts/fail2ban/uninstall.sh`,
	Update:      `bash /www/panel/scripts/fail2ban/update.sh`,
}

var PluginPodman = Plugin{
	Name:        "Podman",
	Description: "Podman（POD MANager）是一款用于管理容器和镜像、挂载到这些容器中的卷以及由容器组构成的 Pod 的工具",
	Slug:        "podman",
	Version:     "4.0.0",
	Requires:    []string{},
	Excludes:    []string{"docker"},
	Install:     `bash /www/panel/scripts/podman/install.sh`,
	Uninstall:   `bash /www/panel/scripts/podman/uninstall.sh`,
	Update:      `bash /www/panel/scripts/podman/update.sh`,
}

var PluginFrp = Plugin{
	Name:        "Frp",
	Description: "frp 是一个专注于内网穿透的高性能的反向代理应用",
	Slug:        "frp",
	Version:     "0.58.0",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/frp/install.sh`,
	Uninstall:   `bash /www/panel/scripts/frp/uninstall.sh`,
	Update:      `bash /www/panel/scripts/frp/update.sh`,
}

var PluginGitea = Plugin{
	Name:        "Gitea",
	Description: "Gitea 是一款极易搭建的自助 Git 服务，它包括 Git 托管、代码审查、团队协作、软件包注册和 CI/CD",
	Slug:        "gitea",
	Version:     "1.22.0",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `bash /www/panel/scripts/gitea/install.sh`,
	Uninstall:   `bash /www/panel/scripts/gitea/uninstall.sh`,
	Update:      `bash /www/panel/scripts/gitea/update.sh`,
}

var PluginToolBox = Plugin{
	Name:        "系统工具箱",
	Description: "可视化调整一些常用的配置项，如 DNS、SWAP、时区等",
	Slug:        "toolbox",
	Version:     "1.0.0",
	Requires:    []string{},
	Excludes:    []string{},
	Install:     `panel writePlugin toolbox 1.0.0`,
	Uninstall:   `panel deletePlugin toolbox`,
	Update:      `panel writePlugin toolbox 1.0.0`,
}
