package acme

import (
	"time"

	"github.com/go-acme/lego/v4/acme"
	"github.com/go-acme/lego/v4/acme/api"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/providers/dns/dnspod"
	"github.com/go-acme/lego/v4/providers/http/webroot"
)

type Client struct {
	Config *lego.Config
	Client *lego.Client
	User   *User
}

type DnsType string

const (
	DnsPod     DnsType = "dnspod"
	AliYun     DnsType = "aliyun"
	CloudFlare DnsType = "cloudflare"
)

type DNSParam struct {
	ID        string `json:"id"`
	Token     string `json:"token"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Email     string `json:"email"`
	APIkey    string `json:"api_key"`
}

// UseDns 使用 DNS 接口验证
func (c *Client) UseDns(dnsType DnsType, param DNSParam) error {
	var p challenge.Provider
	var err error
	if dnsType == DnsPod {
		dnsPodConfig := dnspod.NewDefaultConfig()
		dnsPodConfig.LoginToken = param.ID + "," + param.Token
		p, err = dnspod.NewDNSProviderConfig(dnsPodConfig)
		if err != nil {
			return err
		}
	}
	if dnsType == AliYun {
		aliyunConfig := alidns.NewDefaultConfig()
		aliyunConfig.SecretKey = param.SecretKey
		aliyunConfig.APIKey = param.AccessKey
		p, err = alidns.NewDNSProviderConfig(aliyunConfig)
		if err != nil {
			return err
		}
	}
	if dnsType == CloudFlare {
		cloudflareConfig := cloudflare.NewDefaultConfig()
		cloudflareConfig.AuthEmail = param.Email
		cloudflareConfig.AuthKey = param.APIkey
		p, err = cloudflare.NewDNSProviderConfig(cloudflareConfig)
		if err != nil {
			return err
		}
	}

	return c.Client.Challenge.SetDNS01Provider(p, dns01.AddDNSTimeout(3*time.Minute))
}

// UseManualDns 使用手动 DNS 验证
func (c *Client) UseManualDns(checkDns ...bool) error {
	p := &manualDnsProvider{}
	var err error

	if len(checkDns) > 0 && !checkDns[0] {
		err = c.Client.Challenge.SetDNS01Provider(p, dns01.DisableCompletePropagationRequirement())
	} else {
		err = c.Client.Challenge.SetDNS01Provider(p, dns01.AddDNSTimeout(3*time.Minute))
	}

	return err
}

// UseHTTP 使用 HTTP 验证
func (c *Client) UseHTTP(path string) error {
	httpProvider, err := webroot.NewHTTPProvider(path)
	if err != nil {
		return err
	}

	err = c.Client.Challenge.SetHTTP01Provider(httpProvider)
	if err != nil {
		return err
	}
	return nil
}

// ObtainSSL 签发 SSL 证书
func (c *Client) ObtainSSL(domains []string) (certificate.Resource, error) {
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}

	certificates, err := c.Client.Certificate.Obtain(request)
	if err != nil {
		return certificate.Resource{}, err
	}

	return *certificates, nil
}

// RenewSSL 续签 SSL 证书
func (c *Client) RenewSSL(certUrl string) (certificate.Resource, error) {
	certificates, err := c.Client.Certificate.Get(certUrl, true)
	if err != nil {
		return certificate.Resource{}, err
	}

	certificates, err = c.Client.Certificate.RenewWithOptions(*certificates, &certificate.RenewOptions{
		Bundle:     true,
		MustStaple: true,
	})
	if err != nil {
		return certificate.Resource{}, err
	}

	return *certificates, nil
}

// GetDNSResolve 获取 DNS 解析（手动设置）
func (c *Client) GetDNSResolve(domains []string) (map[string]Resolve, error) {
	core, err := api.New(c.Config.HTTPClient, c.Config.UserAgent, c.Config.CADirURL, c.User.Registration.URI, c.User.Key)
	if err != nil {
		return nil, err
	}
	order, err := core.Orders.New(domains)
	if err != nil {
		return nil, err
	}
	resolves := make(map[string]Resolve)
	resChan, errChan := make(chan acme.Authorization), make(chan domainError)
	for _, authzURL := range order.Authorizations {
		go func(authzURL string) {
			authz, err := core.Authorizations.Get(authzURL)
			if err != nil {
				errChan <- domainError{Domain: authz.Identifier.Value, Error: err}
				return
			}
			resChan <- authz
		}(authzURL)
	}

	var responses []acme.Authorization
	for i := 0; i < len(order.Authorizations); i++ {
		select {
		case res := <-resChan:
			responses = append(responses, res)
		case err := <-errChan:
			resolves[err.Domain] = Resolve{Err: err.Error.Error()}
		}
	}
	close(resChan)
	close(errChan)

	for _, auth := range responses {
		domain := challenge.GetTargetedDomain(auth)
		acmeChallenge, err := challenge.FindChallenge(challenge.DNS01, auth)
		if err != nil {
			resolves[domain] = Resolve{Err: err.Error()}
			continue
		}
		keyAuth, err := core.GetKeyAuthorization(acmeChallenge.Token)
		if err != nil {
			resolves[domain] = Resolve{Err: err.Error()}
			continue
		}
		challengeInfo := dns01.GetChallengeInfo(domain, keyAuth)
		resolves[domain] = Resolve{
			Key:   challengeInfo.FQDN,
			Value: challengeInfo.Value,
		}
	}

	return resolves, nil
}
