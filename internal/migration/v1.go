package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/biz"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20240812-init",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&biz.Cache{},
				&biz.Cert{},
				&biz.CertDNS{},
				&biz.CertAccount{},
				&biz.Cron{},
				&biz.Monitor{},
				&biz.App{},
				&biz.Setting{},
				&biz.Task{},
				&biz.User{},
				&biz.Website{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&biz.Cert{},
				&biz.CertDNS{},
				&biz.CertAccount{},
				&biz.Cron{},
				&biz.Monitor{},
				&biz.App{},
				&biz.Setting{},
				&biz.Task{},
				&biz.User{},
				&biz.Website{},
			)
		},
	})
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20241022-ssh",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&biz.SSH{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&biz.SSH{},
			)
		},
	})
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20241107-database",
		Migrate: func(tx *gorm.DB) error {
			_ = tx.Migrator().DropTable(&biz.Database{})
			return tx.AutoMigrate(
				&biz.DatabaseServer{},
				&biz.Database{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&biz.DatabaseServer{},
				&biz.Database{},
			)
		},
	})
}
