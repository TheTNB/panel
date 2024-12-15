package rule

import (
	"fmt"

	"gorm.io/gorm"
)

// NotExists 验证一个值在某个表中的字段中不存在，支持同时判断多个字段
// NotExists verify a value does not exist in a table field, support judging multiple fields at the same time
// 用法：not_exists:表名称,字段名称,字段名称,字段名称
// Usage: not_exists:table_name,field_name,field_name,field_name
// 例子：not_exists:users,phone,email
// Example: not_exists:users,phone,email
type NotExists struct {
	DB *gorm.DB
}

func NewNotExists(db *gorm.DB) *NotExists {
	return &NotExists{DB: db}
}

func (r *NotExists) Passes(val any, options ...any) bool {
	if len(options) < 2 {
		return false
	}

	tableName := options[0].(string)
	fieldNames := options[1:]

	query := r.DB.Table(tableName).Where(fmt.Sprintf("%s = ?", fieldNames[0]), val)
	for _, fieldName := range fieldNames[1:] {
		query = query.Or(fmt.Sprintf("%s = ?", fieldName), val)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false
	}

	return count == 0
}
