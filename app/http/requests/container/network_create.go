package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type NetworkCreate struct {
	Name    string   `form:"name" json:"name"`
	Driver  string   `form:"driver" json:"driver"`
	Ipv4    Network  `form:"ipv4" json:"ipv4"`
	Ipv6    Network  `form:"ipv6" json:"ipv6"`
	Labels  []string `form:"labels" json:"labels"`
	Options []string `form:"options" json:"options"`
}

func (r *NetworkCreate) Authorize(ctx http.Context) error {
	return nil
}

func (r *NetworkCreate) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":    "required|string",
		"driver":  "required|string|in:bridge,overlay,macvlan,ipvlan",
		"ipv4":    "required",
		"ipv6":    "required",
		"labels":  "slice",
		"options": "slice",
	}
}

func (r *NetworkCreate) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *NetworkCreate) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *NetworkCreate) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
