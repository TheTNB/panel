package controllers

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/spf13/cast"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/http/controllers"
	"panel/packages/helpers"
)

type OpenRestyController struct {
	//Dependent services
}

func NewOpenrestyController() *OpenRestyController {
	return &OpenRestyController{
		//Inject services
	}
}

// Status 获取运行状态
func (r *OpenRestyController) Status(ctx http.Context) {
	cmd := exec.Command("bash", "-c", "systemctl status openresty | grep Active | grep -v grep | awk '{print $2}'")
	out, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 获取OpenResty状态失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	status := strings.TrimSpace(string(out))
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, true)
	} else {
		controllers.Success(ctx, false)
	}
}

// Reload 重载配置
func (r *OpenRestyController) Reload(ctx http.Context) {
	cmd := exec.Command("bash", "-c", "systemctl reload openresty")
	_, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 重载OpenResty失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "重载OpenResty失败")
		return
	}

	cmd = exec.Command("bash", "-c", "systemctl status openresty | grep Active | grep -v grep | awk '{print $2}'")
	out, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 获取OpenResty状态失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty状态失败")
		return
	}

	status := strings.TrimSpace(string(out))
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, "重载OpenResty成功")
	} else {
		controllers.Error(ctx, 1, "重载OpenResty失败: "+string(out))
	}
}

// Start 启动OpenResty
func (r *OpenRestyController) Start(ctx http.Context) {
	cmd := exec.Command("bash", "-c", "systemctl start openresty")
	_, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 启动OpenResty失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "启动OpenResty失败")
		return
	}

	cmd = exec.Command("bash", "-c", "systemctl status openresty | grep Active | grep -v grep | awk '{print $2}'")
	out, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 获取OpenResty状态失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty状态失败")
		return
	}

	status := strings.TrimSpace(string(out))
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, "启动OpenResty成功")
	} else {
		controllers.Error(ctx, 1, "启动OpenResty失败: "+string(out))
	}
}

// Stop 停止OpenResty
func (r *OpenRestyController) Stop(ctx http.Context) {
	cmd := exec.Command("bash", "-c", "systemctl stop openresty")
	_, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 停止OpenResty失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "停止OpenResty失败")
		return
	}

	cmd = exec.Command("bash", "-c", "systemctl status openresty | grep Active | grep -v grep | awk '{print $2}'")
	out, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 获取OpenResty状态失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty状态失败")
		return
	}

	status := strings.TrimSpace(string(out))
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, "停止OpenResty成功")
	} else {
		controllers.Error(ctx, 1, "停止OpenResty失败: "+string(out))
	}
}

// Restart 重启OpenResty
func (r *OpenRestyController) Restart(ctx http.Context) {
	cmd := exec.Command("bash", "-c", "systemctl restart openresty")
	out, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 重启OpenResty失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "重启OpenResty失败")
		return
	}

	cmd = exec.Command("bash", "-c", "systemctl status openresty | grep Active | grep -v grep | awk '{print $2}'")
	out, err = cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 获取OpenResty状态失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty状态失败")
		return
	}

	status := strings.TrimSpace(string(out))
	if len(status) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty状态失败")
		return
	}

	if status == "active" {
		controllers.Success(ctx, "重启OpenResty成功")
	} else {
		controllers.Error(ctx, 1, "重启OpenResty失败: "+string(out))
	}
}

// GetConfig 获取配置
func (r *OpenRestyController) GetConfig(ctx http.Context) {
	config, err := os.ReadFile("/www/server/openresty/conf/nginx.conf")
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty配置失败")
		return
	}

	controllers.Success(ctx, string(config))
}

// SaveConfig 保存配置
func (r *OpenRestyController) SaveConfig(ctx http.Context) {
	config := ctx.Request().Input("config")
	if len(config) == 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "配置不能为空")
		return
	}

	oldConfig, err := os.ReadFile("/www/server/openresty/conf/nginx.conf")
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty配置失败")
		return
	}

	err = os.WriteFile("/www/server/openresty/conf/nginx.conf", []byte(config), 0644)
	if err != nil {
		facades.Log().Error("[OpenResty] 保存OpenResty配置失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "保存OpenResty配置失败")
		return
	}

	cmd := exec.Command("bash", "-c", "openresty -t 2>&1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 检查OpenResty配置失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "检查OpenResty配置失败")
		return
	}

	if !strings.Contains(string(out), "test is successful") {
		err = os.WriteFile("/www/server/openresty/conf/nginx.conf", oldConfig, 0644)
		if err != nil {
			facades.Log().Error("[OpenResty] 保存OpenResty配置失败: " + err.Error())
			controllers.Error(ctx, http.StatusInternalServerError, "保存OpenResty配置失败")
			return
		}

		controllers.Error(ctx, http.StatusForbidden, "OpenResty配置有误，请修正后再保存: "+string(out))
	}

	cmd = exec.Command("bash", "-c", "systemctl reload openresty")
	_, err = cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 重载OpenResty配置失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "重载OpenResty配置失败")
		return
	}

	controllers.Success(ctx, "保存OpenResty配置成功")
}

// ErrorLog 获取错误日志
func (r *OpenRestyController) ErrorLog(ctx http.Context) {
	cmd := exec.Command("bash", "-c", "tail -n 100 /www/wwwlogs/nginx_error.log")
	out, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 获取OpenResty错误日志失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty错误日志失败")
		return
	}

	controllers.Success(ctx, string(out))
}

// ClearErrorLog 清空错误日志
func (r *OpenRestyController) ClearErrorLog(ctx http.Context) {
	cmd := exec.Command("bash", "-c", "echo '' > /www/wwwlogs/nginx_error.log")
	_, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 清空OpenResty错误日志失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "清空OpenResty错误日志失败")
		return
	}

	controllers.Success(ctx, "清空OpenResty错误日志成功")
}

// Load 获取负载
func (r *OpenRestyController) Load(ctx http.Context) {
	client := req.C().SetTimeout(10 * time.Second)
	resp, err := client.R().Get("http://127.0.0.1/nginx_status")
	if err != nil || !resp.IsSuccessState() {
		facades.Log().Error("[OpenResty] 获取OpenResty负载失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty负载失败")
		return
	}

	raw := resp.String()
	var data map[int]map[string]any

	cmd := exec.Command("bash", "-c", "ps aux | grep nginx | grep 'worker process' | wc -l")
	out, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 获取OpenResty工作进程失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty工作进程失败")
		return
	}
	workers := strings.TrimSpace(string(out))
	data[0]["name"] = "工作进程"
	data[0]["value"] = workers

	cmd = exec.Command("bash", "-c", "ps aux | grep nginx | grep 'worker process' | awk '{memsum+=$6};END {print memsum}'")
	out, err = cmd.CombinedOutput()
	if err != nil {
		facades.Log().Error("[OpenResty] 获取OpenResty内存占用失败: " + err.Error())
		controllers.Error(ctx, http.StatusInternalServerError, "获取OpenResty内存占用失败")
		return
	}
	mem := helpers.FormatBytes(cast.ToFloat64(strings.TrimSpace(string(out))))
	data[1]["name"] = "内存占用"
	data[1]["value"] = mem

	match := regexp.MustCompile(`Active connections:\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) == 2 {
		data[2]["name"] = "活跃连接数"
		data[2]["value"] = match[1]
	}

	match = regexp.MustCompile(`server accepts handled requests\s+(\d+)\s+(\d+)\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) == 4 {
		data[3]["name"] = "总连接次数"
		data[3]["value"] = match[1]
		data[4]["name"] = "总握手次数"
		data[4]["value"] = match[2]
		data[5]["name"] = "总请求次数"
		data[5]["value"] = match[3]
	}

	match = regexp.MustCompile(`Reading:\s+(\d+)\s+Writing:\s+(\d+)\s+Waiting:\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) == 4 {
		data[6]["name"] = "请求数"
		data[6]["value"] = match[1]
		data[7]["name"] = "响应数"
		data[7]["value"] = match[2]
		data[8]["name"] = "驻留进程"
		data[8]["value"] = match[3]
	}

	controllers.Success(ctx, data)
}
