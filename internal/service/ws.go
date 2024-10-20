package service

import (
	"bufio"
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/ssh"
)

type WsService struct {
	sshRepo biz.SSHRepo
}

func NewWsService() *WsService {
	return &WsService{
		sshRepo: data.NewSSHRepo(),
	}
}

func (s *WsService) Session(w http.ResponseWriter, r *http.Request) {
	info, err := s.sshRepo.GetInfo()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	ws, err := s.upgrade(w, r)
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
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return
	}
	defer client.Close()

	turn, err := ssh.NewTurn(ws, client)
	if err != nil {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return
	}
	defer turn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		_ = turn.Handle(ctx)
	}()
	go func() {
		defer wg.Done()
		_ = turn.Wait()
	}()

	wg.Wait()
	cancel()
}

func (s *WsService) Exec(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrade(w, r)
	if err != nil {
		ErrorSystem(w)
		return
	}
	defer ws.Close()

	// 第一条消息是命令
	_, cmd, err := ws.ReadMessage()
	if err != nil {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "failed to read command"))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	out, err := shell.ExecfWithPipe(ctx, string(cmd))
	if err != nil {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "failed to run command"))
		cancel()
		return
	}

	go func() {
		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			line := scanner.Text()
			_ = ws.WriteMessage(websocket.TextMessage, []byte(line))
		}
		if err = scanner.Err(); err != nil {
			_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "failed to read command output"))
		}
	}()

	s.readLoop(ws)
	cancel()
}

func (s *WsService) upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upGrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}

	// debug 模式下不校验 origin，方便 vite 代理调试
	if app.Conf.Bool("app.debug") {
		upGrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	return upGrader.Upgrade(w, r, nil)
}

// readLoop 阻塞直到客户端关闭连接
func (s *WsService) readLoop(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			break
		}
	}
}
