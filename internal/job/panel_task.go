package job

import (
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"

	"go.uber.org/zap"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/types"
)

// PanelTask 面板每日任务
type PanelTask struct {
	appRepo biz.AppRepo
}

func NewPanelTask() *PanelTask {
	return &PanelTask{
		appRepo: data.NewAppRepo(),
	}
}

func (receiver *PanelTask) Run() {
	types.Status = types.StatusMaintain

	// 优化数据库
	if err := app.Orm.Exec("VACUUM").Error; err != nil {
		types.Status = types.StatusFailed
		app.Logger.Error("优化面板数据库失败", zap.Error(err))
	}

	// 备份面板
	if err := io.Compress([]string{"/www/panel"}, filepath.Join(app.Root, "backup", "panel", "panel-"+time.Now().Format(time.DateOnly)+".zip"), io.Zip); err != nil {
		types.Status = types.StatusFailed
		app.Logger.Error("备份面板失败", zap.Error(err))
	}

	// 清理 7 天前的备份
	if _, err := shell.Execf(`find %s -mtime +7 -name "*.zip" -exec rm -rf {} \;`, filepath.Join(app.Root, "backup", "panel")); err != nil {
		types.Status = types.StatusFailed
		app.Logger.Error("清理面板备份失败", zap.Error(err))
	}

	// 更新商店缓存
	if err := receiver.appRepo.UpdateCache(); err != nil {
		app.Logger.Error("更新商店缓存失败", zap.Error(err))
	}

	// 回收内存
	runtime.GC()
	debug.FreeOSMemory()

	types.Status = types.StatusNormal
}
