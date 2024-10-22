package request

type SSHCreate struct {
	Name       string `json:"name" form:"name"`
	Host       string `json:"host" form:"host" validate:"required"`
	Port       uint   `json:"port" form:"port" validate:"required,number,gte=1,lte=65535"`
	AuthMethod string `json:"auth_method" form:"auth_method" validate:"required,oneof=password publickey"`
	User       string `json:"user" form:"user" validate:"required_if=AuthMethod password"`
	Password   string `json:"password" form:"password" validate:"required_if=AuthMethod password"`
	Key        string `json:"key" form:"key" validate:"required_if=AuthMethod publickey"`
	Remark     string `json:"remark" form:"remark"`
}

type SSHUpdate struct {
	ID         uint   `form:"id" json:"id" validate:"required,exists=sshes id"`
	Name       string `json:"name" form:"name"`
	Host       string `json:"host" form:"host" validate:"required"`
	Port       uint   `json:"port" form:"port" validate:"required,number,gte=1,lte=65535"`
	AuthMethod string `json:"auth_method" form:"auth_method" validate:"required,oneof=password publickey"`
	User       string `json:"user" form:"user" validate:"required_if=AuthMethod password"`
	Password   string `json:"password" form:"password" validate:"required_if=AuthMethod password"`
	Key        string `json:"key" form:"key" validate:"required_if=AuthMethod publickey"`
	Remark     string `json:"remark" form:"remark"`
}
