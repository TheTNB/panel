package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/v2/app/models"
)

var Init = &gormigrate.Migration{
	ID: "20240624-init",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(
			&models.Cert{},
			&models.CertDNS{},
			&models.CertUser{},
			&models.Cron{},
			&models.Database{},
			&models.Monitor{},
			&models.Plugin{},
			&models.Setting{},
			&models.Task{},
			&models.User{},
			&models.Website{},
		)
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable(
			&models.Cert{},
			&models.CertDNS{},
			&models.CertUser{},
			&models.Cron{},
			&models.Database{},
			&models.Monitor{},
			&models.Plugin{},
			&models.Setting{},
			&models.Task{},
			&models.User{},
			&models.Website{},
		)
	},
}
