package service

import "net/http"

type CronService struct {
}

func NewCronService() *CronService {
	return &CronService{}
}

func (s *CronService) List(w http.ResponseWriter, r *http.Request) {

}

func (s *CronService) Create(w http.ResponseWriter, r *http.Request) {

}

func (s *CronService) Update(w http.ResponseWriter, r *http.Request) {

}

func (s *CronService) Get(w http.ResponseWriter, r *http.Request) {

}

func (s *CronService) Delete(w http.ResponseWriter, r *http.Request) {

}

func (s *CronService) Status(w http.ResponseWriter, r *http.Request) {

}

func (s *CronService) Log(w http.ResponseWriter, r *http.Request) {

}
