package bootstrap

import (
	"github.com/gookit/validate"
	"github.com/gookit/validate/locales/zhcn"
)

func init() {
	zhcn.RegisterGlobal()
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
		opt.SkipOnEmpty = true
		opt.FieldTag = "form"
	})
}
