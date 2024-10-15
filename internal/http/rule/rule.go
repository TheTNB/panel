package rule

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/TheTNB/panel/internal/app"
)

func RegisterRules(v *validator.Validate) error {
	if err := v.RegisterValidation("exists", NewExists(app.Orm).Exists); err != nil {
		return err
	}
	if err := v.RegisterValidation("not_exists", NewNotExists(app.Orm).NotExists); err != nil {
		return err
	}
	if err := v.RegisterValidation("regexp", NewRegexp().Regexp); err != nil {
		return err
	}

	if err := v.RegisterTranslation("exists", *app.Translator,
		func(ut ut.Translator) error {
			return ut.Add("exists", "{0} 不存在", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("exists", fe.Field())
			return t
		}); err != nil {
		return err
	}
	if err := v.RegisterTranslation("not_exists", *app.Translator,
		func(ut ut.Translator) error {
			return ut.Add("not_exists", "{0} 已存在", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("not_exists", fe.Field())
			return t
		}); err != nil {
		return err
	}
	if err := v.RegisterTranslation("regexp", *app.Translator,
		func(ut ut.Translator) error {
			return ut.Add("regexp", "{0} 格式不正确", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("regexp", fe.Field())
			return t
		}); err != nil {
		return err
	}

	return nil
}
