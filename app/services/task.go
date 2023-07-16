package services

import (
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"

	"panel/app/jobs"
)

type Task interface {
	Process(taskID uint)
}

type TaskImpl struct {
}

func NewTaskImpl() *TaskImpl {
	return &TaskImpl{}
}

func (r *TaskImpl) Process(taskID uint) {
	go func() {
		err := facades.Queue().Job(&jobs.ProcessTask{}, []queue.Arg{
			{Type: "uint", Value: taskID},
		}).Dispatch()
		if err != nil {
			facades.Log().Error("[面板][TaskService] 运行任务失败: " + err.Error())
			return
		}
	}()
}
