package nginx

import (
	"errors"
	"slices"
	"strings"

	"github.com/tufanbarisyildirim/gonginx/config"
	"github.com/tufanbarisyildirim/gonginx/dumper"
	"github.com/tufanbarisyildirim/gonginx/parser"
)

// Parser Nginx vhost 配置解析器
type Parser struct {
	c          *config.Config
	orderIndex map[string]int
}

func NewParser(str ...string) (*Parser, error) {
	if len(str) == 0 {
		str = append(str, defaultConf)
	}
	p := parser.NewStringParser(str[0], parser.WithSkipIncludeParsingErr(), parser.WithSkipValidDirectivesErr())
	c, err := p.Parse()
	if err != nil {
		return nil, err
	}

	orderIndex := make(map[string]int)
	for i, name := range order {
		orderIndex[name] = i
	}

	return &Parser{c: c, orderIndex: orderIndex}, nil
}

func (p *Parser) Config() *config.Config {
	return p.c
}

// Find 通过表达式查找配置
// e.g. Find("server.listen")
func (p *Parser) Find(key string) ([]config.IDirective, error) {
	parts := strings.Split(key, ".")
	var block *config.Block
	var ok bool
	block = p.c.Block
	for i := 0; i < len(parts)-1; i++ {
		key = parts[i]
		directives := block.FindDirectives(key)
		if len(directives) == 0 {
			return nil, errors.New("given key not found")
		}
		if len(directives) > 1 {
			return nil, errors.New("multiple directives found")
		}
		block, ok = directives[0].GetBlock().(*config.Block)
		if !ok {
			return nil, errors.New("block is not *config.Block")
		}
	}

	var result []config.IDirective
	for _, dir := range block.GetDirectives() {
		if dir.GetName() == parts[len(parts)-1] {
			result = append(result, dir)
		}
	}

	return result, nil
}

// FindOne 通过表达式查找一个配置
// e.g. FindOne("server.server_name")
func (p *Parser) FindOne(key string) (config.IDirective, error) {
	directives, err := p.Find(key)
	if err != nil {
		return nil, err
	}
	if len(directives) == 0 {
		return nil, errors.New("given key not found")
	}

	return directives[0], nil
}

// Clear 通过表达式移除配置
// e.g. Clear("server.server_name")
func (p *Parser) Clear(key string) error {
	parts := strings.Split(key, ".")
	last := parts[len(parts)-1]
	parts = parts[:len(parts)-1]

	var block *config.Block
	var ok bool
	block = p.c.Block
	for i := 0; i < len(parts); i++ {
		directives := block.FindDirectives(parts[i])
		if len(directives) == 0 {
			return errors.New("given key not found")
		}
		if len(directives) > 1 {
			return errors.New("multiple directives found")
		}
		block, ok = directives[0].GetBlock().(*config.Block)
		if !ok {
			return errors.New("block is not *config.Block")
		}
	}

	var newDirectives []config.IDirective
	for _, directive := range block.GetDirectives() {
		if directive.GetName() != last {
			newDirectives = append(newDirectives, directive)
		}
	}
	block.Directives = newDirectives

	return nil
}

// Set 通过表达式设置配置
// e.g. Set("server.server_name", []directive)
func (p *Parser) Set(key string, directives []*config.Directive) error {
	parts := strings.Split(key, ".")

	var block *config.Block
	var blockDirective config.IDirective
	var ok bool
	block = p.c.Block
	for i := 0; i < len(parts); i++ {
		sub := block.FindDirectives(parts[i])
		if len(sub) == 0 {
			return errors.New("given key not found")
		}
		if len(sub) > 1 {
			return errors.New("multiple directives found")
		}
		block, ok = sub[0].GetBlock().(*config.Block)
		if !ok {
			return errors.New("block is not *config.Block")
		}
		blockDirective = sub[0]
	}

	for _, directive := range directives {
		directive.SetParent(blockDirective)
		block.Directives = append(block.Directives, directive)
	}

	return nil
}

func (p *Parser) Sort() {
	p.sortDirectives(p.c.Directives, p.orderIndex)
}

func (p *Parser) Dump() string {
	p.Sort()
	return dumper.DumpConfig(p.c, dumper.IndentedStyle)
}

func (p *Parser) sortDirectives(directives []config.IDirective, orderIndex map[string]int) {
	slices.SortFunc(directives, func(a config.IDirective, b config.IDirective) int {
		if orderIndex[a.GetName()] != orderIndex[b.GetName()] {
			return orderIndex[a.GetName()] - orderIndex[b.GetName()]
		}
		return slices.Compare(p.parameters2Slices(a.GetParameters()), p.parameters2Slices(b.GetParameters()))
	})

	for _, directive := range directives {
		if block, ok := directive.GetBlock().(*config.Block); ok {
			p.sortDirectives(block.Directives, orderIndex)
		}
	}
}

func (p *Parser) slices2Parameters(slices []string) []config.Parameter {
	var parameters []config.Parameter
	for _, slice := range slices {
		parameters = append(parameters, config.Parameter{Value: slice})
	}
	return parameters
}

func (p *Parser) parameters2Slices(parameters []config.Parameter) []string {
	var s []string
	for _, parameter := range parameters {
		s = append(s, parameter.Value)
	}
	return s
}
