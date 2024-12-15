//go:build !linux

// 这个文件只是为了在 Windows 下能编译通过，实际上并没有任何卵用

package service

import (
	"fmt"
	stdio "io"
	"net/http"
	stdos "os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-rat/chix"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/tools"
)

type FileService struct {
	taskRepo biz.TaskRepo
}

func NewFileService(task biz.TaskRepo) *FileService {
	return &FileService{
		taskRepo: task,
	}
}

func (s *FileService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FileCreate](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if !req.Dir {
		if _, err = shell.Execf("touch %s", req.Path); err != nil {
			Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
	} else {
		if err = io.Mkdir(req.Path, 0755); err != nil {
			Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
	}

	s.setPermission(req.Path, 0755, "www", "www")
	Success(w, nil)
}

func (s *FileService) Content(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FilePath](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	fileInfo, err := io.FileInfo(req.Path)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	if fileInfo.IsDir() {
		Error(w, http.StatusInternalServerError, "目标路径不是文件")
		return
	}
	if fileInfo.Size() > 10*1024*1024 {
		Error(w, http.StatusInternalServerError, "文件大小超过 10 M，不支持在线编辑")
		return
	}

	content, err := io.Read(req.Path)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, content)
}

func (s *FileService) Save(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FileSave](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	fileInfo, err := io.FileInfo(req.Path)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = io.Write(req.Path, req.Content, fileInfo.Mode()); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	s.setPermission(req.Path, 0755, "www", "www")
	Success(w, nil)
}

func (s *FileService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FilePath](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = io.Remove(req.Path); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *FileService) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 30); err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	path := r.FormValue("path")
	_, handler, err := r.FormFile("file")
	if err != nil {
		Error(w, http.StatusInternalServerError, "上传文件失败：%v", err)
		return
	}
	if io.Exists(path) {
		Error(w, http.StatusForbidden, "目标路径 %s 已存在", path)
		return
	}

	if !io.Exists(filepath.Dir(path)) {
		if err = io.Mkdir(filepath.Dir(path), 0755); err != nil {
			Error(w, http.StatusInternalServerError, "创建文件夹失败：%v", err)
			return
		}
	}

	src, _ := handler.Open()
	out, err := stdos.OpenFile(path, stdos.O_CREATE|stdos.O_RDWR|stdos.O_TRUNC, 0644)
	if err != nil {
		Error(w, http.StatusInternalServerError, "打开文件失败：%v", err)
		return
	}

	if _, err = stdio.Copy(out, src); err != nil {
		Error(w, http.StatusInternalServerError, "写入文件失败：%v", err)
		return
	}

	_ = src.Close()
	s.setPermission(path, 0755, "www", "www")
	Success(w, nil)
}

func (s *FileService) Exist(w http.ResponseWriter, r *http.Request) {
	binder := chix.NewBind(r)
	defer binder.Release()

	var paths []string
	if err := binder.Body(&paths); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	var results []bool
	for item := range slices.Values(paths) {
		results = append(results, io.Exists(item))
	}

	Success(w, results)
}

func (s *FileService) Move(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FileControl](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if io.Exists(req.Target) && !req.Force {
		Error(w, http.StatusForbidden, "目标路径 %s 已存在", req.Target)
		return
	}

	if err = io.Mv(req.Source, req.Target); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *FileService) Copy(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FileControl](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if io.Exists(req.Target) && !req.Force {
		Error(w, http.StatusForbidden, "目标路径 %s 已存在", req.Target)
		return
	}

	if err = io.Cp(req.Source, req.Target); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *FileService) Download(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FilePath](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	info, err := io.FileInfo(req.Path)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	if info.IsDir() {
		Error(w, http.StatusInternalServerError, "不能下载目录")
		return
	}

	render := chix.NewRender(w, r)
	defer render.Release()
	render.Download(req.Path, info.Name())
}

func (s *FileService) RemoteDownload(w http.ResponseWriter, r *http.Request) {
}

func (s *FileService) Info(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FilePath](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	info, err := io.FileInfo(req.Path)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, chix.M{
		"name":     info.Name(),
		"size":     tools.FormatBytes(float64(info.Size())),
		"mode_str": info.Mode().String(),
		"mode":     fmt.Sprintf("%04o", info.Mode().Perm()),
		"dir":      info.IsDir(),
		"modify":   info.ModTime().Format(time.DateTime),
	})
}

func (s *FileService) Permission(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FilePermission](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	// 解析成8进制
	mode, err := strconv.ParseUint(req.Mode, 8, 64)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = io.Chmod(req.Path, stdos.FileMode(mode)); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	if err = io.Chown(req.Path, req.Owner, req.Group); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *FileService) Compress(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FileCompress](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = io.Compress(filepath.Dir(req.Paths[0]), req.Paths, req.File); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	s.setPermission(req.File, 0755, "www", "www")
	Success(w, nil)
}

func (s *FileService) UnCompress(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FileUnCompress](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = io.UnCompress(req.File, req.Path); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	s.setPermission(req.Path, 0755, "www", "www")
	Success(w, nil)
}

func (s *FileService) Search(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FileSearch](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	results, err := io.SearchX(req.Path, req.Keyword, req.Sub)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	paged, total := Paginate(r, s.formatInfo(results))

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

func (s *FileService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FileList](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	list, err := io.ReadDir(req.Path)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if req.Sort == "asc" {
		slices.SortFunc(list, func(a, b stdos.DirEntry) int {
			return strings.Compare(strings.ToLower(b.Name()), strings.ToLower(a.Name()))
		})
	} else if req.Sort == "desc" {
		slices.SortFunc(list, func(a, b stdos.DirEntry) int {
			return strings.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
		})
	} else {
		slices.SortFunc(list, func(a, b stdos.DirEntry) int {
			if a.IsDir() && !b.IsDir() {
				return -1
			}
			if !a.IsDir() && b.IsDir() {
				return 1
			}
			return strings.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
		})
	}

	paged, total := Paginate(r, s.formatDir(req.Path, list))

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

// formatDir 格式化目录信息
func (s *FileService) formatDir(base string, entries []stdos.DirEntry) []any {
	var paths []any
	for _, file := range entries {
		info, _ := file.Info()

		paths = append(paths, map[string]any{
			"name":     info.Name(),
			"full":     filepath.Join(base, info.Name()),
			"size":     tools.FormatBytes(float64(info.Size())),
			"mode_str": info.Mode().String(),
			"mode":     fmt.Sprintf("%04o", info.Mode().Perm()),
			"owner":    "",
			"group":    "",
			"uid":      0,
			"gid":      0,
			"hidden":   io.IsHidden(info.Name()),
			"symlink":  io.IsSymlink(info.Mode()),
			"link":     io.GetSymlink(filepath.Join(base, info.Name())),
			"dir":      info.IsDir(),
			"modify":   info.ModTime().Format(time.DateTime),
		})
	}

	return paths
}

// formatInfo 格式化文件信息
func (s *FileService) formatInfo(infos map[string]stdos.FileInfo) []map[string]any {
	var paths []map[string]any
	for path, info := range infos {
		paths = append(paths, map[string]any{
			"name":     info.Name(),
			"full":     path,
			"size":     tools.FormatBytes(float64(info.Size())),
			"mode_str": info.Mode().String(),
			"mode":     fmt.Sprintf("%04o", info.Mode().Perm()),
			"owner":    "",
			"group":    "",
			"uid":      0,
			"gid":      0,
			"hidden":   io.IsHidden(info.Name()),
			"symlink":  io.IsSymlink(info.Mode()),
			"link":     io.GetSymlink(path),
			"dir":      info.IsDir(),
			"modify":   info.ModTime().Format(time.DateTime),
		})
	}

	slices.SortFunc(paths, func(a, b map[string]any) int {
		if cast.ToBool(a["dir"]) && !cast.ToBool(b["dir"]) {
			return -1
		}
		if !cast.ToBool(a["dir"]) && cast.ToBool(b["dir"]) {
			return 1
		}
		return strings.Compare(strings.ToLower(cast.ToString(a["name"])), strings.ToLower(cast.ToString(b["name"])))
	})

	return paths
}

// setPermission
func (s *FileService) setPermission(path string, mode stdos.FileMode, owner, group string) {
	_ = io.Chmod(path, mode)
	_ = io.Chown(path, owner, group)
}
