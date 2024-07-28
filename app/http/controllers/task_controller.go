package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/shell"
)

type TaskController struct {
	// Dependent services
}

func NewTaskController() *TaskController {
	return &TaskController{
		// Inject services
	}
}

// Status 获取当前任务状态
func (r *TaskController) Status(ctx http.Context) http.Response {
	var task models.Task
	err := facades.Orm().Query().Where("status", models.TaskStatusWaiting).OrWhere("status", models.TaskStatusRunning).FirstOrFail(&task)
	if err == nil {
		return h.Success(ctx, http.Json{
			"task": true,
		})
	}

	return h.Success(ctx, http.Json{
		"task": false,
	})
}

// List 获取任务列表
func (r *TaskController) List(ctx http.Context) http.Response {
	var tasks []models.Task
	var total int64
	err := facades.Orm().Query().Order("id desc").Paginate(ctx.Request().QueryInt("page", 1), ctx.Request().QueryInt("limit", 10), &tasks, &total)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "任务中心").With(map[string]any{
			"error": err.Error(),
		}).Info("查询任务列表失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, http.Json{
		"total": total,
		"items": tasks,
	})
}

// Log 获取任务日志
func (r *TaskController) Log(ctx http.Context) http.Response {
	var task models.Task
	err := facades.Orm().Query().Where("id", ctx.Request().QueryInt("id")).FirstOrFail(&task)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "任务中心").With(map[string]any{
			"id":    ctx.Request().QueryInt("id"),
			"error": err.Error(),
		}).Info("查询任务失败")
		return h.ErrorSystem(ctx)
	}

	log, err := shell.Execf(`tail -n 500 '` + task.Log + `'`)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "日志已被清理")
	}

	return h.Success(ctx, log)
}

// Delete 删除任务
func (r *TaskController) Delete(ctx http.Context) http.Response {
	var task models.Task
	_, err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).Delete(&task)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "任务中心").With(map[string]any{
			"id":    ctx.Request().QueryInt("id"),
			"error": err.Error(),
		}).Info("删除任务失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, nil)
}
