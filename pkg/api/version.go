package api

import (
	"fmt"
	"time"

	"github.com/TheTNB/panel/internal/panel"
)

type Version struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
}

type Versions []Version

// LatestVersion 返回最新版本
func (r *API) LatestVersion() (*Version, error) {
	resp, err := r.client.R().SetResult(&Response{}).Get("/versions/latest")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get latest version: %s", resp.String())
	}

	version, err := getResponseData[Version](resp)
	if err != nil {
		return nil, err
	}

	return version, nil
}

// IntermediateVersions 返回当前版本之后的所有版本
func (r *API) IntermediateVersions() (*Versions, error) {
	resp, err := r.client.R().
		SetQueryParam("start", panel.Version).
		SetResult(&Response{}).Get("/versions/intermediate")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get intermediate versions: %s", resp.String())
	}

	versions, err := getResponseData[Versions](resp)
	if err != nil {
		return nil, err
	}

	return versions, nil
}
