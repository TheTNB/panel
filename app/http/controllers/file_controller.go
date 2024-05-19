package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
	commonrequests "panel/app/http/requests/common"

	requests "panel/app/http/requests/file"
	"panel/pkg/tools"
)

type FileController struct {
}

func NewFileController() *FileController {
	return &FileController{}
}

// Create
//
//	@Summary		创建文件/目录
//	@Description	创建文件/目录到给定路径
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.NotExist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/create [post]
func (r *FileController) Create(ctx http.Context) http.Response {
	var request requests.NotExist
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	isDir := ctx.Request().InputBool("dir")
	if !isDir {
		if out, err := tools.Exec("touch " + request.Path); err != nil {
			return Error(ctx, http.StatusInternalServerError, out)
		}
	} else {
		if err := tools.Mkdir(request.Path, 0755); err != nil {
			return Error(ctx, http.StatusInternalServerError, err.Error())
		}
	}

	r.setPermission(request.Path, 0755, "www", "www")
	return Success(ctx, nil)
}

// Content
//
//	@Summary		获取文件内容
//	@Description	获取给定路径的文件内容
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.Exist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/content [get]
func (r *FileController) Content(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	fileInfo, err := tools.FileInfo(request.Path)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if fileInfo.IsDir() {
		return Error(ctx, http.StatusInternalServerError, "目标路径不是文件")
	}
	if fileInfo.Size() > 10*1024*1024 {
		return Error(ctx, http.StatusInternalServerError, "文件大小超过 10 M，不支持在线编辑")
	}

	content, err := tools.Read(request.Path)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, content)
}

// Save
//
//	@Summary		保存文件内容
//	@Description	保存给定路径的文件内容
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Save	true	"request"
//	@Param			content	body		string			true	"content"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/save [post]
func (r *FileController) Save(ctx http.Context) http.Response {
	var request requests.Save
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	fileInfo, err := tools.FileInfo(request.Path)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err = tools.Write(request.Path, request.Content, fileInfo.Mode()); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.Path, 0755, "www", "www")
	return Success(ctx, nil)
}

// Delete
//
//	@Summary		删除文件/目录
//	@Description	删除给定路径的文件/目录
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Exist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/delete [post]
func (r *FileController) Delete(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.Remove(request.Path); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

// Upload
//
//	@Summary		上传文件
//	@Description	上传文件到给定路径
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			file	formData	file	true	"file"
//	@Param			path	formData	string	true	"path"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/upload [post]
func (r *FileController) Upload(ctx http.Context) http.Response {
	var request requests.Upload
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	src, err := request.File.Open()
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	if tools.Exists(request.Path) && !ctx.Request().InputBool("force") {
		return Error(ctx, http.StatusForbidden, "目标路径已存在，是否覆盖？")
	}

	data, err := io.ReadAll(src)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err = tools.Write(request.Path, string(data), 0755); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.Path, 0755, "www", "www")
	return Success(ctx, nil)
}

// Move
//
//	@Summary		移动文件/目录
//	@Description	移动文件/目录到给定路径，等效于重命名
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Move	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/move [post]
func (r *FileController) Move(ctx http.Context) http.Response {
	var request requests.Move
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	if tools.Exists(request.Target) && !ctx.Request().InputBool("force") {
		return Error(ctx, http.StatusForbidden, "目标路径已存在，是否覆盖？")
	}

	if err := tools.Mv(request.Source, request.Target); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.Target, 0755, "www", "www")
	return Success(ctx, nil)
}

// Copy
//
//	@Summary		复制文件/目录
//	@Description	复制文件/目录到给定路径
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Copy	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/copy [post]
func (r *FileController) Copy(ctx http.Context) http.Response {
	var request requests.Copy
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	if tools.Exists(request.New) && !ctx.Request().InputBool("force") {
		return Error(ctx, http.StatusForbidden, "目标路径已存在，是否覆盖？")
	}

	if err := tools.Cp(request.Old, request.New); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.New, 0755, "www", "www")
	return Success(ctx, nil)
}

// Download
//
//	@Summary		下载文件
//	@Description	下载给定路径的文件
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.NotExist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/download [get]
func (r *FileController) Download(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	info, err := tools.FileInfo(request.Path)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if info.IsDir() {
		return Error(ctx, http.StatusInternalServerError, "不能下载目录")
	}

	return ctx.Response().Download(request.Path, info.Name())
}

// RemoteDownload
//
//	@Summary		下载远程文件
//	@Description	下载远程文件到给定路径
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.NotExist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/remoteDownload [post]
func (r *FileController) RemoteDownload(ctx http.Context) http.Response {
	var request requests.NotExist
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	// TODO 使用异步任务下载文件
	return nil
}

// Info
//
//	@Summary		获取文件/目录信息
//	@Description	获取给定路径的文件/目录信息
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.Exist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/info [get]
func (r *FileController) Info(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	fileInfo, err := tools.FileInfo(request.Path)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, http.Json{
		"name":     fileInfo.Name(),
		"size":     tools.FormatBytes(float64(fileInfo.Size())),
		"mode_str": fileInfo.Mode().String(),
		"mode":     fmt.Sprintf("%04o", fileInfo.Mode().Perm()),
		"dir":      fileInfo.IsDir(),
		"modify":   carbon.FromStdTime(fileInfo.ModTime()).ToDateTimeString(),
	})
}

// Permission
//
//	@Summary		修改文件/目录权限
//	@Description	修改给定路径的文件/目录权限
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Permission	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/permission [post]
func (r *FileController) Permission(ctx http.Context) http.Response {
	var request requests.Permission
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.Chmod(request.Path, os.FileMode(request.Mode)); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err := tools.Chown(request.Path, request.Owner, request.Group); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, nil)
}

// Archive
//
//	@Summary		压缩文件/目录
//	@Description	压缩文件/目录到给定路径
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Archive	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/archive [post]
func (r *FileController) Archive(ctx http.Context) http.Response {
	var request requests.Archive
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.Archive(request.Paths, request.File); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.File, 0755, "www", "www")
	return Success(ctx, nil)
}

// UnArchive
//
//	@Summary		解压文件/目录
//	@Description	解压文件/目录到给定路径
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.UnArchive	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/unArchive [post]
func (r *FileController) UnArchive(ctx http.Context) http.Response {
	var request requests.UnArchive
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	if err := tools.UnArchive(request.File, request.Path); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.Path, 0755, "www", "www")
	return Success(ctx, nil)
}

// Search
//
//	@Summary		搜索文件/目录
//	@Description	通过关键词搜索给定路径的文件/目录
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Param			data	body		requests.Search	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/search [post]
func (r *FileController) Search(ctx http.Context) http.Response {
	var request requests.Search
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	paths := make(map[string]os.FileInfo)
	err := filepath.Walk(request.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(info.Name(), request.KeyWord) {
			paths[path] = info
		}
		return nil
	})
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, paths)
}

// List
//
//	@Summary		获取文件/目录列表
//	@Description	获取给定路径的文件/目录列表
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Param			data	query		requests.Exist			true	"request"
//	@Param			data	query		commonrequests.Paginate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/list [get]
func (r *FileController) List(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	var paginate commonrequests.Paginate
	paginateSanitize := Sanitize(ctx, &paginate)
	if paginateSanitize != nil {
		return paginateSanitize
	}

	fileInfoList, err := os.ReadDir(request.Path)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	var paths []any
	for _, fileInfo := range fileInfoList {
		info, _ := fileInfo.Info()
		stat := info.Sys().(*syscall.Stat_t)

		paths = append(paths, map[string]any{
			"name":     info.Name(),
			"full":     filepath.Join(request.Path, info.Name()),
			"size":     tools.FormatBytes(float64(info.Size())),
			"mode_str": info.Mode().String(),
			"mode":     fmt.Sprintf("%04o", info.Mode().Perm()),
			"owner":    tools.GetUser(stat.Uid),
			"group":    tools.GetGroup(stat.Gid),
			"uid":      stat.Uid,
			"gid":      stat.Gid,
			"hidden":   tools.IsHidden(info.Name()),
			"symlink":  tools.IsSymlink(info.Mode()),
			"link":     tools.GetSymlink(filepath.Join(request.Path, info.Name())),
			"dir":      info.IsDir(),
			"modify":   carbon.FromStdTime(info.ModTime()).ToDateTimeString(),
		})
	}

	start := paginate.Limit * (paginate.Page - 1)
	end := paginate.Limit * paginate.Page
	if start > len(paths) {
		start = len(paths)
	}
	if end > len(paths) {
		end = len(paths)
	}

	paged := paths[start:end]
	if paged == nil {
		paged = []any{}
	}

	return Success(ctx, http.Json{
		"total": len(paths),
		"items": paged,
	})
}

// setPermission
func (r *FileController) setPermission(path string, mode uint, owner, group string) {
	_ = tools.Chmod(path, os.FileMode(mode))
	_ = tools.Chown(path, owner, group)
}
