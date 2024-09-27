package supervisor

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/apploader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	apploader.Register(&types.App{
		Slug: "supervisor",
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/service", service.Service)
			r.Get("/log", service.Log)
			r.Post("/clearLog", service.ClearLog)
			r.Get("/config", service.GetConfig)
			r.Post("/config", service.UpdateConfig)
			r.Get("/processes", service.Processes)
			r.Post("/processes/{name}/start", service.StartProcess)
			r.Post("/processes/{name}/stop", service.StopProcess)
			r.Post("/processes/{name}/restart", service.RestartProcess)
			r.Get("/processes/{name}/log", service.ProcessLog)
			r.Post("/processes/{name}/clearLog", service.ClearProcessLog)
			r.Get("/processes/{name}", service.ProcessConfig)
			r.Post("/processes/{name}", service.UpdateProcessConfig)
			r.Delete("/processes/{name}", service.DeleteProcess)
			r.Post("/processes", service.CreateProcess)
		},
	})
}
