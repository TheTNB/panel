package app

import (
	"gorm.io/gorm"
)

var (
	Orm *gorm.DB
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
