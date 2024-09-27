package phpmyadmin

type UpdateConfig struct {
	Config string `form:"config" json:"config"`
}

type UpdatePort struct {
	Port uint `form:"port" json:"port"`
}
