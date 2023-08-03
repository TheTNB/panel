package pureftpd

import (
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"

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
func (c *PureFtpdController) Status(ctx http.Context) {
	if !controllers.Check(ctx, "pureftpd") {
		return
	}

	status := tools.ExecShell("systemctl status pure-ftpd | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Reload 重载配置
func (c *PureFtpdController) Reload(ctx http.Context) {
	if !controllers.Check(ctx, "pureftpd") {
		return
	}

	tools.ExecShell("systemctl reload pure-ftpd")
	status := tools.ExecShell("systemctl status pure-ftpd | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Restart 重启服务
func (c *PureFtpdController) Restart(ctx http.Context) {
	if !controllers.Check(ctx, "pureftpd") {
		return
	}

	tools.ExecShell("systemctl restart pure-ftpd")
	status := tools.ExecShell("systemctl status pure-ftpd | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Start 启动服务
func (c *PureFtpdController) Start(ctx http.Context) {
	if !controllers.Check(ctx, "pureftpd") {
		return
	}

	tools.ExecShell("systemctl start pure-ftpd")
	status := tools.ExecShell("systemctl status pure-ftpd | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Stop 停止服务
func (c *PureFtpdController) Stop(ctx http.Context) {
	if !controllers.Check(ctx, "pureftpd") {
		return
	}

	tools.ExecShell("systemctl stop pure-ftpd")
	status := tools.ExecShell("systemctl status pure-ftpd | grep Active | grep -v grep | awk '{print $2}'")
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd状态失败")
		return
	}

	if status != "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// List 获取用户列表
func (c *PureFtpdController) List(ctx http.Context) {
	if !controllers.Check(ctx, "pureftpd") {
		return
	}

	listRaw := tools.ExecShell("pure-pw list")
	if len(listRaw) == 0 {
		controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []User{},
		})
		return
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
		controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []User{},
		})
		return
	}
	if endIndex > len(users) {
		endIndex = len(users)
	}
	pagedUsers := users[startIndex:endIndex]

	controllers.Success(ctx, http.Json{
		"total": len(users),
		"items": pagedUsers,
	})
}

// Add 添加用户
func (c *PureFtpdController) Add(ctx http.Context) {
	if !controllers.Check(ctx, "pureftpd") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"username": "required",
		"password": "required|min_len:6",
		"path":     "required",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")
	path := ctx.Request().Input("path")

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !tools.Exists(path) {
		controllers.Error(ctx, http.StatusBadRequest, "目录不存在")
		return
	}

	tools.Chmod(path, 755)
	tools.Chown(path, "www", "www")
	tools.ExecShell(`yes '` + password + `' | pure-pw useradd ` + username + ` -u www -g www -d ` + path)
	tools.ExecShell("pure-pw mkdb")

	controllers.Success(ctx, nil)
}

// Delete 删除用户
func (c *PureFtpdController) Delete(ctx http.Context) {
	if !controllers.Check(ctx, "pureftpd") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"username": "required",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	username := ctx.Request().Input("username")

	tools.ExecShell("pure-pw userdel " + username + " -m")
	tools.ExecShell("pure-pw mkdb")

	controllers.Success(ctx, nil)
}

// ChangePassword 修改密码
func (c *PureFtpdController) ChangePassword(ctx http.Context) {
	if !controllers.Check(ctx, "pureftpd") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"username": "required",
		"password": "required|min_len:6",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")

	tools.ExecShell(`yes '` + password + `' | pure-pw passwd ` + username + ` -m`)
	tools.ExecShell("pure-pw mkdb")

	controllers.Success(ctx, nil)
}

// GetPort 获取端口
func (c *PureFtpdController) GetPort(ctx http.Context) {
	if !controllers.Check(ctx, "pureftpd") {
		return
	}

	port := tools.ExecShell(`cat /www/server/pure-ftpd/etc/pure-ftpd.conf | grep "Bind" | awk '{print $2}' | awk -F "," '{print $2}'`)
	if len(port) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取PureFtpd端口失败")
		return
	}

	controllers.Success(ctx, port)
}

// SetPort 设置端口
func (c *PureFtpdController) SetPort(ctx http.Context) {
	if !controllers.Check(ctx, "pureftpd") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"port": "required",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	port := ctx.Request().Input("port")
	tools.ExecShell(`sed -i "s/Bind.*/Bind 0.0.0.0,` + port + `/g" /www/server/pure-ftpd/etc/pure-ftpd.conf`)
	if tools.IsRHEL() {
		tools.ExecShell("firewall-cmd --zone=public --add-port=" + port + "/tcp --permanent")
		tools.ExecShell("firewall-cmd --reload")
	} else {
		tools.ExecShell("ufw allow " + port + "/tcp")
		tools.ExecShell("ufw reload")
	}
	tools.ExecShell("systemctl restart pure-ftpd")

	controllers.Success(ctx, nil)
}
