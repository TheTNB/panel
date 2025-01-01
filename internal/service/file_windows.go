//go:build !linux

// 这个文件只是为了在 Windows 下能编译通过，实际上并没有任何卵用

package service

import (
	"net/http"

	"github.com/tnb-labs/panel/internal/biz"
)

type FileService struct{}

func NewFileService(_ biz.TaskRepo) *FileService {
	return &FileService{}
}

func (s *FileService) Create(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) Content(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) Save(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) Delete(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) Upload(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) Exist(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) Move(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) Copy(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) Download(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) RemoteDownload(w http.ResponseWriter, r *http.Request) {
}

func (s *FileService) Info(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) Permission(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) Compress(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) UnCompress(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) Search(w http.ResponseWriter, r *http.Request) {}

func (s *FileService) List(w http.ResponseWriter, r *http.Request) {}
