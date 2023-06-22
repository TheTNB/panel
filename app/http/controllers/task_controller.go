package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"panel/app/models"
)

type TaskController struct {
	//Dependent services
}

func NewTaskController() *TaskController {
	return &TaskController{
		//Inject services
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
