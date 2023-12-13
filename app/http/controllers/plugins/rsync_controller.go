package plugins

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"panel/app/http/controllers"
	commonrequests "panel/app/http/requests/common"
	requests "panel/app/http/requests/plugins/rsync"
	"panel/app/internal"
	"panel/app/internal/services"
	"panel/pkg/tools"
)

type RsyncController struct {
	setting internal.Setting
}

func NewRsyncController() *RsyncController {
	return &RsyncController{
		setting: services.NewSettingImpl(),
	}
}

// Status
//
//	@Summary		服务状态
//	@Description	获取 Rsync 服务状态
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/status [get]
func (r *RsyncController) Status(ctx http.Context) http.Response {
	status, err := tools.ServiceStatus("rsyncd")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取服务运行状态失败")
	}

	return controllers.Success(ctx, status)
}

// Restart
//
//	@Summary		重启服务
//	@Description	重启 Rsync 服务
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/restart [post]
func (r *RsyncController) Restart(ctx http.Context) http.Response {
	if err := tools.ServiceRestart("rsyncd"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "重启服务失败")
	}

	return controllers.Success(ctx, nil)
}

// Start
//
//	@Summary		启动服务
//	@Description	启动 Rsync 服务
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/start [post]
func (r *RsyncController) Start(ctx http.Context) http.Response {
	if err := tools.ServiceStart("rsyncd"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "启动服务失败")
	}

	return controllers.Success(ctx, nil)
}

// Stop
//
//	@Summary		停止服务
//	@Description	停止 Rsync 服务
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/stop [post]
func (r *RsyncController) Stop(ctx http.Context) http.Response {
	if err := tools.ServiceStop("rsyncd"); err != nil {
		return nil
	}
	status, err := tools.ServiceStatus("rsyncd")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取服务运行状态失败")
	}

	return controllers.Success(ctx, !status)
}

// List
//
//	@Summary		列出模块
//	@Description	列出所有 Rsync 模块
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		commonrequests.Paginate	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/modules [get]
func (r *RsyncController) List(ctx http.Context) http.Response {
	var paginateRequest commonrequests.Paginate
	sanitize := controllers.Sanitize(ctx, &paginateRequest)
	if sanitize != nil {
		return sanitize
	}

	config, err := tools.Read("/etc/rsyncd.conf")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	var modules []RsyncModule
	lines := strings.Split(config, "\n")
	var currentModule *RsyncModule

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if currentModule != nil {
				modules = append(modules, *currentModule)
			}
			moduleName := line[1 : len(line)-1]
			currentModule = &RsyncModule{
				Name: moduleName,
			}
		} else if currentModule != nil {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				switch key {
				case "path":
					currentModule.Path = value
				case "comment":
					currentModule.Comment = value
				case "read only":
					currentModule.ReadOnly = value == "yes" || value == "true"
				case "auth users":
					currentModule.AuthUser = value
					currentModule.Secret, err = tools.Exec("grep -E '^" + currentModule.AuthUser + ":.*$' /etc/rsyncd.secrets | awk -F ':' '{print $2}'")
					if err != nil {
						return controllers.Error(ctx, http.StatusInternalServerError, "获取模块"+currentModule.AuthUser+"的密钥失败")
					}
				case "hosts allow":
					currentModule.HostsAllow = value
				}
			}
		}
	}

	if currentModule != nil {
		modules = append(modules, *currentModule)
	}

	startIndex := (paginateRequest.Page - 1) * paginateRequest.Limit
	endIndex := paginateRequest.Page * paginateRequest.Limit
	if startIndex > len(modules) {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []RsyncModule{},
		})
	}
	if endIndex > len(modules) {
		endIndex = len(modules)
	}
	pagedModules := modules[startIndex:endIndex]
	if pagedModules == nil {
		pagedModules = []RsyncModule{}
	}

	return controllers.Success(ctx, http.Json{
		"total": len(modules),
		"items": pagedModules,
	})
}

// Create
//
//	@Summary		添加模块
//	@Description	添加 Rsync 模块
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Create	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/modules [post]
func (r *RsyncController) Create(ctx http.Context) http.Response {
	var createRequest requests.Create
	sanitize := controllers.Sanitize(ctx, &createRequest)
	if sanitize != nil {
		return sanitize
	}

	config, err := tools.Read("/etc/rsyncd.conf")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if strings.Contains(config, "["+createRequest.Name+"]") {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "模块 "+createRequest.Name+" 已存在")
	}

	conf := `# ` + createRequest.Name + `-START
[` + createRequest.Name + `]
path = ` + createRequest.Path + `
comment = ` + createRequest.Comment + `
read only = no
auth users = ` + createRequest.AuthUser + `
hosts allow = ` + createRequest.HostsAllow + `
secrets file = /etc/rsyncd.secrets
# ` + createRequest.Name + `-END
`

	if err := tools.WriteAppend("/etc/rsyncd.conf", conf); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if out, err := tools.Exec("echo '" + createRequest.AuthUser + ":" + createRequest.Secret + "' >> /etc/rsyncd.secrets"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	if err := tools.ServiceRestart("rsyncd"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

// Destroy
//
//	@Summary		删除模块
//	@Description	删除 Rsync 模块
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Param			name	path		string	true	"模块名称"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/modules/{name} [delete]
func (r *RsyncController) Destroy(ctx http.Context) http.Response {
	name := ctx.Request().Input("name")
	if len(name) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "name 不能为空")
	}

	config, err := tools.Read("/etc/rsyncd.conf")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if !strings.Contains(config, "["+name+"]") {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "模块 "+name+" 不存在")
	}

	module := tools.Cut(config, "# "+name+"-START", "# "+name+"-END")
	config = strings.Replace(config, "\n# "+name+"-START"+module+"# "+name+"-END", "", -1)

	match := regexp.MustCompile(`auth users = ([^\n]+)`).FindStringSubmatch(module)
	if len(match) == 2 {
		authUser := match[1]
		if out, err := tools.Exec("sed -i '/^" + authUser + ":.*$/d' /etc/rsyncd.secrets"); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	if err = tools.Write("/etc/rsyncd.conf", config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err = tools.ServiceRestart("rsyncd"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

// Update
//
//	@Summary		更新模块
//	@Description	更新 Rsync 模块
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Param			name	path		string			true	"模块名称"
//	@Param			data	body		requests.Update	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/modules/{name} [post]
func (r *RsyncController) Update(ctx http.Context) http.Response {
	var updateRequest requests.Update
	sanitize := controllers.Sanitize(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	config, err := tools.Read("/etc/rsyncd.conf")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if !strings.Contains(config, "["+updateRequest.Name+"]") {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "模块 "+updateRequest.Name+" 不存在")
	}

	newConf := `# ` + updateRequest.Name + `-START
[` + updateRequest.Name + `]
path = ` + updateRequest.Path + `
comment = ` + updateRequest.Comment + `
read only = no
auth users = ` + updateRequest.AuthUser + `
hosts allow = ` + updateRequest.HostsAllow + `
secrets file = /etc/rsyncd.secrets
# ` + updateRequest.Name + `-END`

	module := tools.Cut(config, "# "+updateRequest.Name+"-START", "# "+updateRequest.Name+"-END")
	config = strings.Replace(config, "# "+updateRequest.Name+"-START"+module+"# "+updateRequest.Name+"-END", newConf, -1)

	match := regexp.MustCompile(`auth users = ([^\n]+)`).FindStringSubmatch(module)
	if len(match) == 2 {
		authUser := match[1]
		if out, err := tools.Exec("sed -i '/^" + authUser + ":.*$/d' /etc/rsyncd.secrets"); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	if err = tools.Write("/etc/rsyncd.conf", config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if out, err := tools.Exec("echo '" + updateRequest.AuthUser + ":" + updateRequest.Secret + "' >> /etc/rsyncd.secrets"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	if err = tools.ServiceRestart("rsyncd"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}

// GetConfig
//
//	@Summary		获取配置
//	@Description	获取 Rsync 配置
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/config [get]
func (r *RsyncController) GetConfig(ctx http.Context) http.Response {
	config, err := tools.Read("/etc/rsyncd.conf")
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, config)
}

// UpdateConfig
//
//	@Summary		更新配置
//	@Description	更新 Rsync 配置
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.UpdateConfig	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/config [post]
func (r *RsyncController) UpdateConfig(ctx http.Context) http.Response {
	var updateRequest requests.UpdateConfig
	sanitize := controllers.Sanitize(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.Write("/etc/rsyncd.conf", updateRequest.Config, 0644); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := tools.ServiceRestart("rsyncd"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}
