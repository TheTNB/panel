package plugins

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/app/http/controllers"
	"github.com/TheTNB/panel/pkg/tools"
	"github.com/TheTNB/panel/types"
)

type PureFtpdController struct {
}

func NewPureFtpdController() *PureFtpdController {
	return &PureFtpdController{}
}

// List 获取用户列表
func (r *PureFtpdController) List(ctx http.Context) http.Response {
	listRaw, err := tools.Exec("pure-pw list")
	if err != nil {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []types.PureFtpdUser{},
		})
	}

	listArr := strings.Split(listRaw, "\n")
	var users []types.PureFtpdUser
	for _, v := range listArr {
		if len(v) == 0 {
			continue
		}

		match := regexp.MustCompile(`(\S+)\s+(\S+)`).FindStringSubmatch(v)
		users = append(users, types.PureFtpdUser{
			Username: match[1],
			Path:     strings.Replace(match[2], "/./", "/", 1),
		})
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)
	startIndex := (page - 1) * limit
	endIndex := page * limit
	if startIndex > len(users) {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []types.PureFtpdUser{},
		})
	}
	if endIndex > len(users) {
		endIndex = len(users)
	}
	pagedUsers := users[startIndex:endIndex]

	if pagedUsers == nil {
		pagedUsers = []types.PureFtpdUser{}
	}

	return controllers.Success(ctx, http.Json{
		"total": len(users),
		"items": pagedUsers,
	})
}

// Add 添加用户
func (r *PureFtpdController) Add(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"username": "required",
		"password": "required|min_len:6",
		"path":     "required",
	}); sanitize != nil {
		return sanitize
	}

	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")
	path := ctx.Request().Input("path")

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !tools.Exists(path) {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "目录不存在")
	}

	if err := tools.Chmod(path, 0755); err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "修改目录权限失败")
	}
	if err := tools.Chown(path, "www", "www"); err != nil {
		return nil
	}
	if out, err := tools.Exec(`yes '` + password + `' | pure-pw useradd ` + username + ` -u www -g www -d ` + path); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := tools.Exec("pure-pw mkdb"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// Delete 删除用户
func (r *PureFtpdController) Delete(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"username": "required",
	}); sanitize != nil {
		return sanitize
	}

	username := ctx.Request().Input("username")

	if out, err := tools.Exec("pure-pw userdel " + username + " -m"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := tools.Exec("pure-pw mkdb"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// ChangePassword 修改密码
func (r *PureFtpdController) ChangePassword(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"username": "required",
		"password": "required|min_len:6",
	}); sanitize != nil {
		return sanitize
	}

	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")

	if out, err := tools.Exec(`yes '` + password + `' | pure-pw passwd ` + username + ` -m`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := tools.Exec("pure-pw mkdb"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// GetPort 获取端口
func (r *PureFtpdController) GetPort(ctx http.Context) http.Response {
	port, err := tools.Exec(`cat /www/server/pure-ftpd/etc/pure-ftpd.conf | grep "Bind" | awk '{print $2}' | awk -F "," '{print $2}'`)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd端口失败")
	}

	return controllers.Success(ctx, cast.ToInt(port))
}

// SetPort 设置端口
func (r *PureFtpdController) SetPort(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"port": "required",
	}); sanitize != nil {
		return sanitize
	}

	port := ctx.Request().Input("port")
	if out, err := tools.Exec(`sed -i "s/Bind.*/Bind 0.0.0.0,` + port + `/g" /www/server/pure-ftpd/etc/pure-ftpd.conf`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}
	if tools.IsRHEL() {
		if out, err := tools.Exec("firewall-cmd --zone=public --add-port=" + port + "/tcp --permanent"); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := tools.Exec("firewall-cmd --reload"); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, out)
		}
	} else {
		if out, err := tools.Exec("ufw allow " + port + "/tcp"); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := tools.Exec("ufw reload"); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	if err := tools.ServiceRestart("pure-ftpd"); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, nil)
}
