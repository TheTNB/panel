package plugins

import (
	"context"
	"runtime"
	"strings"
	"sync"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
)

var plugins sync.Map

type Meta[ctx any] struct {
	Name        string    // 插件名称
	Description string    // 插件描述
	Slug        string    // 插件标识
	Version     string    // 插件版本
	Requires    []string  // 依赖插件
	Excludes    []string  // 排除插件
	Install     string    // 安装命令
	Uninstall   string    // 卸载命令
	Update      string    // 更新命令
	OnEnable    func(ctx) // 启用插件后执行的命令
	OnDisable   func(ctx) // 禁用插件后执行的命令
}

func AutoRegister(o *Meta[*context.Context]) foundation.Application {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("unable to get caller")
	}
	name := runtime.FuncForPC(pc).Name()
	a := strings.LastIndex(name, "/")
	if a < 0 {
		panic("invalid package name: " + name)
	}
	name = name[a+1:]
	b := strings.Index(name, ".")
	if b < 0 {
		panic("invalid package name: " + name)
	}
	name = name[:b]
	return Register(name, o)
}

// Register 注册插件控制器
func Register(service string, o *Meta[*context.Context]) foundation.Application {
	plugins.Store(service, o)
	return facades.App()
}
