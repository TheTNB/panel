package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/tnb-labs/panel/internal/biz"
	"github.com/tnb-labs/panel/internal/http/request"
)

type ContainerVolumeService struct {
	containerVolumeRepo biz.ContainerVolumeRepo
}

func NewContainerVolumeService(containerVolume biz.ContainerVolumeRepo) *ContainerVolumeService {
	return &ContainerVolumeService{
		containerVolumeRepo: containerVolume,
	}
}

func (s *ContainerVolumeService) List(w http.ResponseWriter, r *http.Request) {
	volumes, err := s.containerVolumeRepo.List()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	paged, total := Paginate(r, volumes)

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

func (s *ContainerVolumeService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerVolumeCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	name, err := s.containerVolumeRepo.Create(req)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, name)

}

func (s *ContainerVolumeService) Remove(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerVolumeID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerVolumeRepo.Remove(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerVolumeService) Prune(w http.ResponseWriter, r *http.Request) {
	if err := s.containerVolumeRepo.Prune(); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}
