package data

import (
	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
)

type taskRepo struct{}

func NewTaskRepo() biz.TaskRepo {
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
