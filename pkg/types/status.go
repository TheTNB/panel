package types

// 定义面板状态常量
const (
	StatusNormal = iota
	StatusMaintain
	StatusClosed
	StatusUpgrade
	StatusFailed
)

var Status = StatusNormal
