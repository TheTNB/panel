package app

import (
	"github.com/go-chi/chi/v5"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-rat/sessions"
	"github.com/knadh/koanf/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/pkg/queue"
)

var (
	Conf       *koanf.Koanf
	Http       *chi.Mux
	Orm        *gorm.DB
	Validator  *validator.Validate
	Translator *ut.Translator
	Session    *sessions.Manager
	Cron       *cron.Cron
	Queue      *queue.Queue
	Logger     *zap.Logger
)

// 定义面板状态常量
const (
	StatusNormal = iota
	StatusMaintain
	StatusClosed
	StatusUpgrade
	StatusFailed
)

// 面板全局变量
var (
	Root    string
	Version string
	Locale  string
	IsCli   bool
	Status  = StatusNormal
)
