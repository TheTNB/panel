package plugins

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/spf13/cast"

	"panel/app/http/controllers"
	"panel/pkg/tools"
)

type PureFtpdController struct {
}

type User struct {
	Username string `json:"username"`
	Path     string `json:"path"`
}

func NewPureFtpdController() *PureFtpdController {
	return &PureFtpdController{}
}

// Status 获取运行状态
func (r *PureFtpdController) Status(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "pureftpd")
	if check != nil {
		return check
	}

	status := tools.Exec("systemctl status pure-ftpd | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Restart 重启服务
func (r *PureFtpdController) Restart(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "pureftpd")
	if check != nil {
		return check
	}

	tools.Exec("systemctl restart pure-ftpd")
	status := tools.Exec("systemctl status pure-ftpd | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Start 启动服务
func (r *PureFtpdController) Start(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "pureftpd")
	if check != nil {
		return check
	}

	tools.Exec("systemctl start pure-ftpd")
	status := tools.Exec("systemctl status pure-ftpd | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd状态失败")
	}

	if status == "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// Stop 停止服务
func (r *PureFtpdController) Stop(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "pureftpd")
	if check != nil {
		return check
	}

	tools.Exec("systemctl stop pure-ftpd")
	status := tools.Exec("systemctl status pure-ftpd | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd状态失败")
	}

	if status != "active" {
		return controllers.Success(ctx, true)
	} else {
		return controllers.Success(ctx, false)
	}
}

// List 获取用户列表
func (r *PureFtpdController) List(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "pureftpd")
	if check != nil {
		return check
	}

	listRaw := tools.Exec("pure-pw list")
	if len(listRaw) == 0 {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []User{},
		})
	}

	listArr := strings.Split(listRaw, "\n")
	var users []User
	for _, v := range listArr {
		if len(v) == 0 {
			continue
		}

		match := regexp.MustCompile(`(\S+)\s+(\S+)`).FindStringSubmatch(v)
		users = append(users, User{
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
			"items": []User{},
		})
	}
	if endIndex > len(users) {
		endIndex = len(users)
	}
	pagedUsers := users[startIndex:endIndex]

	return controllers.Success(ctx, http.Json{
		"total": len(users),
		"items": pagedUsers,
	})
}

// Add 添加用户
func (r *PureFtpdController) Add(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "pureftpd")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"username": "required",
		"password": "required|min_len:6",
		"path":     "required",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
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

	tools.Chmod(path, 0755)
	tools.Chown(path, "www", "www")
	tools.Exec(`yes '` + password + `' | pure-pw useradd ` + username + ` -u www -g www -d ` + path)
	tools.Exec("pure-pw mkdb")

	return controllers.Success(ctx, nil)
}

// Delete 删除用户
func (r *PureFtpdController) Delete(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "pureftpd")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"username": "required",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	username := ctx.Request().Input("username")

	tools.Exec("pure-pw userdel " + username + " -m")
	tools.Exec("pure-pw mkdb")

	return controllers.Success(ctx, nil)
}

// ChangePassword 修改密码
func (r *PureFtpdController) ChangePassword(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "pureftpd")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"username": "required",
		"password": "required|min_len:6",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")

	tools.Exec(`yes '` + password + `' | pure-pw passwd ` + username + ` -m`)
	tools.Exec("pure-pw mkdb")

	return controllers.Success(ctx, nil)
}

// GetPort 获取端口
func (r *PureFtpdController) GetPort(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "pureftpd")
	if check != nil {
		return check
	}

	port := tools.Exec(`cat /www/server/pure-ftpd/etc/pure-ftpd.conf | grep "Bind" | awk '{print $2}' | awk -F "," '{print $2}'`)
	if len(port) == 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd端口失败")
	}

	return controllers.Success(ctx, cast.ToInt(port))
}

// SetPort 设置端口
func (r *PureFtpdController) SetPort(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "pureftpd")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"port": "required",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	port := ctx.Request().Input("port")
	tools.Exec(`sed -i "s/Bind.*/Bind 0.0.0.0,` + port + `/g" /www/server/pure-ftpd/etc/pure-ftpd.conf`)
	if tools.IsRHEL() {
		tools.Exec("firewall-cmd --zone=public --add-port=" + port + "/tcp --permanent")
		tools.Exec("firewall-cmd --reload")
	} else {
		tools.Exec("ufw allow " + port + "/tcp")
		tools.Exec("ufw reload")
	}
	tools.Exec("systemctl restart pure-ftpd")

	return controllers.Success(ctx, nil)
}
