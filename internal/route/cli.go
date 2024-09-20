package route

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func Cli() []*cli.Command {
	//cliService := service.NewCliService()
	return []*cli.Command{
		{
			Name:  "restart",
			Usage: "重启面板服务",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "停止面板服务",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:  "start",
			Usage: "启动面板服务",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:  "update",
			Usage: "升级面板主程序",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:  "update-cli",
			Usage: "升级面板命令行工具",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:  "info",
			Usage: "输出面板基本信息",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:  "user",
			Usage: "操作面板用户",
			Commands: []*cli.Command{
				{
					Name:  "list",
					Usage: "列出所有用户",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("completed task: ", cmd.Args().First())
						return nil
					},
				},
				{
					Name:  "username",
					Usage: "修改用户名",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("completed task: ", cmd.Args().First())
						return nil
					},
				},
				{
					Name:  "password",
					Usage: "修改用户密码",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("completed task: ", cmd.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:  "https",
			Usage: "操作面板HTTPS",
			Commands: []*cli.Command{
				{
					Name:  "on",
					Usage: "开启HTTPS",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("completed task: ", cmd.Args().First())
						return nil
					},
				},
				{
					Name:  "off",
					Usage: "关闭HTTPS",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("completed task: ", cmd.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:  "entrance",
			Usage: "操作面板访问入口",
			Commands: []*cli.Command{
				{
					Name:  "on",
					Usage: "开启访问入口",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("completed task: ", cmd.Args().First())
						return nil
					},
				},
				{
					Name:  "off",
					Usage: "关闭访问入口",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("completed task: ", cmd.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:  "port",
			Usage: "修改面板端口",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:  "website",
			Usage: "网站管理",
			Commands: []*cli.Command{
				{
					Name:  "add",
					Usage: "创建新站点",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("new task template: ", cmd.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "移除站点",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("removed task template: ", cmd.Args().First())
						return nil
					},
				},
				{
					Name:  "delete",
					Usage: "删除站点（包括站点目录、同名数据库）",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("removed task template: ", cmd.Args().First())
						return nil
					},
				},
				{
					Name:   "write",
					Usage:  "写入网站数据（仅限指导下使用）",
					Hidden: true,
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("removed task template: ", cmd.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:  "backup",
			Usage: "备份数据",
			Commands: []*cli.Command{
				{
					Name:  "website",
					Usage: "备份网站",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("new task template: ", cmd.Args().First())
						return nil
					},
				},
				{
					Name:  "database",
					Usage: "备份数据库",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("removed task template: ", cmd.Args().First())
						return nil
					},
				},
				{
					Name:  "panel",
					Usage: "备份面板",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("removed task template: ", cmd.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:  "cutoff",
			Usage: "日志切割",
			Commands: []*cli.Command{
				{
					Name:  "website",
					Usage: "网站",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("new task template: ", cmd.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:   "writeApp",
			Usage:  "添加面板应用标记（仅限指导下使用）",
			Hidden: true,
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:   "removeApp",
			Usage:  "移除面板应用标记（仅限指导下使用）",
			Hidden: true,
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:   "cleanTask",
			Usage:  "清理面板任务队列（仅限指导下使用）",
			Hidden: true,
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:   "writeSetting",
			Usage:  "写入面板设置（仅限指导下使用）",
			Hidden: true,
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:   "removeSetting",
			Usage:  "移除面板设置（仅限指导下使用）",
			Hidden: true,
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("completed task: ", cmd.Args().First())
				return nil
			},
		},
		{
			Name:   "init",
			Usage:  "初始化面板（仅限指导下使用）",
			Hidden: true,
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("added task: ", cmd.Args().First())
				return nil
			},
		},
	}
}
