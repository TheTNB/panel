package api

import (
	"fmt"
	"time"
)

type App struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Slug        string    `json:"slug"`
	Icon        string    `json:"icon"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Categories  []string  `json:"categories"`
	Requires    []string  `json:"requires"`
	Excludes    []string  `json:"excludes"`
	Versions    []struct {
		Url          string `json:"url"`
		Checksum     string `json:"checksum"`
		PanelVersion string `json:"panel_version"`
	} `json:"versions"`
	Order int `json:"order"`
}

type Apps []App

// GetApps 返回所有应用
func (r *API) GetApps() (*Apps, error) {
	resp, err := r.client.R().SetResult(&Response{}).Get("/apps")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get apps: %s", resp.String())
	}

	apps, err := getResponseData[Apps](resp)
	if err != nil {
		return nil, err
	}

	return apps, nil
}
