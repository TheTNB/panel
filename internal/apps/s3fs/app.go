package s3fs

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-rat/chix"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
)

type App struct {
	settingRepo biz.SettingRepo
}

func NewApp(setting biz.SettingRepo) *App {
	return &App{
		settingRepo: setting,
	}
}

func (s *App) Route(r chi.Router) {
	r.Get("/mounts", s.List)
	r.Post("/mounts", s.Create)
	r.Delete("/mounts", s.Delete)
}

// List 所有 S3fs 挂载
func (s *App) List(w http.ResponseWriter, r *http.Request) {
	var s3fsList []Mount
	list, err := s.settingRepo.Get("s3fs", "[]")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to get s3fs list")
		return
	}

	if err = json.Unmarshal([]byte(list), &s3fsList); err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to get s3fs list")
		return
	}

	paged, total := service.Paginate(r, s3fsList)

	service.Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

// Create 添加 S3fs 挂载
func (s *App) Create(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Create](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	// 检查下地域节点中是否包含bucket，如果包含了，肯定是错误的
	if strings.Contains(req.URL, req.Bucket) {
		service.Error(w, http.StatusUnprocessableEntity, "endpoint should not contain bucket")
		return
	}

	// 检查挂载目录是否存在且为空
	if !io.Exists(req.Path) {
		if err = io.Mkdir(req.Path, 0755); err != nil {
			service.Error(w, http.StatusUnprocessableEntity, "failed to create mount path: %v", err)
			return
		}
	}
	if !io.Empty(req.Path) {
		service.Error(w, http.StatusUnprocessableEntity, "mount path is not empty")
		return
	}

	var s3fsList []Mount
	list, err := s.settingRepo.Get("s3fs", "[]")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to get s3fs list")
		return
	}
	if err = json.Unmarshal([]byte(list), &s3fsList); err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to get s3fs list")
		return
	}

	for _, s := range s3fsList {
		if s.Path == req.Path {
			service.Error(w, http.StatusUnprocessableEntity, "mount path already exists")
			return
		}
	}

	id := time.Now().UnixMicro()
	password := req.Ak + ":" + req.Sk
	if err = io.Write("/etc/passwd-s3fs-"+cast.ToString(id), password, 0600); err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to create passwd file: %v", err)
		return
	}
	if _, err = shell.Execf(`echo 's3fs#%s %s fuse _netdev,allow_other,nonempty,url=%s,passwd_file=/etc/passwd-s3fs-%s 0 0' >> /etc/fstab`, req.Bucket, req.Path, req.URL, cast.ToString(id)); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	if _, err = shell.Execf("mount -a"); err != nil {
		_, _ = shell.Execf(`sed -i 's@^s3fs#%s\s%s.*$@@g' /etc/fstab`, req.Bucket, req.Path)
		service.Error(w, http.StatusInternalServerError, "invalid /etc/fstab: %v", err)
		return
	}
	if _, err = shell.Execf(`df -h | grep '%s'`, req.Path); err != nil {
		_, _ = shell.Execf(`sed -i 's@^s3fs#%s\s%s.*$@@g' /etc/fstab`, req.Bucket, req.Path)
		service.Error(w, http.StatusInternalServerError, "mount failed: %v", err)
		return
	}

	s3fsList = append(s3fsList, Mount{
		ID:     id,
		Path:   req.Path,
		Bucket: req.Bucket,
		Url:    req.URL,
	})
	encoded, err := json.Marshal(s3fsList)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to add s3fs mount")
		return
	}

	if err = s.settingRepo.Set("s3fs", string(encoded)); err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to add s3fs mount")
		return
	}

	service.Success(w, nil)
}

// Delete 删除 S3fs 挂载
func (s *App) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[Delete](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	var s3fsList []Mount
	list, err := s.settingRepo.Get("s3fs", "[]")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to get s3fs list")
		return
	}
	if err = json.Unmarshal([]byte(list), &s3fsList); err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to get s3fs list")
		return
	}

	var mount Mount
	for _, item := range s3fsList {
		if item.ID == req.ID {
			mount = item
			break
		}
	}
	if mount.ID == 0 {
		service.Error(w, http.StatusUnprocessableEntity, "mount not found")
		return
	}

	if _, err = shell.Execf(`fusermount -uz '%s'`, mount.Path); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	if _, err = shell.Execf(`umount -lfq '%s'`, mount.Path); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	if _, err = shell.Execf(`sed -i 's@^s3fs#%s\s%s.*$@@g' /etc/fstab`, mount.Bucket, mount.Path); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}
	if _, err = shell.Execf("mount -a"); err != nil {
		service.Error(w, http.StatusInternalServerError, "invalid /etc/fstab: %v", err)
		return
	}
	if err = io.Remove("/etc/passwd-s3fs-" + cast.ToString(mount.ID)); err != nil {
		service.Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	var newS3fsList []Mount
	for _, item := range s3fsList {
		if item.ID != mount.ID {
			newS3fsList = append(newS3fsList, item)
		}
	}
	encoded, err := json.Marshal(newS3fsList)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to delete s3fs mount")
		return
	}
	if err = s.settingRepo.Set("s3fs", string(encoded)); err != nil {
		service.Error(w, http.StatusInternalServerError, "failed to delete s3fs mount")
		return
	}

	service.Success(w, nil)
}
