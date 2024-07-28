package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/path"
	"github.com/spf13/cast"

	requests "github.com/TheTNB/panel/v2/app/http/requests/setting"
	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/cert"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/os"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/tools"
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
//	@Summary	设置列表
//	@Tags		面板设置
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/setting/list [get]
func (r *SettingController) List(ctx http.Context) http.Response {
	var settings []models.Setting
	err := facades.Orm().Query().Get(&settings)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("获取面板设置列表失败")
		return h.ErrorSystem(ctx)
	}

	userID := cast.ToUint(ctx.Value("user_id"))
	var user models.User
	if err = facades.Orm().Query().Where("id", userID).Get(&user); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取用户信息失败")
	}

	port, err := shell.Execf(`cat /www/panel/panel.conf | grep APP_PORT | awk -F '=' '{print $2}' | tr -d '\n'`)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("获取面板端口失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, http.Json{
		"name":         r.setting.Get(models.SettingKeyName),
		"language":     facades.Config().GetString("app.locale"),
		"entrance":     facades.Config().GetString("panel.entrance"),
		"ssl":          facades.Config().GetBool("panel.ssl"),
		"website_path": r.setting.Get(models.SettingKeyWebsitePath),
		"backup_path":  r.setting.Get(models.SettingKeyBackupPath),
		"username":     user.Username,
		"password":     "",
		"email":        user.Email,
		"port":         port,
	})
}

// Update
//
//	@Summary	更新设置
//	@Tags		面板设置
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		requests.Update	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/setting/update [post]
func (r *SettingController) Update(ctx http.Context) http.Response {
	var updateRequest requests.Update
	sanitize := h.SanitizeRequest(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.setting.Set(models.SettingKeyName, updateRequest.Name)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("保存面板名称失败")
		return h.ErrorSystem(ctx)
	}

	if !io.Exists(updateRequest.BackupPath) {
		if err = io.Mkdir(updateRequest.BackupPath, 0644); err != nil {
			return h.ErrorSystem(ctx)
		}
	}
	err = r.setting.Set(models.SettingKeyBackupPath, updateRequest.BackupPath)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("保存备份目录失败")
		return h.ErrorSystem(ctx)
	}
	if !io.Exists(updateRequest.WebsitePath) {
		if err = io.Mkdir(updateRequest.WebsitePath, 0755); err != nil {
			return h.ErrorSystem(ctx)
		}
		if err = io.Chown(updateRequest.WebsitePath, "www", "www"); err != nil {
			return h.ErrorSystem(ctx)
		}
	}
	err = r.setting.Set(models.SettingKeyWebsitePath, updateRequest.WebsitePath)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("保存建站目录失败")
		return h.ErrorSystem(ctx)
	}

	userID := cast.ToUint(ctx.Value("user_id"))
	var user models.User
	if err = facades.Orm().Query().Where("id", userID).Get(&user); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取用户信息失败")
	}

	user.Username = updateRequest.UserName
	user.Email = updateRequest.Email
	if len(updateRequest.Password) > 0 {
		hash, err := facades.Hash().Make(updateRequest.Password)
		if err != nil {
			return h.ErrorSystem(ctx)
		}
		user.Password = hash
	}
	if err = facades.Orm().Query().Save(&user); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("保存用户信息失败")
		return h.ErrorSystem(ctx)
	}

	oldPort, err := shell.Execf(`cat /www/panel/panel.conf | grep APP_PORT | awk -F '=' '{print $2}' | tr -d '\n'`)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("获取面板端口失败")
		return h.ErrorSystem(ctx)
	}

	port := cast.ToString(updateRequest.Port)
	if oldPort != port {
		if out, err := shell.Execf("sed -i 's/APP_PORT=%s/APP_PORT=%s/g' /www/panel/panel.conf", oldPort, port); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
		if os.IsRHEL() {
			if out, err := shell.Execf("firewall-cmd --remove-port=%s/tcp --permanent", oldPort); err != nil {
				return h.Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := shell.Execf("firewall-cmd --add-port=%s/tcp --permanent", port); err != nil {
				return h.Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := shell.Execf("firewall-cmd --reload"); err != nil {
				return h.Error(ctx, http.StatusInternalServerError, out)
			}
		} else {
			if out, err := shell.Execf("ufw delete allow %s/tcp", oldPort); err != nil {
				return h.Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := shell.Execf("ufw allow %s/tcp", port); err != nil {
				return h.Error(ctx, http.StatusInternalServerError, out)
			}
			if out, err := shell.Execf("ufw reload"); err != nil {
				return h.Error(ctx, http.StatusInternalServerError, out)
			}
		}
	}

	oldEntrance, err := shell.Execf(`cat /www/panel/panel.conf | grep APP_ENTRANCE | awk -F '=' '{print $2}' | tr -d '\n'`)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("获取面板入口失败")
		return h.ErrorSystem(ctx)
	}
	entrance := cast.ToString(updateRequest.Entrance)
	if oldEntrance != entrance {
		if out, err := shell.Execf("sed -i 's!APP_ENTRANCE=" + oldEntrance + "!APP_ENTRANCE=" + entrance + "!g' /www/panel/panel.conf"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	oldLanguage, err := shell.Execf(`cat /www/panel/panel.conf | grep APP_LOCALE | awk -F '=' '{print $2}' | tr -d '\n'`)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "面板设置").With(map[string]any{
			"error": err.Error(),
		}).Info("获取面板语言失败")
		return h.ErrorSystem(ctx)
	}
	if oldLanguage != updateRequest.Language {
		if out, err := shell.Execf("sed -i 's/APP_LOCALE=" + oldLanguage + "/APP_LOCALE=" + updateRequest.Language + "/g' /www/panel/panel.conf"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	if oldPort != port || oldEntrance != entrance || oldLanguage != updateRequest.Language {
		tools.RestartPanel()
	}

	return h.Success(ctx, nil)
}

// GetHttps
//
//	@Summary	获取面板 HTTPS 设置
//	@Tags		面板设置
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/setting/https [get]
func (r *SettingController) GetHttps(ctx http.Context) http.Response {
	certPath := facades.Config().GetString("http.tls.ssl.cert")
	keyPath := facades.Config().GetString("http.tls.ssl.key")
	crt, err := io.Read(certPath)
	if err != nil {
		return h.ErrorSystem(ctx)
	}
	key, err := io.Read(keyPath)
	if err != nil {
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, http.Json{
		"https": facades.Config().GetBool("panel.ssl"),
		"cert":  crt,
		"key":   key,
	})
}

// UpdateHttps
//
//	@Summary	更新面板 HTTPS 设置
//	@Tags		面板设置
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		requests.Https	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/setting/https [post]
func (r *SettingController) UpdateHttps(ctx http.Context) http.Response {
	var httpsRequest requests.Https
	sanitize := h.SanitizeRequest(ctx, &httpsRequest)
	if sanitize != nil {
		return sanitize
	}

	if httpsRequest.Https {
		if _, err := cert.ParseCert(httpsRequest.Cert); err != nil {
			return h.Error(ctx, http.StatusBadRequest, "证书格式错误")
		}
		if _, err := cert.ParseKey(httpsRequest.Key); err != nil {
			return h.Error(ctx, http.StatusBadRequest, "密钥格式错误")
		}
		if err := io.Write(path.Executable("storage/ssl.crt"), httpsRequest.Cert, 0700); err != nil {
			return h.ErrorSystem(ctx)
		}
		if err := io.Write(path.Executable("storage/ssl.key"), httpsRequest.Key, 0700); err != nil {
			return h.ErrorSystem(ctx)
		}
		if out, err := shell.Execf("sed -i 's/APP_SSL=false/APP_SSL=true/g' /www/panel/panel.conf"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	} else {
		if out, err := shell.Execf("sed -i 's/APP_SSL=true/APP_SSL=false/g' /www/panel/panel.conf"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	tools.RestartPanel()
	return h.Success(ctx, nil)
}
