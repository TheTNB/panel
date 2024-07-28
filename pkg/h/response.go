package h

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/app/http/requests/common"
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
		Message: message,
	})
}

// ErrorSystem 响应系统错误
func ErrorSystem(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusInternalServerError, &ErrorResponse{
		Message: facades.Lang(ctx).Get("errors.internal"),
	})
}

// Paginate 取分页条目
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
