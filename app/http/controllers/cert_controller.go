package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	requests "github.com/TheTNB/panel/v2/app/http/requests/cert"
	commonrequests "github.com/TheTNB/panel/v2/app/http/requests/common"
	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/acme"
	"github.com/TheTNB/panel/v2/pkg/h"
)

type CertController struct {
	cron internal.Cron
	cert internal.Cert
}

func NewCertController() *CertController {
	return &CertController{
		cron: services.NewCronImpl(),
		cert: services.NewCertImpl(),
	}
}

// CAProviders
//
//	@Summary		获取 CA 提供商
//	@Description	获取面板证书管理支持的 CA 提供商
//	@Tags			TLS证书
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	SuccessResponse
//	@Router			/panel/cert/caProviders [get]
func (r *CertController) CAProviders(ctx http.Context) http.Response {
	return h.Success(ctx, []map[string]string{
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
//	@Description	获取面板证书管理支持的 DNS 提供商
//	@Tags			TLS证书
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	SuccessResponse
//	@Router			/panel/cert/dnsProviders [get]
func (r *CertController) DNSProviders(ctx http.Context) http.Response {
	return h.Success(ctx, []map[string]any{
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
//	@Description	获取面板证书管理支持的算法列表
//	@Tags			TLS证书
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	SuccessResponse
//	@Router			/panel/cert/algorithms [get]
func (r *CertController) Algorithms(ctx http.Context) http.Response {
	return h.Success(ctx, []map[string]any{
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

// UserList
//
//	@Summary		获取用户列表
//	@Description	获取面板证书管理的 ACME 用户列表
//	@Tags			TLS证书
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		commonrequests.Paginate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/users [get]
func (r *CertController) UserList(ctx http.Context) http.Response {
	var paginateRequest commonrequests.Paginate
	sanitize := h.SanitizeRequest(ctx, &paginateRequest)
	if sanitize != nil {
		return sanitize
	}

	var users []models.CertUser
	var total int64
	err := facades.Orm().Query().Paginate(paginateRequest.Page, paginateRequest.Limit, &users, &total)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Info("获取ACME用户列表失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, http.Json{
		"total": total,
		"items": users,
	})
}

// UserStore
//
//	@Summary		添加 ACME 用户
//	@Description	添加 ACME 用户到面板证书管理
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.UserStore	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/users [post]
func (r *CertController) UserStore(ctx http.Context) http.Response {
	var storeRequest requests.UserStore
	sanitize := h.SanitizeRequest(ctx, &storeRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.UserStore(storeRequest)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Info("添加ACME用户失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// UserUpdate
//
//	@Summary		更新 ACME 用户
//	@Description	更新面板证书管理的 ACME 用户
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			id		path		int					true	"用户 ID"
//	@Param			data	body		requests.UserUpdate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/users/{id} [put]
func (r *CertController) UserUpdate(ctx http.Context) http.Response {
	var updateRequest requests.UserUpdate
	sanitize := h.SanitizeRequest(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.UserUpdate(updateRequest)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"userID": updateRequest.ID,
			"error":  err.Error(),
		}).Info("更新ACME用户失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// UserShow
//
//	@Summary		获取 ACME 用户
//	@Description	获取面板证书管理的 ACME 用户
//	@Tags			TLS证书
//	@Produce		json
//	@Security		BearerToken
//	@Param			id	path		int	true	"用户 ID"
//	@Success		200	{object}	SuccessResponse{data=models.CertUser}
//	@Router			/panel/cert/users/{id} [get]
func (r *CertController) UserShow(ctx http.Context) http.Response {
	var showAndDestroyRequest requests.UserShowAndDestroy
	sanitize := h.SanitizeRequest(ctx, &showAndDestroyRequest)
	if sanitize != nil {
		return sanitize
	}

	user, err := r.cert.UserShow(showAndDestroyRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"userID": showAndDestroyRequest.ID,
			"error":  err.Error(),
		}).Info("获取ACME用户失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, user)
}

// UserDestroy
//
//	@Summary		删除 ACME 用户
//	@Description	删除面板证书管理的 ACME 用户
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			id	path		int	true	"用户 ID"
//	@Success		200	{object}	SuccessResponse
//	@Router			/panel/cert/users/{id} [delete]
func (r *CertController) UserDestroy(ctx http.Context) http.Response {
	var showAndDestroyRequest requests.UserShowAndDestroy
	sanitize := h.SanitizeRequest(ctx, &showAndDestroyRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.UserDestroy(showAndDestroyRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"userID": showAndDestroyRequest.ID,
			"error":  err.Error(),
		}).Info("删除ACME用户失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// DNSList
//
//	@Summary		获取 DNS 接口列表
//	@Description	获取面板证书管理的 DNS 接口列表
//	@Tags			TLS证书
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		commonrequests.Paginate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/dns [get]
func (r *CertController) DNSList(ctx http.Context) http.Response {
	var paginateRequest commonrequests.Paginate
	sanitize := h.SanitizeRequest(ctx, &paginateRequest)
	if sanitize != nil {
		return sanitize
	}

	var dns []models.CertDNS
	var total int64
	err := facades.Orm().Query().Paginate(paginateRequest.Page, paginateRequest.Limit, &dns, &total)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Info("获取DNS接口列表失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, http.Json{
		"total": total,
		"items": dns,
	})
}

// DNSStore
//
//	@Summary		添加 DNS 接口
//	@Description	添加 DNS 接口到面板证书管理
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.DNSStore	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/dns [post]
func (r *CertController) DNSStore(ctx http.Context) http.Response {
	var storeRequest requests.DNSStore
	sanitize := h.SanitizeRequest(ctx, &storeRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.DNSStore(storeRequest)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Info("添加DNS接口失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, nil)
}

// DNSShow
//
//	@Summary		获取 DNS 接口
//	@Description	获取面板证书管理的 DNS 接口
//	@Tags			TLS证书
//	@Produce		json
//	@Security		BearerToken
//	@Param			id	path		int	true	"DNS 接口 ID"
//	@Success		200	{object}	SuccessResponse{data=models.CertDNS}
//	@Router			/panel/cert/dns/{id} [get]
func (r *CertController) DNSShow(ctx http.Context) http.Response {
	var showAndDestroyRequest requests.DNSShowAndDestroy
	sanitize := h.SanitizeRequest(ctx, &showAndDestroyRequest)
	if sanitize != nil {
		return sanitize
	}

	dns, err := r.cert.DNSShow(showAndDestroyRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"dnsID": showAndDestroyRequest.ID,
			"error": err.Error(),
		}).Info("获取DNS接口失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, dns)
}

// DNSUpdate
//
//	@Summary		更新 DNS 接口
//	@Description	更新面板证书管理的 DNS 接口
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			id		path		int					true	"DNS 接口 ID"
//	@Param			data	body		requests.DNSUpdate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/dns/{id} [put]
func (r *CertController) DNSUpdate(ctx http.Context) http.Response {
	var updateRequest requests.DNSUpdate
	sanitize := h.SanitizeRequest(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.DNSUpdate(updateRequest)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"dnsID": updateRequest.ID,
			"error": err.Error(),
		}).Info("更新DNS接口失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, nil)
}

// DNSDestroy
//
//	@Summary		删除 DNS 接口
//	@Description	删除面板证书管理的 DNS 接口
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			id	path		int	true	"DNS 接口 ID"
//	@Success		200	{object}	SuccessResponse
//	@Router			/panel/cert/dns/{id} [delete]
func (r *CertController) DNSDestroy(ctx http.Context) http.Response {
	var showAndDestroyRequest requests.DNSShowAndDestroy
	sanitize := h.SanitizeRequest(ctx, &showAndDestroyRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.DNSDestroy(showAndDestroyRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"dnsID": showAndDestroyRequest.ID,
			"error": err.Error(),
		}).Info("删除DNS接口失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// CertList
//
//	@Summary		获取证书列表
//	@Description	获取面板证书管理的证书列表
//	@Tags			TLS证书
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		commonrequests.Paginate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/certs [get]
func (r *CertController) CertList(ctx http.Context) http.Response {
	var paginateRequest commonrequests.Paginate
	sanitize := h.SanitizeRequest(ctx, &paginateRequest)
	if sanitize != nil {
		return sanitize
	}

	var certs []models.Cert
	var total int64
	err := facades.Orm().Query().With("Website").With("User").With("DNS").Paginate(paginateRequest.Page, paginateRequest.Limit, &certs, &total)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Info("获取证书列表失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, http.Json{
		"total": total,
		"items": certs,
	})
}

// CertStore
//
//	@Summary		添加证书
//	@Description	添加证书到面板证书管理
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.CertStore	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/certs [post]
func (r *CertController) CertStore(ctx http.Context) http.Response {
	var storeRequest requests.CertStore
	sanitize := h.SanitizeRequest(ctx, &storeRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.CertStore(storeRequest)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Info("添加证书失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, nil)
}

// CertUpdate
//
//	@Summary		更新证书
//	@Description	更新面板证书管理的证书
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			id		path		int					true	"证书 ID"
//	@Param			data	body		requests.CertUpdate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/certs/{id} [put]
func (r *CertController) CertUpdate(ctx http.Context) http.Response {
	var updateRequest requests.CertUpdate
	sanitize := h.SanitizeRequest(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.CertUpdate(updateRequest)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"certID": updateRequest.ID,
			"error":  err.Error(),
		}).Info("更新证书失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, nil)
}

// CertShow
//
//	@Summary		获取证书
//	@Description	获取面板证书管理的证书
//	@Tags			TLS证书
//	@Produce		json
//	@Security		BearerToken
//	@Param			id	path		int	true	"证书 ID"
//	@Success		200	{object}	SuccessResponse{data=models.Cert}
//	@Router			/panel/cert/certs/{id} [get]
func (r *CertController) CertShow(ctx http.Context) http.Response {
	var showAndDestroyRequest requests.CertShowAndDestroy
	sanitize := h.SanitizeRequest(ctx, &showAndDestroyRequest)
	if sanitize != nil {
		return sanitize
	}

	cert, err := r.cert.CertShow(showAndDestroyRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"certID": showAndDestroyRequest.ID,
			"error":  err.Error(),
		}).Info("获取证书失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, cert)
}

// CertDestroy
//
//	@Summary		删除证书
//	@Description	删除面板证书管理的证书
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			id	path		int	true	"证书 ID"
//	@Success		200	{object}	SuccessResponse
//	@Router			/panel/cert/certs/{id} [delete]
func (r *CertController) CertDestroy(ctx http.Context) http.Response {
	var showAndDestroyRequest requests.CertShowAndDestroy
	sanitize := h.SanitizeRequest(ctx, &showAndDestroyRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.CertDestroy(showAndDestroyRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"certID": showAndDestroyRequest.ID,
			"error":  err.Error(),
		}).Info("删除证书失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, nil)
}

// Obtain
//
//	@Summary		签发证书
//	@Description	签发面板证书管理的证书
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Obtain	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/obtain [post]
func (r *CertController) Obtain(ctx http.Context) http.Response {
	var obtainRequest requests.Obtain
	sanitize := h.SanitizeRequest(ctx, &obtainRequest)
	if sanitize != nil {
		return sanitize
	}

	cert, err := r.cert.CertShow(obtainRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"certID": obtainRequest.ID,
			"error":  err.Error(),
		}).Info("获取证书失败")
		return h.ErrorSystem(ctx)
	}

	if cert.DNS != nil || cert.Website != nil {
		_, err = r.cert.ObtainAuto(obtainRequest.ID)
	} else {
		_, err = r.cert.ObtainManual(obtainRequest.ID)
	}
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Info("签发证书失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// Renew
//
//	@Summary		续签证书
//	@Description	续签面板证书管理的证书
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Renew	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/renew [post]
func (r *CertController) Renew(ctx http.Context) http.Response {
	var renewRequest requests.Renew
	sanitize := h.SanitizeRequest(ctx, &renewRequest)
	if sanitize != nil {
		return sanitize
	}

	_, err := r.cert.Renew(renewRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Info("续签证书失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ManualDNS
//
//	@Summary		获取手动 DNS 记录
//	@Description	获取签发证书所需的 DNS 记录
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Obtain	true	"request"
//	@Success		200		{object}	SuccessResponse{data=[]acme.DNSRecord}
//	@Router			/panel/cert/manualDNS [post]
func (r *CertController) ManualDNS(ctx http.Context) http.Response {
	var obtainRequest requests.Obtain
	sanitize := h.SanitizeRequest(ctx, &obtainRequest)
	if sanitize != nil {
		return sanitize
	}

	resolves, err := r.cert.ManualDNS(obtainRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Info("获取手动DNS记录失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, resolves)
}

// Deploy
//
//	@Summary		部署证书
//	@Description	部署面板证书管理的证书
//	@Tags			TLS证书
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.CertDeploy	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/cert/deploy [post]
func (r *CertController) Deploy(ctx http.Context) http.Response {
	var deployRequest requests.CertDeploy
	sanitize := h.SanitizeRequest(ctx, &deployRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.Deploy(deployRequest.ID, deployRequest.WebsiteID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"certID": deployRequest.ID,
			"error":  err.Error(),
		}).Info("部署证书失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}
