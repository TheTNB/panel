package supervisor

type UpdateConfig struct {
	Config string `form:"config" json:"config"`
}

type UpdateProcessConfig struct {
	Process string `form:"config" json:"process"`
	Config  string `form:"config" json:"config"`
}

type ProcessName struct {
	Process string `form:"config" json:"process"`
}

type CreateProcess struct {
	Name    string `form:"name" json:"name"`
	User    string `form:"user" json:"user"`
	Path    string `form:"path" json:"path"`
	Command string `form:"command" json:"command"`
	Num     int    `form:"num" json:"num"`
}
