package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/app/http/requests/user"
	"github.com/TheTNB/panel/v2/app/models"
	"github.com/TheTNB/panel/v2/pkg/h"
)

type UserController struct {
	// Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		// Inject services
	}
}

// Login
//
//	@Summary	登录
//	@Tags		用户鉴权
//	@Accept		json
//	@Produce	json
//	@Param		data	body		requests.Login	true	"request"
//	@Success	200		{object}	SuccessResponse
//	@Failure	403		{object}	ErrorResponse	"用户名或密码错误"
//	@Failure	500		{object}	ErrorResponse	"系统内部错误
//	@Router		/panel/user/login [post]
func (r *UserController) Login(ctx http.Context) http.Response {
	var loginRequest requests.Login
	sanitize := h.SanitizeRequest(ctx, &loginRequest)
	if sanitize != nil {
		return sanitize
	}

	var user models.User
	err := facades.Orm().Query().Where("username", loginRequest.Username).First(&user)
	if err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "用户").With(map[string]any{
			"error": err.Error(),
		}).Info("查询用户失败")
		return h.ErrorSystem(ctx)
	}

	if user.ID == 0 || !facades.Hash().Check(loginRequest.Password, user.Password) {
		return h.Error(ctx, http.StatusForbidden, "用户名或密码错误")
	}

	if facades.Hash().NeedsRehash(user.Password) {
		user.Password, err = facades.Hash().Make(loginRequest.Password)
		if err != nil {
			facades.Log().Request(ctx.Request()).Tags("面板", "用户").With(map[string]any{
				"error": err.Error(),
			}).Info("更新密码失败")
			return h.ErrorSystem(ctx)
		}
	}

	ctx.Request().Session().Put("user_id", user.ID)
	return h.Success(ctx, nil)
}

// Logout
//
//	@Summary	登出
//	@Tags		用户鉴权
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/user/logout [post]
func (r *UserController) Logout(ctx http.Context) http.Response {
	if ctx.Request().HasSession() {
		ctx.Request().Session().Forget("user_id")
	}

	return h.Success(ctx, nil)
}

// IsLogin
//
//	@Summary	是否登录
//	@Tags		用户鉴权
//	@Produce	json
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/user/isLogin [get]
func (r *UserController) IsLogin(ctx http.Context) http.Response {
	if !ctx.Request().HasSession() {
		return h.Success(ctx, false)
	}

	return h.Success(ctx, ctx.Request().Session().Has("user_id"))
}

// Info
//
//	@Summary	用户信息
//	@Tags		用户鉴权
//	@Produce	json
//	@Security	BearerToken
//	@Success	200	{object}	SuccessResponse
//	@Router		/panel/user/info [get]
func (r *UserController) Info(ctx http.Context) http.Response {
	userID := cast.ToUint(ctx.Value("user_id"))
	var user models.User
	if err := facades.Orm().Query().Where("id", userID).Get(&user); err != nil {
		facades.Log().Request(ctx.Request()).Tags("面板", "用户").With(map[string]any{
			"error": err.Error(),
		}).Info("获取用户信息失败")
		return h.ErrorSystem(ctx)
	}

	return h.Success(ctx, http.Json{
		"id":       user.ID,
		"role":     []string{"admin"},
		"username": user.Username,
		"email":    user.Email,
	})
}
