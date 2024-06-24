package migrate

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	options := &gormigrate.Options{
		TableName:    "new_migrations",
		IDColumnName: "id",
		IDColumnSize: 255,
	}
	migrator := gormigrate.New(db, options, []*gormigrate.Migration{
		Init,
	})
	if err := migrator.Migrate(); err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}
}
