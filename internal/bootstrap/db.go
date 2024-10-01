package bootstrap

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/migration"
)

func initOrm() {
	logLevel := logger.Error
	if app.Conf.Bool("database.debug") {
		logLevel = logger.Info
	}
	zapLogger := zapgorm2.New(app.Logger)
	zapLogger.LogMode(logLevel)
	zapLogger.SetAsDefault()

	// You can use any other database, like MySQL or PostgreSQL.
	db, err := gorm.Open(sqlite.Open("storage/panel.db"), &gorm.Config{
		Logger:                                   zapLogger,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}
	app.Orm = db
}

func runMigrate() {
	migrator := gormigrate.New(app.Orm, &gormigrate.Options{
		UseTransaction:            true, // Note: MySQL not support DDL transaction
		ValidateUnknownMigrations: true,
	}, migration.Migrations)
	if err := migrator.Migrate(); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}
}
