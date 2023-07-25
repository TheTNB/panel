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

func (r *SshController) GetInfo(ctx http.Context) {
	host := r.setting.Get(models.SettingKeySshHost)
	port := r.setting.Get(models.SettingKeySshPort)
	user := r.setting.Get(models.SettingKeySshUser)
	password := r.setting.Get(models.SettingKeySshPassword)
	if len(host) == 0 || len(user) == 0 || len(password) == 0 {
		Error(ctx, http.StatusInternalServerError, "SSH 配置不完整")
		return
	}

	Success(ctx, http.Json{
		"host":     host,
		"port":     port,
		"user":     user,
		"password": password,
	})
}

func (r *SshController) UpdateInfo(ctx http.Context) {
	validator, err := ctx.Request().Validate(map[string]string{
		"host":     "required",
		"port":     "required",
		"user":     "required",
		"password": "required",
	})
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if validator.Fails() {
		Error(ctx, http.StatusBadRequest, validator.Errors().One())
		return
	}

	host := ctx.Request().Input("host")
	port := ctx.Request().Input("port")
	user := ctx.Request().Input("user")
	password := ctx.Request().Input("password")
	err = r.setting.Set(models.SettingKeySshHost, host)
	if err != nil {
		facades.Log().Error("[面板][SSH] 更新配置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	err = r.setting.Set(models.SettingKeySshPort, port)
	if err != nil {
		facades.Log().Error("[面板][SSH] 更新配置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	err = r.setting.Set(models.SettingKeySshUser, user)
	if err != nil {
		facades.Log().Error("[面板][SSH] 更新配置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	err = r.setting.Set(models.SettingKeySshPassword, password)
	if err != nil {
		facades.Log().Error("[面板][SSH] 更新配置失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, nil)
}

func (r *SshController) Session(ctx http.Context) {
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
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	defer ws.Close()

	var config *ssh.SSHClientConfig
	config = ssh.SSHClientConfigPassword(
		r.setting.Get(models.SettingKeySshHost)+":"+r.setting.Get(models.SettingKeySshPort),
		r.setting.Get(models.SettingKeySshUser),
		r.setting.Get(models.SettingKeySshPassword),
	)

	client, err := ssh.NewSSHClient(config)
	if err != nil {
		ws.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return
	}
	defer client.Close()

	turn, err := ssh.NewTurn(ws, client)
	if err != nil {
		ws.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return
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
}
