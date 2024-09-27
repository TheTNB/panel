package rsync

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/str"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
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
func (s *Service) List(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read("/etc/rsyncd.conf")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var modules []Module
	lines := strings.Split(config, "\n")
	var currentModule *Module

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
			currentModule = &Module{
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
						service.Error(w, http.StatusInternalServerError, "获取模块"+currentModule.AuthUser+"的密钥失败")
						return
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

	paged, total := service.Paginate(r, modules)

	service.Success(w, chix.M{
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
func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Create](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	config, err := io.Read("/etc/rsyncd.conf")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if strings.Contains(config, "["+req.Name+"]") {
		service.Error(w, http.StatusUnprocessableEntity, "模块 "+req.Name+" 已存在")
		return
	}

	conf := `# ` + req.Name + `-START
[` + req.Name + `]
path = ` + req.Path + `
comment = ` + req.Comment + `
read only = no
auth users = ` + req.AuthUser + `
hosts allow = ` + req.HostsAllow + `
secrets file = /etc/rsyncd.secrets
# ` + req.Name + `-END
`

	if err = io.WriteAppend("/etc/rsyncd.conf", conf); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if out, err := shell.Execf("echo '" + req.AuthUser + ":" + req.Secret + "' >> /etc/rsyncd.secrets"); err != nil {
		service.Error(w, http.StatusInternalServerError, out)
		return
	}

	if err = systemctl.Restart("rsyncd"); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
}

// Delete
//
//	@Summary		删除模块
//	@Description	删除 Rsync 模块
//	@Tags			插件-Rsync
//	@Produce		json
//	@Security		BearerToken
//	@Param			name	path		string	true	"模块名称"
//	@Success		200		{object}	controllers.SuccessResponse
//	@Router			/plugins/rsync/modules/{name} [delete]
func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Delete](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	config, err := io.Read("/etc/rsyncd.conf")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !strings.Contains(config, "["+req.Name+"]") {
		service.Error(w, http.StatusUnprocessableEntity, "模块 "+req.Name+" 不存在")
		return
	}

	module := str.Cut(config, "# "+req.Name+"-START", "# "+req.Name+"-END")
	config = strings.Replace(config, "\n# "+req.Name+"-START"+module+"# "+req.Name+"-END", "", -1)

	match := regexp.MustCompile(`auth users = ([^\n]+)`).FindStringSubmatch(module)
	if len(match) == 2 {
		authUser := match[1]
		if out, err := shell.Execf("sed -i '/^" + authUser + ":.*$/d' /etc/rsyncd.secrets"); err != nil {
			service.Error(w, http.StatusInternalServerError, out)
			return
		}
	}

	if err = io.Write("/etc/rsyncd.conf", config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = systemctl.Restart("rsyncd"); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
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
func (s *Service) Update(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Update](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	config, err := io.Read("/etc/rsyncd.conf")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !strings.Contains(config, "["+req.Name+"]") {
		service.Error(w, http.StatusUnprocessableEntity, "模块 "+req.Name+" 不存在")
		return
	}

	newConf := `# ` + req.Name + `-START
[` + req.Name + `]
path = ` + req.Path + `
comment = ` + req.Comment + `
read only = no
auth users = ` + req.AuthUser + `
hosts allow = ` + req.HostsAllow + `
secrets file = /etc/rsyncd.secrets
# ` + req.Name + `-END`

	module := str.Cut(config, "# "+req.Name+"-START", "# "+req.Name+"-END")
	config = strings.Replace(config, "# "+req.Name+"-START"+module+"# "+req.Name+"-END", newConf, -1)

	match := regexp.MustCompile(`auth users = ([^\n]+)`).FindStringSubmatch(module)
	if len(match) == 2 {
		authUser := match[1]
		if out, err := shell.Execf("sed -i '/^" + authUser + ":.*$/d' /etc/rsyncd.secrets"); err != nil {
			service.Error(w, http.StatusInternalServerError, out)
			return
		}
	}

	if err = io.Write("/etc/rsyncd.conf", config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if out, err := shell.Execf("echo '" + req.AuthUser + ":" + req.Secret + "' >> /etc/rsyncd.secrets"); err != nil {
		service.Error(w, http.StatusInternalServerError, out)
		return
	}

	if err = systemctl.Restart("rsyncd"); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
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
func (s *Service) GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read("/etc/rsyncd.conf")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, config)
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
func (s *Service) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = io.Write("/etc/rsyncd.conf", req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = systemctl.Restart("rsyncd"); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
}
