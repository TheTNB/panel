package service

import (
	"net/http"
	"slices"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/firewall"
	"github.com/TheTNB/panel/pkg/os"
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
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, running)
}

func (s *FirewallService) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FirewallStatus](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
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
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *FirewallService) GetRules(w http.ResponseWriter, r *http.Request) {
	rules, err := s.firewall.ListRule()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	var filledRules []map[string]any
	for rule := range slices.Values(rules) {
		// 去除IP规则
		if rule.PortStart == 1 && rule.PortEnd == 65535 {
			continue
		}
		isUse := false
		for port := rule.PortStart; port <= rule.PortEnd; port++ {
			if rule.Protocol == firewall.ProtocolTCP {
				isUse = os.TCPPortInUse(port)
			} else if rule.Protocol == firewall.ProtocolUDP {
				isUse = os.UDPPortInUse(port)
			} else {
				isUse = os.TCPPortInUse(port) || os.UDPPortInUse(port)
			}
			if isUse {
				break
			}
		}
		filledRules = append(filledRules, map[string]any{
			"family":     rule.Family,
			"port_start": rule.PortStart,
			"port_end":   rule.PortEnd,
			"protocol":   rule.Protocol,
			"address":    rule.Address,
			"strategy":   rule.Strategy,
			"direction":  rule.Direction,
			"in_use":     isUse,
		})
	}

	paged, total := Paginate(r, filledRules)

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

func (s *FirewallService) CreateRule(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FirewallRule](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.firewall.Port(firewall.FireInfo{
		Family: req.Family, PortStart: req.PortStart, PortEnd: req.PortEnd, Protocol: firewall.Protocol(req.Protocol), Address: req.Address, Strategy: firewall.Strategy(req.Strategy), Direction: firewall.Direction(req.Direction),
	}, firewall.OperationAdd); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *FirewallService) DeleteRule(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FirewallRule](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.firewall.Port(firewall.FireInfo{
		Family: req.Family, PortStart: req.PortStart, PortEnd: req.PortEnd, Protocol: firewall.Protocol(req.Protocol), Address: req.Address, Strategy: firewall.Strategy(req.Strategy), Direction: firewall.Direction(req.Direction),
	}, firewall.OperationRemove); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *FirewallService) GetIPRules(w http.ResponseWriter, r *http.Request) {
	rules, err := s.firewall.ListRule()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	var filledRules []map[string]any
	for rule := range slices.Values(rules) {
		// 保留IP规则
		if rule.PortStart != 1 || rule.PortEnd != 65535 || rule.Address == "" {
			continue
		}
		filledRules = append(filledRules, map[string]any{
			"family":    rule.Family,
			"protocol":  rule.Protocol,
			"address":   rule.Address,
			"strategy":  rule.Strategy,
			"direction": rule.Direction,
		})
	}

	paged, total := Paginate(r, filledRules)

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

func (s *FirewallService) CreateIPRule(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FirewallIPRule](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.firewall.RichRules(firewall.FireInfo{
		Family: req.Family, Address: req.Address, Protocol: firewall.Protocol(req.Protocol), Strategy: firewall.Strategy(req.Strategy), Direction: firewall.Direction(req.Direction),
	}, firewall.OperationAdd); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *FirewallService) DeleteIPRule(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FirewallIPRule](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.firewall.RichRules(firewall.FireInfo{
		Family: req.Family, Address: req.Address, Protocol: firewall.Protocol(req.Protocol), Strategy: firewall.Strategy(req.Strategy), Direction: firewall.Direction(req.Direction),
	}, firewall.OperationRemove); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *FirewallService) GetForwards(w http.ResponseWriter, r *http.Request) {
	forwards, err := s.firewall.ListForward()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	paged, total := Paginate(r, forwards)

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

func (s *FirewallService) CreateForward(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FirewallForward](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.firewall.Forward(firewall.Forward{
		Protocol: firewall.Protocol(req.Protocol), Port: req.Port, TargetIP: req.TargetIP, TargetPort: req.TargetPort,
	}, firewall.OperationAdd); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *FirewallService) DeleteForward(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FirewallForward](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.firewall.Forward(firewall.Forward{
		Protocol: firewall.Protocol(req.Protocol), Port: req.Port, TargetIP: req.TargetIP, TargetPort: req.TargetPort,
	}, firewall.OperationRemove); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}
