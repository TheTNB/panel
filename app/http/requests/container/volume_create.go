package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type VolumeCreate struct {
	Name    string   `form:"name" json:"name"`
	Driver  string   `form:"driver" json:"driver"`
	Labels  []string `form:"labels" json:"labels"`
	Options []string `form:"options" json:"options"`
}

func (r *VolumeCreate) Authorize(ctx http.Context) error {
	return nil
}

func (r *VolumeCreate) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":    "required|string",
		"driver":  "required|string|in:local",
		"labels":  "slice",
		"options": "slice",
	}
}

func (r *VolumeCreate) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *VolumeCreate) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *VolumeCreate) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
