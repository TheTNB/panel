package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/shell"
)

type TaskService struct {
	taskRepo biz.TaskRepo
}

func NewTaskService() *TaskService {
	return &TaskService{
		taskRepo: data.NewTaskRepo(),
	}
}

// Status
//
//	@Summary	是否有任务正在运行
//	@Tags		任务服务
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	SuccessResponse
//	@Router		/tasks/status [get]
func (s *TaskService) Status(w http.ResponseWriter, r *http.Request) {
	Success(w, chix.M{
		"task": s.taskRepo.HasRunningTask(),
	})
}

// List
//
//	@Summary	任务列表
//	@Tags		任务服务
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	SuccessResponse
//	@Router		/tasks [get]
func (s *TaskService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	tasks, total, err := s.taskRepo.List(req.Page, req.Limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, chix.M{
		"total": total,
		"items": tasks,
	})
}

// Get
//
//	@Summary	任务详情
//	@Tags		任务服务
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	SuccessResponse
//	@Router		/task/log [get]
func (s *TaskService) Get(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	task, err := s.taskRepo.Get(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log, err := shell.Execf(`tail -n 500 '` + task.Log + `'`)
	if err == nil {
		task.Log = log
	}

	Success(w, task)
}

// Delete
//
//	@Summary	删除任务
//	@Tags		任务服务
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	SuccessResponse
//	@Router		/task/delete [post]
func (s *TaskService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = s.taskRepo.Delete(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}
