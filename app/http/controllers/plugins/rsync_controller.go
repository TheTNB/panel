package plugins

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"

	requests "github.com/TheTNB/panel/v2/app/http/requests/plugins/rsync"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/str"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type RsyncController struct {
}

func NewRsyncController() *RsyncController {
	return &RsyncController{}
}

// List
//
//	@Summary		列出模块
//	@Description	列出所有 Rsync 模块
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		commonrequests.Paginate	true	"request"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/modules [get]
func (r *RsyncController) List(ctx http.Context) http.Response {
	config, err := io.Read("/etc/rsyncd.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	var modules []types.RsyncModule
	lines := strings.Split(config, "\n")
	var currentModule *types.RsyncModule

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
			currentModule = &types.RsyncModule{
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
					currentModule.Secret, err = shell.Execf("grep -E '^" + currentModule.AuthUser + ":.*$' /etc/rsyncd.secrets | awk -F ':' '{print $2}'")
					if err != nil {
						return h.Error(ctx, http.StatusInternalServerError, "获取模块"+currentModule.AuthUser+"的密钥失败")
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

	paged, total := h.Paginate(ctx, modules)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
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
	sanitize := h.SanitizeRequest(ctx, &createRequest)
	if sanitize != nil {
		return sanitize
	}

	config, err := io.Read("/etc/rsyncd.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if strings.Contains(config, "["+createRequest.Name+"]") {
		return h.Error(ctx, http.StatusUnprocessableEntity, "模块 "+createRequest.Name+" 已存在")
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

	if err := io.WriteAppend("/etc/rsyncd.conf", conf); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if out, err := shell.Execf("echo '" + createRequest.AuthUser + ":" + createRequest.Secret + "' >> /etc/rsyncd.secrets"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	if err := systemctl.Restart("rsyncd"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
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
		return h.Error(ctx, http.StatusUnprocessableEntity, "name 不能为空")
	}

	config, err := io.Read("/etc/rsyncd.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if !strings.Contains(config, "["+name+"]") {
		return h.Error(ctx, http.StatusUnprocessableEntity, "模块 "+name+" 不存在")
	}

	module := str.Cut(config, "# "+name+"-START", "# "+name+"-END")
	config = strings.Replace(config, "\n# "+name+"-START"+module+"# "+name+"-END", "", -1)

	match := regexp.MustCompile(`auth users = ([^\n]+)`).FindStringSubmatch(module)
	if len(match) == 2 {
		authUser := match[1]
		if out, err := shell.Execf("sed -i '/^" + authUser + ":.*$/d' /etc/rsyncd.secrets"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	if err = io.Write("/etc/rsyncd.conf", config, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err = systemctl.Restart("rsyncd"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
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
	sanitize := h.SanitizeRequest(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	config, err := io.Read("/etc/rsyncd.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if !strings.Contains(config, "["+updateRequest.Name+"]") {
		return h.Error(ctx, http.StatusUnprocessableEntity, "模块 "+updateRequest.Name+" 不存在")
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

	module := str.Cut(config, "# "+updateRequest.Name+"-START", "# "+updateRequest.Name+"-END")
	config = strings.Replace(config, "# "+updateRequest.Name+"-START"+module+"# "+updateRequest.Name+"-END", newConf, -1)

	match := regexp.MustCompile(`auth users = ([^\n]+)`).FindStringSubmatch(module)
	if len(match) == 2 {
		authUser := match[1]
		if out, err := shell.Execf("sed -i '/^" + authUser + ":.*$/d' /etc/rsyncd.secrets"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	if err = io.Write("/etc/rsyncd.conf", config, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if out, err := shell.Execf("echo '" + updateRequest.AuthUser + ":" + updateRequest.Secret + "' >> /etc/rsyncd.secrets"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	if err = systemctl.Restart("rsyncd"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
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
	config, err := io.Read("/etc/rsyncd.conf")
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, config)
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
	sanitize := h.SanitizeRequest(ctx, &updateRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := io.Write("/etc/rsyncd.conf", updateRequest.Config, 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := systemctl.Restart("rsyncd"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}
