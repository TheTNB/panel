package request

import (
	"net/http"

	"github.com/spf13/cast"
)

type FileList struct {
	Path string `json:"path" form:"path" validate:"required|isUnixPath"`
	Sort string `json:"sort" form:"sort"`
}

type FilePath struct {
	Path string `json:"path" form:"path" validate:"required|isUnixPath"`
}

type FileCreate struct {
	Dir  bool   `json:"dir" form:"dir"`
	Path string `json:"path" form:"path" validate:"required|isUnixPath"`
}

type FileSave struct {
	Path    string `form:"path" json:"path" validate:"required|isUnixPath"`
	Content string `form:"content" json:"content"`
}

type FileControl struct {
	Source string `form:"source" json:"source" validate:"required|isUnixPath"`
	Target string `form:"target" json:"target" validate:"required|isUnixPath"`
	Force  bool   `form:"force" json:"force"`
}

type FileRemoteDownload struct {
	Path string `form:"path" json:"path" validate:"required|isUnixPath"`
	URL  string `form:"url" json:"url" validate:"required|isFullURL"`
}

type FilePermission struct {
	Path  string `form:"path" json:"path" validate:"required|isUnixPath"`
	Mode  string `form:"mode" json:"mode" validate:"required"`
	Owner string `form:"owner" json:"owner" validate:"required"`
	Group string `form:"group" json:"group" validate:"required"`
}

type FileCompress struct {
	Dir   string   `form:"dir" json:"dir" validate:"required|isUnixPath"`
	Paths []string `form:"paths" json:"paths" validate:"required|isSlice"`
	File  string   `form:"file" json:"file" validate:"required|isUnixPath"`
}

type FileUnCompress struct {
	File string `form:"file" json:"file" validate:"required|isUnixPath"`
	Path string `form:"path" json:"path" validate:"required|isUnixPath"`
}

type FileSearch struct {
	Path    string `form:"path" json:"path" validate:"required|isUnixPath"`
	Keyword string `form:"keyword" json:"keyword" validate:"required"`
	Sub     bool   `form:"sub" json:"sub"`
}

func (r *FileSearch) Prepare(req *http.Request) error {
	r.Sub = cast.ToBool(req.FormValue("sub"))
	return nil
}
