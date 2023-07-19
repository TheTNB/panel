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

func (r *TaskController) Status(ctx http.Context) {
	var task models.Task
	err := facades.Orm().Query().Where("status", models.TaskStatusWaiting).OrWhere("status", models.TaskStatusRunning).FirstOrFail(&task)
	if err == nil {
		Success(ctx, http.Json{
			"task": true,
		})
		return
	}

	Success(ctx, http.Json{
		"task": false,
	})
}

func (r *TaskController) List(ctx http.Context) {
	status := ctx.Request().Query("status")
	if len(status) == 0 {
		status = models.TaskStatusWaiting
	}

	var tasks []models.Task
	var total int64
	err := facades.Orm().Query().Where("status", status).Paginate(ctx.Request().QueryInt("page"), ctx.Request().QueryInt("limit"), &tasks, &total)
	if err != nil {
		facades.Log().Error("[面板][TaskController] 查询任务列表失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, http.Json{
		"total": total,
		"items": tasks,
	})
}

func (r *TaskController) Log(ctx http.Context) {
	var task models.Task
	err := facades.Orm().Query().Where("id", ctx.Request().QueryInt("id")).FirstOrFail(&task)
	if err != nil {
		facades.Log().Error("[面板][TaskController] 查询任务失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	log := tools.ExecShell("tail -n 30 " + task.Log)

	Success(ctx, log)
}

func (r *TaskController) Delete(ctx http.Context) {
	var task models.Task
	_, err := facades.Orm().Query().Where("id", ctx.Request().QueryInt("id")).Delete(&task)
	if err != nil {
		facades.Log().Error("[面板][TaskController] 删除任务失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}
