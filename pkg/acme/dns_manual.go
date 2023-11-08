package acme

type Resolve struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Err   string `json:"err"`
}

type manualDnsProvider struct {
	Resolve *Resolve
}

func (p *manualDnsProvider) Present(domain, token, keyAuth string) error {
	return nil
}

func (p *manualDnsProvider) CleanUp(domain, token, keyAuth string) error {
	return nil
}
