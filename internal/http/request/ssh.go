package request

type SSHUpdateInfo struct {
	Host     string `json:"host" form:"host"`
	Port     int    `json:"port" form:"port"`
	User     string `json:"user" form:"user"`
	Password string `json:"password" form:"password"`
}
