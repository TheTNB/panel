package service

import "net/http"

type ContainerNetworkService struct {
}

func NewContainerNetworkService() *ContainerNetworkService {
	return &ContainerNetworkService{}
}

func (s *ContainerNetworkService) List(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerNetworkService) Create(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerNetworkService) Remove(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerNetworkService) Exist(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerNetworkService) Inspect(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerNetworkService) Connect(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerNetworkService) Disconnect(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerNetworkService) Prune(w http.ResponseWriter, r *http.Request) {

}
