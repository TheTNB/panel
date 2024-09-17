package data

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"

	"github.com/TheTNB/panel/internal/biz"
)

type firewallRule struct {
	Protocol string
	Port     string
}

type firewallRepo struct {
	ipt *iptables.IPTables
}

func NewFirewallRepo() biz.FirewallRepo {
	ipt, err := iptables.New()
	if err != nil {
		panic(err)
	}

	return &firewallRepo{
		ipt: ipt,
	}
}

func (r *firewallRepo) GetRules() ([]firewallRule, error) {
	raw, err := r.ipt.List("filter", "INPUT")
	if err != nil {
		return nil, err
	}

	var rules []firewallRule
	for _, line := range raw {
		fmt.Println(line)
	}

	return rules, nil
}
