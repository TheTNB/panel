package rule

import (
	"github.com/gookit/validate"
	"gorm.io/gorm"
)

func GlobalRules(db *gorm.DB) {
	validate.AddValidators(validate.M{
		"exists":    NewExists(db).Passes,
		"notExists": NewNotExists(db).Passes,
		"password":  NewPassword().Passes,
		"cron":      NewCron().Passes,
	})
	validate.AddGlobalMessages(map[string]string{
		"exists":    "{field} 不存在",
		"notExists": "{field} 已存在",
		"password":  "密码不满足要求（8-20位，至少包含字母、数字、特殊字符中的两种）",
		"cron":      "Cron 表达式不合法",
	})
}
