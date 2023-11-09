package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"

	requests "panel/app/http/requests/setting"
	responses "panel/app/http/responses/setting"
	"panel/app/models"
	"panel/app/services"
	"panel/pkg/tools"
)

type SettingController struct {
	setting services.Setting
}

func NewSettingController() *SettingController {
	return &SettingController{
		setting: services.NewSettingImpl(),
	}
}

// List
// @Summary 设置列表
// @Description 获取面板设置列表
// @Tags 面板设置
// @Produce json
// @Security BearerToken
// @Success 200 {object} SuccessResponse{data=responses.Settings}
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/setting/list [get]
func (r *SettingController) List(ctx http.Context) http.Response {
	var settings []models.Setting
	err := facades.Orm().Query().Get(&settings)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 查询设置列表失败 ", err)
		return ErrorSystem(ctx)
	}

	var result responses.Settings
	result.Name = r.setting.Get(models.SettingKeyName)
	result.Entrance = facades.Config().GetString("http.entrance")
	result.WebsitePath = r.setting.Get(models.SettingKeyWebsitePath)
	result.BackupPath = r.setting.Get(models.SettingKeyBackupPath)

	var user models.User
	err = facades.Auth().User(ctx, &user)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 获取用户失败 ", err)
		return ErrorSystem(ctx)
	}
	result.Username = user.Username
	result.Email = user.Email

	result.Port = tools.Exec(`cat /www/panel/panel.conf | grep APP_PORT | awk -F '=' '{print $2}' | tr -d '\n'`)

	return Success(ctx, result)
}

// Update
// @Summary 更新设置
// @Description 更新面板设置
// @Tags 面板设置
// @Accept json
// @Produce json
// @Security BearerToken
// @Param data body requests.Update true "更新设置"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse "登录已过期"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /panel/setting/update [post]
func (r *SettingController) Update(ctx http.Context) http.Response {
	var updateRequest requests.Update
	sanitize := Sanitize(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.setting.Set(models.SettingKeyName, updateRequest.Name)
	if err != nil {
		facades.Log().Request(ctx.Request()).With(map[string]any{
			"error": err.Error(),
		}).Tags("面板", "面板设置").Error("保存面板名称失败")
		return ErrorSystem(ctx)
	}

	if !tools.Exists(updateRequest.BackupPath) {
		tools.Mkdir(updateRequest.BackupPath, 0644)
	}
	err = r.setting.Set(models.SettingKeyBackupPath, updateRequest.BackupPath)
	if err != nil {
		facades.Log().Request(ctx.Request()).With(map[string]any{
			"error": err.Error(),
		}).Tags("面板", "面板设置").Error("保存备份目录失败")
		return ErrorSystem(ctx)
	}
	if !tools.Exists(updateRequest.WebsitePath) {
		tools.Mkdir(updateRequest.WebsitePath, 0755)
		tools.Chown(updateRequest.WebsitePath, "www", "www")
	}
	err = r.setting.Set(models.SettingKeyWebsitePath, updateRequest.WebsitePath)
	if err != nil {
		facades.Log().Request(ctx.Request()).With(map[string]any{
			"error": err.Error(),
		}).Tags("面板", "面板设置").Error("保存建站目录失败")
		return ErrorSystem(ctx)
	}

	var user models.User
	err = facades.Auth().User(ctx, &user)
	if err != nil {
		return ErrorSystem(ctx)
	}
	user.Username = updateRequest.UserName
	user.Email = updateRequest.Email
	if len(updateRequest.Password) > 0 {
		hash, err := facades.Hash().Make(updateRequest.Password)
		if err != nil {
			return ErrorSystem(ctx)
		}
		user.Password = hash
	}
	if err = facades.Orm().Query().Save(&user); err != nil {
		facades.Log().Request(ctx.Request()).With(map[string]any{
			"error": err.Error(),
		}).Tags("面板", "面板设置").Error("保存用户信息失败")
		return ErrorSystem(ctx)
	}

	oldPort := tools.Exec(`cat /www/panel/panel.conf | grep APP_PORT | awk -F '=' '{print $2}' | tr -d '\n'`)
	port := cast.ToString(updateRequest.Port)
	if oldPort != port {
		tools.Exec("sed -i 's/APP_PORT=" + oldPort + "/APP_PORT=" + port + "/g' /www/panel/panel.conf")
	}
	oldEntrance := tools.Exec(`cat /www/panel/panel.conf | grep APP_ENTRANCE | awk -F '=' '{print $2}' | tr -d '\n'`)
	entrance := cast.ToString(updateRequest.Entrance)
	if oldEntrance != entrance {
		tools.Exec("sed -i 's!APP_ENTRANCE=" + oldEntrance + "!APP_ENTRANCE=" + entrance + "!g' /www/panel/panel.conf")
	}

	if oldPort != port || oldEntrance != entrance {
		tools.RestartPanel()
	}

	return Success(ctx, nil)
}
