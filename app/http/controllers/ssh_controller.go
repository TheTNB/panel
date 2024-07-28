package controllers

import (
	"bytes"
	"context"
	nethttp "net/http"
	"sync"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/ssh"
)

type SshController struct {
	AuthMethod ssh.AuthMethod
	setting    internal.Setting
}

func NewSshController() *SshController {
	return &SshController{
		AuthMethod: ssh.PASSWORD,
		setting:    services.NewSettingImpl(),
	}
}

// GetInfo 获取 SSH 配置
func (r *SshController) GetInfo(ctx http.Context) http.Response {
	host := r.setting.Get(models.SettingKeySshHost)
	port := r.setting.Get(models.SettingKeySshPort)
	user := r.setting.Get(models.SettingKeySshUser)
	password := r.setting.Get(models.SettingKeySshPassword)
	if len(host) == 0 || len(user) == 0 || len(password) == 0 {
		return h.Error(ctx, http.StatusInternalServerError, "SSH 配置不完整")
	}

	return h.Success(ctx, http.Json{
		"host":     host,
		"port":     cast.ToInt(port),
		"user":     user,
		"password": password,
	})
}

// UpdateInfo 更新 SSH 配置
func (r *SshController) UpdateInfo(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"host":     "required",
		"port":     "required",
		"user":     "required",
		"password": "required",
	}); sanitize != nil {
		return sanitize
	}

	host := ctx.Request().Input("host")
	port := ctx.Request().Input("port")
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	if err := r.setting.Set(models.SettingKeySshHost, host); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err := r.setting.Set(models.SettingKeySshPort, port); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err := r.setting.Set(models.SettingKeySshUser, user); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err := r.setting.Set(models.SettingKeySshPassword, password); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// Session SSH 会话
func (r *SshController) Session(ctx http.Context) http.Response {
	upGrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *nethttp.Request) bool {
			return true
		},
		Subprotocols: []string{ctx.Request().Header("Sec-WebSocket-Protocol")},
	}

	ws, err := upGrader.Upgrade(ctx.Response().Writer(), ctx.Request().Origin(), nil)
	if err != nil {
		facades.Log().Tags("面板", "SSH").With(map[string]any{
			"error": err.Error(),
		}).Infof("建立连接失败")
		return h.ErrorSystem(ctx)
	}
	defer ws.Close()

	config := ssh.ClientConfigPassword(
		r.setting.Get(models.SettingKeySshHost)+":"+r.setting.Get(models.SettingKeySshPort),
		r.setting.Get(models.SettingKeySshUser),
		r.setting.Get(models.SettingKeySshPassword),
	)
	client, err := ssh.NewSSHClient(config)

	if err != nil {
		_ = ws.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return h.ErrorSystem(ctx)
	}
	defer client.Close()

	turn, err := ssh.NewTurn(ws, client)
	if err != nil {
		_ = ws.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return h.ErrorSystem(ctx)
	}
	defer turn.Close()

	var bufPool = sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}
	var logBuff = bufPool.Get().(*bytes.Buffer)
	logBuff.Reset()
	defer bufPool.Put(logBuff)

	sshCtx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err = turn.LoopRead(logBuff, sshCtx); err != nil {
			facades.Log().Tags("面板", "SSH").With(map[string]any{
				"error": err.Error(),
			}).Infof("读取数据失败")
		}
	}()
	go func() {
		defer wg.Done()
		if err = turn.SessionWait(); err != nil {
			facades.Log().Tags("面板", "SSH").With(map[string]any{
				"error": err.Error(),
			}).Infof("会话错误")
		}
		cancel()
	}()
	wg.Wait()

	return nil
}
