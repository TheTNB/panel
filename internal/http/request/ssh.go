package request

type SSHUpdateInfo struct {
	Host     string `json:"host" form:"host" validate:"required"`
	Port     int    `json:"port" form:"port" validate:"required,number,gte=1,lte=65535"`
	User     string `json:"user" form:"user" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}
