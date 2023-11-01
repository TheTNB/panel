package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	requests "panel/app/http/requests/cert"

	commonrequests "panel/app/http/requests/common"
	responses "panel/app/http/responses/cert"
	"panel/app/models"
	"panel/app/services"
	"panel/pkg/acme"
)

type CertController struct {
	cron services.Cron
	cert services.Cert
}

func NewCertController() *CertController {
	return &CertController{
		cron: services.NewCronImpl(),
		cert: services.NewCertImpl(),
	}
}

// CAProviders
// @Summary 获取 CA 提供商
// @Description 获取面板证书管理支持的 CA 提供商
// @Tags 证书
// @Produce json
// @Security BearerToken
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Router /panel/cert/caProviders [get]
func (r *CertController) CAProviders(ctx http.Context) http.Response {
	return Success(ctx, []map[string]string{
		{
			"name": "Let's Encrypt",
			"ca":   acme.CALetEncrypt,
		},
		{
			"name": "ZeroSSL",
			"ca":   acme.CAZeroSSL,
		},
		{
			"name": "SSL.com",
			"ca":   acme.CASSLcom,
		},
		{
			"name": "Google",
			"ca":   acme.CAGoogle,
		},
		{
			"name": "Buypass",
			"ca":   acme.CABuypass,
		},
	})
}

// DNSProviders
// @Summary 获取 DNS 提供商
// @Description 获取面板证书管理支持的 DNS 提供商
// @Tags 证书
// @Produce json
// @Security BearerToken
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Router /panel/cert/dnsProviders [get]
func (r *CertController) DNSProviders(ctx http.Context) http.Response {
	return Success(ctx, []map[string]any{
		{
			"name": "DNSPod",
			"dns":  acme.DnsPod,
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
// @Summary 获取算法列表
// @Description 获取面板证书管理支持的算法列表
// @Tags 证书
// @Produce json
// @Security BearerToken
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Router /panel/cert/algorithms [get]
func (r *CertController) Algorithms(ctx http.Context) http.Response {
	return Success(ctx, []map[string]any{
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
// @Summary 获取用户列表
// @Description 获取面板证书管理的 ACME 用户列表
// @Tags 证书
// @Produce json
// @Security BearerToken
// @Success 200 {object} SuccessResponse{data=responses.CertList}
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/users [get]
func (r *CertController) UserList(ctx http.Context) http.Response {
	var updateProfileRequest commonrequests.Paginate
	sanitize := Sanitize(ctx, &updateProfileRequest)
	if sanitize != nil {
		return sanitize
	}

	var users []models.CertUser
	var total int64
	err := facades.Orm().Query().Paginate(updateProfileRequest.Page, updateProfileRequest.Limit, &users, &total)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Error("获取ACME用户列表失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, &responses.UserList{
		Total: total,
		Items: users,
	})
}

// UserAdd
// @Summary 添加 ACME 用户
// @Description 添加 ACME 用户到面板证书管理
// @Tags 证书
// @Accept json
// @Produce json
// @Security BearerToken
// @Param data body requests.UserAdd true "用户信息"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/users [post]
func (r *CertController) UserAdd(ctx http.Context) http.Response {
	var addRequest requests.UserAdd
	sanitize := Sanitize(ctx, &addRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.UserAdd(addRequest)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Error("添加ACME用户失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, nil)
}

// UserDelete
// @Summary 删除 ACME 用户
// @Description 删除面板证书管理的 ACME 用户
// @Tags 证书
// @Accept json
// @Produce json
// @Security BearerToken
// @Param id path int true "用户 ID"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/users/{id} [delete]
func (r *CertController) UserDelete(ctx http.Context) http.Response {
	userID := ctx.Request().InputInt("id")

	err := r.cert.UserDelete(uint(userID))
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"userID": userID,
			"error":  err.Error(),
		}).Error("删除ACME用户失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, nil)
}

// DNSList
// @Summary 获取 DNS 接口列表
// @Description 获取面板证书管理的 DNS 接口列表
// @Tags 证书
// @Produce json
// @Security BearerToken
// @Success 200 {object} SuccessResponse{data=responses.DNSList}
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/dns [get]
func (r *CertController) DNSList(ctx http.Context) http.Response {
	var updateProfileRequest commonrequests.Paginate
	sanitize := Sanitize(ctx, &updateProfileRequest)
	if sanitize != nil {
		return sanitize
	}

	var dns []models.CertDNS
	var total int64
	err := facades.Orm().Query().Paginate(updateProfileRequest.Page, updateProfileRequest.Limit, &dns, &total)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Error("获取DNS接口列表失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, &responses.DNSList{
		Total: total,
		Items: dns,
	})
}

// DNSAdd
// @Summary 添加 DNS 接口
// @Description 添加 DNS 接口到面板证书管理
// @Tags 证书
// @Accept json
// @Produce json
// @Security BearerToken
// @Param data body requests.DNSAdd true "DNS 接口信息"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/dns [post]
func (r *CertController) DNSAdd(ctx http.Context) http.Response {
	var addRequest requests.DNSAdd
	sanitize := Sanitize(ctx, &addRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.DNSAdd(addRequest)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Error("添加DNS接口失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, nil)
}

// DNSDelete
// @Summary 删除 DNS 接口
// @Description 删除面板证书管理的 DNS 接口
// @Tags 证书
// @Accept json
// @Produce json
// @Security BearerToken
// @Param id path int true "DNS 接口 ID"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/dns/{id} [delete]
func (r *CertController) DNSDelete(ctx http.Context) http.Response {
	dnsID := ctx.Request().InputInt("id")

	err := r.cert.DNSDelete(uint(dnsID))
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"dnsID": dnsID,
			"error": err.Error(),
		}).Error("删除DNS接口失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, nil)
}

// CertList
// @Summary 获取证书列表
// @Description 获取面板证书管理的证书列表
// @Tags 证书
// @Produce json
// @Security BearerToken
// @Success 200 {object} SuccessResponse{data=responses.CertList}
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/certs [get]
func (r *CertController) CertList(ctx http.Context) http.Response {
	var updateProfileRequest commonrequests.Paginate
	sanitize := Sanitize(ctx, &updateProfileRequest)
	if sanitize != nil {
		return sanitize
	}

	var certs []models.Cert
	var total int64
	err := facades.Orm().Query().Paginate(updateProfileRequest.Page, updateProfileRequest.Limit, &certs, &total)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Error("获取证书列表失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, &responses.CertList{
		Total: total,
		Items: certs,
	})
}

// CertAdd
// @Summary 添加证书
// @Description 添加证书到面板证书管理
// @Tags 证书
// @Accept json
// @Produce json
// @Security BearerToken
// @Param data body requests.CertAdd true "证书信息"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/certs [post]
func (r *CertController) CertAdd(ctx http.Context) http.Response {
	var addRequest requests.CertAdd
	sanitize := Sanitize(ctx, &addRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.cert.CertAdd(addRequest)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Error("添加证书失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, nil)
}

// CertDelete
// @Summary 删除证书
// @Description 删除面板证书管理的证书
// @Tags 证书
// @Accept json
// @Produce json
// @Security BearerToken
// @Param id path int true "证书 ID"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/certs/{id} [delete]
func (r *CertController) CertDelete(ctx http.Context) http.Response {
	certID := ctx.Request().InputInt("id")

	err := r.cert.CertDelete(uint(certID))
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"certID": certID,
			"error":  err.Error(),
		}).Error("删除证书失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, nil)
}

// Obtain
// @Summary 签发证书
// @Description 签发面板证书管理的证书
// @Tags 证书
// @Accept json
// @Produce json
// @Security BearerToken
// @Param data body requests.Obtain true "证书信息"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/obtain [post]
func (r *CertController) Obtain(ctx http.Context) http.Response {
	var obtainRequest requests.Obtain
	sanitize := Sanitize(ctx, &obtainRequest)
	if sanitize != nil {
		return sanitize
	}

	cert, err := r.cert.GetByID(obtainRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"certID": obtainRequest.ID,
			"error":  err.Error(),
		}).Error("获取证书失败")
		return ErrorSystem(ctx)
	}

	if cert.DNS != nil || cert.Website != nil {
		_, err = r.cert.ObtainAuto(obtainRequest.ID)
	} else {
		_, err = r.cert.ObtainManual(obtainRequest.ID)
	}
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Error("签发证书失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, nil)
}

// Renew
// @Summary 续签证书
// @Description 续签面板证书管理的证书
// @Tags 证书
// @Accept json
// @Produce json
// @Security BearerToken
// @Param data body requests.Renew true "证书信息"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/renew [post]
func (r *CertController) Renew(ctx http.Context) http.Response {
	var renewRequest requests.Renew
	sanitize := Sanitize(ctx, &renewRequest)
	if sanitize != nil {
		return sanitize
	}

	_, err := r.cert.Renew(renewRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Error("续签证书失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, nil)
}

// ManualDNS
// @Summary 获取手动 DNS 记录
// @Description 获取签发证书所需的 DNS 记录
// @Tags 证书
// @Accept json
// @Produce json
// @Security BearerToken
// @Param data body requests.Obtain true "证书信息"
// @Success 200 {object} SuccessResponse{data=map[string]acme.Resolve}
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/cert/manualDNS [post]
func (r *CertController) ManualDNS(ctx http.Context) http.Response {
	var obtainRequest requests.Obtain
	sanitize := Sanitize(ctx, &obtainRequest)
	if sanitize != nil {
		return sanitize
	}

	resolves, err := r.cert.ManualDNS(obtainRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "证书管理").With(map[string]any{
			"error": err.Error(),
		}).Error("获取手动DNS记录失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, resolves)
}
