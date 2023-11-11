package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/http/requests/user"
	responses "panel/app/http/responses/user"
	"panel/app/models"
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
//	@Summary		登录
//	@Description	通过用户名和密码获取访问令牌
//	@Tags			用户鉴权
//	@Accept			json
//	@Produce		json
//	@Param			data	body		requests.Login	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Failure		403		{object}	ErrorResponse	"用户名或密码错误"
//	@Failure		500		{object}	ErrorResponse	"系统内部错误
//	@Router			/panel/user/login [post]
func (r *UserController) Login(ctx http.Context) http.Response {
	var loginRequest requests.Login
	sanitize := Sanitize(ctx, &loginRequest)
	if sanitize != nil {
		return sanitize
	}

	var user models.User
	err := facades.Orm().Query().Where("username", loginRequest.Username).First(&user)
	if err != nil {
		facades.Log().Request(ctx.Request()).With(map[string]any{
			"error": err.Error(),
		}).Tags("面板", "用户").Info("查询用户失败")
		return ErrorSystem(ctx)
	}

	if user.ID == 0 || !facades.Hash().Check(loginRequest.Password, user.Password) {
		return Error(ctx, http.StatusForbidden, "用户名或密码错误")
	}

	if facades.Hash().NeedsRehash(user.Password) {
		user.Password, err = facades.Hash().Make(loginRequest.Password)
		if err != nil {
			facades.Log().Request(ctx.Request()).With(map[string]any{
				"error": err.Error(),
			}).Tags("面板", "用户").Info("更新密码失败")
			return ErrorSystem(ctx)
		}
	}

	token, loginErr := facades.Auth().LoginUsingID(ctx, user.ID)
	if loginErr != nil {
		facades.Log().Request(ctx.Request()).With(map[string]any{
			"error": err.Error(),
		}).Tags("面板", "用户").Info("登录失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, http.Json{
		"access_token": token,
	})
}

// Info
//
//	@Summary		用户信息
//	@Description	获取当前登录用户信息
//	@Tags			用户鉴权
//	@Produce		json
//	@Security		BearerToken
//	@Success		200	{object}	SuccessResponse{data=responses.Info}
//	@Router			/panel/user/info [get]
func (r *UserController) Info(ctx http.Context) http.Response {
	var user models.User
	err := facades.Auth().User(ctx, &user)
	if err != nil {
		facades.Log().Request(ctx.Request()).With(map[string]any{
			"error": err.Error(),
		}).Tags("面板", "用户").Info("获取用户信息失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, responses.Info{
		ID:       user.ID,
		Role:     []string{"admin"},
		Username: user.Username,
		Email:    user.Email,
	})
}
