package supervisor

type UpdateConfig struct {
	Config string `form:"config" json:"config" validate:"required"`
}

type UpdateProcessConfig struct {
	Process string `form:"process" json:"process" validate:"required"`
	Config  string `form:"config" json:"config" validate:"required"`
}

type ProcessName struct {
	Process string `form:"process" json:"process" validate:"required"`
}

type CreateProcess struct {
	Name    string `form:"name" json:"name" validate:"required"`
	User    string `form:"user" json:"user" validate:"required"`
	Path    string `form:"path" json:"path" validate:"required"`
	Command string `form:"command" json:"command" validate:"required"`
	Num     int    `form:"num" json:"num" validate:"required|min:1"`
}
