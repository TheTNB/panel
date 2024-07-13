package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"

	"github.com/TheTNB/panel/v2/pkg/types"
)

type NetworkCreate struct {
	Name    string                 `form:"name" json:"name"`
	Driver  string                 `form:"driver" json:"driver"`
	Ipv4    types.ContainerNetwork `form:"ipv4" json:"ipv4"`
	Ipv6    types.ContainerNetwork `form:"ipv6" json:"ipv6"`
	Labels  []types.KV             `form:"labels" json:"labels"`
	Options []types.KV             `form:"options" json:"options"`
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

func (r *NetworkCreate) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
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
