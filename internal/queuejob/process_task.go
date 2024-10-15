package queuejob

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

func (r *ProcessTask) Handle(args ...any) error {
	taskID, ok := args[0].(uint)
	if !ok {
		return errors.New("参数错误")
	}
	r.taskID = taskID

	task, err := r.taskRepo.Get(taskID)
	if err != nil {
		return err
	}

	if err = r.taskRepo.UpdateStatus(taskID, biz.TaskStatusRunning); err != nil {
		return err
	}

	if _, err = shell.Execf(task.Shell); err != nil { // nolint: govet
		return err
	}

	if err = r.taskRepo.UpdateStatus(taskID, biz.TaskStatusSuccess); err != nil {
		return err
	}

	return nil
}

func (r *ProcessTask) ErrHandle(err error) {
	_ = r.taskRepo.UpdateStatus(r.taskID, biz.TaskStatusFailed)
}
