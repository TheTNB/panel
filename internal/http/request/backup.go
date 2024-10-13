package request

type BackupList struct {
	Type string `json:"type" form:"type" validate:"required,oneof=path website mysql postgres redis panel"`
}

type BackupCreate struct {
	Type   string `json:"type" form:"type" validate:"required,oneof=website mysql postgres redis panel"`
	Target string `json:"target" form:"target" validate:"required"`
	Path   string `json:"path" form:"path"`
}

type BackupFile struct {
	Type string `json:"type" form:"type" validate:"required,oneof=website mysql postgres redis panel"`
	File string `json:"file" form:"file" validate:"required"`
}

type BackupRestore struct {
	Type   string `json:"type" form:"type" validate:"required,oneof=website mysql postgres redis panel"`
	File   string `json:"file" form:"file" validate:"required"`
	Target string `json:"target" form:"target" validate:"required"`
}
