package bootstrap

import (
	"log/slog"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/knadh/koanf/v2"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"

	"github.com/tnb-labs/panel/internal/app"
	"github.com/tnb-labs/panel/internal/migration"
)

func NewDB(conf *koanf.Koanf, log *slog.Logger) (*gorm.DB, error) {
	// You can use any other database, like MySQL or PostgreSQL.
	return gorm.Open(sqlite.Open(filepath.Join(app.Root, "panel/storage/app.db")), &gorm.Config{
		Logger:                                   slogGorm.New(slogGorm.WithHandler(log.Handler())),
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

func NewMigrate(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, &gormigrate.Options{
		UseTransaction: true, // Note: MySQL not support DDL transaction
	}, migration.Migrations)
}
