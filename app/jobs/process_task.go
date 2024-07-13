package jobs

import (
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/pkg/shell"
)

// ProcessTask 处理面板任务
type ProcessTask struct {
}

// Signature The name and signature of the job.
func (receiver *ProcessTask) Signature() string {
	return "process_task"
}

// Handle Execute the job.
func (receiver *ProcessTask) Handle(args ...any) error {
	taskID, ok := args[0].(uint)
	if !ok {
		facades.Log().Tags("面板", "异步任务").With(map[string]any{
			"args": args,
		}).Infof("参数错误")
		return nil
	}

	var task models.Task
	_ = facades.Orm().Query().Where("id = ?", taskID).Get(&task)
	if task.ID == 0 {
		facades.Log().Tags("面板", "异步任务").With(map[string]any{
			"task_id": taskID,
		}).Infof("任务不存在")
		return nil
	}

	facades.Log().Tags("面板", "异步任务").With(map[string]any{
		"task_id": taskID,
	}).Infof("开始执行任务")

	task.Status = models.TaskStatusRunning
	if err := facades.Orm().Query().Save(&task); err != nil {
		facades.Log().Tags("面板", "异步任务").With(map[string]any{
			"task_id": taskID,
			"error":   err.Error(),
		}).Infof("更新任务状态失败")
		return nil
	}

	if _, err := shell.Execf(task.Shell); err != nil {
		task.Status = models.TaskStatusFailed
		_ = facades.Orm().Query().Save(&task)
		facades.Log().Tags("面板", "异步任务").With(map[string]any{
			"task_id": taskID,
			"error":   err.Error(),
		}).Infof("执行任务失败")
		return nil
	}

	task.Status = models.TaskStatusSuccess
	if err := facades.Orm().Query().Save(&task); err != nil {
		facades.Log().Tags("面板", "异步任务").With(map[string]any{
			"task_id": taskID,
			"error":   err.Error(),
		}).Infof("更新任务状态失败")
		return nil
	}

	facades.Log().Tags("面板", "异步任务").With(map[string]any{
		"task_id": taskID,
	}).Infof("执行任务成功")
	return nil
}
