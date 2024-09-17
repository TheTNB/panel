package service

import (
	"net/http"

	"github.com/go-rat/chix"
	"github.com/golang-module/carbon/v2"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/str"
)

type ContainerVolumeService struct {
	containerVolumeRepo biz.ContainerVolumeRepo
}

func NewContainerVolumeService() *ContainerVolumeService {
	return &ContainerVolumeService{
		containerVolumeRepo: data.NewContainerVolumeRepo(),
	}
}

func (s *ContainerVolumeService) List(w http.ResponseWriter, r *http.Request) {
	volumes, err := s.containerVolumeRepo.List()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	paged, total := Paginate(r, volumes)

	items := make([]any, len(paged))
	for _, item := range paged {
		var usage any
		if item.UsageData != nil {
			usage = map[string]any{
				"ref_count": item.UsageData.RefCount,
				"size":      str.FormatBytes(float64(item.UsageData.Size)),
			}
		}
		items = append(items, map[string]any{
			"id":      item.Name,
			"created": carbon.Parse(item.CreatedAt).ToDateTimeString(),
			"driver":  item.Driver,
			"mount":   item.Mountpoint,
			"labels":  item.Labels,
			"options": item.Options,
			"scope":   item.Scope,
			"status":  item.Status,
			"usage":   usage,
		})
	}

	Success(w, chix.M{
		"total": total,
		"items": items,
	})
}

func (s *ContainerVolumeService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerVolumeCreate](r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	volume, err := s.containerVolumeRepo.Create(req)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, volume.Name)

}

func (s *ContainerVolumeService) Exist(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerVolumeID](r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	exist, err := s.containerVolumeRepo.Exist(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, exist)
}

func (s *ContainerVolumeService) Remove(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerVolumeID](r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = s.containerVolumeRepo.Remove(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *ContainerVolumeService) Inspect(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerVolumeID](r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	volume, err := s.containerVolumeRepo.Inspect(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, volume)
}

func (s *ContainerVolumeService) Prune(w http.ResponseWriter, r *http.Request) {
	if err := s.containerVolumeRepo.Prune(); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}
