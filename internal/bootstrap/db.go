package bootstrap

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/TheTNB/panel/internal/migration"
	"github.com/TheTNB/panel/internal/panel"
)

func initOrm() {
	logLevel := logger.Error
	if panel.Conf.Bool("database.debug") {
		logLevel = logger.Info
	}
	// You can use any other database, like MySQL or PostgreSQL.
	db, err := gorm.Open(sqlite.Open("storage/panel.db"), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}
	panel.Orm = db
}

func runMigrate() {
	migrator := gormigrate.New(panel.Orm, &gormigrate.Options{
		UseTransaction:            true, // Note: MySQL not support DDL transaction
		ValidateUnknownMigrations: true,
	}, migration.Migrations)
	if err := migrator.Migrate(); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}
}
