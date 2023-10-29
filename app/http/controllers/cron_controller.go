package controllers

import (
	"regexp"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/spf13/cast"
	"panel/app/models"
	"panel/app/services"
	"panel/pkg/tools"
)

type CronController struct {
	cron    services.Cron
	setting services.Setting
}

func NewCronController() *CronController {
	return &CronController{
		cron:    services.NewCronImpl(),
		setting: services.NewSettingImpl(),
	}
}

// List 获取计划任务列表
func (c *CronController) List(ctx http.Context) http.Response {
	limit := ctx.Request().QueryInt("limit", 10)
	page := ctx.Request().QueryInt("page", 1)

	var crons []models.Cron
	var total int64
	err := facades.Orm().Query().Paginate(page, limit, &crons, &total)
	if err != nil {
		facades.Log().Error("[面板][CronController] 查询计划任务列表失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	return Success(ctx, http.Json{
		"total": total,
		"items": crons,
	})
}

// Add 添加计划任务
func (c *CronController) Add(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"name":        "required|min_len:1|max_len:255",
		"time":        "required",
		"script":      "required",
		"type":        "required|in:shell,backup,cutoff",
		"backup_type": "required_if:type,backup|in:website,mysql,postgresql",
	})
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	// 单独验证时间格式
	if !regexp.MustCompile(`^((\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+)(,(\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+))*\s?){5}$`).MatchString(ctx.Request().Input("time")) {
		return Error(ctx, http.StatusUnprocessableEntity, "时间格式错误")
	}

	shell := ctx.Request().Input("script")
	cronType := ctx.Request().Input("type")
	if cronType == "backup" {
		backupType := ctx.Request().Input("backup_type")
		backupName := ctx.Request().Input("backup_database")
		if backupType == "website" {
			backupName = ctx.Request().Input("website")
		}
		backupPath := ctx.Request().Input("backup_path")
		if len(backupPath) == 0 {
			backupPath = c.setting.Get(models.SettingKeyBackupPath) + "/" + backupType
		}
		backupSave := ctx.Request().InputInt("save", 10)
		shell = `#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

# 耗子Linux面板 - 数据备份脚本

type=` + backupType + `
path=` + backupPath + `
name=` + backupName + `
save=` + cast.ToString(backupSave) + `

# 执行备份
panel backup ${type} ${name} ${path} ${save} 2>&1
`
	}
	if cronType == "cutoff" {
		website := ctx.Request().Input("website")
		save := ctx.Request().InputInt("save", 180)
		shell = `#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

# 耗子Linux面板 - 日志切割脚本

name=` + website + `
save=` + cast.ToString(save) + `

# 执行切割
panel cutoff ${name} ${save} 2>&1
`
	}

	shellDir := "/www/server/cron/"
	shellLogDir := "/www/server/cron/logs/"
	if !tools.Exists(shellDir) {
		facades.Log().Error("[面板][CronController] 计划任务目录不存在")
		return Error(ctx, http.StatusInternalServerError, "计划任务目录不存在")
	}
	if !tools.Exists(shellLogDir) {
		facades.Log().Error("[面板][CronController] 计划任务日志目录不存在")
		return Error(ctx, http.StatusInternalServerError, "计划任务日志目录不存在")
	}
	shellFile := strconv.Itoa(int(carbon.Now().Timestamp())) + tools.RandomString(16)
	if !tools.Write(shellDir+shellFile+".sh", shell, 0700) {
		facades.Log().Error("[面板][CronController] 创建计划任务脚本失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}
	tools.Exec("dos2unix " + shellDir + shellFile + ".sh")

	var cron models.Cron
	cron.Name = ctx.Request().Input("name")
	cron.Type = ctx.Request().Input("type")
	cron.Status = true
	cron.Time = ctx.Request().Input("time")
	cron.Shell = shellDir + shellFile + ".sh"
	cron.Log = shellLogDir + shellFile + ".log"

	err = facades.Orm().Query().Create(&cron)
	if err != nil {
		facades.Log().Error("[面板][CronController] 创建计划任务失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	c.cron.AddToSystem(cron)

	return Success(ctx, http.Json{
		"id": cron.ID,
	})
}

// Script 获取脚本内容
func (c *CronController) Script(ctx http.Context) http.Response {
	var cron models.Cron
	err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, "计划任务不存在")
	}

	return Success(ctx, tools.Read(cron.Shell))
}

// Update 更新计划任务
func (c *CronController) Update(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"name":   "required|min_len:1|max_len:255",
		"time":   "required",
		"script": "required",
	})
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	// 单独验证时间格式
	if !regexp.MustCompile(`^((\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+)(,(\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+))*\s?){5}$`).MatchString(ctx.Request().Input("time")) {
		return Error(ctx, http.StatusUnprocessableEntity, "时间格式错误")
	}

	var cron models.Cron
	err = facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, "计划任务不存在")
	}

	if !cron.Status {
		return Error(ctx, http.StatusUnprocessableEntity, "计划任务已禁用")
	}

	cron.Time = ctx.Request().Input("time")
	cron.Name = ctx.Request().Input("name")
	err = facades.Orm().Query().Save(&cron)
	if err != nil {
		facades.Log().Error("[面板][CronController] 更新计划任务失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if !tools.Write(cron.Shell, ctx.Request().Input("script"), 0644) {
		facades.Log().Error("[面板][CronController] 更新计划任务脚本失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}
	tools.Exec("dos2unix " + cron.Shell)

	c.cron.DeleteFromSystem(cron)
	if cron.Status {
		c.cron.AddToSystem(cron)
	}

	return Success(ctx, nil)
}

// Delete 删除计划任务
func (c *CronController) Delete(ctx http.Context) http.Response {
	var cron models.Cron
	err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, "计划任务不存在")
	}

	c.cron.DeleteFromSystem(cron)
	tools.Remove(cron.Shell)

	_, err = facades.Orm().Query().Delete(&cron)
	if err != nil {
		facades.Log().Error("[面板][CronController] 删除计划任务失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	return Success(ctx, nil)
}

// Status 更新计划任务状态
func (c *CronController) Status(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"status": "bool",
	})
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	var cron models.Cron
	err = facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, "计划任务不存在")
	}

	cron.Status = ctx.Request().InputBool("status")
	err = facades.Orm().Query().Save(&cron)
	if err != nil {
		facades.Log().Error("[面板][CronController] 更新计划任务状态失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	c.cron.DeleteFromSystem(cron)
	if cron.Status {
		c.cron.AddToSystem(cron)
	}

	return Success(ctx, nil)
}

// Log 获取计划任务日志
func (c *CronController) Log(ctx http.Context) http.Response {
	var cron models.Cron
	err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, "计划任务不存在")
	}

	if !tools.Exists(cron.Log) {
		return Error(ctx, http.StatusUnprocessableEntity, "日志文件不存在")
	}

	log := tools.Exec("tail -n 1000 " + cron.Log)

	return Success(ctx, log)
}
