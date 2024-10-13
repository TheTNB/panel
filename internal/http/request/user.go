package request

type UserLogin struct {
	Username string `json:"username" form:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=255"`
}
