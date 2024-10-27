package request

import (
	"net/http"

	"github.com/spf13/cast"
)

type FileList struct {
	Path string `json:"path" form:"path" validate:"required"`
	Sort string `json:"sort" form:"sort"`
}

type FilePath struct {
	Path string `json:"path" form:"path" validate:"required"`
}

type FileCreate struct {
	Dir  bool   `json:"dir" form:"dir"`
	Path string `json:"path" form:"path" validate:"required"`
}

type FileSave struct {
	Path    string `form:"path" json:"path" validate:"required"`
	Content string `form:"content" json:"content"`
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

type FileRemoteDownload struct {
	Path string `form:"path" json:"path" validate:"required"`
	URL  string `form:"url" json:"url" validate:"required"`
}

type FilePermission struct {
	Path  string `form:"path" json:"path" validate:"required"`
	Mode  string `form:"mode" json:"mode" validate:"required"`
	Owner string `form:"owner" json:"owner" validate:"required"`
	Group string `form:"group" json:"group" validate:"required"`
}

type FileCompress struct {
	Dir   string   `form:"dir" json:"dir" validate:"required"`
	Paths []string `form:"paths" json:"paths" validate:"min=1,dive,required"`
	File  string   `form:"file" json:"file" validate:"required"`
}

type FileUnCompress struct {
	File string `form:"file" json:"file" validate:"required"`
	Path string `form:"path" json:"path" validate:"required"`
}

type FileSearch struct {
	Path    string `form:"path" json:"path" validate:"required"`
	Keyword string `form:"keyword" json:"keyword" validate:"required"`
	Sub     bool   `form:"sub" json:"sub"`
}

func (r *FileSearch) Prepare(req *http.Request) error {
	r.Sub = cast.ToBool(req.FormValue("sub"))
	return nil
}
