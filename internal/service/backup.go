package service

import (
	stdio "io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/io"
)

type BackupService struct {
	backupRepo biz.BackupRepo
}

func NewBackupService() *BackupService {
	return &BackupService{
		backupRepo: data.NewBackupRepo(),
	}
}

func (s *BackupService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.BackupList](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	list, _ := s.backupRepo.List(biz.BackupType(req.Type))
	paged, total := Paginate(r, list)

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

func (s *BackupService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.BackupCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.backupRepo.Create(biz.BackupType(req.Type), req.Target, req.Path); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *BackupService) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 30); err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	_, handler, err := r.FormFile("file")
	if err != nil {
		Error(w, http.StatusInternalServerError, "上传文件失败：%v", err)
		return
	}
	path, err := s.backupRepo.GetPath(biz.BackupType(r.FormValue("type")))
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if !io.Exists(filepath.Dir(path)) {
		if err = io.Mkdir(filepath.Dir(path), 0755); err != nil {
			Error(w, http.StatusInternalServerError, "创建文件夹失败：%v", err)
			return
		}
	}

	src, _ := handler.Open()
	out, err := os.OpenFile(filepath.Join(path, handler.Filename), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		Error(w, http.StatusInternalServerError, "打开文件失败：%v", err)
		return
	}

	if _, err = stdio.Copy(out, src); err != nil {
		Error(w, http.StatusInternalServerError, "写入文件失败：%v", err)
		return
	}

	_ = src.Close()
	Success(w, nil)
}

func (s *BackupService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.BackupFile](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.backupRepo.Delete(biz.BackupType(req.Type), req.File); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *BackupService) Restore(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.BackupRestore](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.backupRepo.Restore(biz.BackupType(req.Type), req.File, req.Target); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}
