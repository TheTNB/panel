package models

import (
	"fmt"

	"github.com/goravel/framework/facades"
)

func init() {
	if err := facades.Orm().Query().AutoMigrate(
		&Cert{},
		&CertDNS{},
		&CertUser{},
		&Cron{},
		&Database{},
		&Monitor{},
		&Plugin{},
		&Setting{},
		&Task{},
		&User{},
		&Website{},
	); err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}
}
