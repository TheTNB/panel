package app

import (
	"context"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/urfave/cli/v3"
)

type Cli struct {
	cmd      *cli.Command
	migrator *gormigrate.Gormigrate
}

func NewCli(cmd *cli.Command, migrator *gormigrate.Gormigrate) *Cli {
	IsCli = true
	return &Cli{
		cmd:      cmd,
		migrator: migrator,
	}
}

func (r *Cli) Run() error {
	// migrate database
	// 这里不处理错误，这么做是为了在异常时用户可以用 fix 命令尝试修复
	_ = r.migrator.Migrate()

	return r.cmd.Run(context.Background(), os.Args)
}
