package controllers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/http"

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

	if err := tools.Chmod(request.Path, 0755); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err := tools.Chown(request.Path, "www", "www"); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

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
//	@Param			data	body		requests.Exist	true	"request"
//	@Param			content	body		string			true	"content"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/save [post]
func (r *FileController) Save(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	content := ctx.Request().Input("content")

	fileInfo, err := tools.FileInfo(request.Path)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if err = tools.Write(request.Path, content, fileInfo.Mode()); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

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
//	@Param			data	body		requests.NotExist	true	"request"
//	@Param			file	formData	file				true	"file"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/upload [post]
func (r *FileController) Upload(ctx http.Context) http.Response {
	var request requests.NotExist
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	pathInfo, err := tools.FileInfo(request.Path)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}
	if !pathInfo.IsDir() {
		return Error(ctx, http.StatusInternalServerError, "目标路径不是目录")
	}

	file, err := ctx.Request().File("file")
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, "上传文件失败")
	}

	_, err = file.Store(request.Path)
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

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

	if tools.Exists(request.New) && !ctx.Request().InputBool("force") {
		return Error(ctx, http.StatusInternalServerError, "目标路径已存在，是否覆盖？")
	}

	if err := tools.Mv(request.Old, request.New); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

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
		return Error(ctx, http.StatusInternalServerError, "目标路径已存在，是否覆盖？")
	}

	if err := tools.Cp(request.Old, request.New); err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

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
	var request requests.NotExist
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

	return Success(ctx, fileInfo)
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
//	@Param			data	query		requests.Exist	true	"request"
//	@Success		200		{object}	SuccessResponse
//	@Router			/panel/file/list [get]
func (r *FileController) List(ctx http.Context) http.Response {
	var request requests.Exist
	sanitize := Sanitize(ctx, &request)
	if sanitize != nil {
		return sanitize
	}

	paths := make(map[string]os.FileInfo)
	err := filepath.Walk(request.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		paths[path] = info
		return nil
	})
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return Success(ctx, paths)
}
