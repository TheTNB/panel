package service

import (
	"net/http"
	"path/filepath"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/io"
)

type WebsiteService struct {
	websiteRepo biz.WebsiteRepo
	settingRepo biz.SettingRepo
}

func NewWebsiteService() *WebsiteService {
	return &WebsiteService{
		websiteRepo: data.NewWebsiteRepo(),
		settingRepo: data.NewSettingRepo(),
	}
}

// GetDefaultConfig
//
//	@Summary	获取默认配置
//	@Tags		网站服务
//	@Produce	json
//	@Success	200	{object}	SuccessResponse{data=map[string]string}
//	@Router		/panel/website/defaultConfig [get]
func (s *WebsiteService) GetDefaultConfig(w http.ResponseWriter, r *http.Request) {
	index, err := io.Read(filepath.Join(app.Root, "server/openresty/html/index.html"))
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	stop, err := io.Read(filepath.Join(app.Root, "server/openresty/html/stop.html"))
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, chix.M{
		"index": index,
		"stop":  stop,
	})
}

// UpdateDefaultConfig
//
//	@Summary	更新默认配置
//	@Tags		网站服务
//	@Accept		json
//	@Produce	json
//	@Param		data	body		map[string]string	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/website/defaultConfig [post]
func (s *WebsiteService) UpdateDefaultConfig(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteDefaultConfig](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.UpdateDefaultConfig(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

// List
//
//	@Summary	网站列表
//	@Tags		网站服务
//	@Produce	json
//	@Param		data	query		commonrequests.Paginate	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/websites [get]
func (s *WebsiteService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	websites, total, err := s.websiteRepo.List(req.Page, req.Limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, chix.M{
		"total": total,
		"items": websites,
	})
}

// Create
//
//	@Summary	创建网站
//	@Tags		网站服务
//	@Accept		json
//	@Produce	json
//	@Param		data	body		requests.Add	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/websites [post]
func (s *WebsiteService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if len(req.Path) == 0 {
		req.Path, _ = s.settingRepo.Get(biz.SettingKeyWebsitePath)
		req.Path = filepath.Join(req.Path, req.Name)
	}

	if _, err = s.websiteRepo.Create(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

// Get
//
//	@Summary	获取网站
//	@Tags		网站服务
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse{data=types.WebsiteAdd}
//	@Router		/panel/websites/{id}/config [get]
func (s *WebsiteService) Get(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	config, err := s.websiteRepo.Get(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, config)
}

// Update
//
//	@Summary	更新网站
//	@Tags		网站服务
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int					true	"网站 ID"
//	@Param		data	body		requests.SaveConfig	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/websites/{id}/config [post]
func (s *WebsiteService) Update(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteUpdate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.Update(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

// Delete
//
//	@Summary	删除网站
//	@Tags		网站服务
//	@Accept		json
//	@Produce	json
//	@Param		data	body		requests.Delete	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/websites/delete [post]
func (s *WebsiteService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteDelete](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.Delete(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

// ClearLog
//
//	@Summary	清空网站日志
//	@Tags		网站服务
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/websites/{id}/log [delete]
func (s *WebsiteService) ClearLog(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.ClearLog(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

// UpdateRemark
//
//	@Summary	更新网站备注
//	@Tags		网站服务
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/websites/{id}/updateRemark [post]
func (s *WebsiteService) UpdateRemark(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteUpdateRemark](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.UpdateRemark(req.ID, req.Remark); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

// ResetConfig
//
//	@Summary	重置网站配置
//	@Tags		网站服务
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/websites/{id}/resetConfig [post]
func (s *WebsiteService) ResetConfig(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.ResetConfig(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

// UpdateStatus
//
//	@Summary	更新网站状态
//	@Tags		网站服务
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/websites/{id}/status [post]
func (s *WebsiteService) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.WebsiteUpdateStatus](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.websiteRepo.UpdateStatus(req.ID, req.Status); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}
