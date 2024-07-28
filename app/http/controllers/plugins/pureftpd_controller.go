package plugins

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/os"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type PureFtpdController struct {
}

func NewPureFtpdController() *PureFtpdController {
	return &PureFtpdController{}
}

// List 获取用户列表
func (r *PureFtpdController) List(ctx http.Context) http.Response {
	listRaw, err := shell.Execf("pure-pw list")
	if err != nil {
		return h.Success(ctx, http.Json{
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

	paged, total := h.Paginate(ctx, users)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// Add 添加用户
func (r *PureFtpdController) Add(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
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
	if !io.Exists(path) {
		return h.Error(ctx, http.StatusUnprocessableEntity, "目录不存在")
	}

	if err := io.Chmod(path, 0755); err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "修改目录权限失败")
	}
	if err := io.Chown(path, "www", "www"); err != nil {
		return nil
	}
	if out, err := shell.Execf(`yes '` + password + `' | pure-pw useradd ` + username + ` -u www -g www -d ` + path); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf("pure-pw mkdb"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}

// Delete 删除用户
func (r *PureFtpdController) Delete(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"username": "required",
	}); sanitize != nil {
		return sanitize
	}

	username := ctx.Request().Input("username")

	if out, err := shell.Execf("pure-pw userdel " + username + " -m"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf("pure-pw mkdb"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}

// ChangePassword 修改密码
func (r *PureFtpdController) ChangePassword(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"username": "required",
		"password": "required|min_len:6",
	}); sanitize != nil {
		return sanitize
	}

	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")

	if out, err := shell.Execf(`yes '` + password + `' | pure-pw passwd ` + username + ` -m`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf("pure-pw mkdb"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	return h.Success(ctx, nil)
}

// GetPort 获取端口
func (r *PureFtpdController) GetPort(ctx http.Context) http.Response {
	port, err := shell.Execf(`cat /www/server/pure-ftpd/etc/pure-ftpd.conf | grep "Bind" | awk '{print $2}' | awk -F "," '{print $2}'`)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取PureFtpd端口失败")
	}

	return h.Success(ctx, cast.ToInt(port))
}

// SetPort 设置端口
func (r *PureFtpdController) SetPort(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"port": "required",
	}); sanitize != nil {
		return sanitize
	}

	port := ctx.Request().Input("port")
	if out, err := shell.Execf(`sed -i "s/Bind.*/Bind 0.0.0.0,%s/g" /www/server/pure-ftpd/etc/pure-ftpd.conf`, port); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if os.IsRHEL() {
		if out, err := shell.Execf("firewall-cmd --zone=public --add-port=%s/tcp --permanent", port); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := shell.Execf("firewall-cmd --reload"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	} else {
		if out, err := shell.Execf("ufw allow %s/tcp", port); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
		if out, err := shell.Execf("ufw reload"); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	}

	if err := systemctl.Restart("pure-ftpd"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}
