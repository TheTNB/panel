package service

import (
	"fmt"
	"net/http"

	"github.com/go-rat/chix"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
)

type UserService struct {
	repo biz.UserRepo
}

func NewUserService() *UserService {
	return &UserService{
		repo: data.NewUserRepo(),
	}
}

// Login
//
//	@Summary	登录
//	@Tags		用户服务
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.UserLogin	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Router		/user/login [post]
func (s *UserService) Login(w http.ResponseWriter, r *http.Request) {
	sess, err := app.Session.GetSession(r)
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

// Logout
//
//	@Summary	登出
//	@Tags		用户服务
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	SuccessResponse
//	@Router		/user/logout [post]
func (s *UserService) Logout(w http.ResponseWriter, r *http.Request) {
	sess, err := app.Session.GetSession(r)
	if err == nil {
		sess.Forget("user_id")
	}
	Success(w, nil)
}

// IsLogin
//
//	@Summary	是否登录
//	@Tags		用户服务
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	SuccessResponse
//	@Router		/user/isLogin [get]
func (s *UserService) IsLogin(w http.ResponseWriter, r *http.Request) {
	sess, err := app.Session.GetSession(r)
	if err != nil {
		Success(w, false)
		return
	}
	Success(w, sess.Has("user_id"))
}

// Info
//
//	@Summary	用户信息
//	@Tags		用户服务
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	SuccessResponse
//	@Router		/user/info/{id} [get]
func (s *UserService) Info(w http.ResponseWriter, r *http.Request) {
	userID := cast.ToUint(r.Context().Value("user_id"))
	fmt.Println(userID)
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
