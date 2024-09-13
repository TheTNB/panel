package service

import "net/http"

type TaskService struct {
}

func NewTaskService() *TaskService {
	return &TaskService{}
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

}

// Index
//
//	@Summary	任务列表
//	@Tags		任务服务
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	SuccessResponse
//	@Router		/tasks [get]
func (s *TaskService) Index(w http.ResponseWriter, r *http.Request) {

}

// Show
//
//	@Summary	任务详情
//	@Tags		任务服务
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	SuccessResponse
//	@Router		/task/log [get]
func (s *TaskService) Show(w http.ResponseWriter, r *http.Request) {

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

}
