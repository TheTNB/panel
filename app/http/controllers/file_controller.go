package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type FileController struct {
	// Dependent services
}

func NewFileController() *FileController {
	return &FileController{
		// Inject services
	}
}

func (r *FileController) Index(ctx http.Context) {
}
