package data

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/golang-module/carbon/v2"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/os"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/str"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type cronRepo struct {
	settingRepo biz.SettingRepo
}

func NewCronRepo() biz.CronRepo {
	return &cronRepo{
		settingRepo: NewSettingRepo(),
	}
}

func (r *cronRepo) Count() (int64, error) {
	var count int64
	if err := app.Orm.Model(&biz.Cron{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *cronRepo) List(page, limit uint) ([]*biz.Cron, int64, error) {
	var cron []*biz.Cron
	var total int64
	err := app.Orm.Model(&biz.Cert{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&cron).Error
	return cron, total, err
}

func (r *cronRepo) Get(id uint) (*biz.Cron, error) {
	cron := new(biz.Cron)
	if err := app.Orm.Where("id = ?", id).First(cron).Error; err != nil {
		return nil, err
	}

	return cron, nil
}

func (r *cronRepo) Create(req *request.CronCreate) error {
	if !regexp.MustCompile(`^((\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+)(,(\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+))*\s?){5}$`).MatchString(req.Time) {
		return errors.New("时间格式错误")
	}

	var script string
	if req.Type == "backup" {
		if len(req.BackupPath) == 0 {
			req.BackupPath, _ = r.settingRepo.Get(biz.SettingKeyBackupPath)
			if len(req.BackupPath) == 0 {
				return errors.New("备份路径不能为空")
			}
			req.BackupPath = filepath.Join(req.BackupPath, req.BackupType)
		}
		script = fmt.Sprintf(`#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

# 耗子面板 - 数据备份脚本

type=%s
path=%s
name=%s
save=%d

# 执行备份
panel backup ${type} ${name} ${path} ${save} 2>&1
`, req.BackupType, req.BackupPath, req.Target, req.Save)
	}
	if req.Type == "cutoff" {
		script = fmt.Sprintf(`#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

# 耗子面板 - 日志切割脚本

name=%s
save=%d

# 执行切割
panel cutoff ${name} ${save} 2>&1
`, req.Target, req.Save)
	}

	shellDir := fmt.Sprintf("%s/server/cron/", app.Root)
	shellLogDir := fmt.Sprintf("%s/server/cron/logs/", app.Root)
	if !io.Exists(shellDir) {
		return errors.New("计划任务目录不存在")
	}
	if !io.Exists(shellLogDir) {
		return errors.New("计划任务日志目录不存在")
	}
	shellFile := strconv.Itoa(int(carbon.Now().Timestamp())) + str.RandomString(16)
	if err := io.Write(filepath.Join(shellDir, shellFile+".sh"), script, 0700); err != nil {
		return errors.New(err.Error())
	}
	if out, err := shell.Execf("dos2unix %s%s.sh", shellDir, shellFile); err != nil {
		return errors.New(out)
	}

	cron := new(biz.Cron)
	cron.Name = req.Name
	cron.Type = req.Type
	cron.Status = true
	cron.Time = req.Time
	cron.Shell = shellDir + shellFile + ".sh"
	cron.Log = shellLogDir + shellFile + ".log"

	if err := app.Orm.Create(cron).Error; err != nil {
		return err
	}
	if err := r.addToSystem(cron); err != nil {
		return err
	}

	return nil
}

func (r *cronRepo) Update(req *request.CronUpdate) error {
	cron, err := r.Get(req.ID)
	if err != nil {
		return err
	}

	if !regexp.MustCompile(`^((\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+)(,(\*|\d+|\d+-\d+|\d+/\d+|\d+-\d+/\d+|\*/\d+))*\s?){5}$`).MatchString(req.Time) {
		return errors.New("时间格式错误")
	}

	if !cron.Status {
		return errors.New("计划任务已禁用")
	}

	cron.Time = req.Time
	cron.Name = req.Name
	if err = app.Orm.Save(cron).Error; err != nil {
		return err
	}

	if err = io.Write(cron.Shell, req.Script, 0700); err != nil {
		return err
	}
	if out, err := shell.Execf("dos2unix %s", cron.Shell); err != nil {
		return errors.New(out)
	}

	if err = r.deleteFromSystem(cron); err != nil {
		return err
	}
	if cron.Status {
		if err = r.addToSystem(cron); err != nil {
			return err
		}
	}

	return nil
}

func (r *cronRepo) Delete(id uint) error {
	cron, err := r.Get(id)
	if err != nil {
		return err
	}

	if err = r.deleteFromSystem(cron); err != nil {
		return err
	}
	if err = io.Remove(cron.Shell); err != nil {
		return err
	}

	return app.Orm.Delete(cron).Error
}

func (r *cronRepo) Status(id uint, status bool) error {
	cron, err := r.Get(id)
	if err != nil {
		return err
	}

	if err = r.deleteFromSystem(cron); err != nil {
		return err
	}
	if status {
		return r.addToSystem(cron)
	}

	cron.Status = status

	return app.Orm.Save(cron).Error
}

func (r *cronRepo) Log(id uint) (string, error) {
	cron, err := r.Get(id)
	if err != nil {
		return "", err
	}

	if !io.Exists(cron.Log) {
		return "", errors.New("日志文件不存在")
	}

	log, err := shell.Execf("tail -n 1000 '%s'", cron.Log)
	if err != nil {
		return "", err
	}

	return log, nil
}

// addToSystem 添加到系统
func (r *cronRepo) addToSystem(cron *biz.Cron) error {
	if _, err := shell.Execf(`( crontab -l; echo "%s %s >> %s 2>&1" ) | sort - | uniq - | crontab -`, cron.Time, cron.Shell, cron.Log); err != nil {
		return err
	}

	return r.restartCron()
}

// deleteFromSystem 从系统中删除
func (r *cronRepo) deleteFromSystem(cron *biz.Cron) error {
	if _, err := shell.Execf(`( crontab -l | grep -v -F "%s %s >> %s 2>&1" ) | crontab -`, cron.Time, cron.Shell, cron.Log); err != nil {
		return err
	}

	return r.restartCron()
}

// restartCron 重启 cron 服务
func (r *cronRepo) restartCron() error {
	if os.IsRHEL() {
		return systemctl.Restart("crond")
	}

	if os.IsDebian() || os.IsUbuntu() {
		return systemctl.Restart("cron")
	}

	return errors.New("不支持的系统")
}
