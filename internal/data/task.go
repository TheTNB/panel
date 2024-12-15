package data

import (
	"fmt"
	"github.com/TheTNB/panel/pkg/queue"
	"gorm.io/gorm"
	"log/slog"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/queuejob"
)

type taskRepo struct {
	db    *gorm.DB
	log   *slog.Logger
	queue *queue.Queue
}

func NewTaskRepo(db *gorm.DB, log *slog.Logger, queue *queue.Queue) biz.TaskRepo {
	return &taskRepo{
		db:    db,
		log:   log,
		queue: queue,
	}
}

func (r *taskRepo) HasRunningTask() bool {
	var count int64
	r.db.Model(&biz.Task{}).Where("status = ?", biz.TaskStatusRunning).Or("status = ?", biz.TaskStatusWaiting).Count(&count)
	return count > 0
}

func (r *taskRepo) List(page, limit uint) ([]*biz.Task, int64, error) {
	var tasks []*biz.Task
	var total int64
	err := r.db.Model(&biz.Task{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&tasks).Error
	return tasks, total, err
}

func (r *taskRepo) Get(id uint) (*biz.Task, error) {
	task := new(biz.Task)
	err := r.db.Model(&biz.Task{}).Where("id = ?", id).First(task).Error
	return task, err
}

func (r *taskRepo) Delete(id uint) error {
	return r.db.Model(&biz.Task{}).Where("id = ?", id).Delete(&biz.Task{}).Error
}

func (r *taskRepo) UpdateStatus(id uint, status biz.TaskStatus) error {
	return r.db.Model(&biz.Task{}).Where("id = ?", id).Update("status", status).Error
}

func (r *taskRepo) Push(task *biz.Task) error {
	var count int64
	if err := r.db.Model(&biz.Task{}).Where("shell = ? and (status = ? or status = ?)", task.Shell, biz.TaskStatusWaiting, biz.TaskStatusRunning).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("duplicate submission, please wait for the previous task to end")
	}

	if err := r.db.Create(task).Error; err != nil {
		return err
	}

	return r.queue.Push(queuejob.NewProcessTask(r.log, r), []any{
		task.ID,
	})
}

// TODO fix
func (r *taskRepo) DispatchWaiting() {
	// cli下不处理
	if app.IsCli {
		return
	}

	if err := r.db.Model(&biz.Task{}).Where("status = ?", biz.TaskStatusRunning).Update("status", biz.TaskStatusFailed).Error; err != nil {
		r.log.Warn("failed to mark running tasks as failed", slog.Any("err", err))
		return
	}

	var tasks []biz.Task
	if err := r.db.Where("status = ?", biz.TaskStatusWaiting).Find(&tasks).Error; err != nil {
		r.log.Warn("failed to get pending tasks", slog.Any("err", err))
		return
	}

	for _, task := range tasks {
		if err := r.queue.Push(queuejob.NewProcessTask(r.log, r), []any{
			task.ID,
		}); err != nil {
			r.log.Warn("failed to push task", slog.Any("err", err))
			return
		}
	}
}
