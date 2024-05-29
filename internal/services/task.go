package services

import (
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/app/jobs"
	"github.com/TheTNB/panel/app/models"
)

type TaskImpl struct {
}

func NewTaskImpl() *TaskImpl {
	return &TaskImpl{}
}

func (r *TaskImpl) Process(taskID uint) error {
	return facades.Queue().Job(&jobs.ProcessTask{}, []any{taskID}).Dispatch()
}

func (r *TaskImpl) DispatchWaiting() error {
	var tasks []models.Task
	if err := facades.Orm().Query().Where("status = ?", models.TaskStatusWaiting).Find(&tasks); err != nil {
		return err
	}

	for _, task := range tasks {
		if err := r.Process(task.ID); err != nil {
			return err
		}
	}

	return nil
}
