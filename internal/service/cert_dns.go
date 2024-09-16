package service

import (
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/go-rat/chix"
	"net/http"
)

type CertDNSService struct {
	certDNSRepo biz.CertDNSRepo
}

func NewCertDNSService() *CertDNSService {
	return &CertDNSService{
		certDNSRepo: data.NewCertDNSRepo(),
	}
}

func (s *CertDNSService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	certDNS, total, err := s.certDNSRepo.List(req.Page, req.Limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, chix.M{
		"total": total,
		"items": certDNS,
	})
}

func (s *CertDNSService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.CertDNSCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	certDNS, err := s.certDNSRepo.Create(req)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, certDNS)
}

func (s *CertDNSService) Update(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.CertDNSUpdate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.certDNSRepo.Update(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *CertDNSService) Get(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	certDNS, err := s.certDNSRepo.Get(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, certDNS)
}

func (s *CertDNSService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = s.certDNSRepo.Delete(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}
