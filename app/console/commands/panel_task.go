package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	"panel/internal"
	"panel/pkg/tools"
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
	return "[面板] 每日任务"
}

// Extend The console command extend.
func (receiver *PanelTask) Extend() command.Extend {
	return command.Extend{
		Category: "panel",
	}
}

// Handle Execute the console command.
func (receiver *PanelTask) Handle(ctx console.Context) error {
	internal.Status = internal.StatusMaintain

	// 优化数据库
	if _, err := facades.Orm().Query().Exec("VACUUM"); err != nil {
		facades.Log().Tags("面板", "每日任务").
			With(map[string]any{
				"error": err.Error(),
			}).Error("优化面板数据库失败")
		return err
	}

	// 备份面板
	if err := tools.Archive([]string{"/www/panel"}, "/www/backup/panel/panel-"+carbon.Now().ToShortDateTimeString()+".zip"); err != nil {
		facades.Log().Tags("面板", "每日任务").
			With(map[string]any{
				"error": err.Error(),
			}).Error("备份面板失败")
		return err
	}

	// 清理 7 天前的备份
	if _, err := tools.Exec(`find /www/backup/panel -mtime +7 -name "*.zip" -exec rm -rf {} \;`); err != nil {
		facades.Log().Tags("面板", "每日任务").
			With(map[string]any{
				"error": err.Error(),
			}).Error("清理面板备份失败")
		return err
	}

	internal.Status = internal.StatusNormal
	return nil
}
