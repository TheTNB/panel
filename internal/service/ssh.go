package service

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"

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
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
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
		ErrorSystem(w)
		return
	}
	defer client.Close()

	turn, err := ssh.NewTurn(ws, client)
	if err != nil {
		_ = ws.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		ErrorSystem(w)
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

	sshCtx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err = turn.LoopRead(logBuff, sshCtx); err != nil {
			ErrorSystem(w)
			return
		}
	}()
	go func() {
		defer wg.Done()
		if err = turn.SessionWait(); err != nil {
			ErrorSystem(w)
			return
		}
		cancel()
	}()
	wg.Wait()

}
