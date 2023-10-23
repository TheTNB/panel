package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"panel/pkg/tools"

	"panel/app/models"
	"panel/app/services"
)

type SettingController struct {
	setting services.Setting
}

func NewSettingController() *SettingController {
	return &SettingController{
		setting: services.NewSettingImpl(),
	}
}

// List 获取设置列表
func (r *SettingController) List(ctx http.Context) http.Response {
	var settings []models.Setting
	err := facades.Orm().Query().Get(&settings)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 查询设置列表失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	type data struct {
		Name        string `json:"name"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		Email       string `json:"email"`
		Port        string `json:"port"`
		Entrance    string `json:"entrance"`
		WebsitePath string `json:"website_path"`
		BackupPath  string `json:"backup_path"`
	}

	var result data
	result.Name = r.setting.Get(models.SettingKeyName)
	result.Entrance = r.setting.Get(models.SettingKeyEntrance)
	result.WebsitePath = r.setting.Get(models.SettingKeyWebsitePath)
	result.BackupPath = r.setting.Get(models.SettingKeyBackupPath)

	var user models.User
	err = facades.Auth().User(ctx, &user)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 获取用户失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}
	result.Username = user.Username
	result.Email = user.Email

	result.Port = tools.Exec(`cat /www/panel/panel.conf | grep APP_PORT | awk -F '=' '{print $2}' | tr -d '\n'`)

	return Success(ctx, result)
}

// Save 保存设置
func (r *SettingController) Save(ctx http.Context) http.Response {
	name := ctx.Request().Input("name")
	port := ctx.Request().Input("port")
	backupPath := ctx.Request().Input("backup_path")
	websitePath := ctx.Request().Input("website_path")
	entrance := ctx.Request().Input("entrance")
	username := ctx.Request().Input("username")
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")

	err := r.setting.Set(models.SettingKeyName, name)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}
	oldPort := tools.Exec(`cat /www/panel/panel.conf | grep APP_PORT | awk -F '=' '{print $2}' | tr -d '\n'`)
	if oldPort != port {
		tools.Exec("sed -i 's/APP_PORT=" + oldPort + "/APP_PORT=" + port + "/g' /www/panel/panel.conf")
	}
	if !tools.Exists(backupPath) {
		tools.Mkdir(backupPath, 0644)
	}
	err = r.setting.Set(models.SettingKeyBackupPath, backupPath)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}
	if !tools.Exists(websitePath) {
		tools.Mkdir(websitePath, 0755)
		tools.Chown(websitePath, "www", "www")
	}
	err = r.setting.Set(models.SettingKeyWebsitePath, websitePath)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}
	err = r.setting.Set(models.SettingKeyEntrance, entrance)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	var user models.User
	err = facades.Auth().User(ctx, &user)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 获取用户失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if len(username) > 0 {
		user.Username = username
	}
	if len(email) > 0 {
		user.Email = email
	}
	if len(password) > 0 {
		hash, err := facades.Hash().Make(password)
		if err != nil {
			facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
			return Error(ctx, http.StatusInternalServerError, "系统内部错误")
		}
		user.Password = hash
	}

	if err = facades.Orm().Query().Save(&user); err != nil {
		facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	return Success(ctx, nil)
}
