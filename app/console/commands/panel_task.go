package commands

import (
	"context"
	"runtime"
	"runtime/debug"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/types"
)

// PanelTask 面板每日任务
type PanelTask struct {
}

// Signature The name and signature of the console command.
func (receiver *PanelTask) Signature() string {
	return "panel:task"
}

// Description The console command description.
func (receiver *PanelTask) Description() string {
	return facades.Lang(context.Background()).Get("commands.panel:task.description")
}

// Extend The console command extend.
func (receiver *PanelTask) Extend() command.Extend {
	return command.Extend{
		Category: "panel",
	}
}

// Handle Execute the console command.
func (receiver *PanelTask) Handle(console.Context) error {
	types.Status = types.StatusMaintain

	// 优化数据库
	if _, err := facades.Orm().Query().Exec("VACUUM"); err != nil {
		types.Status = types.StatusFailed
		facades.Log().Tags("面板", "每日任务").
			With(map[string]any{
				"error": err.Error(),
			}).Error("优化面板数据库失败")
		return err
	}

	// 备份面板
	if err := io.Archive([]string{"/www/panel"}, "/www/backup/panel/panel-"+carbon.Now().ToShortDateTimeString()+".zip"); err != nil {
		types.Status = types.StatusFailed
		facades.Log().Tags("面板", "每日任务").
			With(map[string]any{
				"error": err.Error(),
			}).Error("备份面板失败")
		return err
	}

	// 清理 7 天前的备份
	if _, err := shell.Execf(`find /www/backup/panel -mtime +7 -name "*.zip" -exec rm -rf {} \;`); err != nil {
		types.Status = types.StatusFailed
		facades.Log().Tags("面板", "每日任务").
			With(map[string]any{
				"error": err.Error(),
			}).Error("清理面板备份失败")
		return err
	}

	// 回收内存
	runtime.GC()
	debug.FreeOSMemory()

	types.Status = types.StatusNormal
	return nil
}
