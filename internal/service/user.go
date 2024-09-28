package service

import (
	"net/http"

	"github.com/go-rat/chix"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/internal/panel"
)

type UserService struct {
	repo biz.UserRepo
}

func NewUserService() *UserService {
	return &UserService{
		repo: data.NewUserRepo(),
	}
}

func (s *UserService) Login(w http.ResponseWriter, r *http.Request) {
	sess, err := panel.Session.GetSession(r)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	req, err := Bind[request.UserLogin](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := s.repo.CheckPassword(req.Username, req.Password)
	if err != nil {
		Error(w, http.StatusForbidden, err.Error())
		return
	}

	sess.Put("user_id", user.ID)
	Success(w, nil)
}

func (s *UserService) Logout(w http.ResponseWriter, r *http.Request) {
	sess, err := panel.Session.GetSession(r)
	if err == nil {
		sess.Forget("user_id")
	}
	Success(w, nil)
}

func (s *UserService) IsLogin(w http.ResponseWriter, r *http.Request) {
	sess, err := panel.Session.GetSession(r)
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
