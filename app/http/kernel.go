package http

import (
	"github.com/goravel/framework/contracts/http"
	sessionmiddleware "github.com/goravel/framework/session/middleware"

	"github.com/TheTNB/panel/v2/app/http/middleware"
)

type Kernel struct {
}

// The application's global HTTP middleware stack.
// These middleware are run during every request to your application.
func (kernel Kernel) Middleware() []http.Middleware {
	return []http.Middleware{
		sessionmiddleware.StartSession(),
		middleware.Log(),
		middleware.Status(),
		middleware.Entrance(),
		middleware.Static(),
	}
}
