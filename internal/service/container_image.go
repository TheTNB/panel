package service

import "net/http"

type ContainerImageService struct {
}

func NewContainerImageService() *ContainerImageService {
	return &ContainerImageService{}
}

func (s *ContainerImageService) List(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerImageService) Exist(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerImageService) Pull(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerImageService) Remove(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerImageService) Inspect(w http.ResponseWriter, r *http.Request) {

}

func (s *ContainerImageService) Prune(w http.ResponseWriter, r *http.Request) {

}
