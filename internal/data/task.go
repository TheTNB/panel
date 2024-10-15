package data

import (
	"sync"

	"go.uber.org/zap"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/queuejob"
)

var taskDispatchOnce sync.Once

type taskRepo struct{}

func NewTaskRepo() biz.TaskRepo {
	task := &taskRepo{}
	taskDispatchOnce.Do(task.DispatchWaiting)
	return &taskRepo{}
}

func (r *taskRepo) HasRunningTask() bool {
	var count int64
	app.Orm.Model(&biz.Task{}).Where("status = ?", biz.TaskStatusRunning).Or("status = ?", biz.TaskStatusWaiting).Count(&count)
	return count > 0
}

func (r *taskRepo) List(page, limit uint) ([]*biz.Task, int64, error) {
	var tasks []*biz.Task
	var total int64
	err := app.Orm.Model(&biz.Task{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&tasks).Error
	return tasks, total, err
}

func (r *taskRepo) Get(id uint) (*biz.Task, error) {
	task := new(biz.Task)
	err := app.Orm.Model(&biz.Task{}).Where("id = ?", id).First(task).Error
	return task, err
}

func (r *taskRepo) Delete(id uint) error {
	return app.Orm.Model(&biz.Task{}).Where("id = ?", id).Delete(&biz.Task{}).Error
}

func (r *taskRepo) UpdateStatus(id uint, status biz.TaskStatus) error {
	return app.Orm.Model(&biz.Task{}).Where("id = ?", id).Update("status", status).Error
}

func (r *taskRepo) Push(task *biz.Task) error {
	if err := app.Orm.Create(task).Error; err != nil {
		return err
	}
	return app.Queue.Push(queuejob.NewProcessTask(r), []any{
		task.ID,
	})
}

func (r *taskRepo) DispatchWaiting() {
	var tasks []biz.Task
	if err := app.Orm.Where("status = ?", biz.TaskStatusWaiting).Find(&tasks).Error; err != nil {
		app.Logger.Error("获取待处理任务失败", zap.Error(err))
		return
	}

	for _, task := range tasks {
		if err := app.Queue.Push(queuejob.NewProcessTask(r), []any{
			task.ID,
		}); err != nil {
			app.Logger.Error("推送任务失败", zap.Error(err))
			return
		}
	}
}
