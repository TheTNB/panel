package frp

type Name struct {
	Name string `form:"name" json:"name" validate:"required"`
}

type UpdateConfig struct {
	Name   string `form:"name" json:"name" validate:"required"`
	Config string `form:"config" json:"config" validate:"required"`
}
