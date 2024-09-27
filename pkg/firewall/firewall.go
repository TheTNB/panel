package firewall

import (
	"errors"
	"fmt"
	"regexp"
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
		richRuleRegex:    regexp.MustCompile(`^rule family="([^"]+)" (?:source address="([^"]+)" )?(?:port port="([^"]+)" )?(?:protocol="([^"]+)" )?(accept|drop|reject)$`),
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
			var itemPort FireInfo
			if strings.Contains(port, "/") {
				itemPort.Port = cast.ToUint(strings.Split(port, "/")[0])
				itemPort.Protocol = strings.Split(port, "/")[1]
			}
			itemPort.Strategy = "accept"
			data = append(data, itemPort)
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

func (r *Firewall) ListForward() ([]FireInfo, error) {
	out, err := shell.Execf("firewall-cmd --zone=public --list-forward-ports")
	if err != nil {
		return nil, err
	}

	var data []FireInfo
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
			data = append(data, FireInfo{
				Port:       cast.ToUint(match[1]),
				Protocol:   match[2],
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
		if itemRule, err := r.parseRichRule(rule); err == nil {
			data = append(data, *itemRule)
		}
	}

	return data, nil
}

func (r *Firewall) Port(port FireInfo, operation string) error {
	stdout, err := shell.Execf("firewall-cmd --zone=public --%s-port=%d/%s --permanent", operation, port.Port, port.Protocol)
	if err != nil {
		return fmt.Errorf("%s port %d/%s failed, err: %s", operation, port.Port, port.Protocol, stdout)
	}

	_, err = shell.Execf("firewall-cmd --reload")
	return err
}

func (r *Firewall) RichRules(rule FireInfo, operation string) error {
	families := strings.Split(rule.Family, "/") // ipv4 ipv6

	for _, family := range families {
		var ruleStr strings.Builder
		ruleStr.WriteString(fmt.Sprintf(`rule family="%s" `, family))
		if len(rule.Address) != 0 {
			ruleStr.WriteString(fmt.Sprintf(`source address="%s" `, rule.Address))
		}
		if rule.Port != 0 {
			ruleStr.WriteString(fmt.Sprintf(`port port="%d" `, rule.Port))
		}
		if len(rule.Protocol) != 0 {
			ruleStr.WriteString(fmt.Sprintf(`protocol="%s" `, rule.Protocol))
		}

		ruleStr.WriteString(rule.Strategy)
		out, err := shell.Execf("firewall-cmd --zone=public --%s-rich-rule '%s' --permanent", operation, ruleStr.String())
		if err != nil {
			return fmt.Errorf("%s rich rules (%s) failed, err: %s", operation, ruleStr.String(), out)
		}
	}

	_, err := shell.Execf("firewall-cmd --reload")
	return err
}

func (r *Firewall) PortForward(info Forward, operation string) error {
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

	out, err := shell.Execf(ruleStr.String()) // nolint: govet
	if err != nil {
		return fmt.Errorf("%s port forward failed, err: %s", operation, out)
	}

	_, err = shell.Execf("firewall-cmd --reload")
	return err
}

func (r *Firewall) parseRichRule(line string) (*FireInfo, error) {
	itemRule := new(FireInfo)
	if r.richRuleRegex.MatchString(line) {
		match := r.richRuleRegex.FindStringSubmatch(line)
		if len(match) < 6 {
			return nil, errors.New("invalid rich rule")
		}

		itemRule.Family = match[1]
		itemRule.Address = match[2]
		itemRule.Port = cast.ToUint(match[3])
		itemRule.Protocol = match[4]
		itemRule.Strategy = match[5]
	}

	return itemRule, nil
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
