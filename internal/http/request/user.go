package request

type UserLogin struct {
	Username  string `json:"username" form:"username" validate:"required"`
	Password  string `json:"password" form:"password" validate:"required"`
	SafeLogin bool   `json:"safe_login" form:"safe_login"`
}
