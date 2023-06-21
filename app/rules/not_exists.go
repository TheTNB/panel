package rules

import (
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type NotExists struct {
}

// Signature The name of the rule.
func (receiver *NotExists) Signature() string {
	return "not_exists"
}

// Passes Determine if the validation rule passes.
func (receiver *NotExists) Passes(_ validation.Data, val any, options ...any) bool {

	// 第一个参数，表名称，如 categories
	tableName := options[0].(string)
	// 第二个参数，字段名称，如 id
	fieldName := options[1].(string)
	// 用户请求过来的数据
	requestValue, ok := val.(string)
	if !ok {
		return false
	}

	// 判断是否为空
	if len(requestValue) == 0 {
		return false
	}

	// 判断是否存在
	var count int64
	query := facades.Orm().Query().Table(tableName).Where(fieldName, requestValue)
	// 判断第三个参数及之后的参数是否存在
	if len(options) > 2 {
		for i := 2; i < len(options); i++ {
			query = query.OrWhere(options[i].(string), requestValue)
		}
	}
	err := query.Count(&count)
	if err != nil {
		return false
	}

	return count == 0
}

// Message Get the validation error message.
func (receiver *NotExists) Message() string {
	return "记录已存在"
}
