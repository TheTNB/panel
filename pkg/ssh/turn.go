package ssh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type MessageResize struct {
	Resize  bool `json:"resize"`
	Columns int  `json:"columns"`
	Rows    int  `json:"rows"`
}

type Turn struct {
	stdin   io.WriteCloser
	session *ssh.Session
	ws      *websocket.Conn
}

func NewTurn(ws *websocket.Conn, client *ssh.Client) (*Turn, error) {
	sess, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	stdin, err := sess.StdinPipe()
	if err != nil {
		return nil, err
	}

	turn := &Turn{stdin: stdin, session: sess, ws: ws}
	sess.Stdout = turn
	sess.Stderr = turn

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err = sess.RequestPty("xterm", 150, 80, modes); err != nil {
		return nil, err
	}
	if err = sess.Shell(); err != nil {
		return nil, err
	}

	return turn, nil
}

func (t *Turn) Write(p []byte) (n int, err error) {
	writer, err := t.ws.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return 0, err
	}
	defer writer.Close()

	return writer.Write(p)
}

func (t *Turn) Close() error {
	if t.session != nil {
		_ = t.session.Close()
	}

	return t.ws.Close()
}

func (t *Turn) Handle(context context.Context) error {
	var resize MessageResize
	for {
		select {
		case <-context.Done():
			return errors.New("ssh context done exit")
		default:
			_, data, err := t.ws.ReadMessage()
			if err != nil {
				return fmt.Errorf("reading ws message err: %v", err)
			}

			// 判断是否是 resize 消息
			if err = json.Unmarshal(data, &resize); err == nil {
				if resize.Resize && resize.Columns > 0 && resize.Rows > 0 {
					if err = t.session.WindowChange(resize.Rows, resize.Columns); err != nil {
						return fmt.Errorf("change window size err: %v", err)
					}
				}
				continue
			}

			if _, err = t.stdin.Write(data); err != nil {
				return fmt.Errorf("writing ws message to stdin err: %v", err)
			}
		}
	}
}

func (t *Turn) Wait() error {
	return t.session.Wait()
}
