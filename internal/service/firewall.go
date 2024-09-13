package service

import "net/http"

type FirewallService struct {
}

func NewFirewallService() *FirewallService {
	return &FirewallService{}
}

func (s *FirewallService) GetStatus(w http.ResponseWriter, r *http.Request) {

}

func (s *FirewallService) UpdateStatus(w http.ResponseWriter, r *http.Request) {

}

func (s *FirewallService) GetRules(w http.ResponseWriter, r *http.Request) {

}

func (s *FirewallService) CreateRule(w http.ResponseWriter, r *http.Request) {

}

func (s *FirewallService) DeleteRule(w http.ResponseWriter, r *http.Request) {

}
