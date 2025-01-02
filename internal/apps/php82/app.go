package php82

import (
	"github.com/go-chi/chi/v5"

	"github.com/tnb-labs/panel/internal/apps/php"
	"github.com/tnb-labs/panel/internal/biz"
)

type App struct {
	php *php.App
}

func NewApp(task biz.TaskRepo) *App {
	return &App{
		php: php.NewApp(task),
	}
}

func (s *App) Route(r chi.Router) {
	s.php.Route(82)(r)
}
