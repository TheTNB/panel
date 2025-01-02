package job

import (
	"log/slog"
	"math/rand/v2"
	"runtime"
	"runtime/debug"
	"time"

	"gorm.io/gorm"

	"github.com/tnb-labs/panel/internal/app"
	"github.com/tnb-labs/panel/internal/biz"
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
		r.log.Warn("[Panel Task] failed to vacuum database", slog.Any("err", err))
	}
	if err := r.db.Exec("PRAGMA wal_checkpoint(TRUNCATE);").Error; err != nil {
		app.Status = app.StatusFailed
		r.log.Warn("[Panel Task] failed to wal checkpoint database", slog.Any("err", err))
	}

	// 备份面板
	if err := r.backupRepo.Create(biz.BackupTypePanel, ""); err != nil {
		r.log.Warn("备份面板失败", slog.Any("err", err))
	}

	// 清理备份
	path, err := r.backupRepo.GetPath("panel")
	if err == nil {
		if err = r.backupRepo.ClearExpired(path, "panel_", 10); err != nil {
			r.log.Warn("[Panel Task] failed to clear backup", slog.Any("err", err))
		}
	}

	// 更新商店缓存
	time.AfterFunc(time.Duration(rand.IntN(300))*time.Second, func() {
		if offline, err := r.settingRepo.GetBool(biz.SettingKeyOfflineMode); err == nil && !offline {
			if err = r.cacheRepo.UpdateApps(); err != nil {
				r.log.Warn("[Panel Task] failed to update apps cache", slog.Any("err", err))
			}
		}
	})

	// 更新伪静态缓存
	time.AfterFunc(time.Duration(rand.IntN(300))*time.Second, func() {
		if offline, err := r.settingRepo.GetBool(biz.SettingKeyOfflineMode); err == nil && !offline {
			if err = r.cacheRepo.UpdateRewrites(); err != nil {
				r.log.Warn("[Panel Task] failed to update rewrites cache", slog.Any("err", err))
			}
		}
	})

	// 回收内存
	runtime.GC()
	debug.FreeOSMemory()

	app.Status = app.StatusNormal
}
