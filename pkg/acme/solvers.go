package acme

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/libdns/alidns"
	"github.com/libdns/cloudflare"
	"github.com/libdns/huaweicloud"
	"github.com/libdns/libdns"
	"github.com/libdns/tencentcloud"
	"github.com/mholt/acmez/v3/acme"
	"golang.org/x/net/publicsuffix"

	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type httpSolver struct {
	conf string
}

func (s httpSolver) Present(_ context.Context, challenge acme.Challenge) error {
	conf := fmt.Sprintf(`location = %s {
    default_type text/plain;
    return 200 %q;
}
`, challenge.HTTP01ResourcePath(), challenge.KeyAuthorization)
	if err := os.WriteFile(s.conf, []byte(conf), 0644); err != nil {
		return fmt.Errorf("无法写入 Nginx 配置文件: %w", err)
	}
	if err := systemctl.Reload("nginx"); err != nil {
		_, err = shell.Execf("nginx -t")
		return fmt.Errorf("无法重载 Nginx: %w", err)
	}

	return nil
}

// CleanUp cleans up the HTTP server if it is the last one to finish.
func (s httpSolver) CleanUp(_ context.Context, challenge acme.Challenge) error {
	_ = os.WriteFile(s.conf, []byte{}, 0644)
	_ = systemctl.Reload("nginx")
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

	results, err := provider.SetRecords(ctx, zone+".", []libdns.Record{rec})
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

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	zone, err := publicsuffix.EffectiveTLDPlusOne(dnsName)
	if err != nil {
		return fmt.Errorf("获取域名 %q 的顶级域失败: %w", dnsName, err)
	}
	_, _ = provider.DeleteRecords(ctx, zone+".", *s.records)
	return nil
}

func (s dnsSolver) getDNSProvider() (DNSProvider, error) {
	var dns DNSProvider

	switch s.dns {
	case AliYun:
		dns = &alidns.Provider{
			AccKeyID:     s.param.AK,
			AccKeySecret: s.param.SK,
		}
	case Tencent:
		dns = &tencentcloud.Provider{
			SecretId:  s.param.AK,
			SecretKey: s.param.SK,
		}
	case Huawei:
		dns = &huaweicloud.Provider{
			AccessKeyId:     s.param.AK,
			SecretAccessKey: s.param.SK,
		}
	case CloudFlare:
		dns = &cloudflare.Provider{
			APIToken: s.param.AK,
		}
	default:
		return nil, fmt.Errorf("未知的DNS提供商 %q", s.dns)
	}

	return dns, nil
}

type DnsType string

const (
	Tencent    DnsType = "tencent"
	AliYun     DnsType = "aliyun"
	Huawei     DnsType = "huawei"
	CloudFlare DnsType = "cloudflare"
)

type DNSParam struct {
	AK string `form:"ak" json:"ak"`
	SK string `form:"sk" json:"sk"`
}

type DNSProvider interface {
	libdns.RecordSetter
	libdns.RecordDeleter
}

type manualDNSSolver struct {
	check       bool
	controlChan chan struct{}
	dataChan    chan any
	records     *[]DNSRecord
}

func (s manualDNSSolver) Present(ctx context.Context, challenge acme.Challenge) error {
	full := challenge.DNS01TXTRecordName()
	keyAuth := challenge.DNS01KeyAuthorization()
	domain, err := publicsuffix.EffectiveTLDPlusOne(full)
	if err != nil {
		return fmt.Errorf("获取 %q 的顶级域失败: %w", full, err)
	}

	*s.records = append(*s.records, DNSRecord{
		Name:   strings.TrimSuffix(full, "."+domain),
		Domain: domain,
		Value:  keyAuth,
	})
	s.dataChan <- *s.records

	<-s.controlChan
	return nil
}

func (s manualDNSSolver) CleanUp(_ context.Context, _ acme.Challenge) error {
	return nil
}

type DNSRecord struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Value  string `json:"value"`
}
