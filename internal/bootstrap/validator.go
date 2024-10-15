package bootstrap

import (
	"log"

	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/http/rule"
)

func initValidator() {
	translator := zh_Hans_CN.New()
	uni := ut.New(translator, translator)
	trans, _ := uni.GetTranslator("zh_Hans_CN")

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := zh.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatalf("failed to register validator translations: %v", err)
	}

	app.Translator = &trans
	app.Validator = validate
	if err := rule.RegisterRules(validate); err != nil {
		log.Fatalf("failed to register validator rules: %v", err)
	}
}
