package controllers

import (
	"regexp"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	"panel/app/models"
	"panel/app/services"
	"panel/pkg/tools"
)

type CronController struct {
	cron services.Cron
}

func NewCronController() *CronController {
	return &CronController{
		cron: services.NewCronImpl(),
	}
}

func (r *CronController) List(ctx http.Context) {
	limit := ctx.Request().QueryInt("limit")
	page := ctx.Request().QueryInt("page")

	var crons []models.Cron
	var total int64
	err := facades.Orm().Query().Paginate(page, limit, &crons, &total)
	if err != nil {
		facades.Log().Error("[面板][CronController] 查询计划任务列表失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, http.Json{
		"total": total,
		"items": crons,
	})
}

func (r *CronController) Add(ctx http.Context) {
	validator, err := ctx.Request().Validate(map[string]string{
		"name":   "required|min_len:1|max_len:255",
		"time":   "required",
		"script": "required",
	})
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	// 单独验证时间格式
	if !regexp.MustCompile(`^((\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+)(,(\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+))*\s?){5}$`).MatchString(ctx.Request().Input("time")) {
		Error(ctx, http.StatusBadRequest, "时间格式错误")
		return
	}

	// 写入shell
	shellDir := "/www/server/cron/"
	shellLogDir := "/www/server/cron/logs/"
	if !tools.Exists(shellDir) {
		facades.Log().Error("[面板][CronController] 计划任务目录不存在")
		Error(ctx, http.StatusInternalServerError, "计划任务目录不存在")
		return
	}
	if !tools.Exists(shellLogDir) {
		facades.Log().Error("[面板][CronController] 计划任务日志目录不存在")
		Error(ctx, http.StatusInternalServerError, "计划任务日志目录不存在")
		return
	}
	shellFile := strconv.Itoa(int(carbon.Now().Timestamp())) + tools.RandomString(16)
	if !tools.WriteFile(shellDir+shellFile+".sh", ctx.Request().Input("script"), 0700) {
		facades.Log().Error("[面板][CronController] 创建计划任务脚本失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	tools.ExecShell("dos2unix " + shellDir + shellFile + ".sh")

	var cron models.Cron
	cron.Name = ctx.Request().Input("name")
	cron.Type = "shell"
	cron.Status = true
	cron.Time = ctx.Request().Input("time")
	cron.Shell = shellDir + shellFile + ".sh"
	cron.Log = shellLogDir + shellFile + ".log"

	err = facades.Orm().Query().Create(&cron)
	if err != nil {
		facades.Log().Error("[面板][CronController] 创建计划任务失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	r.cron.AddToSystem(cron)

	Success(ctx, http.Json{
		"id": cron.ID,
	})
}

// Script 获取脚本内容
func (r *CronController) Script(ctx http.Context) {
	var cron models.Cron
	err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		Error(ctx, http.StatusBadRequest, "计划任务不存在")
		return
	}

	Success(ctx, tools.ReadFile(cron.Shell))
}

func (r *CronController) Update(ctx http.Context) {
	validator, err := ctx.Request().Validate(map[string]string{
		"name":   "required|min_len:1|max_len:255",
		"time":   "required",
		"script": "required",
	})
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		Error(ctx, http.StatusBadRequest, validator.Errors().All())
		return
	}

	// 单独验证时间格式
	if !regexp.MustCompile(`^((\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+)(,(\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+))*\s?){5}$`).MatchString(ctx.Request().Input("time")) {
		Error(ctx, http.StatusBadRequest, "时间格式错误")
		return
	}

	var cron models.Cron
	err = facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		Error(ctx, http.StatusBadRequest, "计划任务不存在")
		return
	}

	if !cron.Status {
		Error(ctx, http.StatusBadRequest, "计划任务已禁用")
		return
	}

	cron.Time = ctx.Request().Input("time")
	cron.Name = ctx.Request().Input("name")
	err = facades.Orm().Query().Save(&cron)
	if err != nil {
		facades.Log().Error("[面板][CronController] 更新计划任务失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	if !tools.WriteFile(cron.Shell, ctx.Request().Input("script"), 0644) {
		facades.Log().Error("[面板][CronController] 更新计划任务脚本失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	tools.ExecShell("dos2unix " + cron.Shell)

	r.cron.DeleteFromSystem(cron)
	if cron.Status {
		r.cron.AddToSystem(cron)
	}

	Success(ctx, nil)
}

func (r *CronController) Delete(ctx http.Context) {
	var cron models.Cron
	err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		Error(ctx, http.StatusBadRequest, "计划任务不存在")
		return
	}

	r.cron.DeleteFromSystem(cron)
	tools.RemoveFile(cron.Shell)

	_, err = facades.Orm().Query().Delete(&cron)
	if err != nil {
		facades.Log().Error("[面板][CronController] 删除计划任务失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}

func (r *CronController) Status(ctx http.Context) {
	validator, err := ctx.Request().Validate(map[string]string{
		"status": "bool",
	})
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	var cron models.Cron
	err = facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		Error(ctx, http.StatusBadRequest, "计划任务不存在")
		return
	}

	cron.Status = ctx.Request().InputBool("status")
	err = facades.Orm().Query().Save(&cron)
	if err != nil {
		facades.Log().Error("[面板][CronController] 更新计划任务状态失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	r.cron.DeleteFromSystem(cron)
	if cron.Status {
		r.cron.AddToSystem(cron)
	}

	Success(ctx, nil)
}

func (r *CronController) Log(ctx http.Context) {
	var cron models.Cron
	err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		Error(ctx, http.StatusBadRequest, "计划任务不存在")
		return
	}

	if !tools.Exists(cron.Log) {
		Error(ctx, http.StatusBadRequest, "日志文件不存在")
		return
	}

	Success(ctx, tools.ReadFile(cron.Log))
}
