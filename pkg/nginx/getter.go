package nginx

import (
	"fmt"
	"slices"
	"strings"
)

func (p *Parser) GetListen() ([][]string, error) {
	directives, err := p.Find("server.listen")
	if err != nil {
		return nil, err
	}

	var result [][]string
	for _, dir := range directives {
		result = append(result, dir.GetParameters())
	}

	return result, nil
}

func (p *Parser) GetServerName() ([]string, error) {
	directive, err := p.FindOne("server.server_name")
	if err != nil {
		return nil, err
	}

	return directive.GetParameters(), nil
}

func (p *Parser) GetIndex() ([]string, error) {
	directive, err := p.FindOne("server.index")
	if err != nil {
		return nil, err
	}

	return directive.GetParameters(), nil
}

func (p *Parser) GetRoot() (string, error) {
	directive, err := p.FindOne("server.root")
	if err != nil {
		return "", err
	}
	if len(directive.GetParameters()) == 0 {
		return "", nil
	}

	return directive.GetParameters()[0], nil
}

func (p *Parser) GetIncludes() (includes []string, comments [][]string, err error) {
	directives, err := p.Find("server.include")
	if err != nil {
		return nil, nil, err
	}

	for _, dir := range directives {
		if len(dir.GetParameters()) != 1 {
			return nil, nil, fmt.Errorf("invalid include directive, expected 1 parameter but got %d", len(dir.GetParameters()))
		}
		includes = append(includes, dir.GetParameters()[0])
		comments = append(comments, dir.GetComment())
	}

	return includes, comments, nil
}

func (p *Parser) GetPHP() (int, error) {
	directives, err := p.Find("server.include")
	if err != nil {
		return 0, err
	}

	var result int
	for _, dir := range directives {
		if slices.ContainsFunc(dir.GetParameters(), func(s string) bool {
			return strings.HasPrefix(s, "enable-php-") && strings.HasSuffix(s, ".conf")
		}) {
			_, err = fmt.Sscanf(dir.GetParameters()[0], "enable-php-%d.conf", &result)
		}
	}

	return result, err
}

func (p *Parser) GetHTTPS() bool {
	directive, err := p.FindOne("server.ssl_certificate")
	if err != nil {
		return false
	}
	if len(directive.GetParameters()) == 0 {
		return false
	}

	return true
}

func (p *Parser) GetHTTPSProtocols() []string {
	directive, err := p.FindOne("server.ssl_protocols")
	if err != nil {
		return nil
	}

	return directive.GetParameters()
}

func (p *Parser) GetHTTPSCiphers() string {
	directive, err := p.FindOne("server.ssl_ciphers")
	if err != nil {
		return ""
	}
	if len(directive.GetParameters()) == 0 {
		return ""
	}

	return directive.GetParameters()[0]
}

func (p *Parser) GetOCSP() bool {
	directive, err := p.FindOne("server.ssl_stapling")
	if err != nil {
		return false
	}
	if len(directive.GetParameters()) == 0 {
		return false
	}

	return directive.GetParameters()[0] == "on"
}

func (p *Parser) GetHSTS() bool {
	directives, err := p.Find("server.add_header")
	if err != nil {
		return false
	}

	for _, dir := range directives {
		if slices.Contains(dir.GetParameters(), "Strict-Transport-Security") {
			return true
		}
	}

	return false
}

func (p *Parser) GetHTTPSRedirect() (bool, error) {
	directives, err := p.Find("server.if")
	if err != nil {
		return false, err
	}

	for _, dir := range directives {
		for _, dir2 := range dir.GetBlock().GetDirectives() {
			if dir2.GetName() == "return" && slices.Contains(dir2.GetParameters(), "https://$host$request_uri") {
				return true, nil
			}
		}
	}

	return false, nil
}

func (p *Parser) GetAccessLog() (string, error) {
	directive, err := p.FindOne("server.access_log")
	if err != nil {
		return "", err
	}
	if len(directive.GetParameters()) == 0 {
		return "", nil
	}

	return directive.GetParameters()[0], nil
}

func (p *Parser) GetErrorLog() (string, error) {
	directive, err := p.FindOne("server.error_log")
	if err != nil {
		return "", err
	}
	if len(directive.GetParameters()) == 0 {
		return "", nil
	}

	return directive.GetParameters()[0], nil
}
