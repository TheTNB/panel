package providers

import (
	"fmt"

	"github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/database/gorm"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/app/models"
)

type DatabaseServiceProvider struct {
}

func (receiver *DatabaseServiceProvider) Register(app foundation.Application) {

}

func (receiver *DatabaseServiceProvider) Boot(app foundation.Application) {
	facades.Seeder().Register([]seeder.Seeder{})
	if err := facades.Orm().Query().(*gorm.QueryImpl).Instance().AutoMigrate(
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
	); err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}
}
