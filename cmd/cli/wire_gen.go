// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/tnb-labs/panel/internal/app"
	"github.com/tnb-labs/panel/internal/bootstrap"
	"github.com/tnb-labs/panel/internal/data"
	"github.com/tnb-labs/panel/internal/route"
	"github.com/tnb-labs/panel/internal/service"
)

import (
	_ "time/tzdata"
)

// Injectors from wire.go:

// initCli init command line.
func initCli() (*app.Cli, error) {
	koanf, err := bootstrap.NewConf()
	if err != nil {
		return nil, err
	}
	logger := bootstrap.NewLog(koanf)
	db, err := bootstrap.NewDB(koanf, logger)
	if err != nil {
		return nil, err
	}
	cacheRepo := data.NewCacheRepo(db)
	queue := bootstrap.NewQueue()
	taskRepo := data.NewTaskRepo(db, logger, queue)
	appRepo := data.NewAppRepo(db, cacheRepo, taskRepo)
	userRepo := data.NewUserRepo(db)
	settingRepo := data.NewSettingRepo(db, koanf, taskRepo)
	databaseServerRepo := data.NewDatabaseServerRepo(db, logger)
	databaseUserRepo := data.NewDatabaseUserRepo(db, databaseServerRepo)
	databaseRepo := data.NewDatabaseRepo(db, databaseServerRepo, databaseUserRepo)
	certRepo := data.NewCertRepo(db)
	certAccountRepo := data.NewCertAccountRepo(db, userRepo)
	websiteRepo := data.NewWebsiteRepo(db, cacheRepo, databaseRepo, databaseServerRepo, databaseUserRepo, certRepo, certAccountRepo)
	backupRepo := data.NewBackupRepo(db, settingRepo, websiteRepo)
	cliService := service.NewCliService(koanf, db, appRepo, cacheRepo, userRepo, settingRepo, backupRepo, websiteRepo, databaseServerRepo)
	cli := route.NewCli(cliService)
	command := bootstrap.NewCli(cli)
	gormigrate := bootstrap.NewMigrate(db)
	appCli := app.NewCli(command, gormigrate)
	return appCli, nil
}
