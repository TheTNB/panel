package memcached

import (
	"bufio"
	"net"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"

	"github.com/tnb-labs/panel/internal/service"
	"github.com/tnb-labs/panel/pkg/io"
	"github.com/tnb-labs/panel/pkg/systemctl"
	"github.com/tnb-labs/panel/pkg/types"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (s *App) Route(r chi.Router) {
	r.Get("/load", s.Load)
	r.Get("/config", s.GetConfig)
	r.Post("/config", s.UpdateConfig)
}

func (s *App) Load(w http.ResponseWriter, r *http.Request) {
	status, err := systemctl.Status("memcached")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to get Memcached status: %v", err)
		return
	}
	if !status {
		service.Success(w, []types.NV{})
		return
	}

	conn, err := net.Dial("tcp", "127.0.0.1:11211")
	if err != nil {
		service.Success(w, []types.NV{})
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("stats\nquit\n"))
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to write to Memcached: %v", err)
		return
	}

	data := make([]types.NV, 0)
	re := regexp.MustCompile(`STAT\s(\S+)\s(\S+)`)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); len(matches) == 3 {
			data = append(data, types.NV{
				Name:  matches[1],
				Value: matches[2],
			})
		}
		if line == "END" {
			break
		}
	}

	if err = scanner.Err(); err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to read from Memcached: %v", err)
		return
	}

	service.Success(w, data)
}

func (s *App) GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read("/etc/systemd/system/memcached.service")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, config)
}

func (s *App) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = io.Write("/etc/systemd/system/memcached.service", req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = systemctl.Restart("memcached"); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	service.Success(w, nil)
}
