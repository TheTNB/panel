package pureftpd

type Create struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Path     string `form:"path" json:"path"`
}

type Delete struct {
	Username string `form:"username" json:"username"`
}

type ChangePassword struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type UpdatePort struct {
	Port uint `form:"port" json:"port"`
}
