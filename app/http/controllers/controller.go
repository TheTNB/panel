package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

// SuccessResponse 通用成功响应
type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// ErrorResponse 通用错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Success 响应成功
func Success(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(&SuccessResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 响应错误
func Error(ctx http.Context, code int, message string) http.Response {
	return ctx.Response().Json(http.StatusOK, &ErrorResponse{
		Code:    code,
		Message: "错误: " + message,
	})
}

// ErrorSystem 响应系统错误
func ErrorSystem(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusOK, &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "系统内部错误",
	})
}

// Sanitize 消毒请求参数
func Sanitize(ctx http.Context, request http.FormRequest) http.Response {
	errors, err := ctx.Request().ValidateRequest(request)
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if errors != nil {
		return Error(ctx, http.StatusUnprocessableEntity, errors.One())
	}

	return nil
}
