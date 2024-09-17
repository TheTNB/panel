package service

import "net/http"

type BackupService struct {
}

func NewBackupService() *BackupService {
	return &BackupService{}
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
