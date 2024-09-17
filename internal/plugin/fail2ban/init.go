package fail2ban

import (
	"github.com/go-chi/chi/v5"

	"github.com/TheTNB/panel/pkg/pluginloader"
	"github.com/TheTNB/panel/pkg/types"
)

func init() {
	pluginloader.Register(&types.Plugin{
		Slug:        "fail2ban",
		Name:        "Fail2ban",
		Description: "Fail2ban 扫描系统日志文件并从中找出多次尝试失败的IP地址，将该IP地址加入防火墙的拒绝访问列表中",
		Version:     "1.0.2",
		Requires:    []string{},
		Excludes:    []string{},
		Install:     `bash /www/panel/scripts/fail2ban/install.sh`,
		Uninstall:   `bash /www/panel/scripts/fail2ban/uninstall.sh`,
		Update:      `bash /www/panel/scripts/fail2ban/update.sh`,
		Route: func(r chi.Router) {
			service := NewService()
			r.Get("/jails", service.List)
			r.Post("/jails", service.Add)
			r.Delete("/jails", service.Delete)
			r.Get("/jails/{name}", service.BanList)
			r.Post("/unban", service.Unban)
			r.Post("/whiteList", service.SetWhiteList)
			r.Get("/whiteList", service.GetWhiteList)
		},
	})
}
