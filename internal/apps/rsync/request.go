package rsync

type Create struct {
	Name       string `form:"name" json:"name" validate:"required"`
	Path       string `form:"path" json:"path" validate:"required"`
	Comment    string `form:"comment" json:"comment"`
	AuthUser   string `form:"auth_user" json:"auth_user" validate:"required"`
	Secret     string `form:"secret" json:"secret" validate:"required"`
	HostsAllow string `form:"hosts_allow" json:"hosts_allow"`
}

type Delete struct {
	Name string `form:"name" json:"name" validate:"required"`
}

type Update struct {
	Name       string `form:"name" json:"name" validate:"required"`
	Path       string `form:"path" json:"path" validate:"required"`
	Comment    string `form:"comment" json:"comment"`
	AuthUser   string `form:"auth_user" json:"auth_user" validate:"required"`
	Secret     string `form:"secret" json:"secret" validate:"required"`
	HostsAllow string `form:"hosts_allow" json:"hosts_allow"`
}

type UpdateConfig struct {
	Config string `form:"config" json:"config" validate:"required"`
}
