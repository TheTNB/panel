package service

import "net/http"

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

func (s *FileService) Create(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) Content(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) Save(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) Delete(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) Upload(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) Move(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) Copy(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) Download(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) RemoteDownload(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) Info(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) Permission(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) Archive(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) UnArchive(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) Search(w http.ResponseWriter, r *http.Request) {

}

func (s *FileService) List(w http.ResponseWriter, r *http.Request) {

}
