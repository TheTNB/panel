package plugins

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/json"
	"github.com/spf13/cast"

	"panel/app/http/controllers"
	"panel/app/services"
	"panel/pkg/tools"
)

type S3fsController struct {
	setting services.Setting
}

func NewS3fsController() *S3fsController {
	return &S3fsController{
		setting: services.NewSettingImpl(),
	}
}

// List 所有 S3fs 挂载
func (r *S3fsController) List(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "s3fs")
	if check != nil {
		return check
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)

	var s3fsList []S3fsMount
	err := json.UnmarshalString(r.setting.Get("s3fs", "[]"), &s3fsList)
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "获取 S3fs 挂载失败")
	}

	startIndex := (page - 1) * limit
	endIndex := page * limit
	if startIndex > len(s3fsList) {
		return controllers.Success(ctx, http.Json{
			"total": 0,
			"items": []S3fsMount{},
		})
	}
	if endIndex > len(s3fsList) {
		endIndex = len(s3fsList)
	}
	pagedS3fsList := s3fsList[startIndex:endIndex]
	if pagedS3fsList == nil {
		pagedS3fsList = []S3fsMount{}
	}

	return controllers.Success(ctx, http.Json{
		"total": len(s3fsList),
		"items": pagedS3fsList,
	})
}

// Add 添加 S3fs 挂载
func (r *S3fsController) Add(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "s3fs")
	if check != nil {
		return check
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"ak":     "required|regex:^[a-zA-Z0-9]*$",
		"sk":     "required|regex:^[a-zA-Z0-9]*$",
		"bucket": "required|regex:^[a-zA-Z0-9_-]*$",
		"url":    "required|full_url",
		"path":   "required|regex:^/[a-zA-Z0-9_-]+$",
	})
	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if validator.Fails() {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, validator.Errors().One())
	}

	ak := ctx.Request().Input("ak")
	sk := ctx.Request().Input("sk")
	path := ctx.Request().Input("path")
	bucket := ctx.Request().Input("bucket")
	url := ctx.Request().Input("url")

	// 检查下地域节点中是否包含bucket，如果包含了，肯定是错误的
	if strings.Contains(url, bucket) {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "地域节点不能包含 Bucket 名称")
	}

	// 检查挂载目录是否存在且为空
	if !tools.Exists(path) {
		tools.Mkdir(path, 0755)
	}
	if !tools.Empty(path) {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "挂载目录必须为空")
	}

	var s3fsList []S3fsMount
	err = json.UnmarshalString(r.setting.Get("s3fs", "[]"), &s3fsList)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 S3fs 挂载失败")
	}

	for _, s := range s3fsList {
		if s.Path == path {
			return controllers.Error(ctx, http.StatusUnprocessableEntity, "路径已存在")
		}
	}

	id := carbon.Now().TimestampMilli()
	password := ak + ":" + sk
	tools.Write("/etc/passwd-s3fs-"+cast.ToString(id), password, 0600)
	tools.Exec(`echo 's3fs#` + bucket + ` ` + path + ` fuse _netdev,allow_other,nonempty,url=` + url + `,passwd_file=/etc/passwd-s3fs-` + cast.ToString(id) + ` 0 0' >> /etc/fstab`)
	mountCheck := tools.Exec("mount -a 2>&1")
	if len(mountCheck) != 0 {
		tools.Exec(`sed -i 's@^s3fs#` + bucket + `\s` + path + `.*$@@g' /etc/fstab`)
		return controllers.Error(ctx, http.StatusInternalServerError, "检测到/etc/fstab有误: "+mountCheck)
	}
	dfCheck := tools.Exec("df -h | grep " + path + " 2>&1")
	if len(dfCheck) == 0 {
		tools.Exec(`sed -i 's@^s3fs#` + bucket + `\s` + path + `.*$@@g' /etc/fstab`)
		return controllers.Error(ctx, http.StatusInternalServerError, "挂载失败，请检查配置是否正确")
	}

	s3fsList = append(s3fsList, S3fsMount{
		ID:     id,
		Path:   path,
		Bucket: bucket,
		Url:    url,
	})
	encoded, err := json.MarshalString(s3fsList)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "添加 S3fs 挂载失败")
	}
	err = r.setting.Set("s3fs", encoded)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "添加 S3fs 挂载失败")
	}

	return controllers.Success(ctx, nil)
}

// Delete 删除 S3fs 挂载
func (r *S3fsController) Delete(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "s3fs")
	if check != nil {
		return check
	}

	id := ctx.Request().Input("id")
	if len(id) == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "挂载ID不能为空")
	}

	var s3fsList []S3fsMount
	err := json.UnmarshalString(r.setting.Get("s3fs", "[]"), &s3fsList)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "获取 S3fs 挂载失败")
	}

	var mount S3fsMount
	for _, s := range s3fsList {
		if cast.ToString(s.ID) == id {
			mount = s
			break
		}
	}
	if mount.ID == 0 {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, "挂载ID不存在")
	}

	tools.Exec(`fusermount -u '` + mount.Path + `' 2>&1`)
	tools.Exec(`umount '` + mount.Path + `' 2>&1`)
	tools.Exec(`sed -i 's@^s3fs#` + mount.Bucket + `\s` + mount.Path + `.*$@@g' /etc/fstab`)
	mountCheck := tools.Exec("mount -a 2>&1")
	if len(mountCheck) != 0 {
		return controllers.Error(ctx, http.StatusInternalServerError, "检测到/etc/fstab有误: "+mountCheck)
	}
	tools.Remove("/etc/passwd-s3fs-" + cast.ToString(mount.ID))

	var newS3fsList []S3fsMount
	for _, s := range s3fsList {
		if s.ID != mount.ID {
			newS3fsList = append(newS3fsList, s)
		}
	}
	encoded, err := json.MarshalString(newS3fsList)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "删除 S3fs 挂载失败")
	}
	err = r.setting.Set("s3fs", encoded)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "删除 S3fs 挂载失败")
	}

	return controllers.Success(ctx, nil)
}
