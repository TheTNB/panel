package service

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/http/request"
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
func Success(w http.ResponseWriter, data any) {
	render := chix.NewRender(w)
	defer render.Release()
	render.JSON(&SuccessResponse{
		Message: "success",
		Data:    data,
	})
}

// Error 响应错误
func Error(w http.ResponseWriter, code int, message string) {
	render := chix.NewRender(w)
	defer render.Release()
	render.Status(code)
	render.JSON(&ErrorResponse{
		Message: message,
	})
}

// ErrorSystem 响应系统错误
func ErrorSystem(w http.ResponseWriter) {
	render := chix.NewRender(w)
	defer render.Release()
	render.Status(http.StatusInternalServerError)
	render.JSON(&ErrorResponse{
		Message: "系统内部错误",
	})
}

// Bind 验证并绑定请求参数
func Bind[T any](r *http.Request) (*T, error) {
	req := new(T)

	// 绑定参数
	binder := chix.NewBind(r)
	defer binder.Release()
	if err := binder.URI(req); err != nil {
		return nil, err
	}
	if err := binder.Query(req); err != nil {
		return nil, err
	}
	if slices.Contains([]string{"POST", "PUT", "PATCH"}, strings.ToUpper(r.Method)) {
		if err := binder.Body(req); err != nil {
			return nil, err
		}
	}

	// 准备验证
	if reqWithPrepare, ok := any(req).(request.WithPrepare); ok {
		if err := reqWithPrepare.Prepare(r); err != nil {
			return nil, err
		}
	}
	if reqWithAuthorize, ok := any(req).(request.WithAuthorize); ok {
		if err := reqWithAuthorize.Authorize(r); err != nil {
			return nil, err
		}
	}
	if reqWithRules, ok := any(req).(request.WithRules); ok {
		if rules := reqWithRules.Rules(r); rules != nil {
			app.Validator.RegisterStructValidationMapRules(rules, req)
		}
	}

	// 验证参数
	err := app.Validator.Struct(req)
	if err == nil {
		return req, nil
	}

	// 翻译错误信息
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, e := range errs {
			if reqWithMessages, ok := any(req).(request.WithMessages); ok {
				if msg, found := reqWithMessages.Messages(r)[fmt.Sprintf("%s.%s", e.Field(), e.Tag())]; found {
					return nil, errors.New(msg)
				}
			}
			return nil, errors.New(e.Translate(*app.Translator))
		}
	}

	return nil, err
}

// Paginate 取分页条目
func Paginate[T any](r *http.Request, allItems []T) (pagedItems []T, total uint) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		req.Page = 1
		req.Limit = 10
	}
	total = uint(len(allItems))
	startIndex := (req.Page - 1) * req.Limit
	endIndex := req.Page * req.Limit

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

// removeTopStruct 移除验证器返回中的顶层结构
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
