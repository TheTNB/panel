package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"panel/app/models"
	"panel/pkg/tools"
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
		return Success(ctx, http.Json{
			"task": true,
		})
	}

	return Success(ctx, http.Json{
		"task": false,
	})
}

// List 获取任务列表
func (r *TaskController) List(ctx http.Context) http.Response {
	status := ctx.Request().Query("status")
	if len(status) == 0 {
		status = models.TaskStatusWaiting
	}

	var tasks []models.Task
	var total int64
	err := facades.Orm().Query().Where("status", status).Paginate(ctx.Request().QueryInt("page"), ctx.Request().QueryInt("limit"), &tasks, &total)
	if err != nil {
		facades.Log().Error("[面板][TaskController] 查询任务列表失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	return Success(ctx, http.Json{
		"total": total,
		"items": tasks,
	})
}

// Log 获取任务日志
func (r *TaskController) Log(ctx http.Context) http.Response {
	var task models.Task
	err := facades.Orm().Query().Where("id", ctx.Request().QueryInt("id")).FirstOrFail(&task)
	if err != nil {
		facades.Log().Error("[面板][TaskController] 查询任务失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	log := tools.Exec("tail -n 30 " + task.Log)

	return Success(ctx, log)
}

// Delete 删除任务
func (r *TaskController) Delete(ctx http.Context) http.Response {
	var task models.Task
	_, err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).Delete(&task)
	if err != nil {
		facades.Log().Error("[面板][TaskController] 删除任务失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	return Success(ctx, nil)
}
