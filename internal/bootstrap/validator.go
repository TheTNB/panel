package bootstrap

import (
	"github.com/gookit/validate"
	"github.com/gookit/validate/locales/zhcn"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/internal/http/rule"
)

// NewValidator just for register global rules
func NewValidator(db *gorm.DB) *validate.Validation {
	zhcn.RegisterGlobal()
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
		opt.SkipOnEmpty = true
		opt.FieldTag = "form"
	})

	// register global rules
	rule.GlobalRules(db)

	return validate.NewEmpty()
}
