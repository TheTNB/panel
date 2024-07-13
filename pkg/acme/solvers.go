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
	"github.com/libdns/tencentcloud"
	"github.com/mholt/acmez/v2/acme"
	"golang.org/x/net/publicsuffix"

	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
)

type httpSolver struct {
	conf string
	path string
}

func (s httpSolver) Present(_ context.Context, challenge acme.Challenge) error {
	var err error
	if s.path == "" {
		return nil
	}

	challengeFilePath := filepath.Join(s.path, challenge.HTTP01ResourcePath())
	if err = os.MkdirAll(filepath.Dir(challengeFilePath), 0755); err != nil {
		return fmt.Errorf("无法在网站目录创建HTTP挑战所需的目录: %w", err)
	}

	if err = os.WriteFile(challengeFilePath, []byte(challenge.KeyAuthorization), 0644); err != nil {
		return fmt.Errorf("无法在网站目录创建HTTP挑战所需的文件: %w", err)
	}

	conf := fmt.Sprintf(`location = /.well-known/acme-challenge/%s {
    default_type text/plain;
    return 200 %q;
}
`, challenge.Token, challenge.KeyAuthorization)
	if err = os.WriteFile(s.conf, []byte(conf), 0644); err != nil {
		return fmt.Errorf("无法写入OpenResty配置文件: %w", err)
	}
	if err = systemctl.Reload("openresty"); err != nil {
		_, err = shell.Execf("openresty -t")
		return fmt.Errorf("无法重载OpenResty: %w", err)
	}

	return nil
}

// CleanUp cleans up the HTTP server if it is the last one to finish.
func (s httpSolver) CleanUp(_ context.Context, challenge acme.Challenge) error {
	if s.path == "" {
		return nil
	}

	_ = os.Remove(filepath.Join(s.path, challenge.HTTP01ResourcePath()))
	_ = os.WriteFile(s.conf, []byte{}, 0644)
	_ = systemctl.Reload("openresty")
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
		return fmt.Errorf("获取DNS提供商失败: %w", err)
	}
	zone, err := publicsuffix.EffectiveTLDPlusOne(dnsName)
	if err != nil {
		return fmt.Errorf("获取域名%q的顶级域失败: %w", dnsName, err)
	}

	rec := libdns.Record{
		Type:  "TXT",
		Name:  libdns.RelativeName(dnsName+".", zone+"."),
		Value: keyAuth,
	}

	results, err := provider.AppendRecords(ctx, zone+".", []libdns.Record{rec})
	if err != nil {
		return fmt.Errorf("域名%q添加临时记录%q失败: %w", zone, dnsName, err)
	}
	if len(results) != 1 {
		return fmt.Errorf("预期添加1条记录，但实际添加了%d条记录", len(results))
	}

	s.records = &results
	return nil
}

func (s dnsSolver) CleanUp(ctx context.Context, challenge acme.Challenge) error {
	dnsName := challenge.DNS01TXTRecordName()
	provider, err := s.getDNSProvider()
	if err != nil {
		return fmt.Errorf("获取DNS提供商失败: %w", err)
	}
	zone, _ := publicsuffix.EffectiveTLDPlusOne(dnsName)

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	_, _ = provider.DeleteRecords(ctx, zone+".", *s.records)
	return nil
}

func (s dnsSolver) getDNSProvider() (DNSProvider, error) {
	var dns DNSProvider

	switch s.dns {
	case DnsPod:
		dns = &dnspod.Provider{
			APIToken: s.param.ID + "," + s.param.Token,
		}
	case Tencent:
		dns = &tencentcloud.Provider{
			SecretId:  s.param.AccessKey,
			SecretKey: s.param.SecretKey,
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
		return nil, fmt.Errorf("未知的DNS提供商 %q", s.dns)
	}

	return dns, nil
}

type DnsType string

const (
	DnsPod     DnsType = "dnspod"
	Tencent    DnsType = "tencent"
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
