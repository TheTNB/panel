package service

import "github.com/urfave/cli/v2"

type CliService struct {
}

func NewCliService() *CliService {
	return &CliService{}
}

func (s *CliService) Test(c *cli.Context) error {
	println("Hello, World!")
	return nil
}
