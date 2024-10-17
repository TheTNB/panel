package rule

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Regex struct{}

func NewRegex() *Regex {
	return &Regex{}
}

func (r *Regex) Regex(fl validator.FieldLevel) bool {
	// 从标签中获取正则，格式类似于 `regex=^[a-zA-Z0-9_]+$`
	pattern := fl.Param()
	// 替换转义字符
	pattern = strings.ReplaceAll(pattern, "，", ",")

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	value := fl.Field().String()
	return re.MatchString(value)
}
