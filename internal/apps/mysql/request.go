package mysql

type UpdateConfig struct {
	Config string `form:"config" json:"config" validate:"required"`
}

type SetRootPassword struct {
	Password string `form:"password" json:"password" validate:"required|password"`
}
