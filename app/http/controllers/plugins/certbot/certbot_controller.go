package redis

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/json"

	"panel/app/http/controllers"
	"panel/app/services"
	"panel/pkg/acme"
)

type CertbotController struct {
	setting services.Setting
}

type User struct {
	Email       string
	CA          string // CA 提供商 (letsencrypt, zerossl, sslcom, google, buypass)
	Kid         string
	HmacEncoded string
	PrivateKey  string
}

type DNS struct {
	Type string // DNS 类型 (dnspod, aliyun, cloudflare)
	acme.DNSParam
}

type Cert struct {
	ID      int64    `json:"id"`
	CronID  int64    `json:"cron_id"`
	Type    string   `json:"type"`
	Domains []string `json:"domains"`
}

func NewCertbotController() *CertbotController {
	return &CertbotController{
		setting: services.NewSettingImpl(),
	}
}

// CAProviders 获取 CA 提供商
func (c *CertbotController) CAProviders(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "certbot")
	if check != nil {
		return check
	}

	return controllers.Success(ctx, []map[string]string{
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

// Algorithms 获取算法列表
func (c *CertbotController) Algorithms(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "certbot")
	if check != nil {
		return check
	}

	return controllers.Success(ctx, []map[string]any{
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

// UserList 获取用户列表
func (c *CertbotController) UserList(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "certbot")
	if check != nil {
		return check
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)

	var userList []User
	err := json.UnmarshalString(c.setting.Get("certbot_user", "[]"), &userList)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "获取证书列表失败")
	}

	startIndex := (page - 1) * limit
	endIndex := page * limit
	if startIndex > len(userList) {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []User{},
		})
	}
	if endIndex > len(userList) {
		endIndex = len(userList)
	}
	pagedCertList := userList[startIndex:endIndex]
	if pagedCertList == nil {
		pagedCertList = []User{}
	}

	return controllers.Success(ctx, http.Json{
		"total": len(userList),
		"items": pagedCertList,
	})
}

// DNSList 获取 DNS 列表
func (c *CertbotController) DNSList(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "certbot")
	if check != nil {
		return check
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)

	var dnsList []DNS
	err := json.UnmarshalString(c.setting.Get("certbot_dns", "[]"), &dnsList)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "获取证书列表失败")
	}

	startIndex := (page - 1) * limit
	endIndex := page * limit
	if startIndex > len(dnsList) {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []DNS{},
		})
	}
	if endIndex > len(dnsList) {
		endIndex = len(dnsList)
	}
	pagedCertList := dnsList[startIndex:endIndex]
	if pagedCertList == nil {
		pagedCertList = []DNS{}
	}

	return controllers.Success(ctx, http.Json{
		"total": len(dnsList),
		"items": pagedCertList,
	})
}

// CertList 所有证书
func (c *CertbotController) CertList(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "certbot")
	if check != nil {
		return check
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)

	var certList []Cert
	err := json.UnmarshalString(c.setting.Get("certbot_cert", "[]"), &certList)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "获取证书列表失败")
	}

	startIndex := (page - 1) * limit
	endIndex := page * limit
	if startIndex > len(certList) {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []Cert{},
		})
	}
	if endIndex > len(certList) {
		endIndex = len(certList)
	}
	pagedCertList := certList[startIndex:endIndex]
	if pagedCertList == nil {
		pagedCertList = []Cert{}
	}

	return controllers.Success(ctx, http.Json{
		"total": len(certList),
		"items": pagedCertList,
	})
}

// CertAdd 添加证书
func (c *CertbotController) CertAdd(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "certbot")
	if check != nil {
		return check
	}

	var certList []Cert
	err := json.UnmarshalString(c.setting.Get("certbot_cert", "[]"), &certList)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "获取证书列表失败")
	}

	var cert Cert
	cert.ID = carbon.Now().TimestampMilli()
	certList = append(certList, cert)
	encoded, err := json.MarshalString(certList)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "添加证书失败")
	}
	err = c.setting.Set("certbot", encoded)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "添加证书失败")
	}

	return controllers.Success(ctx, nil)
}

// CertDelete 删除证书
func (c *CertbotController) CertDelete(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "certbot")
	if check != nil {
		return check
	}

	var certList []Cert
	err := json.UnmarshalString(c.setting.Get("certbot_cert", "[]"), &certList)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "获取证书列表失败")
	}

	var cert Cert
	cert.ID = ctx.Request().InputInt64("id")
	for i, item := range certList {
		if item.ID == cert.ID {
			certList = append(certList[:i], certList[i+1:]...)
			break
		}
	}
	encoded, err := json.MarshalString(certList)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "删除证书失败")
	}
	err = c.setting.Set("certbot", encoded)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "删除证书失败")
	}

	return controllers.Success(ctx, nil)
}
