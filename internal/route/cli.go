package route

import (
	"github.com/urfave/cli/v3"

	"github.com/TheTNB/panel/internal/service"
)

func Cli() []*cli.Command {
	cliService := service.NewCliService()
	return []*cli.Command{
		{
			Name:   "restart",
			Usage:  "重启面板服务",
			Action: cliService.Restart,
		},
		{
			Name:   "stop",
			Usage:  "停止面板服务",
			Action: cliService.Stop,
		},
		{
			Name:   "start",
			Usage:  "启动面板服务",
			Action: cliService.Start,
		},
		{
			Name:   "update",
			Usage:  "升级面板",
			Action: cliService.Update,
		},
		{
			Name:   "info",
			Usage:  "输出面板基本信息",
			Action: cliService.Info,
		},
		{
			Name:  "user",
			Usage: "操作面板用户",
			Commands: []*cli.Command{
				{
					Name:   "list",
					Usage:  "列出所有用户",
					Action: cliService.UserList,
				},
				{
					Name:   "username",
					Usage:  "修改用户名",
					Action: cliService.UserName,
				},
				{
					Name:   "password",
					Usage:  "修改用户密码",
					Action: cliService.UserPassword,
				},
			},
		},
		{
			Name:  "https",
			Usage: "操作面板HTTPS",
			Commands: []*cli.Command{
				{
					Name:   "on",
					Usage:  "开启HTTPS",
					Action: cliService.HTTPSOn,
				},
				{
					Name:   "off",
					Usage:  "关闭HTTPS",
					Action: cliService.HTTPSOff,
				},
			},
		},
		{
			Name:  "entrance",
			Usage: "操作面板访问入口",
			Commands: []*cli.Command{
				{
					Name:   "on",
					Usage:  "开启访问入口",
					Action: cliService.EntranceOn,
				},
				{
					Name:   "off",
					Usage:  "关闭访问入口",
					Action: cliService.EntranceOff,
				},
			},
		},
		{
			Name:   "port",
			Usage:  "修改面板端口",
			Action: cliService.Port,
		},
		{
			Name:  "website",
			Usage: "网站管理",
			Commands: []*cli.Command{
				{
					Name:   "create",
					Usage:  "创建新站点",
					Action: cliService.WebsiteCreate,
				},
				{
					Name:   "remove",
					Usage:  "移除站点",
					Action: cliService.WebsiteRemove,
				},
				{
					Name:   "delete",
					Usage:  "删除站点（包括站点目录、同名数据库）",
					Action: cliService.WebsiteDelete,
				},
				{
					Name:   "write",
					Usage:  "写入网站数据（仅限指导下使用）",
					Hidden: true,
					Action: cliService.WebsiteWrite,
				},
			},
		},
		{
			Name:  "backup",
			Usage: "备份数据",
			Commands: []*cli.Command{
				{
					Name:   "website",
					Usage:  "备份网站",
					Action: cliService.BackupWebsite,
				},
				{
					Name:   "database",
					Usage:  "备份数据库",
					Action: cliService.BackupDatabase,
				},
				{
					Name:   "panel",
					Usage:  "备份面板",
					Action: cliService.BackupPanel,
				},
			},
		},
		{
			Name:  "cutoff",
			Usage: "日志切割",
			Commands: []*cli.Command{
				{
					Name:   "website",
					Usage:  "网站",
					Action: cliService.CutoffWebsite,
				},
			},
		},
		{
			Name:  "app",
			Usage: "应用管理",
			Commands: []*cli.Command{
				{
					Name:   "install",
					Usage:  "安装应用",
					Action: cliService.AppInstall,
				},
				{
					Name:   "uninstall",
					Usage:  "卸载应用",
					Action: cliService.AppUnInstall,
				},
				{
					Name:   "write",
					Usage:  "添加面板应用标记（仅限指导下使用）",
					Hidden: true,
					Action: cliService.AppWrite,
				},
				{
					Name:   "remove",
					Usage:  "移除面板应用标记（仅限指导下使用）",
					Hidden: true,
					Action: cliService.AppRemove,
				},
			},
		},
		{
			Name:   "setting",
			Usage:  "设置管理",
			Hidden: true,
			Commands: []*cli.Command{
				{
					Name:   "write",
					Usage:  "写入面板设置（仅限指导下使用）",
					Hidden: true,
					Action: cliService.WriteSetting,
				},
				{
					Name:   "remove",
					Usage:  "移除面板设置（仅限指导下使用）",
					Hidden: true,
					Action: cliService.RemoveSetting,
				},
			},
		},
		{
			Name:   "clearTask",
			Usage:  "清理面板任务队列（仅限指导下使用）",
			Hidden: true,
			Action: cliService.ClearTask,
		},
		{
			Name:   "init",
			Usage:  "初始化面板（仅限指导下使用）",
			Hidden: true,
			Action: cliService.Init,
		},
	}
}
