package data

import (
	"errors"

	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
)

type sshRepo struct {
	settingRepo biz.SettingRepo
}

func NewSSHRepo() biz.SSHRepo {
	return &sshRepo{
		settingRepo: NewSettingRepo(),
	}
}

func (r *sshRepo) GetInfo() (map[string]any, error) {
	host, _ := r.settingRepo.Get(biz.SettingKeySshHost)
	port, _ := r.settingRepo.Get(biz.SettingKeySshPort)
	user, _ := r.settingRepo.Get(biz.SettingKeySshUser)
	password, _ := r.settingRepo.Get(biz.SettingKeySshPassword)
	if len(host) == 0 || len(user) == 0 || len(password) == 0 {
		return nil, errors.New("SSH 配置不完整")
	}

	return map[string]any{
		"host":     host,
		"port":     cast.ToInt(port),
		"user":     user,
		"password": password,
	}, nil
}

func (r *sshRepo) UpdateInfo(req *request.SSHUpdateInfo) error {
	if err := r.settingRepo.Set(biz.SettingKeySshHost, req.Host); err != nil {
		return err
	}
	if err := r.settingRepo.Set(biz.SettingKeySshPort, req.Port); err != nil {
		return err
	}
	if err := r.settingRepo.Set(biz.SettingKeySshUser, req.User); err != nil {
		return err
	}
	if err := r.settingRepo.Set(biz.SettingKeySshPassword, req.Password); err != nil {
		return err
	}

	return nil
}
