package rules

import (
	"github.com/goravel/framework/contracts/validation"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/pkg/io"
)

type PathNotExists struct {
}

// Signature The name of the rule.
func (receiver *PathNotExists) Signature() string {
	return "path_not_exists"
}

// Passes Determine if the validation rule passes.
func (receiver *PathNotExists) Passes(_ validation.Data, val any, options ...any) bool {
	// 用户请求过来的数据
	requestValue, err := cast.ToStringE(val)
	if err != nil {
		return false
	}

	// 判断是否为空
	if len(requestValue) == 0 {
		return false
	}

	return !io.Exists(requestValue)
}

// Message Get the validation error message.
func (receiver *PathNotExists) Message() string {
	return "路径已存在"
}
