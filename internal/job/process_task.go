package job

import (
	"errors"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/pkg/shell"
)

// ProcessTask 处理面板任务
type ProcessTask struct {
	taskRepo biz.TaskRepo
	taskID   uint
}

// NewProcessTask 实例化 ProcessTask
func NewProcessTask(taskRepo biz.TaskRepo) *ProcessTask {
	return &ProcessTask{
		taskRepo: taskRepo,
	}
}

func (receiver *ProcessTask) Handle(args ...any) error {
	taskID, ok := args[0].(uint)
	if !ok {
		return errors.New("参数错误")
	}
	receiver.taskID = taskID

	task, err := receiver.taskRepo.Get(taskID)
	if err != nil {
		return err
	}

	if err = receiver.taskRepo.UpdateStatus(taskID, biz.TaskStatusRunning); err != nil {
		return err
	}

	if _, err = shell.Execf(task.Shell); err != nil {
		return err
	}

	if err = receiver.taskRepo.UpdateStatus(taskID, biz.TaskStatusSuccess); err != nil {
		return err
	}

	return nil
}

func (receiver *ProcessTask) ErrHandle(err error) {
	_ = receiver.taskRepo.UpdateStatus(receiver.taskID, biz.TaskStatusFailed)
}
