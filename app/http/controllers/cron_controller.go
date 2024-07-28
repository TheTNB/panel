package controllers

import (
	"regexp"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/str"
)

type CronController struct {
	cron    internal.Cron
	setting internal.Setting
}

func NewCronController() *CronController {
	return &CronController{
		cron:    services.NewCronImpl(),
		setting: services.NewSettingImpl(),
	}
}

// List 获取计划任务列表
func (r *CronController) List(ctx http.Context) http.Response {
	limit := ctx.Request().QueryInt("limit", 10)
	page := ctx.Request().QueryInt("page", 1)

	var crons []models.Cron
	var total int64
	err := facades.Orm().Query().Paginate(page, limit, &crons, &total)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "计划任务").With(map[string]any{
			"error": err.Error(),
		}).Info("查询计划任务列表失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, http.Json{
		"total": total,
		"items": crons,
	})
}

// Add 添加计划任务
func (r *CronController) Add(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"name":        "required|min_len:1|max_len:255",
		"time":        "required",
		"script":      "required",
		"type":        "required|in:shell,backup,cutoff",
		"backup_type": "required_if:type,backup|in:website,mysql,postgresql",
	}); sanitize != nil {
		return sanitize
	}

	// 单独验证时间格式
	if !regexp.MustCompile(`^((\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+)(,(\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+))*\s?){5}$`).MatchString(ctx.Request().Input("time")) {
		return h.Error(ctx, http.StatusUnprocessableEntity, "时间格式错误")
	}

	script := ctx.Request().Input("script")
	cronType := ctx.Request().Input("type")
	if cronType == "backup" {
		backupType := ctx.Request().Input("backup_type")
		backupName := ctx.Request().Input("database")
		if backupType == "website" {
			backupName = ctx.Request().Input("website")
		}
		backupPath := ctx.Request().Input("backup_path")
		if len(backupPath) == 0 {
			backupPath = r.setting.Get(models.SettingKeyBackupPath) + "/" + backupType
		}
		backupSave := ctx.Request().InputInt("save", 10)
		script = `#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

# 耗子面板 - 数据备份脚本

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
		script = `#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

# 耗子面板 - 日志切割脚本

name=` + website + `
save=` + cast.ToString(save) + `

# 执行切割
panel cutoff ${name} ${save} 2>&1
`
	}

	shellDir := "/www/server/cron/"
	shellLogDir := "/www/server/cron/logs/"
	if !io.Exists(shellDir) {
		return h.Error(ctx, http.StatusInternalServerError, "计划任务目录不存在")
	}
	if !io.Exists(shellLogDir) {
		return h.Error(ctx, http.StatusInternalServerError, "计划任务日志目录不存在")
	}
	shellFile := strconv.Itoa(int(carbon.Now().Timestamp())) + str.RandomString(16)
	if err := io.Write(shellDir+shellFile+".sh", script, 0700); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if out, err := shell.Execf("dos2unix " + shellDir + shellFile + ".sh"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	var cron models.Cron
	cron.Name = ctx.Request().Input("name")
	cron.Type = ctx.Request().Input("type")
	cron.Status = true
	cron.Time = ctx.Request().Input("time")
	cron.Shell = shellDir + shellFile + ".sh"
	cron.Log = shellLogDir + shellFile + ".log"

	if err := facades.Orm().Query().Create(&cron); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "计划任务").With(map[string]any{
			"error": err.Error(),
		}).Info("保存计划任务失败")
		return h.ErrorSystem(ctx)
	}

	if err := r.cron.AddToSystem(cron); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, http.Json{
		"id": cron.ID,
	})
}

// Script 获取脚本内容
func (r *CronController) Script(ctx http.Context) http.Response {
	var cron models.Cron
	err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron)
	if err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "计划任务不存在")
	}

	script, err := io.Read(cron.Shell)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, script)
}

// Update 更新计划任务
func (r *CronController) Update(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"name":   "required|min_len:1|max_len:255",
		"time":   "required",
		"script": "required",
	}); sanitize != nil {
		return sanitize
	}

	// 单独验证时间格式
	if !regexp.MustCompile(`^((\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+)(,(\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+))*\s?){5}$`).MatchString(ctx.Request().Input("time")) {
		return h.Error(ctx, http.StatusUnprocessableEntity, "时间格式错误")
	}

	var cron models.Cron
	if err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron); err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "计划任务不存在")
	}

	if !cron.Status {
		return h.Error(ctx, http.StatusUnprocessableEntity, "计划任务已禁用")
	}

	cron.Time = ctx.Request().Input("time")
	cron.Name = ctx.Request().Input("name")
	if err := facades.Orm().Query().Save(&cron); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "计划任务").With(map[string]any{
			"error": err.Error(),
		}).Info("更新计划任务失败")
		return h.ErrorSystem(ctx)
	}

	if err := io.Write(cron.Shell, ctx.Request().Input("script"), 0644); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if out, err := shell.Execf("dos2unix " + cron.Shell); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}

	if err := r.cron.DeleteFromSystem(cron); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if cron.Status {
		if err := r.cron.AddToSystem(cron); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}
	}

	return h.Success(ctx, nil)
}

// Delete 删除计划任务
func (r *CronController) Delete(ctx http.Context) http.Response {
	var cron models.Cron
	if err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron); err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "计划任务不存在")
	}

	if err := r.cron.DeleteFromSystem(cron); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err := io.Remove(cron.Shell); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if _, err := facades.Orm().Query().Delete(&cron); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "计划任务").With(map[string]any{
			"error": err.Error(),
		}).Info("删除计划任务失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, nil)
}

// Status 更新计划任务状态
func (r *CronController) Status(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"status": "bool",
	}); sanitize != nil {
		return sanitize
	}

	var cron models.Cron
	if err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron); err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "计划任务不存在")
	}

	cron.Status = ctx.Request().InputBool("status")
	if err := facades.Orm().Query().Save(&cron); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "计划任务").With(map[string]any{
			"error": err.Error(),
		}).Info("更新计划任务状态失败")
		return h.ErrorSystem(ctx)
	}

	if err := r.cron.DeleteFromSystem(cron); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if cron.Status {
		if err := r.cron.AddToSystem(cron); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}
	}

	return h.Success(ctx, nil)
}

// Log 获取计划任务日志
func (r *CronController) Log(ctx http.Context) http.Response {
	var cron models.Cron
	if err := facades.Orm().Query().Where("id", ctx.Request().Input("id")).FirstOrFail(&cron); err != nil {
		return h.Error(ctx, http.StatusUnprocessableEntity, "计划任务不存在")
	}

	if !io.Exists(cron.Log) {
		return h.Error(ctx, http.StatusUnprocessableEntity, "日志文件不存在")
	}

	log, err := shell.Execf("tail -n 1000 " + cron.Log)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, log)
}
