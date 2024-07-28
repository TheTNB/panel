package controllers

import (
	"fmt"
	stdio "io"
	stdos "os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"

	requests "github.com/TheTNB/panel/v2/app/http/requests/file"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/os"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/str"
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
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.NotExist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/create [post]
func (r *FileController) Create(ctx http.Context) http.Response {
	var request requests.NotExist
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	isDir := ctx.Request().InputBool("dir")
	if !isDir {
		if out, err := shell.Execf("touch " + request.Path); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, out)
		}
	} else {
		if err := io.Mkdir(request.Path, 0755); err != nil {
			return h.Error(ctx, http.StatusInternalServerError, err.Error())
		}
	}

	r.setPermission(request.Path, 0755, "www", "www")
	return h.Success(ctx, nil)
}

// Content
//
//	@Summary		获取文件内容
//	@Description	获取给定路径的文件内容
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.Exist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/content [get]
func (r *FileController) Content(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	fileInfo, err := io.FileInfo(request.Path)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if fileInfo.IsDir() {
		return h.Error(ctx, http.StatusInternalServerError, "目标路径不是文件")
	}
	if fileInfo.Size() > 10*1024*1024 {
		return h.Error(ctx, http.StatusInternalServerError, "文件大小超过 10 M，不支持在线编辑")
	}

	content, err := io.Read(request.Path)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, content)
}

// Save
//
//	@Summary		保存文件内容
//	@Description	保存给定路径的文件内容
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Save	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/save [post]
func (r *FileController) Save(ctx http.Context) http.Response {
	var request requests.Save
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	fileInfo, err := io.FileInfo(request.Path)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err = io.Write(request.Path, request.Content, fileInfo.Mode()); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.Path, 0755, "www", "www")
	return h.Success(ctx, nil)
}

// Delete
//
//	@Summary		删除文件/目录
//	@Description	删除给定路径的文件/目录
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Exist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/delete [post]
func (r *FileController) Delete(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	if err := io.Remove(request.Path); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// Upload
//
//	@Summary		上传文件
//	@Description	上传文件到给定路径
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			file	formData	file	true	"file"
//	@Param			path	formData	string	true	"path"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/upload [post]
func (r *FileController) Upload(ctx http.Context) http.Response {
	var request requests.Upload
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	src, err := request.File.Open()
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	if io.Exists(request.Path) && !ctx.Request().InputBool("force") {
		return h.Error(ctx, http.StatusForbidden, "目标路径已存在，是否覆盖？")
	}

	data, err := stdio.ReadAll(src)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err = io.Write(request.Path, string(data), 0755); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.Path, 0755, "www", "www")
	return h.Success(ctx, nil)
}

// Move
//
//	@Summary		移动文件/目录
//	@Description	移动文件/目录到给定路径，等效于重命名
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Move	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/move [post]
func (r *FileController) Move(ctx http.Context) http.Response {
	var request requests.Move
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	if io.Exists(request.Target) && !ctx.Request().InputBool("force") {
		return h.Error(ctx, http.StatusForbidden, "目标路径"+request.Target+"已存在")
	}

	if err := io.Mv(request.Source, request.Target); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.Target, 0755, "www", "www")
	return h.Success(ctx, nil)
}

// Copy
//
//	@Summary		复制文件/目录
//	@Description	复制文件/目录到给定路径
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Copy	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/copy [post]
func (r *FileController) Copy(ctx http.Context) http.Response {
	var request requests.Copy
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	if io.Exists(request.Target) && !ctx.Request().InputBool("force") {
		return h.Error(ctx, http.StatusForbidden, "目标路径"+request.Target+"已存在")
	}

	if err := io.Cp(request.Source, request.Target); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.Source, 0755, "www", "www")
	return h.Success(ctx, nil)
}

// Download
//
//	@Summary		下载文件
//	@Description	下载给定路径的文件
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.NotExist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/download [get]
func (r *FileController) Download(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	info, err := io.FileInfo(request.Path)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if info.IsDir() {
		return h.Error(ctx, http.StatusInternalServerError, "不能下载目录")
	}

	return ctx.Response().Download(request.Path, info.Name())
}

// RemoteDownload
//
//	@Summary		下载远程文件
//	@Description	下载远程文件到给定路径
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.NotExist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/remoteDownload [post]
func (r *FileController) RemoteDownload(ctx http.Context) http.Response {
	var request requests.NotExist
	sanitize := h.SanitizeRequest(ctx, &request)
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
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	query		requests.Exist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/info [get]
func (r *FileController) Info(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	fileInfo, err := io.FileInfo(request.Path)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, http.Json{
		"name":     fileInfo.Name(),
		"size":     str.FormatBytes(float64(fileInfo.Size())),
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
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Permission	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/permission [post]
func (r *FileController) Permission(ctx http.Context) http.Response {
	var request requests.Permission
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	// 解析成8进制
	mode, err := strconv.ParseUint(request.Mode, 8, 64)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err = io.Chmod(request.Path, stdos.FileMode(mode)); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err = io.Chown(request.Path, request.Owner, request.Group); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, nil)
}

// Archive
//
//	@Summary		压缩文件/目录
//	@Description	压缩文件/目录到给定路径
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.Archive	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/archive [post]
func (r *FileController) Archive(ctx http.Context) http.Response {
	var request requests.Archive
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	if err := io.Archive(request.Paths, request.File); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.File, 0755, "www", "www")
	return h.Success(ctx, nil)
}

// UnArchive
//
//	@Summary		解压文件/目录
//	@Description	解压文件/目录到给定路径
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			data	body		requests.UnArchive	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/unArchive [post]
func (r *FileController) UnArchive(ctx http.Context) http.Response {
	var request requests.UnArchive
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	if err := io.UnArchive(request.File, request.Path); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	r.setPermission(request.Path, 0755, "www", "www")
	return h.Success(ctx, nil)
}

// Search
//
//	@Summary		搜索文件/目录
//	@Description	通过关键词搜索给定路径的文件/目录
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Param			data	body		requests.Search	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/search [post]
func (r *FileController) Search(ctx http.Context) http.Response {
	var request requests.Search
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	paths := make(map[string]stdos.FileInfo)
	err := filepath.Walk(request.Path, func(path string, info stdos.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(info.Name(), request.KeyWord) {
			paths[path] = info
		}
		return nil
	})
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return h.Success(ctx, paths)
}

// List
//
//	@Summary		获取文件/目录列表
//	@Description	获取给定路径的文件/目录列表
//	@Tags			文件
//	@Accept			json
//	@Produce		json
//	@Param			data	query		requests.Exist			true	"request"
//	@Param			data	query		commonrequests.Paginate	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/list [get]
func (r *FileController) List(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := h.SanitizeRequest(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	fileInfoList, err := io.ReadDir(request.Path)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	var paths []any
	for _, fileInfo := range fileInfoList {
		info, _ := fileInfo.Info()
		stat := info.Sys().(*syscall.Stat_t)

		paths = append(paths, map[string]any{
			"name":     info.Name(),
			"full":     filepath.Join(request.Path, info.Name()),
			"size":     str.FormatBytes(float64(info.Size())),
			"mode_str": info.Mode().String(),
			"mode":     fmt.Sprintf("%04o", info.Mode().Perm()),
			"owner":    os.GetUser(stat.Uid),
			"group":    os.GetGroup(stat.Gid),
			"uid":      stat.Uid,
			"gid":      stat.Gid,
			"hidden":   io.IsHidden(info.Name()),
			"symlink":  io.IsSymlink(info.Mode()),
			"link":     io.GetSymlink(filepath.Join(request.Path, info.Name())),
			"dir":      info.IsDir(),
			"modify":   carbon.FromStdTime(info.ModTime()).ToDateTimeString(),
		})
	}

	paged, total := h.Paginate(ctx, paths)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// setPermission
func (r *FileController) setPermission(path string, mode stdos.FileMode, owner, group string) {
	_ = io.Chmod(path, mode)
	_ = io.Chown(path, owner, group)
}
