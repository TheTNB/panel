package bootstrap

import (
	"log"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/migration"
)

func initOrm() {
	db, err := gorm.Open(sqlite.Open(filepath.Join(app.Root, "panel/storage/app.db")), &gorm.Config{
		Logger:                                   slogGorm.New(slogGorm.WithHandler(app.Logger.Handler())),
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	app.Orm = db
}

func runMigrate() {
	migrator := gormigrate.New(app.Orm, &gormigrate.Options{
		UseTransaction: true, // Note: MySQL not support DDL transaction
	}, migration.Migrations)
	if err := migrator.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
