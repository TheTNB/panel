package bootstrap

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/TheTNB/panel/internal/panel"
	"github.com/TheTNB/panel/internal/route"
)

func initCli() {
	app := &cli.App{
		Name:        "panel-cli",
		HelpName:    fmt.Sprintf("耗子面板 %s", panel.Version),
		Usage:       "命令行工具",
		UsageText:   "panel-cli [global options] command [command options] [arguments...]",
		HideVersion: true,
		Commands:    route.Cli(),
	}

	if err := app.Run(os.Args); err != nil {
		panic(fmt.Sprintf("failed to run cli: %v", err))
	}
}
