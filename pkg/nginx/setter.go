package nginx

import (
	"fmt"
	"slices"
	"strings"

	"github.com/tufanbarisyildirim/gonginx/config"
)

func (p *Parser) SetListen(listen [][]string) error {
	var directives []*config.Directive
	for _, l := range listen {
		directives = append(directives, &config.Directive{
			Name:       "listen",
			Parameters: p.slices2Parameters(l),
		})
	}

	if err := p.Clear("server.listen"); err != nil {
		return err
	}

	return p.Set("server", directives)
}

func (p *Parser) SetServerName(serverName []string) error {
	if err := p.Clear("server.server_name"); err != nil {
		return err
	}

	return p.Set("server", []*config.Directive{
		{
			Name:       "server_name",
			Parameters: p.slices2Parameters(serverName),
		},
	})
}

func (p *Parser) SetIndex(index []string) error {
	if err := p.Clear("server.index"); err != nil {
		return err
	}

	return p.Set("server", []*config.Directive{
		{
			Name:       "index",
			Parameters: p.slices2Parameters(index),
		},
	})
}

func (p *Parser) SetIndexWithComment(index []string, comment []string) error {
	if err := p.Clear("server.index"); err != nil {
		return err
	}

	return p.Set("server", []*config.Directive{
		{
			Name:       "index",
			Parameters: p.slices2Parameters(index),
			Comment:    comment,
		},
	})
}

func (p *Parser) SetRoot(root string) error {
	if err := p.Clear("server.root"); err != nil {
		return err
	}

	return p.Set("server", []*config.Directive{
		{
			Name:       "root",
			Parameters: []config.Parameter{{Value: root}},
		},
	})
}

func (p *Parser) SetRootWithComment(root string, comment []string) error {
	if err := p.Clear("server.root"); err != nil {
		return err
	}

	return p.Set("server", []*config.Directive{
		{
			Name:       "root",
			Parameters: []config.Parameter{{Value: root}},
			Comment:    comment,
		},
	})
}

func (p *Parser) SetIncludes(includes []string, comments [][]string) error {
	if err := p.Clear("server.include"); err != nil {
		return err
	}

	var directives []*config.Directive
	for i, item := range includes {
		var comment []string
		if i < len(comments) {
			comment = comments[i]
		}
		directives = append(directives, &config.Directive{
			Name:       "include",
			Parameters: []config.Parameter{{Value: item}},
			Comment:    comment,
		})
	}

	return p.Set("server", directives)
}

func (p *Parser) SetPHP(php int) error {
	old, err := p.Find("server.include")
	if err != nil {
		return err
	}
	if err = p.Clear("server.include"); err != nil {
		return err
	}

	var directives []*config.Directive
	var foundFlag bool
	for _, item := range old {
		// 查找enable-php的配置
		if slices.ContainsFunc(p.parameters2Slices(item.GetParameters()), func(s string) bool {
			return strings.HasPrefix(s, "enable-php-") && strings.HasSuffix(s, ".conf")
		}) {
			foundFlag = true
			directives = append(directives, &config.Directive{
				Name:       item.GetName(),
				Parameters: []config.Parameter{{Value: fmt.Sprintf("enable-php-%d.conf", php)}},
				Comment:    item.GetComment(),
			})
		} else {
			// 其余的原样保留
			directives = append(directives, &config.Directive{
				Name:       item.GetName(),
				Parameters: item.GetParameters(),
				Comment:    item.GetComment(),
			})
		}
	}

	// 如果没有找到enable-php的配置，直接添加一个
	if !foundFlag {
		directives = append(directives, &config.Directive{
			Name:       "include",
			Parameters: []config.Parameter{{Value: fmt.Sprintf("enable-php-%d.conf", php)}},
		})
	}

	return p.Set("server", directives)
}

func (p *Parser) ClearSetHTTPS() error {
	if err := p.Clear("server.ssl_certificate"); err != nil {
		return err
	}
	if err := p.Clear("server.ssl_certificate_key"); err != nil {
		return err
	}
	if err := p.Clear("server.ssl_session_timeout"); err != nil {
		return err
	}
	if err := p.Clear("server.ssl_session_cache"); err != nil {
		return err
	}
	if err := p.Clear("server.ssl_protocols"); err != nil {
		return err
	}
	if err := p.Clear("server.ssl_ciphers"); err != nil {
		return err
	}
	if err := p.Clear("server.ssl_prefer_server_ciphers"); err != nil {
		return err
	}
	if err := p.Clear("server.ssl_early_data"); err != nil {
		return err
	}

	return nil
}

func (p *Parser) SetHTTPS(cert, key string) error {
	if err := p.ClearSetHTTPS(); err != nil {
		return err
	}

	return p.Set("server", []*config.Directive{
		{
			Name:       "ssl_certificate",
			Parameters: []config.Parameter{{Value: cert}},
		},
		{
			Name:       "ssl_certificate_key",
			Parameters: []config.Parameter{{Value: key}},
		},
		{
			Name:       "ssl_session_timeout",
			Parameters: []config.Parameter{{Value: "1d"}},
		},
		{
			Name:       "ssl_session_cache",
			Parameters: []config.Parameter{{Value: "shared:SSL:10m"}},
		},
		{
			Name:       "ssl_protocols",
			Parameters: []config.Parameter{{Value: "TLSv1.2"}, {Value: "TLSv1.3"}},
		},
		{
			Name:       "ssl_ciphers",
			Parameters: []config.Parameter{{Value: "ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-CHACHA20-POLY1305"}},
		},
		{
			Name:       "ssl_prefer_server_ciphers",
			Parameters: []config.Parameter{{Value: "off"}},
		},
		{
			Name:       "ssl_early_data",
			Parameters: []config.Parameter{{Value: "on"}},
		},
	})
}

func (p *Parser) SetHTTPSProtocols(protocols []string) error {
	if err := p.Clear("server.ssl_protocols"); err != nil {
		return err
	}

	return p.Set("server", []*config.Directive{
		{
			Name:       "ssl_protocols",
			Parameters: p.slices2Parameters(protocols),
		},
	})
}

func (p *Parser) SetHTTPSCiphers(ciphers string) error {
	if err := p.Clear("server.ssl_ciphers"); err != nil {
		return err
	}

	return p.Set("server", []*config.Directive{
		{
			Name:       "ssl_ciphers",
			Parameters: []config.Parameter{{Value: ciphers}},
		},
	})
}

func (p *Parser) SetOCSP(ocsp bool) error {
	if err := p.Clear("server.ssl_stapling"); err != nil {
		return err
	}
	if err := p.Clear("server.ssl_stapling_verify"); err != nil {
		return err
	}

	if ocsp {
		return p.Set("server", []*config.Directive{
			{
				Name:       "ssl_stapling",
				Parameters: []config.Parameter{{Value: "on"}},
			},
			{
				Name:       "ssl_stapling_verify",
				Parameters: []config.Parameter{{Value: "on"}},
			},
		})
	}

	return nil
}

func (p *Parser) SetHSTS(hsts bool) error {
	old, err := p.Find("server.add_header")
	if err != nil {
		return err
	}
	if err = p.Clear("server.add_header"); err != nil {
		return err
	}

	var directives []*config.Directive
	var foundFlag bool
	for _, dir := range old {
		if slices.Contains(p.parameters2Slices(dir.GetParameters()), "Strict-Transport-Security") {
			foundFlag = true
			if hsts {
				directives = append(directives, &config.Directive{
					Name:       dir.GetName(),
					Parameters: []config.Parameter{{Value: "Strict-Transport-Security"}, {Value: "max-age=31536000"}},
					Comment:    dir.GetComment(),
				})
			}
		} else {
			directives = append(directives, &config.Directive{
				Name:       dir.GetName(),
				Parameters: dir.GetParameters(),
				Comment:    dir.GetComment(),
			})
		}
	}

	if !foundFlag && hsts {
		directives = append(directives, &config.Directive{
			Name:       "add_header",
			Parameters: []config.Parameter{{Value: "Strict-Transport-Security"}, {Value: "max-age=31536000"}},
		})
	}

	return p.Set("server", directives)
}

func (p *Parser) SetHTTPRedirect(httpRedirect bool) error {
	// if 重定向
	ifs, err := p.Find("server.if")
	if err != nil {
		return err
	}
	if err = p.Clear("server.if"); err != nil {
		return err
	}

	var directives []*config.Directive
	var foundFlag bool
	for _, dir := range ifs { // 所有 if
		if !httpRedirect {
			if len(dir.GetParameters()) == 3 && dir.GetParameters()[0].GetValue() == "($scheme" && dir.GetParameters()[1].GetValue() == "=" && dir.GetParameters()[2].GetValue() == "http)" {
				continue
			}
		}
		var ifDirectives []config.IDirective
		for _, dir2 := range dir.GetBlock().GetDirectives() { // 每个 if 中所有指令
			if !httpRedirect {
				// 不启用http重定向，则判断并移除特定的return指令
				if dir2.GetName() != "return" && !slices.Contains(p.parameters2Slices(dir2.GetParameters()), "https://$host$request_uri") {
					ifDirectives = append(ifDirectives, dir2)
				}
			} else {
				// 启用http重定向，需要检查防止重复添加
				if dir2.GetName() == "return" && slices.Contains(p.parameters2Slices(dir2.GetParameters()), "https://$host$request_uri") {
					foundFlag = true
				}
				ifDirectives = append(ifDirectives, dir2)
			}
		}
		// 写回 if 指令
		if block, ok := dir.GetBlock().(*config.Block); ok {
			block.Directives = ifDirectives
		}
		directives = append(directives, &config.Directive{
			Block:      dir.GetBlock(),
			Name:       dir.GetName(),
			Parameters: dir.GetParameters(),
			Comment:    dir.GetComment(),
		})
	}

	if !foundFlag && httpRedirect {
		ifDir := &config.Directive{
			Name:       "if",
			Block:      &config.Block{},
			Parameters: []config.Parameter{{Value: "($scheme"}, {Value: "="}, {Value: "http)"}},
		}
		redirectDir := &config.Directive{
			Name:       "return",
			Parameters: []config.Parameter{{Value: "308"}, {Value: "https://$host$request_uri"}},
		}
		redirectDir.SetParent(ifDir.GetParent())
		ifBlock := ifDir.GetBlock().(*config.Block)
		ifBlock.Directives = append(ifBlock.Directives, redirectDir)
		directives = append(directives, ifDir)
	}

	if err = p.Set("server", directives); err != nil {
		return err
	}

	// error_page 497 重定向
	directives = nil
	errorPages, err := p.Find("server.error_page")
	if err != nil {
		return err
	}
	if err = p.Clear("server.error_page"); err != nil {
		return err
	}
	var found497 bool
	for _, dir := range errorPages {
		if !httpRedirect {
			// 不启用https重定向，则判断并移除特定的return指令
			if !slices.Contains(p.parameters2Slices(dir.GetParameters()), "497") && !slices.Contains(p.parameters2Slices(dir.GetParameters()), "https://$host:$server_port$request_uri") {
				directives = append(directives, &config.Directive{
					Block:      dir.GetBlock(),
					Name:       dir.GetName(),
					Parameters: dir.GetParameters(),
					Comment:    dir.GetComment(),
				})
			}
		} else {
			// 启用https重定向，需要检查防止重复添加
			if slices.Contains(p.parameters2Slices(dir.GetParameters()), "497") && slices.Contains(p.parameters2Slices(dir.GetParameters()), "https://$host:$server_port$request_uri") {
				found497 = true
			}
			directives = append(directives, &config.Directive{
				Block:      dir.GetBlock(),
				Name:       dir.GetName(),
				Parameters: dir.GetParameters(),
				Comment:    dir.GetComment(),
			})
		}
	}

	if !found497 && httpRedirect {
		directives = append(directives, &config.Directive{
			Name:       "error_page",
			Parameters: []config.Parameter{{Value: "497"}, {Value: "=308"}, {Value: "https://$host:$server_port$request_uri"}},
		})
	}

	return p.Set("server", directives)
}

func (p *Parser) SetAltSvc(altSvc string) error {
	old, err := p.Find("server.add_header")
	if err != nil {
		return err
	}
	if err = p.Clear("server.add_header"); err != nil {
		return err
	}

	var directives []*config.Directive
	var foundFlag bool
	for _, dir := range old {
		if slices.Contains(p.parameters2Slices(dir.GetParameters()), "Alt-Svc") {
			foundFlag = true
			if altSvc != "" { // 为空表示要删除
				directives = append(directives, &config.Directive{
					Name:       dir.GetName(),
					Parameters: []config.Parameter{{Value: "Alt-Svc"}, {Value: altSvc}},
					Comment:    dir.GetComment(),
				})
			}
		} else {
			directives = append(directives, &config.Directive{
				Name:       dir.GetName(),
				Parameters: dir.GetParameters(),
				Comment:    dir.GetComment(),
			})
		}
	}

	if !foundFlag && altSvc != "" {
		directives = append(directives, &config.Directive{
			Name:       "add_header",
			Parameters: []config.Parameter{{Value: "Alt-Svc"}, {Value: altSvc}},
		})
	}

	return p.Set("server", directives)
}

func (p *Parser) SetAccessLog(accessLog string) error {
	if err := p.Clear("server.access_log"); err != nil {
		return err
	}

	return p.Set("server", []*config.Directive{
		{
			Name:       "access_log",
			Parameters: []config.Parameter{{Value: accessLog}},
		},
	})
}

func (p *Parser) SetErrorLog(errorLog string) error {
	if err := p.Clear("server.error_log"); err != nil {
		return err
	}

	return p.Set("server", []*config.Directive{
		{
			Name:       "error_log",
			Parameters: []config.Parameter{{Value: errorLog}},
		},
	})
}
