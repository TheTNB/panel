package request

type SSHCreate struct {
	Name       string `json:"name" form:"name"`
	Host       string `json:"host" form:"host" validate:"required"`
	Port       uint   `json:"port" form:"port" validate:"required|min:1|max:65535"`
	AuthMethod string `json:"auth_method" form:"auth_method" validate:"required|in:password,publickey"`
	User       string `json:"user" form:"user" validate:"requiredIf:AuthMethod,password"`
	Password   string `json:"password" form:"password" validate:"requiredIf:AuthMethod,password"`
	Key        string `json:"key" form:"key" validate:"requiredIf:AuthMethod,publickey"`
	Remark     string `json:"remark" form:"remark"`
}

type SSHUpdate struct {
	ID         uint   `form:"id" json:"id" validate:"required|exists:sshes,id"`
	Name       string `json:"name" form:"name"`
	Host       string `json:"host" form:"host" validate:"required"`
	Port       uint   `json:"port" form:"port" validate:"required|min:1|max:65535"`
	AuthMethod string `json:"auth_method" form:"auth_method" validate:"required|in:password,publickey"`
	User       string `json:"user" form:"user" validate:"requiredIf:AuthMethod,password"`
	Password   string `json:"password" form:"password" validate:"requiredIf:AuthMethod,password"`
	Key        string `json:"key" form:"key" validate:"requiredIf:AuthMethod,publickey"`
	Remark     string `json:"remark" form:"remark"`
}
