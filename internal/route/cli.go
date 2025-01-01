package route

import (
	"github.com/urfave/cli/v3"

	"github.com/tnb-labs/panel/internal/service"
)

type Cli struct {
	cli *service.CliService
}

func NewCli(cli *service.CliService) *Cli {
	return &Cli{
		cli: cli,
	}
}

func (route *Cli) Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "restart",
			Usage:  "重启面板服务",
			Action: route.cli.Restart,
		},
		{
			Name:   "stop",
			Usage:  "停止面板服务",
			Action: route.cli.Stop,
		},
		{
			Name:   "start",
			Usage:  "启动面板服务",
			Action: route.cli.Start,
		},
		{
			Name:   "update",
			Usage:  "更新面板",
			Action: route.cli.Update,
		},
		{
			Name:   "sync",
			Usage:  "同步数据",
			Action: route.cli.Sync,
		},
		{
			Name:   "fix",
			Usage:  "修复面板",
			Action: route.cli.Fix,
		},
		{
			Name:   "info",
			Usage:  "输出面板基本信息",
			Action: route.cli.Info,
		},
		{
			Name:  "user",
			Usage: "操作面板用户",
			Commands: []*cli.Command{
				{
					Name:   "list",
					Usage:  "列出所有用户",
					Action: route.cli.UserList,
				},
				{
					Name:   "username",
					Usage:  "修改用户名",
					Action: route.cli.UserName,
				},
				{
					Name:   "password",
					Usage:  "修改用户密码",
					Action: route.cli.UserPassword,
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
					Action: route.cli.HTTPSOn,
				},
				{
					Name:   "off",
					Usage:  "关闭HTTPS",
					Action: route.cli.HTTPSOff,
				},
				{
					Name:   "generate",
					Usage:  "生成HTTPS证书",
					Action: route.cli.HTTPSGenerate,
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
					Action: route.cli.EntranceOn,
				},
				{
					Name:   "off",
					Usage:  "关闭访问入口",
					Action: route.cli.EntranceOff,
				},
			},
		},
		{
			Name:   "port",
			Usage:  "修改面板端口",
			Action: route.cli.Port,
		},
		{
			Name:  "website",
			Usage: "网站管理",
			Commands: []*cli.Command{
				{
					Name:   "create",
					Usage:  "创建新网站",
					Action: route.cli.WebsiteCreate,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Usage:    "网站名称",
							Aliases:  []string{"n"},
							Required: true,
						},
						&cli.StringSliceFlag{
							Name:     "domains",
							Usage:    "与网站关联的域名列表",
							Aliases:  []string{"d"},
							Required: true,
						},
						&cli.StringSliceFlag{
							Name:     "listens",
							Usage:    "与网站关联的监听地址列表",
							Aliases:  []string{"l"},
							Required: true,
						},
						&cli.StringFlag{
							Name:  "path",
							Usage: "网站托管的路径（不填则默认路径）",
						},
						&cli.IntFlag{
							Name:  "php",
							Usage: "网站使用的 PHP 版本（不填不使用）",
						},
					},
				},
				{
					Name:   "remove",
					Usage:  "移除网站",
					Action: route.cli.WebsiteRemove,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Usage:    "网站名称",
							Aliases:  []string{"n"},
							Required: true,
						},
					},
				},
				{
					Name:   "delete",
					Usage:  "删除网站（包括网站目录、同名数据库）",
					Action: route.cli.WebsiteDelete,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Usage:    "网站名称",
							Aliases:  []string{"n"},
							Required: true,
						},
					},
				},
				{
					Name:   "write",
					Usage:  "写入网站数据（仅限指导下使用）",
					Hidden: true,
					Action: route.cli.WebsiteWrite,
				},
			},
		},
		{
			Name:  "database",
			Usage: "数据库管理",
			Commands: []*cli.Command{
				{
					Name:   "add-server",
					Usage:  "添加数据库服务器",
					Action: route.cli.DatabaseAddServer,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "type",
							Usage:    "服务器类型",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "name",
							Usage:    "服务器名称",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "host",
							Usage:    "服务器地址",
							Required: true,
						},
						&cli.UintFlag{
							Name:     "port",
							Usage:    "服务器端口",
							Required: true,
						},
						&cli.StringFlag{
							Name:  "username",
							Usage: "服务器用户名",
						},
						&cli.StringFlag{
							Name:  "password",
							Usage: "服务器密码",
						},
						&cli.StringFlag{
							Name:  "remark",
							Usage: "服务器备注",
						},
					},
				},
				{
					Name:   "delete-server",
					Usage:  "删除数据库服务器",
					Action: route.cli.DatabaseDeleteServer,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Usage:    "服务器名称",
							Aliases:  []string{"n"},
							Required: true,
						},
					},
				},
			},
		},
		{
			Name:  "backup",
			Usage: "数据备份",
			Commands: []*cli.Command{
				{
					Name:   "website",
					Usage:  "备份网站",
					Action: route.cli.BackupWebsite,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Aliases:  []string{"n"},
							Usage:    "网站名称",
							Required: true,
						},
						&cli.StringFlag{
							Name:    "path",
							Aliases: []string{"p"},
							Usage:   "保存目录（不填则默认路径）",
						},
					},
				},
				{
					Name:   "database",
					Usage:  "备份数据库",
					Action: route.cli.BackupDatabase,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "type",
							Aliases:  []string{"t"},
							Usage:    "数据库类型",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "name",
							Aliases:  []string{"n"},
							Usage:    "数据库名称",
							Required: true,
						},
						&cli.StringFlag{
							Name:    "path",
							Aliases: []string{"p"},
							Usage:   "保存目录（不填则默认路径）",
						},
					},
				},
				{
					Name:   "panel",
					Usage:  "备份面板",
					Action: route.cli.BackupPanel,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "path",
							Aliases: []string{"p"},
							Usage:   "保存目录（不填则默认路径）",
						},
					},
				},
				{
					Name:   "clear",
					Usage:  "清理备份",
					Action: route.cli.BackupClear,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "type",
							Aliases:  []string{"t"},
							Usage:    "备份类型",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "file",
							Aliases:  []string{"f"},
							Usage:    "备份文件",
							Required: true,
						},
						&cli.IntFlag{
							Name:     "save",
							Aliases:  []string{"s"},
							Usage:    "保存份数",
							Required: true,
						},
						&cli.StringFlag{
							Name:    "path",
							Aliases: []string{"p"},
							Usage:   "备份目录（不填则默认路径）",
						},
					},
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
					Action: route.cli.CutoffWebsite,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Aliases:  []string{"n"},
							Usage:    "网站名称",
							Required: true,
						},

						&cli.StringFlag{
							Name:    "path",
							Aliases: []string{"p"},
							Usage:   "保存目录（不填则默认路径）",
						},
					},
				},
				{
					Name:   "clear",
					Usage:  "清理切割的日志",
					Action: route.cli.CutoffClear,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "type",
							Aliases:  []string{"t"},
							Usage:    "切割类型",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "file",
							Aliases:  []string{"f"},
							Usage:    "切割文件",
							Required: true,
						},
						&cli.IntFlag{
							Name:     "save",
							Aliases:  []string{"s"},
							Usage:    "保存份数",
							Required: true,
						},
						&cli.StringFlag{
							Name:    "path",
							Aliases: []string{"p"},
							Usage:   "切割目录（不填则默认路径）",
						},
					},
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
					Action: route.cli.AppInstall,
				},
				{
					Name:   "uninstall",
					Usage:  "卸载应用",
					Action: route.cli.AppUnInstall,
				},
				{
					Name:   "update",
					Usage:  "更新应用",
					Action: route.cli.AppUpdate,
				},
				{
					Name:   "write",
					Usage:  "添加面板应用标记（仅限指导下使用）",
					Hidden: true,
					Action: route.cli.AppWrite,
				},
				{
					Name:   "remove",
					Usage:  "移除面板应用标记（仅限指导下使用）",
					Hidden: true,
					Action: route.cli.AppRemove,
				},
			},
		},
		{
			Name:   "setting",
			Usage:  "设置管理",
			Hidden: true,
			Commands: []*cli.Command{
				{
					Name:   "get",
					Usage:  "获取面板设置（仅限指导下使用）",
					Hidden: true,
					Action: route.cli.GetSetting,
				},
				{
					Name:   "write",
					Usage:  "写入面板设置（仅限指导下使用）",
					Hidden: true,
					Action: route.cli.WriteSetting,
				},
				{
					Name:   "remove",
					Usage:  "移除面板设置（仅限指导下使用）",
					Hidden: true,
					Action: route.cli.RemoveSetting,
				},
			},
		},
		{
			Name:   "sync-time",
			Usage:  "同步系统时间",
			Action: route.cli.SyncTime,
		},
		{
			Name:   "clear-task",
			Usage:  "清理面板任务队列（仅限指导下使用）",
			Hidden: true,
			Action: route.cli.ClearTask,
		},
		{
			Name:   "init",
			Usage:  "初始化面板（仅限指导下使用）",
			Hidden: true,
			Action: route.cli.Init,
		},
	}
}
