package service

import (
	"net/http"

	"github.com/go-rat/chix"
	"github.com/golang-module/carbon/v2"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
)

type ContainerNetworkService struct {
	containerNetworkRepo biz.ContainerNetworkRepo
}

func NewContainerNetworkService() *ContainerNetworkService {
	return &ContainerNetworkService{
		containerNetworkRepo: data.NewContainerNetworkRepo(),
	}
}

func (s *ContainerNetworkService) List(w http.ResponseWriter, r *http.Request) {
	networks, err := s.containerNetworkRepo.List()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	paged, total := Paginate(r, networks)

	items := make([]any, len(paged))
	for _, item := range paged {
		var ipamConfig []any
		for _, v := range item.IPAM.Config {
			ipamConfig = append(ipamConfig, map[string]any{
				"subnet":      v.Subnet,
				"gateway":     v.Gateway,
				"ip_range":    v.IPRange,
				"aux_address": v.AuxAddress,
			})
		}
		items = append(items, map[string]any{
			"id":         item.ID,
			"name":       item.Name,
			"driver":     item.Driver,
			"ipv6":       item.EnableIPv6,
			"scope":      item.Scope,
			"internal":   item.Internal,
			"attachable": item.Attachable,
			"ingress":    item.Ingress,
			"labels":     item.Labels,
			"options":    item.Options,
			"ipam": map[string]any{
				"config":  ipamConfig,
				"driver":  item.IPAM.Driver,
				"options": item.IPAM.Options,
			},
			"created": carbon.CreateFromStdTime(item.Created).ToDateTimeString(),
		})
	}

	Success(w, chix.M{
		"total": total,
		"items": items,
	})
}

func (s *ContainerNetworkService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerNetworkCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	id, err := s.containerNetworkRepo.Create(req)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, id)
}

func (s *ContainerNetworkService) Remove(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerNetworkID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.containerNetworkRepo.Remove(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *ContainerNetworkService) Exist(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerNetworkID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	exist, err := s.containerNetworkRepo.Exist(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, exist)
}

func (s *ContainerNetworkService) Inspect(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerNetworkID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	network, err := s.containerNetworkRepo.Inspect(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, network)
}

func (s *ContainerNetworkService) Connect(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerNetworkConnect](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.containerNetworkRepo.Connect(req.Network, req.Container); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *ContainerNetworkService) Disconnect(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerNetworkConnect](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.containerNetworkRepo.Disconnect(req.Network, req.Container); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *ContainerNetworkService) Prune(w http.ResponseWriter, r *http.Request) {
	if err := s.containerNetworkRepo.Prune(); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}
