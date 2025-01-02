package app

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
	Key    string         // 密钥
	Root   string         // 根目录
	Locale string         // 语言
	IsCli  bool           // 是否命令行
	Status = StatusNormal // 面板状态
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
