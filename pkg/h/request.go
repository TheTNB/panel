package h

import "github.com/goravel/framework/contracts/http"

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
