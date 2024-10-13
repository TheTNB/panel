package request

type FilePath struct {
	Path string `json:"path" form:"path" validate:"required"`
}

type FileCreate struct {
	Dir  bool   `json:"dir" form:"dir"`
	Path string `json:"path" form:"path" validate:"required"`
}

type FileSave struct {
	Path    string `form:"path" json:"path" validate:"required"`
	Content string `form:"content" json:"content" validate:"required"`
}

type FileMove struct {
	Source string `form:"source" json:"source" validate:"required"`
	Target string `form:"target" json:"target" validate:"required"`
	Force  bool   `form:"force" json:"force"`
}

type FileCopy struct {
	Source string `form:"source" json:"source" validate:"required"`
	Target string `form:"target" json:"target" validate:"required"`
	Force  bool   `form:"force" json:"force"`
}

type FilePermission struct {
	Path  string `form:"path" json:"path" validate:"required"`
	Mode  string `form:"mode" json:"mode" validate:"required"`
	Owner string `form:"owner" json:"owner" validate:"required"`
	Group string `form:"group" json:"group" validate:"required"`
}

type FileCompress struct {
	Paths []string `form:"paths" json:"paths" validate:"required"`
	File  string   `form:"file" json:"file" validate:"required"`
}

type FileUnCompress struct {
	File string `form:"file" json:"file" validate:"required"`
	Path string `form:"path" json:"path" validate:"required"`
}

type FileSearch struct {
	Path    string `form:"path" json:"path" validate:"required"`
	KeyWord string `form:"keyword" json:"keyword" validate:"required"`
}
