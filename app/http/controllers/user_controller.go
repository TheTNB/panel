package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	models "panel/app/Models"
	"panel/app/http/requests"
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
		ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"code":    422,
			"message": err.Error(),
		})
		return
	}
	if errors != nil {
		ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"code":    422,
			"message": errors.All(),
		})
		return
	}

	var user models.User
	err = facades.Orm().Query().Where("username", loginRequest.Username).First(&user)
	if err != nil {
		facades.Log().Error("[面板][UserController] 查询用户失败 ", err)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": "系统内部错误",
		})
		return
	}

	if user.ID == 0 || !facades.Hash().Check(loginRequest.Password, user.Password) {
		ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 检查密码是否需要重新哈希
	if facades.Hash().NeedsRehash(user.Password) {
		// 更新密码
		user.Password, err = facades.Hash().Make(loginRequest.Password)
		if err != nil {
			facades.Log().Error("[面板][UserController] 更新密码失败 ", err)
			ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"code":    500,
				"message": "系统内部错误",
			})
			return
		}
	}

	token, loginErr := facades.Auth().LoginUsingID(ctx, user.ID)
	if loginErr != nil {
		facades.Log().Error("[面板][UserController] 登录失败 ", loginErr)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": loginErr.Error(),
		})
		return
	}

	ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "登录成功",
		"data": http.Json{
			"access_token": token,
		},
	})
}

func (r *UserController) Info(ctx http.Context) {
	ctx.Response().Success().Json(http.Json{
		"Hello": "Goravel",
	})
}
