package s3fs

import (
	"strings"

	"github.com/bytedance/sonic"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
	"github.com/spf13/cast"

	"panel/app/http/controllers"
	"panel/app/services"
	"panel/pkg/tools"
)

type S3fsController struct {
	setting services.Setting
}

type s3fs struct {
	ID     int64  `json:"id"`
	Path   string `json:"path"`
	Bucket string `json:"bucket"`
	Url    string `json:"url"`
}

func NewS3fsController() *S3fsController {
	return &S3fsController{
		setting: services.NewSettingImpl(),
	}
}

// List 所有 S3fs 挂载
func (c *S3fsController) List(ctx http.Context) {
	if !controllers.Check(ctx, "s3fs") {
		return
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)

	var s3fsList []s3fs
	err := sonic.UnmarshalString(c.setting.Get("s3fs", "[]"), &s3fsList)
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "获取 S3fs 挂载失败")
		return
	}

	startIndex := (page - 1) * limit
	endIndex := page * limit
	if startIndex > len(s3fsList) {
		controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []s3fs{},
		})
		return
	}
	if endIndex > len(s3fsList) {
		endIndex = len(s3fsList)
	}
	pagedS3fsList := s3fsList[startIndex:endIndex]

	controllers.Success(ctx, http.Json{
		"total": len(s3fsList),
		"items": pagedS3fsList,
	})
}

// Add 添加 S3fs 挂载
func (c *S3fsController) Add(ctx http.Context) {
	if !controllers.Check(ctx, "s3fs") {
		return
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"ak":     "required|regex:^[a-zA-Z0-9]*$",
		"sk":     "required|regex:^[a-zA-Z0-9]*$",
		"bucket": "required|regex:^[a-zA-Z0-9_-]*$",
		"url":    "required|full_url",
		"path":   "required|regex:^/[a-zA-Z0-9_-]+$",
	})
	if err != nil {
		controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if validator.Fails() {
		controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
		return
	}

	ak := ctx.Request().Input("ak")
	sk := ctx.Request().Input("sk")
	path := ctx.Request().Input("path")
	bucket := ctx.Request().Input("bucket")
	url := ctx.Request().Input("url")

	// 检查下地域节点中是否包含bucket，如果包含了，肯定是错误的
	if strings.Contains(url, bucket) {
		controllers.Error(ctx, http.StatusUnprocessableEntity, "地域节点不能包含 Bucket 名称")
		return
	}

	// 检查挂载目录是否存在且为空
	if !tools.Exists(path) {
		tools.Mkdir(path, 0755)
	}
	if !tools.Empty(path) {
		controllers.Error(ctx, http.StatusUnprocessableEntity, "挂载目录必须为空")
		return
	}

	var s3fsList []s3fs
	err = sonic.UnmarshalString(c.setting.Get("s3fs", "[]"), &s3fsList)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, "获取 S3fs 挂载失败")
		return
	}

	for _, s := range s3fsList {
		if s.Path == path {
			controllers.Error(ctx, http.StatusUnprocessableEntity, "路径已存在")
			return
		}
	}

	id := carbon.Now().TimestampMilli()
	password := ak + ":" + sk
	tools.WriteFile("/etc/passwd-s3fs-"+cast.ToString(id), password, 0600)
	tools.ExecShell(`echo 's3fs#` + bucket + ` ` + path + ` fuse _netdev,allow_other,nonempty,url=` + url + `,passwd_file=/etc/passwd-s3fs-` + cast.ToString(id) + ` 0 0' >> /etc/fstab`)
	check := tools.ExecShell("mount -a 2>&1")
	if len(check) != 0 {
		tools.ExecShell(`sed -i 's@^s3fs#` + bucket + `\s` + path + `.*$@@g' /etc/fstab`)
		controllers.Error(ctx, http.StatusInternalServerError, "检测到/etc/fstab有误: "+check)
		return
	}
	check2 := tools.ExecShell("df -h | grep " + path + " 2>&1")
	if len(check2) == 0 {
		tools.ExecShell(`sed -i 's@^s3fs#` + bucket + `\s` + path + `.*$@@g' /etc/fstab`)
		controllers.Error(ctx, http.StatusInternalServerError, "挂载失败，请检查配置是否正确")
		return
	}

	s3fsList = append(s3fsList, s3fs{
		ID:     id,
		Path:   path,
		Bucket: bucket,
		Url:    url,
	})
	json, err := sonic.MarshalString(s3fsList)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, "添加 S3fs 挂载失败")
		return
	}
	err = c.setting.Set("s3fs", json)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, "添加 S3fs 挂载失败")
		return
	}

	controllers.Success(ctx, nil)
}

// Delete 删除 S3fs 挂载
func (c *S3fsController) Delete(ctx http.Context) {
	if !controllers.Check(ctx, "s3fs") {
		return
	}

	id := ctx.Request().Input("id")
	if len(id) == 0 {
		controllers.Error(ctx, http.StatusUnprocessableEntity, "挂载ID不能为空")
		return
	}

	var s3fsList []s3fs
	err := sonic.UnmarshalString(c.setting.Get("s3fs", "[]"), &s3fsList)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, "获取 S3fs 挂载失败")
		return
	}

	var mount s3fs
	for _, s := range s3fsList {
		if cast.ToString(s.ID) == id {
			mount = s
			break
		}
	}
	if mount.ID == 0 {
		controllers.Error(ctx, http.StatusUnprocessableEntity, "挂载ID不存在")
		return
	}

	tools.ExecShell(`fusermount -u '` + mount.Path + `'`)
	tools.ExecShell(`umount '` + mount.Path + `'`)
	tools.ExecShell(`sed -i 's@^s3fs#` + mount.Bucket + `\s` + mount.Path + `.*$@@g' /etc/fstab`)
	check := tools.ExecShell("mount -a 2>&1")
	if len(check) != 0 {
		controllers.Error(ctx, http.StatusInternalServerError, "检测到/etc/fstab有误: "+check)
		return
	}
	tools.RemoveFile("/etc/passwd-s3fs-" + cast.ToString(mount.ID))

	var newS3fsList []s3fs
	for _, s := range s3fsList {
		if s.ID != mount.ID {
			newS3fsList = append(newS3fsList, s)
		}
	}
	json, err := sonic.MarshalString(newS3fsList)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, "删除 S3fs 挂载失败")
		return
	}
	err = c.setting.Set("s3fs", json)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, "删除 S3fs 挂载失败")
		return
	}

	controllers.Success(ctx, nil)
}
