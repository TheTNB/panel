package app

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-rat/sessions"
	"github.com/go-rat/utils/crypt"
	"github.com/knadh/koanf/v2"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"github.com/TheTNB/panel/pkg/queue"
)

var (
	Conf    *koanf.Koanf
	Http    *chi.Mux
	Orm     *gorm.DB
	Session *sessions.Manager
	Cron    *cron.Cron
	Queue   *queue.Queue
	Logger  *slog.Logger
	Crypter crypt.Crypter
)

// 面板状态常量
const (
	StatusNormal = iota
	StatusMaintain
	StatusClosed
	StatusUpgrade
	StatusFailed
)

// 面板全局变量
var (
	Key    string
	Root   string
	Locale string
	IsCli  bool
	Status = StatusNormal
)

// 自动注入
var (
	Version    = "0.0.0"
	BuildTime  string
	CommitHash string
	GoVersion  string
	BuildID    string
	BuildUser  string
	BuildHost  string
)
