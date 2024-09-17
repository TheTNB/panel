package request

type UserLogin struct {
	Username string `json:"username" form:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=255"`
}

type UserID struct {
	ID uint `uri:"id" validate:"required,number"`
}

type AddUser struct {
	Name string `json:"name" form:"name" validate:"required,min=3,max=255" comment:"用户名"`
}
