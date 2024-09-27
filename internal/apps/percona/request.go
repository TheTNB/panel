package percona

type UpdateConfig struct {
	Config string `form:"config" json:"config"`
}

type SetRootPassword struct {
	Password string `form:"password" json:"password"`
}
