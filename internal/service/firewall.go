package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/firewall"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type FirewallService struct {
	firewall *firewall.Firewall
}

func NewFirewallService() *FirewallService {

	return &FirewallService{
		firewall: firewall.NewFirewall(),
	}
}

func (s *FirewallService) GetStatus(w http.ResponseWriter, r *http.Request) {
	running, err := s.firewall.Status()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, running)
}

func (s *FirewallService) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FirewallStatus](r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Status {
		err = systemctl.Start("firewalld")
		if err == nil {
			err = systemctl.Enable("firewalld")
		}
	} else {
		err = systemctl.Stop("firewalld")
		if err == nil {
			err = systemctl.Disable("firewalld")
		}
	}

	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *FirewallService) GetRules(w http.ResponseWriter, r *http.Request) {
	rules, err := s.firewall.ListRule()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	paged, total := Paginate(r, rules)

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

func (s *FirewallService) CreateRule(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FirewallCreateRule](r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.firewall.Port(firewall.FireInfo{Port: req.Port, Protocol: req.Protocol}, "add")
}

func (s *FirewallService) DeleteRule(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FirewallCreateRule](r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.firewall.Port(firewall.FireInfo{Port: req.Port, Protocol: req.Protocol}, "remove")
}
