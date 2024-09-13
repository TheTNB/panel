package service

import "net/http"

type ContainerService struct {
}

func NewContainerService() *ContainerService {
	return &ContainerService{}
}

func (s *ContainerService) List(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Search(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Create(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Remove(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Start(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Stop(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Restart(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Pause(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Unpause(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Inspect(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Kill(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Rename(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Stats(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Exist(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Logs(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerService) Prune(w http.ResponseWriter, r *http.Request) {

}
