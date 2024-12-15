package job

import (
	"gorm.io/gorm"
	"log/slog"
	"math/rand/v2"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
)

// PanelTask 面板每日任务
type PanelTask struct {
	db          *gorm.DB
	log         *slog.Logger
	backupRepo  biz.BackupRepo
	cacheRepo   biz.CacheRepo
	settingRepo biz.SettingRepo
}

func NewPanelTask(db *gorm.DB, log *slog.Logger, backup biz.BackupRepo, cache biz.CacheRepo, setting biz.SettingRepo) *PanelTask {
	return &PanelTask{
		db:          db,
		log:         log,
		backupRepo:  backup,
		cacheRepo:   cache,
		settingRepo: setting,
	}
}

func (r *PanelTask) Run() {
	app.Status = app.StatusMaintain

	// 优化数据库
	if err := r.db.Exec("VACUUM").Error; err != nil {
		app.Status = app.StatusFailed
		r.log.Warn("优化面板数据库失败", slog.Any("err", err))
	}
	if err := r.db.Exec("PRAGMA wal_checkpoint(TRUNCATE);").Error; err != nil {
		app.Status = app.StatusFailed
		r.log.Warn("优化面板数据库失败", slog.Any("err", err))
	}

	// 备份面板
	if err := r.backupRepo.Create(biz.BackupTypePanel, ""); err != nil {
		r.log.Warn("备份面板失败", slog.Any("err", err))
	}

	// 清理备份
	path, err := r.backupRepo.GetPath("panel")
	if err == nil {
		if err = r.backupRepo.ClearExpired(path, "panel_", 10); err != nil {
			r.log.Warn("清理面板备份失败", slog.Any("err", err))
		}
	}

	// 更新商店缓存
	time.AfterFunc(time.Duration(rand.IntN(300))*time.Second, func() {
		if offline, err := r.settingRepo.GetBool(biz.SettingKeyOfflineMode); err == nil && !offline {
			if err = r.cacheRepo.UpdateApps(); err != nil {
				r.log.Warn("更新商店缓存失败", slog.Any("err", err))
			}
		}
	})

	// 更新伪静态缓存
	time.AfterFunc(time.Duration(rand.IntN(300))*time.Second, func() {
		if offline, err := r.settingRepo.GetBool(biz.SettingKeyOfflineMode); err == nil && !offline {
			if err = r.cacheRepo.UpdateRewrites(); err != nil {
				r.log.Warn("更新伪静态缓存失败", slog.Any("err", err))
			}
		}
	})

	// 回收内存
	runtime.GC()
	debug.FreeOSMemory()

	app.Status = app.StatusNormal
}
