package job

import (
	"runtime"
	"runtime/debug"

	"go.uber.org/zap"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
)

// PanelTask 面板每日任务
type PanelTask struct {
	appRepo     biz.AppRepo
	backupRepo  biz.BackupRepo
	settingRepo biz.SettingRepo
}

func NewPanelTask() *PanelTask {
	return &PanelTask{
		appRepo:     data.NewAppRepo(),
		backupRepo:  data.NewBackupRepo(),
		settingRepo: data.NewSettingRepo(),
	}
}

func (r *PanelTask) Run() {
	app.Status = app.StatusMaintain

	// 优化数据库
	if err := app.Orm.Exec("VACUUM").Error; err != nil {
		app.Status = app.StatusFailed
		app.Logger.Error("优化面板数据库失败", zap.Error(err))
	}
	if err := app.Orm.Exec("PRAGMA wal_checkpoint(TRUNCATE);").Error; err != nil {
		app.Status = app.StatusFailed
		app.Logger.Error("优化面板数据库失败", zap.Error(err))
	}

	// 备份面板
	if err := r.backupRepo.Create(biz.BackupTypePanel, ""); err != nil {
		app.Logger.Error("备份面板失败", zap.Error(err))
	}

	// 清理备份
	path, err := r.backupRepo.GetPath("panel")
	if err == nil {
		if err = r.backupRepo.ClearExpired(path, "panel_", 10); err != nil {
			app.Logger.Error("清理面板备份失败", zap.Error(err))
		}
	}

	// 更新商店缓存
	if offline, err := r.settingRepo.GetBool(biz.SettingKeyOfflineMode); err == nil && !offline {
		if err = r.appRepo.UpdateCache(); err != nil {
			app.Logger.Error("更新商店缓存失败", zap.Error(err))
		}
	}

	// 回收内存
	runtime.GC()
	debug.FreeOSMemory()

	app.Status = app.StatusNormal
}
