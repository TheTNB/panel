package app

import (
	"github.com/go-chi/chi/v5"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-rat/sessions"
	"github.com/knadh/koanf/v2"
	"gorm.io/gorm"
)

var (
	Conf       *koanf.Koanf
	Http       *chi.Mux
	Orm        *gorm.DB
	Validator  *validator.Validate
	Translator *ut.Translator
	Session    *sessions.Manager
)

// 面板全局变量
var (
	Root    string
	Version string
	Locale  string
)
