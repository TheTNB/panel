package controllers

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/goravel/fiber"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
	"github.com/valyala/fasthttp"

	"panel/app/internal"
	"panel/app/internal/services"
	"panel/app/models"
	"panel/pkg/ssh"
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
	if err = r.setting.Set(models.SettingKeySshHost, host); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err = r.setting.Set(models.SettingKeySshPort, port); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err = r.setting.Set(models.SettingKeySshUser, user); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err = r.setting.Set(models.SettingKeySshPassword, password); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

// Session SSH 会话
func (r *SshController) Session(ctx http.Context) http.Response {
	upGrader := websocket.FastHTTPUpgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
			return true
		},
		Subprotocols: []string{ctx.Request().Header("Sec-WebSocket-Protocol")},
	}

	config := ssh.SSHClientConfigPassword(
		r.setting.Get(models.SettingKeySshHost)+":"+r.setting.Get(models.SettingKeySshPort),
		r.setting.Get(models.SettingKeySshUser),
		r.setting.Get(models.SettingKeySshPassword),
	)
	client, err := ssh.NewSSHClient(config)

	err = upGrader.Upgrade(ctx.(*fiber.Context).Instance().Context(), func(conn *websocket.Conn) {
		if err != nil {
			_ = conn.WriteControl(websocket.CloseMessage,
				[]byte(err.Error()), time.Now().Add(time.Second))
			return
		}

		defer client.Close()

		turn, err := ssh.NewTurn(conn, client)
		if err != nil {
			_ = conn.WriteControl(websocket.CloseMessage,
				[]byte(err.Error()), time.Now().Add(time.Second))
			return
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

		ctx2, cancel := context.WithCancel(context.Background())
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			err := turn.LoopRead(logBuff, ctx2)
			if err != nil {
				facades.Log().Request(ctx.Request()).Tags("面板", "SSH").With(map[string]any{
					"error": err.Error(),
				}).Info("SSH 会话错误")
			}
		}()
		go func() {
			defer wg.Done()
			err := turn.SessionWait()
			if err != nil {
				facades.Log().Request(ctx.Request()).Tags("面板", "SSH").With(map[string]any{
					"error": err.Error(),
				}).Info("SSH 会话错误")
			}
			cancel()
		}()
		wg.Wait()
	})
	if err != nil {
		return ErrorSystem(ctx)
	}

	return nil
}
