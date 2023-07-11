package rules

import (
	"regexp"

	"github.com/goravel/framework/contracts/validation"
)

type Regex struct {
}

// Signature The name of the rule.
func (receiver *Regex) Signature() string {
	return "regex"
}

// Passes Determine if the validation rule passes.
func (receiver *Regex) Passes(data validation.Data, val any, options ...any) bool {
	// 第一个参数，正则表达式
	regex := options[0].(string)
	// 用户请求过来的数据
	requestValue, ok := val.(string)
	if !ok {
		return false
	}

	// 判断是否为空
	if len(requestValue) == 0 {
		return false
	}

	// 判断是否匹配
	return regexp.MustCompile(regex).MatchString(requestValue)
}

// Message Get the validation error message.
func (receiver *Regex) Message() string {
	return "格式不正确"
}
