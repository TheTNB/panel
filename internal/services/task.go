package services

import (
	"sync"

	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/app/jobs"
	"github.com/TheTNB/panel/v2/app/models"
)

var taskMap sync.Map

type TaskImpl struct {
}

func NewTaskImpl() *TaskImpl {
	return &TaskImpl{}
}

func (r *TaskImpl) Process(taskID uint) error {
	taskMap.Store(taskID, true)
	return facades.Queue().Job(&jobs.ProcessTask{}, []any{taskID}).Dispatch()
}

func (r *TaskImpl) DispatchWaiting() error {
	var tasks []models.Task
	if err := facades.Orm().Query().Where("status = ?", models.TaskStatusWaiting).Find(&tasks); err != nil {
		return err
	}

	for _, task := range tasks {
		if _, ok := taskMap.Load(task.ID); ok {
			continue
		}
		if err := r.Process(task.ID); err != nil {
			return err
		}
	}

	return nil
}
