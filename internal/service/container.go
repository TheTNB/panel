package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
)

type ContainerService struct {
	containerRepo biz.ContainerRepo
}

func NewContainerService() *ContainerService {
	return &ContainerService{
		containerRepo: data.NewContainerRepo(),
	}
}

func (s *ContainerService) List(w http.ResponseWriter, r *http.Request) {
	containers, err := s.containerRepo.ListAll()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
	}

	paged, total := Paginate(r, containers)
	items := make([]any, 0)
	for _, item := range paged {
		items = append(items, map[string]any{
			"id":         item.ID,
			"name":       item.Name,
			"image":      item.Image,
			"command":    item.Command,
			"created_at": item.CreatedAt,
			"ports":      item.Ports,
			"labels":     item.Labels,
			"state":      item.State,
			"status":     item.Status,
		})
	}

	Success(w, chix.M{
		"total": total,
		"items": items,
	})
}

func (s *ContainerService) Search(w http.ResponseWriter, r *http.Request) {
	containers, err := s.containerRepo.ListByName(r.FormValue("name"))
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, chix.M{
		"total": len(containers),
		"items": containers,
	})
}

func (s *ContainerService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	id, err := s.containerRepo.Create(req)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, id)
}

func (s *ContainerService) Remove(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerRepo.Remove(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerService) Start(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerRepo.Start(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerService) Stop(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerRepo.Stop(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerService) Restart(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerRepo.Restart(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerService) Pause(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerRepo.Pause(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerService) Unpause(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerRepo.Unpause(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerService) Kill(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerRepo.Kill(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerService) Rename(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerRename](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerRepo.Rename(req.ID, req.Name); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerService) Logs(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	logs, err := s.containerRepo.Logs(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, logs)
}

func (s *ContainerService) Prune(w http.ResponseWriter, r *http.Request) {
	if err := s.containerRepo.Prune(); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}
