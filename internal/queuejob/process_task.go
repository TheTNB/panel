package queuejob

import (
	"errors"
	"log/slog"

	"github.com/tnb-labs/panel/internal/biz"
	"github.com/tnb-labs/panel/pkg/shell"
)

// ProcessTask 处理面板任务
type ProcessTask struct {
	log      *slog.Logger
	taskRepo biz.TaskRepo
	taskID   uint
}

// NewProcessTask 实例化 ProcessTask
func NewProcessTask(log *slog.Logger, taskRepo biz.TaskRepo) *ProcessTask {
	return &ProcessTask{
		log:      log,
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

	if _, err = shell.Exec(task.Shell); err != nil {
		return err
	}

	if err = r.taskRepo.UpdateStatus(taskID, biz.TaskStatusSuccess); err != nil {
		return err
	}

	return nil
}

func (r *ProcessTask) ErrHandle(err error) {
	r.log.Warn("background task failed", slog.Any("task_id", r.taskID), slog.Any("err", err))
	_ = r.taskRepo.UpdateStatus(r.taskID, biz.TaskStatusFailed)
}
