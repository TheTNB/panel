package providers

import (
	"github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/database/gorm"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/pkg/migrate"
)

type DatabaseServiceProvider struct {
}

func (receiver *DatabaseServiceProvider) Register(app foundation.Application) {

}

func (receiver *DatabaseServiceProvider) Boot(app foundation.Application) {
	facades.Seeder().Register([]seeder.Seeder{})
	migrate.Migrate(facades.Orm().Query().(*gorm.QueryImpl).Instance())
}
