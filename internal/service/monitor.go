package service

import "net/http"

type MonitorService struct {
}

func NewMonitorService() *MonitorService {
	return &MonitorService{}
}

func (s *MonitorService) GetSetting(w http.ResponseWriter, r *http.Request) {

}

func (s *MonitorService) UpdateSetting(w http.ResponseWriter, r *http.Request) {

}

func (s *MonitorService) Clear(w http.ResponseWriter, r *http.Request) {

}

func (s *MonitorService) List(w http.ResponseWriter, r *http.Request) {

}
