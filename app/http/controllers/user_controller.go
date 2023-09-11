package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/http/requests"
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

// Login 用户登录
func (r *UserController) Login(ctx http.Context) http.Response {
	var loginRequest requests.LoginRequest
	errors, err := ctx.Request().ValidateRequest(&loginRequest)
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if errors != nil {
		return Error(ctx, http.StatusUnprocessableEntity, errors.All())
	}

	var user models.User
	err = facades.Orm().Query().Where("username", loginRequest.Username).First(&user)
	if err != nil {
		facades.Log().Error("[面板][UserController] 查询用户失败 ", err)
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if user.ID == 0 || !facades.Hash().Check(loginRequest.Password, user.Password) {
		return Error(ctx, http.StatusForbidden, "用户名或密码错误")
	}

	if facades.Hash().NeedsRehash(user.Password) {
		user.Password, err = facades.Hash().Make(loginRequest.Password)
		if err != nil {
			facades.Log().Error("[面板][UserController] 更新密码失败 ", err)
			return Error(ctx, http.StatusInternalServerError, "系统内部错误")
		}
	}

	token, loginErr := facades.Auth().LoginUsingID(ctx, user.ID)
	if loginErr != nil {
		facades.Log().Error("[面板][UserController] 登录失败 ", loginErr)
		return Error(ctx, http.StatusInternalServerError, loginErr.Error())
	}

	return Success(ctx, http.Json{
		"access_token": token,
	})
}

// Info 用户信息
func (r *UserController) Info(ctx http.Context) http.Response {
	user, ok := ctx.Value("user").(models.User)
	if !ok {
		return Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	return Success(ctx, http.Json{
		"username": user.Username,
		"email":    user.Email,
	})
}
