package rule

import (
	"regexp"

	"github.com/spf13/cast"
)

// Cron 校验规则
type Cron struct {
	re *regexp.Regexp
}

func NewCron() *Cron {
	return &Cron{
		re: regexp.MustCompile(`(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|((\*|\d+)(\/|-)\d+)|\d+|\*) ?){5,7})`),
	}
}

func (s *Cron) Passes(val any, options ...any) bool {
	return s.re.MatchString(cast.ToString(val))
}
