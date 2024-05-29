package services

import (
	"github.com/goravel/framework/facades"

	"panel/app/jobs"
	"panel/app/models"
)

type TaskImpl struct {
}

func NewTaskImpl() *TaskImpl {
	return &TaskImpl{}
}

func (r *TaskImpl) Process(taskID uint) {
	err := facades.Queue().Job(&jobs.ProcessTask{}, []any{taskID}).Dispatch()
	if err != nil {
		facades.Log().Info("[面板][TaskService] 运行任务失败: " + err.Error())
		return
	}
}

func (r *TaskImpl) DispatchWaiting() error {
	var tasks []models.Task
	if err := facades.Orm().Query().Where("status = ?", models.TaskStatusWaiting).Find(&tasks); err != nil {
		return err
	}

	for _, task := range tasks {
		r.Process(task.ID)
	}

	return nil
}
