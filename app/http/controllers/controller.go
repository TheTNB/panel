package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	commonrequests "github.com/TheTNB/panel/app/http/requests/common"
)

// SuccessResponse 通用成功响应
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// ErrorResponse 通用错误响应
type ErrorResponse struct {
	Message string `json:"message"`
}

// Success 响应成功
func Success(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(&SuccessResponse{
		Message: "success",
		Data:    data,
	})
}

// Error 响应错误
func Error(ctx http.Context, code int, message string) http.Response {
	return ctx.Response().Json(code, &ErrorResponse{
		Message: facades.Lang(ctx).Get("messages.mistake") + ": " + message,
	})
}

// ErrorSystem 响应系统错误
func ErrorSystem(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusInternalServerError, &ErrorResponse{
		Message: facades.Lang(ctx).Get("errors.internal"),
	})
}

func Paginate[T any](ctx http.Context, allItems []T) (pagedItems []T, total int) {
	var paginateRequest commonrequests.Paginate
	sanitize := SanitizeRequest(ctx, &paginateRequest)
	if sanitize != nil {
		return []T{}, 0
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)
	total = len(allItems)
	startIndex := (page - 1) * limit
	endIndex := page * limit

	if total == 0 {
		return []T{}, 0
	}
	if startIndex > total {
		return []T{}, total
	}
	if endIndex > total {
		endIndex = total
	}

	return allItems[startIndex:endIndex], total
}

// SanitizeRequest 消毒请求参数
func SanitizeRequest(ctx http.Context, request http.FormRequest) http.Response {
	errors, err := ctx.Request().ValidateRequest(request)
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if errors != nil {
		return Error(ctx, http.StatusUnprocessableEntity, errors.One())
	}

	return nil
}

// Sanitize 消毒参数
func Sanitize(ctx http.Context, rules map[string]string) http.Response {
	validator, err := ctx.Request().Validate(rules)
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	return nil
}
