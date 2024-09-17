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

type ContainerImageService struct {
	containerImageRepo biz.ContainerImageRepo
}

func NewContainerImageService() *ContainerImageService {
	return &ContainerImageService{
		containerImageRepo: data.NewContainerImageRepo(),
	}
}

func (s *ContainerImageService) List(w http.ResponseWriter, r *http.Request) {
	images, err := s.containerImageRepo.List()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	paged, total := Paginate(r, images)

	items := make([]any, len(paged))
	for _, item := range paged {
		items = append(items, map[string]any{
			"id":           item.ID,
			"created":      carbon.CreateFromTimestamp(item.Created).ToDateTimeString(),
			"containers":   item.Containers,
			"size":         str.FormatBytes(float64(item.Size)),
			"labels":       item.Labels,
			"repo_tags":    item.RepoTags,
			"repo_digests": item.RepoDigests,
		})
	}

	Success(w, chix.M{
		"total": total,
		"items": items,
	})
}

func (s *ContainerImageService) Exist(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerImageID](r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	exist, err := s.containerImageRepo.Exist(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, exist)
}

func (s *ContainerImageService) Pull(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerImagePull](r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = s.containerImageRepo.Pull(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *ContainerImageService) Remove(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerImageID](r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = s.containerImageRepo.Remove(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *ContainerImageService) Inspect(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerImageID](r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	inspect, err := s.containerImageRepo.Inspect(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, inspect)
}

func (s *ContainerImageService) Prune(w http.ResponseWriter, r *http.Request) {
	if err := s.containerImageRepo.Prune(); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}
