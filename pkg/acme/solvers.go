package acme

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/libdns/alidns"
	"github.com/libdns/cloudflare"
	"github.com/libdns/dnspod"
	"github.com/libdns/libdns"
	"github.com/mholt/acmez/v2/acme"
	"golang.org/x/net/publicsuffix"
)

type httpSolver struct {
	path string
}

func (s httpSolver) Present(_ context.Context, challenge acme.Challenge) error {
	var err error
	if s.path == "" {
		return nil
	}

	challengeFilePath := filepath.Join(s.path, challenge.HTTP01ResourcePath())
	if err = os.MkdirAll(filepath.Dir(challengeFilePath), 0o755); err != nil {
		return fmt.Errorf("无法在网站目录创建 HTTP 挑战所需的目录: %w", err)
	}

	if err = os.WriteFile(challengeFilePath, []byte(challenge.KeyAuthorization), 0o644); err != nil {
		return fmt.Errorf("无法在网站目录创建 HTTP 挑战所需的文件: %w", err)
	}

	return nil
}

// CleanUp cleans up the HTTP server if it is the last one to finish.
func (s httpSolver) CleanUp(_ context.Context, challenge acme.Challenge) error {
	if s.path == "" {
		return nil
	}

	if err := os.Remove(filepath.Join(s.path, challenge.HTTP01ResourcePath())); err != nil {
		return fmt.Errorf("无法删除 HTTP 挑战文件: %w", err)
	}

	return nil
}

type dnsSolver struct {
	dns     DnsType
	param   DNSParam
	records *[]libdns.Record
}

func (s dnsSolver) Present(ctx context.Context, challenge acme.Challenge) error {
	dnsName := challenge.DNS01TXTRecordName()
	keyAuth := challenge.DNS01KeyAuthorization()
	provider, err := s.getDNSProvider()
	if err != nil {
		return fmt.Errorf("获取 DNS 提供商失败: %w", err)
	}
	zone, err := publicsuffix.EffectiveTLDPlusOne(dnsName)
	if err != nil {
		return fmt.Errorf("获取域名 %q 的顶级域失败: %w", dnsName, err)
	}

	rec := libdns.Record{
		Type:  "TXT",
		Name:  libdns.RelativeName(dnsName+".", zone+"."),
		Value: keyAuth,
	}

	results, err := provider.AppendRecords(ctx, zone+".", []libdns.Record{rec})
	if err != nil {
		return fmt.Errorf("域名 %q 添加临时记录 %q 失败: %w", zone, dnsName, err)
	}
	if len(results) != 1 {
		return fmt.Errorf("预期添加 1 条记录，但实际添加了 %d 条记录", len(results))
	}

	s.records = &results
	return nil
}

func (s dnsSolver) CleanUp(ctx context.Context, challenge acme.Challenge) error {
	dnsName := challenge.DNS01TXTRecordName()
	provider, err := s.getDNSProvider()
	if err != nil {
		return fmt.Errorf("获取 DNS 提供商失败: %w", err)
	}
	zone, err := publicsuffix.EffectiveTLDPlusOne(dnsName)
	if err != nil {
		return fmt.Errorf("获取域名 %q 的顶级域失败: %w", dnsName, err)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	_, err = provider.DeleteRecords(ctx, zone+".", *s.records)
	if err != nil {
		return fmt.Errorf("域名 %q 删除临时记录 %q 失败: %w", zone, dnsName, err)
	}

	return nil
}

func (s dnsSolver) getDNSProvider() (DNSProvider, error) {
	var dns DNSProvider

	switch s.dns {
	case DnsPod:
		dns = &dnspod.Provider{
			APIToken: s.param.ID + "," + s.param.Token,
		}
	case AliYun:
		dns = &alidns.Provider{
			AccKeyID:     s.param.AccessKey,
			AccKeySecret: s.param.SecretKey,
		}
	case CloudFlare:
		dns = &cloudflare.Provider{
			APIToken: s.param.APIkey,
		}
	default:
		return nil, fmt.Errorf("未知的 DNS 提供商 %q", s.dns)
	}

	return dns, nil
}

type DnsType string

const (
	DnsPod     DnsType = "dnspod"
	AliYun     DnsType = "aliyun"
	CloudFlare DnsType = "cloudflare"
)

type DNSParam struct {
	ID        string `form:"id" json:"id"`
	Token     string `form:"token" json:"token"`
	AccessKey string `form:"access_key" json:"access_key"`
	SecretKey string `form:"secret_key" json:"secret_key"`
	APIkey    string `form:"api_key" json:"api_key"`
}

type DNSProvider interface {
	libdns.RecordAppender
	libdns.RecordDeleter
}

type manualDNSSolver struct {
	check       bool
	controlChan chan struct{}
	dataChan    chan any
	records     *[]DNSRecord
}

func (s manualDNSSolver) Present(ctx context.Context, challenge acme.Challenge) error {
	dnsName := challenge.DNS01TXTRecordName()
	keyAuth := challenge.DNS01KeyAuthorization()

	*s.records = append(*s.records, DNSRecord{
		Key:   dnsName,
		Value: keyAuth,
	})
	s.dataChan <- *s.records

	_, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	<-s.controlChan
	return nil
}

func (s manualDNSSolver) CleanUp(_ context.Context, _ acme.Challenge) error {
	return nil
}

type DNSRecord struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
