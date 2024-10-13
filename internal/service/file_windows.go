//go:build !linux

// 这个文件只是为了在 Windows 下能编译通过，实际上并没有任何卵用

package service

import (
	"fmt"
	stdio "io"
	"net/http"
	stdos "os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/str"
)

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
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
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	if io.Exists(path) {
		Error(w, http.StatusForbidden, "目标路径%s已存在", path)
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

func (s *FileService) Move(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FileMove](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if io.Exists(req.Target) && !req.Force {
		Error(w, http.StatusForbidden, "目标路径"+req.Target+"已存在")
	}

	if err = io.Mv(req.Source, req.Target); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *FileService) Copy(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FileCopy](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if io.Exists(req.Target) && !req.Force {
		Error(w, http.StatusForbidden, "目标路径"+req.Target+"已存在")
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
		"size":     str.FormatBytes(float64(info.Size())),
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

	if err = io.Compress(req.Paths, req.File, io.Zip); err != nil {
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

	if err = io.UnCompress(req.File, req.Path, io.Zip); err != nil {
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

	paths := make(map[string]stdos.FileInfo)
	err = filepath.Walk(req.Path, func(path string, info stdos.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(info.Name(), req.KeyWord) {
			paths[path] = info
		}
		return nil
	})
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, paths)
}

func (s *FileService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.FilePath](r)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	fileInfoList, err := io.ReadDir(req.Path)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	var paths []any
	for _, fileInfo := range fileInfoList {
		info, _ := fileInfo.Info()

		paths = append(paths, map[string]any{
			"name":     info.Name(),
			"full":     filepath.Join(req.Path, info.Name()),
			"size":     str.FormatBytes(float64(info.Size())),
			"mode_str": info.Mode().String(),
			"mode":     fmt.Sprintf("%04o", info.Mode().Perm()),
			"owner":    "",
			"group":    "",
			"uid":      0,
			"gid":      0,
			"hidden":   io.IsHidden(info.Name()),
			"symlink":  io.IsSymlink(info.Mode()),
			"link":     io.GetSymlink(filepath.Join(req.Path, info.Name())),
			"dir":      info.IsDir(),
			"modify":   info.ModTime().Format(time.DateTime),
		})
	}

	paged, total := Paginate(r, paths)

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

// setPermission
func (s *FileService) setPermission(path string, mode stdos.FileMode, owner, group string) {
	_ = io.Chmod(path, mode)
	_ = io.Chown(path, owner, group)
}
