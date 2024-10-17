package firewall

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/cast"

	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type Firewall struct {
	forwardListRegex *regexp.Regexp
	richRuleRegex    *regexp.Regexp
}

func NewFirewall() *Firewall {
	firewall := &Firewall{
		forwardListRegex: regexp.MustCompile(`^port=(\d{1,5}):proto=(.+?):toport=(\d{1,5}):toaddr=(.*)$`),
		richRuleRegex:    regexp.MustCompile(`^rule family="([^"]+)"(?: .*?(source|destination) address="([^"]+)")?(?: .*?port port="([^"]+)")?(?: .*?protocol="([^"]+)")?.*?(accept|drop|reject)$`),
	}

	return firewall
}

func (r *Firewall) Status() (bool, error) {
	return systemctl.Status("firewalld")
}

func (r *Firewall) Version() (string, error) {
	return shell.Execf("firewall-cmd --version")
}

func (r *Firewall) ListRule() ([]FireInfo, error) {
	var wg sync.WaitGroup
	var data []FireInfo
	wg.Add(2)

	go func() {
		defer wg.Done()
		out, err := shell.Execf("firewall-cmd --zone=public --list-ports")
		if err != nil {
			return
		}
		ports := strings.Split(out, " ")
		for _, port := range ports {
			if len(port) == 0 {
				continue
			}
			var item FireInfo
			if strings.Contains(port, "/") {
				ruleItem := strings.Split(port, "/")
				portItem := strings.Split(ruleItem[0], "-")
				if len(portItem) > 1 {
					item.PortStart = cast.ToUint(portItem[0])
					item.PortEnd = cast.ToUint(portItem[1])
				} else {
					item.PortStart = cast.ToUint(ruleItem[0])
					item.PortEnd = cast.ToUint(ruleItem[0])
				}
				item.Protocol = Protocol(ruleItem[1])
			}
			item.Family = "ipv4"
			item.Strategy = "accept"
			item.Direction = "in"
			data = append(data, item)
		}
	}()
	go func() {
		defer wg.Done()
		rich, err := r.ListRichRule()
		if err != nil {
			return
		}
		data = append(data, rich...)
	}()

	wg.Wait()
	return data, nil
}

func (r *Firewall) ListForward() ([]FireForwardInfo, error) {
	out, err := shell.Execf("firewall-cmd --zone=public --list-forward-ports")
	if err != nil {
		return nil, err
	}

	var data []FireForwardInfo
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimFunc(line, func(r rune) bool {
			return r <= 32
		})
		if r.forwardListRegex.MatchString(line) {
			match := r.forwardListRegex.FindStringSubmatch(line)
			if len(match) < 4 {
				continue
			}
			if len(match[4]) == 0 {
				match[4] = "127.0.0.1"
			}
			data = append(data, FireForwardInfo{
				Port:       cast.ToUint(match[1]),
				Protocol:   Protocol(match[2]),
				TargetIP:   match[4],
				TargetPort: match[3],
			})
		}
	}

	return data, nil
}

func (r *Firewall) ListRichRule() ([]FireInfo, error) {
	out, err := shell.Execf("firewall-cmd --zone=public --list-rich-rules")
	if err != nil {
		return nil, err
	}

	var data []FireInfo
	rules := strings.Split(out, "\n")
	for _, rule := range rules {
		if len(rule) == 0 {
			continue
		}
		if richRules, err := r.parseRichRule(rule); err == nil {
			data = append(data, richRules)
		}
	}

	return data, nil
}

func (r *Firewall) Port(rule FireInfo, operation Operation) error {
	if rule.PortEnd == 0 {
		rule.PortEnd = rule.PortStart
	}
	if rule.PortStart > rule.PortEnd {
		return fmt.Errorf("invalid port range: %d-%d", rule.PortStart, rule.PortEnd)
	}
	// 不支持的切换使用rich rules
	if (rule.Family != "" && rule.Family != "ipv4") || rule.Direction != "in" || rule.Address != "" || rule.Strategy != "accept" {
		return r.RichRules(rule, operation)
	}

	protocols := strings.Split(string(rule.Protocol), "/")
	for protocol := range slices.Values(protocols) {
		stdout, err := shell.Execf("firewall-cmd --zone=public --%s-port=%d-%d/%s --permanent", operation, rule.PortStart, rule.PortEnd, protocol)
		if err != nil {
			return fmt.Errorf("%s port %d-%d/%s failed, err: %s", operation, rule.PortStart, rule.PortEnd, protocol, stdout)
		}
	}

	_, err := shell.Execf("firewall-cmd --reload")
	return err
}

func (r *Firewall) RichRules(rule FireInfo, operation Operation) error {
	protocols := strings.Split(string(rule.Protocol), "/")
	for protocol := range slices.Values(protocols) {
		var ruleBuilder strings.Builder
		ruleBuilder.WriteString(fmt.Sprintf(`rule family="%s" `, rule.Family))

		if len(rule.Address) != 0 {
			if rule.Direction == "in" {
				ruleBuilder.WriteString(fmt.Sprintf(`source address="%s" `, rule.Address))
			} else if rule.Direction == "out" {
				ruleBuilder.WriteString(fmt.Sprintf(`destination address="%s" `, rule.Address))
			} else if rule.Direction != "" {
				return fmt.Errorf("invalid direction: %s", rule.Direction)
			}
		}
		if rule.PortStart != 0 && rule.PortEnd != 0 && (rule.PortStart != 1 && rule.PortEnd != 65535) { // 1-65535是解析出来无端口规则的情况
			ruleBuilder.WriteString(fmt.Sprintf(`port port="%d-%d" `, rule.PortStart, rule.PortEnd))
		}
		if operation == OperationRemove && protocol != "" && rule.Protocol != "tcp/udp" { // 删除操作，可以不指定协议
			ruleBuilder.WriteString(fmt.Sprintf(`protocol="%s" `, protocol))
		}
		if operation == OperationAdd && protocol != "" {
			ruleBuilder.WriteString(fmt.Sprintf(`protocol="%s" `, protocol))
		}

		ruleBuilder.WriteString(string(rule.Strategy))
		_, err := shell.Execf("firewall-cmd --zone=public --%s-rich-rule '%s' --permanent", operation, ruleBuilder.String())
		if err != nil {
			return fmt.Errorf("%s rich rules (%s) failed, err: %v", operation, ruleBuilder.String(), err)
		}
	}

	_, err := shell.Execf("firewall-cmd --reload")
	return err
}

func (r *Firewall) PortForward(info Forward, operation Operation) error {
	if err := r.enableForward(); err != nil {
		return err
	}

	var ruleStr strings.Builder
	ruleStr.WriteString(fmt.Sprintf("firewall-cmd --zone=public --%s-forward-port=port=%d:proto=%s:", operation, info.Port, info.Protocol))
	if info.TargetIP != "" && info.TargetIP != "127.0.0.1" && info.TargetIP != "localhost" {
		ruleStr.WriteString(fmt.Sprintf("toaddr=%s:toport=%d", info.TargetIP, info.TargetPort))
	} else {
		ruleStr.WriteString(fmt.Sprintf("toport=%d", info.TargetPort))
	}
	ruleStr.WriteString(" --permanent")

	_, err := shell.Execf(ruleStr.String()) // nolint: govet
	if err != nil {
		return fmt.Errorf("%s port forward failed, err: %v", operation, err)
	}

	_, err = shell.Execf("firewall-cmd --reload")
	return err
}

func (r *Firewall) parseRichRule(line string) (FireInfo, error) {
	if !r.richRuleRegex.MatchString(line) {
		return FireInfo{}, errors.New("invalid rich rule format")
	}

	match := r.richRuleRegex.FindStringSubmatch(line)
	if len(match) < 7 {
		return FireInfo{}, errors.New("invalid rich rule")
	}

	fireInfo := FireInfo{
		Family:   match[1],
		Address:  match[3],
		Protocol: Protocol(match[5]),
		Strategy: Strategy(match[6]),
	}

	if match[2] == "destination" {
		fireInfo.Direction = "out"
	} else {
		fireInfo.Direction = "in"
	}
	if fireInfo.Protocol == "" {
		fireInfo.Protocol = "tcp/udp"
	}

	ports := strings.Split(match[4], "-")
	if len(ports) == 2 { // 添加端口范围
		fireInfo.PortStart = cast.ToUint(ports[0])
		fireInfo.PortEnd = cast.ToUint(ports[1])
	} else if len(ports) == 1 && ports[0] != "" { // 添加单个端口
		port := cast.ToUint(ports[0])
		fireInfo.PortStart = port
		fireInfo.PortEnd = port
	} else if len(ports) == 1 && ports[0] == "" { // 未添加端口规则，表示所有端口
		fireInfo.PortStart = 1
		fireInfo.PortEnd = 65535
	}

	return fireInfo, nil
}

func (r *Firewall) enableForward() error {
	out, err := shell.Execf("firewall-cmd --zone=public --query-masquerade")
	if err != nil {
		if out == "no" {
			out, err = shell.Execf("firewall-cmd --zone=public --add-masquerade --permanent")
			if err != nil {
				return fmt.Errorf("%s: %s", err, out)
			}

			_, err = shell.Execf("firewall-cmd --reload")
			return err
		}

		return fmt.Errorf("%v: %s", err, out)
	}

	return nil
}
