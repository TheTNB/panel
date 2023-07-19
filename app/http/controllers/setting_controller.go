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

func (r *SettingController) List(ctx http.Context) {
	var settings []models.Setting
	err := facades.Orm().Query().Get(&settings)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 查询设置列表失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")

		return
	}

	var result = make(map[string]string)
	for _, setting := range settings {
		if setting.Key == models.SettingKeyMysqlRootPassword {
			continue
		}

		result[setting.Key] = setting.Value
	}

	var user models.User
	err = facades.Auth().User(ctx, &user)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 获取用户失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")

		return
	}

	result["username"] = user.Username
	result["email"] = user.Email

	Success(ctx, result)
}

func (r *SettingController) Save(ctx http.Context) {
	name := ctx.Request().Input("name")
	port := ctx.Request().Input("port")
	backupPath := ctx.Request().Input("backup_path")
	websitePath := ctx.Request().Input("website_path")
	panelEntrance := ctx.Request().Input("panel_entrance")
	username := ctx.Request().Input("username")
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")

	err := r.setting.Set(models.SettingKeyName, name)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")

		return
	}
	oldPort := tools.ExecShell("cat /www/panel/panel.conf | grep APP_PORT | awk -F '=' '{print $2}'")
	if oldPort != port {
		tools.ExecShell("sed -i 's/APP_PORT=" + oldPort + "/APP_PORT=" + port + "/g' /www/panel/panel.conf")
	}
	err = r.setting.Set(models.SettingKeyBackupPath, backupPath)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")

		return
	}
	err = r.setting.Set(models.SettingKeyWebsitePath, websitePath)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")

		return
	}
	err = r.setting.Set(models.SettingKeyPanelEntrance, panelEntrance)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")

		return
	}

	var user models.User
	err = facades.Auth().User(ctx, &user)
	if err != nil {
		facades.Log().Error("[面板][SettingController] 获取用户失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")

		return
	}

	user.Username = username
	user.Email = email
	if len(password) > 0 {
		hash, err := facades.Hash().Make(password)
		if err != nil {
			facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
			Error(ctx, http.StatusInternalServerError, "系统内部错误")

			return
		}
		user.Password = hash
	}

	if err = facades.Orm().Query().Save(&user); err != nil {
		facades.Log().Error("[面板][SettingController] 保存设置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")

		return
	}

	Success(ctx, nil)
}
