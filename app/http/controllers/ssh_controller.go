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

	"panel/app/models"
	"panel/app/services"
	"panel/pkg/ssh"
)

type SshController struct {
	AuthMethod ssh.AuthMethod
	setting    services.Setting
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
		return Error(ctx, http.StatusInternalServerError, "SSH 配置不完整")
	}

	return Success(ctx, http.Json{
		"host":     host,
		"port":     cast.ToInt(port),
		"user":     user,
		"password": password,
	})
}

// UpdateInfo 更新 SSH 配置
func (r *SshController) UpdateInfo(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"host":     "required",
		"port":     "required",
		"user":     "required",
		"password": "required",
	})
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	host := ctx.Request().Input("host")
	port := ctx.Request().Input("port")
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	err = r.setting.Set(models.SettingKeySshHost, host)
	if err != nil {
		facades.Log().Error("[面板][SSH] 更新配置失败 ", err)
		return ErrorSystem(ctx)
	}
	err = r.setting.Set(models.SettingKeySshPort, port)
	if err != nil {
		facades.Log().Error("[面板][SSH] 更新配置失败 ", err)
		return ErrorSystem(ctx)
	}
	err = r.setting.Set(models.SettingKeySshUser, user)
	if err != nil {
		facades.Log().Error("[面板][SSH] 更新配置失败 ", err)
		return ErrorSystem(ctx)
	}
	err = r.setting.Set(models.SettingKeySshPassword, password)
	if err != nil {
		facades.Log().Error("[面板][SSH] 更新配置失败 ", err)
		return ErrorSystem(ctx)
	}

	return Success(ctx, nil)
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
		facades.Log().Error("[面板][SSH] 建立连接失败 ", err)
		return ErrorSystem(ctx)
	}
	defer ws.Close()

	config := ssh.SSHClientConfigPassword(
		r.setting.Get(models.SettingKeySshHost)+":"+r.setting.Get(models.SettingKeySshPort),
		r.setting.Get(models.SettingKeySshUser),
		r.setting.Get(models.SettingKeySshPassword),
	)
	client, err := ssh.NewSSHClient(config)
	if err != nil {
		_ = ws.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return ErrorSystem(ctx)
	}
	defer client.Close()

	turn, err := ssh.NewTurn(ws, client)
	if err != nil {
		_ = ws.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return ErrorSystem(ctx)
	}
	defer turn.Close()

	var bufPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	var logBuff = bufPool.Get().(*bytes.Buffer)
	logBuff.Reset()
	defer bufPool.Put(logBuff)

	ctx2, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := turn.LoopRead(logBuff, ctx2)
		if err != nil {
			facades.Log().Error("[面板][SSH] 读取数据失败 ", err.Error())
		}
	}()
	go func() {
		defer wg.Done()
		err := turn.SessionWait()
		if err != nil {
			facades.Log().Error("[面板][SSH] 会话失败 ", err.Error())
		}
		cancel()
	}()
	wg.Wait()

	return Success(ctx, logBuff.String())
}
