package app

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

type Cli struct {
	cmd *cli.Command
}

func NewCli(cmd *cli.Command) *Cli {
	IsCli = true
	return &Cli{
		cmd: cmd,
	}
}

func (r *Cli) Run() error {
	return r.cmd.Run(context.Background(), os.Args)
}
