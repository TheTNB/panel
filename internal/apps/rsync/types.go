package rsync

type Module struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Comment    string `json:"comment"`
	ReadOnly   bool   `json:"read_only"`
	AuthUser   string `json:"auth_user"`
	Secret     string `json:"secret"`
	HostsAllow string `json:"hosts_allow"`
}
