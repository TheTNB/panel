package rule

import (
	"unicode"

	"github.com/spf13/cast"
)

// Password 密码复杂度校验
type Password struct{}

func NewPassword() *Password {
	return &Password{}
}

func (r *Password) Passes(val any, options ...any) bool {
	password := cast.ToString(val)
	// 不对空密码进行校验，有需要可以使用 required 标签
	if password == "" {
		return true
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	if len(password) < 8 || len(password) > 20 {
		return false
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// 至少包含两类字符组合
	valid := (hasUpper && hasLower) ||
		(hasUpper && hasNumber) ||
		(hasUpper && hasSpecial) ||
		(hasLower && hasNumber) ||
		(hasLower && hasSpecial) ||
		(hasNumber && hasSpecial)

	return valid
}
