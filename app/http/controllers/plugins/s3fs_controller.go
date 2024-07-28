package plugins

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/json"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/v2/internal"
	"github.com/TheTNB/panel/v2/internal/services"
	"github.com/TheTNB/panel/v2/pkg/h"
	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type S3fsController struct {
	setting internal.Setting
}

func NewS3fsController() *S3fsController {
	return &S3fsController{
		setting: services.NewSettingImpl(),
	}
}

// List 所有 S3fs 挂载
func (r *S3fsController) List(ctx http.Context) http.Response {
	var s3fsList []types.S3fsMount
	err := json.UnmarshalString(r.setting.Get("s3fs", "[]"), &s3fsList)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取 S3fs 挂载失败")
	}

	paged, total := h.Paginate(ctx, s3fsList)

	return h.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// Add 添加 S3fs 挂载
func (r *S3fsController) Add(ctx http.Context) http.Response {
	if sanitize := h.Sanitize(ctx, map[string]string{
		"ak":     "required|regex:^[a-zA-Z0-9]*$",
		"sk":     "required|regex:^[a-zA-Z0-9]*$",
		"bucket": "required|regex:^[a-zA-Z0-9_-]*$",
		"url":    "required|full_url",
		"path":   "required|regex:^/[a-zA-Z0-9_-]+$",
	}); sanitize != nil {
		return sanitize
	}

	ak := ctx.Request().Input("ak")
	sk := ctx.Request().Input("sk")
	path := ctx.Request().Input("path")
	bucket := ctx.Request().Input("bucket")
	url := ctx.Request().Input("url")

	// 检查下地域节点中是否包含bucket，如果包含了，肯定是错误的
	if strings.Contains(url, bucket) {
		return h.Error(ctx, http.StatusUnprocessableEntity, "地域节点不能包含 Bucket 名称")
	}

	// 检查挂载目录是否存在且为空
	if !io.Exists(path) {
		if err := io.Mkdir(path, 0755); err != nil {
			return h.Error(ctx, http.StatusUnprocessableEntity, "挂载目录创建失败")
		}
	}
	if !io.Empty(path) {
		return h.Error(ctx, http.StatusUnprocessableEntity, "挂载目录必须为空")
	}

	var s3fsList []types.S3fsMount
	if err := json.UnmarshalString(r.setting.Get("s3fs", "[]"), &s3fsList); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取 S3fs 挂载失败")
	}

	for _, s := range s3fsList {
		if s.Path == path {
			return h.Error(ctx, http.StatusUnprocessableEntity, "路径已存在")
		}
	}

	id := carbon.Now().TimestampMilli()
	password := ak + ":" + sk
	if err := io.Write("/etc/passwd-s3fs-"+cast.ToString(id), password, 0600); err != nil {
		return nil
	}
	out, err := shell.Execf(`echo 's3fs#` + bucket + ` ` + path + ` fuse _netdev,allow_other,nonempty,url=` + url + `,passwd_file=/etc/passwd-s3fs-` + cast.ToString(id) + ` 0 0' >> /etc/fstab`)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if mountCheck, err := shell.Execf("mount -a 2>&1"); err != nil {
		_, _ = shell.Execf(`sed -i 's@^s3fs#` + bucket + `\s` + path + `.*$@@g' /etc/fstab`)
		return h.Error(ctx, http.StatusInternalServerError, "检测到/etc/fstab有误: "+mountCheck)
	}
	if _, err := shell.Execf("df -h | grep " + path + " 2>&1"); err != nil {
		_, _ = shell.Execf(`sed -i 's@^s3fs#` + bucket + `\s` + path + `.*$@@g' /etc/fstab`)
		return h.Error(ctx, http.StatusInternalServerError, "挂载失败，请检查配置是否正确")
	}

	s3fsList = append(s3fsList, types.S3fsMount{
		ID:     id,
		Path:   path,
		Bucket: bucket,
		Url:    url,
	})
	encoded, err := json.MarshalString(s3fsList)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "添加 S3fs 挂载失败")
	}
	err = r.setting.Set("s3fs", encoded)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "添加 S3fs 挂载失败")
	}

	return h.Success(ctx, nil)
}

// Delete 删除 S3fs 挂载
func (r *S3fsController) Delete(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	if len(id) == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "挂载ID不能为空")
	}

	var s3fsList []types.S3fsMount
	err := json.UnmarshalString(r.setting.Get("s3fs", "[]"), &s3fsList)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "获取 S3fs 挂载失败")
	}

	var mount types.S3fsMount
	for _, s := range s3fsList {
		if cast.ToString(s.ID) == id {
			mount = s
			break
		}
	}
	if mount.ID == 0 {
		return h.Error(ctx, http.StatusUnprocessableEntity, "挂载ID不存在")
	}

	if out, err := shell.Execf(`fusermount -u '` + mount.Path + `' 2>&1`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf(`umount '` + mount.Path + `' 2>&1`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if out, err := shell.Execf(`sed -i 's@^s3fs#` + mount.Bucket + `\s` + mount.Path + `.*$@@g' /etc/fstab`); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, out)
	}
	if mountCheck, err := shell.Execf("mount -a 2>&1"); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "检测到/etc/fstab有误: "+mountCheck)
	}
	if err := io.Remove("/etc/passwd-s3fs-" + cast.ToString(mount.ID)); err != nil {
		return h.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	var newS3fsList []types.S3fsMount
	for _, s := range s3fsList {
		if s.ID != mount.ID {
			newS3fsList = append(newS3fsList, s)
		}
	}
	encoded, err := json.MarshalString(newS3fsList)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "删除 S3fs 挂载失败")
	}
	err = r.setting.Set("s3fs", encoded)
	if err != nil {
		return h.Error(ctx, http.StatusInternalServerError, "删除 S3fs 挂载失败")
	}

	return h.Success(ctx, nil)
}
