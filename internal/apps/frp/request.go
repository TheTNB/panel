package frp

type Name struct {
	Name string `form:"name" json:"name"`
}

type UpdateConfig struct {
	Name   string `form:"name" json:"name"`
	Config string `form:"config" json:"config"`
}
