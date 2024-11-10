package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
)

type DatabaseServer struct {
	databaseServerRepo biz.DatabaseServerRepo
}

func NewDatabaseServerService() *DatabaseServer {
	return &DatabaseServer{
		databaseServerRepo: data.NewDatabaseServerRepo(),
	}
}

func (s *DatabaseServer) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	certs, total, err := s.databaseServerRepo.List(req.Page, req.Limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, chix.M{
		"total": total,
		"items": certs,
	})
}

func (s *DatabaseServer) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.DatabaseServerCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.databaseServerRepo.Create(req); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *DatabaseServer) Update(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.DatabaseServerUpdate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.databaseServerRepo.Update(req); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *DatabaseServer) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.databaseServerRepo.Delete(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}
