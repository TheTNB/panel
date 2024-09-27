package rsync

type Create struct {
	Name       string `form:"name" json:"name"`
	Path       string `form:"path" json:"path"`
	Comment    string `form:"comment" json:"comment"`
	AuthUser   string `form:"auth_user" json:"auth_user"`
	Secret     string `form:"secret" json:"secret"`
	HostsAllow string `form:"hosts_allow" json:"hosts_allow"`
}

type Delete struct {
	Name string `form:"name" json:"name"`
}

type Update struct {
	Name       string `form:"name" json:"name"`
	Path       string `form:"path" json:"path"`
	Comment    string `form:"comment" json:"comment"`
	AuthUser   string `form:"auth_user" json:"auth_user"`
	Secret     string `form:"secret" json:"secret"`
	HostsAllow string `form:"hosts_allow" json:"hosts_allow"`
}

type UpdateConfig struct {
	Config string `form:"config" json:"config"`
}
