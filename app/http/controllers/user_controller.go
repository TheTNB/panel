package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/http/requests"
	"panel/app/models"
)

type UserController struct {
	//Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		//Inject services
	}
}

func (r *UserController) Login(ctx http.Context) {
	var loginRequest requests.LoginRequest
	errors, err := ctx.Request().ValidateRequest(&loginRequest)
	if err != nil {
		Error(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if errors != nil {
		Error(ctx, http.StatusUnprocessableEntity, errors.All())
		return
	}

	var user models.User
	err = facades.Orm().Query().Where("username", loginRequest.Username).First(&user)
	if err != nil {
		facades.Log().Error("[面板][UserController] 查询用户失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	if user.ID == 0 || !facades.Hash().Check(loginRequest.Password, user.Password) {
		Error(ctx, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	if facades.Hash().NeedsRehash(user.Password) {
		user.Password, err = facades.Hash().Make(loginRequest.Password)
		if err != nil {
			facades.Log().Error("[面板][UserController] 更新密码失败 ", err)
			Error(ctx, http.StatusInternalServerError, "系统内部错误")
			return
		}
	}

	token, loginErr := facades.Auth().LoginUsingID(ctx, user.ID)
	if loginErr != nil {
		facades.Log().Error("[面板][UserController] 登录失败 ", loginErr)
		Error(ctx, http.StatusInternalServerError, loginErr.Error())
		return
	}

	Success(ctx, http.Json{
		"access_token": token,
	})
}

func (r *UserController) Info(ctx http.Context) {
	user, ok := ctx.Value("user").(models.User)
	if !ok {
		Error(ctx, http.StatusUnauthorized, "登录已过期")
		return
	}

	Success(ctx, http.Json{
		"username": user.Username,
		"email":    user.Email,
	})
}
