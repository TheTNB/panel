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
	if err := r.migrator.Migrate(); err != nil {
		return err
	}

	return r.cmd.Run(context.Background(), os.Args)
}
