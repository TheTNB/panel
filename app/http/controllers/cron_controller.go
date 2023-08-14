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

func (c *CronController) List(ctx http.Context) {
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

func (c *CronController) Add(ctx http.Context) {
	validator, err := ctx.Request().Validate(map[string]string{
		"name":        "required|min_len:1|max_len:255",
		"time":        "required",
		"script":      "required",
		"type":        "required|in:shell,backup,cutoff",
		"backup_type": "required_if:type,backup|in:website,mysql,postgresql",
	})
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	// 单独验证时间格式
	if !regexp.MustCompile(`^((\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+)(,(\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+))*\s?){5}$`).MatchString(ctx.Request().Input("time")) {
		Error(ctx, http.StatusBadRequest, "时间格式错误")
		return
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
		if len(backupName) == 0 {
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
		Error(ctx, http.StatusInternalServerError, "计划任务目录不存在")
		return
	}
	if !tools.Exists(shellLogDir) {
		facades.Log().Error("[面板][CronController] 计划任务日志目录不存在")
		Error(ctx, http.StatusInternalServerError, "计划任务日志目录不存在")
		return
	}
	shellFile := strconv.Itoa(int(carbon.Now().Timestamp())) + tools.RandomString(16)
	if !tools.Write(shellDir+shellFile+".sh", shell, 0700) {
		facades.Log().Error("[面板][CronController] 创建计划任务脚本失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
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
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	c.cron.AddToSystem(cron)

	Success(ctx, http.Json{
		"id": cron.ID,
	})
}

// Script 获取脚本内容
func (c *CronController) Script(ctx http.Context) {
	var cron models.Cron
	err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		Error(ctx, http.StatusBadRequest, "计划任务不存在")
		return
	}

	Success(ctx, tools.Read(cron.Shell))
}

func (c *CronController) Update(ctx http.Context) {
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
		Error(ctx, http.StatusBadRequest, validator.Errors().One())
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

	if !tools.Write(cron.Shell, ctx.Request().Input("script"), 0644) {
		facades.Log().Error("[面板][CronController] 更新计划任务脚本失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	tools.Exec("dos2unix " + cron.Shell)

	c.cron.DeleteFromSystem(cron)
	if cron.Status {
		c.cron.AddToSystem(cron)
	}

	Success(ctx, nil)
}

func (c *CronController) Delete(ctx http.Context) {
	var cron models.Cron
	err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		Error(ctx, http.StatusBadRequest, "计划任务不存在")
		return
	}

	c.cron.DeleteFromSystem(cron)
	tools.Remove(cron.Shell)

	_, err = facades.Orm().Query().Delete(&cron)
	if err != nil {
		facades.Log().Error("[面板][CronController] 删除计划任务失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}

func (c *CronController) Status(ctx http.Context) {
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

	c.cron.DeleteFromSystem(cron)
	if cron.Status {
		c.cron.AddToSystem(cron)
	}

	Success(ctx, nil)
}

func (c *CronController) Log(ctx http.Context) {
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

	Success(ctx, tools.Read(cron.Log))
}
