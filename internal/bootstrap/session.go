package bootstrap

import (
	"log"

	"github.com/go-rat/gormstore"
	"github.com/go-rat/sessions"

	"github.com/TheTNB/panel/internal/app"
)

func initSession() {
	// initialize session manager
	manager, err := sessions.NewManager(&sessions.ManagerOptions{
		Key:                  app.Conf.String("app.key"),
		Lifetime:             120,
		GcInterval:           30,
		DisableDefaultDriver: true,
	})
	if err != nil {
		log.Fatalf("failed to initialize session manager: %v", err)
	}

	// extend gorm store driver
	store := gormstore.New(app.Orm)
	if err = manager.Extend("default", store); err != nil {
		log.Fatalf("failed to extend session manager: %v", err)
	}

	app.Session = manager
}
