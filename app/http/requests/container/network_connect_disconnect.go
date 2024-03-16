package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type NetworkConnectDisConnect struct {
	Network   string `form:"network" json:"network"`
	Container string `form:"container" json:"container"`
}

func (r *NetworkConnectDisConnect) Authorize(ctx http.Context) error {
	return nil
}

func (r *NetworkConnectDisConnect) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"network":   "required|string",
		"container": "required|string",
	}
}

func (r *NetworkConnectDisConnect) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *NetworkConnectDisConnect) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *NetworkConnectDisConnect) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
