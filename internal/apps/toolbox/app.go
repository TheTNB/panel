package toolbox

import (
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-rat/chix"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/ntp"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/pkg/types"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (s *App) Route(r chi.Router) {
	r.Get("/dns", s.GetDNS)
	r.Post("/dns", s.UpdateDNS)
	r.Get("/swap", s.GetSWAP)
	r.Post("/swap", s.UpdateSWAP)
	r.Get("/timezone", s.GetTimezone)
	r.Post("/timezone", s.UpdateTimezone)
	r.Post("/time", s.UpdateTime)
	r.Post("/syncTime", s.SyncTime)
	r.Get("/hostname", s.GetHostname)
	r.Post("/hostname", s.UpdateHostname)
	r.Get("/hosts", s.GetHosts)
	r.Post("/hosts", s.UpdateHosts)
	r.Post("/rootPassword", s.UpdateRootPassword)
}

// GetDNS 获取 DNS 信息
func (s *App) GetDNS(w http.ResponseWriter, r *http.Request) {
	raw, err := io.Read("/etc/resolv.conf")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	match := regexp.MustCompile(`nameserver\s+(\S+)`).FindAllStringSubmatch(raw, -1)
	if len(match) == 0 {
		service.Error(w, http.StatusInternalServerError, "找不到 DNS 信息")
		return
	}

	var dns []string
	for _, m := range match {
		dns = append(dns, m[1])
	}

	service.Success(w, dns)
}

// UpdateDNS 设置 DNS 信息
func (s *App) UpdateDNS(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[DNS](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	var dns string
	dns += "nameserver " + req.DNS1 + "\n"
	dns += "nameserver " + req.DNS2 + "\n"

	if err := io.Write("/etc/resolv.conf", dns, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "写入 DNS 信息失败")
		return
	}

	service.Success(w, nil)
}

// GetSWAP 获取 SWAP 信息
func (s *App) GetSWAP(w http.ResponseWriter, r *http.Request) {
	var total, used, free string
	var size int64
	if io.Exists(filepath.Join(app.Root, "swap")) {
		file, err := io.FileInfo(filepath.Join(app.Root, "swap"))
		if err != nil {
			service.Error(w, http.StatusInternalServerError, "获取 SWAP 信息失败")
			return
		}

		size = file.Size() / 1024 / 1024
		total = tools.FormatBytes(float64(file.Size()))
	} else {
		size = 0
		total = "0.00 B"
	}

	raw, err := shell.Execf("free | grep Swap")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取 SWAP 信息失败")
		return
	}

	match := regexp.MustCompile(`Swap:\s+(\d+)\s+(\d+)\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) >= 4 {
		used = tools.FormatBytes(cast.ToFloat64(match[2]) * 1024)
		free = tools.FormatBytes(cast.ToFloat64(match[3]) * 1024)
	}

	service.Success(w, chix.M{
		"total": total,
		"size":  size,
		"used":  used,
		"free":  free,
	})
}

// UpdateSWAP 设置 SWAP 信息
func (s *App) UpdateSWAP(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[SWAP](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if io.Exists(filepath.Join(app.Root, "swap")) {
		if _, err = shell.Execf("swapoff '%s'", filepath.Join(app.Root, "swap")); err != nil {
			service.Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
		if _, err = shell.Execf("rm -f '%s'", filepath.Join(app.Root, "swap")); err != nil {
			service.Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
		if _, err = shell.Execf(`sed -i "\|^%s|d" /etc/fstab`, filepath.Join(app.Root, "swap")); err != nil {
			service.Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
	}

	if req.Size > 1 {
		var free string
		free, err = shell.Execf("df -k %s | awk '{print $4}' | tail -n 1", app.Root)
		if err != nil {
			service.Error(w, http.StatusInternalServerError, "获取硬盘空间失败")
			return
		}
		if cast.ToInt64(free)*1024 < req.Size*1024*1024 {
			service.Error(w, http.StatusInternalServerError, "硬盘空间不足，当前剩余%s", tools.FormatBytes(cast.ToFloat64(free)))
			return
		}

		btrfsCheck, _ := shell.Execf("df -T %s | awk '{print $2}' | tail -n 1", app.Root)
		if strings.Contains(btrfsCheck, "btrfs") {
			if _, err = shell.Execf("btrfs filesystem mkswapfile --size %dM --uuid clear %s", req.Size, filepath.Join(app.Root, "swap")); err != nil {
				service.Error(w, http.StatusInternalServerError, "%v", err)
				return
			}
		} else {
			if _, err = shell.Execf("dd if=/dev/zero of=%s bs=1M count=%d", filepath.Join(app.Root, "swap"), req.Size); err != nil {
				service.Error(w, http.StatusInternalServerError, "%v", err)
				return
			}
			if _, err = shell.Execf("mkswap -f '%s'", filepath.Join(app.Root, "swap")); err != nil {
				service.Error(w, http.StatusInternalServerError, "%v", err)
				return
			}
			if err = io.Chmod(filepath.Join(app.Root, "swap"), 0600); err != nil {
				service.Error(w, http.StatusInternalServerError, "设置 SWAP 权限失败")
				return
			}
		}
		if _, err = shell.Execf("swapon '%s'", filepath.Join(app.Root, "swap")); err != nil {
			service.Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
		if _, err = shell.Execf("echo '%s    swap    swap    defaults    0 0' >> /etc/fstab", filepath.Join(app.Root, "swap")); err != nil {
			service.Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
	}

	service.Success(w, nil)
}

// GetTimezone 获取时区
func (s *App) GetTimezone(w http.ResponseWriter, r *http.Request) {
	raw, err := shell.Execf("timedatectl | grep zone")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取时区信息失败")
		return
	}

	match := regexp.MustCompile(`zone:\s+(\S+)`).FindStringSubmatch(raw)
	if len(match) == 0 {
		service.Error(w, http.StatusInternalServerError, "找不到时区信息")
		return
	}

	zonesRaw, err := shell.Execf("timedatectl list-timezones")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取时区列表失败")
		return
	}
	zones := strings.Split(zonesRaw, "\n")

	var zonesList []types.LV
	for _, z := range zones {
		zonesList = append(zonesList, types.LV{
			Label: z,
			Value: z,
		})
	}

	service.Success(w, chix.M{
		"timezone":  match[1],
		"timezones": zonesList,
	})
}

// UpdateTimezone 设置时区
func (s *App) UpdateTimezone(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Timezone](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if _, err = shell.Execf("timedatectl set-timezone '%s'", req.Timezone); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}

// UpdateTime 设置时间
func (s *App) UpdateTime(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Time](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = ntp.UpdateSystemTime(req.Time); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)

}

// SyncTime 同步时间
func (s *App) SyncTime(w http.ResponseWriter, r *http.Request) {
	now, err := ntp.Now()
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = ntp.UpdateSystemTime(now); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}

// GetHostname 获取主机名
func (s *App) GetHostname(w http.ResponseWriter, r *http.Request) {
	hostname, _ := io.Read("/etc/hostname")
	service.Success(w, strings.TrimSpace(hostname))
}

// UpdateHostname 设置主机名
func (s *App) UpdateHostname(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Hostname](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if _, err = shell.Execf("hostnamectl set-hostname '%s'", req.Hostname); err != nil {
		// 直接写 /etc/hostname
		if err = io.Write("/etc/hostname", req.Hostname, 0644); err != nil {
			service.Error(w, http.StatusInternalServerError, "写入主机名失败")
			return
		}
	}

	service.Success(w, nil)
}

// GetHosts 获取 hosts 信息
func (s *App) GetHosts(w http.ResponseWriter, r *http.Request) {
	hosts, _ := io.Read("/etc/hosts")
	service.Success(w, hosts)
}

// UpdateHosts 设置 hosts 信息
func (s *App) UpdateHosts(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Hosts](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = io.Write("/etc/hosts", req.Hosts, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "写入 hosts 信息失败")
		return
	}

	service.Success(w, nil)
}

// UpdateRootPassword 设置 root 密码
func (s *App) UpdateRootPassword(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Password](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	req.Password = strings.ReplaceAll(req.Password, `'`, `\'`)
	if _, err = shell.Execf(`yes '%s' | passwd root`, req.Password); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}
