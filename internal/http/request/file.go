package request

type FilePath struct {
	Path string `json:"path" form:"path"`
}

type FileCreate struct {
	Dir  bool   `json:"dir" form:"dir"`
	Path string `json:"path" form:"path"`
}

type FileSave struct {
	Path    string `form:"path" json:"path"`
	Content string `form:"content" json:"content"`
}

type FileMove struct {
	Source string `form:"source" json:"source"`
	Target string `form:"target" json:"target"`
	Force  bool   `form:"force" json:"force"`
}

type FileCopy struct {
	Source string `form:"source" json:"source"`
	Target string `form:"target" json:"target"`
	Force  bool   `form:"force" json:"force"`
}

type FilePermission struct {
	Path  string `form:"path" json:"path"`
	Mode  string `form:"mode" json:"mode"`
	Owner string `form:"owner" json:"owner"`
	Group string `form:"group" json:"group"`
}

type FileCompress struct {
	Paths []string `form:"paths" json:"paths"`
	File  string   `form:"file" json:"file"`
}

type FileUnCompress struct {
	File string `form:"file" json:"file"`
	Path string `form:"path" json:"path"`
}

type FileSearch struct {
	Path    string `form:"path" json:"path"`
	KeyWord string `form:"keyword" json:"keyword"`
}
