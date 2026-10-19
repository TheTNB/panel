package rule

import (
	"regexp"
	"slices"
	"strings"

	"github.com/leonelquinteros/gotext"
	"github.com/libtnb/validator"
	"github.com/spf13/cast"
)

// Entrance 校验面板访问入口
type Entrance struct {
	t  *gotext.Locale
	re *regexp.Regexp
}

func NewEntrance(t *gotext.Locale) *Entrance {
	return &Entrance{
		t:  t,
		re: regexp.MustCompile(`^[\w.~-]+(/[\w.~-]+)*$`),
	}
}

func (r *Entrance) Signature() string { return "entrance" }

func (r *Entrance) Message() string {
	return r.t.Get("{field} can only contain letters, numbers, -, _, . and ~ separated by /, and cannot start with /api, /login, /webhook or /download")
}

func (r *Entrance) Passes(f *validator.Field) bool {
	entrance := strings.Trim(cast.ToString(f.Reflect().Interface()), "/")
	if entrance == "" {
		return true
	}
	if !r.re.MatchString(entrance) {
		return false
	}
	// 与内置路由同名会互相遮挡
	first, _, _ := strings.Cut(entrance, "/")
	return !slices.Contains([]string{"api", "login", "webhook", "download"}, first)
}
