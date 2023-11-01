package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"panel/app/http/requests/user"
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
// @Summary 用户登录
// @Description 通过用户名和密码获取访问令牌
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body requests.Login true "登录信息"
// @Success 200 {object} SuccessResponse
// @Failure 403 {object} ErrorResponse "用户名或密码错误"
// @Failure 500 {object} ErrorResponse "系统内部错误
// @Router /panel/user/login [post]
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
		}).Tags("面板", "用户").Error("查询用户失败")
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
			}).Tags("面板", "用户").Error("更新密码失败")
			return ErrorSystem(ctx)
		}
	}

	token, loginErr := facades.Auth().LoginUsingID(ctx, user.ID)
	if loginErr != nil {
		facades.Log().Request(ctx.Request()).With(map[string]any{
			"error": err.Error(),
		}).Tags("面板", "用户").Error("登录失败")
		return ErrorSystem(ctx)
	}

	return Success(ctx, http.Json{
		"access_token": token,
	})
}

// Info 用户信息
func (r *UserController) Info(ctx http.Context) http.Response {
	var user models.User
	err := facades.Auth().User(ctx, &user)
	if err != nil {
		facades.Log().With(map[string]any{
			"error": err.Error(),
		}).Error("[面板][UserController] 查询用户信息失败")
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	return Success(ctx, http.Json{
		"id":       user.ID,
		"role":     []string{"admin"},
		"username": user.Username,
		"email":    user.Email,
	})
}
