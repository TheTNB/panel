package rule

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Regexp struct{}

func NewRegexp() *Regexp {
	return &Regexp{}
}

func (r *Regexp) Regexp(fl validator.FieldLevel) bool {
	// 从标签中获取正则，格式类似于 `regexp=^[a-zA-Z0-9_]+$`
	pattern := fl.Param()
	value := fl.Field().String()

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return re.MatchString(value)
}
