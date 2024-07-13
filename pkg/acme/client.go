package acme

import (
	"context"
	"sort"

	"github.com/libdns/libdns"
	"github.com/mholt/acmez/v2"
	"github.com/mholt/acmez/v2/acme"

	"github.com/TheTNB/panel/v2/pkg/cert"
)

type Certificate struct {
	PrivateKey []byte
	acme.Certificate
}

type Client struct {
	Account acme.Account
	zClient acmez.Client
	// 手动 DNS 所需的信号通道
	manualDNSSolver
}

// UseDns 使用 DNS 接口验证
func (c *Client) UseDns(dnsType DnsType, param DNSParam) {
	c.zClient.ChallengeSolvers = map[string]acmez.Solver{
		acme.ChallengeTypeDNS01: dnsSolver{
			dns:     dnsType,
			param:   param,
			records: &[]libdns.Record{},
		},
	}
}

// UseManualDns 使用手动 DNS 验证
func (c *Client) UseManualDns(total int, check ...bool) {
	c.controlChan = make(chan struct{})
	c.dataChan = make(chan any)
	c.zClient.ChallengeSolvers = map[string]acmez.Solver{
		acme.ChallengeTypeDNS01: manualDNSSolver{
			check:       len(check) > 0 && check[0],
			controlChan: c.controlChan,
			dataChan:    c.dataChan,
			records:     &[]DNSRecord{},
		},
	}
}

// UseHTTP 使用 HTTP 验证
// conf openresty 配置文件路径
// path 验证文件存放路径
func (c *Client) UseHTTP(conf, path string) {
	c.zClient.ChallengeSolvers = map[string]acmez.Solver{
		acme.ChallengeTypeHTTP01: httpSolver{
			conf: conf,
			path: path,
		},
	}
}

// ObtainSSL 签发 SSL 证书
func (c *Client) ObtainSSL(ctx context.Context, domains []string, keyType KeyType) (Certificate, error) {
	certPrivateKey, err := generatePrivateKey(keyType)
	if err != nil {
		return Certificate{}, err
	}
	pemPrivateKey, err := cert.EncodeKey(certPrivateKey)
	if err != nil {
		return Certificate{}, err
	}

	certs, err := c.zClient.ObtainCertificateForSANs(ctx, c.Account, certPrivateKey, domains)
	if err != nil {
		return Certificate{}, err
	}

	cert := c.selectPreferredChain(certs)
	return Certificate{PrivateKey: pemPrivateKey, Certificate: cert}, nil
}

// ObtainSSLManual 手动验证 SSL 证书
func (c *Client) ObtainSSLManual() (Certificate, error) {
	// 发送信号，开始验证
	c.controlChan <- struct{}{}
	// 等待验证完成
	data := <-c.dataChan

	if err, ok := data.(error); ok {
		return Certificate{}, err
	}

	return data.(Certificate), nil
}

// RenewSSL 续签 SSL 证书
func (c *Client) RenewSSL(ctx context.Context, certUrl string, domains []string, keyType KeyType) (Certificate, error) {
	_, err := c.zClient.GetCertificateChain(ctx, c.Account, certUrl)
	if err != nil {
		return Certificate{}, err
	}

	return c.ObtainSSL(ctx, domains, keyType)
}

// GetDNSRecords 获取 DNS 解析（手动设置）
func (c *Client) GetDNSRecords(ctx context.Context, domains []string, keyType KeyType) ([]DNSRecord, error) {
	go func(ctx context.Context, domains []string, keyType KeyType) {
		certs, err := c.ObtainSSL(ctx, domains, keyType)
		// 将证书和错误信息发送到 dataChan
		if err != nil {
			c.dataChan <- err
			return
		}
		c.dataChan <- certs
	}(ctx, domains, keyType)

	// 这里要少一次循环，因为需要卡住最后一次的 dataChan，等待手动 DNS 验证完成
	for i := 1; i < len(domains); i++ {
		<-c.dataChan
		c.controlChan <- struct{}{}
	}

	// 因为上面少了一次循环，所以这里接收到的即为完整的 DNS 记录切片
	data := <-c.dataChan
	if err, ok := data.(error); ok {
		return nil, err
	}

	return data.([]DNSRecord), nil
}

func (c *Client) selectPreferredChain(certChains []acme.Certificate) acme.Certificate {
	if len(certChains) == 1 {
		return certChains[0]
	}

	sort.Slice(certChains, func(i, j int) bool {
		return len(certChains[i].ChainPEM) < len(certChains[j].ChainPEM)
	})

	return certChains[0]
}
