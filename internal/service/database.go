package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
)

type Database struct {
	databaseRepo biz.DatabaseRepo
}

func NewDatabaseService() *Database {
	return &Database{
		databaseRepo: data.NewDatabaseRepo(),
	}
}

func (s *Database) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	certs, total, err := s.databaseRepo.List(req.Page, req.Limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, chix.M{
		"total": total,
		"items": certs,
	})
}

func (s *Database) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.DatabaseCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.databaseRepo.Create(req); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *Database) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.DatabaseDelete](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.databaseRepo.Delete(req.ServerID, req.Name); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}
