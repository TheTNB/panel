package service

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/go-rat/chix"
	"github.com/spf13/cast"
	"golang.org/x/crypto/sha3"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/rsacrypto"
)

type UserService struct {
	repo biz.UserRepo
}

func NewUserService() *UserService {
	return &UserService{
		repo: data.NewUserRepo(),
	}
}

func (s *UserService) GetKey(w http.ResponseWriter, r *http.Request) {
	key, err := rsacrypto.GenerateKey()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	sess, err := app.Session.GetSession(r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	encoded, err := json.Marshal(key)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	sess.Put("key", encoded)

	pk, err := rsacrypto.PublicKeyToString(&key.PublicKey)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, pk)
}

func (s *UserService) Login(w http.ResponseWriter, r *http.Request) {
	sess, err := app.Session.GetSession(r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	req, err := Bind[request.UserLogin](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	key := new(rsa.PrivateKey)
	if err = json.Unmarshal(sess.Get("key").([]byte), key); err != nil {
		Error(w, http.StatusForbidden, "invalid key, please refresh the page")
		return
	}

	decryptedUsername, _ := rsacrypto.DecryptData(key, req.Username)
	decryptedPassword, _ := rsacrypto.DecryptData(key, req.Password)
	user, err := s.repo.CheckPassword(string(decryptedUsername), string(decryptedPassword))
	if err != nil {
		Error(w, http.StatusForbidden, "%v", err)
		return
	}

	// 安全登录模式下，将当前客户端与会话绑定
	// 安全登录模式只在未启用TLS时生效，因为TLS本身就是安全的
	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	if req.SafeLogin && !app.Conf.Bool("http.tls") {
		ua := r.Header.Get("User-Agent")
		sess.Put("safe_login", true)
		sess.Put("safe_client", fmt.Sprintf("%x", sha3.Sum256([]byte(ip+"|"+ua))))
	}

	sess.Put("user_id", user.ID)
	sess.Forget("key")
	Success(w, nil)
}

func (s *UserService) Logout(w http.ResponseWriter, r *http.Request) {
	sess, err := app.Session.GetSession(r)
	if err == nil {
		if err = sess.Invalidate(); err != nil {
			Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
	}

	Success(w, nil)
}

func (s *UserService) IsLogin(w http.ResponseWriter, r *http.Request) {
	sess, err := app.Session.GetSession(r)
	if err != nil {
		Success(w, false)
		return
	}
	Success(w, sess.Has("user_id"))
}

func (s *UserService) Info(w http.ResponseWriter, r *http.Request) {
	userID := cast.ToUint(r.Context().Value("user_id"))
	if userID == 0 {
		ErrorSystem(w)
		return
	}

	user, err := s.repo.Get(userID)
	if err != nil {
		ErrorSystem(w)
		return
	}

	Success(w, chix.M{
		"id":       user.ID,
		"role":     []string{"admin"},
		"username": user.Username,
		"email":    user.Email,
	})
}
