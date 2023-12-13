package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"

	requests "panel/app/http/requests/setting"
	"panel/app/models"
	"panel/internal"
	"panel/internal/services"
	"panel/pkg/tools"
)

type SettingController struct {
	setting internal.Setting
}

func NewSettingController() *SettingController {
	return &SettingController{
		setting: services.NewSettingImpl(),
	}
}

// List
//
//	@Summary		设置列表
//	@Description	获取面板设置列表
//	@Tags			面板设置
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	SuccessResponse
//	@Router			/panel/setting/list [get]
func (r *SettingController) List(ctx http.Context) http.Response {
	var settings []models.Setting
	err := facades.Orm().Query().Get(&settings)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("获取面板设置列表失败")
		return ErrorSystem(ctx)
	}

	var user models.User
	err = facades.Auth().User(ctx, &user)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("获取用户信息失败")
		return ErrorSystem(ctx)
	}

	port, err := tools.Exec(`cat /www/panel/panel.conf | grep APP_PORT | awk -F '=' '{print $2}' | tr -d '\n'`)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("获取面板端口失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, http.Json{
		"name":         r.setting.Get(models.SettingKeyName),
		"entrance":     facades.Config().GetString("http.entrance"),
		"website_path": r.setting.Get(models.SettingKeyWebsitePath),
		"backup_path":  r.setting.Get(models.SettingKeyBackupPath),
		"user_name":    user.Username,
		"password":     "",
		"email":        user.Email,
		"port":         port,
	})
}

// Update
//
//	@Summary		更新设置
//	@Description	更新面板设置
//	@Tags			面板设置
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Update	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/setting/update [post]
func (r *SettingController) Update(ctx http.Context) http.Response {
	var updateRequest requests.Update
	sanitize := Sanitize(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.setting.Set(models.SettingKeyName, updateRequest.Name)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("保存面板名称失败")
		return ErrorSystem(ctx)
	}

	if !tools.Exists(updateRequest.BackupPath) {
		if err = tools.Mkdir(updateRequest.BackupPath, 0644); err != nil {
			return ErrorSystem(ctx)
		}
	}
	err = r.setting.Set(models.SettingKeyBackupPath, updateRequest.BackupPath)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("保存备份目录失败")
		return ErrorSystem(ctx)
	}
	if !tools.Exists(updateRequest.WebsitePath) {
		if err = tools.Mkdir(updateRequest.WebsitePath, 0755); err != nil {
			return ErrorSystem(ctx)
		}
		if err = tools.Chown(updateRequest.WebsitePath, "www", "www"); err != nil {
			return ErrorSystem(ctx)
		}
	}
	err = r.setting.Set(models.SettingKeyWebsitePath, updateRequest.WebsitePath)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("保存建站目录失败")
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
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("保存用户信息失败")
		return ErrorSystem(ctx)
	}

	oldPort, err := tools.Exec(`cat /www/panel/panel.conf | grep APP_PORT | awk -F '=' '{print $2}' | tr -d '\n'`)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("获取面板端口失败")
		return ErrorSystem(ctx)
	}

	port := cast.ToString(updateRequest.Port)
	if oldPort != port {
		if out, err := tools.Exec("sed -i 's/APP_PORT=" + oldPort + "/APP_PORT=" + port + "/g' /www/panel/panel.conf"); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
		if tools.IsRHEL() {
			if out, err := tools.Exec("firewall-cmd --remove-port=" + cast.ToString(port) + "/tcp --permanent 2>&1"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := tools.Exec("firewall-cmd --add-port=" + cast.ToString(port) + "/tcp --permanent 2>&1"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := tools.Exec("firewall-cmd --reload"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
		} else {
			if out, err := tools.Exec("ufw delete allow " + cast.ToString(port) + "/tcp"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := tools.Exec("ufw allow " + cast.ToString(port) + "/tcp"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := tools.Exec("ufw reload"); err != nil {
				return Error(ctx, http.StatusInternalServerError, out)
			}
		}
	}

	oldEntrance, err := tools.Exec(`cat /www/panel/panel.conf | grep APP_ENTRANCE | awk -F '=' '{print $2}' | tr -d '\n'`)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("获取面板入口失败")
		return ErrorSystem(ctx)
	}

	entrance := cast.ToString(updateRequest.Entrance)
	if oldEntrance != entrance {
		if out, err := tools.Exec("sed -i 's!APP_ENTRANCE=" + oldEntrance + "!APP_ENTRANCE=" + entrance + "!g' /www/panel/panel.conf"); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
	}

	if oldPort != port || oldEntrance != entrance {
		tools.RestartPanel()
	}

	return Success(ctx, nil)
}
