package app

import (
	"context"
	"fmt"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/urfave/cli/v3"

	"github.com/tnb-labs/panel/pkg/apploader"
)

type Cli struct {
	cmd      *cli.Command
	migrator *gormigrate.Gormigrate
}

func NewCli(cmd *cli.Command, migrator *gormigrate.Gormigrate, _ *apploader.Loader) *Cli {
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

	if err := r.cmd.Run(context.TODO(), os.Args); err != nil {
		fmt.Printf("|-%v\n", err)
	}

	return nil
}
