package acme

type Resolve struct {
	Key   string
	Value string
	Err   string
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
