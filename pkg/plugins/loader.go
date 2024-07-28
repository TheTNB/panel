package plugins

import (
	"errors"
	"runtime"
	"strings"
	"sync"

	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/app/http/middleware"
)

var plugins sync.Map

type Plugin struct {
	Name        string       // 插件名称
	Description string       // 插件描述
	Slug        string       // 插件标识
	Version     string       // 插件版本
	Requires    []string     // 依赖插件
	Excludes    []string     // 排除插件
	Install     string       // 安装命令
	Uninstall   string       // 卸载命令
	Update      string       // 更新命令
	OnEnable    func() error // 启用插件后执行的命令
	OnDisable   func() error // 禁用插件后执行的命令
}

func NewPlugin() (*Plugin, error) {
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
	slug := name[:b]

	if _, ok = plugins.Load(slug); ok {
		return nil, errors.New("plugin already exists")
	}

	instance := &Plugin{
		Slug: slug,
	}
	plugins.Store(slug, instance)

	return instance, instance.OnEnable()
}

// Route 注册路由
func (r *Plugin) Route(group func(router route.Router)) {
	facades.Route().Prefix("api/plugins/"+r.Slug).Middleware(middleware.Session(), middleware.MustInstall()).Group(group)
}
