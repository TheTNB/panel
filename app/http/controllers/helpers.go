package controllers

import "github.com/goravel/framework/contracts/http"

func Success(ctx http.Context, data any) {
	ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

func Error(ctx http.Context, code int, message any) {
	ctx.Response().Json(http.StatusOK, http.Json{
		"code":    code,
		"message": message,
	})
}
