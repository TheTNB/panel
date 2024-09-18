package route

import (
	"github.com/urfave/cli/v2"

	"github.com/TheTNB/panel/internal/service"
)

func Cli() []*cli.Command {
	cliService := service.NewCliService()
	return []*cli.Command{
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "print a test message",
			Action:  cliService.Test,
		},
	}
}
