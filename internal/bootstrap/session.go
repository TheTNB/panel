package bootstrap

import (
	"fmt"

	"github.com/go-rat/gormstore"
	"github.com/go-rat/sessions"

	"github.com/TheTNB/panel/internal/panel"
)

func initSession() {
	// initialize session manager
	manager, err := sessions.NewManager(&sessions.ManagerOptions{
		Key:                  panel.Conf.String("app.key"),
		Lifetime:             120,
		GcInterval:           30,
		DisableDefaultDriver: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to initialize session manager: %v", err))
	}

	// extend gorm store driver
	store := gormstore.New(panel.Orm)
	if err = manager.Extend("default", store); err != nil {
		panic(fmt.Sprintf("failed to extend session manager: %v", err))
	}

	panel.Session = manager
}
