package service

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/TheTNB/panel/pkg/systemctl"
)

type CliService struct {
}

func NewCliService() *CliService {
	return &CliService{}
}

func (s *CliService) Restart(ctx context.Context, cmd *cli.Command) error {
	return systemctl.Restart("panel")
}

func (s *CliService) Stop(ctx context.Context, cmd *cli.Command) error {
	return systemctl.Stop("panel")
}

func (s *CliService) Start(ctx context.Context, cmd *cli.Command) error {
	return systemctl.Start("panel")
}

func (s *CliService) Update(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) Info(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) UserList(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) UserName(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) UserPassword(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) HTTPSOn(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) HTTPSOff(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) EntranceOn(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) EntranceOff(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) Port(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) WebsiteCreate(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) WebsiteRemove(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) WebsiteDelete(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) WebsiteWrite(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) BackupWebsite(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) BackupDatabase(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) BackupPanel(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) CutoffWebsite(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) AppInstall(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) AppUnInstall(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) AppWrite(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) AppRemove(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) ClearTask(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) WriteSetting(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) RemoveSetting(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}

func (s *CliService) Init(ctx context.Context, cmd *cli.Command) error {
	println("Hello, World!")
	return nil
}
