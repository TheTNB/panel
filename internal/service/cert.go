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

// CAProviders
//
//	@Summary		获取 CA 提供商
//	@Tags			证书服务
//	@Produce		json
//	@Success		200	{object}	SuccessResponse
//	@Router			/cert/caProviders	[get]
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

// DNSProviders
//
//	@Summary		获取 DNS 提供商
//	@Tags			证书服务
//	@Produce		json
//	@Success		200	{object}	SuccessResponse
//	@Router			/cert/dnsProviders	[get]
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

// Algorithms
//
//	@Summary		获取算法列表
//	@Tags			证书服务
//	@Produce		json
//	@Success		200	{object}	SuccessResponse
//	@Router			/cert/algorithms	[get]
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

// List
//
//	@Summary		证书列表
//	@Tags			证书服务
//	@Produce		json
//	@Success		200		{object}	SuccessResponse
//	@Router			/cert/cert [get]
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

// Create
//
//	@Summary		创建证书
//	@Tags			证书服务
//	@Produce		json
//	@Success		200		{object}	SuccessResponse
//	@Router			/cert/cert [post]
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

// Update
//
//	@Summary		更新证书
//	@Tags			证书服务
//	@Produce		json
//	@Param			id		path		int					true	"证书 ID"
//	@Success		200		{object}	SuccessResponse
//	@Router			/cert/cert/{id} [post]
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

// Get
//
//	@Summary		获取证书
//	@Tags			证书服务
//	@Produce		json
//	@Param			id		path		int					true	"证书 ID"
//	@Success		200		{object}	SuccessResponse
//	@Router			/cert/cert/{id} [get]
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

// Delete
//
//	@Summary		删除证书
//	@Tags			证书服务
//	@Produce		json
//	@Param			id		path		int					true	"证书 ID"
//	@Success		200		{object}	SuccessResponse
//	@Router			/cert/cert/{id} [delete]
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

// Obtain
//
//	@Summary		签发证书
//	@Tags			证书服务
//	@Produce		json
//	@Param			id		path		int					true	"证书 ID"
//	@Success		200		{object}	SuccessResponse
//	@Router			/cert/{id}/obtain [post]
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

// Renew
//
//	@Summary		续签证书
//	@Tags			证书服务
//	@Produce		json
//	@Param			id		path		int					true	"证书 ID"
//	@Success		200		{object}	SuccessResponse
//	@Router			/cert/{id}/renew [post]
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

// ManualDNS
//
//	@Summary		手动 DNS
//	@Tags			证书服务
//	@Produce		json
//	@Param			id		path		int					true	"证书 ID"
//	@Success		200		{object}	SuccessResponse
//	@Router			/cert/{id}/manualDNS [post]
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

// Deploy
//
//	@Summary		部署证书
//	@Tags			证书服务
//	@Produce		json
//	@Param			id			path		int					true	"证书 ID"
//	@Param			websiteID	query		int					true	"网站 ID"
//	@Success		200		{object}	SuccessResponse
//	@Router			/cert/{id}/deploy [post]
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
