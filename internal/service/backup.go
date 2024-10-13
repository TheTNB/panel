package service

import (
	"net/http"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
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

}

func (s *BackupService) Create(w http.ResponseWriter, r *http.Request) {

}

func (s *BackupService) Update(w http.ResponseWriter, r *http.Request) {

}

func (s *BackupService) Get(w http.ResponseWriter, r *http.Request) {

}

func (s *BackupService) Delete(w http.ResponseWriter, r *http.Request) {

}

func (s *BackupService) Restore(w http.ResponseWriter, r *http.Request) {

}
