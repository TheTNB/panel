package service

import "net/http"

type ContainerVolumeService struct {
}

func NewContainerVolumeService() *ContainerVolumeService {
	return &ContainerVolumeService{}
}

func (s *ContainerVolumeService) List(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerVolumeService) Create(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerVolumeService) Exist(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerVolumeService) Remove(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerVolumeService) Inspect(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerVolumeService) Prune(w http.ResponseWriter, r *http.Request) {

}
