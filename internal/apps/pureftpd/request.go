package pureftpd

type Create struct {
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required|password"`
	Path     string `form:"path" json:"path" validate:"required"`
}

type Delete struct {
	Username string `form:"username" json:"username" validate:"required"`
}

type ChangePassword struct {
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required|password"`
}

type UpdatePort struct {
	Port uint `form:"port" json:"port" validate:"required|number|min:1|max:65535"`
}
