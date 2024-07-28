package controllers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	commonrequests "github.com/TheTNB/panel/v2/app/http/requests/common"
	requests "github.com/TheTNB/panel/v2/app/http/requests/website"
	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/str"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
)

type WebsiteController struct {
	website internal.Website
	setting internal.Setting
	backup  internal.Backup
}

func NewWebsiteController() *WebsiteController {
	return &WebsiteController{
		website: services.NewWebsiteImpl(),
		setting: services.NewSettingImpl(),
		backup:  services.NewBackupImpl(),
	}
}

// List
//
//	@Summary	获取网站列表
//	@Tags		网站
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	query		commonrequests.Paginate	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/websites [get]
func (r *WebsiteController) List(ctx http.Context) http.Response {
	var paginateRequest commonrequests.Paginate
	sanitize := h.SanitizeRequest(ctx, &paginateRequest)
	if sanitize != nil {
		return sanitize
	}

	total, websites, err := r.website.List(paginateRequest.Page, paginateRequest.Limit)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"error": err.Error(),
		}).Info("获取网站列表失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, http.Json{
		"total": total,
		"items": websites,
	})
}

// Add
//
//	@Summary	添加网站
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		requests.Add	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/websites [post]
func (r *WebsiteController) Add(ctx http.Context) http.Response {
	var addRequest requests.Add
	sanitize := h.SanitizeRequest(ctx, &addRequest)
	if sanitize != nil {
		return sanitize
	}

	if len(addRequest.Path) == 0 {
		addRequest.Path = r.setting.Get(models.SettingKeyWebsitePath) + "/" + addRequest.Name
	}

	_, err := r.website.Add(addRequest)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"error": err.Error(),
		}).Info("添加网站失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// Delete
//
//	@Summary	删除网站
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		requests.Delete	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/websites/delete [post]
func (r *WebsiteController) Delete(ctx http.Context) http.Response {
	var deleteRequest requests.Delete
	sanitize := h.SanitizeRequest(ctx, &deleteRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := r.website.Delete(deleteRequest); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"id":    deleteRequest.ID,
			"error": err.Error(),
		}).Info("删除网站失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// GetDefaultConfig
//
//	@Summary	获取默认配置
//	@Tags		网站
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	SuccessResponse{data=map[string]string}
//	@Router		/panel/website/defaultConfig [get]
func (r *WebsiteController) GetDefaultConfig(ctx http.Context) http.Response {
	index, err := io.Read("/www/server/openresty/html/index.html")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	stop, err := io.Read("/www/server/openresty/html/stop.html")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, http.Json{
		"index": index,
		"stop":  stop,
	})
}

// SaveDefaultConfig
//
//	@Summary	保存默认配置
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		map[string]string	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/website/defaultConfig [post]
func (r *WebsiteController) SaveDefaultConfig(ctx http.Context) http.Response {
	index := ctx.Request().Input("index")
	stop := ctx.Request().Input("stop")

	if err := io.Write("/www/server/openresty/html/index.html", index, 0644); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"error": err.Error(),
		}).Info("保存默认首页配置失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := io.Write("/www/server/openresty/html/stop.html", stop, 0644); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"error": err.Error(),
		}).Info("保存默认停止页配置失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// GetConfig
//
//	@Summary	获取网站配置
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse{data=types.WebsiteAdd}
//	@Router		/panel/websites/{id}/config [get]
func (r *WebsiteController) GetConfig(ctx http.Context) http.Response {
	var idRequest requests.ID
	sanitize := h.SanitizeRequest(ctx, &idRequest)
	if sanitize != nil {
		return sanitize
	}

	config, err := r.website.GetConfig(idRequest.ID)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"id":    idRequest.ID,
			"error": err.Error(),
		}).Info("获取网站配置失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, config)
}

// SaveConfig
//
//	@Summary	保存网站配置
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		id		path		int					true	"网站 ID"
//	@Param		data	body		requests.SaveConfig	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/websites/{id}/config [post]
func (r *WebsiteController) SaveConfig(ctx http.Context) http.Response {
	var saveConfigRequest requests.SaveConfig
	sanitize := h.SanitizeRequest(ctx, &saveConfigRequest)
	if sanitize != nil {
		return sanitize
	}

	err := r.website.SaveConfig(saveConfigRequest)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ClearLog
//
//	@Summary	清空网站日志
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/websites/{id}/log [delete]
func (r *WebsiteController) ClearLog(ctx http.Context) http.Response {
	var idRequest requests.ID
	sanitize := h.SanitizeRequest(ctx, &idRequest)
	if sanitize != nil {
		return sanitize
	}

	website := models.Website{}
	err := facades.Orm().Query().Where("id", idRequest.ID).Get(&website)
	if err != nil {
		return h.ErrorSystem(ctx)
	}

	if err := io.Remove("/www/wwwlogs/" + website.Name + ".log"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// UpdateRemark
//
//	@Summary	更新网站备注
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/websites/{id}/updateRemark [post]
func (r *WebsiteController) UpdateRemark(ctx http.Context) http.Response {
	var idRequest requests.ID
	sanitize := h.SanitizeRequest(ctx, &idRequest)
	if sanitize != nil {
		return sanitize
	}

	website := models.Website{}
	err := facades.Orm().Query().Where("id", idRequest.ID).Get(&website)
	if err != nil {
		return h.ErrorSystem(ctx)
	}

	website.Remark = ctx.Request().Input("remark")
	if err = facades.Orm().Query().Save(&website); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"id":    idRequest.ID,
			"error": err.Error(),
		}).Info("更新网站备注失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, nil)
}

// BackupList
//
//	@Summary	获取网站备份列表
//	@Tags		网站
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	query		commonrequests.Paginate	true	"request"
//	@Success	200		{object}	SuccessResponse{data=[]types.BackupFile}
//	@Router		/panel/website/backupList [get]
func (r *WebsiteController) BackupList(ctx http.Context) http.Response {
	var paginateRequest commonrequests.Paginate
	sanitize := h.SanitizeRequest(ctx, &paginateRequest)
	if sanitize != nil {
		return sanitize
	}

	backups, err := r.backup.WebsiteList()
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"error": err.Error(),
		}).Info("获取备份列表失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	paged, total := h.Paginate(ctx, backups)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// CreateBackup
//
//	@Summary	创建网站备份
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/websites/{id}/createBackup [post]
func (r *WebsiteController) CreateBackup(ctx http.Context) http.Response {
	var idRequest requests.ID
	sanitize := h.SanitizeRequest(ctx, &idRequest)
	if sanitize != nil {
		return sanitize
	}

	website := models.Website{}
	if err := facades.Orm().Query().Where("id", idRequest.ID).Get(&website); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"id":    idRequest.ID,
			"error": err.Error(),
		}).Info("获取网站信息失败")
		return h.ErrorSystem(ctx)
	}

	if err := r.backup.WebSiteBackup(website); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"id":    idRequest.ID,
			"error": err.Error(),
		}).Info("备份网站失败")
		return h.Error(ctx, http.StatusInternalServerError, "备份网站失败: "+err.Error())
	}

	return h.Success(ctx, nil)
}

// UploadBackup
//
//	@Summary	上传网站备份
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		file	formData	file	true	"备份文件"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/website/uploadBackup [put]
func (r *WebsiteController) UploadBackup(ctx http.Context) http.Response {
	file, err := ctx.Request().File("file")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "上传文件失败")
	}

	backupPath := r.setting.Get(models.SettingKeyBackupPath) + "/website"
	if !io.Exists(backupPath) {
		if err = io.Mkdir(backupPath, 0644); err != nil {
			return nil
		}
	}

	name := file.GetClientOriginalName()
	_, err = file.StoreAs(backupPath, name)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"error": err.Error(),
		}).Info("上传备份失败")
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// RestoreBackup
//
//	@Summary	还原网站备份
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/websites/{id}/restoreBackup [post]
func (r *WebsiteController) RestoreBackup(ctx http.Context) http.Response {
	var restoreBackupRequest requests.RestoreBackup
	sanitize := h.SanitizeRequest(ctx, &restoreBackupRequest)
	if sanitize != nil {
		return sanitize
	}

	website := models.Website{}
	if err := facades.Orm().Query().Where("id", restoreBackupRequest.ID).Get(&website); err != nil {
		return h.ErrorSystem(ctx)
	}

	if err := r.backup.WebsiteRestore(website, restoreBackupRequest.Name); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"id":    restoreBackupRequest.ID,
			"file":  restoreBackupRequest.Name,
			"error": err.Error(),
		}).Info("还原网站失败")
		return h.Error(ctx, http.StatusInternalServerError, "还原网站失败: "+err.Error())
	}

	return h.Success(ctx, nil)
}

// DeleteBackup
//
//	@Summary	删除网站备份
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		data	body		requests.DeleteBackup	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/panel/website/deleteBackup [delete]
func (r *WebsiteController) DeleteBackup(ctx http.Context) http.Response {
	var deleteBackupRequest requests.DeleteBackup
	sanitize := h.SanitizeRequest(ctx, &deleteBackupRequest)
	if sanitize != nil {
		return sanitize
	}

	backupPath := r.setting.Get(models.SettingKeyBackupPath) + "/website"
	if !io.Exists(backupPath) {
		if err := io.Mkdir(backupPath, 0644); err != nil {
			return nil
		}
	}

	if err := io.Remove(backupPath + "/" + deleteBackupRequest.Name); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// ResetConfig
//
//	@Summary	重置网站配置
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/websites/{id}/resetConfig [post]
func (r *WebsiteController) ResetConfig(ctx http.Context) http.Response {
	var idRequest requests.ID
	sanitize := h.SanitizeRequest(ctx, &idRequest)
	if sanitize != nil {
		return sanitize
	}

	website := models.Website{}
	if err := facades.Orm().Query().Where("id", idRequest.ID).Get(&website); err != nil {
		return h.ErrorSystem(ctx)
	}

	website.Status = true
	website.SSL = false
	if err := facades.Orm().Query().Save(&website); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "网站管理").With(map[string]any{
			"id":    idRequest.ID,
			"error": err.Error(),
		}).Info("保存网站配置失败")
		return h.ErrorSystem(ctx)
	}

	raw := fmt.Sprintf(`
# 配置文件中的标记位请勿随意修改，改错将导致面板无法识别！
# 有自定义配置需求的，请将自定义的配置写在各标记位下方。
server
{
    # port标记位开始
    listen 80;
    # port标记位结束
    # server_name标记位开始
    server_name localhost;
    # server_name标记位结束
    # index标记位开始
    index index.php index.html;
    # index标记位结束
    # root标记位开始
    root %s;
    # root标记位结束

    # ssl标记位开始
    # ssl标记位结束

    # php标记位开始
    include enable-php-%d.conf;
    # php标记位结束

    # waf标记位开始
    waf off;
    waf_rule_path /www/server/openresty/ngx_waf/assets/rules/;
    waf_mode DYNAMIC;
    waf_cc_deny rate=1000r/m duration=60m;
    waf_cache capacity=50;
    # waf标记位结束

    # 错误页配置，可自行设置
    error_page 404 /404.html;
    #error_page 502 /502.html;

    # acme证书签发配置，不可修改
    include /www/server/vhost/acme/%s.conf;

    # 伪静态规则引入，修改后将导致面板设置的伪静态规则失效
    include /www/server/vhost/rewrite/%s.conf;

    # 面板默认禁止访问部分敏感目录，可自行修改
    location ~ ^/(\.user.ini|\.htaccess|\.git|\.svn)
    {
        return 404;
    }
    # 面板默认不记录静态资源的访问日志并开启1小时浏览器缓存，可自行修改
    location ~ .*\.(js|css)$
    {
        expires 1h;
        error_log /dev/null;
        access_log /dev/null;
    }

    access_log /www/wwwlogs/%s.log;
    error_log /www/wwwlogs/%s.log;
}

`, website.Path, website.PHP, website.Name, website.Name, website.Name, website.Name)
	if err := io.Write("/www/server/vhost/"+website.Name+".conf", raw, 0644); err != nil {
		return nil
	}
	if err := io.Write("/www/server/vhost/rewrite/"+website.Name+".conf", "", 0644); err != nil {
		return nil
	}
	if err := io.Write("/www/server/vhost/acme/"+website.Name+".conf", "", 0644); err != nil {
		return nil
	}
	if err := systemctl.Reload("openresty"); err != nil {
		_, err = shell.Execf("openresty -t")
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("重载OpenResty失败: %v", err))
	}

	return h.Success(ctx, nil)
}

// Status
//
//	@Summary	获取网站状态
//	@Tags		网站
//	@Accept		json
//	@Produce	json
//	@Security	BearerToken
//	@Param		id	path		int	true	"网站 ID"
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/websites/{id}/status [post]
func (r *WebsiteController) Status(ctx http.Context) http.Response {
	var idRequest requests.ID
	sanitize := h.SanitizeRequest(ctx, &idRequest)
	if sanitize != nil {
		return sanitize
	}

	website := models.Website{}
	if err := facades.Orm().Query().Where("id", idRequest.ID).Get(&website); err != nil {
		return h.ErrorSystem(ctx)
	}

	website.Status = ctx.Request().InputBool("status")
	if err := facades.Orm().Query().Save(&website); err != nil {
		return h.ErrorSystem(ctx)
	}

	raw, err := io.Read("/www/server/vhost/" + website.Name + ".conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	// 运行目录
	rootConfig := str.Cut(raw, "# root标记位开始\n", "# root标记位结束")
	match := regexp.MustCompile(`root\s+(.+);`).FindStringSubmatch(rootConfig)
	if len(match) == 2 {
		if website.Status {
			root := regexp.MustCompile(`# root\s+(.+);`).FindStringSubmatch(rootConfig)
			raw = strings.ReplaceAll(raw, rootConfig, "    root "+root[1]+";\n    ")
		} else {
			raw = strings.ReplaceAll(raw, rootConfig, "    root /www/server/openresty/html;\n    # root "+match[1]+";\n    ")
		}
	}

	// 默认文件
	indexConfig := str.Cut(raw, "# index标记位开始\n", "# index标记位结束")
	match = regexp.MustCompile(`index\s+(.+);`).FindStringSubmatch(indexConfig)
	if len(match) == 2 {
		if website.Status {
			index := regexp.MustCompile(`# index\s+(.+);`).FindStringSubmatch(indexConfig)
			raw = strings.ReplaceAll(raw, indexConfig, "    index "+index[1]+";\n    ")
		} else {
			raw = strings.ReplaceAll(raw, indexConfig, "    index stop.html;\n    # index "+match[1]+";\n    ")
		}
	}

	if err = io.Write("/www/server/vhost/"+website.Name+".conf", raw, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err = systemctl.Reload("openresty"); err != nil {
		_, err = shell.Execf("openresty -t")
		return h.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("重载OpenResty失败: %v", err))
	}

	return h.Success(ctx, nil)
}
