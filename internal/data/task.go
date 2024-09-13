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
