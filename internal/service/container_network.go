package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/tnb-labs/panel/internal/biz"
	"github.com/tnb-labs/panel/internal/http/request"
)

type ContainerNetworkService struct {
	containerNetworkRepo biz.ContainerNetworkRepo
}

func NewContainerNetworkService(containerNetwork biz.ContainerNetworkRepo) *ContainerNetworkService {
	return &ContainerNetworkService{
		containerNetworkRepo: containerNetwork,
	}
}

func (s *ContainerNetworkService) List(w http.ResponseWriter, r *http.Request) {
	networks, err := s.containerNetworkRepo.List()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	paged, total := Paginate(r, networks)

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

func (s *ContainerNetworkService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerNetworkCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	id, err := s.containerNetworkRepo.Create(req)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, id)
}

func (s *ContainerNetworkService) Remove(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerNetworkID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerNetworkRepo.Remove(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerNetworkService) Prune(w http.ResponseWriter, r *http.Request) {
	if err := s.containerNetworkRepo.Prune(); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}
