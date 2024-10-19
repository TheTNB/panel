package service

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/ssh"
)

type SSHService struct {
	sshRepo biz.SSHRepo
}

func NewSSHService() *SSHService {
	return &SSHService{
		sshRepo: data.NewSSHRepo(),
	}
}

func (s *SSHService) GetInfo(w http.ResponseWriter, r *http.Request) {
	info, err := s.sshRepo.GetInfo()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, info)
}

func (s *SSHService) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SSHUpdateInfo](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.sshRepo.UpdateInfo(req); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
}

func (s *SSHService) Session(w http.ResponseWriter, r *http.Request) {
	info, err := s.sshRepo.GetInfo()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	upGrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}

	ws, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		ErrorSystem(w)
		return
	}
	defer ws.Close()

	config := ssh.ClientConfigPassword(
		cast.ToString(info["host"])+":"+cast.ToString(info["port"]),
		cast.ToString(info["user"]),
		cast.ToString(info["password"]),
	)
	client, err := ssh.NewSSHClient(config)
	if err != nil {
		_ = ws.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return
	}
	defer client.Close()

	turn, err := ssh.NewTurn(ws, client)
	if err != nil {
		_ = ws.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return
	}
	defer turn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err = turn.Handle(ctx); err != nil {
			app.Logger.Error("读取 ssh 数据失败", zap.Error(err))
			return
		}
	}()
	go func() {
		defer wg.Done()
		if err = turn.Wait(); err != nil {
			app.Logger.Error("保持 ssh 会话失败", zap.Error(err))
		}
		cancel()
	}()

	wg.Wait()
}
