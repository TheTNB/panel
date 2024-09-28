package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/acme"
)

type CertService struct {
	certRepo biz.CertRepo
}

func NewCertService() *CertService {
	return &CertService{
		certRepo: data.NewCertRepo(),
	}
}

func (s *CertService) CAProviders(w http.ResponseWriter, r *http.Request) {
	Success(w, []map[string]string{
		{
			"name": "Let's Encrypt",
			"ca":   "letsencrypt",
		},
		{
			"name": "ZeroSSL",
			"ca":   "zerossl",
		},
		{
			"name": "SSL.com",
			"ca":   "sslcom",
		},
		{
			"name": "Google",
			"ca":   "google",
		},
		{
			"name": "Buypass",
			"ca":   "buypass",
		},
	})

}

func (s *CertService) DNSProviders(w http.ResponseWriter, r *http.Request) {
	Success(w, []map[string]any{
		{
			"name": "DNSPod",
			"dns":  acme.DnsPod,
		},
		{
			"name": "腾讯云",
			"dns":  acme.Tencent,
		},
		{
			"name": "阿里云",
			"dns":  acme.AliYun,
		},
		{
			"name": "CloudFlare",
			"dns":  acme.CloudFlare,
		},
	})

}

func (s *CertService) Algorithms(w http.ResponseWriter, r *http.Request) {
	Success(w, []map[string]any{
		{
			"name": "EC256",
			"key":  acme.KeyEC256,
		},
		{
			"name": "EC384",
			"key":  acme.KeyEC384,
		},
		{
			"name": "RSA2048",
			"key":  acme.KeyRSA2048,
		},
		{
			"name": "RSA4096",
			"key":  acme.KeyRSA4096,
		},
	})

}

func (s *CertService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	certs, total, err := s.certRepo.List(req.Page, req.Limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, chix.M{
		"total": total,
		"items": certs,
	})
}

func (s *CertService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.CertCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	cert, err := s.certRepo.Create(req)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, cert)
}

func (s *CertService) Update(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.CertUpdate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.certRepo.Update(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *CertService) Get(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	cert, err := s.certRepo.Get(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, cert)
}

func (s *CertService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = s.certRepo.Delete(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *CertService) Obtain(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	cert, err := s.certRepo.Get(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if cert.DNS != nil || cert.Website != nil {
		_, err = s.certRepo.ObtainAuto(req.ID)
	} else {
		_, err = s.certRepo.ObtainManual(req.ID)
	}

	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *CertService) Renew(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, err = s.certRepo.Renew(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *CertService) ManualDNS(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	dns, err := s.certRepo.ManualDNS(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, dns)
}

func (s *CertService) Deploy(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.CertDeploy](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = s.certRepo.Deploy(req.ID, req.WebsiteID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}
