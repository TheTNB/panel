package phpmyadmin

type UpdateConfig struct {
	Config string `form:"config" json:"config" validate:"required"`
}

type UpdatePort struct {
	Port uint `form:"port" json:"port" validate:"required,number,gte=1,lte=65535"`
}
