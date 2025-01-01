package service

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/go-rat/chix"
	"github.com/gookit/validate"

	"github.com/tnb-labs/panel/internal/http/request"
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
func Error(w http.ResponseWriter, code int, format string, args ...any) {
	render := chix.NewRender(w)
	defer render.Release()
	render.Header(chix.HeaderContentType, chix.MIMEApplicationJSONCharsetUTF8) // must before Status()
	render.Status(code)
	render.JSON(&ErrorResponse{
		Message: fmt.Sprintf(format, args...),
	})
}

// ErrorSystem 响应系统错误
func ErrorSystem(w http.ResponseWriter) {
	render := chix.NewRender(w)
	defer render.Release()
	render.Header(chix.HeaderContentType, chix.MIMEApplicationJSONCharsetUTF8) // must before Status()
	render.Status(http.StatusInternalServerError)
	render.JSON(&ErrorResponse{
		Message: http.StatusText(http.StatusInternalServerError),
	})
}

// Bind 验证并绑定请求参数
func Bind[T any](r *http.Request) (*T, error) {
	req := new(T)

	// 绑定参数
	binder := chix.NewBind(r)
	defer binder.Release()
	if slices.Contains([]string{"POST", "PUT", "PATCH", "DELETE"}, strings.ToUpper(r.Method)) {
		if r.ContentLength > 0 {
			if err := binder.Body(req); err != nil {
				return nil, err
			}
		}
	}
	if err := binder.Query(req); err != nil {
		return nil, err
	}
	if err := binder.URI(req); err != nil {
		return nil, err
	}

	// 准备验证
	df, err := validate.FromStruct(req)
	if err != nil {
		return nil, err
	}
	v := df.Create()

	if reqWithPrepare, ok := any(req).(request.WithPrepare); ok {
		if err = reqWithPrepare.Prepare(r); err != nil {
			return nil, err
		}
	}
	if reqWithAuthorize, ok := any(req).(request.WithAuthorize); ok {
		if err = reqWithAuthorize.Authorize(r); err != nil {
			return nil, err
		}
	}
	if reqWithRules, ok := any(req).(request.WithRules); ok {
		if rules := reqWithRules.Rules(r); rules != nil {
			for key, value := range rules {
				v.StringRule(key, value)
			}
		}
	}
	if reqWithFilters, ok := any(req).(request.WithFilters); ok {
		if filters := reqWithFilters.Filters(r); filters != nil {
			v.FilterRules(filters)
		}
	}
	if reqWithMessages, ok := any(req).(request.WithMessages); ok {
		if messages := reqWithMessages.Messages(r); messages != nil {
			v.AddMessages(messages)
		}
	}

	// 开始验证
	if v.Validate() && v.IsSuccess() {
		return req, nil
	}

	return nil, v.Errors.OneError()
}

// Paginate 取分页条目
func Paginate[T any](r *http.Request, items []T) (pagedItems []T, total uint) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		req = &request.Paginate{
			Page:  1,
			Limit: 10,
		}
	}
	total = uint(len(items))
	start := (req.Page - 1) * req.Limit
	end := req.Page * req.Limit

	if total == 0 {
		return []T{}, 0
	}
	if start > total {
		return []T{}, total
	}
	if end > total {
		end = total
	}

	return items[start:end], total
}
